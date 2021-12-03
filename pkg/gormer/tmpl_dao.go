package main

const (
	daoHeader = `// Code generated by gormer. DO NOT EDIT.
package %s

import (
	"context"
	"time"
	
	"github.com/pkg/errors"
	"gorm.io/gorm"
	
	"%s/%s"
	"%s/%s"
)
`
	daoExtHeader = `package %s
	
// Implement ext method here
`
)

var (
	daoTmplRepo = `
type {{.StructName}}Repo struct {
	db *gorm.DB
}

func New{{.StructName}}Repo(db *gorm.DB) *{{.StructName}}Repo {
	return &{{.StructName}}Repo{db: db}
}

var _ model.{{.StructName}}Model = New{{.StructName}}Repo(nil)

`
	daoTmplAdd = `func (repo *{{.StructName}}Repo) Add{{.StructName}}(ctx context.Context, {{.StructSmallCamelName}} *gormer.{{.StructName}}) (err error) {
{{if ne .FieldCreateTime "" }}
    if {{.StructSmallCamelName}}.{{.FieldCreateTime}}.IsZero() {
		{{.StructSmallCamelName}}.{{.FieldCreateTime}} = time.Now()
	}
{{end}}
{{if ne .FieldUpdateTime "" }}
    if {{.StructSmallCamelName}}.{{.FieldUpdateTime}}.IsZero() {
		{{.StructSmallCamelName}}.{{.FieldUpdateTime}} = time.Now()
	}
{{end}}
	repo.db.WithContext(ctx).
		Table(gormer.{{.StructName}}TableName).
		Create({{.StructSmallCamelName}})
	err = repo.db.Error
	return
}

`
	daoTmplQuery = `func (repo *{{.StructName}}Repo) Query{{.StructName}}s(ctx context.Context, pageNumber, pageSize int, condition *gormer.{{.StructName}}Options) ({{.StructSmallCamelName}}s []gormer.{{.StructName}}, err error) {
	db := repo.db
	if condition != nil {
		db = db.Where(condition.{{.StructName}}, condition.Fields)
	}
{{if ne .FieldSoftDeleteKey "" }}
	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db.WithContext(ctx).
		Table(gormer.{{.StructName}}TableName).
		Limit(pageSize).
		Offset((pageNumber - 1) * pageSize).
		Find(&{{.StructSmallCamelName}}s)
	err = repo.db.Error
	return
}

`
	daoTmplCount = `func (repo *{{.StructName}}Repo) Count{{.StructName}}s(ctx context.Context, condition *gormer.{{.StructName}}Options) (count int64, err error) {
	db := repo.db
	if condition != nil {
		db = db.Where(condition.{{.StructName}}, condition.Fields)
	}
{{if ne .FieldSoftDeleteKey "" }}
	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db.WithContext(ctx).
		Table(gormer.{{.StructName}}TableName).
		Count(&count)
	err = repo.db.Error
	return
}

`
	daoTmplUpdate = `func (repo *{{.StructName}}Repo) Update{{.StructName}}(ctx context.Context, updated, condition *gormer.{{.StructName}}Options) (err error) {
	if updated == nil || len(updated.Fields) == 0 {
		return errors.New("update must choose certain fields")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}
{{if ne .FieldUpdateTime "" }}
    if updated.{{.StructName}}.{{.FieldUpdateTime}}.IsZero() {
		updated.{{.StructName}}.{{.FieldUpdateTime}} = time.Now()
		updated.Fields = append(updated.Fields, "{{.TableUpdateTime}}")
	}
{{end}}
	repo.db.WithContext(ctx).
		Table(gormer.{{.StructName}}TableName).
		Where(condition.{{.StructName}}, condition.Fields).
		Select(updated.Fields).
		Updates(updated.{{.StructName}})
	err = repo.db.Error
	return
}

`
	daoTmplDelete = `func (repo *{{.StructName}}Repo) Delete{{.StructName}}(ctx context.Context, condition *gormer.{{.StructName}}Options) (err error) {
	if condition == nil {
		return errors.New("delete must include where condition")
	}

	repo.db.WithContext(ctx).
        Table(gormer.{{.StructName}}TableName).
		Where(condition.{{.StructName}}, condition.Fields).
{{if eq .FieldSoftDeleteKey "" }} Delete(&gormer.{{.StructName}}{})
{{ else }}  {{if eq .FieldUpdateTime "" }}
				Select("{{.TableSoftDeleteKey}}").
				Updates(&gormer.{{.StructName}}{
					{{.FieldSoftDeleteKey}}:{{.TableSoftDeleteValue}},
				})
            {{ else }}
                Select("{{.TableSoftDeleteKey}}","{{.TableUpdateTime}}").
				Updates(&gormer.{{.StructName}}{
					{{.FieldSoftDeleteKey}}:{{.TableSoftDeleteValue}},
					{{.FieldUpdateTime}} : time.Now(),
				})
            {{ end }}
{{ end }}
	err = repo.db.Error
	return
}

`
	daoTmpl = daoTmplRepo + daoTmplAdd + daoTmplQuery + daoTmplCount + daoTmplUpdate + daoTmplDelete
)
