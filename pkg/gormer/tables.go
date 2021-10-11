package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

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

type StructLevel struct {
	// table -> struct
	TableName            string
	StructName           string
	StructSmallCamelName string

	// table column -> struct field
	Columns []FieldLevel

	// create time
	TableCreateTime string
	FieldCreateTime string

	// update time
	TableUpdateTime string
	FieldUpdateTime string

	// soft delete
	TableSoftDeleteKey   string
	TableSoftDeleteValue int
	FieldSoftDeleteKey   string
}

type FieldLevel struct {
	FieldName string
	FieldType string
	// gorm tag for field
	GormName string
	// comment from create table sql
	Comment string
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
		TableName:  table,
		StructName: camelCase(matchInfo.GoStruct),
	}

	structData.StructSmallCamelName = string(unicode.ToLower(rune(structData.StructName[0]))) + structData.StructName[1:]

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
