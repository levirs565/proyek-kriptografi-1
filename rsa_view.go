package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createRsaTab() fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Generate Key", createKeyGenSubTab()),
		container.NewTabItem("Enkripsi", createEncryptSubTab()),
		container.NewTabItem("Dekripsi", createDecryptSubTab()),
	)
}

func createKeyGenSubTab() fyne.CanvasObject {
	keySizeSelect := widget.NewSelect([]string{"512 bit (Tidak Aman)", "1024 bit", "2048 bit", "3072 bit", "4096 bit"}, nil)
	keySizeSelect.SetSelected("2048 bit")

	pOut := widget.NewMultiLineEntry(); pOut.SetPlaceHolder("Bilangan prima P..."); pOut.Disable()
	qOut := widget.NewMultiLineEntry(); qOut.SetPlaceHolder("Bilangan prima Q..."); qOut.Disable()
	nOut := widget.NewMultiLineEntry(); nOut.SetPlaceHolder("Modulus N..."); nOut.Disable()
	eOut := widget.NewMultiLineEntry(); eOut.SetPlaceHolder("Public Exponent E..."); eOut.Disable()
	dOut := widget.NewMultiLineEntry(); dOut.SetPlaceHolder("Private Exponent D..."); dOut.Disable()

	generateButton := widget.NewButton("GENERATE RSA KEY PAIR", func() {
		pOut.SetText(""); qOut.SetText(""); nOut.SetText(""); eOut.SetText(""); dOut.SetText("")
		
		selected := keySizeSelect.Selected
		bitStr := strings.Split(selected, " ")[0]
		bits, err := strconv.Atoi(bitStr)
		if err != nil {
			log.Printf("Error parsing bit size: %v", err)
			return
		}

		pOut.SetText(fmt.Sprintf("Menghasilkan kunci %d bit... (mungkin butuh beberapa detik)", bits))
		
		go func() {
			vals, err := GenerateKeys(bits)
			if err != nil {
				dOut.SetText(fmt.Sprintf("Error: %v", err))
				return
			}
			
			pOut.SetText(vals.P.String())
			qOut.SetText(vals.Q.String())
			nOut.SetText(vals.N.String())
			eOut.SetText(vals.E.String())
			dOut.SetText(vals.D.String())
		}()
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Ukuran Kunci", Widget: keySizeSelect},
			{Text: "Bilangan Prima P", Widget: pOut},
			{Text: "Bilangan Prima Q", Widget: qOut},
			{Text: "Modulus N (Publik)", Widget: nOut},
			{Text: "Exponent E (Publik)", Widget: eOut},
			{Text: "Exponent D (Privat)", Widget: dOut},
		},
	}
	
	return container.NewVBox(
		widget.NewLabel("Pilih ukuran kunci dan klik generate untuk membuat pasangan kunci baru."),
		generateButton,
		form,
	)
}

func createEncryptSubTab() fyne.CanvasObject {
	messageEntry := widget.NewEntry(); messageEntry.SetPlaceHolder("Tulis pesan Anda di sini")
	eEntry := widget.NewEntry(); eEntry.SetPlaceHolder("Kunci Publik E")
	nEntry := widget.NewEntry(); nEntry.SetPlaceHolder("Modulus N")
	outputArea := widget.NewMultiLineEntry(); outputArea.SetPlaceHolder("Ciphertext (hasil enkripsi) akan ditampilkan di sini..."); outputArea.Wrapping = fyne.TextWrapWord; outputArea.Disable()
	parseInput := func(entry *widget.Entry) *big.Int { if entry.Text == "" { return nil }; val, ok := new(big.Int).SetString(entry.Text, 10); if !ok { log.Printf("Input tidak valid: %s", entry.Text); return nil }; return val }
	encryptButton := widget.NewButton("ENKRIPSI", func() {
		message := messageEntry.Text; e := parseInput(eEntry); n := parseInput(nEntry)
		if message == "" || e == nil || n == nil { outputArea.SetText("Error: Pesan, E, dan N harus diisi."); return }
		m := StringToBigInt(message)
		c, err := Encrypt(m, e, n); if err != nil { outputArea.SetText(fmt.Sprintf("Error saat enkripsi: %v", err)); return }
		outputArea.SetText(fmt.Sprintf("Ciphertext (Integer):\n%s", c.String()))
	})
	form := &widget.Form{ Items: []*widget.FormItem{ {Text: "PESAN (M)", Widget: messageEntry}, {Text: "PUBLIC KEY (E)", Widget: eEntry}, {Text: "MODULUS (N)", Widget: nEntry}, }}
	return container.NewBorder( container.NewVBox( widget.NewLabel("Masukkan pesan dan Kunci Publik (E, N) untuk mengenkripsi."), form, encryptButton, ), nil, nil, nil, outputArea,)
}

func createDecryptSubTab() fyne.CanvasObject {
	cEntry := widget.NewEntry(); cEntry.SetPlaceHolder("Ciphertext (Integer)")
	eEntry := widget.NewEntry(); eEntry.SetPlaceHolder("Contoh: 65537 atau 17")
	nEntry := widget.NewEntry(); nEntry.SetPlaceHolder("Integer (Hasil p * q)")
	dEntry := widget.NewEntry(); dEntry.SetPlaceHolder("Dihitung dari E dan PHI")
	pEntry := widget.NewEntry(); pEntry.SetPlaceHolder("Bilangan prima")
	qEntry := widget.NewEntry(); qEntry.SetPlaceHolder("Bilangan prima")
	phiEntry := widget.NewEntry(); phiEntry.SetPlaceHolder("Dihitung dari P dan Q")
	outputArea := widget.NewMultiLineEntry(); outputArea.SetPlaceHolder("Hasil akan ditampilkan di sini..."); outputArea.Wrapping = fyne.TextWrapWord; outputArea.Disable()
	displayMode := widget.NewRadioGroup([]string{"String", "Integer", "Hex", "Computed Values"}, nil); displayMode.SetSelected("String")
	parseInput := func(entry *widget.Entry) *big.Int { if entry.Text == "" { return nil }; val, ok := new(big.Int).SetString(entry.Text, 10); if !ok { log.Printf("Input tidak valid: %s", entry.Text); return nil }; return val }
	calculateButton := widget.NewButton("► KALKULASI / DEKRIPSI", func() {
		outputArea.SetText("Menghitung..."); vals := RSAValues{ C: parseInput(cEntry), E: parseInput(eEntry), N: parseInput(nEntry), D: parseInput(dEntry), P: parseInput(pEntry), Q: parseInput(qEntry)}; if err := vals.CalculateMissingValues(); err != nil { outputArea.SetText(fmt.Sprintf("Error: %v", err)); return }
		if vals.N != nil { nEntry.SetText(vals.N.String()) }; if vals.D != nil { dEntry.SetText(vals.D.String()) }; if vals.Phi != nil { phiEntry.SetText(vals.Phi.String()) }
		m, err := Decrypt(vals.C, vals.D, vals.N)
		if err != nil { computedVals := "Input tidak cukup untuk dekripsi.\nNilai yang berhasil dihitung:\n"; if vals.N != nil { computedVals += fmt.Sprintf("  N = %s\n", vals.N.String()) }; if vals.Phi != nil { computedVals += fmt.Sprintf("  Φ = %s\n", vals.Phi.String()) }; if vals.D != nil { computedVals += fmt.Sprintf("  D = %s\n", vals.D.String()) }; outputArea.SetText(computedVals); return }
		var result string
		switch displayMode.Selected { case "String": result = fmt.Sprintf("Plaintext (String):\n%s", BigIntToString(m)); case "Integer": result = fmt.Sprintf("Plaintext (Integer):\n%s", m.String()); case "Hex": result = fmt.Sprintf("Plaintext (Hexadecimal):\n0x%x", m); case "Computed Values": result = fmt.Sprintf("Nilai yang Dihitung:\nN = %s\nΦ = %s\nD = %s", vals.N.String(), vals.Phi.String(), vals.D.String()) }
		outputArea.SetText(result)
	})
	form := &widget.Form{ Items: []*widget.FormItem{ {Text: "CIPHERTEXT (C)", Widget: cEntry}, {Text: "PUBLIC KEY (E)", Widget: eEntry}, {Text: "MODULUS (N)", Widget: nEntry}, {Text: "PRIVATE KEY (D)", Widget: dEntry}, {Text: "FAKTOR 1 (P)", Widget: pEntry}, {Text: "FAKTOR 2 (Q)", Widget: qEntry}, {Text: "NILAI PHI (Φ)", Widget: phiEntry}, }}
	return container.NewBorder( container.NewVBox( widget.NewLabel("Masukkan nilai yang diketahui untuk melakukan dekripsi/kalkulasi."), form, displayMode, calculateButton, ), nil, nil, nil, outputArea,)
}
