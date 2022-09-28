package frame

import "github.com/gorilla/sessions"

func (fa *FrameApplication) CookieSessions(sessionName string, maxAge int) *FrameApplication {
	sessionStorage := sessions.NewCookieStore([]byte(fa.Config.SessionKey))
	sessionStorage.MaxAge(maxAge)
	sessionStorage.Options.Path = "/"
	sessionStorage.Options.HttpOnly = true
	sessionStorage.Options.Secure = fa.Config.Debug

	fa.sessionName = sessionName
	fa.sessionStore = sessionStorage
	return fa
}
