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
