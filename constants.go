package frame

const (
	// Error paths
	UnexpectedErrorPath string = "/errors/unexpected"

	// Site Auth endpoints
	SiteAuthLoginPath          string = "/member/login"
	SiteAuthLogoutPath         string = "/member/logout"
	SiteAuthAccountPendingPath string = "/member/account-pending"

	// Member API endpoints
	MemberApiCurrentMember string = "/api/member/current"
	MemberApiLogOut        string = "/api/member/logout"
)
