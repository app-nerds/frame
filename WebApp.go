package frame

import (
	"html/template"
	"io/fs"
	"net/http"

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
	FrameConfig        *Config
	InternalTemplateFS fs.FS
	MemberService      *MemberService
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
	frameConfig        *Config
	internalTemplateFS fs.FS
	logger             *logrus.Entry
	memberService      *MemberService
	primaryLayoutName  string
	sessionName        string
	sessionStore       sessions.Store
	sessionType        FrameSessionType
	templateFS         fs.FS
	templates          map[string]*template.Template
	templateManifest   TemplateCollection
	webAppConfig       *WebAppConfig
	version            string
}

func NewWebApp(internalConfig InternalWebAppConfig, webAppConfig *WebAppConfig) *WebApp {
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
	http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
}

func (wa *WebApp) setupSessions() {
	if wa.sessionType == CookieSessionType {
		wa.sessionName, wa.sessionStore = CookieSessions(wa.frameConfig)
		wa.adminSessionName, wa.adminSessionStore = AdminCookieSessions(wa.frameConfig)
	}
}

func (wa *WebApp) registerAdminTemplates() TemplateCollection {
	manifest := TemplateCollection{}
	manifest = append(manifest, Template{Name: "admin-layout.tmpl", IsLayout: true})
	manifest = append(manifest, Template{Name: "admin-login.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-dashboard.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-members-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-members-edit.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-roles-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-roles-create.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-roles-edit.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})

	return manifest
}

func (wa *WebApp) registerInternalTemplates() TemplateCollection {
	wa.templateManifest = append(wa.templateManifest, Template{Name: "account-pending.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "login.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "unexpected-error.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "sign-up.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "member-profile.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "member-edit-avatar.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})

	return wa.templateManifest
}
