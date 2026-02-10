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