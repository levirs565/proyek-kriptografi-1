package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Layout = (*flexibleLayout)(nil)

type FlexibleItem struct {
	object fyne.CanvasObject
	grow   bool
}

func NewFlexibleItem(grow bool, object fyne.CanvasObject) FlexibleItem {
	return FlexibleItem{
		object: object,
		grow:   grow,
	}
}

type flexibleLayout struct {
	vertical bool
	items    []FlexibleItem
}

func NewFlexibleContainer(vertical bool, items ...FlexibleItem) *fyne.Container {
	var layout = &flexibleLayout{
		vertical: vertical,
		items:    items,
	}
	objects := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		objects[i] = item.object
	}
	return container.New(layout, objects...)
}

func NewFlexibleColumn(items ...FlexibleItem) *fyne.Container {
	return NewFlexibleContainer(false, items...)
}

func NewFlexibleRow(items ...FlexibleItem) *fyne.Container {
	return NewFlexibleContainer(true, items...)
}

func (f *flexibleLayout) countVisibleItems() int {
	count := 0
	for _, item := range f.items {
		if item.object.Visible() {
			count++
		}
	}
	return count
}

func (f *flexibleLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	columns := f.countVisibleItems()
	padding := theme.Padding()

	padSize := float32(columns-1) * padding
	var fullSize float32
	if !f.vertical {
		fullSize = size.Width
	} else {
		fullSize = size.Height
	}

	fixedSize := float32(0)
	growCount := 0
	for _, item := range f.items {
		if item.object.Visible() {
			if item.grow {
				growCount++
				continue
			}
			minSize := item.object.MinSize()
			if !f.vertical {
				fixedSize += minSize.Width
			} else {
				fixedSize += minSize.Height
			}
		}
	}

	cellSize := float32(fullSize-padSize-fixedSize) / float32(growCount)

	cellPos := float32(0)
	for _, item := range f.items {
		if !item.object.Visible() {
			continue
		}

		currentCellSize := cellSize
		if !item.grow {
			minSize := item.object.MinSize()
			if !f.vertical {
				currentCellSize = minSize.Width
			} else {
				currentCellSize = minSize.Height
			}
		}

		if !f.vertical {
			item.object.Move(fyne.NewPos(cellPos, 0))
			item.object.Resize(fyne.NewSize(currentCellSize, size.Height))
		} else {
			item.object.Move(fyne.NewPos(0, cellPos))
			item.object.Resize(fyne.NewSize(size.Width, currentCellSize))

		}

		cellPos += currentCellSize + padding
	}
}

func (f *flexibleLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	columns := f.countVisibleItems()
	minSize := fyne.NewSize(0, 0)
	fixedSize := fyne.NewSize(0, 0)
	growCount := 0

	for _, item := range f.items {
		if item.object.Visible() {
			if item.grow {
				minSize = minSize.Max(item.object.MinSize())
				growCount++
			} else {
				fixedSize = fixedSize.Add(item.object.MinSize())
			}
		}
	}

	padding := theme.Padding()

	if !f.vertical {
		width := minSize.Width*float32(growCount) + fixedSize.Width
		xpad := padding * fyne.Max(float32(columns-1), 0)

		return fyne.NewSize(width+xpad, minSize.Height)
	} else {
		height := minSize.Height*float32(growCount) + fixedSize.Height
		ypad := padding * fyne.Max(float32(columns-1), 0)

		return fyne.NewSize(minSize.Width, height+ypad)
	}
}
