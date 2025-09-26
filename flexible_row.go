package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Layout = (*flexibleLayout)(nil)

type flexibleLayout struct {
}

func NewFlexibleLayout() fyne.Layout {
	return &flexibleLayout{}
}

func (f *flexibleLayout) countColumns(objects []fyne.CanvasObject) int {
	count := 0
	for _, object := range objects {
		if object.Visible() {
			count++
		}
	}
	return count
}

// Layout implements fyne.Layout.
func (f *flexibleLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	columns := f.countColumns(objects)
	padding := theme.Padding()

	padWidth := float32(columns-1) * padding
	cellWidth := float32(size.Width-padWidth) / float32(columns)

	i := 0
	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		x0 := (cellWidth + padding) * float32(i)

		object.Move(fyne.NewPos(x0, 0))
		object.Resize(fyne.NewSize(cellWidth, size.Height))

		i++
	}
}

// MinSize implements fyne.Layout.
func (f *flexibleLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	columns := f.countColumns(objects)

	minSize := fyne.NewSize(0, 0)
	for _, object := range objects {
		if object.Visible() {
			minSize = minSize.Max(object.MinSize())
		}
	}

	padding := theme.Padding()
	width := minSize.Width * float32(columns)
	xpad := padding * fyne.Max(float32(columns-1), 0)

	return fyne.NewSize(width+xpad, minSize.Height)
}
