package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
)

const (
	showTablesSQL      = "SHOW TABLES"
	showCreateTableSQL = "SHOW CREATE TABLE %s"
	commentMark        = "COMMENT"
)

var tableToGoStruct = map[string]string{
	"bigint":    "int64",
	"int":       "int",
	"smallint":  "int",
	"tinyint":   "int",
	"mediumint": "int",
	"decimal":   "float64",
	"numeric":   "float64",
	"float":     "float64",
	"datetime":  "time.Time",
	"date":      "time.Time",
	"timestamp": "time.Time",
	"varchar":   "string",
	"char":      "string",
	"text":      "string",
}

func getAllTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(showTablesSQL)
	if err != nil {
		return nil, errors.Wrapf(err, "query tables failed")
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, errors.Wrapf(err, "scan result failed")
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func Generate(db *sql.DB, table string, matchInfo TableInfo) (StructLevel, error) {
	var structData = StructLevel{
		TableName: table,
	}
	structName := camelCase(matchInfo.GoStruct)
	camelStructName := string(unicode.ToUpper(rune(structName[0]))) + structName[1:]
	structData.StructName.UpperS = inflection.Singular(structName)
	structData.StructName.UpperP = inflection.Plural(structName)

	camelStructName = string(unicode.ToLower(rune(structName[0]))) + structName[1:]
	structData.StructName.LowerS = inflection.Singular(camelStructName)
	structData.StructName.LowerP = inflection.Plural(camelStructName)

	rows, err := db.Query(fmt.Sprintf(showCreateTableSQL, table))
	if err != nil {
		return structData, errors.Wrapf(err, "show table failed")
	}
	defer rows.Close()

	var t, s string
	for rows.Next() {
		err = rows.Scan(&t, &s)
		if err != nil {
			return structData, err
		}
		structData.Columns = parseTable(s)
	}

	for _, v := range structData.Columns {
		switch v.GormName {
		case matchInfo.CreateTime:
			structData.TableCreateTime = matchInfo.CreateTime
			structData.FieldCreateTime = v.FieldName
		case matchInfo.UpdateTime:
			structData.TableUpdateTime = matchInfo.UpdateTime
			structData.FieldUpdateTime = v.FieldName
		case matchInfo.SoftDeleteKey:
			structData.TableSoftDeleteKey = matchInfo.SoftDeleteKey
			structData.TableSoftDeleteValue = matchInfo.SoftDeleteValue
			structData.FieldSoftDeleteKey = v.FieldName
		}
	}
	structData.GenQueries = matchInfo.GenQueries
	return structData, nil
}

func parseTable(s string) []FieldLevel {
	lines := strings.Split(s, "\n")
	var columns []FieldLevel
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, "`") {
			p := strings.Split(line, " ")
			name := strings.Trim(p[0], "`")
			dataType := p[1]
			comment := ""
			if strings.ToUpper(p[len(p)-2]) == commentMark {
				comment = strings.Trim(strings.Trim(p[len(p)-1], ","), "'")
			}
			columns = append(columns, FieldLevel{
				FieldName: camelCase(name),
				FieldType: fieldType(dataType),
				GormName:  name,
				Comment:   comment,
			})
		}
	}
	return columns
}

func fieldType(dataType string) string {
	for tableType, goType := range tableToGoStruct {
		if strings.HasPrefix(dataType, tableType) {
			return goType
		}
	}
	return "unknown"
}

func camelCase(s string) string {
	var buf bytes.Buffer
	var flag = false
	for i, c := range s {
		if c == '_' {
			flag = true
			continue
		}
		if i == 0 {
			buf.WriteRune(unicode.ToUpper(c))
			continue
		}
		if flag {
			buf.WriteRune(unicode.ToUpper(c))
			flag = false
			continue
		}
		buf.WriteRune(c)
	}
	return buf.String()
}
