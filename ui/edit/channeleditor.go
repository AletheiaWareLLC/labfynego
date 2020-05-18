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
	"encoding/base64"
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/labgo"
	"github.com/golang/protobuf/proto"
	"log"
	"sort"
)

type ChannelEditor struct {
	DeltaEditor
	Node     *bcgo.Node
	Listener bcgo.MiningListener
	Channel  *bcgo.Channel
	Entries  map[string]*bcgo.BlockEntry
}

func NewChannelEditor(node *bcgo.Node, listener bcgo.MiningListener, channel *bcgo.Channel) *ChannelEditor {
	e := &ChannelEditor{
		DeltaEditor: DeltaEditor{
			Editor: Editor{
				//TextAlign: fyne.TextAlignTrailing,
				TextColor: theme.TextColor(),
				TextSize:  theme.TextSize(),
				TextStyle: fyne.TextStyle{},
				TextWrap:  fyne.TextWrapWord,
			},
			Deltas: make(map[string]*labgo.Delta),
		},
		Node:     node,
		Listener: listener,
		Channel:  channel,
		Entries:  make(map[string]*bcgo.BlockEntry),
	}
	e.ExtendBaseWidget(e)
	e.AddShortcuts()
	if channel != nil {
		channel.AddTrigger(e.Read)
		defer e.Read()
	}
	e.OnDelta = e.Write
	return e
}

func (e *ChannelEditor) Read() {
	log.Println("Read")
	e.Lock()
	if err := bcgo.Read(e.Channel.Name, e.Channel.Head, nil, e.Node.Cache, e.Node.Network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		id := base64.RawURLEncoding.EncodeToString(entry.RecordHash)
		e.Entries[id] = entry
		// Unmarshal as Delta
		delta := &labgo.Delta{}
		if err := proto.Unmarshal(data, delta); err != nil {
			return err
		}
		_, ok := e.Deltas[id]
		if ok {
			return nil
		}
		e.Deltas[id] = delta
		e.Entries[id] = entry
		e.Order = append(e.Order, id)
		sort.Slice(e.Order, func(i, j int) bool {
			return e.Entries[e.Order[i]].Record.Timestamp < e.Entries[e.Order[j]].Record.Timestamp //.Record.Creator
		})
		buffer := []byte{}
		for _, id := range e.Order {
			delta := e.Deltas[id]
			log.Println("Edit:", id, e.Entries[id].Record.Creator, delta)
			buffer = labgo.DeltaToBuffer(delta, buffer)
			if e.Entries[id].Record.Creator == e.Node.Alias {
				e.Cursor = delta.Offset + uint64(len(delta.Add))
			}
		}
		e.Buffer = []rune(string(buffer))
		log.Println("Buffer:", string(e.Buffer))
		return nil
	}); err != nil {
		log.Println(err)
	}
	e.Unlock()
	e.Refresh()
}

func (e *ChannelEditor) Write(parent string, delta *labgo.Delta) {
	log.Println("Write:", parent, delta)
	// Create protobuf record
	hash, record, err := labgo.ProtoToRecord(e.Node.Alias, e.Node.Key, bcgo.Timestamp(), delta)
	if err != nil {
		log.Println(err)
		return
	}
	//h := base64.RawURLEncoding.EncodeToString(hash)

	// Create entries
	entries := []*bcgo.BlockEntry{
		&bcgo.BlockEntry{
			RecordHash: hash,
			Record:     record,
		},
	}
	// Mine Channel
	if _, _, err := e.Node.MineEntries(e.Channel, labgo.CHANNEL_THRESHOLD, e.Listener, entries); err != nil {
		log.Println(err)
		return
	}

	if e.Node.Network != nil {
		// Push Update to Peers
		if err := e.Channel.Push(e.Node.Cache, e.Node.Network); err != nil {
			log.Println(err)
			return
		}
	}
}
