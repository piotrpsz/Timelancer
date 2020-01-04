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

package company

import (
	"Timelancer/shared/tr"
	"Timelancer/sqlite"
	"Timelancer/sqlite/field"
	"Timelancer/sqlite/row"
)

/*
CREATE TABLE company
(
id       INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
shortcut TEXT NOT NULL COLLATE NOCASE UNIQUE,
name     TEXT NOT NULL COLLATE NOCASE,
used     INTEGER NOT NULL CHECK(used==0 OR used==1) DEFAULT 1
);
*/

type Company struct {
	id       int
	shortcut string
	name     string
	used     bool
}

func New() *Company {
	return &Company{id: 0, used: true}
}

func NewWithRow(r row.Row) *Company {
	c := &Company{}
	ok := false

	if value, exists := r["id"]; exists {
		if id, err := value.Int64(); tr.IsOK(err) {
			c.id = id
			ok = true
		}
	}
	if ok {
		ok = false
		if value, exists := r["shortcut"]; exists {
			if shortcut, err := value.Text(); tr.IsOK(err) {
				c.shortcut = shortcut
				ok = true
			}
		}
	}
	if ok {
		ok = false
		if value, exists := r["name"]; exists {
			if name, err := value.Text(); tr.IsOK(err) {
				c.name = name
				ok = true
			}
		}
	}
	if ok {
		ok = false
		if value, exists := r["used"]; exists {
			if used, err := value.Bool(); tr.IsOK(err) {
				c.used = used
				ok = true
			}
		}
	}

	if ok {
		return c
	}
	return nil
}

func (c *Company) ID() int {
	return c.id
}

func (c *Company) Shortcut() string {
	return c.shortcut
}

func (c *Company) Name() string {
	return c.name
}

func (c *Company) Used() bool {
	return c.used
}

func (c *Company) SetShortcut(value string) {
	c.shortcut = value
}

func (c *Company) SetName(value string) {
	c.name = value
}

func (c *Company) SetUsed(value bool) {
	c.used = value
}

func (c *Company) Valid() bool {
	return c.name != "" && c.shortcut != ""
}

func (c *Company) Save() bool {
	if c.id == 0 {
		return c.insert()
	}
	return c.update()
}

func (c *Company) fields() []*field.Field {
	var data []*field.Field

	if c.id > 0 {
		data = append(data, field.NewWithValue("id", c.id))
	}
	data = append(data, field.NewWithValue("shortcut", c.shortcut))
	data = append(data, field.NewWithValue("name", c.name))
	data = append(data, field.NewWithValue("used", c.used))

	return data
}

func (c *Company) insert() bool {
	fields := c.fields()
	if id, ok := sqlite.SQLite().Insert("company", fields); ok {
		c.id = id
		return true
	}
	return false
}

func (c *Company) update() bool {
	fields := c.fields()
	return sqlite.SQLite().Update("company", fields)
}

func CompaniesInUse() []*Company {
	if n := sqlite.SQLite().CountWhereInt("company", "used", 1); n > 0 {
		var data []*Company
		query := "SELECT * FROM company WHERE used=1 ORDER BY shortcut ASC"
		if result := sqlite.SQLite().Select(query); len(result) > 0 {
			for _, r := range result {
				if c := NewWithRow(r); c != nil {
					data = append(data, c)
				}
			}
			if len(data) > 0 {
				return data
			}
		}
	}
	return nil
}
