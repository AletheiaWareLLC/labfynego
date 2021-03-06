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

package labfynego

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego"
	bcui "github.com/AletheiaWareLLC/bcfynego/ui"
	bcdata "github.com/AletheiaWareLLC/bcfynego/ui/data"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/labclientgo"
	"github.com/AletheiaWareLLC/labfynego/ui/data"
	"github.com/AletheiaWareLLC/labfynego/ui/experiment"
	"github.com/AletheiaWareLLC/labgo"
	"log"
	"os"
)

type LabFyneClient struct {
	bcfynego.BCFyneClient
	labclientgo.LabClient
	Experiment *labgo.Experiment
}

func (c *LabFyneClient) GetExperiment() *labgo.Experiment {
	if c.Experiment == nil {
		ec := make(chan *labgo.Experiment, 1)
		go c.ShowExperimentDialog(func(e *labgo.Experiment) {
			if node := c.GetNode(); node != nil {
				if net, ok := node.Network.(*bcgo.TCPNetwork); ok {
					go labgo.Serve(node, node.Cache, net)
				}
			}
			ec <- e
		})
		c.Experiment = <-ec
		c.Window.SetTitle("LAB - " + c.Experiment.ID)
	}
	return c.Experiment
}

func (c *LabFyneClient) GetLogo() fyne.CanvasObject {
	return &canvas.Image{
		Resource: bcdata.NewThemedResource(data.LogoUnmasked),
		//FillMode: canvas.ImageFillContain,
		FillMode: canvas.ImageFillOriginal,
	}
}

func (c *LabFyneClient) ShowExperiment(n *bcgo.Node, e *labgo.Experiment) {
	log.Println("ShowExperiment")
	ui := experiment.NewExperiment(
		n,
		&bcgo.PrintingMiningListener{Output: os.Stdout},
		n.Cache,
		n.Network,
		e,
		c.Window)
	c.Window.SetContent(ui.CanvasObject())
	c.Window.SetMainMenu(ui.MainMenu())
}

func (c *LabFyneClient) ShowExperimentDialog(callback func(*labgo.Experiment)) {
	log.Println("ShowExperimentDialog")
	create := experiment.NewCreateExperiment(c.Window)
	join := experiment.NewJoinExperiment()

	c.Dialog = dialog.NewCustom("Experiment Access", "Cancel",
		widget.NewAccordionContainer(
			&widget.AccordionItem{Title: "Create", Detail: create.CanvasObject(), Open: true},
			widget.NewAccordionItem("Join", join.CanvasObject()),
		), c.Window)

	create.CreateButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Create Tapped")
		uri := create.Path.Text
		go func() {
			var reader fyne.FileReadCloser
			if uri != "" {
				r, err := storage.OpenFileFromURI(storage.NewURI(uri))
				if err != nil {
					dialog.ShowError(err, c.Window)
					return
				}
				defer r.Close()
				reader = r
			}
			progress := dialog.NewProgress("Creating", "message", c.Window)
			defer progress.Hide()
			listener := &bcui.ProgressMiningListener{Func: progress.SetValue}
			experiment, err := labgo.CreateFromReader(c.GetNode(), listener, uri, reader)
			if err != nil {
				dialog.ShowError(err, c.Window)
				return
			}
			// TODO callback should add experiment.Path to channels
			callback(experiment)
		}()
	}
	join.JoinButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Join Tapped")
		host := join.Host.Text
		id := join.ID.Text
		go func() {
			// Create channel
			p := labgo.OpenPathChannel(id)
			cache, err := c.GetCache()
			if err != nil {
				log.Println(err)
			} else {
				// Load channel
				if err := p.LoadCachedHead(cache); err != nil {
					log.Println(err)
				}
				net, err := c.GetNetwork()
				if err != nil {
					log.Println(err)
				} else {
					// Connect to host
					if host != "" && host != "localhost" {
						if n, ok := net.(*bcgo.TCPNetwork); ok {
							n.Connect(host, []byte("test"))
						}
					}
					// Pull channel from network
					if err := p.Pull(cache, net); err != nil {
						log.Println(err)
					}
				}
			}
			// TODO callback should add experiment.Path to channels
			callback(&labgo.Experiment{
				ID:   id,
				Path: p,
			})
		}()
	}
	c.Dialog.Show()
}

/*
func (c *LabFyneClient) ShowExperimentList(fn func(i int, b binding.Binding)) {
		experimentItems := binding.NewStringList()
		go func() {
			// TODO read experiments from chain
		}()
		var experimentCells int
		experimentList := &widget.List{
			Items: experimentItems,
			OnCreateCell: func() fyne.CanvasObject {
				experimentCells++
				log.Println("Created Label Cell:", experimentCells)
				return &widget.Label{
					Wrapping: fyne.TextWrapBreak,
				}
			},
			OnBindCell: func(object fyne.CanvasObject, data binding.Binding) {
				t, ok := object.(*widget.Label)
				if ok {
					s, ok := data.(binding.String)
					if ok {
						t.Text = s.Get()
					}
					t.Show()
				}
			},
			OnSelected: func(i int, b binding.Binding) {
				log.Println("Selected:", i, b)
				if fn != nil {
					fn(i, b)
				}
			},
		}
		c.Window.SetContent(experimentList)
}
*/
