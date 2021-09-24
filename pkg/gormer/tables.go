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
	TableName      string
	Name           string
	SmallCamelName string
	Columns        []FieldLevel
}

type FieldLevel struct {
	FieldName string
	FieldType string
	GormName  string
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

func Generate(db *sql.DB, table, structName string) (StructLevel, error) {
	var structData = StructLevel{
		TableName: table,
		Name:      camelCase(structName),
	}

	structData.SmallCamelName = string(unicode.ToLower(rune(structData.Name[0]))) + structData.Name[1:]

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
			columns = append(columns, FieldLevel{
				FieldName: camelCase(name),
				FieldType: fieldType(dataType),
				GormName:  name,
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
