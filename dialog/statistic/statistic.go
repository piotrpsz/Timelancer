package statistic

import (
	"Timelancer/shared/tr"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogTitle   = "working time statistic"
	cancelBtnText = "cancel"
	exportBtnText = "export"

	idColumnIdx       = 0
	shortcutColumnIdx = 1
	nameColumnIdx     = 2
	startColumnIdx    = 3
	stopColumnIdx     = 4
)

type Dialog struct {
	self         *gtk.Dialog
	parent       *gtk.Window
	companiesBox *gtk.ComboBoxText
	startBox     *gtk.ComboBoxText
	endBox       *gtk.ComboBoxText
}

func New(parent *gtk.Window) *Dialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(parent)
		dialog.SetBorderWidth(6)
		dialog.SetTitle(dialogTitle)
		dialog.SetSizeRequest(400, 200)

		instance := &Dialog{self: dialog, parent: parent}
		return instance

		//if contentArea, err := dialog.GetContentArea(); tr.IsOK(err) {
		//	if buttonBox := instance.createButtons(); buttonBox != nil {
		//		if separator, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); tr.IsOK(err) {
		//			if instance.createTable() {
		//
		//				contentArea.PackEnd(buttonBox, false, false, 1)
		//				contentArea.PackEnd(separator, true, false, 1)
		//				contentArea.PackEnd(instance.treeView, true, true, 1)
		//
		//				return instance
		//			}
		//		}
		//	}
		//}
	}
	return nil
}
