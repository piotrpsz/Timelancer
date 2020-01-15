/*
 * BSD 2-Clause License
 *
 *	Copyright (c) 2019, Piotr Pszczółkowski (beesoft software)
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

package statistic

import (
	"fmt"
	"time"

	"Timelancer/model/company"
	"Timelancer/shared"
	"Timelancer/shared/tr"
	"Timelancer/sqlite"
	"Timelancer/sqlite/row"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogTitle      = "working time statistic"
	companyLabelText = "company:"
	companyTooltip   = "companies you work for"
	periodLabelText  = "period:"
	periodTooltip    = "predefined periods of time"
	//startLabelText     = "start date"
	//endLabelText       = "end date"
	cancelBtnText    = "return"
	exportBtnText    = "export"
	cancelBtnTooltip = "close this dialog"
	exportBtnTooltip = "save records to csv file"

	idColumnIdx      = 0
	idColumnName     = "id"
	nameColumnIdx    = 1
	nameColumnName   = "name"
	startColumnIdx   = 2
	startColumnName  = "start"
	finishColumnIdx  = 3
	finishColumnName = "finish"
	periodColumnIdx  = 4
	perionColumnName = "period"
)

var (
	periods = []string{"all", "today", "yesterday", "this week", "previous week", "this month", "previous month", "this year", "previous year"}
)

type Dialog struct {
	self            *gtk.Dialog
	parent          *gtk.Window
	companyLabel    *gtk.Label
	companyComboBox *gtk.ComboBoxText
	periodLabel     *gtk.Label
	periodComboBox  *gtk.ComboBoxText
	cancelBtn       *gtk.Button
	exportBtn       *gtk.Button
	treeView        *gtk.TreeView
	listStore       *gtk.ListStore

	ids []int
}

func New(parent *gtk.Window) *Dialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(parent)
		dialog.SetBorderWidth(6)
		dialog.SetTitle(dialogTitle)
		dialog.SetSizeRequest(400, 200)

		instance := &Dialog{self: dialog, parent: parent}

		if contentArea, err := dialog.GetContentArea(); tr.IsOK(err) {
			if buttonBox := instance.createButtons(); buttonBox != nil {
				if separatorBottom, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
					if scroll := instance.createTable(); scroll != nil {
						if separatorTop, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
							if toolbarGrid := instance.createToolbar(); toolbarGrid != nil {

								contentArea.PackEnd(buttonBox, false, false, 1)
								contentArea.PackEnd(separatorBottom, true, false, 1)
								contentArea.PackEnd(scroll, true, true, 1)
								contentArea.PackEnd(separatorTop, true, false, 1)
								contentArea.PackEnd(toolbarGrid, true, true, 1)

								return instance
							}
						}
					}
				}
			}
		}

	}
	return nil
}

func (d *Dialog) ShowAll() {
	d.populateCompanyComboBox()
	d.populatePeriodComboBox()
	d.DidSelectAllCompanies()

	d.self.ShowAll()
	d.self.SetResizable(false)
}

func (d *Dialog) Run() gtk.ResponseType {
	return d.self.Run()
}

func (d *Dialog) Destroy() {
	d.self.Destroy()
}

func (d *Dialog) DidSelectAllCompanies() {
	query := "SELECT timer.id, timer.start, timer.finish, company.name FROM timer,company WHERE timer.company_id=company.id ORDER BY timer.id DESC"
	d.updateTable(query)
}

func (d *Dialog) DidSelectecCompanyWithID(id int) {
	tr.Info("id: %d", id)
	query := fmt.Sprintf("SELECT timer.id, timer.start, timer.finish, company.name FROM timer,company WHERE timer.company_id=%d AND timer.company_id=company.id ORDER BY timer.id DESC", id)
	d.updateTable(query)
}

func (d *Dialog) updateTable(query string) {
	d.listStore.Clear()

	sqlite.SQLite().SelectAndHandle(query, func(r row.Row) {
		//fmt.Printf("%+v\n", r)
		if iter := d.listStore.Append(); iter != nil {
			if id, ok := getID(r); ok {
				if name, ok := getName(r); ok {
					if start, ok := getStart(r); ok {
						if finish, ok := getFinish(r); ok {
							d.listStore.SetValue(iter, idColumnIdx, id)
							d.listStore.SetValue(iter, nameColumnIdx, name)
							d.listStore.SetValue(iter, startColumnIdx, shared.TimeAsString(start))
							d.listStore.SetValue(iter, finishColumnIdx, shared.TimeAsString(finish))
							d.listStore.SetValue(iter, periodColumnIdx, getPeriod(start, finish))
						}
					}
				}
			}
		}
	})
}

func getID(r row.Row) (int64, bool) {
	if id, ok := r["id"]; ok {
		if id, ok := id.Value.(int64); ok {
			return id, true
		}
	}
	return -1, false
}
func getName(r row.Row) (string, bool) {
	if name, ok := r["name"]; ok {
		if name, ok := name.Value.(string); ok {
			return name, true
		}
	}
	return "", false
}
func getStart(r row.Row) (time.Time, bool) {
	if start, ok := r["start"]; ok {
		if start, ok := start.Value.(int64); ok {
			if t := time.Unix(start, 0); !t.IsZero() {
				return t, true
			}
		}
	}
	return time.Time{}, false
}
func getFinish(r row.Row) (time.Time, bool) {
	if start, ok := r["finish"]; ok {
		if start, ok := start.Value.(int64); ok {
			if t := time.Unix(start, 0); !t.IsZero() {
				return t, true
			}
		}
	}
	return time.Time{}, false
}
func getPeriod(start, finish time.Time) string {
	duration := finish.Sub(start)
	seconds := uint(duration.Seconds())
	h, m, s := shared.DurationComponents(seconds)
	if s >= 30 {
		m += 1
		if m > 60 {
			h += 1
			m -= 60
		}
	}

	return fmt.Sprintf("%dh %02dmin", h, m)
}

func (d *Dialog) createButtons() *gtk.Box {
	var err error

	if d.cancelBtn, err = gtk.ButtonNewWithLabel(cancelBtnText); tr.IsOK(err) {
		if d.exportBtn, err = gtk.ButtonNewWithLabel(exportBtnText); tr.IsOK(err) {
			if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1); tr.IsOK(err) {
				d.cancelBtn.SetTooltipText(cancelBtnTooltip)
				d.exportBtn.SetTooltipText(exportBtnTooltip)

				box.PackEnd(d.cancelBtn, false, false, 2)
				box.PackEnd(d.exportBtn, false, false, 2)

				d.cancelBtn.Connect("clicked", func() {
					d.self.Response(gtk.RESPONSE_OK)
				})
				d.exportBtn.Connect("clicked", func() {
					fmt.Println("export data")
				})

				return box
			}
		}
	}
	return nil
}

func (d *Dialog) selectedCompanyChanged() {
	if row := d.companyComboBox.GetActive(); row > -1 {
		if row < len(d.ids) {
			if id := d.ids[row]; id == -1 {
				d.DidSelectAllCompanies()
			} else {
				d.DidSelectecCompanyWithID(id)
			}
		}
	}
}

func (d *Dialog) selectedPersonChanged() {
	fmt.Println("selectedPersonCahnged")
}

func (d *Dialog) createToolbar() *gtk.Grid {
	if grid, err := gtk.GridNew(); tr.IsOK(err) {
		if companiesBox := d.createCompanyBox(); companiesBox != nil {
			if periodBox := d.createPeriodBox(); periodBox != nil {
				grid.SetColumnSpacing(10)
				grid.Attach(companiesBox, 0, 0, 1, 1)
				grid.Attach(periodBox, 1, 0, 1, 1)

				return grid
			}
		}
	}
	return nil
}

func (d *Dialog) createCompanyBox() *gtk.Box {
	var err error

	if d.companyLabel, err = gtk.LabelNew(companyLabelText); tr.IsOK(err) {
		if d.companyComboBox, err = gtk.ComboBoxTextNew(); tr.IsOK(err) {
			if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2); tr.IsOK(err) {
				d.companyComboBox.SetTooltipText(companyTooltip)
				d.companyComboBox.Connect("changed", d.selectedCompanyChanged)

				box.PackStart(d.companyLabel, false, false, 2)
				box.PackStart(d.companyComboBox, true, false, 2)

				return box
			}
		}
	}
	return nil
}

func (d *Dialog) createPeriodBox() *gtk.Box {
	var err error

	if d.periodLabel, err = gtk.LabelNew(periodLabelText); tr.IsOK(err) {
		if d.periodComboBox, err = gtk.ComboBoxTextNew(); tr.IsOK(err) {
			if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2); tr.IsOK(err) {
				d.periodComboBox.SetTooltipText(periodTooltip)
				d.periodComboBox.Connect("changed", d.selectedPersonChanged)

				box.PackStart(d.periodLabel, false, false, 2)
				box.PackStart(d.periodComboBox, true, false, 2)

				return box
			}
		}
	}
	return nil
}

func (d *Dialog) populateCompanyComboBox() {
	var ids []int

	d.companyComboBox.RemoveAll()
	d.companyComboBox.AppendText("All")
	ids = append(ids, -1)

	companies := company.CompaniesInUse()
	for _, c := range companies {
		d.companyComboBox.AppendText(c.Name())
		ids = append(ids, c.ID())
	}
	d.companyComboBox.SetActive(0)
	d.ids = ids
}

func (d *Dialog) populatePeriodComboBox() {
	d.periodComboBox.RemoveAll()
	for _, text := range periods {
		d.periodComboBox.AppendText(text)
	}
	d.periodComboBox.SetActive(0)
}

func (d *Dialog) createTable() *gtk.ScrolledWindow {
	if scroll, err := gtk.ScrolledWindowNew(nil, nil); tr.IsOK(err) {
		if treeView, err := gtk.TreeViewNew(); tr.IsOK(err) {
			if appendColumns(treeView) {
				if store, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING); tr.IsOK(err) {
					treeView.SetModel(store)
					if selection, err := treeView.GetSelection(); tr.IsOK(err) {
						selection.SetMode(gtk.SELECTION_SINGLE)

						d.treeView = treeView
						d.listStore = store

						scroll.SetSizeRequest(500, 250)
						scroll.Add(d.treeView)
						return scroll
					}
				}
			}
		}
	}
	return nil
}

func appendColumns(treeView *gtk.TreeView) bool {
	if idColumn := createTextColumn(idColumnName, idColumnIdx); idColumn != nil {
		if nameColumn := createTextColumn(nameColumnName, nameColumnIdx); nameColumn != nil {
			if startColumn := createTextColumn(startColumnName, startColumnIdx); startColumn != nil {
				if finishColumn := createTextColumn(finishColumnName, finishColumnIdx); finishColumn != nil {
					if periodColumn := createTextColumn(perionColumnName, periodColumnIdx); periodColumn != nil {
						idColumn.SetVisible(false)

						treeView.AppendColumn(idColumn)
						treeView.AppendColumn(nameColumn)
						treeView.AppendColumn(startColumn)
						treeView.AppendColumn(finishColumn)
						treeView.AppendColumn(periodColumn)
						treeView.ColumnsAutosize()

						return true
					}
				}
			}
		}
	}
	return false
}

func createTextColumn(title string, idx int) *gtk.TreeViewColumn {
	if renderer, err := gtk.CellRendererTextNew(); tr.IsOK(err) {
		if column, err := gtk.TreeViewColumnNewWithAttribute(title, renderer, "text", idx); tr.IsOK(err) {
			column.SetResizable(true)
			return column
		}
	}
	return nil
}
