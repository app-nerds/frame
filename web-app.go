package frame

func (fa *FrameApplication) WebAppFolder(path string) *FrameApplication {
	fa.webAppFolder = path
	return fa
}

func (fa *FrameApplication) PrimaryLayoutName(name string) *FrameApplication {
	fa.primaryLayoutName = name
	return fa
}
