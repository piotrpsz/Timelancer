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
	"strconv"

	"Timelancer/dialog/company"
	companyData "Timelancer/model/company"

	"Timelancer/shared/tr"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogTitle      = "all companies table"
	cancelBtnText    = "cancel"
	addBtnText       = "add new"
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
	self        *gtk.Dialog
	cancelBtn   *gtk.Button
	addBtn      *gtk.Button
	editBtn     *gtk.Button
	deleteBtn   *gtk.Button
	scroll      *gtk.ScrolledWindow
	treeView    *gtk.TreeView
	listStore   *gtk.ListStore
	parent      *gtk.Window
	selectedRow int
}

func New(parent *gtk.Window) *Dialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(parent)
		dialog.SetBorderWidth(6)
		dialog.SetTitle(dialogTitle)
		//dialog.SetSizeRequest(400, 200)

		instance := &Dialog{self: dialog, parent: parent, selectedRow: -1}

		if contentArea, err := dialog.GetContentArea(); tr.IsOK(err) {
			if buttonBox := instance.createButtons(); buttonBox != nil {
				if separator, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
					if instance.createTable() {

						contentArea.PackEnd(buttonBox, false, false, 1)
						contentArea.PackEnd(separator, true, false, 1)
						contentArea.PackEnd(instance.scroll, true, true, 1)

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
	d.self.SetResizable(false)
}

func (d *Dialog) Run() gtk.ResponseType {
	return d.self.Run()
}

func (d *Dialog) Destroy() {
	d.self.Destroy()
}

func (d *Dialog) UpdateTable() {
	d.listStore.Clear()
	if companiesData := companyData.Companies(); len(companiesData) > 0 {
		for _, c := range companiesData {
			d.updateDataAtIter(d.listStore.Append(), c)
		}
	}
	d.treeView.GrabFocus()
	d.updateButtonStates()
}

func (d *Dialog) updateButtonStates() {
	if _, ok := d.listStore.GetIterFirst(); ok {
		d.deleteBtn.SetSensitive(true)
		d.editBtn.SetSensitive(true)
		return
	}
	d.deleteBtn.SetSensitive(false)
	d.editBtn.SetSensitive(false)
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
						d.addBtn.Connect("clicked", d.addActionHandler)
						d.editBtn.Connect("clicked", d.editActionHandler)
						d.deleteBtn.Connect("clicked", d.deleteActionHandler)

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
*                B U T T O N   H A N D L E R S                      *
*                                                                   *
********************************************************************/

func (d *Dialog) addActionHandler() {
	if dialog := company.New(&d.self.Window, nil); dialog != nil {
		defer dialog.Destroy()

		dialog.ShowAll()
		if dialog.Run() == gtk.RESPONSE_OK {
			if c := dialog.Company(); c != nil && c.Valid() {
				if c.Save() {
					d.UpdateTable()
					d.selectRowWithID(c.ID())
					return
				}
			}
			d.saveFailure()
		}
	}
}

func (d *Dialog) editActionHandler() {
	if c := d.selectedCompany(); c != nil {
		if dialog := company.New(&d.self.Window, c); dialog != nil {
			defer dialog.Destroy()

			dialog.ShowAll()
			if dialog.Run() == gtk.RESPONSE_OK {
				if c := dialog.Company(); c != nil && c.Valid() {
					if c.Save() {
						d.updateDataInSelectedRow(c)
						return
					}
				}
				d.saveFailure()
			}
		}
	}
}

/// Remove selected in table company from database and update table.
func (d *Dialog) deleteActionHandler() {
	// TODO: check if it is possible (maybe company was alrady used)
	if iter := d.currentSelectionIter(); iter != nil {
		if c := d.companyAtIter(iter); c != nil {
			if c.Remove() {
				d.listStore.Remove(iter)
			}
		}
	}
}

func (d *Dialog) saveFailure() {
	if dialog := gtk.MessageDialogNew(&d.self.Window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, "error"); dialog != nil {
		defer dialog.Destroy()
		dialog.FormatSecondaryText("can't save company data to database.")
		dialog.Run()
	}
}

/********************************************************************
*                                                                   *
*                             T A B L E                             *
*                                                                   *
********************************************************************/

func (d *Dialog) createTable() bool {
	if scroll, err := gtk.ScrolledWindowNew(nil, nil); tr.IsOK(err) {
		if treeView, listStore := d.setupTreeView(); treeView != nil {
			d.scroll = scroll
			d.treeView = treeView
			d.listStore = listStore

			d.scroll.SetSizeRequest(500, 250)
			d.scroll.Add(d.treeView)
			return true
		}
	}
	return false
}

func (d *Dialog) setupTreeView() (*gtk.TreeView, *gtk.ListStore) {
	if treeView, err := gtk.TreeViewNew(); tr.IsOK(err) {
		if idColumn := d.createTextColumn("id", idColumnIdx); idColumn != nil {
			if shortcutColumn := d.createTextColumn("shortcut", shortcutColumnIdx); shortcutColumn != nil {
				if nameColumn := d.createTextColumn("name", nameColumnIdx); nameColumn != nil {
					if useColumn := d.createToggleColumn("in use", useColumnIdx); useColumn != nil {
						idColumn.SetVisible(false)

						treeView.AppendColumn(idColumn)
						treeView.AppendColumn(shortcutColumn)
						treeView.AppendColumn(nameColumn)
						treeView.AppendColumn(useColumn)
						treeView.ColumnsAutosize()

						if listStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_BOOLEAN); tr.IsOK(err) {
							treeView.SetModel(listStore)

							if selection, err := treeView.GetSelection(); tr.IsOK(err) {
								selection.SetMode(gtk.SELECTION_SINGLE)
								selection.Connect("changed", d.selectionChanged)
								return treeView, listStore
							}
						}

					}
				}
			}
		}
	}
	return nil, nil
}

func (d *Dialog) createTextColumn(title string, idx int) *gtk.TreeViewColumn {
	if cellRenderer, err := gtk.CellRendererTextNew(); tr.IsOK(err) {
		if column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", idx); tr.IsOK(err) {
			column.SetResizable(true)
			return column
		}
	}
	return nil
}

func (d *Dialog) createToggleColumn(title string, idx int) *gtk.TreeViewColumn {
	if renderer, err := gtk.CellRendererToggleNew(); tr.IsOK(err) {
		renderer.SetActivatable(true)
		renderer.Connect("toggled", func(p *gtk.CellRendererToggle, rowAsString string) {
			if row, err := strconv.Atoi(rowAsString); tr.IsOK(err) {
				if path, err := gtk.TreePathNewFromIndicesv([]int{row}); tr.IsOK(err) {
					if iter, err := d.listStore.GetIter(path); tr.IsOK(err) {
						if id, ok := d.getID(iter); ok {
							if use, ok := d.getUse(iter); ok {
								if c := companyData.CompanyWithID(id); c != nil {
									c.SetUsed(!use)
									if c.Save() {
										d.listStore.SetValue(iter, useColumnIdx, !use)
									}
								}
							}
						}
					}
				}
			}
		})
		if column, err := gtk.TreeViewColumnNewWithAttribute(title, renderer, "active", idx); tr.IsOK(err) {
			return column
		}
	}
	return nil
}

func (d *Dialog) selectionChanged(s *gtk.TreeSelection) {
	if rows := s.GetSelectedRows(d.listStore); rows != nil {
		if path, ok := rows.Data().(*gtk.TreePath); ok {
			if indices := path.GetIndices(); len(indices) > 0 {
				d.selectedRow = indices[0]
			}
		}
	}
}

func (d *Dialog) currentSelectionIter() *gtk.TreeIter {
	if selection, err := d.treeView.GetSelection(); tr.IsOK(err) {
		if _, iter, ok := selection.GetSelected(); ok {
			return iter
		}
	}
	return nil
}

func (d *Dialog) selectedCompany() *companyData.Company {
	if iter := d.currentSelectionIter(); iter != nil {
		return d.companyAtIter(iter)
	}
	return nil
}

func (d *Dialog) companyAtIter(iter *gtk.TreeIter) *companyData.Company {
	if id, ok := d.getID(iter); ok {
		return companyData.CompanyWithID(id)
	}
	return nil
}

func (d *Dialog) getID(iter *gtk.TreeIter) (int, bool) {
	if value, err := d.listStore.GetValue(iter, idColumnIdx); tr.IsOK(err) {
		if idValue, err := value.GoValue(); tr.IsOK(err) {
			if id, ok := idValue.(int); ok {
				return id, true
			}
		}
	}
	return -1, false
}

func (d *Dialog) getUse(iter *gtk.TreeIter) (bool, bool) {
	if value, err := d.listStore.GetValue(iter, useColumnIdx); tr.IsOK(err) {
		if useValue, err := value.GoValue(); tr.IsOK(err) {
			if use, ok := useValue.(bool); ok {
				return use, true
			}
		}
	}
	return false, false
}

func (d *Dialog) iterForID(id int) *gtk.TreeIter {
	if iter, ok := d.listStore.GetIterFirst(); ok {
		if v, ok := d.getID(iter); ok && v == id {
			return iter
		}
		for {
			if d.listStore.IterNext(iter) {
				if v, ok := d.getID(iter); ok && v == id {
					return iter
				}
			}
		}
	}
	return nil
}

func (d *Dialog) selectRowWithID(id int) {
	if iter := d.iterForID(id); iter != nil {
		if selection, err := d.treeView.GetSelection(); tr.IsOK(err) {
			selection.SelectIter(iter)
			d.scrollToIter(iter)
		}
	}
}

func (d *Dialog) scrollToIter(iter *gtk.TreeIter) {
	if model, err := d.treeView.GetModel(); tr.IsOK(err) {
		if path, err := model.GetPath(iter); tr.IsOK(err) {
			if column := d.treeView.GetColumn(nameColumnIdx); column != nil {
				d.treeView.ScrollToCell(path, column, false, 0.0, 0.0)
			}
		}
	}
}

func (d *Dialog) updateDataInSelectedRow(c *companyData.Company) {
	d.updateDataAtIter(d.currentSelectionIter(), c)
}

func (d *Dialog) updateDataAtIter(iter *gtk.TreeIter, c *companyData.Company) {
	if iter != nil && c != nil {
		d.listStore.SetValue(iter, idColumnIdx, c.ID())
		d.listStore.SetValue(iter, shortcutColumnIdx, c.Shortcut())
		d.listStore.SetValue(iter, nameColumnIdx, c.Name())
		d.listStore.SetValue(iter, useColumnIdx, c.Used())
	}
}
