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

const (
	Ok         = iota // Successful result
	Error             // SQL error or missing database
	Internal          // Internal logic error in SQLite
	Perm              // Access permission denied
	Abort             // Callback routine requested an abort
	Busy              // The database file is locked
	Locked            // A table in the database is locked
	NoMem             // A malloc() failed
	ReadOnly          // Attempt to write a readonly database
	Interrupt         // Operation terminated by sqlite3_interrupt()
	IoErr             // Some kind of disk I/O error occurred
	Corrupt           // The database disk image is malformed
	NotFound          // NOT USED. Table or record not found
	Full              // Insertion failed because database is full
	CantOpen          // Unable to open the database file
	Protocol          // NOT USED. Database lock protocol error
	Empty             // Database is empty
	Schema            // The database schema changed
	TooBig            // String or BLOB exceeds size limit
	Constraint        // Abort due to constraint violation
	Mismatch          // Data type mismatch
	Misuse            // Library used incorrectly
	NoLfs             // Uses OS features not supported on host
	Auth              // Authorization denied
	Format            // Auxiliary database format error
	Range             // 2nd parameter to sqlite3_bind out of range
	NotADb            // File opened that is not a database file
	StatusRow  = 100  // sqlite3_step() has another row ready
	StatusDone = 101  // sqlite3_step() has finished executing
)

type ValueType uint8
const (
	Null  ValueType = iota // 0
	Int                    // 1
	Float                  // 2
	Text                   // 3
	Blob                   // 4
	Bool
)

var (
	databaseHeader  = []byte{0x53, 0x51, 0x4c, 0x69, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x20, 0x33, 0x00}
	valueTypeString = [6]string{"Null", "Int", "Float", "Text", "Blob"}
)


