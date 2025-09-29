package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func createSuperTab(w fyne.Window) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Generate Key", createSuperKeyGenSubTab(w)),
		container.NewTabItem("Encrypt", createSuperEncryptSubTab(w)),
		container.NewTabItem("Decrypt", createSuperDecryptSubTab(w)),
	)
}

func createSuperKeyGenSubTab(w fyne.Window) fyne.CanvasObject {
	privateKeyEntry := widget.NewMultiLineEntry()
	publicKeyEntry := widget.NewMultiLineEntry()

	privateKeyEntry.Wrapping = fyne.TextWrapBreak
	publicKeyEntry.Wrapping = fyne.TextWrapBreak

	privateKeyEntry.Disable()
	publicKeyEntry.Disable()

	generateButton := widget.NewButton("Generate Key", func() {
		go func() {
			key, err := SuperGenerateKey()
			fyne.Do(func() {
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				privateKeyEntry.SetText(key.private)
				publicKeyEntry.SetText(key.public)
			})
		}()
	})

	return container.NewPadded(NewFlexibleRow(
		NewFlexibleItem(false, generateButton),
		NewFlexibleItem(false, widget.NewLabel("Private Key")),
		NewFlexibleItem(true, privateKeyEntry),
		NewFlexibleItem(false, widget.NewLabel("Public Key")),
		NewFlexibleItem(true, publicKeyEntry),
	))
}

func createSuperEncryptSubTab(w fyne.Window) fyne.CanvasObject {
	publicKeyEntry := widget.NewEntry()
	plainEntry := widget.NewMultiLineEntry()
	cipherEntry := widget.NewMultiLineEntry()

	plainEntry.Wrapping = fyne.TextWrapBreak
	cipherEntry.Wrapping = fyne.TextWrapBreak
	cipherEntry.Disable()

	buttonEncrypt := widget.NewButton("Encrypt", func() {
		key, err := SuperDecodePublicKey(publicKeyEntry.Text)
		if err != nil {
			dialog.NewError(err, w).Show()
			return
		}

		encrypted, err := SuperEncrypt(key, []uint8(plainEntry.Text))
		if err != nil {
			dialog.NewError(errors.Join(ErrEncrypt, err), w).Show()
			return
		}

		hex, err := encodeHexString(encrypted)
		if err != nil {
			dialog.NewError(errors.Join(ErrCipherEncode, err), w).Show()
			return
		}

		cipherEntry.SetText(hex)
	})

	return container.NewPadded(NewFlexibleRow(
		NewFlexibleItem(false, widget.NewForm(
			widget.NewFormItem("Kunci Publik", publicKeyEntry),
		)),
		NewFlexibleItem(true, container.NewBorder(
			widget.NewLabel("Plain"), nil, nil, nil,
			plainEntry,
		)),
		NewFlexibleItem(false, buttonEncrypt),
		NewFlexibleItem(true, container.NewBorder(
			widget.NewLabel("Cipher (Hex)"), nil, nil, nil,
			cipherEntry,
		)),
	))
}
func createSuperDecryptSubTab(w fyne.Window) fyne.CanvasObject {
	privateKeyEntry := widget.NewEntry()
	plainEntry := widget.NewMultiLineEntry()
	cipherEntry := widget.NewMultiLineEntry()

	plainEntry.Wrapping = fyne.TextWrapBreak
	cipherEntry.Wrapping = fyne.TextWrapBreak
	plainEntry.Disable()

	buttonEncrypt := widget.NewButton("Decrypt", func() {
		key, err := SuperDecodePrivateKey(privateKeyEntry.Text)
		if err != nil {
			dialog.NewError(err, w).Show()
			return
		}

		cipher, err := decodeHexString(cipherEntry.Text)
		if err != nil {
			dialog.NewError(errors.Join(ErrCipherDecode, err), w).Show()
			return
		}

		decrypted, err := SuperDecrypt(key, cipher)
		if err != nil {
			dialog.NewError(errors.Join(ErrEncrypt, err), w).Show()
			return
		}

		plainEntry.SetText(string(decrypted))
	})

	return container.NewPadded(NewFlexibleRow(
		NewFlexibleItem(false, widget.NewForm(
			widget.NewFormItem("Kunci Private", privateKeyEntry),
		)),
		NewFlexibleItem(true, container.NewBorder(
			widget.NewLabel("Cipher (Hex)"), nil, nil, nil,
			cipherEntry,
		)),
		NewFlexibleItem(false, buttonEncrypt),
		NewFlexibleItem(true, container.NewBorder(
			widget.NewLabel("Plain"), nil, nil, nil,
			plainEntry,
		)),
	))
}
