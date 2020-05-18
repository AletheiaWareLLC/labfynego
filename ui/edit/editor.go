/*
 * Copyright 2020 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package edit

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"log"
	"math"
	"sync"
	"unicode"
)

type Editor struct {
	widget.BaseWidget
	sync.Mutex
	Cursor      uint64
	Selection   uint64
	IsSelecting bool
	IsFocused   bool
	TextAlign   fyne.TextAlign
	TextColor   color.Color
	TextSize    int
	TextStyle   fyne.TextStyle
	TextWrap    fyne.TextWrap
	Buffer      []rune
	Lines       []*Line

	shortcut fyne.ShortcutHandler
}

func NewEditor() *Editor {
	e := &Editor{
		TextAlign: fyne.TextAlignLeading,
		TextColor: theme.TextColor(),
		TextSize:  theme.TextSize(),
		TextStyle: fyne.TextStyle{},
		TextWrap:  fyne.TextWrapWord,
	}
	e.ExtendBaseWidget(e)
	e.AddShortcuts()
	return e
}

func (e *Editor) AddShortcuts() {
	e.shortcut.AddShortcut(&fyne.ShortcutCut{}, func(se fyne.Shortcut) {
		cut := se.(*fyne.ShortcutCut)
		e.CutToClipboard(cut.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutCopy{}, func(se fyne.Shortcut) {
		cpy := se.(*fyne.ShortcutCopy)
		e.CopyToClipboard(cpy.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutPaste{}, func(se fyne.Shortcut) {
		paste := se.(*fyne.ShortcutPaste)
		e.PasteFromClipboard(paste.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutSelectAll{}, func(se fyne.Shortcut) {
		e.SelectAll()
	})
}

func (e *Editor) Refresh() {
	log.Println("Editor.Refresh")
	e.Lock()
	textWrap := e.TextWrap
	textStyle := e.TextStyle
	textSize := e.TextSize
	maxWidth := e.Size().Width - 2*theme.Padding()

	// TODO only redo this if Buffer, textWrap, or maxWidth has changed
	e.Lines = lineBounds(e.Buffer, textWrap, maxWidth, func(text []rune) int {
		return fyne.MeasureText(string(text), textSize, textStyle).Width
	})

	log.Println("Lines:", e.Lines)
	e.Unlock()
	e.BaseWidget.Refresh()
}

func (e *Editor) SetText(text string) {
	e.Buffer = []rune(text)
	e.Refresh()
}

// splitLines accepts a slice of runes and returns a slice containing the
// start and end indicies of each line delimited by the newline character.
func splitLines(text []rune) []*Line {
	var low, high int
	var Lines []*Line
	length := len(text)
	for i := 0; i < length; i++ {
		if text[i] == '\n' {
			high = i
			Lines = append(Lines, &Line{start: low, end: high})
			low = i + 1
		}
	}
	return append(Lines, &Line{start: low, end: length})
}

// binarySearch accepts a function that checks if the text width less the maximum width and the start and end rune index
// binarySearch returns the index of rune located as close to the maximum line width as possible
func binarySearch(lessMaxWidth func(int, int) bool, low int, maxHigh int) int {
	if low >= maxHigh {
		return low
	}
	if lessMaxWidth(low, maxHigh) {
		return maxHigh
	}
	high := low
	delta := maxHigh - low
	for delta > 0 {
		delta /= 2
		if lessMaxWidth(low, high+delta) {
			high += delta
		}
	}
	for (high < maxHigh) && lessMaxWidth(low, high+1) {
		high++
	}
	return high
}

// findSpaceIndex accepts a slice of runes and a fallback index
// findSpaceIndex returns the index of the last space in the text, or fallback if there are no spaces
func findSpaceIndex(text []rune, fallback int) int {
	curIndex := fallback
	for ; curIndex >= 0; curIndex-- {
		if unicode.IsSpace(text[curIndex]) {
			break
		}
	}
	if curIndex < 0 {
		return fallback
	}
	return curIndex
}

// lineBounds accepts a slice of runes, a wrapping mode, a maximum line width and a function to measure line width.
// lineBounds returns a slice containing the start and end indicies of each line with the given wrapping applied.
func lineBounds(text []rune, wrap fyne.TextWrap, maxWidth int, measurer func([]rune) int) []*Line {

	Lines := splitLines(text)
	if maxWidth <= 0 || wrap == fyne.TextWrapOff {
		return Lines
	}

	checker := func(low int, high int) bool {
		return measurer(text[low:high]) <= maxWidth
	}

	var bounds []*Line
	for _, l := range Lines {
		low := l.start
		high := l.end
		if low == high {
			bounds = append(bounds, l)
			continue
		}
		switch wrap {
		case fyne.TextTruncate:
			high = binarySearch(checker, low, high)
			bounds = append(bounds, &Line{start: low, end: high})
		case fyne.TextWrapBreak:
			for low < high {
				if measurer(text[low:high]) <= maxWidth {
					bounds = append(bounds, &Line{start: low, end: high})
					low = high
					high = l.end
				} else {
					high = binarySearch(checker, low, high)
				}
			}
		case fyne.TextWrapWord:
			for low < high {
				sub := text[low:high]
				if measurer(sub) <= maxWidth {
					bounds = append(bounds, &Line{start: low, end: high})
					low = high
					high = l.end
					if low < high && unicode.IsSpace(text[low]) {
						low++
					}
				} else {
					last := low + len(sub) - 1
					high = low + findSpaceIndex(sub, binarySearch(checker, low, last)-low)
				}
			}
		}
	}
	return bounds
}

func (e *Editor) FocusGained() {
	e.IsFocused = true
	e.Refresh()
}

func (e *Editor) FocusLost() {
	e.IsFocused = false
	e.Refresh()
}

func (e *Editor) Focused() bool {
	return e.IsFocused
}

func (e *Editor) Tapped(event *fyne.PointEvent) {
	log.Println("Editor.Tapped:", event)
	e.updateCursor(event)
}

func (e *Editor) TappedSecondary(event *fyne.PointEvent) {
	log.Println("Editor.TappedSecondary:", event)
	log.Println("TODO")
	// TODO show popup menu
	//  - cut
	//  - copy
	//  - paste
	//  - select all
	//  - reveal in tree
}

func (e *Editor) DoubleTapped(event *fyne.PointEvent) {
	log.Println("Editor.DoubleTapped:", event)
	log.Println("TODO")
	// TODO select tapped text
}

func (e *Editor) MouseDown(event *desktop.MouseEvent) {
	log.Println("Editor.MouseDown:", event)
	e.updateCursor(&event.PointEvent)
}

func (e *Editor) MouseUp(event *desktop.MouseEvent) {
	log.Println("Editor.MouseUp:", event)
	// TODO
}

func (e *Editor) Dragged(event *fyne.DragEvent) {
	log.Println("Editor.Dragged:", event)
	log.Println("TODO")
	// TODO select text
}

func (e *Editor) DragEnd() {
	log.Println("Editor.DragEnd")
	log.Println("TODO")
	// TODO stop Selecting
}

func (e *Editor) KeyDown(event *fyne.KeyEvent) {
	log.Println("Editor.KeyDown:", event)
	// TODO
}

func (e *Editor) KeyUp(event *fyne.KeyEvent) {
	log.Println("Editor.KeyUp:", event)
	// TODO
	e.Refresh()
}

func (e *Editor) TypedKey(event *fyne.KeyEvent) {
	log.Println("Editor.TypedKey:", event)
	// TODO
}

func (e *Editor) TypedRune(r rune) {
	log.Println("Editor.TypedRune:", r)
	// TODO
}

func (e *Editor) TypedShortcut(shortcut fyne.Shortcut) {
	log.Println("Editor.TypedShortcut:", shortcut)
	e.shortcut.TypedShortcut(shortcut)
}

func (e *Editor) CutToClipboard(clipboard fyne.Clipboard) {
	log.Println("Editor.CutToClipboard:", clipboard)
	if !e.IsSelecting {
		return
	}
	clipboard.SetContent(e.SelectedText())
	e.EraseSelection()
}

func (e *Editor) CopyToClipboard(clipboard fyne.Clipboard) {
	log.Println("Editor.CopyToClipboard:", clipboard)
	if !e.IsSelecting {
		return
	}
	clipboard.SetContent(e.SelectedText())
}

func (e *Editor) PasteFromClipboard(clipboard fyne.Clipboard) {
	log.Println("Editor.PasteFromClipboard:", clipboard)
	// TODO
}

func (e *Editor) SelectAll() {
	log.Println("Editor.SelectAll")
	// TODO
}

func (e *Editor) EraseSelection() {
	log.Println("Editor.EraseSelection")
	// TODO
}

func (e *Editor) SelectedText() string {
	if !e.IsSelecting {
		return ""
	}
	e.Lock()
	start := e.Cursor
	end := e.Selection
	e.Unlock()
	if end < start {
		start, end = end, start
	}
	return string(e.Buffer[start:end])
}

func (e *Editor) updateCursor(event *fyne.PointEvent) {
	e.Lock()
	rowHeight := e.charMinSize().Height
	row := int(math.Floor(float64(event.Position.Y-theme.Padding()) / float64(rowHeight)))
	if row < 0 {
		row = 0
	} else if row >= len(e.Lines) {
		row = len(e.Lines) - 1
	}
	line := e.Lines[row]
	e.Cursor = uint64(line.start)
	text := e.Buffer[line.start:line.end]
	style := e.TextStyle
	size := e.TextSize
	for i := 0; i < len(text); i++ {
		width := fyne.MeasureText(string(text[0:i]), size, style).Width
		if width+theme.Padding() > event.Position.X {
			break
		} else {
			e.Cursor++
		}
	}
	e.Unlock()
	e.Refresh()
}

func (e *Editor) charMinSize() fyne.Size {
	return fyne.MeasureText("M", e.TextSize, e.TextStyle)
}

func (e *Editor) CreateRenderer() fyne.WidgetRenderer {
	cursor := canvas.NewRectangle(theme.FocusColor())
	cursor.Hide()
	return &EditorRenderer{
		editor:  e,
		cursor:  cursor,
		objects: []fyne.CanvasObject{cursor},
	}
}

type EditorRenderer struct {
	editor    *Editor
	cursor    *canvas.Rectangle
	texts     []*canvas.Text
	selection []fyne.CanvasObject
	objects   []fyne.CanvasObject
}

func (r *EditorRenderer) Layout(size fyne.Size) {
	//log.Println("EditorRenderer.Layout:", size)
	if r.cursor.Visible() {
		cursor := r.editor.Cursor
		for i, line := range r.editor.Lines {
			if uint64(line.start) <= cursor && uint64(line.end) >= cursor {
				text := string(r.editor.Buffer[line.start:cursor])
				size := fyne.MeasureText(text, r.editor.TextSize, r.editor.TextStyle)
				r.cursor.Resize(fyne.NewSize(2, size.Height))
				r.cursor.Move(fyne.NewPos(size.Width-1+theme.Padding(), size.Height*i+theme.Padding()))
				break
			}
		}
	}

	y := theme.Padding()
	rowHeight := r.editor.charMinSize().Height
	lineSize := fyne.NewSize(size.Width-theme.Padding()*2, rowHeight)
	for _, t := range r.texts {
		t.Resize(lineSize)
		t.Move(fyne.NewPos(theme.Padding(), y))
		y += rowHeight
	}
}

func (r *EditorRenderer) MinSize() (size fyne.Size) {
	charMinSize := r.editor.charMinSize()
	for i := 0; i < fyne.Min(len(r.texts), len(r.editor.Lines)); i++ {
		min := charMinSize
		if r.texts[i].Text != "" {
			min = r.texts[i].MinSize()
		}
		size.Height += min.Height
		size.Width = fyne.Max(size.Width, min.Width)
	}
	size.Width += theme.Padding() * 2
	size.Height += theme.Padding() * 2
	//log.Println("EditorRenderer.MinSize:", size)
	return
}

func (r *EditorRenderer) Refresh() {
	//log.Println("EditorRenderer.Refresh")
	r.editor.Lock()
	index := 0
	for ; index < len(r.editor.Lines); index++ {
		var textCanvas *canvas.Text
		if index < len(r.texts) {
			textCanvas = r.texts[index]
		} else {
			textCanvas = &canvas.Text{}
			r.texts = append(r.texts, textCanvas)
			r.objects = append(r.objects, textCanvas)
		}
		line := r.editor.Lines[index]
		textCanvas.Text = string(r.editor.Buffer[line.start:line.end])
		log.Println("Text:", textCanvas.Text)
		textCanvas.Show()
	}
	r.editor.Unlock()

	for ; index < len(r.texts); index++ {
		r.texts[index].Text = ""
	}

	for _, t := range r.texts {
		t.Alignment = r.editor.TextAlign
		t.Color = r.editor.TextColor
		t.TextSize = r.editor.TextSize
		t.TextStyle = r.editor.TextStyle
		t.Hidden = r.editor.Hidden
	}

	if r.editor.Focused() {
		r.cursor.Show()
	} else {
		r.cursor.Hide()
	}

	r.Layout(r.editor.Size())
	log.Println("canvas.Refresh")
	canvas.Refresh(r.editor)
	for _, t := range r.texts {
		canvas.Refresh(t)
	}
}

func (r *EditorRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *EditorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *EditorRenderer) Destroy() {
}

type Line struct {
	start, end int
}
