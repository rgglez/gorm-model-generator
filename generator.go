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
)

// ----------------------------------------------------------------------------

func generateStruct(outputPath, table, structName string, columns []Column, foreignKeys []ForeignKey, includeBaseModel bool) string {
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
	if needsTime || needsDataTypes || includeBaseModel {
		file.WriteString("import (\n")
		if needsTime {
			file.WriteString("\t\"time\"\n")
		}
		if needsDataTypes {
			file.WriteString("\t\"gorm.io/datatypes\"\n")
		}
		if includeBaseModel {
			file.WriteString("\t\"gorm.io/gorm\"\n")
		}
		file.WriteString(")\n\n")
	}

	// Add table name constant
	constName := fmt.Sprintf("TableName_%s", structName)
	file.WriteString(fmt.Sprintf("const %s = \"%s\"\n\n", constName, table))

	file.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	if includeBaseModel {
		file.WriteString("\tgorm.Model\n")
	}

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

	usedNames := make(map[string]struct{}, len(columns)+len(foreignKeys))
	for _, col := range columns {
		usedNames[toPascalCase(col.Name)] = struct{}{}
	}

	for _, fk := range foreignKeys {
		relationshipName := toPascalCase(strings.TrimSuffix(fk.Column, "_id"))
		if relationshipName == "" || relationshipName == toPascalCase(fk.Column) {
			relationshipName = toPascalCase(fk.ReferencedTable)
		}

		if _, exists := usedNames[relationshipName]; exists {
			relationshipName = relationshipName + toPascalCase(fk.ReferencedTable)
		}
		if _, exists := usedNames[relationshipName]; exists {
			continue
		}
		usedNames[relationshipName] = struct{}{}

		referencedStruct := toPascalCase(fk.ReferencedTable)
		foreignKeyField := toPascalCase(fk.Column)
		referencedField := toPascalCase(fk.ReferencedColumn)
		tags := fmt.Sprintf("`gorm:\"foreignKey:%s;references:%s\"`", foreignKeyField, referencedField)

		file.WriteString(fmt.Sprintf("\t%s %s %s\n", relationshipName, referencedStruct, tags))
	}

	file.WriteString("}\n\n")
	file.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", structName))
	file.WriteString(fmt.Sprintf("\treturn %s\n", constName))
	file.WriteString("}\n")

	return filename
}
