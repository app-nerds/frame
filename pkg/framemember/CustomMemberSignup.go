package framemember

import "net/http"

type CustomMemberSignupConfig struct {
	Handler      http.HandlerFunc
	LayoutName   string
	TemplateName string
}
