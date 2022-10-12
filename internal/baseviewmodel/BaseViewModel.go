package baseviewmodel

import webapp "github.com/app-nerds/frame/pkg/web-app"

type BaseViewModel struct {
	webapp.JavascriptIncludes
	AppName     string
	Stylesheets []string
}
