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
	"os/exec"
	"strings"
)

// ----------------------------------------------------------------------------

// toPascalCase converts snake_case to PascalCase
func toPascalCase(s string) string {
	// Split by underscore
	parts := strings.Split(s, "_")
	result := ""

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		// Capitalize first letter and append the rest
		result += strings.ToUpper(part[0:1]) + strings.ToLower(part[1:])
	}

	return result
}

// ----------------------------------------------------------------------------

func matchesAny(parts ...string) func(string) bool {
	return func(value string) bool {
		for _, part := range parts {
			if strings.Contains(value, part) {
				return true
			}
		}
		return false
	}
}

// ----------------------------------------------------------------------------

func cleanString(value string) string {
	cleanValue := strings.ReplaceAll(value, "\n", " ")
	cleanValue = strings.ReplaceAll(cleanValue, "\r", " ")
	cleanValue = strings.ReplaceAll(cleanValue, "\t", " ")
	cleanValue = strings.Join(strings.Fields(cleanValue), " ")

	return cleanValue
}

// ----------------------------------------------------------------------------

func formatGoFile(filename string) error {
	// Try goimports first (it's better as it also organizes imports)
	cmd := exec.Command("goimports", "-w", filename)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to gofmt
	cmd = exec.Command("gofmt", "-w", filename)
	return cmd.Run()
}
