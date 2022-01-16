package mysql

import (
	"fmt"
	
	`github.com/DATA-DOG/go-sqlmock`
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Junedayday/micro_web_service/internal/zlog"
)

var GormDB *gorm.DB

func InitGorm(user, password, addr string, dbname string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, addr, dbname)
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 结束后
	_ = GormDB.Callback().Create().After("gorm:after_create").Register(callBackLogName, afterLog)
	_ = GormDB.Callback().Query().After("gorm:after_query").Register(callBackLogName, afterLog)
	_ = GormDB.Callback().Delete().After("gorm:after_delete").Register(callBackLogName, afterLog)
	_ = GormDB.Callback().Update().After("gorm:after_update").Register(callBackLogName, afterLog)
	_ = GormDB.Callback().Row().After("gorm:row").Register(callBackLogName, afterLog)
	_ = GormDB.Callback().Raw().After("gorm:raw").Register(callBackLogName, afterLog)
	return
}

const callBackLogName = "zlog"

func afterLog(db *gorm.DB) {
	err := db.Error
	ctx := db.Statement.Context

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	if err != nil {
		zlog.WithTrace(ctx).Errorf("sql=%s || error=%v", sql, err)
		return
	}
	zlog.WithTrace(ctx).Infof("sql=%s", sql)
}


func NewMockDB() *gorm.DB{
	db,_, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic(err)
	}
	
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	return gormDB
}