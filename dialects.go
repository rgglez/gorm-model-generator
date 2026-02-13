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
	"fmt"
	"strings"

	"gorm.io/gorm"
)

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
