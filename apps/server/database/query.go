package database

import "fmt"

func PrefixQuery(table string, columns []string) []string {
	prefixed := make([]string, len(columns))
	for i, col := range columns {
		prefixed[i] = fmt.Sprintf("public.%s.%s", table, col)
	}
	return prefixed
}
