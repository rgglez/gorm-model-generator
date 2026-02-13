package main

import (
	"database/sql"

	"gorm.io/gorm"
)

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

type Column struct {
	Name       string
	Type       string
	Nullable   bool
	IsPrimary  bool
	IsAutoIncr bool
	IsUnsigned bool
	Default    sql.NullString
	Comment    string
	EnumValues string
}

// ----------------------------------------------------------------------------

func getColumns(db *gorm.DB, table string, d dialect) ([]Column, error) {
	var columns []Column
	var rows *sql.Rows
	var err error

	query, args := d.ColumnsQuery(table)
	if args != nil {
		rows, err = db.Raw(query, args...).Rows()
	} else {
		rows, err = db.Raw(query).Rows()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		col, err := d.ScanColumn(rows)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}
