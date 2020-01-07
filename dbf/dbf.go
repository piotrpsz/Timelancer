package dbf

import (
	"GamerHash/shared"
	"Timelancer/shared/tr"
	"Timelancer/sqlite"
)

const (
	scheme = `
CREATE TABLE company
(
	id       INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	shortcut TEXT NOT NULL COLLATE NOCASE UNIQUE,
	name     TEXT NOT NULL COLLATE NOCASE UNIQUE,
	used     INTEGER NOT NULL CHECK(used==0 OR used==1) DEFAULT 1
);
CREATE TABLE timer
(
	id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	company_id INTEGER NOT NULL,
	start      INTEGER NOT NULL,
	finish     INTEGER NOT NULL,
	FOREIGN KEY (company_id) REFERENCES company(id)
)
`
)

var db *sqlite.Database = sqlite.SQLite()

func OpenOrCreate(filePath string) bool {
	if shared.FileExists(filePath) {
		if db.Open(filePath) {
			return true
		}
		tr.Error("can't open database: %v", filePath)
		return false
	}

	if db.Create(filePath, scheme) {
		return true
	}
	tr.Error("can't create database: %v", filePath)
	return false
}
