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

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"github.com/AletheiaWareLLC/bcfynego/ui/account"
	"github.com/AletheiaWareLLC/labfynego/ui/edit"
	"github.com/AletheiaWareLLC/labfynego/ui/experiment"
)

func main() {
	a := app.New()
	w := a.NewWindow("LAB UI")
	w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(3),
		//fyne.NewContainerWithLayout(layout.NewCenterLayout(), account.NewExportKey().CanvasObject()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), account.NewImportKey().CanvasObject()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), account.NewSignIn().CanvasObject()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), account.NewSignUp().CanvasObject()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), edit.NewEditor()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), edit.NewDeltaEditor(nil)),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), edit.NewChannelEditor(nil, nil, nil)),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), experiment.NewExperiment(nil, nil, nil, nil, nil, nil).CanvasObject()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), experiment.NewCreateExperiment(w).CanvasObject()),
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), experiment.NewJoinExperiment().CanvasObject()),
	))
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}
