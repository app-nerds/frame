package frame

import "github.com/app-nerds/gobucket/v2/cmd/gobucketgo"

func (fa *FrameApplication) withGobucket() {
	if fa.Config.GobucketURL != "" {
		fa.gobucketClient = gobucketgo.New(gobucketgo.Config{
			AppKey:     fa.Config.GobucketAppKey,
			BaseURL:    fa.Config.GobucketURL,
			ClientCode: fa.Config.GobucketClientCode,
		})
	}
}
