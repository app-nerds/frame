package webapp

import (
	"io/fs"

	"github.com/app-nerds/frame/pkg/framesessions"
)

type WebAppConfig struct {
	AppFolder         string
	AppFS             fs.FS
	PrimaryLayoutName string
	TemplateFS        fs.FS
	TemplateManifest  TemplateCollection
	SessionType       framesessions.FrameSessionType
}
