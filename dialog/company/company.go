package company

import (
	"fmt"
	"strings"

	"Timelancer/model/company"
	"Timelancer/shared/tr"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogTitle = "company data"
)

type Dialog struct {
	self          *gtk.Dialog
	shortcutEntry *gtk.Entry
	nameEntry     *gtk.Entry
	usedBox       *gtk.CheckButton
	company       *company.Company
}

func New(app *gtk.Application, company *company.Company) *Dialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(app.GetActiveWindow())
		dialog.SetBorderWidth(6)
		dialog.SetTitle(dialogTitle)
		//dialog.SetSizeRequest(400, 200)

		instance := &Dialog{self: dialog, company: company}

		if contentArea, err := dialog.GetContentArea(); tr.IsOK(err) {
			if buttonBox := instance.createButtons(); buttonBox != nil {
				if separator, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
					if contentGrid := instance.createContent(); contentGrid != nil {
						contentArea.SetBorderWidth(4)
						contentArea.SetSpacing(4)

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
	//d.minSpin.GrabFocus()
}

func (d *Dialog) Run() gtk.ResponseType {
	return d.self.Run()
}

func (d *Dialog) Destroy() {
	d.self.Destroy()
}

func (d *Dialog) createButtons() *gtk.Box {
	if okBtn, err := gtk.ButtonNewWithLabel("Save"); tr.IsOK(err) {
		if cancelBtn, err := gtk.ButtonNewWithLabel("Cancel"); tr.IsOK(err) {
			if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1); tr.IsOK(err) {
				okBtn.SetTooltipText("Save data to database")
				cancelBtn.SetTooltipText("Do nothing")

				box.PackEnd(okBtn, false, true, 2)
				box.PackEnd(cancelBtn, false, true, 2)

				okBtn.Connect("clicked", func() {
					if d.company == nil {
						d.company = company.New()
					}
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

		if shortcutLabel, err := gtk.LabelNew("Shortcut:"); tr.IsOK(err) {
			if nameLabel, err := gtk.LabelNew("Name:"); tr.IsOK(err) {
				if usedLabel, err := gtk.LabelNew("In use:"); tr.IsOK(err) {
					if d.shortcutEntry, err = gtk.EntryNew(); tr.IsOK(err) {
						if d.nameEntry, err = gtk.EntryNew(); tr.IsOK(err) {
							if d.usedBox, err = gtk.CheckButtonNew(); tr.IsOK(err) {
								shortcutLabel.SetHAlign(gtk.ALIGN_END)
								nameLabel.SetHAlign(gtk.ALIGN_END)
								usedLabel.SetHAlign(gtk.ALIGN_END)
								d.shortcutEntry.SetMaxWidthChars(5)
								d.nameEntry.SetWidthChars(35)

								grid.Attach(shortcutLabel, 0, 0, 1, 1)
								grid.Attach(d.shortcutEntry, 1, 0, 1, 1)
								grid.Attach(nameLabel, 0, 1, 1, 1)
								grid.Attach(d.nameEntry, 1, 1, 1, 1)
								grid.Attach(usedLabel, 0, 2, 1, 1)
								grid.Attach(d.usedBox, 1, 2, 1, 1)

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

}

func (d *Dialog) canNotBeEmpty(name string) {
	if dialog := gtk.MessageDialogNew(d.self, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, ""); dialog != nil {
		defer dialog.Destroy()
		dialog.FormatSecondaryText(fmt.Sprintf("field '%s' can not be empty!", name))
		dialog.Run()
	}
}
