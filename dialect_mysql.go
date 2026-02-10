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

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------

type mysqlDialect struct{}

// ----------------------------------------------------------------------------

func (mysqlDialect) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// ----------------------------------------------------------------------------

func (mysqlDialect) TablesQuery() string {
	return "SHOW TABLES"
}

// ----------------------------------------------------------------------------

func (mysqlDialect) ColumnsQuery(table string) (string, []interface{}) {
	return "SHOW COLUMNS FROM " + table, nil
}

// ----------------------------------------------------------------------------

func (mysqlDialect) ScanColumn(rows *sql.Rows) (Column, error) {
	var col Column
	var null, key, extra string
	var dflt interface{}

	if err := rows.Scan(&col.Name, &col.Type, &null, &key, &dflt, &extra); err != nil {
		return Column{}, err
	}

	col.Nullable = null == "YES"
	col.IsPrimary = key == "PRI"
	col.IsAutoIncr = strings.Contains(extra, "auto_increment")
	col.IsUnsigned = strings.Contains(strings.ToUpper(col.Type), "UNSIGNED")

	return col, nil
}
