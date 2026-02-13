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
	var dfltValue sql.NullString
	var notNull int
	var pk int

	err := rows.Scan(&cid, &col.Name, &col.Type, &notNull, &dfltValue, &pk)
	if err != nil {
		return col, err
	}

	col.Nullable = notNull == 0
	col.IsPrimary = pk == 1
	col.IsAutoIncr = pk == 1 && strings.Contains(strings.ToUpper(col.Type), "INTEGER")
	col.IsUnsigned = strings.Contains(strings.ToUpper(col.Type), "UNSIGNED")
	col.Default = dfltValue
	col.Comment = ""
	col.EnumValues = ""

	return col, nil
}

// ----------------------------------------------------------------------------

func (sqliteDialect) ForeignKeysQuery(table string) (string, []interface{}) {
	return "PRAGMA foreign_key_list(" + table + ")", nil
}

// ----------------------------------------------------------------------------

func (sqliteDialect) ScanForeignKey(rows *sql.Rows) (ForeignKey, error) {
	var fk ForeignKey
	var id, seq int
	var onUpdate, onDelete, match string

	err := rows.Scan(&id, &seq, &fk.ReferencedTable, &fk.Column, &fk.ReferencedColumn, &onUpdate, &onDelete, &match)
	if err != nil {
		return fk, err
	}
	return fk, nil
}
