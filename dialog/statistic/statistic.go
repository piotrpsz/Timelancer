package statistic

import (
	"fmt"

	"Timelancer/model/company"
	"Timelancer/shared/tr"
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

	d.self.ShowAll()
	d.self.SetResizable(false)
}

func (d *Dialog) Run() gtk.ResponseType {
	return d.self.Run()
}

func (d *Dialog) Destroy() {
	d.self.Destroy()
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

				box.PackStart(d.periodLabel, false, false, 2)
				box.PackStart(d.periodComboBox, true, false, 2)

				return box
			}
		}
	}
	return nil
}

func (d *Dialog) populateCompanyComboBox() {
	d.companyComboBox.RemoveAll()
	d.companyComboBox.AppendText("All")

	companies := company.CompaniesInUse()
	for _, c := range companies {
		d.companyComboBox.AppendText(c.Name())
	}
	d.companyComboBox.SetActive(0)
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
