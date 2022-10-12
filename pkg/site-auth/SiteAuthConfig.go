package siteauth

type SiteAuthConfig struct {
	ContentTemplateName   string
	HtmlPaths             []string
	LayoutName            string
	PathsExcludedFromAuth []string
}
