package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/Junedayday/micro_web_service/internal/model"
)

/*
  gorm.io/gorm 指的是gorm V2版本，详细可参考 https://gorm.io/zh_CN/docs/v2_release_note.html
  github.com/jinzhu/gorm 一般指V1版本
*/

/*
CREATE TABLE orders
(
id bigint PRIMARY KEY AUTO_INCREMENT,
name varchar(255),
price decimal(15,3)
)
*/

type OrderRepo struct {
	db *gorm.DB
}

// 将gorm.DB作为一个参数，在初始化时赋值：方便测试时，放一个mock的db
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (repo *OrderRepo) AddOrder(order *model.Order) (err error) {
	err = repo.db.Create(order).Error
	return
}

func (repo *OrderRepo) QueryOrders(pageNumber, pageSize int, condition *model.OrderFields) (orders []model.Order, err error) {
	db := repo.db
	// condition非nil的话，追加条件
	if condition != nil {
		// 这里的field指定了order中生效的字段，这些字段会被放在SQL的where条件中
		db = db.Where(condition.Order, condition.Fields)
	}
	err = db.
		// Select("id","name").
		Limit(pageSize).
		Offset((pageNumber - 1) * pageSize).
		Find(&orders).Error
	return
}

func (repo *OrderRepo) UpdateOrder(updated, condition *model.OrderFields) (err error) {
	if updated == nil || len(updated.Fields) == 0 {
		return errors.New("update must choose certain fields")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}

	err = repo.db.
		Model(&model.Order{}).
		// 这里的field指定了被更新的where条件中的字段
		Where(condition.Order, condition.Fields).
		// 这里的field指定了order中被更新的字段
		Select(updated.Fields).
		Updates(updated.Order).
		Error
	return
}
