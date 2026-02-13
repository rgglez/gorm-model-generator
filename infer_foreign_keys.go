package main

import "strings"

func inferForeignKeys(table string, columns []Column, tables []string) []ForeignKey {
	tableSet := make(map[string]struct{}, len(tables))
	for _, name := range tables {
		tableSet[name] = struct{}{}
	}

	var inferred []ForeignKey
	for _, col := range columns {
		if !strings.HasSuffix(col.Name, "_id") {
			continue
		}

		base := strings.TrimSuffix(col.Name, "_id")
		candidates := []string{
			base,
			base + "s",
		}
		if strings.HasSuffix(base, "y") {
			candidates = append(candidates, strings.TrimSuffix(base, "y")+"ies")
		}

		for _, candidate := range candidates {
			if candidate == table {
				continue
			}
			if _, ok := tableSet[candidate]; !ok {
				continue
			}

			inferred = append(inferred, ForeignKey{
				Column:           col.Name,
				ReferencedTable:  candidate,
				ReferencedColumn: "id",
			})
			break
		}
	}

	return inferred
}

func mergeForeignKeys(existing []ForeignKey, inferred []ForeignKey) []ForeignKey {
	if len(inferred) == 0 {
		return existing
	}

	seen := make(map[string]struct{}, len(existing)+len(inferred))
	for _, fk := range existing {
		key := fk.Column + ":" + fk.ReferencedTable + ":" + fk.ReferencedColumn
		seen[key] = struct{}{}
	}

	merged := make([]ForeignKey, 0, len(existing)+len(inferred))
	merged = append(merged, existing...)
	for _, fk := range inferred {
		key := fk.Column + ":" + fk.ReferencedTable + ":" + fk.ReferencedColumn
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, fk)
	}

	return merged
}
