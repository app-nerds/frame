package framesessions

import (
	"github.com/app-nerds/frame/pkg/config"
	"github.com/gorilla/sessions"
)

func CookieSessions(config *config.Config) (string, sessions.Store) {
	sessionStorage := sessions.NewCookieStore([]byte(config.SessionKey))
	sessionStorage.MaxAge(config.SessionMaxAge)
	sessionStorage.Options.Path = "/"
	sessionStorage.Options.HttpOnly = true
	sessionStorage.Options.Secure = config.Debug

	return config.SessionName, sessionStorage
}

func AdminCookieSessions(config *config.Config) (string, sessions.Store) {
	sessionStorage := sessions.NewCookieStore([]byte(config.AdminSessionKey))
	sessionStorage.MaxAge(config.AdminSessionMaxAge)
	sessionStorage.Options.Path = "/"
	sessionStorage.Options.HttpOnly = true
	sessionStorage.Options.Secure = config.Debug

	return config.AdminSessionName, sessionStorage

}
