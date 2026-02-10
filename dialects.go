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
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------

type Column struct {
	Name       string
	Type       string
	Nullable   bool
	IsPrimary  bool
	IsAutoIncr bool
	IsUnsigned bool
}

// ----------------------------------------------------------------------------

type dialect interface {
	Open(dsn string) (*gorm.DB, error)
	TablesQuery() string
	ColumnsQuery(table string) (string, []interface{})
	ScanColumn(rows *sql.Rows) (Column, error)
}

// ----------------------------------------------------------------------------

var dialectFactory = map[string]func() dialect{
	"mysql":      func() dialect { return mysqlDialect{} },
	"postgres":   func() dialect { return postgresDialect{} },
	"postgresql": func() dialect { return postgresDialect{} },
	"sqlite":     func() dialect { return sqliteDialect{} },
}

// ----------------------------------------------------------------------------

func newDialect(dbType string) (dialect, error) {
	factory, ok := dialectFactory[strings.ToLower(dbType)]
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	return factory(), nil
}