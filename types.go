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
	"strings"
)

// ----------------------------------------------------------------------------

type goTypeFactory func(unsigned bool) string

// ----------------------------------------------------------------------------

type sqlTypeRule struct {
	matches      func(string) bool
	factory      goTypeFactory
	skipNullWrap bool
}

// ----------------------------------------------------------------------------

func constantTypeFactory(dataType string) goTypeFactory {
	return func(bool) string {
		return dataType
	}
}

// ----------------------------------------------------------------------------

func signedUnsignedTypeFactory(signedType, unsignedType string) goTypeFactory {
	return func(unsigned bool) string {
		if unsigned {
			return unsignedType
		}
		return signedType
	}
}

// ----------------------------------------------------------------------------

var sqlTypeRules = []sqlTypeRule{
	{matches: matchesAny("tinyint(1)"), factory: constantTypeFactory("bool")},
	{matches: matchesAny("bigint"), factory: signedUnsignedTypeFactory("int64", "uint64")},
	{matches: matchesAny("mediumint"), factory: signedUnsignedTypeFactory("int32", "uint32")},
	{matches: matchesAny("smallint"), factory: signedUnsignedTypeFactory("int16", "uint16")},
	{matches: matchesAny("tinyint"), factory: signedUnsignedTypeFactory("int8", "uint8")},
	{matches: matchesAny("int"), factory: signedUnsignedTypeFactory("int", "uint")},
	{matches: matchesAny("varchar", "text", "char", "character"), factory: constantTypeFactory("string")},
	{matches: matchesAny("decimal", "numeric"), factory: constantTypeFactory("float64")},
	{matches: matchesAny("float", "double"), factory: constantTypeFactory("float64")},
	{matches: matchesAny("bool"), factory: constantTypeFactory("bool")},
	{matches: matchesAny("date", "time"), factory: constantTypeFactory("time.Time"), skipNullWrap: true},
	{matches: matchesAny("json"), factory: constantTypeFactory("datatypes.JSON"), skipNullWrap: true},
}

// ----------------------------------------------------------------------------

func resolveGoType(sqlType string, unsigned bool) (string, bool) {
	for _, rule := range sqlTypeRules {
		if rule.matches(sqlType) {
			return rule.factory(unsigned), rule.skipNullWrap
		}
	}
	return "string", false
}

// ----------------------------------------------------------------------------

func mapSQLTypeToGo(sqlType string, nullable bool, unsigned bool) string {
	baseType, skipNullWrap := resolveGoType(strings.ToLower(sqlType), unsigned)
	if skipNullWrap {
		return baseType
	}

	if nullable && baseType != "string" {
		return "*" + baseType
	}
	return baseType
}