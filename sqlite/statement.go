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

/*
#include <stdlib.h>
#include <sqlite3.h>

#cgo LDFLAGS: -lsqlite3

int bind_text(sqlite3_stmt *stmt, int index, const char* txt) {
	return sqlite3_bind_text(stmt, index, txt, -1, SQLITE_TRANSIENT);
}

const char* column_text(sqlite3_stmt *stmt, int index) {
	return (const char *)sqlite3_column_text(stmt, index);
}
*/
import "C"
import (
	"unsafe"

	"Timelancer/shared/tr"
	"Timelancer/sqlite/field"
	"Timelancer/sqlite/row"
)

func (db *Database) step() int {
	return int(C.sqlite3_step(db.stmt))
}

func (db *Database) reset() int {
	if retv := C.sqlite3_reset(db.stmt); retv == C.SQLITE_OK {
		return int(C.sqlite3_clear_bindings(db.stmt))
	} else {
		return int(retv)
	}
}

func (db *Database) prepare(query string) bool {
	cstr := C.CString(query)
	defer C.free(unsafe.Pointer(cstr))

	if retv := C.sqlite3_prepare_v2(db.ptr, cstr, -1, &db.stmt, nil); retv == C.SQLITE_OK {
		return true
	}
	db.checkError()
	return false
}

func (db *Database) finalize() {
	if retv := C.sqlite3_finalize(db.stmt); retv == C.SQLITE_OK {
		return
	}
	db.checkError()
}

func (db *Database) columnType(idx int) ValueType {
	ct := C.sqlite3_column_type(db.stmt, C.int(idx))
	switch ct {
	case C.SQLITE_INTEGER:
		return Int
	case C.SQLITE_FLOAT:
		return Float
	case C.SQLITE_TEXT:
		return Text
	case C.SQLITE_BLOB:
		return Blob
	}
	return Null
}

func (db *Database) columnCount() int {
	return int(C.sqlite3_column_count(db.stmt))
}

func (db *Database) columnIndex(columnName string) int {
	cstr := C.CString(columnName)
	defer C.free(unsafe.Pointer(cstr))
	return int(C.sqlite3_bind_parameter_index(db.stmt, cstr))
}

func (db *Database) columnName(cindex int) string {
	return C.GoString(C.sqlite3_column_name(db.stmt, C.int(cindex)))
}

func (db *Database) bindFields(fields []*field.Field) {
	for _, f := range fields {
		columnIndex := db.columnIndex(":" + f.Name)
		switch f.ValueType {
		case Null:
			db.bindNull(columnIndex)
		case Int:
			if value, err := f.Int64(); tr.IsOK(err) {
				db.bindInt(columnIndex, value)
			}
		case Float:
			if value, err := f.Float64(); tr.IsOK(err) {
				db.bindFloat(columnIndex, value)
			}
		case Text:
			if value, err := f.Text(); tr.IsOK(err) {
				db.bindText(columnIndex, value)
			}
		case Blob:
			if value, err := f.Blob(); tr.IsOK(err) {
				db.bindBlob(columnIndex, value)
			}
		}
	}
}

func (db *Database) fetchResult() row.Result {
	var result row.Result

	n := db.columnCount()
	if n > 0 {
		for db.step() == StatusRow {
			oneRow := row.New()
			for i := 0; i < n; i++ {
				f := field.New(db.columnName(i))
				switch db.columnType(i) {
				case Null:
					f.SetValue(nil)
				case Int:
					f.SetValue(db.fetchInt(i))
				case Float:
					f.SetValue(db.fetchFloat(i))
				case Text:
					f.SetValue(db.fetchText(i))
				case Blob:
					f.SetValue(db.fetchBlob(i))
				}
				oneRow.Append(f)
			}
			result = append(result, oneRow)
		}
	}

	return result
}

/********************************************************************
*                                                                   *
*                          S E T T E R S                            *
*                                                                   *
********************************************************************/

func (db *Database) bindNull(colIndex int) int {
	return int(C.sqlite3_bind_null(db.stmt, C.int(colIndex)))
}

func (db *Database) bindInt(colIndex int, v int) int {
	return int(C.sqlite3_bind_int64(db.stmt, C.int(colIndex), C.sqlite3_int64(v)))
}

func (db *Database) bindFloat(colIndex int, v float64) int {
	return int(C.sqlite3_bind_double(db.stmt, C.int(colIndex), C.double(v)))
}

func (db *Database) bindText(colIndex int, v string) int {
	cstr := C.CString(v)
	defer C.free(unsafe.Pointer(cstr))
	return int(C.bind_text(db.stmt, C.int(colIndex), cstr))
}

func (db *Database) bindBlob(colIndex int, v []byte) int {
	return int(C.sqlite3_bind_blob(db.stmt, C.int(colIndex), unsafe.Pointer(&v[0]), C.int(len(v)), nil))
}

/********************************************************************
*                                                                   *
*                          G E T T E R S                            *
*                                                                   *
********************************************************************/

func (db *Database) fetchInt(colIndex int) int64 {
	return int64(C.sqlite3_column_int64(db.stmt, C.int(colIndex)))
}

func (db *Database) fetchFloat(colIndex int) float64 {
	return float64(C.sqlite3_column_double(db.stmt, C.int(colIndex)))
}

func (db *Database) fetchText(colIndex int) string {
	return C.GoString(C.column_text(db.stmt, C.int(colIndex)))
}

func (db *Database) fetchBlob(colIndex int) []byte {
	n := C.int(C.sqlite3_column_bytes(db.stmt, C.int(colIndex)))
	ptr := C.sqlite3_column_blob(db.stmt, C.int(colIndex))
	return C.GoBytes(ptr, n)
}
