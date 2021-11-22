package main

const (
	daoBaseHeader = `// Code generated by gormer. DO NOT EDIT.
package %s

import (
	"context"

	trace "git.xiaojukeji.com/lego/context-go"
	"gorm.io/gorm"
	
	"%s"
)
`
)

var (
	daoBaseTmpl = `
func getError(ctx context.Context, db *gorm.DB) (err error) {
	err = db.Error
{{if eq .LogOn true }}
	if err != nil {
		// add you error log here - db.Statement.SQL.String(), db.Statement.Vars
		return
	}
	// add you info log here - db.Statement.SQL.String(), db.Statement.Vars
{{end}}
	return
}
`
)
