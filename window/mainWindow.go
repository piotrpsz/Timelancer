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

package window

import (
	"context"
	"fmt"
	"sync"
	"time"

	"Timelancer/dialog/alarm"
	"Timelancer/dialog/companies"
	"Timelancer/dialog/company"
	"Timelancer/dialog/statistic"
	"Timelancer/model/timer"
	"Timelancer/shared"
	"Timelancer/shared/tr"
	"Timelancer/sound"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	licence = `BSD 2-Clause License

Copyright (c) 2019, Piotr Pszczółkowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

	dateFormat = "%04d-%02d-%02d  %s"
	timeFormat = "<span font_desc='16' foreground='#AAA555'>%02d:%02d</span>" +
		"<span font_desc='10' foreground='#AAA555'>.%02d</span>"

	workTimeActiveFormat   = "<span font_desc='18' foreground='#AAA555'>%02d:%02d:%02d</span>"
	workTimeInactiveFormat = "<span font_desc='18' foreground='#999999'>%02d:%02d:%02d</span>"

	alarmAfterActiveFormat   = "<span font_desc='18' foreground='#AAA555'>%02d:%02d:%02d</span>"
	alarmAfterInactiveFormat = "<span font_desc='18' foreground='#999999'>%02d:%02d:%02d</span>"

	// you worked %dh%2dmin.\nWould you like to save this information to database?

	workedTimeFormat = "<span font_desc='12' foreground='#BBBBBBBB'>you worked </span>" +
		"<span font_desc='16' foreground='#CCC777'> %d</span><span font_desc='10' foreground='#BBBBBBBB'>h </span>" +
		"<span font_desc='16' foreground='#CCC777'>%02d</span><span font_desc='10' foreground='#BBBBBBBB'>min </span>" +
		"<span font_desc='12' foreground='#BBBBBBBB'>\nwould you like to save this information to database?</span>"
)

type MainWindow struct {
	app                *gtk.Application
	win                *gtk.ApplicationWindow
	timeLabel          *gtk.Label
	headerBar          *gtk.HeaderBar
	companyLabel       *gtk.Label
	companyCombo       *gtk.ComboBoxText
	companyAddBtn      *gtk.Button
	timerLabel         *gtk.Label
	timerValue         *gtk.Label
	timerStartBtn      *gtk.Button
	timerStopBtn       *gtk.Button
	alarmAfterLabel    *gtk.Label
	alarmAfterValue    *gtk.Label
	alarmAfterStartBtn *gtk.Button
	alarmAfterSetBtn   *gtk.Button
	alarmAfterStopBtn  *gtk.Button
	alarmAtLabel       *gtk.Label
	alarmAtValue       *gtk.Label
	alarmAtStartBtn    *gtk.Button
	alarmAtSetBtn      *gtk.Button
	alarmAtStopBtn     *gtk.Button

	wg     sync.WaitGroup
	cancel context.CancelFunc

	lastTime              time.Time
	workTimeStart         time.Time
	workTimeRunned        bool
	alarmAfterDuration    uint
	alarmAfterDurationPrv uint
	alarmAfterRunned      bool
	alarmAt               time.Time
	alarmAtRunned         bool
	companyIndex          int
}

func New(app *gtk.Application) *MainWindow {
	if win, err := gtk.ApplicationWindowNew(app); tr.IsOK(err) {
		mw := &MainWindow{app: app, win: win}
		if mw.setupHeaderBar() && mw.setupMenu() && mw.setupContent() {
			ctx, cancel := context.WithCancel(context.Background())
			mw.cancel = cancel
			ticker := time.NewTicker(1 * time.Second)

			mw.updateCurrentTime(time.Now())
			mw.updateWorkTime(uint(0))
			mw.resetAlarmAfter()
			mw.resetAlarmAt()
			mw.selectedCompanyChanged()

			mw.companyCombo.Connect("changed", mw.selectedCompanyChanged)

			mw.wg.Add(1)
			go mw.timeHandler(ctx, &mw.wg, ticker)

			win.SetPosition(gtk.WIN_POS_CENTER)
			//win.SetDefaultSize(400, 200)
			return mw
		}
	}
	return nil
}

func (mw *MainWindow) ShowAll() {
	mw.win.ShowAll()
}

func (mw *MainWindow) setupHeaderBar() bool {
	if headerBar, err := gtk.HeaderBarNew(); tr.IsOK(err) {
		if timeLabel, err := gtk.LabelNew(""); tr.IsOK(err) {
			headerBar.SetShowCloseButton(false)
			headerBar.SetTitle(shared.AppName)
			headerBar.PackStart(timeLabel)

			mw.headerBar = headerBar
			mw.timeLabel = timeLabel
			mw.win.SetTitlebar(headerBar)
			return true
		}
	}
	return false
}

func (mw *MainWindow) setupMenu() bool {
	if menuButton, err := gtk.MenuButtonNew(); tr.IsOK(err) {
		if menu := glib.MenuNew(); menu != nil {
			menu.Append("companies...", "custom.companies")
			menu.Append("working time statistic...", "custom.statistic")
			menu.Append("settings...", "custom.settings")
			menu.Append("about...", "custom.about")
			menu.Append("quit", "custom.quit")

			companiesAction := glib.SimpleActionNew("companies", nil)
			companiesAction.Connect("activate", mw.companiesActionHandler)

			statisticAction := glib.SimpleActionNew("statistic", nil)
			statisticAction.Connect("activate", mw.statisticActionHandler)

			settingsAction := glib.SimpleActionNew("settings", nil)
			settingsAction.Connect("activate", func() {
				fmt.Println("Settings...")
			})
			aboutAction := glib.SimpleActionNew("about", nil)
			aboutAction.Connect("activate", mw.aboutActionHandler)

			quitAction := glib.SimpleActionNew("quit", nil)
			quitAction.Connect("activate", func() {
				mw.cancel()
				mw.wg.Wait()
				mw.app.Quit()
			})

			customGroup := glib.SimpleActionGroupNew()
			customGroup.AddAction(companiesAction)
			customGroup.AddAction(statisticAction)
			customGroup.AddAction(settingsAction)
			customGroup.AddAction(aboutAction)
			customGroup.AddAction(quitAction)
			mw.win.InsertActionGroup("custom", customGroup)

			menuButton.SetMenuModel(&menu.MenuModel)
			mw.headerBar.PackEnd(menuButton)
			return true
		}
	}
	return false
}

func (mw *MainWindow) setupContent() bool {
	if grid, err := gtk.GridNew(); tr.IsOK(err) {
		grid.SetBorderWidth(8)
		grid.SetRowSpacing(8)
		grid.SetColumnSpacing(8)

		if mw.createCompanyWidgets(grid) {
			if mw.createTimerWidgets(grid) {
				if mw.createAlarmAfterWidgets(grid) {
					if mw.createAlarmAtWidgets(grid) {
						mw.win.Container.Add(grid)
						return true
					}
				}
			}
		}
	}
	return false
}

func (mw *MainWindow) createCompanyWidgets(grid *gtk.Grid) bool {
	var err error

	if mw.companyLabel, err = gtk.LabelNew("Company:"); tr.IsOK(err) {
		if mw.companyCombo, err = gtk.ComboBoxTextNew(); tr.IsOK(err) {
			if mw.companyAddBtn, err = gtk.ButtonNewWithLabel("Add"); tr.IsOK(err) {
				mw.companyLabel.SetHAlign(gtk.ALIGN_END)
				mw.companyAddBtn.SetTooltipText("Add new company")

				mw.companyAddBtn.Connect("clicked", mw.addCompanyHandler)

				grid.Attach(mw.companyLabel, 0, 0, 1, 1)
				grid.Attach(mw.companyCombo, 1, 0, 3, 1)
				grid.Attach(mw.companyAddBtn, 4, 0, 1, 1)

				mw.populateCompanyCombo()
				return true
			}
		}
	}
	return false
}

func (mw *MainWindow) createTimerWidgets(grid *gtk.Grid) bool {
	var err error

	if mw.timerLabel, err = gtk.LabelNew("Working time:"); tr.IsOK(err) {
		if mw.timerValue, err = gtk.LabelNew(""); tr.IsOK(err) {
			if mw.timerStartBtn, err = gtk.ButtonNewWithLabel("Start"); tr.IsOK(err) {
				if mw.timerStopBtn, err = gtk.ButtonNewWithLabel("Stop"); tr.IsOK(err) {
					mw.timerLabel.SetSensitive(false)
					mw.timerStopBtn.SetSensitive(false)
					mw.timerLabel.SetHAlign(gtk.ALIGN_END)
					mw.timerValue.SetHAlign(gtk.ALIGN_START)
					mw.timerStartBtn.SetTooltipText("Start of work")
					mw.timerStopBtn.SetTooltipText("End of work")

					mw.timerStartBtn.Connect("clicked", func() {
						mw.companyCombo.SetSensitive(false)
						mw.companyAddBtn.SetSensitive(false)
						mw.timerLabel.SetSensitive(true)
						mw.timerStopBtn.SetSensitive(true)
						mw.timerStartBtn.SetSensitive(false)
						mw.workTimeStart = time.Now()
						mw.workTimeRunned = true
					})
					mw.timerStopBtn.Connect("clicked", func() {
						mw.workTimeRunned = false
						mw.saveTimerAfterStopIfNeeded()

						mw.companyCombo.SetSensitive(true)
						mw.companyAddBtn.SetSensitive(true)
						mw.timerLabel.SetSensitive(false)
						mw.timerStopBtn.SetSensitive(false)
						mw.timerStartBtn.SetSensitive(true)
						mw.updateWorkTime(uint(0))
					})

					grid.Attach(mw.timerLabel, 0, 1, 1, 1)
					grid.Attach(mw.timerValue, 1, 1, 1, 1)
					grid.Attach(mw.timerStartBtn, 2, 1, 1, 1)
					grid.Attach(mw.timerStopBtn, 3, 1, 1, 1)

					return true
				}
			}
		}
	}
	return false
}

func (mw *MainWindow) createAlarmAfterWidgets(grid *gtk.Grid) bool {
	var err error

	if mw.alarmAfterLabel, err = gtk.LabelNew("Alarm after:"); tr.IsOK(err) {
		if mw.alarmAfterValue, err = gtk.LabelNew(""); tr.IsOK(err) {
			if mw.alarmAfterStartBtn, err = gtk.ButtonNewWithLabel("Start"); tr.IsOK(err) {
				if mw.alarmAfterSetBtn, err = gtk.ButtonNewWithLabel("Set"); tr.IsOK(err) {
					if mw.alarmAfterStopBtn, err = gtk.ButtonNewWithLabel("Stop"); tr.IsOK(err) {
						mw.alarmAfterLabel.SetSensitive(false)
						mw.alarmAfterStartBtn.SetSensitive(false)
						mw.alarmAfterStopBtn.SetSensitive(false)

						mw.alarmAfterLabel.SetHAlign(gtk.ALIGN_END)
						mw.alarmAfterValue.SetHAlign(gtk.ALIGN_START)
						mw.alarmAfterStartBtn.SetTooltipText("Start the timer")
						mw.alarmAfterSetBtn.SetTooltipText("Setup the timer")
						mw.alarmAfterStopBtn.SetTooltipText("Stop the timer")

						mw.alarmAfterSetBtn.Connect("clicked", mw.alarmAfterSetHandler)
						mw.alarmAfterStartBtn.Connect("clicked", mw.alarmAfterStartHandler)
						mw.alarmAfterStopBtn.Connect("clicked", mw.alarmAfterStopHandler)

						separator, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
						grid.Attach(separator, 0, 2, 5, 1)

						grid.Attach(mw.alarmAfterLabel, 0, 3, 1, 1)
						grid.Attach(mw.alarmAfterValue, 1, 3, 1, 1)
						grid.Attach(mw.alarmAfterStartBtn, 2, 3, 1, 1)
						grid.Attach(mw.alarmAfterSetBtn, 3, 3, 1, 1)
						grid.Attach(mw.alarmAfterStopBtn, 4, 3, 1, 1)

						return true
					}
				}
			}
		}
	}
	return false
}

func (mw *MainWindow) createAlarmAtWidgets(grid *gtk.Grid) bool {
	var err error

	if mw.alarmAtLabel, err = gtk.LabelNew("Alarm at:"); tr.IsOK(err) {
		if mw.alarmAtValue, err = gtk.LabelNew(""); tr.IsOK(err) {
			if mw.alarmAtStartBtn, err = gtk.ButtonNewWithLabel("Start"); tr.IsOK(err) {
				if mw.alarmAtSetBtn, err = gtk.ButtonNewWithLabel("Set"); tr.IsOK(err) {
					if mw.alarmAtStopBtn, err = gtk.ButtonNewWithLabel("Stop"); tr.IsOK(err) {
						mw.alarmAtLabel.SetSensitive(false)
						mw.alarmAtStopBtn.SetSensitive(false)
						mw.alarmAtStartBtn.SetSensitive(false)

						mw.alarmAtLabel.SetHAlign(gtk.ALIGN_END)
						mw.alarmAtValue.SetHAlign(gtk.ALIGN_START)
						mw.alarmAtStartBtn.SetTooltipText("Start the timer")
						mw.alarmAtSetBtn.SetTooltipText("Setup the timer")
						mw.alarmAtStopBtn.SetTooltipText("Stop the timer")

						mw.alarmAtSetBtn.Connect("clicked", mw.alarmAtSetHandler)
						mw.alarmAtStartBtn.Connect("clicked", mw.alarmAtStartHandler)
						mw.alarmAtStopBtn.Connect("clicked", mw.alarmAtStopHandler)

						grid.Attach(mw.alarmAtLabel, 0, 4, 1, 1)
						grid.Attach(mw.alarmAtValue, 1, 4, 1, 1)
						grid.Attach(mw.alarmAtStartBtn, 2, 4, 1, 1)
						grid.Attach(mw.alarmAtSetBtn, 3, 4, 1, 1)
						grid.Attach(mw.alarmAtStopBtn, 4, 4, 1, 1)

						return true
					}
				}
			}
		}
	}
	return false
}

func (mw *MainWindow) timeHandler(ctx context.Context, wg *sync.WaitGroup, ticker *time.Ticker) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case t := <-ticker.C:
			mw.updateCurrentTime(t)
			if mw.workTimeRunned {
				mw.lastTime = t
				sub := t.Sub(mw.workTimeStart)
				mw.updateWorkTime(uint(sub.Seconds()))
			}
			if mw.alarmAfterRunned {
				mw.alarmAfterDuration -= 1
				mw.updateAlarmAfter(mw.alarmAfterDuration)
			}
			if mw.alarmAtRunned {
				mw.updateAlarmAt(t)
			}
		}
	}
}

func (mw *MainWindow) saveTimerAfterStopIfNeeded() {
	if !mw.workTimeStart.IsZero() && mw.lastTime.After(mw.workTimeStart) {
		duration := mw.lastTime.Sub(mw.workTimeStart)
		sec := uint(duration.Seconds())
		if sec > 5 /*(5 * 60)*/ { // we can save after 5 minutes
			h, m, s := shared.DurationComponents(sec)
			if s >= 30 {
				m += 1
				if m > 60 {
					h += 1
					m -= 60
				}
			}

			if dialog := gtk.MessageDialogNew(mw.app.GetActiveWindow(), gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "working time"); dialog != nil {
				defer dialog.Destroy()

				dialog.FormatSecondaryMarkup(workedTimeFormat, h, m)
				if dialog.Run() == gtk.RESPONSE_YES {
					if id := mw.selectedCompanyID(); id != -1 {
						if timer.NewWithData(int64(id), mw.workTimeStart.Unix(), mw.lastTime.Unix()).Save() {
							return
						}
					}
				}
			}

		}
	}
	// TODO: not saved dialog
}

func (mw *MainWindow) updateCurrentTime(t time.Time) {
	y, m, d, h, min, s := shared.DateTimeComponents(t)
	dtMarkup := fmt.Sprintf(timeFormat, h, min, s)
	tmMarkup := fmt.Sprintf(dateFormat, y, m, d, t.Weekday())
	glib.IdleAdd(func() {
		mw.timeLabel.SetMarkup(dtMarkup)
		mw.headerBar.SetSubtitle(tmMarkup)
	})
}

func (mw *MainWindow) updateWorkTime(duration uint) {
	h, m, s := shared.DurationComponents(duration)

	glib.IdleAdd(func() {
		if mw.workTimeRunned {
			mw.timerValue.SetMarkup(fmt.Sprintf(workTimeActiveFormat, h, m, s))
		} else {
			mw.timerValue.SetMarkup(fmt.Sprintf(workTimeInactiveFormat, h, m, s))
		}
	})
}

func (mw *MainWindow) updateAlarmAfter(duration uint) {
	if duration <= 0 {
		mw.alarmAfterStopHandler()
		mw.alarmAtFinished()
		return
	}

	h, m, s := shared.DurationComponents(duration)
	glib.IdleAdd(func() {
		if mw.alarmAfterRunned {
			mw.alarmAfterValue.SetMarkup(fmt.Sprintf(alarmAfterActiveFormat, h, m, s))
		} else {
			mw.alarmAfterValue.SetMarkup(fmt.Sprintf(alarmAfterInactiveFormat, h, m, s))
		}
	})
}
func (mw *MainWindow) alarmAtFinished() {
	glib.IdleAdd(func() {
		sound.PlayDrip(3)
		if dialog := gtk.MessageDialogNew(mw.app.GetActiveWindow(), gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, ""); dialog != nil {
			defer dialog.Destroy()
			h, m, s := shared.DurationComponents(mw.alarmAfterDurationPrv)
			dialog.FormatSecondaryText(fmt.Sprintf("Alarm after %d:%02d:%02d finished", h, m, s))
			dialog.Run()
		}
	})
}

func (mw *MainWindow) resetAlarmAfter() {
	glib.IdleAdd(func() {
		mw.alarmAfterValue.SetMarkup(fmt.Sprintf(alarmAfterInactiveFormat, 0, 0, 0))
	})
}

func (mw *MainWindow) updateAlarmAt(t time.Time) {
	if t.Equal(mw.alarmAt) || t.After(mw.alarmAt) {
		mw.alarmAtStopHandler()
		return
	}

	_, _, _, h, m, s := shared.DateTimeComponents(mw.alarmAt)

	glib.IdleAdd(func() {
		if mw.alarmAtRunned {
			mw.alarmAtValue.SetMarkup(fmt.Sprintf(alarmAfterActiveFormat, h, m, s))
		} else {
			mw.alarmAtValue.SetMarkup(fmt.Sprintf(alarmAfterInactiveFormat, h, m, s))
		}
	})
}

func (mw *MainWindow) resetAlarmAt() {
	glib.IdleAdd(func() {
		mw.alarmAtValue.SetMarkup(fmt.Sprintf(alarmAfterInactiveFormat, 0, 0, 0))
	})
}

func (mw *MainWindow) setAlarmAt() {
	_, _, _, h, m, s := shared.DateTimeComponents(mw.alarmAt)
	glib.IdleAdd(func() {
		mw.alarmAtValue.SetMarkup(fmt.Sprintf(alarmAfterInactiveFormat, h, m, s))
	})
}

func (mw *MainWindow) alarmAfterStartHandler() {
	mw.alarmAfterLabel.SetSensitive(true)
	mw.alarmAfterStopBtn.SetSensitive(true)
	mw.alarmAfterSetBtn.SetSensitive(false)
	mw.alarmAfterStartBtn.SetSensitive(false)
	mw.alarmAfterRunned = true
}

func (mw *MainWindow) addCompanyHandler() {
	if dialog := company.New(mw.app.GetActiveWindow(), nil); dialog != nil {
		defer dialog.Destroy()

		dialog.ShowAll()
		if dialog.Run() == gtk.RESPONSE_OK {
			if c := dialog.Company(); c != nil && c.Valid() {
				if c.Save() {
					mw.populateCompanyCombo()
					mw.selectCompanyWithID(c.ID())
					return
				}
				tr.Error("can't save the company data")
			}
		}
	}
}

func (mw *MainWindow) alarmAfterSetHandler() {
	if dialog := alarm.New(mw.app, true); dialog != nil {
		defer dialog.Destroy()

		dialog.ShowAll()
		if dialog.Run() == gtk.RESPONSE_OK {
			h, m, s := dialog.Selection()
			mw.alarmAfterDuration = s + 60*m + 60*60*h
			mw.alarmAfterDurationPrv = mw.alarmAfterDuration
			mw.updateAlarmAfter(mw.alarmAfterDuration)
			mw.alarmAfterStartBtn.SetSensitive(true)
		}
	}
}

func (mw *MainWindow) alarmAfterStopHandler() {
	glib.IdleAdd(func() {
		mw.alarmAfterLabel.SetSensitive(false)
		mw.alarmAfterStopBtn.SetSensitive(false)
		mw.alarmAfterSetBtn.SetSensitive(true)
		mw.alarmAfterStartBtn.SetSensitive(true)
		mw.alarmAfterRunned = false
		mw.alarmAfterDuration = mw.alarmAfterDurationPrv
		mw.updateAlarmAfter(mw.alarmAfterDuration)
	})
}

func (mw *MainWindow) alarmAtStartHandler() {
	mw.alarmAtLabel.SetSensitive(true)
	mw.alarmAtStopBtn.SetSensitive(true)
	mw.alarmAtSetBtn.SetSensitive(false)
	mw.alarmAtStartBtn.SetSensitive(false)
	mw.alarmAtRunned = true
}

func (mw *MainWindow) alarmAtSetHandler() {
	if dialog := alarm.New(mw.app, false); dialog != nil {
		defer dialog.Destroy()
		dialog.ShowAll()
		if dialog.Run() == gtk.RESPONSE_OK {
			h, min, s := dialog.Selection()
			y, m, d, _, _, _ := shared.DateTimeComponents(time.Now())
			mw.alarmAt = time.Date(y, time.Month(m), d, int(h), int(min), int(s), 0, time.Local)
			mw.setAlarmAt()
			mw.alarmAtStartBtn.SetSensitive(true)
		}
	}
}

func (mw *MainWindow) alarmAtStopHandler() {
	glib.IdleAdd(func() {
		mw.alarmAtRunned = false
		sound.PlayDrip(3)
		if dialog := gtk.MessageDialogNew(mw.app.GetActiveWindow(), gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, ""); dialog != nil {
			defer dialog.Destroy()
			_, _, _, h, m, s := shared.DateTimeComponents(mw.alarmAt)
			dialog.FormatSecondaryText(fmt.Sprintf("Alarm at  %d:%02d:%02d  finished", h, m, s))
			dialog.Run()
		}
		mw.alarmAtLabel.SetSensitive(false)
		mw.alarmAtStopBtn.SetSensitive(false)
		mw.alarmAtSetBtn.SetSensitive(true)
		mw.alarmAtStartBtn.SetSensitive(false)
		mw.resetAlarmAt()
	})
}

func (mw *MainWindow) companiesActionHandler() {
	if dialog := companies.New(mw.app.GetActiveWindow()); dialog != nil {
		defer dialog.Destroy()

		dialog.UpdateTable()
		dialog.ShowAll()
		dialog.Run()
	}
}

func (mw *MainWindow) statisticActionHandler() {
	if dialog := statistic.New(mw.app.GetActiveWindow()); dialog != nil {
		defer dialog.Destroy()

		dialog.ShowAll()
		dialog.Run()
	}
}

func (mw *MainWindow) aboutActionHandler() {
	if dialog, err := gtk.AboutDialogNew(); tr.IsOK(err) {
		defer dialog.Destroy()

		dialog.SetTransientFor(mw.app.GetActiveWindow())
		dialog.SetProgramName(shared.AppName)
		dialog.SetVersion(shared.AppVersion)
		dialog.SetCopyright("Copyright (c) 2019, Beesoft Software")
		dialog.SetAuthors([]string{"Piotr Pszczółkowski (piotr@beesoft.pl)"})
		dialog.SetWebsite("http://www.beesoft.pl/Timelancer")
		dialog.SetWebsiteLabel("Timelancer home page")
		dialog.SetLicense(licence)
		dialog.SetLogo(nil)
		dialog.Run()
	}
}
