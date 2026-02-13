package main

import (
	"database/sql"

	"gorm.io/gorm"
)

type ForeignKey struct {
	Column           string
	ReferencedTable  string
	ReferencedColumn string
}

func getForeignKeys(db *gorm.DB, table string, d dialect) ([]ForeignKey, error) {
	var foreignKeys []ForeignKey
	query, args := d.ForeignKeysQuery(table)

	var rows *sql.Rows
	var err error
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
		fk, err := d.ScanForeignKey(rows)
		if err != nil {
			return nil, err
		}
		foreignKeys = append(foreignKeys, fk)
	}

	return foreignKeys, nil
}
