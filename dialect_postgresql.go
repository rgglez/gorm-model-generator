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