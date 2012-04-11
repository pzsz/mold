package wobject

type WModuleRenderer interface {
	WModule
	Render()
}

type WModuleRendererFactory interface {
	SetupManager(*WObjectManager)
	CreateWObjectRenderer(*WObject) WModuleRenderer
}
