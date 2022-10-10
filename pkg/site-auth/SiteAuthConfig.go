package siteauth

type SiteAuthConfig struct {
	BaseData              map[string]interface{}
	ContentTemplateName   string
	HtmlPaths             []string
	LayoutName            string
	PathsExcludedFromAuth []string
}
