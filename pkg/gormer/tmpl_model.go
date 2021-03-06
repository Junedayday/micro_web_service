package main

const (
	modelHeader = `// Code generated by gormer. DO NOT EDIT.
package %s

import (
	"context"
	"time"
	
	"%s/%s"
)
`

	modelExtHeader = `package %s

// write you method here
`
)

var (
	modelTmpl = `
type {{.StructName.UpperS}}Model interface {
	Add{{.StructName.UpperS}}(ctx context.Context, {{.StructName.LowerS}} *gormer.{{.StructName.UpperS}}) (err error)
	Add{{.StructName.UpperP}}(ctx context.Context, {{.StructName.LowerP}} []*gormer.{{.StructName.UpperS}}) (err error)
	Query{{.StructName.UpperS}}s(ctx context.Context, pageNumber, pageSize int, condition *gormer.{{.StructName.UpperS}}Options) ({{.StructName.LowerS}}s []gormer.{{.StructName.UpperS}}, err error)
	Count{{.StructName.UpperS}}s(ctx context.Context, condition *gormer.{{.StructName.UpperS}}Options) (count int64, err error)
	Update{{.StructName.UpperS}}(ctx context.Context, updated, condition *gormer.{{.StructName.UpperS}}Options) (err error)
	Delete{{.StructName.UpperS}}(ctx context.Context, condition *gormer.{{.StructName.UpperS}}Options) (err error)
	
	// Defined in genQueries
	
	{{range $item := .GenQueries}} // Query{{$item.Method}} {{$item.Desc}}
	Query{{$item.Method}}(ctx context.Context,{{range $match := $item.Args}} {{$match.Name}} {{$match.Type}}, {{end}}pageNumber, pageSize int, condition *gormer.{{$.StructName.UpperS}}Options) ({{$.StructName.LowerP}} []gormer.{{$.StructName.UpperS}}, err error)
	// Count{{$item.Method}} {{$item.Desc}}
	Count{{$item.Method}}(ctx context.Context,{{range $match := $item.Args}} {{$match.Name}} {{$match.Type}}, {{end}}condition *gormer.{{$.StructName.UpperS}}Options) (count int64, err error)
	{{end}}
	// Implement Your Method in ext model
	{{.StructName.UpperS}}ExtModel
}
`
	modelExtTmpl = `
type %sExtModel interface {
}
`
)
