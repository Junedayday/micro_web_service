package main

import "strings"

func getTableMatcher(tableMatcher string, tables []string) map[string]string {
	var tMatcher = make(map[string]string)
	matchers := strings.Split(tableMatcher, ",")
	for _, matcher := range matchers {
		m := strings.Split(matcher, ":")
		if len(m) == 2 {
			tMatcher[m[0]] = m[1]
		}
	}
	for _, table := range tables {
		if _, ok := tMatcher[table]; ok {
			continue
		}
		tMatcher[table] = table
	}
	return tMatcher
}
