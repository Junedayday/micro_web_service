package main

type StructLevel struct {
	// table -> struct
	TableName string

	StructName StructName

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

	GenQueries []GenQuery
}

// 复数Plural
// 单数Singular
type StructName struct {
	// first letter upper/lower
	UpperS, UpperP string // Order,Orders
	LowerS, LowerP string // order,orders
}

type FieldLevel struct {
	FieldName string
	FieldType string
	// gorm tag for field
	GormName string
	// comment from create table sql
	Comment string
}
