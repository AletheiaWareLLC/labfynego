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
	"fyne.io/fyne/theme"
	"github.com/AletheiaWareLLC/labgo"
	"log"
)

type DeltaEditor struct {
	Editor

	Deltas map[string]*labgo.Delta
	Order  []string

	OnDelta func(string, *labgo.Delta)
}

func NewDeltaEditor(callback func(parent string, delta *labgo.Delta)) *DeltaEditor {
	e := &DeltaEditor{
		Editor: Editor{
			TextAlign: fyne.TextAlignLeading,
			TextColor: theme.TextColor(),
			TextSize:  theme.TextSize(),
			TextStyle: fyne.TextStyle{},
			TextWrap:  fyne.TextWrapWord,
		},
		OnDelta: callback,
		Deltas:  make(map[string]*labgo.Delta),
	}
	e.ExtendBaseWidget(e)
	e.AddShortcuts()
	return e
}

func (e *DeltaEditor) TypedRune(r rune) {
	log.Println("DeltaEditor.TypedRune:", r)
	// TODO add runes to list until timeout or cursor is moved elsewhere, then create file delta
	var parentRecordId string
	e.Lock()
	delta := &labgo.Delta{
		Offset: e.Cursor,
		Add:    []byte(string(r)),
	}
	if len(e.Order) > 0 {
		parentRecordId = e.Order[len(e.Order)-1]
	}
	e.Unlock()
	if e.OnDelta != nil {
		e.OnDelta(parentRecordId, delta)
	}
}

func (e *DeltaEditor) PasteFromClipboard(clipboard fyne.Clipboard) {
	log.Println("DeltaEditor.PasteFromClipboard:", clipboard)
	var parentRecordId string
	e.Lock()
	delta := &labgo.Delta{
		Offset: e.Cursor,
		Add:    []byte(clipboard.Content()),
	}
	if e.IsSelecting {
		if e.Selection < delta.Offset {
			delta.Offset = e.Selection
		}
		delta.Remove = []byte(e.SelectedText())
		e.IsSelecting = false
	}
	if len(e.Order) > 0 {
		parentRecordId = e.Order[len(e.Order)-1]
	}
	e.Unlock()
	if e.OnDelta != nil {
		e.OnDelta(parentRecordId, delta)
	}
}

func (e *DeltaEditor) EraseSelection() {
	log.Println("DeltaEditor.EraseSelection")
	/* TODO
	if e.OnDelta != nil {
		e.OnDelta
	}
	*/
}
