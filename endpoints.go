package frame

import (
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
Endpoint defines a single HTTP endpoint. Each endpoint is used
to configure a Gorilla Mux route.
*/
type Endpoint struct {
	Path        string
	Methods     []string
	HandlerFunc http.HandlerFunc
	Handler     http.Handler
}

/*
Endpoints represents an Endpoint slice.
*/
type Endpoints []Endpoint

func (a Endpoints) Len() int {
	return len(a)
}

func (a Endpoints) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Endpoints) Less(i, j int) bool {

	firstDynamic := isDynamic(a[i].Path)
	secondDynamic := isDynamic(a[j].Path)

	if firstDynamic && !secondDynamic {
		return false
	}

	if !firstDynamic && secondDynamic {
		return true
	}

	if len(a[i].Path) != len(a[j].Path) {
		return len(a[i].Path) > len(a[j].Path)
	}

	if a[i].Path == a[j].Path {
		panic("Two endpoints can't be same")
	}

	return true
}

func isDynamic(url string) bool {
	return strings.Contains(url, "{") && strings.Contains(url, "}")
}

func (fa *FrameApplication) SetupEndpoints(webAppFS fs.FS, endpoints Endpoints) *FrameApplication {
	fa.webAppFS = webAppFS
	fa.router.Use(accessControlMiddleware(AllowAllOrigins, AllowAllMethods, AllowAllHeaders))

	fs := http.FileServer(fa.getStaticFileSystem())
	adminFs := http.FileServer(fa.getAdminStaticFileSystem())

	sort.Sort(endpoints)

	for _, e := range endpoints {
		if fa.Config.Debug {
			fa.Logger.WithFields(logrus.Fields{
				"path":    e.Path,
				"methods": e.Methods,
			}).Info("registering endpoint")
		}

		if e.HandlerFunc != nil {
			fa.router.HandleFunc(e.Path, e.HandlerFunc).Methods(e.Methods...)
		} else {
			fa.router.Handle(e.Path, e.Handler).Methods(e.Methods...)
		}
	}

	if fa.Config.Debug {
		fa.Logger.Info("registering /static endpoint")
	}

	fa.router.HandleFunc("/errors/unexpected", fa.handleUnexpectedError)
	fa.router.PathPrefix("/static/").Handler(fs).Methods(http.MethodGet)
	fa.router.PathPrefix("/admin-static/").Handler(adminFs).Methods(http.MethodGet)

	return fa
}

func (fa *FrameApplication) getStaticFileSystem() http.FileSystem {
	if fa.Config.Version == "development" {
		if fa.Config.Debug {
			fa.Logger.Infof("serving static assets from filesystem out of '%s'", fa.webAppFolder)
		}

		return http.FS(os.DirFS(fa.webAppFolder))
	}

	if fa.Config.Debug {
		fa.Logger.Infof("serving static assets from embedded content at '%s'", fa.webAppFolder)
	}

	fsys, err := fs.Sub(fa.webAppFS, fa.webAppFolder)

	if err != nil {
		fa.Logger.WithError(err).Fatal("error loading static asset filesystem")
	}

	return http.FS(fsys)
}

func (fa *FrameApplication) getAdminStaticFileSystem() http.FileSystem {
	return http.FS(adminStaticFS)
}
