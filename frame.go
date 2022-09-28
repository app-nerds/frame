package frame

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FrameApplication struct {
	appName       string
	externalAuths []goth.Provider
	router        *mux.Router
	sessionName   string
	sessionStore  sessions.Store
	templateFS    fs.FS
	templates     map[string]*template.Template
	version       string
	webAppFS      fs.FS

	// Paths
	accountAwaitingApprovalPath string
	unauthorizedPath            string
	unexpectedErrorPath         string
	webAppFolder                string

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
		appName: appName,
		Logger: logrus.New().WithFields(logrus.Fields{
			"who":     appName,
			"version": version,
		}),
		router:  mux.NewRouter(),
		version: version,
	}

	config := result.setupConfig()
	result.Logger.Logger.SetLevel(config.GetLogLevel())
	result.Config = config

	// Attach services
	result.MemberService = newMemberService(result)

	return result
}

func (fa *FrameApplication) Start() chan os.Signal {
	fa.Logger.WithFields(logrus.Fields{
		"host": fa.Config.ServerHost,
	}).Info("starting HTTP server...")

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
