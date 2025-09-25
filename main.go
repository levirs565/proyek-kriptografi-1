package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Kriptografi")

	w.SetContent(widget.NewLabel("Halo dunia"))
	w.ShowAndRun()
}
