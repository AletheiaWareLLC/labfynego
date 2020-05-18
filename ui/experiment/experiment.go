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

package experiment

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego/ui"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/labfynego/ui/edit"
	"github.com/AletheiaWareLLC/labgo"
	"log"
	"os"
	"strings"
)

type Experiment struct {
	Node       *bcgo.Node
	Listener   bcgo.MiningListener
	Cache      bcgo.Cache
	Network    bcgo.Network
	Experiment *labgo.Experiment
	Window     fyne.Window

	Chat    *widget.Label
	Editors map[string]*edit.ChannelEditor
	Items   map[string]*widget.TabItem
	Status  *widget.Label
	Tabber  *widget.TabContainer
	Tree    fyne.CanvasObject
}

func NewExperiment(node *bcgo.Node, listener bcgo.MiningListener, cache bcgo.Cache, network bcgo.Network, experiment *labgo.Experiment, window fyne.Window) *Experiment {
	e := &Experiment{
		Node:       node,
		Listener:   listener,
		Cache:      cache,
		Network:    network,
		Experiment: experiment,
		Window:     window,
		Tabber:     widget.NewTabContainer(),
		Items:      make(map[string]*widget.TabItem),
		Editors:    make(map[string]*edit.ChannelEditor),
		Chat:       widget.NewLabel("Chat"),
		Status:     widget.NewLabel("Ready"),
	}
	var channel *bcgo.Channel
	if experiment != nil {
		channel = experiment.Path
	}
	e.Tree = edit.NewTree(channel, cache, network, e.SelectPath)
	return e
}

func (e *Experiment) GetOrOpenDeltaChannel(fileId string) *bcgo.Channel {
	channel, err := e.Node.GetChannel(labgo.LAB_PREFIX_FILE + fileId)
	if err != nil {
		log.Println(err)
		channel = labgo.OpenFileChannel(fileId)
		// Load channel
		if err := channel.LoadCachedHead(e.Cache); err != nil {
			log.Println(err)
		}
		if e.Network != nil {
			// Pull channel from network
			if err := channel.Pull(e.Cache, e.Network); err != nil {
				log.Println(err)
			}
		}
		// Add channel to node
		e.Node.AddChannel(channel)
	}
	return channel
}

func (e *Experiment) SelectPath(id string, path ...string) {
	log.Println("Selected:", id, path)
	go func() {
		editor, ok := e.Editors[id]
		if !ok {
			editor = edit.NewChannelEditor(e.Node, e.Listener, e.GetOrOpenDeltaChannel(id))
			e.Editors[id] = editor
		}
		item, ok := e.Items[id]
		if !ok {
			name := id
			if len(path) > 0 {
				name = path[len(path)-1]
			}
			item = widget.NewTabItem(name, widget.NewVScrollContainer(editor))
			e.Items[id] = item
			e.Tabber.Append(item)
		}
		e.Tabber.SelectTab(item)
		if len(e.Items) == 1 {
			// First tab, resize tabber
			e.Tabber.Resize(e.Tabber.MinSize())
		}
	}()
}

func (e *Experiment) CanvasObject() fyne.CanvasObject {
	left := widget.NewVScrollContainer(e.Tree)
	center := e.Tabber
	right := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, e.Status, nil, nil), e.Status, e.Chat)
	splitter := widget.NewHSplitContainer(left, center)
	splitter.Offset = 0.25
	//splitter = widget.NewVSplitContainer(splitter, bottom)
	//splitter.Offset = 0.75
	splitter = widget.NewHSplitContainer(splitter, right)
	splitter.Offset = 0.75
	return splitter
}

func (e *Experiment) MainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New", func() {
				fmt.Println("Menu File->New")
				filepath := widget.NewEntry()
				filepath.SetPlaceHolder("/path/to/new/file")
				dialog.ShowCustomConfirm("New File", "Create", "Cancel", filepath, func(b bool) {
					if !b {
						return
					}
					path := filepath.Text
					log.Println(path)
					id, _, err := labgo.CreatePath(e.Node, e.Listener, e.Experiment.Path, strings.Split(path, string(os.PathSeparator)))
					if err != nil {
						dialog.ShowError(err, e.Window)
						return
					}
					e.SelectPath(id, path)
				}, e.Window)
			}),
			fyne.NewMenuItem("Import", func() {
				fmt.Println("Menu File->Import")
				dialog.ShowFileOpen(func(reader fyne.FileReadCloser, err error) {
					if err != nil {
						dialog.ShowError(err, e.Window)
						return
					}
					uri := reader.URI().String()
					// TODO truncate uri to remove file:///Users/foobar/...
					path := uri
					log.Println(path)
					id, _, err := labgo.CreatePathFromReader(e.Node, e.Listener, e.Experiment.Path, strings.Split(path, string(os.PathSeparator)), reader)
					if err != nil {
						dialog.ShowError(err, e.Window)
						return
					}
					e.SelectPath(id, path)
				}, e.Window)
			}),
			fyne.NewMenuItem("Export", func() {
				fmt.Println("Menu File->Export")
				dialog.ShowFileSave(func(writer fyne.FileWriteCloser, err error) {
					if err != nil {
						dialog.ShowError(err, e.Window)
						return
					}
					log.Println("Exporting", writer.URI(), "Not Yet Supported")
					/* TODO how to write paths to directories?
					if err := labgo.Save(e.Node, e.Experiment, writer); err != nil {
						dialog.ShowError(err, e.Window)
						return
					}
					*/
				}, e.Window)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Settings", func() {
				fmt.Println("Menu Settings")
			})),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Cut", func() {
				ui.ShortcutFocused(&fyne.ShortcutCut{
					Clipboard: e.Window.Clipboard(),
				}, e.Window)
			}),
			fyne.NewMenuItem("Copy", func() {
				ui.ShortcutFocused(&fyne.ShortcutCopy{
					Clipboard: e.Window.Clipboard(),
				}, e.Window)
			}),
			fyne.NewMenuItem("Paste", func() {
				ui.ShortcutFocused(&fyne.ShortcutPaste{
					Clipboard: e.Window.Clipboard(),
				}, e.Window)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Find", func() {
				fmt.Println("Menu Find")
			})),
		fyne.NewMenu("Help", fyne.NewMenuItem("Help", func() {
			fmt.Println("Help Menu")
		})),
	)
}
