package frame

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterBaseAuthRoutes(router *mux.Router, webApp *WebApp) {
	router.HandleFunc(SiteAuthAccountPendingPath, handleAuthAccountPending(webApp)).Methods(http.MethodGet)

}
