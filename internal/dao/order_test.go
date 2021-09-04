package dao

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 注意，我们使用的是gorm 2.0，网上很多例子其实是针对1.0的
var (
	DB   *gorm.DB
	mock sqlmock.Sqlmock
)

// TestMain是在当前package下，最先运行的一个函数，常用于初始化
func TestMain(m *testing.M) {
	var (
		db  *sql.DB
		err error
	)

	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic(err)
	}

	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// m.Run 是真正调用下面各个Test函数的入口
	os.Exit(m.Run())
}

/*
  sqlmock 对语法限制比较大，下面的sql语句必须精确匹配（包括符号和空格）
*/

func TestOrderRepo_AddOrder(t *testing.T) {
	var order = &Order{Name: "order1", Price: 1.1}
	orderRepo := NewOrderRepo(DB)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `orders` (`name`,`price`) VALUES (?,?)").
		WithArgs(order.Name, order.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := orderRepo.AddOrder(order)
	assert.Nil(t, err)
}

func TestOrderRepo_QueryOrders(t *testing.T) {
	var orders = []Order{
		{1, "name1", 1.0},
		{2, "name2", 1.0},
	}
	page, size := 2, 10
	orderRepo := NewOrderRepo(DB)
	condition := NewOrderFields(&Order{Price: 1.0}, []interface{}{"price"})

	mock.ExpectQuery(
		"SELECT * FROM `orders` WHERE `orders`.`price` = ? LIMIT 10 OFFSET 10").
		WithArgs(condition.order.Price).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow(orders[0].Id, orders[0].Name, orders[0].Price).
				AddRow(orders[1].Id, orders[1].Name, orders[1].Price))

	ret, err := orderRepo.QueryOrders(page, size, condition)
	assert.Nil(t, err)
	assert.Equal(t, orders, ret)
}

func TestOrderRepo_UpdateOrder(t *testing.T) {
	orderRepo := NewOrderRepo(DB)
	// 表示要更新的字段为Order对象中的id,name两个字段
	updated := NewOrderFields(&Order{Id: 1, Name: "test_name"}, []interface{}{"id", "name"})
	// 表示更新的条件为Order对象中的price字段
	condition := NewOrderFields(&Order{Price: 1.0}, []interface{}{"price"})

	mock.ExpectBegin()
	mock.ExpectExec(
		"UPDATE `orders` SET `id`=?,`name`=? WHERE `orders`.`price` = ?").
		WithArgs(updated.order.Id, updated.order.Name, condition.order.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := orderRepo.UpdateOrder(updated, condition)
	assert.Nil(t, err)
}
