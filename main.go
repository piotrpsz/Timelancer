/*
 * BSD 2-Clause License
 *
 *	Copyright (c) 2019, Piotr Pszczółkowski
 *	All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"os"
	"path/filepath"

	"Timelancer/dbf"
	"Timelancer/shared"
	"Timelancer/shared/tr"
	"Timelancer/window"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	appID = "pl.beesoft.gtk3.Timelancer"
)

func main() {
	if openDatabase() {
		if app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE); tr.IsOK(err) {
			app.Connect("activate", func() {
				tr.Init()
				if win := window.New(app); win != nil {
					quitAction := glib.SimpleActionNew("quit", nil)
					quitAction.Connect("activate", func() {
						tr.Cancel()
						app.Quit()
					})
					app.AddAction(quitAction)

					win.ShowAll()
				}
			})
			retv := app.Run(os.Args)
			os.Exit(retv)
		}
	}
	os.Exit(1)
}

func openDatabase() bool {
	if dataDir := shared.AppDir(); dataDir != "" {
		filePath := filepath.Join(dataDir, shared.AppName+".sqlite")
		return dbf.OpenOrCreate(filePath)
	}
	return false
}
