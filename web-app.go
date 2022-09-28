package frame

func (fa *FrameApplication) WebAppFolder(path string) *FrameApplication {
	fa.webAppFolder = path
	return fa
}
