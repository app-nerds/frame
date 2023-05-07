package frame

import (
	"context"
	"crypto/tls"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/app-nerds/gobucket/v2/cmd/gobucketgo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

//go:embed frontend-templates
var internalTemplatesFS embed.FS

//go:embed admin-templates
var adminTemplatesFS embed.FS

//go:embed admin-static
var adminStaticFS embed.FS

//go:embed frame-static
var frameStaticFS embed.FS

type FrameApplication struct {
	*sync.Mutex

	appName       string
	cron          *cron.Cron
	externalAuths []goth.Provider
	hasEndpoints  bool
	pageSize      int
	router        *mux.Router
	templateFS    fs.FS
	templates     map[string]*template.Template
	version       string

	// Template setup
	primaryLayoutName string

	// Internal models
	customMemberSignupConfig *CustomMemberSignupConfig

	// Internal services
	gobucketClient   *gobucketgo.GoBucket
	memberManagement *MemberManagement
	siteAuth         *SiteAuth
	webApp           *WebApp

	// Public
	Config        *Config
	DB            *sql.DB
	Logger        *logrus.Entry
	EmailService  EmailServicer
	MemberService MemberService
	NsqPublisher  *nsq.Producer
	NsqConsumers  []*nsq.Consumer
	Server        *http.Server

	// Hooks
	OnAuthSuccess func(w http.ResponseWriter, r *http.Request, member Member)
}

/*
NewFrameApplication creates a new Frame application. This is the main entry point.
*/
func NewFrameApplication(appName, version string) *FrameApplication {
	result := &FrameApplication{
		Mutex: &sync.Mutex{},

		appName: appName,
		cron:    cron.New(),
		Logger: logrus.New().WithFields(logrus.Fields{
			"who":     appName,
			"version": version,
		}),
		pageSize: 25,
		router:   mux.NewRouter(),
		version:  version,
	}

	config := NewConfig(appName, version)
	result.Logger.Logger.SetLevel(config.GetLogLevel())
	result.Config = config

	if !config.Debug {
		result.Logger.Info("setting log format to JSON")
		result.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Attach Fireplace if configured
	result.withFireplace()
	result.withGobucket()

	return result
}

func (fa *FrameApplication) AddSiteAuth(config SiteAuthConfig) *FrameApplication {
	if fa.webApp == nil {
		fa.Logger.Fatalf("please configure web application before site auth by calling AddWebApp()")
	}

	fa.siteAuth = NewSiteAuth(InternalSiteAuthConfig{
		FrameStaticFS: frameStaticFS,
		Logger:        fa.Logger,
		SessionName:   fa.webApp.GetSessionName(),
		SessionStore:  fa.webApp.GetSessionStore(),
	}, config)

	fa.memberManagement = NewMemberManagement(InternalMemberManagementConfig{
		AppName:        fa.appName,
		GobucketClient: fa.gobucketClient,
		Logger:         fa.Logger,
		MemberService:  &fa.MemberService,
		WebApp:         fa.webApp,
	})

	fa.webApp.memberManagement = fa.memberManagement

	return fa
}

/*
AddWebApp configures this Frame application to include a web application. Web Applications
in Frame are simple Go Template-based applications. Templates are embedded into the
final binary.
*/
func (fa *FrameApplication) AddWebApp(config *WebAppConfig) *FrameApplication {
	fa.webApp = NewWebApp(
		InternalWebAppConfig{
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

func (fa *FrameApplication) AddCron(schedule string, cronFunc func(app *FrameApplication)) *FrameApplication {
	fa.cron.AddFunc(schedule, func() {
		cronFunc(fa)
	})

	return fa
}

func (fa *FrameApplication) AddEmailService() *FrameApplication {
	fa.Logger.Info("setting up email service...")
	fa.EmailService = NewEmailService(emailServiceConfig{
		ApiKey: fa.Config.MailApiKey,
	})

	return fa
}

func (fa *FrameApplication) AddNsqConsumer(topic, channel string, handler nsq.Handler) *FrameApplication {
	var (
		err      error
		consumer *nsq.Consumer
	)

	nsqConfig := nsq.NewConfig()

	if consumer, err = nsq.NewConsumer(topic, channel, nsqConfig); err != nil {
		fa.Logger.WithError(err).Fatal("unable to create NSQ consumer")
	}

	consumer.AddHandler(handler)

	if err = consumer.ConnectToNSQLookupd(fa.Config.NsqLookupd); err != nil {
		fa.Logger.WithError(err).WithField("address", fa.Config.NsqLookupd).Fatal("error connecting to nsqlookupd")
	}

	consumer.SetLoggerLevel(nsq.LogLevelError)
	fa.NsqConsumers = append(fa.NsqConsumers, consumer)
	return fa
}

func (fa *FrameApplication) AddNsqPublisher() *FrameApplication {
	var (
		err      error
		producer *nsq.Producer
	)

	fa.Logger.Info("connecting to NSQ...")
	nsqConfig := nsq.NewConfig()

	producer, err = nsq.NewProducer(fa.Config.Nsqd, nsqConfig)

	if err != nil {
		fa.Logger.WithError(err).WithField("address", fa.Config.Nsqd).Fatal("error connecting to nsqd")
	}

	producer.SetLoggerLevel(nsq.LogLevelError)
	fa.NsqPublisher = producer
	return fa
}

func (fa *FrameApplication) Database(migrationDirectory string) *FrameApplication {
	var (
		err error
	)

	/*
	 * Connect to the database
	 */
	if fa.DB, err = sql.Open("postgres", fa.Config.DSN); err != nil {
		fa.Logger.WithError(err).Fatal("error connecting to database")
	}

	fa.DB.SetConnMaxIdleTime(time.Minute * 5)
	fa.DB.SetMaxOpenConns(100)

	if migrationDirectory != "" {
		d, _ := os.Getwd()
		finalPath := filepath.Join(d, migrationDirectory)
		fa.Logger.Infof("auto-migrating database using directory '%s'...", finalPath)

		driver, err := postgres.WithInstance(fa.DB, &postgres.Config{})

		if err != nil {
			fa.Logger.WithError(err).Fatal("error creating database driver")
		}

		m, err := migrate.NewWithDatabaseInstance("file://"+finalPath, "postgres", driver)

		if err != nil {
			fa.Logger.WithError(err).Fatal("error creating migration instance")
		}

		err = m.Up()

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fa.Logger.WithError(err).Fatal("error running migrations")
		}
	}

	if fa.webApp != nil {
		fa.setupServicesThatRequireDB()
	}

	return fa
}

func (fa *FrameApplication) GetMemberSession(r *http.Request) Member {
	firstName := r.Context().Value("firstName").(string)
	lastName := r.Context().Value("lastName").(string)
	email := r.Context().Value("email").(string)
	avatarURL := r.Context().Value("avatarURL").(string)
	memberID := r.Context().Value("memberID").(string)
	status := r.Context().Value("status").(string)

	return Member{
		ID:        memberID,
		AvatarURL: avatarURL,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Status: MembersStatus{
			Status: MemberStatus(status),
		},
	}
}

func (fa *FrameApplication) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	fa.webApp.RenderTemplate(w, name, data)
}

func (fa *FrameApplication) SetLogLevel(level logrus.Level) {
	fa.Logger.Logger.SetLevel(level)
}

func (fa *FrameApplication) SetLogLevelString(level string) {
	parsedLevel, err := logrus.ParseLevel(level)

	if err != nil {
		fa.Logger.WithError(err).Fatal("unable to parse log level")
	}

	fa.Logger.Logger.SetLevel(parsedLevel)
}

func (fa *FrameApplication) Start() chan os.Signal {
	var adminRouter *mux.Router

	/*
	 * Start CRON jobs
	 */
	if len(fa.cron.Entries()) > 0 {
		fa.Logger.Infof("starting %d cron jobs...", len(fa.cron.Entries()))
		fa.cron.Start()
	}

	/*
	 * If we have a web app register the admin routes
	 */
	if fa.webApp != nil {
		adminRouter = fa.router.PathPrefix("/admin").Subrouter()
		adminRouter.Use(adminAuthMiddleware(fa.Logger, fa.Config, fa.webApp.GetAdminSessionStore()))

		fa.webApp.RegisterRoutes(fa.router, adminRouter)
	}

	if fa.siteAuth != nil && fa.webApp != nil {
		fa.siteAuth.RegisterSiteAuthRoutes(fa.router, fa.webApp, &fa.MemberService)
		fa.siteAuth.RegisterStaticFrameAssetsRoute(fa.router)
		fa.memberManagement.RegisterRoutes(fa.router, adminRouter)
	}

	/*
	 * If we have either endpoints or a web app start the HTTP server
	 */
	if fa.hasEndpoints || fa.webApp != nil {
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
			Handler:      handlers.CompressHandler(fa.router),
		}

		go func() {
			var (
				err error
			)

			if fa.Config.AutoSSLEmail != "" && fa.Config.AutoSSLWhitelist != "" {
				autocertManager := &autocert.Manager{
					Prompt: autocert.AcceptTOS,
					Cache:  autocert.DirCache("./certs"),
					Email:  fa.Config.AutoSSLEmail,
					HostPolicy: func(ctx context.Context, host string) error {
						domains := strings.Split(fa.Config.AutoSSLWhitelist, ",")

						for _, domain := range domains {
							if host == domain {
								return nil
							}
						}

						return fmt.Errorf("acme/autocert: %s host not allowed", host)
					},
				}

				fa.Server.TLSConfig = &tls.Config{
					GetCertificate: autocertManager.GetCertificate,
				}

				go func() {
					autocertMux := &http.ServeMux{}
					autocertServer := &http.Server{
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 5 * time.Second,
						IdleTimeout:  120 * time.Second,
						Handler:      autocertMux,
						Addr:         ":80",
					}

					autocertServer.Handler = autocertManager.HTTPHandler(autocertServer.Handler)
					_ = autocertServer.ListenAndServe()
				}()

				err = fa.Server.ListenAndServeTLS("", "")
			} else {
				err = fa.Server.ListenAndServe()
			}

			if err != nil && err != http.ErrServerClosed {
				fa.Logger.WithError(err).Fatal("error starting server")
			}
		}()
	}

	fa.Logger.Info("started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func (fa *FrameApplication) Stop() {
	var (
		err         error
		cronContext context.Context
	)

	/*
	 * Stop any running CRON jobs
	 */
	if len(fa.cron.Entries()) > 0 {
		fa.Logger.Infof("stopping %d cron jobs...", len(fa.cron.Entries()))
		cronContext = fa.cron.Stop()

		<-cronContext.Done()
		fa.Logger.Info("cron jobs stopped.")
	}

	/*
	 * Stop any NSQ consumers
	 */
	if len(fa.NsqConsumers) > 0 {
		fa.Logger.Infof("stopping %d nsq consumers...", len(fa.NsqConsumers))

		for _, consumer := range fa.NsqConsumers {
			consumer.Stop()
			<-consumer.StopChan
		}

		fa.Logger.Info("nsq consumers stopped.")
	}

	/*
	 * If we have a web server stop it
	 */
	if fa.webApp != nil || fa.hasEndpoints {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = fa.Server.Shutdown(ctx); err != nil {
			fa.Logger.WithError(err).Fatal("error shutting down server")
		}
	}

	fa.Logger.Info("server stopped.")
}

// func (fa *FrameApplication) WithCustomSignUpForm(config *CustomMemberSignupConfig) *FrameApplication {
// 	fa.customMemberSignupConfig = config
// 	return fa
// }

func (fa *FrameApplication) setupServicesThatRequireDB() {
	fa.MemberService = NewMemberService(MemberServiceConfig{
		DB:       fa.DB,
		PageSize: fa.Config.PageSize,
	})
}
