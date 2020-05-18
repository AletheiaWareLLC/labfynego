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
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego/ui"
)

type CreateExperiment struct {
	Path         *widget.Entry
	CreateButton *widget.Button
}

func NewCreateExperiment(window fyne.Window) *CreateExperiment {
	c := &CreateExperiment{
		Path: widget.NewEntry(),
		CreateButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Create Experiment",
		},
	}
	c.Path.SetPlaceHolder("Path")
	// TODO Path is single line, handle enter key by moving to button/auto click
	c.Path.ActionItem = ui.NewFilePicker(window, c.Path)
	return c
}

func (c *CreateExperiment) CanvasObject() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		c.Path,
		layout.NewSpacer(),
		c.CreateButton,
	)
}
