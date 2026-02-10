package main

/*
  Copyright 2026 Rodolfo González González

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
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
