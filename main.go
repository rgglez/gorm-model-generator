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

func generateStruct(outputPath, table, structName string, columns []Column, includeBaseModel bool) string {
	filename := fmt.Sprintf("%s/%s.go", outputPath, table)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return ""
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

	// Add table name constant
	constName := fmt.Sprintf("TableName_%s", structName)
	file.WriteString(fmt.Sprintf("const %s = \"%s\"\n\n", constName, table))

	file.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for _, col := range columns {
		fieldName := toPascalCase(col.Name)
		goType := mapSQLTypeToGo(col.Type, col.Nullable, col.IsUnsigned)

		tags := fmt.Sprintf("`gorm:\"column:%s", col.Name)

		// Add type information for ENUM and SET
		if col.EnumValues != "" {
			tags += fmt.Sprintf(";type:%s", col.EnumValues)
		}

		if col.IsPrimary {
			tags += ";primaryKey"
		}
		if col.IsAutoIncr {
			tags += ";autoIncrement"
		}
		if !col.Nullable {
			tags += ";not null"
		}
		if col.IsUnsigned {
			tags += ";unsigned"
		}

		// Add default value
		if col.Default.Valid {
			defaultVal := col.Default.String
			// Remove quotes if present (MySQL returns defaults with quotes for strings/enums)
			defaultVal = strings.Trim(defaultVal, "'\"")
			tags += fmt.Sprintf(";default:%s", defaultVal)
		}

		// Add comment
		if col.Comment != "" {
			cleanComment := cleanString(col.Comment)
			tags += fmt.Sprintf(";comment:%s", cleanComment)
		}

		tags += "\"`"

		file.WriteString(fmt.Sprintf("\t%s %s %s\n", fieldName, goType, tags))
	}

	file.WriteString("}\n\n")
	file.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", structName))
	file.WriteString(fmt.Sprintf("\treturn %s\n", constName))
	file.WriteString("}\n")

	return filename
}

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

		structName := toPascalCase(table)
		filename := generateStruct(*outputPath, table, structName, columns, *includeBaseModel)

		// Format the generated file
		if err := formatGoFile(filename); err != nil {
			fmt.Printf("  Warning: Could not format file: %v\n", err)
		}
	}

	fmt.Printf("\n✓ Structs generated successfully in: %s\n", *outputPath)
}
