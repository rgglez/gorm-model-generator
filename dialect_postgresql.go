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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------

type postgresDialect struct{}

// ----------------------------------------------------------------------------

func (postgresDialect) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// ----------------------------------------------------------------------------

func (postgresDialect) TablesQuery() string {
	return "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
}

// ----------------------------------------------------------------------------

func (postgresDialect) ColumnsQuery(table string) (string, []interface{}) {
	query := `SELECT column_name, data_type, is_nullable,
		CASE WHEN column_name IN (
			SELECT a.attname FROM pg_index i
			JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
			WHERE i.indrelid = $1::regclass AND i.indisprimary
		) THEN 'YES' ELSE 'NO' END as is_primary
		FROM information_schema.columns
		WHERE table_name = $1`
	return query, []interface{}{table}
}

// ----------------------------------------------------------------------------

func (postgresDialect) ScanColumn(rows *sql.Rows) (Column, error) {
	var col Column
	var nullable, isPrimary string

	if err := rows.Scan(&col.Name, &col.Type, &nullable, &isPrimary); err != nil {
		return Column{}, err
	}

	col.Nullable = nullable == "YES"
	col.IsPrimary = isPrimary == "YES"
	col.IsUnsigned = false

	return col, nil
}