package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var ErrRSAEDecode = errors.New("gagal decode E")
var ErrRSADDecode = errors.New("gagal decode D")
var ErrRSANDecode = errors.New("gagal decode N")

func createRsaTab(w fyne.Window) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Generate Key", createKeyGenSubTab()),
		container.NewTabItem("Enkripsi", createEncryptSubTab(w)),
		container.NewTabItem("Dekripsi", createDecryptSubTab(w)),
	)
}

func createKeyGenSubTab() fyne.CanvasObject {
	keySizeSelect := widget.NewSelect([]string{"512 bit (Tidak Aman)", "1024 bit", "2048 bit", "3072 bit", "4096 bit"}, nil)
	keySizeSelect.SetSelected("2048 bit")

	pEntry := widget.NewMultiLineEntry()
	qEntry := widget.NewMultiLineEntry()
	nEntry := widget.NewMultiLineEntry()
	eEntry := widget.NewEntry()
	dEntry := widget.NewMultiLineEntry()

	pEntry.Wrapping = fyne.TextWrapBreak
	qEntry.Wrapping = fyne.TextWrapBreak
	nEntry.Wrapping = fyne.TextWrapBreak
	dEntry.Wrapping = fyne.TextWrapBreak

	pEntry.Disable()
	qEntry.Disable()
	nEntry.Disable()
	eEntry.Disable()
	dEntry.Disable()

	var generateButton *widget.Button
	generateButton = widget.NewButton("GENERATE RSA KEY PAIR", func() {
		pEntry.SetText("")
		qEntry.SetText("")
		nEntry.SetText("")
		eEntry.SetText("")
		dEntry.SetText("")

		selected := keySizeSelect.Selected
		bitStr := strings.Split(selected, " ")[0]
		bits, err := strconv.Atoi(bitStr)
		if err != nil {
			log.Printf("Error parsing bit size: %v", err)
			return
		}

		pEntry.SetText(fmt.Sprintf("Menghasilkan kunci %d bit... (mungkin butuh beberapa detik)", bits))

		generateButton.Disable()
		go func() {
			vals, err := RSAGenerateKeys(bits)
			if err != nil {
				dEntry.SetText(fmt.Sprintf("Error: %v", err))
				return
			}

			fyne.Do(func() {
				pEntry.SetText(vals.P.String())
				qEntry.SetText(vals.Q.String())
				nEntry.SetText(vals.N.String())
				eEntry.SetText(vals.E.String())
				dEntry.SetText(vals.D.String())
				generateButton.Enable()
			})
		}()
	})

	return container.NewPadded(container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Ukuran Kunci"), keySizeSelect, generateButton,
		),
		widget.NewForm(
			widget.NewFormItem("Bilangan Prima P", pEntry),
			widget.NewFormItem("Bilangan Prima Q", qEntry),
			widget.NewFormItem("Modulus N (Publik)", nEntry),
			widget.NewFormItem("Exponent E (Publik)", eEntry),
			widget.NewFormItem("Exponent D (Privat)", dEntry),
		),
	))
}

func parseBigNumberInput(entry *widget.Entry) (*big.Int, error) {
	if entry.Text == "" {
		return nil, fmt.Errorf("input kosong")
	}
	val, ok := new(big.Int).SetString(entry.Text, 10)
	if !ok {
		return nil, fmt.Errorf("input tidak valid: %s", entry.Text)
	}
	return val, nil
}

func createEncryptSubTab(w fyne.Window) fyne.CanvasObject {
	plainEntry := widget.NewMultiLineEntry()
	eEntry := widget.NewEntry()
	nEntry := widget.NewEntry()
	cipherEntry := widget.NewMultiLineEntry()
	cipherEntry.Wrapping = fyne.TextWrapWord
	cipherEntry.Disable()

	plainTypeSelect := widget.NewSelect(inputTypes, func(s string) {})
	plainTypeSelect.SetSelectedIndex(0)

	encryptButton := widget.NewButton("Enkripsi", func() {
		e, err := parseBigNumberInput(eEntry)
		if err != nil {
			dialog.NewError(errors.Join(ErrRSAEDecode, err), w).Show()
		}
		n, err := parseBigNumberInput(nEntry)
		if err != nil {
			dialog.NewError(errors.Join(ErrRSANDecode, err), w).Show()
		}
		m := new(big.Int)
		if plainTypeSelect.SelectedIndex() == 0 {
			m.SetBytes([]byte(plainEntry.Text))
		} else {
			bytes, err := decodeHexString(plainEntry.Text)
			if err != nil {
				dialog.NewError(errors.Join(ErrPlainDecode, err), w).Show()
			}
			m.SetBytes(bytes)
		}
		c, err := RSAEncrypt(m, e, n)
		if err != nil {
			dialog.NewError(errors.Join(ErrEncrypt, err), w).Show()
			return
		}
		cipherEntry.SetText(c.String())
	})

	form := widget.NewForm(
		widget.NewFormItem("Public Key (E)", eEntry),
		widget.NewFormItem("Modulus (N)", nEntry),
	)
	mainContainer := NewFlexibleRow(
		NewFlexibleItem(true, container.NewBorder(
			container.NewHBox(widget.NewLabel("Plain"), plainTypeSelect),
			nil, nil, nil,
			plainEntry,
		)),
		NewFlexibleItem(false, encryptButton),
		NewFlexibleItem(true, container.NewBorder(
			widget.NewLabel("Cipher (Number)"),
			nil, nil, nil,
			cipherEntry,
		)),
	)
	return container.NewPadded(container.NewBorder(
		form,
		nil, nil, nil,
		mainContainer,
	))
}

func createDecryptSubTab(w fyne.Window) fyne.CanvasObject {
	cipherEntry := widget.NewEntry()
	nEntry := widget.NewEntry()
	nEntry.SetPlaceHolder("")
	dEntry := widget.NewEntry()
	dEntry.SetPlaceHolder("")
	plainEntry := widget.NewMultiLineEntry()
	plainEntry.Wrapping = fyne.TextWrapBreak
	plainEntry.Disable()
	plainTypeSelect := widget.NewSelect(inputTypes, nil)
	plainTypeSelect.SetSelectedIndex(0)
	calculateButton := widget.NewButton("Dekripsi", func() {
		D, err := parseBigNumberInput(dEntry)
		if err != nil {
			dialog.NewError(errors.Join(ErrRSADDecode), w).Show()
			return
		}
		N, err := parseBigNumberInput(nEntry)
		if err != nil {
			dialog.NewError(errors.Join(ErrRSANDecode), w).Show()
			return
		}
		cipher, err := parseBigNumberInput(cipherEntry)
		if err != nil {
			dialog.NewError(errors.Join(ErrCipherDecode), w).Show()
			return
		}
		plain, err := RSADecrypt(cipher, D, N)
		if err != nil {
			dialog.NewError(errors.Join(ErrDecrypt, err), w).Show()
			return
		}
		var result string
		if plainTypeSelect.SelectedIndex() == 0 {
			result = string(plain.Bytes())
		} else {
			result, err = encodeHexString(plain.Bytes())
			if err != nil {
				dialog.NewError(errors.Join(ErrPlainEncode, err), w).Show()
				return
			}
		}
		plainEntry.SetText(result)
	})
	form := widget.NewForm(
		widget.NewFormItem("Modulus (N)", nEntry),
		widget.NewFormItem("Private Key (D)", dEntry),
	)
	return container.NewPadded(container.NewBorder(
		form, nil, nil, nil,
		NewFlexibleRow(
			NewFlexibleItem(true, container.NewBorder(
				widget.NewLabel("Ciphertext/C (Integer)"),
				nil, nil, nil,
				cipherEntry,
			)),
			NewFlexibleItem(false, calculateButton),
			NewFlexibleItem(true, container.NewBorder(
				container.NewHBox(widget.NewLabel("Plain"), plainTypeSelect),
				nil, nil, nil,
				plainEntry,
			)),
		),
	))
}
