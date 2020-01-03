package company

import (
	"Timelancer/sqlite"
	"Timelancer/sqlite/field"
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
