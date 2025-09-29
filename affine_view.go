package main

import (
	"errors"
	"fmt"
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
	modeSelect := widget.NewSelect([]string{"Alfabet (A-Z)", "ASCII"}, nil)
	modeMap := []AffineMode{AffineModeAlphabet, AffineModeASCII}
	modeSelect.SetSelectedIndex(0)
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
			return res1, res2, false
		}
		if err2 != nil {
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
		widget.NewButton("Enkripsi", func() {
			keyA, keyB, success := getKeyInt()
			mode := modeMap[modeSelect.SelectedIndex()]
			if !success {
				dialog.NewError(errors.New("kunci a dan kunci b harus berupa angka"), w).Show()
				return
			}
			if !affineIsCoprime(keyA, mode) {
				dialog.NewError(fmt.Errorf("kunci a tidak koprima dengan %d", affineGetModulo(mode)), w).Show()
				return
			}

			plain := []byte(plainEntry.Text)
			encrypted := affineEncryptBytes(plain, keyA, keyB, mode)

			processTitle.SetText("Proses Enkripsi")
			showProcess(plain, encrypted)
			cipherEntry.SetText(string(encrypted))
		}),
		widget.NewButton("Dekripsi", func() {
			keyA, keyB, success := getKeyInt()
			mode := modeMap[modeSelect.SelectedIndex()]

			if !success {
				return
			}

			encrypted := []byte(cipherEntry.Text)
			plain := affineDecryptBytes(encrypted, keyA, keyB, mode)
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
