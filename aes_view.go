package main

import (
	"errors"
	"kriptografi1/aes"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var ErrKunciHex = errors.New("gagal mendecode Hex kunci")
var ErrPlainEncode = errors.New("gagal encode hex plain")
var ErrPlainDecode = errors.New("gagal decode hex plain")
var ErrCipherEncode = errors.New("gagal encode hex cipher")
var ErrCipherDecode = errors.New("gagal decode hex cipher")
var ErrDecrypt = errors.New("gagal decrypt")

var inputTypes = []string{"Raw", "Hex"}

func createAesTab(w fyne.Window) fyne.CanvasObject {
	aes.AESInit()

	variantSelect := widget.NewSelect([]string{"AES-128", "AES-192", "AES-256"}, func(s string) {})
	variantMap := []aes.AESVariant{aes.AES128, aes.AES192, aes.AES256}

	keyTypeSelect := widget.NewSelect(inputTypes, func(s string) {})
	keyEntry := widget.NewEntry()
	plainEntry := widget.NewMultiLineEntry()
	plainTypeSelect := widget.NewSelect(inputTypes, func(s string) {})
	cipherEntry := widget.NewMultiLineEntry()

	mainContainer := container.NewGridWithRows(2,
		container.NewBorder(
			container.NewVBox(widget.NewLabel("Plain Teks"), plainTypeSelect),
			nil, nil, nil, plainEntry,
		),
		container.NewBorder(widget.NewLabel("Cipher Teks"), nil, nil, nil, cipherEntry),
	)

	getContext := func() (*aes.AesContext, bool) {
		variant := variantMap[variantSelect.SelectedIndex()]

		var key []uint8
		if keyTypeSelect.SelectedIndex() == 0 {
			key = []uint8(keyEntry.Text)
		} else {
			k, err := decodeHexString(strings.TrimSpace(keyEntry.Text))
			if err != nil {
				dialog.NewError(errors.Join(ErrKunciHex, err), w).Show()
				return nil, false
			}
			key = k
		}

		ctx, err := aes.NewAesContext(variant, key)
		if err != nil {
			dialog.NewError(err, w).Show()
			return nil, false
		}
		return ctx, true
	}

	variantSelect.SetSelectedIndex(0)
	keyTypeSelect.SetSelectedIndex(0)
	plainTypeSelect.SetSelectedIndex(0)

	leftLayout := container.NewVBox(
		widget.NewLabelWithStyle("Jenis AES", fyne.TextAlignCenter, fyne.TextStyle{}),
		variantSelect,
		widget.NewLabelWithStyle("Kunci", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyTypeSelect,
		keyEntry,
		widget.NewLabelWithStyle("Padding: PKCS#7", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewButton("Enkripsi", func() {
			ctx, success := getContext()
			if !success {
				return
			}

			var plain []uint8
			if plainTypeSelect.SelectedIndex() == 0 {
				plain = []uint8(plainEntry.Text)
			} else {
				decoded, err := decodeHexString(strings.TrimSpace(plainEntry.Text))
				if err != nil {
					dialog.NewError(errors.Join(ErrPlainDecode, err), w).Show()
					return
				}
				plain = decoded
			}

			cipher := ctx.EncryptECB(plain)

			hex, err := encodeHexString(cipher)

			if err != nil {
				dialog.NewError(errors.Join(ErrCipherEncode, err), w).Show()
			}

			cipherEntry.SetText(hex)
		}),
		widget.NewButton("Dekripsi", func() {
			ctx, success := getContext()
			if !success {
				return
			}

			cipher, err := decodeHexString(strings.TrimSpace(cipherEntry.Text))
			if err != nil {
				dialog.NewError(errors.Join(ErrCipherDecode, err), w).Show()
				return
			}

			plain, err := ctx.DecryptECB(cipher)
			if err != nil {
				dialog.NewError(errors.Join(ErrDecrypt, err), w).Show()
				return
			}

			if plainTypeSelect.SelectedIndex() == 0 {
				plainEntry.SetText(string(plain))
			} else {
				hex, err := encodeHexString(plain)

				if err != nil {
					dialog.NewError(errors.Join(ErrPlainEncode, err), w).Show()
				}

				plainEntry.SetText(hex)
			}
		}),
	)

	return container.NewBorder(
		widget.NewLabel("AES"), nil,
		container.NewStack(container.NewGridWrap(fyne.NewSize(300, 0)), leftLayout), nil,
		container.NewStack(container.NewGridWrap(fyne.NewSize(300, 0)), mainContainer),
	)
}
