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

package row

import (
	"log"

	"Timelancer/sqlite/field"
)

type (
	Row map[string]*field.Field
	Result []Row
)

func New() Row {
	return make(Row)
}

func (r Row) Append(f *field.Field) bool {
	if !f.Valid() {
		log.Printf("Field %s is not valid", f.Name)
		return false
	}
	if r.containsField(f.Name) {
		log.Printf("field %s already exists", f.Name)
		return false
	}
	r[f.Name] = f
	return true
}

func (r Row) Field(name string) *field.Field {
	if f, ok := r[name]; ok {
		return f
	}
	return nil
}

func (r Row) containsField(name string) bool {
	if _, ok := r[name]; ok {
		return true
	}
	return false
}

func (r Row) Count() int {
	return len(r)
}

func (r Row) IsNotEmpty() bool {
	return len(r) > 0
}
