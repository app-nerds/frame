package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
)

var (
	templateOptions = []string{"Base", "Admin Left Nav", "Admin Top Nav"}

	//go:embed templates/*
	templates embed.FS

	//go:embed database-migrations/*
	databaseMigrations embed.FS
)

type Context struct {
	AppName            string
	ModulePath         string
	TemplateSelection  string
	DockerRegistryPath string
}

type MappedFile struct {
	Source string
	Dest   string
}

func main() {
	var (
		err error
	)

	ctx := Context{}

	appNamePrompt := promptui.Prompt{
		Label: "App Name",
	}

	modulePathPrompt := promptui.Prompt{
		Label: "Module Path",
	}

	templatePrompt := promptui.Select{
		Label: "Select Template",
		Items: templateOptions,
	}

	ctx.AppName, _ = appNamePrompt.Run()
	ctx.ModulePath, _ = modulePathPrompt.Run()
	_, ctx.TemplateSelection, _ = templatePrompt.Run()
	ctx.DockerRegistryPath = getDockerRegistryPath(ctx.ModulePath)

	/*
	 * Create the app directory
	 */
	if err = os.MkdirAll(fmt.Sprintf("%s/frontend-templates", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/app/static/css", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/app/static/js", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/app/static/images", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/internal/handlers", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/internal/viewmodels", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/database-migrations", ctx.AppName), 0755); err != nil {
		panic(err)
	}

	/*
	 * Copy templates that don't need rendering
	 */
	staticFileMapping := []MappedFile{
		{Source: "templates/gitignore", Dest: fmt.Sprintf("%s/.gitignore", ctx.AppName)},
		{Source: "database-migrations/00000_init.down.sql", Dest: fmt.Sprintf("%s/database-migrations/00000_init.down.sql", ctx.AppName)},
		{Source: "database-migrations/00000_init.up.sql", Dest: fmt.Sprintf("%s/database-migrations/00000_init.up.sql", ctx.AppName)},
		{Source: "templates/jsconfig.json", Dest: fmt.Sprintf("%s/jsconfig.json", ctx.AppName)},
		{Source: "templates/base-layout", Dest: fmt.Sprintf("%s/frontend-templates/layout.tmpl", ctx.AppName)},
		{Source: "templates/base.min.css", Dest: fmt.Sprintf("%s/app/static/css/base.min.css", ctx.AppName)},
		{Source: "templates/components.min.css", Dest: fmt.Sprintf("%s/app/static/css/components.min.css", ctx.AppName)},
		{Source: "templates/frame.min.js", Dest: fmt.Sprintf("%s/app/static/js/frame.min.js", ctx.AppName)},
		{Source: "templates/icons.min.css", Dest: fmt.Sprintf("%s/app/static/css/icons.min.css", ctx.AppName)},
		{Source: "templates/index", Dest: fmt.Sprintf("%s/frontend-templates/index.tmpl", ctx.AppName)},
		{Source: "templates/VERSION", Dest: fmt.Sprintf("%s/VERSION", ctx.AppName)},
	}

	if ctx.TemplateSelection == "Base" {
		staticFileMapping = append(staticFileMapping, MappedFile{Source: "templates/base-layout", Dest: fmt.Sprintf("%s/frontend-templates/layout.tmpl", ctx.AppName)})
	}

	if ctx.TemplateSelection == "Admin Left Nav" {
		staticFileMapping = append(staticFileMapping, MappedFile{Source: "templates/admin-left-nav-layout", Dest: fmt.Sprintf("%s/frontend-templates/layout.tmpl", ctx.AppName)})
		staticFileMapping = append(staticFileMapping, MappedFile{Source: "templates/admin-left-side-nav.min.css", Dest: fmt.Sprintf("%s/app/static/css/admin-left-side-nav.min.css", ctx.AppName)})
	}

	if ctx.TemplateSelection == "Admin Top Nav" {
		staticFileMapping = append(staticFileMapping, MappedFile{Source: "templates/admin-top-nav-layout", Dest: fmt.Sprintf("%s/frontend-templates/layout.tmpl", ctx.AppName)})
		staticFileMapping = append(staticFileMapping, MappedFile{Source: "templates/admin-top-nav.min.css", Dest: fmt.Sprintf("%s/app/static/css/admin-top-nav.min.css", ctx.AppName)})
	}

	for _, staticFile := range staticFileMapping {
		if err = copyToFs(staticFile.Source, staticFile.Dest); err != nil {
			panic(err)
		}
	}

	/*
	 * Render Go templates
	 */
	goTemplates, err := template.ParseFS(templates, "templates/*.tmpl")
	renderTemplates := []MappedFile{
		{Source: ".env.tmpl", Dest: fmt.Sprintf("%s/.env", ctx.AppName)},
		{Source: ".env.template.tmpl", Dest: fmt.Sprintf("%s/.env.template", ctx.AppName)},
		{Source: "docker-compose.yml.tmpl", Dest: fmt.Sprintf("%s/docker-compose.yml", ctx.AppName)},
		{Source: "Dockerfile.tmpl", Dest: fmt.Sprintf("%s/Dockerfile", ctx.AppName)},
		{Source: "IndexHandler.go.tmpl", Dest: fmt.Sprintf("%s/internal/handlers/IndexHandler.go", ctx.AppName)},
		{Source: "IndexPageViewModel.go.tmpl", Dest: fmt.Sprintf("%s/internal/viewmodels/IndexPageViewModel.go", ctx.AppName)},
		{Source: "main.go.tmpl", Dest: fmt.Sprintf("%s/main.go", ctx.AppName)},
		{Source: "Makefile.tmpl", Dest: fmt.Sprintf("%s/Makefile", ctx.AppName)},
		{Source: "README.md.tmpl", Dest: fmt.Sprintf("%s/README.md", ctx.AppName)},
	}

	if err != nil {
		panic(err)
	}

	for _, renderTemplate := range renderTemplates {
		if err = renderToFs(goTemplates, renderTemplate.Source, renderTemplate.Dest, ctx); err != nil {
			panic(err)
		}
	}

	/*
	 * Now, switch to the new directory
	 */
	if err = os.Chdir(ctx.AppName); err != nil {
		panic(err)
	}

	/* Run go mod init */
	cmd := exec.Command("go", "mod", "init", ctx.ModulePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	/* Run go mod tidy */
	cmd = exec.Command("go", "mod", "tidy")

	if err = cmd.Run(); err != nil {
		panic(err)
	}
}

func copyToFs(src, dest string) error {
	var (
		err      error
		srcFile  fs.File
		destFile *os.File
	)

	if strings.HasPrefix(src, "templates") {
		if srcFile, err = templates.Open(src); err != nil {
			return err
		}
	}

	if strings.HasPrefix(src, "database-migrations") {
		if srcFile, err = databaseMigrations.Open(src); err != nil {
			return err
		}
	}

	defer srcFile.Close()

	if destFile, err = os.Create(dest); err != nil {
		return err
	}

	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

func renderToFs(parsedTemplates *template.Template, src, dest string, data any) error {
	var (
		err      error
		destFile *os.File
	)

	if destFile, err = os.Create(dest); err != nil {
		return err
	}

	defer destFile.Close()

	if err = parsedTemplates.ExecuteTemplate(destFile, src, data); err != nil {
		return err
	}

	return nil
}

func getDockerRegistryPath(modulePath string) string {
	split := strings.Split(modulePath, "/")

	if len(split) > 1 {
		suffix := strings.Join(split[1:], "/")
		return fmt.Sprintf("ghcr.io/%s", suffix)
	}

	return modulePath
}
