package main

import (
	"crypto/rand"
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
var ErrEncrypt = errors.New("gagal encrypt")
var ErrIvLength = errors.New("ukuran iv harus 16 byte")
var ErrIvDecode = errors.New("iv gagal di hex decode")

var inputTypes = []string{"Raw", "Hex"}

func createAesTab(w fyne.Window) fyne.CanvasObject {
	aes.AESInit()

	variantSelect := widget.NewSelect([]string{"AES-128", "AES-192", "AES-256"}, func(s string) {})
	variantMap := []aes.AESVariant{aes.AES128, aes.AES192, aes.AES256}

	paddingSelect := widget.NewSelect([]string{"PKCS#7", "No Padding"}, func(s string) {})
	paddingMap := []aes.Padding{aes.PKCS7Padding, aes.NoPadding}

	keyTypeSelect := widget.NewSelect(inputTypes, func(s string) {})
	keyEntry := widget.NewEntry()
	ivLabel := widget.NewLabelWithStyle("Initialization Vector (Hex)", fyne.TextAlignCenter, fyne.TextStyle{})
	ivEntry := widget.NewEntry()
	plainEntry := widget.NewMultiLineEntry()
	plainTypeSelect := widget.NewSelect(inputTypes, func(s string) {})
	cipherEntry := widget.NewMultiLineEntry()

	generateIvButton := widget.NewButton("Random IV", func() {
		var iv [16]uint8
		rand.Read(iv[:])

		hex, err := encodeHexString(iv[:])

		if err != nil {
			dialog.NewError(err, w).Show()
			return
		}

		ivEntry.SetText(hex)
	})

	modeSelect := widget.NewSelect([]string{"ECB", "CBC"}, func(s string) {
		if s == "ECB" {
			ivEntry.Hide()
			ivLabel.Hide()
			generateIvButton.Hide()
		} else {
			ivEntry.Show()
			ivLabel.Show()
			generateIvButton.Show()
		}
	})

	mainContainer := container.NewGridWithRows(2,
		container.NewBorder(
			container.NewVBox(widget.NewLabel("Plain Teks"), plainTypeSelect),
			nil, nil, nil, plainEntry,
		),
		container.NewBorder(widget.NewLabel("Cipher Teks (Hex)"), nil, nil, nil, cipherEntry),
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

		if modeSelect.SelectedIndex() == 1 {
			iv, err := decodeHexString(strings.TrimSpace(ivEntry.Text))
			if err != nil {
				dialog.NewError(errors.Join(ErrIvDecode, err), w).Show()
				return nil, false
			}

			if len(iv) != 16 {
				dialog.NewError(ErrIvLength, w).Show()
				return nil, false
			}

			ctx.SetIv([16]uint8(iv))
		}

		return ctx, true
	}

	variantSelect.SetSelectedIndex(0)
	keyTypeSelect.SetSelectedIndex(0)
	plainTypeSelect.SetSelectedIndex(0)
	paddingSelect.SetSelectedIndex(0)
	modeSelect.SetSelectedIndex(0)

	leftLayout := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Jenis AES", variantSelect),
			widget.NewFormItem("Mode", modeSelect),
			widget.NewFormItem("Padding", paddingSelect),
		),
		widget.NewLabelWithStyle("Kunci", fyne.TextAlignCenter, fyne.TextStyle{}),
		keyTypeSelect,
		keyEntry,
		ivLabel,
		ivEntry,
		generateIvButton,
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

			var cipher []uint8
			var err error
			padding := paddingMap[paddingSelect.SelectedIndex()]
			if modeSelect.SelectedIndex() == 0 {
				cipher, err = ctx.EncryptECB(plain, padding)
			} else {
				cipher, err = ctx.EncryptCBC(plain, padding)
			}

			if err != nil {
				dialog.NewError(errors.Join(ErrEncrypt, err), w).Show()
			}

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

			var plain []uint8
			padding := paddingMap[paddingSelect.SelectedIndex()]

			if modeSelect.SelectedIndex() == 0 {
				plain, err = ctx.DecryptECB(cipher, padding)
			} else {
				plain, err = ctx.DecryptCBC(cipher, padding)
			}

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
		container.NewPadded(container.NewGridWrap(fyne.NewSize(300, 0)), leftLayout), nil,
		container.NewStack(container.NewGridWrap(fyne.NewSize(300, 0)), mainContainer),
	)
}
