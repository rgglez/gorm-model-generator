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

import "strings"

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
