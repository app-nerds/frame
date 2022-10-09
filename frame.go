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

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/sirupsen/logrus"
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
	sessionName   string
	sessionStore  sessions.Store
	templateFS    fs.FS
	templates     map[string]*template.Template
	version       string
	webAppFS      fs.FS

	// Paths
	webAppFolder string

	// Template setup
	primaryLayoutName string

	// Public
	Config        *Config
	DB            *gorm.DB
	Logger        *logrus.Entry
	MemberService MemberService
	Server        *http.Server

	// Hooks
	OnAuthSuccess func(w http.ResponseWriter, r *http.Request, member Member)
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

	config := result.setupConfig()
	result.Logger.Logger.SetLevel(config.GetLogLevel())
	result.Config = config

	// Attach services
	result.MemberService = newMemberService(result)

	return result
}

func (fa *FrameApplication) Start() chan os.Signal {
	fa.setupAdminRoutes()

	fa.Logger.WithFields(logrus.Fields{
		"host":    fa.Config.ServerHost,
		"debug":   fa.Config.Debug,
		"version": fa.Config.Version,
	}).Info("starting HTTP server...")

	if fa.Config.Debug {
		fa.router.Use(requestLoggerMiddleware(fa.Logger))
	}

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
