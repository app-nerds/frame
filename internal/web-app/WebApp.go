package webapp

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/app-nerds/frame/internal/routepaths"
	"github.com/app-nerds/frame/pkg/config"
	"github.com/app-nerds/frame/pkg/framemember"
	"github.com/app-nerds/frame/pkg/framesessions"
	pkgwebapp "github.com/app-nerds/frame/pkg/web-app"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type InternalWebAppConfig struct {
	AppName            string
	AdminTemplateFS    fs.FS
	AdminStaticFS      fs.FS
	Debug              bool
	Logger             *logrus.Entry
	FrameConfig        *config.Config
	InternalTemplateFS fs.FS
	MemberService      *framemember.MemberService
	Version            string
}

type WebApp struct {
	adminStaticFS      fs.FS
	adminTemplateFS    fs.FS
	adminSessionName   string
	adminSessionStore  sessions.Store
	appName            string
	appFS              fs.FS
	appFolder          string
	debug              bool
	frameConfig        *config.Config
	internalTemplateFS fs.FS
	logger             *logrus.Entry
	memberService      *framemember.MemberService
	primaryLayoutName  string
	sessionName        string
	sessionStore       sessions.Store
	sessionType        framesessions.FrameSessionType
	templateFS         fs.FS
	templates          map[string]*template.Template
	templateManifest   pkgwebapp.TemplateCollection
	webAppConfig       *pkgwebapp.WebAppConfig
	version            string
}

func NewWebApp(internalConfig InternalWebAppConfig, webAppConfig *pkgwebapp.WebAppConfig) *WebApp {
	result := &WebApp{
		adminStaticFS:      internalConfig.AdminStaticFS,
		adminTemplateFS:    internalConfig.AdminTemplateFS,
		appName:            internalConfig.AppName,
		appFS:              webAppConfig.AppFS,
		appFolder:          webAppConfig.AppFolder,
		debug:              internalConfig.Debug,
		frameConfig:        internalConfig.FrameConfig,
		internalTemplateFS: internalConfig.InternalTemplateFS,
		logger:             internalConfig.Logger,
		memberService:      internalConfig.MemberService,
		primaryLayoutName:  webAppConfig.PrimaryLayoutName,
		sessionType:        webAppConfig.SessionType,
		templateFS:         webAppConfig.TemplateFS,
		templates:          map[string]*template.Template{},
		templateManifest:   webAppConfig.TemplateManifest,
		webAppConfig:       webAppConfig,
		version:            internalConfig.Version,
	}

	result.setupSessions()
	result.setupTemplateEngine()
	result.setupAdminTemplates()

	return result
}

func (wa *WebApp) GetAdminSessionName() string {
	return wa.adminSessionName
}

func (wa *WebApp) GetAdminSessionStore() sessions.Store {
	return wa.adminSessionStore
}

func (wa *WebApp) GetAppFolder() string {
	return wa.appFolder
}

func (wa *WebApp) GetSessionName() string {
	return wa.sessionName
}

func (wa *WebApp) GetSessionStore() sessions.Store {
	return wa.sessionStore
}

func (wa *WebApp) GetStaticFS() fs.FS {
	return wa.appFS
}

func (wa *WebApp) RegisterRoutes(router *mux.Router, adminRouter *mux.Router) {
	router.HandleFunc("/errors/unexpected", wa.handleUnexpectedError)

	adminRouter.HandleFunc("", wa.handleAdminDashboard)
	adminRouter.HandleFunc("/login", wa.handleAdminLogin)
}

/*
UnexpectedError redirects the user to a page for unexpected errors. This is configured
when calling AddWebApp
*/
func (wa *WebApp) UnexpectedError(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
}

func (wa *WebApp) setupSessions() {
	if wa.sessionType == framesessions.CookieSessionType {
		wa.sessionName, wa.sessionStore = framesessions.CookieSessions(wa.frameConfig)
		wa.adminSessionName, wa.adminSessionStore = framesessions.AdminCookieSessions(wa.frameConfig)
	}
}

func (wa *WebApp) registerAdminTemplates() pkgwebapp.TemplateCollection {
	manifest := pkgwebapp.TemplateCollection{}
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-layout.tmpl", IsLayout: true})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-login.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-dashboard.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-members-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-members-edit.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-roles-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})

	return manifest
}

func (wa *WebApp) registerInternalTemplates() pkgwebapp.TemplateCollection {
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "account-pending.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "login.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "unexpected-error.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "sign-up.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "member-profile.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "member-edit-avatar.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})

	return wa.templateManifest
}
