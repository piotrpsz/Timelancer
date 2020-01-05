/*
 * BSD 2-Clause License
 *
 * Copyright (c) 2019, Piotr Pszczółkowski
 * All rights reserved.
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

package companies

import (
	"Timelancer/model/company"
	"Timelancer/shared/tr"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogTitle      = "all companies table"
	cancelBtnText    = "cancel"
	addBtnText       = "new"
	editBtnText      = "edit"
	deleteBtnText    = "remove"
	cancelBtnTooltip = "close this dialog"
	addBtnTooltip    = "add new company"
	editBtnTooltip   = "edit selected company"
	deleteBtnTooltip = "remove selected company"

	idColumnIdx       = 0
	shortcutColumnIdx = 1
	nameColumnIdx     = 2
	useColumnIdx      = 3
)

type Dialog struct {
	self      *gtk.Dialog
	cancelBtn *gtk.Button
	addBtn    *gtk.Button
	editBtn   *gtk.Button
	deleteBtn *gtk.Button
	treeView  *gtk.TreeView
	listStore *gtk.ListStore
}

func New(app *gtk.Application) *Dialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(app.GetActiveWindow())
		dialog.SetBorderWidth(6)
		dialog.SetTitle(dialogTitle)
		//dialog.SetSizeRequest(400, 200)

		instance := &Dialog{self: dialog}

		if contentArea, err := dialog.GetContentArea(); tr.IsOK(err) {
			if buttonBox := instance.createButtons(); buttonBox != nil {
				if separator, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
					if instance.createTable() {

						contentArea.PackEnd(buttonBox, false, false, 1)
						contentArea.PackEnd(separator, true, false, 1)
						contentArea.PackEnd(instance.treeView, true, true, 1)

						return instance
					}
				}
			}
		}
	}
	return nil
}

func (d *Dialog) ShowAll() {
	d.self.ShowAll()
	//d.self.SetResizable(false)
}

func (d *Dialog) Run() gtk.ResponseType {
	return d.self.Run()
}

func (d *Dialog) Destroy() {
	d.self.Destroy()
}

func (d *Dialog) UpdateTable() {
	d.listStore.Clear()
	if companiesData := company.CompaniesInUse(); len(companiesData) > 0 {
		for _, c := range companiesData {
			iter := d.listStore.Append()
			d.listStore.SetValue(iter, idColumnIdx, c.ID())
			d.listStore.SetValue(iter, shortcutColumnIdx, c.Shortcut())
			d.listStore.SetValue(iter, nameColumnIdx, c.Name())
		}
	}
}

func (d *Dialog) createButtons() *gtk.Box {
	var err error

	if d.cancelBtn, err = gtk.ButtonNewWithLabel(cancelBtnText); tr.IsOK(err) {
		if d.addBtn, err = gtk.ButtonNewWithLabel(addBtnText); tr.IsOK(err) {
			if d.editBtn, err = gtk.ButtonNewWithLabel(editBtnText); tr.IsOK(err) {
				if d.deleteBtn, err = gtk.ButtonNewWithLabel(deleteBtnText); tr.IsOK(err) {
					if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1); tr.IsOK(err) {
						d.cancelBtn.SetTooltipText(cancelBtnTooltip)
						d.addBtn.SetTooltipText(addBtnTooltip)
						d.editBtn.SetTooltipText(editBtnTooltip)
						d.deleteBtn.SetTooltipText(deleteBtnTooltip)

						box.PackEnd(d.cancelBtn, false, false, 2)
						box.PackEnd(d.addBtn, false, false, 2)
						box.PackEnd(d.editBtn, false, false, 2)
						box.PackEnd(d.deleteBtn, false, false, 2)

						d.cancelBtn.Connect("clicked", func() {
							d.self.Response(gtk.RESPONSE_OK)
						})

						return box
					}
				}
			}
		}
	}
	return nil
}

/********************************************************************
*                                                                   *
*                             T A B L E                             *
*                                                                   *
********************************************************************/

func (d *Dialog) createTable() bool {
	if treeView, listStore := setupTreeView(); treeView != nil {
		d.treeView = treeView
		d.listStore = listStore
		return true
	}
	return false
}

func setupTreeView() (*gtk.TreeView, *gtk.ListStore) {
	if treeView, err := gtk.TreeViewNew(); tr.IsOK(err) {
		if idColumn := createColumn("id", idColumnIdx); idColumn != nil {
			if shortcutColumn := createColumn("shortcut", shortcutColumnIdx); shortcutColumn != nil {
				if nameColumn := createColumn("name", nameColumnIdx); nameColumn != nil {
					if useColumn := createColumn("in use", useColumnIdx); useColumn != nil {
						treeView.AppendColumn(idColumn)
						treeView.AppendColumn(shortcutColumn)
						treeView.AppendColumn(nameColumn)
						treeView.AppendColumn(useColumn)

						if listStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING); tr.IsOK(err) {
							treeView.SetModel(listStore)
							return treeView, listStore
						}

					}
				}
			}
		}
	}
	return nil, nil
}

func createColumn(title string, idx int) *gtk.TreeViewColumn {
	if cellRenderer, err := gtk.CellRendererTextNew(); tr.IsOK(err) {
		if column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", idx); tr.IsOK(err) {
			return column
		}
	}
	return nil
}
