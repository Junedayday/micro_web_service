package main

import (
	"bytes"
	"html/template"
)

var gormerTmpl = `
// Table Level Info
const {{.Name}}TableName = "{{.TableName}}"

// Field Level Info
type {{.Name}}Field string
const (
{{range $item := .Columns}}
    {{$.Name}}Field{{$item.FieldName}} {{$.Name}}Field = "{{$item.GormName}}" {{end}}
)

var {{$.Name}}FieldAll = []{{$.Name}}Field{ {{range $k,$item := .Columns}}"{{$item.GormName}}", {{end}}}

// Kernel struct for table for one row
type {{.Name}} struct { {{range $item := .Columns}}
	{{$item.FieldName}}	{{$item.FieldType}}	` + "`" + `gorm:"column:{{$item.GormName}}"` + "`" + ` {{end}}
}

// Kernel struct for table operation
type {{.Name}}Options struct {
    {{.Name}} *{{.Name}}
    Fields []string
}

// Match: case insensitive
var {{$.Name}}FieldMap = map[string]string{
{{range $item := .Columns}}"{{$item.FieldName}}":"{{$item.GormName}}","{{$item.GormName}}":"{{$item.GormName}}",
{{end}}
}

func New{{.Name}}Options(target *{{.Name}}, fields ...{{$.Name}}Field) *{{.Name}}Options{
    options := &{{.Name}}Options{
        {{.Name}}: target,
        Fields: make([]string, len(fields)),
    }
    for index, field := range fields {
        options.Fields[index] = string(field)
    }
    return options
}

func New{{.Name}}OptionsAll(target *{{.Name}}) *{{.Name}}Options{
    return New{{.Name}}Options(target, {{$.Name}}FieldAll...)
}

func New{{.Name}}OptionsRawString(target *{{.Name}}, fields ...string) *{{.Name}}Options{
    options := &{{.Name}}Options{
        {{.Name}}: target,
    }
    for _, field := range fields {
        if f,ok := {{$.Name}}FieldMap[field];ok {
             options.Fields = append(options.Fields, f)
        }
    }
    return options
}
`

var daoTmpl = `
type {{.Name}}Repo struct {
	db *gorm.DB
}

func New{{.Name}}Repo(db *gorm.DB) *{{.Name}}Repo {
	return &{{.Name}}Repo{db: db}
}

func (repo *{{.Name}}Repo) Add{{.Name}}({{.SmallCamelName}} *gormer.{{.Name}}) (err error) {
	err = repo.db.Create({{.SmallCamelName}}).Error
	return
}

func (repo *{{.Name}}Repo) Query{{.Name}}s(pageNumber, pageSize int, condition *gormer.{{.Name}}Options) ({{.SmallCamelName}}s []gormer.{{.Name}}, err error) {
	db := repo.db
	if condition != nil {
		db = db.Where(condition.{{.Name}}, condition.Fields)
	}
	err = db.
		Limit(pageSize).
		Offset((pageNumber - 1) * pageSize).
		Find(&{{.SmallCamelName}}s).Error
	return
}

func (repo *{{.Name}}Repo) Update{{.Name}}(updated, condition *gormer.{{.Name}}Options) (err error) {
	if updated == nil || len(updated.Fields) == 0 {
		return errors.New("update must choose certain fields")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}

	err = repo.db.
		Model(&gormer.{{.Name}}{}).
		Where(condition.{{.Name}}, condition.Fields).
		Select(updated.Fields).
		Updates(updated.{{.Name}}).
		Error
	return
}
`

func parseToGormerTmpl(structData StructLevel) (string, error) {
	tmpl, err := template.New("t").Parse(gormerTmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, structData)
	return buf.String(), nil
}

func parseToDaoTmpl(structData StructLevel) (string, error) {
	tmpl, err := template.New("t").Parse(daoTmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, structData)
	return buf.String(), nil
}
