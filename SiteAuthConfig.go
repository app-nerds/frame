package frame

type SiteAuthConfig struct {
	ContentTemplateName   string
	HtmlPaths             []string
	LayoutName            string
	PathsExcludedFromAuth []string
}
