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
	"Timelancer/model/company"
)

var companiesData []*company.Company

func (mw *MainWindow) populateCompanyCombo() {
	mw.companyCombo.RemoveAll()
	mw.companyCombo.AppendText("Select a company")
	mw.companyCombo.SetActive(0)

	companiesData = company.CompaniesInUse()
	for _, c := range companiesData {
		mw.companyCombo.AppendText(c.Name())
	}
}

func (mw *MainWindow) selectCompanyWithID(id int) bool {
	for index, c := range companiesData {
		if c.ID() == id {
			mw.companyCombo.SetActive(index + 1)
			return true
		}
	}
	return false
}

func (mw *MainWindow) selectedCompanyChanged() {
	mw.companyIndex = mw.companyCombo.GetActive()

	if mw.companyIndex == 0 || mw.companyIndex == -1 {
		mw.companyLabel.SetSensitive(false)
		mw.timerLabel.SetSensitive(false)
		mw.timerValue.SetSensitive(false)
		mw.timerStartBtn.SetSensitive(false)
		mw.timerStopBtn.SetSensitive(false)
	} else {
		mw.companyLabel.SetSensitive(true)
		mw.timerLabel.SetSensitive(true)
		mw.timerValue.SetSensitive(true)
		mw.timerStartBtn.SetSensitive(true)
		mw.timerStopBtn.SetSensitive(false)
	}
}
