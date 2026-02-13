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
	query := `SELECT
		COLUMN_NAME,
		COLUMN_TYPE,
		IS_NULLABLE,
		COLUMN_KEY,
		COLUMN_DEFAULT,
		EXTRA,
		COLUMN_COMMENT
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
	ORDER BY ORDINAL_POSITION`
	return query, []interface{}{table}
}

// ----------------------------------------------------------------------------

func (mysqlDialect) ScanColumn(rows *sql.Rows) (Column, error) {
	var col Column
	var columnType string // This contains the full type like "enum('0','1')"
	var null, key, extra string
	var dfltValue sql.NullString

	err := rows.Scan(&col.Name, &columnType, &null, &key, &dfltValue, &extra, &col.Comment)
	if err != nil {
		return col, err
	}

	col.Nullable = null == "YES"
	col.IsPrimary = key == "PRI"
	col.IsAutoIncr = strings.Contains(extra, "auto_increment")
	col.IsUnsigned = strings.Contains(strings.ToUpper(columnType), "UNSIGNED")
	col.Default = dfltValue

	// Extract ENUM/SET values from COLUMN_TYPE
	// COLUMN_TYPE comes as: "enum('0','1')" or "varchar(100)" or "int(11) unsigned"
	if strings.HasPrefix(columnType, "enum(") {
		col.EnumValues = columnType // Store full enum definition
		col.Type = "enum"
	} else if strings.HasPrefix(columnType, "set(") {
		col.EnumValues = columnType // Store full set definition
		col.Type = "set"
	} else {
		// For other types, extract base type (e.g., "varchar(100)" -> "varchar", "int(11) unsigned" -> "int")
		col.Type = columnType
		// Remove length specifiers and modifiers for type matching
		if idx := strings.Index(col.Type, "("); idx != -1 {
			col.Type = col.Type[:idx]
		}
		col.Type = strings.TrimSpace(strings.ToLower(col.Type))
		// Remove "unsigned" from type for matching
		col.Type = strings.Replace(col.Type, "unsigned", "", 1)
		col.Type = strings.TrimSpace(col.Type)
	}

	return col, nil
}

// ----------------------------------------------------------------------------

func (mysqlDialect) ForeignKeysQuery(table string) (string, []interface{}) {
	query := `SELECT
		kcu.COLUMN_NAME,
		kcu.REFERENCED_TABLE_NAME,
		kcu.REFERENCED_COLUMN_NAME
	FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
	WHERE kcu.TABLE_SCHEMA = DATABASE()
		AND kcu.TABLE_NAME = ?
		AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
	ORDER BY kcu.ORDINAL_POSITION`
	return query, []interface{}{table}
}

// ----------------------------------------------------------------------------

func (mysqlDialect) ScanForeignKey(rows *sql.Rows) (ForeignKey, error) {
	var fk ForeignKey
	err := rows.Scan(&fk.Column, &fk.ReferencedTable, &fk.ReferencedColumn)
	if err != nil {
		return fk, err
	}
	return fk, nil
}
