package timer

import (
	"fmt"
	"time"

	"Timelancer/shared"
	"Timelancer/shared/tr"
	"Timelancer/sqlite"
	"Timelancer/sqlite/field"
	"Timelancer/sqlite/row"
)

/*
CREATE TABLE timer
(
	company_id INTEGER NOT NULL,
	start      INTEGER NOT NULL,
	finish     INTEGER NOT NULL,
	FOREIGN KEY (company_id) REFERENCES company(id)
)
*/

type Timer struct {
	id        int64
	companyID int64
	start     int64
	finish    int64
}

func NewWithData(companyID, start, finish int64) *Timer {
	return &Timer{companyID: companyID, start: start, finish: finish}
}

func NewWithRow(r row.Row) *Timer {
	tm := &Timer{}
	ok := false

	if value, exists := r["id"]; exists {
		if value, err := value.Int64(); tr.IsOK(err) {
			tm.id = int64(value)
			ok = true
		}
	}
	if ok {
		ok = false
		if value, exists := r["company_id"]; exists {
			if value, err := value.Int64(); tr.IsOK(err) {
				tm.companyID = value
				ok = true
			}
		}
	}
	if ok {
		ok = false
		if value, exists := r["start"]; exists {
			if value, err := value.Int64(); tr.IsOK(err) {
				tm.start = value
				ok = true
			}
		}
	}
	if ok {
		ok = false
		if value, exists := r["finish"]; exists {
			if value, err := value.Int64(); tr.IsOK(err) {
				tm.finish = value
				ok = true
			}
		}
	}

	if ok {
		return tm
	}
	return nil
}

func (tm *Timer) ID() int64 {
	return tm.id
}

func (tm *Timer) CompanyID() int64 {
	return tm.companyID
}

func (tm *Timer) Start() string {
	if t := time.Unix(tm.start, 0); !t.IsZero() {
		return shared.TimeAsString(t)
	}
	return ""
}

func (tm *Timer) Finish() string {
	if t := time.Unix(tm.finish, 0); !t.IsZero() {
		return shared.TimeAsString(t)
	}
	return ""
}

func (tm *Timer) StartAsString() string {
	return ""
}

func (tm *Timer) FinishAsString() string {
	return ""
}

func (tm *Timer) Valid() bool {
	return tm.id != 0 && tm.companyID != 0 && tm.start != 0 && tm.finish != 0
}

func (tm *Timer) Remove() bool {
	query := fmt.Sprintf("DELETE FROM timer WHERE id=%d", tm.id)
	return sqlite.SQLite().ExecQuery(query)
}

func (tm *Timer) Save() bool {
	if tm.id == 0 {
		return tm.insert()
	}
	return tm.update()
}

func (tm *Timer) fields() []*field.Field {
	var data []*field.Field

	if tm.id > 0 {
		data = append(data, field.NewWithValue("id", int64(tm.id)))
	}
	data = append(data, field.NewWithValue("company_id", int64(tm.companyID)))
	data = append(data, field.NewWithValue("start", int64(tm.start)))
	data = append(data, field.NewWithValue("finish", int64(tm.finish)))

	return data
}

func (tm *Timer) insert() bool {
	fields := tm.fields()
	if id, ok := sqlite.SQLite().Insert("timer", fields); ok {
		tm.id = id
		return true
	}
	return false
}

func (tm *Timer) update() bool {
	fields := tm.fields()
	return sqlite.SQLite().Update("timer", fields)
}
