# Frame
Frame is a web application framework written for Go, HTML, JavaScript, and CSS. It provides tools for building service applications, CRON applications, and web applications. It is very opinionated and comes with tons of built-in functionality for things you'd need in many web application, such as:

* Member/User/Role management
* Authentication management
* Basic CSS components built from scratch

## 🚀 Getting Started

To use Frame add it to your Go project by running:

```bash
go get -u github.com/app-nerds/frame github.com/app-nerds/kit/v6 github.com/app-nerds/fireplace/v2 github.com/app-nerds/gobucket/v2
```

The most basic Frame application starts with initializing the framework.

```go
package main

import (
  "github.com/app-nerds/frame"
)

func main() {
	app := frame.NewFrameApplication("cronapp", "1.0.0")

	<-app.Start()
	app.Stop()
}
```

This is the most bare-bones application. It doesn't do anything. Once the app is created you can add on additional functionality. Here are a few small examples:

### CRON Application

```go
package main

import (
  "github.com/app-nerds/frame"
)

func main() {
	app := frame.NewFrameApplication("cronapp", "1.0.0").
		AddCron("*/1 * * * *", func(app *frame.FrameApplication) {
			app.Logger.Info("Cron job running")
		})

	<-app.Start()
	app.Stop()
}
```

### Service Application

```go
package main

import (
	"net/http"

	"github.com/app-nerds/frame"
)

func main() {
	app := frame.NewFrameApplication("serviceapp", "1.0.0")

	app = app.SetupEndpoints(frame.Endpoints{
		{Path: "/hello", Methods: []string{http.MethodGet}, Handler: helloHandler(app)},
	})

	<-app.Start()
	app.Stop()
}

func helloHandler(app *frame.FrameApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("Hello world!")
		app.WriteString(w, http.StatusOK, "Hello world!")
	}
}
```

### Web Application

```go
package main

import (
	"embed"
	"net/http"

	"github.com/app-nerds/frame"
)

var (
	//go:embed app
	appFs embed.FS

	//go:embed frontend-templates
	templateFs embed.FS
)

func main() {
	app := frame.NewFrameApplication("webapp", "1.0.0").
		AddWebApp(&frame.WebAppConfig{
			AppFolder:         "app",
			AppFS:             appFs,
			PrimaryLayoutName: "layout",
			TemplateFS:        templateFs,
			TemplateManifest: []frame.Template{
				{Name: "layout.tmpl", IsLayout: true, UseLayout: ""},
				{Name: "home.tmpl", IsLayout: false, UseLayout: "layout.tmpl"},
			},
		})

	app = app.SetupEndpoints(frame.Endpoints{
		{Path: "/", Methods: []string{http.MethodGet}, Handler: homeHandler(app)},
	})

	<-app.Start()
	app.Stop()
}

func homeHandler(app *frame.FrameApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.RenderTemplate(w, "home.tmpl", nil)
	}
}
```

## Frame Application

The core of this framework is the Frame Application. This is what houses the bulk of the functionality. To create an application you always start with the following.

```go
app := frame.NewFrameApplication("appName", "version")
```

From here you can add on more features based on your needs. Then, when ready to run it:

```go
<-app.Start()
app.Stop()
```

### Adding a Web Application
Frame supports embedding Go templates to build web applications. To do this you call `AddWebApp()`. Before that, you'll need to add variables to embed two directories. The first directory is for your Go templates. The second will hold any static assets.

```go
var (
	//go:embed app
	appFS embed.FS

	//go:embed frontend-templates/*
	templateFS embed.FS
)

func main() {
	app := frame.NewFrameApplication("webapp", "1.0.0").
		AddWebApp(&frame.WebAppConfig{
			AppFolder:         "app",
			AppFS:             appFs,
			PrimaryLayoutName: "layout",
			TemplateFS:        templateFs,
			TemplateManifest: []frame.Template{
				{Name: "layout.tmpl", IsLayout: true},
				{Name: "home.tmpl", UseLayout: "layout.tmpl"},
			},
		})

	app = app.SetupEndpoints(frame.Endpoints{
		{Path: "/", Methods: []string{http.MethodGet}, Handler: homeHandler(app)},
	})

	<-app.Start()
	app.Stop()
}

func homeHandler(app *frame.FrameApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.RenderTemplate(w, "home.tmpl", nil)
	}
}
```

Here's the breakdown:

* **AppFolder** provides the name of the folder that contains your static app assets.
* **AppFS** is the variable that contains the embedded static app assets.
* **PrimaryLayoutName** is the name (minus extension) of the main Go template used as a layout.
* **TemplateFS** is the variable that contains the embedded Go templates.
* **TemplateManifest** is a slice of template definitions. Here you need to define your templates and layouts.

After adding the web application you can set up endpoints and their handlers. Handles are just standard HTTP `HandlerFunc` from the Go standard library. The definition needs that path, accepted methods, a handler, and an optional middleware. The path can accept variables in the form of `{varname}`.

## Adding a Database

Frame supports a Postgres database. It also supports using database migration scripts.
