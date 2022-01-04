package main

var tableInfo []TableInfo

type TableInfo struct {
	Name            string     `yaml:"name"`
	GoStruct        string     `yaml:"goStruct"`
	CreateTime      string     `yaml:"createTime"`
	UpdateTime      string     `yaml:"updateTime"`
	SoftDeleteKey   string     `yaml:"softDeleteKey"`
	SoftDeleteValue int        `yaml:"softDeleteValue"`
	GenQueries      []GenQuery `yaml:"genQueries"`
}

type GenQuery struct {
	Method  string
	Desc    string
	Where   string
	Args    []ArgInfo
	OrderBy string
}

type ArgInfo struct {
	Name string
	Type string
}

func getTableMatcher() map[string]TableInfo {
	var tMatcher = make(map[string]TableInfo)
	for _, matcher := range tableInfo {
		tMatcher[matcher.Name] = matcher
	}
	return tMatcher
}
