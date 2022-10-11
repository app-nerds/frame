package frame

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/app-nerds/frame/internal/membermanagement"
	siteauth "github.com/app-nerds/frame/internal/site-auth"
	webapp "github.com/app-nerds/frame/internal/web-app"
	"github.com/app-nerds/frame/pkg/config"
	"github.com/app-nerds/frame/pkg/framemember"
	pkgsiteauth "github.com/app-nerds/frame/pkg/site-auth"
	pkgwebapp "github.com/app-nerds/frame/pkg/web-app"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed templates
var internalTemplatesFS embed.FS

//go:embed admin-templates
var adminTemplatesFS embed.FS

//go:embed admin-static
var adminStaticFS embed.FS

type FrameApplication struct {
	*sync.Mutex

	appName       string
	externalAuths []goth.Provider
	pageSize      int
	router        *mux.Router
	templateFS    fs.FS
	templates     map[string]*template.Template
	version       string

	// Template setup
	primaryLayoutName string

	// Internal services
	memberManagement *membermanagement.MemberManagement
	siteAuth         *siteauth.SiteAuth
	webApp           *webapp.WebApp

	// Public
	Config        *config.Config
	DB            *gorm.DB
	Logger        *logrus.Entry
	MemberService framemember.MemberService
	Server        *http.Server

	// Hooks
	OnAuthSuccess func(w http.ResponseWriter, r *http.Request, member framemember.Member)
}

func NewFrameApplication(appName, version string) *FrameApplication {
	result := &FrameApplication{
		Mutex: &sync.Mutex{},

		appName: appName,
		Logger: logrus.New().WithFields(logrus.Fields{
			"who":     appName,
			"version": version,
		}),
		pageSize: 25,
		router:   mux.NewRouter(),
		version:  version,
	}

	config := config.NewConfig(appName, version)
	result.Logger.Logger.SetLevel(config.GetLogLevel())
	result.Config = config

	// Attach Fireplace if configured
	result.withFireplace()

	return result
}

func (fa *FrameApplication) AddSiteAuth(config pkgsiteauth.SiteAuthConfig) *FrameApplication {
	if fa.webApp == nil {
		fa.Logger.Fatalf("please configure web application before site auth by calling AddWebApp()")
	}

	fa.siteAuth = siteauth.NewSiteAuth(siteauth.InternalSiteAuthConfig{
		Logger:       fa.Logger,
		SessionName:  fa.webApp.GetSessionName(),
		SessionStore: fa.webApp.GetSessionStore(),
	}, config)

	fa.memberManagement = membermanagement.NewMemberManagement(membermanagement.InternalMemberManagementConfig{
		AppName:       fa.appName,
		Logger:        fa.Logger,
		MemberService: &fa.MemberService,
		WebApp:        fa.webApp,
	})

	return fa
}

func (fa *FrameApplication) AddWebApp(config *pkgwebapp.WebAppConfig) *FrameApplication {
	fa.webApp = webapp.NewWebApp(
		webapp.InternalWebAppConfig{
			AdminTemplateFS:    adminTemplatesFS,
			AdminStaticFS:      adminStaticFS,
			AppName:            fa.appName,
			Debug:              fa.Config.Debug,
			Logger:             fa.Logger,
			FrameConfig:        fa.Config,
			InternalTemplateFS: internalTemplatesFS,
			MemberService:      &fa.MemberService,
			Version:            fa.version,
		},
		config,
	)

	return fa
}

func (fa *FrameApplication) Database(dst ...interface{}) *FrameApplication {
	var (
		err error
	)

	if fa.DB, err = gorm.Open(postgres.Open(fa.Config.DSN), &gorm.Config{}); err != nil {
		fa.Logger.WithError(err).Fatal("unable to connect to the database")
	}

	dst = append(dst, &framemember.MembersStatus{}, &framemember.Member{})
	_ = fa.DB.AutoMigrate(dst...)

	if err = fa.seedDataMemberStatuses(); err != nil {
		fa.Logger.WithError(err).Fatal("error seeding database...")
	}

	fa.setupServicesThatRequireDB()
	return fa
}

func (fa *FrameApplication) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	fa.webApp.RenderTemplate(w, name, data)
}

func (fa *FrameApplication) Start() chan os.Signal {
	if fa.webApp != nil {
		fa.webApp.RegisterRoutes(fa.router)
	}

	if fa.siteAuth != nil && fa.webApp != nil {
		fa.siteAuth.RegisterSiteAuthRoutes(fa.router, fa.webApp, &fa.MemberService)
		fa.memberManagement.RegisterRoutes(fa.router)

	}

	fa.Logger.WithFields(logrus.Fields{
		"host":     fa.Config.ServerHost,
		"debug":    fa.Config.Debug,
		"version":  fa.Config.Version,
		"loglevel": fa.Logger.Logger.Level,
	}).Info("starting HTTP server...")

	if fa.Config.Debug {
		fa.router.Use(requestLoggerMiddleware(fa.Logger))
	}

	fa.router.Use(accessControlMiddleware(AllowAllOrigins, AllowAllMethods, AllowAllHeaders))

	fa.Server = &http.Server{
		Addr:         fa.Config.ServerHost,
		WriteTimeout: time.Second * time.Duration(fa.Config.ServerWriteTimeout),
		ReadTimeout:  time.Second * time.Duration(fa.Config.ServerReadTimeout),
		IdleTimeout:  time.Second * time.Duration(fa.Config.ServerIdleTimeout),
		Handler:      fa.router,
	}

	go func() {
		err := fa.Server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			fa.Logger.WithError(err).Fatal("error starting server")
		}
	}()

	fa.Logger.Info("started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func (fa *FrameApplication) Stop() {
	var (
		err error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = fa.Server.Shutdown(ctx); err != nil {
		fa.Logger.WithError(err).Fatal("error shutting down server")
	}

	fa.Logger.Info("server stopped")
}

func (fa *FrameApplication) setupServicesThatRequireDB() {
	fa.MemberService = framemember.NewMemberService(framemember.MemberServiceConfig{
		DB:       fa.DB,
		PageSize: fa.Config.PageSize,
	})
}
