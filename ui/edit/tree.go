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
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/labgo"
	"github.com/golang/protobuf/proto"
	"log"
)

func NewTree(paths *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, callback func(id string, path ...string)) fyne.CanvasObject {
	tree := widget.NewVBox()
	if paths != nil {
		trigger := func() {
			var objects []fyne.CanvasObject
			if err := bcgo.Read(paths.Name, paths.Head, nil, cache, network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
				id := base64.RawURLEncoding.EncodeToString(entry.RecordHash)
				// Unmarshal as Path
				p := &labgo.Path{}
				if err := proto.Unmarshal(data, p); err != nil {
					return err
				}
				name := id
				if len(p.Path) > 0 {
					name = p.Path[len(p.Path)-1]
				}
				objects = append(objects, &widget.Button{
					Text: name,
					OnTapped: func() {
						callback(id, p.Path...)
					},
				})
				return nil
			}); err != nil {
				log.Println(err)
			}
			tree.Children = objects
			tree.Refresh()
		}
		paths.AddTrigger(trigger)
		trigger()
	}
	return tree
}
