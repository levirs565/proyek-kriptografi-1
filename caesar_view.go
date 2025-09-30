package main

import (
	"errors"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var ErrInvalidKey = errors.New("kunci tidak valid")

func createCaesarTab(w fyne.Window) fyne.CanvasObject {
	keyEntry := widget.NewEntry()
	customCharsetEntry := widget.NewEntry()
	customCharsetEntry.SetPlaceHolder("Masukkan karakter custom...")
	customCharsetEntry.Hide()
	modeSelect := widget.NewSelect([]string{"Alfabet (A-Z)", "Alphanum (A-Z dan 0-9)", "ASCII", "Custom karakter"},
		func(value string) {
			if value == "Custom karakter" {
				customCharsetEntry.Show()
			} else {
				customCharsetEntry.Hide()
			}
		})
	modeMap := []CaesarMode{CaesarModeAlphabet, CaesarModeAlphanum, CaesarModeASCII, CaesarModeCustom}
	modeSelect.SetSelected("Alfabet (A-Z)")
	plainEntry := widget.NewMultiLineEntry()
	cipherEntry := widget.NewMultiLineEntry()
	processTitle := widget.NewLabel("Proses")
	processGrid := container.NewGridWrap(fyne.NewSize(30, 70))
	processPanel := container.NewBorder(processTitle, nil, nil, nil, container.NewVScroll(processGrid))
	mainContainer := NewFlexibleColumn(
		NewFlexibleItem(true, container.NewGridWithRows(2,
			container.NewBorder(widget.NewLabel("Plain Teks"), nil, nil, nil, plainEntry),
			container.NewBorder(widget.NewLabel("Cipher Teks"), nil, nil, nil, cipherEntry),
		)),
		NewFlexibleItem(true, processPanel),
	)

	processPanel.Hide()

	showProcessCheck := widget.NewCheck(
		"Show Process",
		func(b bool) {
			if b {
				processPanel.Show()
			} else {
				processPanel.Hide()
			}
			mainContainer.Refresh()
		},
	)

	showProcess := func(from, to []byte) {
		processGrid.RemoveAll()
		for i := range len(from) {
			content := container.NewBorder(
				widget.NewLabelWithStyle(string(from[i]), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(string(to[i]), fyne.TextAlignCenter, fyne.TextStyle{}),
				nil,
				nil,
				container.NewCenter(widget.NewIcon(theme.MoveDownIcon())),
			)
			background := canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground))
			processGrid.Add(container.NewStack(background, content))
		}
	}

	getKeyInt := func() (int, bool) {
		res, err := strconv.Atoi(keyEntry.Text)
		if err != nil {
			dialog.NewError(errors.Join(ErrInvalidKey, err), w).Show()
			return res, false
		}
		return res, true
	}

	middleLayout := container.NewVBox(
		widget.NewLabelWithStyle("Kunci", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyEntry,
		widget.NewLabelWithStyle("Opsi Caesar", fyne.TextAlignCenter, fyne.TextStyle{}), // <<<--- ditambahkan
		modeSelect,
		customCharsetEntry,
		widget.NewButton("Enkripsi", func() {
			key, success := getKeyInt()
			if !success {
				return
			}

			plain := []byte(plainEntry.Text)
			encrypted := caesarEncryptBytes(plain, key, modeMap[modeSelect.SelectedIndex()], customCharsetEntry.Text)
			processTitle.SetText("Proses Enkripsi")
			showProcess(plain, encrypted)
			cipherEntry.SetText(string(encrypted))
		}),
		widget.NewButton("Dekripsi", func() {
			key, success := getKeyInt()
			if !success {
				return
			}
			encrypted := []byte(cipherEntry.Text)
			plain := caesarDecryptBytes(encrypted, key, modeMap[modeSelect.SelectedIndex()], customCharsetEntry.Text)
			processTitle.SetText("Proses Dekripsi")
			showProcess(encrypted, plain)
			plainEntry.SetText(string(plain))
		}),
		showProcessCheck,
	)

	return container.NewPadded(container.NewBorder(
		widget.NewLabel("Caesar"), nil, middleLayout, nil,
		mainContainer,
	))
}
