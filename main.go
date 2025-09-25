package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Kriptografi")

	keyEntry := widget.NewEntry()
	plainEntry := widget.NewEntry()
	cipherEntry := widget.NewEntry()

	getKeyInt := func() (int, bool) {
		res, err := strconv.Atoi(keyEntry.Text)
		if err != nil {
			dialog.NewError(err, w).Show()
			return res, false
		}
		return res, true
	}

	middleLayout := container.NewVBox(
		widget.NewLabelWithStyle("Kunci", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyEntry,
		widget.NewButton("Enkripsi", func() {
			key, success := getKeyInt()
			if !success {
				return
			}
			cipherEntry.SetText(string(EncrpytBytes([]byte(plainEntry.Text), key)))
		}),
		widget.NewButton("Dekripsi", func() {
			key, success := getKeyInt()
			if !success {
				return
			}
			plainEntry.SetText(string(DecryptBytes([]byte(cipherEntry.Text), key)))
		}),
	)

	w.SetContent(
		container.NewBorder(
			widget.NewLabel("Caesar"), nil, middleLayout, nil,
			container.NewGridWithRows(2,
				container.NewBorder(widget.NewLabel("Plain Teks"), nil, nil, nil, plainEntry),
				container.NewBorder(widget.NewLabel("Cipher Teks"), nil, nil, nil, cipherEntry),
			)))
	w.ShowAndRun()
}
