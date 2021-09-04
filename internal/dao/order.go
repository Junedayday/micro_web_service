package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

/*
  gorm.io/gorm 指的是gorm V2版本，详细可参考 https://gorm.io/zh_CN/docs/v2_release_note.html
  github.com/jinzhu/gorm 一般指V1版本
*/

type OrderRepo struct {
	db *gorm.DB
}

// 将gorm.DB作为一个参数，在初始化时赋值：方便测试时，放一个mock的db
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Order针对的是 orders 表中的一行数据
type Order struct {
	Id    int64
	Name  string
	Price float32
}

// OrderFields 作为一个 数据库Order对象+fields字段的组合
// fields用来指定Order中的哪些字段生效
type OrderFields struct {
	order  *Order
	fields []interface{}
}

func NewOrderFields(order *Order, fields []interface{}) *OrderFields {
	return &OrderFields{
		order:  order,
		fields: fields,
	}
}

func (repo *OrderRepo) AddOrder(order *Order) (err error) {
	err = repo.db.Create(order).Error
	return
}

func (repo *OrderRepo) QueryOrders(pageNumber, pageSize int, condition *OrderFields) (orders []Order, err error) {
	db := repo.db
	// condition非nil的话，追加条件
	if condition != nil {
		// 这里的field指定了order中生效的字段，这些字段会被放在SQL的where条件中
		db = db.Where(condition.order, condition.fields...)
	}
	err = db.
		Limit(pageSize).
		Offset((pageNumber - 1) * pageSize).
		Find(&orders).Error
	return
}

func (repo *OrderRepo) UpdateOrder(updated, condition *OrderFields) (err error) {
	if updated == nil || len(updated.fields) == 0 {
		return errors.New("update must choose certain fields")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}

	err = repo.db.
		Model(&Order{}).
		// 这里的field指定了order中被更新的字段
		Select(updated.fields[0], updated.fields[1:]...).
		// 这里的field指定了被更新的where条件中的字段
		Where(condition.order, condition.fields...).
		Updates(updated.order).
		Error
	return
}
