package frame

import (
	"io/fs"
)

type WebAppConfig struct {
	AppFolder         string
	AppFS             fs.FS
	PrimaryLayoutName string
	TemplateFS        fs.FS
	TemplateManifest  TemplateCollection
	SessionType       FrameSessionType
}
