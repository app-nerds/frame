package siteauth

import (
	"net/http"

	"github.com/app-nerds/frame/internal/routepaths"
	webapp "github.com/app-nerds/frame/internal/web-app"
	"github.com/gorilla/mux"
)

func RegisterBaseAuthRoutes(router *mux.Router, webApp *webapp.WebApp) {
	router.HandleFunc(routepaths.SiteAuthAccountPendingPath, handleAuthAccountPending(webApp)).Methods(http.MethodGet)

}
