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
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/labfynego/ui/edit"
	"github.com/AletheiaWareLLC/labgo"
	"log"
	"os"
)

func main() {
	// Load config files (if any)
	err := bcgo.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config: %w", err)
	}

	// Get root directory
	rootDir, err := bcgo.GetRootDirectory()
	if err != nil {
		log.Fatal("Could not get root directory: %w", err)
	}

	// Get cache directory
	cacheDir, err := bcgo.GetCacheDirectory(rootDir)
	if err != nil {
		log.Fatal("Could not get cache directory: %w", err)
	}

	// Create file cache
	cache, err := bcgo.NewFileCache(cacheDir)
	if err != nil {
		log.Fatal("Could not create file cache: %w", err)
	}

	// Create network of peers
	network := bcgo.NewTCPNetwork()

	// Create node
	node, err := bcgo.GetNode(rootDir, cache, network)
	if err != nil {
		log.Fatal("Could not create node: %w", err)
	}

	// Create listener
	listener := &bcgo.PrintingMiningListener{Output: os.Stdout}

	// Create channel
	channel := bcgo.OpenPoWChannel(labgo.LAB_PREFIX_FILE+"LabTest", labgo.CHANNEL_THRESHOLD)

	// Load channel
	if err := channel.LoadCachedHead(node.Cache); err != nil {
		log.Println("Could not load head from cache:", err)
	}
	if node.Network != nil {
		// Pull channel from network
		if err := channel.Pull(node.Cache, node.Network); err != nil {
			log.Println("Could not load head from network:", err)
		}
	}
	// Add channel to node
	node.AddChannel(channel)

	// Create application
	app := app.New()

	// Create editor
	editor := edit.NewChannelEditor(node, listener, channel)

	// Create window
	window := app.NewWindow("LAB")
	window.SetContent(editor)
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()
	window.ShowAndRun()
}
