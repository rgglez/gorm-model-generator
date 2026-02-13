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
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// ----------------------------------------------------------------------------

func main() {
	dsn := flag.String("dsn", "", "Database DSN connection string")
	dbType := flag.StringP("type", "t", "mysql", "Database type (mysql, postgres, sqlite)")
	outputPath := flag.StringP("output", "o", "./models", "Output path for generated files")
	tableName := flag.String("tables", "", "Specific table name (empty for all tables)")
	includeBaseModel := flag.BoolP("include-base", "b", false, "Include base GORM model (gorm.Model)")
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
		fmt.Println("  --include-base (optional, includes gorm.Model in every generated struct)")
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
		foreignKeys, err := getForeignKeys(db, table, d)
		if err != nil {
			fmt.Printf("  Warning: could not read foreign keys: %v\n", err)
			foreignKeys = nil
		}
		foreignKeys = mergeForeignKeys(foreignKeys, inferForeignKeys(table, columns, tables))

		structName := toPascalCase(table)
		filename := generateStruct(*outputPath, table, structName, columns, foreignKeys, *includeBaseModel)

		// Format the generated file
		if err := formatGoFile(filename); err != nil {
			fmt.Printf("  Warning: Could not format file: %v\n", err)
		}
	}

	fmt.Printf("\n✓ Structs generated successfully in: %s\n", *outputPath)
}
