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

package company

import (
	"fmt"
	"strings"

	"Timelancer/model/company"
	"Timelancer/shared/tr"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogTitle       = "company data"
	shortcutLabelText = "shortcut:"
	nameLabelText     = "name:"
	inUseLabelText    = "is use:"
	saveBtnText       = "save"
	cancelBtnText     = "cancel"
	saveTooltip       = "save data to database"
	cancelTooltip     = "do nothing"
)

type Dialog struct {
	self          *gtk.Dialog
	shortcutLabel *gtk.Label
	nameLabel     *gtk.Label
	usedLabel     *gtk.Label
	shortcutEntry *gtk.Entry
	nameEntry     *gtk.Entry
	usedBox       *gtk.CheckButton
	company       *company.Company
}

func New(app *gtk.Application, c *company.Company) *Dialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(app.GetActiveWindow())
		dialog.SetBorderWidth(6)
		dialog.SetTitle(dialogTitle)
		//dialog.SetSizeRequest(400, 200)

		instance := &Dialog{self: dialog, company: c}

		if contentArea, err := dialog.GetContentArea(); tr.IsOK(err) {
			if buttonBox := instance.createButtons(); buttonBox != nil {
				if separator, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
					if contentGrid := instance.createContent(); contentGrid != nil {
						contentArea.SetBorderWidth(4)
						contentArea.SetSpacing(4)

						if instance.company == nil {
							instance.company = company.New()
							instance.usedBox.SetActive(true)
						}

						contentArea.PackEnd(buttonBox, false, false, 0)
						contentArea.PackEnd(separator, true, true, 1)
						contentArea.PackEnd(contentGrid, false, false, 0)
						return instance
					}
				}
			}
		}

		return instance
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

func (d *Dialog) Company() *company.Company {
	return d.company
}

func (d *Dialog) createButtons() *gtk.Box {
	if okBtn, err := gtk.ButtonNewWithLabel(saveBtnText); tr.IsOK(err) {
		if cancelBtn, err := gtk.ButtonNewWithLabel(cancelBtnText); tr.IsOK(err) {
			if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1); tr.IsOK(err) {
				okBtn.SetTooltipText(saveTooltip)
				cancelBtn.SetTooltipText(cancelTooltip)

				box.PackEnd(okBtn, false, true, 2)
				box.PackEnd(cancelBtn, false, true, 2)

				okBtn.Connect("clicked", func() {
					if d.widgetsToCompany() {
						d.self.Response(gtk.RESPONSE_OK)
					}
				})
				cancelBtn.Connect("clicked", func() {
					d.self.Response(gtk.RESPONSE_CANCEL)
				})

				return box
			}
		}
	}
	return nil
}

func (d *Dialog) createContent() *gtk.Grid {
	if grid, err := gtk.GridNew(); tr.IsOK(err) {
		grid.SetBorderWidth(8)
		grid.SetRowSpacing(8)
		grid.SetColumnSpacing(8)

		var err error
		if d.shortcutLabel, err = gtk.LabelNew(shortcutLabelText); tr.IsOK(err) {
			if d.nameLabel, err = gtk.LabelNew(nameLabelText); tr.IsOK(err) {
				if d.usedLabel, err = gtk.LabelNew(inUseLabelText); tr.IsOK(err) {
					if d.shortcutEntry, err = gtk.EntryNew(); tr.IsOK(err) {
						if d.nameEntry, err = gtk.EntryNew(); tr.IsOK(err) {
							if d.usedBox, err = gtk.CheckButtonNew(); tr.IsOK(err) {
								d.shortcutLabel.SetHAlign(gtk.ALIGN_END)
								d.nameLabel.SetHAlign(gtk.ALIGN_END)
								d.usedLabel.SetHAlign(gtk.ALIGN_END)
								d.shortcutEntry.SetMaxWidthChars(5)
								d.nameEntry.SetWidthChars(35)
								d.usedBox.SetCanFocus(false)

								grid.Attach(d.shortcutLabel, 0, 0, 1, 1)
								grid.Attach(d.shortcutEntry, 1, 0, 1, 1)
								grid.Attach(d.nameLabel, 0, 1, 1, 1)
								grid.Attach(d.nameEntry, 1, 1, 1, 1)
								grid.Attach(d.usedLabel, 0, 2, 1, 1)
								grid.Attach(d.usedBox, 1, 2, 1, 1)

								d.usedBox.Connect("toggled", d.updateFocus)

								return grid
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (d *Dialog) widgetsToCompany() bool {
	if shortcut, err := d.shortcutEntry.GetText(); tr.IsOK(err) {
		if strings.TrimSpace(shortcut) == "" {
			d.canNotBeEmpty("shortcut")
			d.shortcutEntry.GrabFocus()
			return false
		}
		if name, err := d.nameEntry.GetText(); tr.IsOK(err) {
			if strings.TrimSpace(name) == "" {
				d.canNotBeEmpty("name")
				d.nameEntry.GrabFocus()
				return false
			}
			d.company.SetShortcut(shortcut)
			d.company.SetName(name)
			d.company.SetUsed(d.usedBox.GetActive())
			return true
		}
	}
	return false
}

func (d *Dialog) companyToWidgets() {
	d.shortcutEntry.SetText(d.company.Shortcut())
	d.nameEntry.SetText(d.company.Name())
	d.usedBox.SetActive(d.company.Used())
	d.updateFocus()
}

func (d *Dialog) canNotBeEmpty(name string) {
	if dialog := gtk.MessageDialogNew(d.self, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, ""); dialog != nil {
		defer dialog.Destroy()
		dialog.FormatSecondaryText(fmt.Sprintf("field '%s' can not be empty!", name))
		dialog.Run()
	}
}

func (d *Dialog) updateFocus() {
	if d.usedBox.GetActive() {
		d.shortcutLabel.SetSensitive(true)
		d.nameLabel.SetSensitive(true)
		d.shortcutEntry.SetSensitive(true)
		d.nameEntry.SetSensitive(true)
		d.shortcutEntry.GrabFocus()
	} else {
		d.shortcutLabel.SetSensitive(false)
		d.nameLabel.SetSensitive(false)
		d.shortcutEntry.SetSensitive(false)
		d.nameEntry.SetSensitive(false)
	}
}
