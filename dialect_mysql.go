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
