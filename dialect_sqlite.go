package main

/*
GORM model generator
Copyright (C) 2026 Rodolfo González González

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"database/sql"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------

type sqliteDialect struct{}

// ----------------------------------------------------------------------------

func (sqliteDialect) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

// ----------------------------------------------------------------------------

func (sqliteDialect) TablesQuery() string {
	return "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
}

// ----------------------------------------------------------------------------

func (sqliteDialect) ColumnsQuery(table string) (string, []interface{}) {
	return "PRAGMA table_info(" + table + ")", nil
}

// ----------------------------------------------------------------------------

func (sqliteDialect) ScanColumn(rows *sql.Rows) (Column, error) {
	var col Column
	var cid int
	var dfltValue interface{}
	var notNull int
	var pk int

	if err := rows.Scan(&cid, &col.Name, &col.Type, &notNull, &dfltValue, &pk); err != nil {
		return Column{}, err
	}

	col.Nullable = notNull == 0
	col.IsPrimary = pk == 1
	col.IsUnsigned = strings.Contains(strings.ToUpper(col.Type), "UNSIGNED")

	return col, nil
}
