package frame

import fireplacehook "github.com/app-nerds/fireplace/v2/cmd/fireplace-hook"

func (fa *FrameApplication) withFireplace() {
	if fa.Config.FireplaceURL != "" {
		fa.Logger.Logger.AddHook(fireplacehook.NewFireplaceHook(&fireplacehook.FireplaceHookConfig{
			Application:  fa.appName,
			FireplaceURL: fa.Config.FireplaceURL,
			Password:     fa.Config.FireplacePassword,
		}))
	}
}
