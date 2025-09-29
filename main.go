package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Kriptografi")

	tabs := container.NewAppTabs(
		container.NewTabItem("Caesar Cipher", createCaesarTab(w)),
		container.NewTabItem("Affine Cipher", createAlfineTab(w)),
		container.NewTabItem("RSA Cipher", createRsaTab(w)),
		container.NewTabItem("Advanced Encryption Standard", createAesTab(w)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.ShowAndRun()
}
