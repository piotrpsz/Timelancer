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

package sqlite

import (
	"fmt"
	"strings"

	"Timelancer/shared/tr"
	"Timelancer/sqlite/field"
	"Timelancer/sqlite/row"
	"Timelancer/sqlite/vtc"
)

func (db *Database) Select(query string) row.Result {
	if db.prepare(query) {
		defer db.finalize()

		if retv := db.fetchResult(); retv != nil {
			return retv
		}
	}

	if err := db.ErrorCode(); err != vtc.StatusDone && err != vtc.Ok {
		db.checkError()
	}
	return nil
}

func (db *Database) Insert(table string, fields []*field.Field) (int, bool) {
	if len(fields) == 0 {
		return -1, false
	}

	var b0, b1 strings.Builder
	endIdx := len(fields) - 1
	for _, f := range fields[:endIdx] {
		fmt.Fprintf(&b0, "%s,", f.Name)
		fmt.Fprintf(&b1, ":%s,", f.Name)
	}
	fmt.Fprintf(&b0, "%s", fields[endIdx].Name)
	fmt.Fprintf(&b1, ":%s", fields[endIdx].Name)

	names := b0.String()
	binds := b1.String()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, names, binds)
	if db.prepare(query) {
		defer db.finalize()

		db.bindFields(fields)
		if retv := db.step(); retv == vtc.Ok || retv == vtc.StatusDone {
			return db.LastInsertedRowID(), true
		}
	}
	db.checkError()
	return -1, false
}

func (db *Database) Update(table string, fields []*field.Field) bool {
	if len(fields) == 0 {
		return false
	}

	var b strings.Builder
	endIdx := len(fields) - 1
	for _, f := range fields[:endIdx] {
		fmt.Fprintf(&b, "%s=:%s,", f.Name, f.Name)
	}
	fmt.Fprintf(&b, "%s=:%s", fields[endIdx].Name, fields[endIdx].Name)
	assigns := b.String()

	// zakładam, że pierwsze pole to primary key (odpowiednik rowid)
	if idValue, err := fields[0].Int64(); tr.IsOK(err) {
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s=%d", table, assigns, fields[0].Name, idValue)
		if db.prepare(query) {
			defer db.finalize()

			db.bindFields(fields)
			if retv := db.step(); retv == vtc.Ok || retv == vtc.StatusDone {
				return true
			}
		}
	}

	db.checkError()
	return false
}

func (db *Database) Delete(table, idColumnName string, idValue int) bool {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s=%d", table, idColumnName, idValue)
	return db.ExecQuery(query)
}

func (db *Database) Count(table string) int {
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM %s", table)

	if result := db.Select(query); len(result) == 1 {
		if f := result[0].Field("count"); f != nil {
			if n, err := f.Int64(); tr.IsOK(err) {
				return n
			}
		}
	}
	return -1
}

func (db *Database) CountWhereInt(table string, field string, value int) int {
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE %s=%d", table, field, value)
	if result := db.Select(query); len(result) == 1 {
		if f := result[0].Field("count"); f != nil {
			if n, err := f.Int64(); tr.IsOK(err) {
				return n
			}
		}
	}
	return 0
}
