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

package dt

import (
	"fmt"
	"time"
)

var location *time.Location

type Datime struct {
	value time.Time
}

func init() {
	location = time.Now().Location()
	fmt.Println("*", location)
}

func New() Datime {
	return Datime{value: time.Now()}
}

func NewUnix(seconds int64) Datime {
	return Datime{value: time.Unix(seconds, 0)}
}

func (dt Datime) Unix() int64 {
	return dt.value.Unix()
}

func NewWithComponents(year, month, day, hour, min, sec int) Datime {
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, location)
	return Datime{value: t}
}

func (dt Datime) String() string {
	year, month, day, hour, min, sec := dt.Components()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func (dt Datime) Components() (int, int, int, int, int, int) {
	year, month, day := dt.value.Date()
	hour := dt.value.Hour()
	min := dt.value.Minute()
	sec := dt.value.Second()
	return year, int(month), day, hour, min, sec
}

func (dt Datime) StartOfDay() Datime {
	year, month, day := dt.value.Date()
	t := time.Date(year, month, day, 0, 0, 0, 0, location)
	return Datime{value: t}
}

func (dt Datime) EndOfDay() Datime {
	year, month, day := dt.value.Date()
	t := time.Date(year, month, day, 23, 59, 59, 0, location)
	return Datime{value: t}
}

func (dt Datime) FirstOfWeek() Datime {
	imap := []int{6, 0, 1, 2, 3, 4, 5}
	dayOfWeek := imap[int(dt.value.Weekday())]
	return dt.AddDay(-dayOfWeek)
}
func (dt Datime) LastOfWeek() Datime {
	imap := []int{6, 0, 1, 2, 3, 4, 5}
	dayOfWeek := imap[int(dt.value.Weekday())]
	return dt.AddDay(6 - dayOfWeek)
}

func (dt Datime) FirstOfMonth() Datime {
	year, month, _ := dt.value.Date()
	t := time.Date(year, month, 1, 0, 0, 0, 0, location)
	return Datime{value: t}
}

func (dt Datime) LastOfMonth() Datime {
	tm := dt.FirstOfMonth().value
	tm = tm.AddDate(0, 1, -1)
	return Datime{value: tm}
}

func (dt Datime) FirstOfYear() Datime {
	year, _, _ := dt.value.Date()
	t := time.Date(year, 1, 1, 0, 0, 0, 0, location)
	return Datime{value: t}
}

func (dt Datime) LastOfYear() Datime {
	year, _, _ := dt.value.Date()
	t := time.Date(year, 12, 31, 0, 0, 0, 0, location)
	return Datime{value: t}
}

func (dt Datime) AddYear(n int) Datime {
	t := dt.value.AddDate(n, 0, 0)
	return Datime{value: t}
}

func (dt Datime) AddMonth(n int) Datime {
	t := dt.value.AddDate(0, n, 0)
	return Datime{value: t}
}

func (dt Datime) AddDay(n int) Datime {
	t := dt.value.AddDate(0, 0, n)
	return Datime{value: t}
}

func (dt Datime) AddHour(n int) Datime {
	t := dt.value.Add(time.Duration(n) * time.Hour)
	return Datime{value: t}
}

func (dt Datime) AddMinute(n int) Datime {
	t := dt.value.Add(time.Duration(n) * time.Minute)
	return Datime{value: t}
}

func (dt Datime) AddSecond(n int) Datime {
	t := dt.value.Add(time.Duration(n) * time.Second)
	return Datime{value: t}
}
