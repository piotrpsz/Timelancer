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

package shared

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"Carmel/shared/tr"
)

const (
	AppName    = "timelancer"
	AppVersion = "0.1.0"
	appDir     = ".timelancer"
)

func AppDir() string {
	if homeDir, err := os.UserHomeDir(); tr.IsOK(err) {
		appDir := filepath.Join(homeDir, appDir)
		if CreateDirIfNeeded(appDir) {
			return appDir
		}
	}
	return ""
}

/********************************************************************
*                                                                   *
*                       D A T E   &   T I M E                       *
*                                                                   *
********************************************************************/

func Now() time.Time {
	t := time.Now().UTC()
	year, month, day := t.Date()
	hour, min, sec := t.Hour(), t.Minute(), t.Second()
	// without miliseconds
	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}

func TimeAsString(t time.Time) string {
	year, month, day, hour, min, sec := DateTimeComponents(t)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func TimeAndDate(t time.Time) (string, string) {
	year, month, day, hour, min, sec := DateTimeComponents(t)
	dt := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	tm := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
	return tm, dt
}

func DateTimeComponents(t time.Time) (int, int, int, int, int, int) {
	year, month, day := t.Date()
	hour, min, sec := t.Hour(), t.Minute(), t.Second()
	return year, int(month), day, hour, min, sec
}

func DurationComponents(seconds uint) (uint, uint, uint) {
	h, m, s := uint(0), uint(0), uint(0)

	if seconds > 0 {
		s = seconds % 60
		m = seconds / 60
		if m > 59 {
			h = m / 60
			m = m % 60
		}
	}

	return h, m, s
}

/********************************************************************
*                                                                   *
*             F I L E S   &   D I R E C T O R I E S                 *
*                                                                   *
********************************************************************/

func ExistsFile(filePath string) bool {
	var err error
	var fi os.FileInfo

	if fi, err = os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
		}
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

func ExistsDir(dirPath string) bool {
	var err error
	var fi os.FileInfo

	if fi, err = os.Stat(dirPath); err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
		}
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func CreateDirIfNeeded(dirPath string) bool {
	if ExistsDir(dirPath) {
		return true
	}
	if err := os.MkdirAll(dirPath, os.ModePerm); tr.IsOK(err) {
		return true
	}
	return false
}

func RemoveFile(filePath string) bool {
	if err := os.Remove(filePath); tr.IsOK(err) {
		return true
	}
	return false
}

func ReadFromFile(filePath string) []byte {
	if handle, err := os.OpenFile(filePath, os.O_RDONLY, 0666); tr.IsOK(err) {
		if buffer, err := ioutil.ReadAll(handle); tr.IsOK(err) {
			return buffer
		}
	}
	return nil
}

func WriteToFile(filePath string, data []byte) bool {
	if handle, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666); tr.IsOK(err) {
		if nbytes, err := handle.Write(data); tr.IsOK(err) {
			return nbytes == len(data)
		}
	}
	return false
}
