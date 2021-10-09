package main

var tableInfo []TableInfo

type TableInfo struct {
	Name            string `yaml:"name"`
	GoStruct        string `yaml:"goStruct"`
	CreateTime      string `yaml:"createTime"`
	UpdateTime      string `yaml:"updateTime"`
	SoftDeleteKey   string `yaml:"softDeleteKey"`
	SoftDeleteValue int    `yaml:"softDeleteValue"`
}

func getTableMatcher() map[string]TableInfo {
	var tMatcher = make(map[string]TableInfo)
	for _, matcher := range tableInfo {
		tMatcher[matcher.Name] = matcher
	}
	return tMatcher
}
