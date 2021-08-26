package dao

import (
	"time"

	"github.com/pkg/errors"

	"github.com/Junedayday/micro_web_service/internal/mysql"
)

type Order struct {
	Id         int64
	Name       string
	Price      float32
	CreateTime time.Time
}

func AddOrder(order *Order) (err error) {
	err = mysql.GormDB.Create(order).Error
	return
}

func QueryOrders(pageNumber, pageSize int, condition *Order, fields ...interface{}) (orders []Order, err error) {
	err = mysql.GormDB.
		Where(condition, fields...).
		Limit(pageSize).
		Offset((pageNumber - 1) * pageSize).
		Find(&orders).Error
	return
}

func UpdateOrder(order *Order, fields ...interface{}) (err error) {
	if len(fields) == 0 {
		return errors.New("update db must choose certain fields")
	}
	err = mysql.GormDB.
		Model(&Order{}).
		Select(fields[0], fields[1:]...).
		Updates(order).
		Error
	return
}
