package main
import (
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
	processTitle := widget.NewLabel("Proses")
	processGrid := container.NewGridWrap(fyne.NewSize(30, 70))
	processPanel := container.NewBorder(processTitle, nil, nil, nil, container.NewVScroll(processGrid))
	mainContainer := container.New(
		NewFlexibleLayout(),
		container.NewGridWithRows(2,
			container.NewBorder(widget.NewLabel("Plain Teks"), nil, nil, nil, plainEntry),
			container.NewBorder(widget.NewLabel("Cipher Teks"), nil, nil, nil, cipherEntry),
		),
		processPanel,
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
		res1, err := strconv.Atoi(keyEntryA.Text)
		res2, err := strconv.Atoi(keyEntryB.Text)
		if err != nil {
			dialog.NewError(err, w).Show()
			return res1, res2, false
		}
		return res1, res2, true
	}

	middleLayout := container.NewVBox(
		widget.NewLabelWithStyle("Kunci A", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyEntryA,
		widget.NewLabelWithStyle("Kunci B", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyEntryB,
		widget.NewButton("Enkripsi", func() {
			keyA, keyB, success := getKeyInt()
			if !success {
				// dialog.show("Kunci A harus berupa angka", w)
				return
			}
			isKoprima(keyA)
			plain := []byte(plainEntry.Text)
			encrypted := affineEncryptBytes(plain, keyA, keyB)
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
