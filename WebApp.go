package frame

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/laher/mergefs"
	"github.com/sirupsen/logrus"
)

/*******************************************************************************
 * Internal Web App
 ******************************************************************************/

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
	memberManagement   *MemberManagement
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

type WebAppConfig struct {
	AppFolder         string
	AppFS             fs.FS
	PrimaryLayoutName string
	TemplateFS        fs.FS
	TemplateManifest  TemplateCollection
	SessionType       FrameSessionType
}

type AdminLoginData struct {
	BaseViewModel
	Message string
}

type Template struct {
	Name      string
	IsLayout  bool
	UseLayout string
}

type TemplateCollection []Template

type JavascriptInclude struct {
	Type string
	Src  string
}

type JavascriptIncludes []JavascriptInclude

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

/*******************************************************************************
 * Registration functions
 ******************************************************************************/

func (wa *WebApp) RegisterRoutes(router *mux.Router, adminRouter *mux.Router) {
	router.HandleFunc("/errors/unexpected", wa.handleUnexpectedError)

	adminRouter.HandleFunc("", wa.handleAdminDashboard)
	adminRouter.HandleFunc("/login", wa.handleAdminLogin)
}

func (wa *WebApp) registerAdminTemplates() TemplateCollection {
	manifest := TemplateCollection{}
	manifest = append(manifest, Template{Name: "admin-layout.tmpl", IsLayout: true})
	manifest = append(manifest, Template{Name: "admin-login.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-dashboard.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, wa.memberManagement.RegisterAdminTemplate()...)

	return manifest
}

func (wa *WebApp) registerInternalTemplates() TemplateCollection {
	wa.templateManifest = append(wa.templateManifest, Template{Name: "account-pending.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "login.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "unexpected-error.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, Template{Name: "sign-up.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, wa.memberManagement.RegisterTemplates()...)

	return wa.templateManifest
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

/*******************************************************************************
 * Admin routes
 ******************************************************************************/

func (wa *WebApp) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"AppName": wa.appName,
	}

	wa.RenderTemplate(w, "admin-dashboard.tmpl", data)
}

func (wa *WebApp) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session *sessions.Session
	)

	data := AdminLoginData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: []JavascriptInclude{},
			AppName:            wa.appName,
			Stylesheets:        []string{},
		},
		Message: "",
	}

	/*
	 * Handle login submission
	 */
	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		// First, is this a root user?
		if r.FormValue("userName") == wa.frameConfig.RootUserName && r.FormValue("password") == wa.frameConfig.RootUserPassword {
			if session, err = wa.adminSessionStore.Get(r, wa.adminSessionName); err != nil {
				wa.logger.WithError(err).Error("error geting session")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			session.Values["adminUserName"] = wa.frameConfig.RootUserName

			if err = wa.adminSessionStore.Save(r, w, session); err != nil {
				wa.logger.WithError(err).Error("error saving session")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			goTo := r.FormValue("referer")

			if goTo == "" {
				goTo = "/admin"
			}

			http.Redirect(w, r, goTo, http.StatusFound)
			return
		}

		// TODO: Add users from a database

		data.Message = "Invalid user name or password"
	}

	wa.RenderTemplate(w, "admin-login.tmpl", data)
}

/*
GET /errors/unexpected
*/
func (wa *WebApp) handleUnexpectedError(w http.ResponseWriter, r *http.Request) {
	wa.RenderTemplate(w, "unexpected-error.tmpl", nil)
}

/*******************************************************************************
 * Templates
 ******************************************************************************/

func (wa *WebApp) templateFuncIsSet(name string, data interface{}) bool {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	return v.FieldByName(name).IsValid()
}

func (wa *WebApp) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var (
		err  error
		tmpl *template.Template
		ok   bool
	)

	if tmpl, ok = wa.templates[name]; !ok {
		wa.logger.Fatalf("template '%s' not found!", name)
	}

	if err = tmpl.Execute(w, data); err != nil {
		wa.logger.WithError(err).Fatalf("error rendering '%s'", name)
	}
}

func (wa *WebApp) setupTemplateEngine() {
	var (
		err            error
		tmplDefinition Template
	)

	wa.templateFS = mergefs.Merge(wa.templateFS, wa.internalTemplateFS)
	wa.templateManifest = wa.registerInternalTemplates()

	templateFuncs := template.FuncMap{
		"IsSet": wa.templateFuncIsSet,
	}

	for _, tmplDefinition = range wa.templateManifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("frontend-templates", tmplDefinition.Name)
		layoutPath := filepath.Join("frontend-templates", tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.New(tmplDefinition.Name).Funcs(templateFuncs).ParseFS(wa.templateFS, tmplPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if tmplDefinition.UseLayout == "" {
				wa.logger.Fatalf("The template '%s' must have a layout specified.", tmplDefinition.Name)
			}

			if parsedTemplate, err = template.New(tmplDefinition.Name).Funcs(templateFuncs).ParseFS(wa.templateFS, tmplPath, layoutPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		wa.templates[tmplDefinition.Name] = parsedTemplate
	}

	/*
	 * Log out some debug information in development mode
	 */
	if wa.debug {
		for k := range wa.templates {
			wa.logger.WithFields(logrus.Fields{
				"templateName": k,
			}).Info("template captured")
		}
	}
}

func (wa *WebApp) setupAdminTemplates() {
	var (
		err            error
		tmplDefinition Template
	)

	manifest := wa.registerAdminTemplates()

	for _, tmplDefinition = range manifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("admin-templates", tmplDefinition.Name)
		layoutPath := filepath.Join("admin-templates", tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.ParseFS(wa.adminTemplateFS, tmplPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing admin layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if parsedTemplate, err = template.ParseFS(wa.adminTemplateFS, tmplPath, layoutPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing admin template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		wa.templates[tmplDefinition.Name] = parsedTemplate
	}
}
