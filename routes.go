package frame

func (fa *FrameApplication) AccountAwaitingApprovalPath(path string) *FrameApplication {
	fa.accountAwaitingApprovalPath = path
	return fa
}

func (fa *FrameApplication) UnauthorizedPath(path string) *FrameApplication {
	fa.unauthorizedPath = path
	return fa
}

func (fa *FrameApplication) UnexpectedErrorPath(path string) *FrameApplication {
	fa.unexpectedErrorPath = path
	return fa
}
