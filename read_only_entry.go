package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ReadOnlyEntry struct {
	widget.Entry
}

var _ fyne.Focusable = (*ReadOnlyEntry)(nil)
var _ fyne.Shortcutable = (*ReadOnlyEntry)(nil)

func NewReadOnlyEntry() *ReadOnlyEntry {
	res := &ReadOnlyEntry{
		widget.Entry{
			Wrapping: fyne.TextWrap(fyne.TextTruncateClip),
		},
	}
	res.ExtendBaseWidget(res)
	return res
}

func NewMultilineReadOnlyEntry() *ReadOnlyEntry {
	res := &ReadOnlyEntry{
		widget.Entry{
			MultiLine: true,
			Wrapping:  fyne.TextWrap(fyne.TextTruncateClip),
		},
	}
	res.ExtendBaseWidget(res)
	return res
}

func (e *ReadOnlyEntry) TypedRune(r rune) {
	// do nothing
}

func (e *ReadOnlyEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyBackspace, fyne.KeyDelete, fyne.KeyReturn, fyne.KeyEnter, fyne.KeyTab:
		return
	default:
		e.Entry.TypedKey(key)
	}
}

func (e *ReadOnlyEntry) TypedShortcut(shortcut fyne.Shortcut) {
	switch shortcut.ShortcutName() {
	case "Paste", "Cut", "Undo", "Redo",
		"CustomDesktop:Control+BackSpace", "CustomDesktop:Control+Delete":
		return
	default:
		// print(shortcut.ShortcutName())
		e.Entry.TypedShortcut(shortcut)
	}
}
