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
*/
import "C"
import (
	"crypto/subtle"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"unsafe"

	"Timelancer/shared/tr"
	"Timelancer/sqlite/vtc"
)

type Database struct {
	ptr   *C.sqlite3
	stmt  *C.sqlite3_stmt
	fpath string
}

var instance *Database
var once sync.Once

func SQLite() *Database {
	once.Do(func() {
		instance = &Database{}
	})
	return instance
}

func (db *Database) Version() string {
	return C.GoString(C.sqlite3_libversion())
}

func (db *Database) ErrorCode() int {
	return int(C.sqlite3_errcode(db.ptr))
}

func (db *Database) ErrorString() string {
	return C.GoString(C.sqlite3_errmsg(db.ptr))
}

func (db *Database) Open(filePath string) bool {
	if db.ptr != nil {
		log.Println("database is already opened")
		return false
	}

	if !databaseExists(filePath) {
		return false
	}

	cstr := C.CString(filePath)
	defer C.free(unsafe.Pointer(cstr))

	C.sqlite3_initialize()
	if C.sqlite3_open_v2(cstr, &db.ptr, C.SQLITE_OPEN_READWRITE, nil) == C.SQLITE_OK {
		db.fpath = filePath
		return true
	}
	db.checkError()
	return false
}

func (db *Database) Create(filePath, scheme string) bool {
	if db.ptr != nil {
		log.Println("database is already opened")
		return false
	}

	if databaseExists(filePath) {
		log.Println("database already exists")
		return false
	}

	cstr := C.CString(filePath)
	defer C.free(unsafe.Pointer(cstr))

	C.sqlite3_initialize()
	if C.sqlite3_open_v2(cstr, &db.ptr, C.SQLITE_OPEN_READWRITE|C.SQLITE_OPEN_CREATE, nil) == C.SQLITE_OK {
		if db.ExecQuery(scheme) {
			db.fpath = filePath
			return true
		}
	}
	db.checkError()
	return false
}

func (db *Database) Close() {
	if db.ptr == nil {
		return
	}
	if retv := C.sqlite3_close(db.ptr); retv == C.SQLITE_OK {
		C.sqlite3_shutdown()
		db.ptr = nil
	}
}

func (db *Database) Remove() bool {
	if err := os.Remove(db.fpath); tr.IsOK(err) {
		return true
	}
	return false
}

func (db *Database) ExecQuery(query string) bool {
	cstr := C.CString(query)
	defer C.free(unsafe.Pointer(cstr))

	if C.sqlite3_exec(db.ptr, cstr, nil, nil, nil) == C.SQLITE_OK {
		return true
	}
	db.checkError()
	return false
}

func (db *Database) LastInsertedRowID() int64 {
	return int64(int(C.sqlite3_last_insert_rowid(db.ptr)))
}

func (db *Database) BeginTransaction() bool {
	return db.ExecQuery("BEGIN IMMEDIATE TRANSACTION")
}

func (db *Database) CommitTransaction() bool {
	return db.ExecQuery("COMMIT TRANSACTION")
}

func (db *Database) RollbackTransaction() bool {
	return db.ExecQuery("ROLLBACK TRANSACTION")
}

func (db *Database) FinishTransaction(success bool) bool {
	if success {
		return db.CommitTransaction()
	}
	return db.RollbackTransaction()
}

func (db *Database) DefaultPragmas() bool {
	return db.ExecQuery("PRAGMA foreign_keys=ON")
}

/********************************************************************
*                                                                   *
*                         P R I V A T E                             *
*                                                                   *
********************************************************************/

func databaseExists(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	nbytes := len(vtc.DatabaseHeader)
	data := make([]byte, nbytes)
	if count, err := f.Read(data); tr.IsOK(err) && count == nbytes {
		return subtle.ConstantTimeCompare(vtc.DatabaseHeader, data) == 1
	}

	return true
}

func (db *Database) checkError() {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("sqlite error: %s (%d): %s (%d)\n", fn, line, db.ErrorString(), db.ErrorCode())
}
