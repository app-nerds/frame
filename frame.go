package frame

import (
	"html/template"
	"io/fs"

	"github.com/sirupsen/logrus"
)

type FrameApplication struct {
	appName    string
	logger     *logrus.Entry
	templateFS fs.FS
	templates  map[string]*template.Template
	version    string
}

func NewFrameApplication(appName, version string) *FrameApplication {
	result := &FrameApplication{
		appName: appName,
		logger: logrus.New().WithFields(logrus.Fields{
			"who":     appName,
			"version": version,
		}),
		version: version,
	}

	return result
}
