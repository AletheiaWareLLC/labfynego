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
)

type JoinExperiment struct {
	Host       *widget.Entry
	ID         *widget.Entry
	JoinButton *widget.Button
}

func NewJoinExperiment() *JoinExperiment {
	j := &JoinExperiment{
		Host: widget.NewEntry(),
		ID:   widget.NewEntry(),
		JoinButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Join Experiment",
		},
	}
	j.Host.SetPlaceHolder("Host")
	j.ID.SetPlaceHolder("ID")
	// TODO Host is single line, handle enter key by moving to id
	// TODO ID is single line, handle enter key by moving to button/auto click
	return j
}

func (j *JoinExperiment) CanvasObject() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		j.Host,
		j.ID,
		layout.NewSpacer(),
		j.JoinButton,
	)
}
