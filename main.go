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
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------

func getTables(db *gorm.DB, d dialect) ([]string, error) {
	var tables []string
	rows, err := db.Raw(d.TablesQuery()).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// ----------------------------------------------------------------------------

func getColumns(db *gorm.DB, table string, d dialect) ([]Column, error) {
	var columns []Column
	query, args := d.ColumnsQuery(table)
	rows, err := db.Raw(query, args...).Rows()
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

// ----------------------------------------------------------------------------

func generateStruct(outputPath, table, structName string, columns []Column) {
	filename := fmt.Sprintf("%s/%s.go", outputPath, table)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Check if we need to import time or datatypes
	needsTime := false
	needsDataTypes := false
	for _, col := range columns {
		goType := mapSQLTypeToGo(col.Type, col.Nullable, col.IsUnsigned)
		if strings.Contains(goType, "time.Time") {
			needsTime = true
		}
		if strings.Contains(goType, "datatypes.") {
			needsDataTypes = true
		}
	}

	// Write package and imports
	file.WriteString("package models\n\n")
	if needsTime || needsDataTypes {
		file.WriteString("import (\n")
		if needsTime {
			file.WriteString("\t\"time\"\n")
		}
		if needsDataTypes {
			file.WriteString("\t\"gorm.io/datatypes\"\n")
		}
		file.WriteString(")\n\n")
	}

	file.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for _, col := range columns {
		fieldName := toPascalCase(col.Name)
		goType := mapSQLTypeToGo(col.Type, col.Nullable, col.IsUnsigned)

		tags := fmt.Sprintf("`gorm:\"column:%s", col.Name)
		if col.IsPrimary {
			tags += ";primaryKey"
		}
		if col.IsAutoIncr {
			tags += ";autoIncrement"
		}
		tags += "\"`"

		file.WriteString(fmt.Sprintf("\t%s %s %s\n", fieldName, goType, tags))
	}

	file.WriteString("}\n\n")
	file.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", structName))
	file.WriteString(fmt.Sprintf("\treturn \"%s\"\n", table))
	file.WriteString("}\n")
}

// ----------------------------------------------------------------------------

func main() {
	dsn := flag.String("dsn", "", "Database DSN connection string")
	dbType := flag.StringP("type", "t", "mysql", "Database type (mysql, postgres, sqlite)")
	outputPath := flag.StringP("output", "o", "./models", "Output path for generated files")
	tableName := flag.String("tables", "", "Specific table name (empty for all tables)")
	flag.Parse()

	if *dsn == "" {
		*dsn = os.Getenv("DATABASE_DSN")
	}

	if *dsn == "" {
		fmt.Println("Error: Database DSN not provided")
		fmt.Println("Usage:")
		fmt.Println("  --dsn=\"user:pass@tcp(localhost:3306)/dbname\"")
		fmt.Println("  --type=mysql (mysql, postgres, sqlite)")
		fmt.Println("  --output=./models")
		fmt.Println("  --tables=users (optional, comma separade names for specific tables)")
		os.Exit(1)
	}

	d, err := newDialect(*dbType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	db, err := d.Open(*dsn)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Successfully connected to database")

	// Get list of tables
	var tables []string
	if *tableName != "" {
		tables = strings.Split(*tableName, ",")
	} else {
		tables, err = getTables(db, d)
		if err != nil {
			fmt.Printf("Error getting tables: %v\n", err)
			os.Exit(1)
		}
	}

	if len(tables) == 0 {
		fmt.Println("No tables found in database")
		os.Exit(0)
	}

	// Create output directory
	if err := os.MkdirAll(*outputPath, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Generate structs for each table
	for _, table := range tables {
		fmt.Printf("Generating struct for table: %s\n", table)
		columns, err := getColumns(db, table, d)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		structName := toPascalCase(table)
		generateStruct(*outputPath, table, structName, columns)
	}

	fmt.Printf("\n✓ Structs generated successfully in: %s\n", *outputPath)
}
