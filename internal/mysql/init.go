package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func Init(user, password, ip string, port int, dbname string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, ip, port, dbname)
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}
