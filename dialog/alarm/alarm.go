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

package alarm

import (
	"time"

	"Timelancer/shared"
	"Timelancer/shared/tr"
	"github.com/gotk3/gotk3/gtk"
)

const (
	dialogAfterTitle = "alarm after duration"
	dialogAtTitle    = "alarm at"
	hoursAfter       = "Hours"
	hoursAt          = "Hour"
	minutesAfter     = "Minutes"
	minutesAt        = "Minute"
	secondsAfter     = "Seconds"
	secondsAt        = "Second"
)

type AlarmDialog struct {
	self     *gtk.Dialog
	hourSpin *gtk.SpinButton
	minSpin  *gtk.SpinButton
	secSpin  *gtk.SpinButton
	after    bool
}

func New(app *gtk.Application, after bool) *AlarmDialog {
	if dialog, err := gtk.DialogNew(); tr.IsOK(err) {
		dialog.SetTransientFor(app.GetActiveWindow())
		if after {
			dialog.SetTitle(dialogAfterTitle)
		} else {
			dialog.SetTitle(dialogAtTitle)
		}
		dialog.SetBorderWidth(6)
		instance := &AlarmDialog{self: dialog, after:after}

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
	}
	return nil
}

func (d *AlarmDialog) ShowAll() {
	d.self.ShowAll()
	d.self.SetResizable(false)
	d.minSpin.GrabFocus()
}

func (d *AlarmDialog) Run() gtk.ResponseType {
	return d.self.Run()
}

func (d *AlarmDialog) Destroy() {
	d.self.Destroy()
}

func (d *AlarmDialog) Selection() (uint, uint, uint) {
	h := d.hourSpin.GetValueAsInt()
	m := d.minSpin.GetValueAsInt()
	s := d.secSpin.GetValueAsInt()
	return uint(h), uint(m), uint(s)
}

func (d *AlarmDialog) createButtons() *gtk.Box {
	if okBtn, err := gtk.ButtonNewWithLabel("OK"); tr.IsOK(err) {
		if cancelBtn, err := gtk.ButtonNewWithLabel("Cancel"); tr.IsOK(err) {
			if box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1); tr.IsOK(err) {
				box.PackEnd(okBtn, false, true, 2)
				box.PackEnd(cancelBtn, false, true, 2)

				okBtn.Connect("clicked", func() {
					d.self.Response(gtk.RESPONSE_OK)
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

func (d *AlarmDialog) createContent() *gtk.Grid {
	if grid, err := gtk.GridNew(); tr.IsOK(err) {
		grid.SetBorderWidth(8)
		grid.SetRowSpacing(8)
		grid.SetColumnSpacing(8)

		if hoursLabel, err := gtk.LabelNew(""); tr.IsOK(err) {
			if minutesLabel, err := gtk.LabelNew(""); tr.IsOK(err) {
				if secondsLabel, err := gtk.LabelNew(""); tr.IsOK(err) {
					if d.hourSpin = createHourSpin(); d.hourSpin != nil {
						if d.minSpin = createMinSpin(); d.minSpin != nil {
							if d.secSpin = createSecSpin(); d.secSpin != nil {
								if d.after {
									hoursLabel.SetText(hoursAfter)
									minutesLabel.SetText(minutesAfter)
									secondsLabel.SetText(secondsAfter)
								} else {
									_, _, _, h, m, s := shared.DateTimeComponents(time.Now())
									d.hourSpin.SetValue(float64(h))
									d.minSpin.SetValue(float64(m))
									d.secSpin.SetValue(float64(s))

									hoursLabel.SetText(hoursAt)
									minutesLabel.SetText(minutesAt)
									secondsLabel.SetText(secondsAt)
								}
								grid.Attach(hoursLabel,   0, 0, 1, 1)
								grid.Attach(minutesLabel, 1, 0, 1, 1)
								grid.Attach(secondsLabel, 2, 0, 1, 1)
								grid.Attach(d.hourSpin,   0, 1, 1, 1)
								grid.Attach(d.minSpin,    1, 1, 1, 1)
								grid.Attach(d.secSpin,    2, 1, 1, 1)

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

func createHourSpin() *gtk.SpinButton {
	if adjustment, err := gtk.AdjustmentNew(0.0, 0.0, 24.0, 1.0, 1.0, 0.0); tr.IsOK(err) {
		if spin, err := gtk.SpinButtonNew(adjustment, 1.0, 0); tr.IsOK(err) {
			return spin
		}
	}
	return nil
}

func createMinSpin() *gtk.SpinButton {
	if adjustment, err := gtk.AdjustmentNew(0.0, 0.0, 60.0, 1.0, 10.0, 0.0); tr.IsOK(err) {
		if spin, err := gtk.SpinButtonNew(adjustment, 1.0, 0); tr.IsOK(err) {
			return spin
		}
	}
	return nil
}

func createSecSpin() *gtk.SpinButton {
	if adjustment, err := gtk.AdjustmentNew(0.0, 0.0, 60.0, 1.0, 10.0, 0.0); tr.IsOK(err) {
		if spin, err := gtk.SpinButtonNew(adjustment, 1.0, 0); tr.IsOK(err) {
			return spin
		}
	}
	return nil
}
