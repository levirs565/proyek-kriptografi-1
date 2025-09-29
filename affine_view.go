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

func createAlfineTab(w fyne.Window) fyne.CanvasObject {
	keyEntryA := widget.NewEntry()
	keyEntryB := widget.NewEntry()

	plainEntry := widget.NewMultiLineEntry()
	cipherEntry := widget.NewMultiLineEntry()
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
	modeSelect.SetSelected("Alfabet (A-Z)")
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

	getKeyInt := func() (int, int, bool) {
		res1, err1 := strconv.Atoi(keyEntryA.Text)
		res2, err2 := strconv.Atoi(keyEntryB.Text)
		if err1 != nil {
			dialog.NewError(err1, w).Show()
			return res1, res2, false
		}
		if err2 != nil {
			dialog.NewError(err2, w).Show()
			return res1, res2, false
		}

		return res1, res2, true
	}

	middleLayout := container.NewVBox(
		widget.NewLabelWithStyle("Kunci A", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyEntryA,
		widget.NewLabelWithStyle("Kunci B", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyEntryB,
		widget.NewLabelWithStyle("Opsi Affine", fyne.TextAlignCenter, fyne.TextStyle{}), // <<<--- ditambahkan
		modeSelect,
		customCharsetEntry,
		widget.NewButton("Enkripsi", func() {
			keyA, keyB, success := getKeyInt()
			if !success {
				dialog.NewError(errors.New("kunci a dan kunci b harus berupa angka"), w).Show()
				return
			}
			if !isKoprima(keyA, modeSelect.Selected) {
				dialog.NewError(errors.New("kunci a tidak koprima dengan 255"), w).Show()
			}

			plain := []byte(plainEntry.Text)
			encrypted, err := affineEncryptBytes(plain, keyA, keyB, modeSelect.Selected)

			if err != nil {
				dialog.NewError(errors.Join(err), w).Show()
			}
			processTitle.SetText("Proses Enkripsi")
			showProcess(plain, encrypted)
			cipherEntry.SetText(string(encrypted))
		}),
		widget.NewButton("Dekripsi", func() {
			keyA, keyB, success := getKeyInt()
			if !success {
				return
			}

			encrypted := []byte(cipherEntry.Text)
			plain := affineDecryptBytes(encrypted, keyA, keyB)
			processTitle.SetText("Proses Dekripsi")
			showProcess(encrypted, plain)
			plainEntry.SetText(string(plain))
		}),
		showProcessCheck,
	)

	return container.NewBorder(
		widget.NewLabel("Affine"), nil, middleLayout, nil,
		mainContainer,
	)
}
