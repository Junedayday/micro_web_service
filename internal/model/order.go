package model

// Order针对的是 orders 表中的一行数据
// 在这里定义，是为了分离Model与Dao
type Order struct {
	Id    int64
	Name  string
	Price float32
}

// OrderFields 作为一个 数据库Order对象+fields字段的组合
// fields用来指定Order中的哪些字段生效
type OrderFields struct {
	Order  *Order
	Fields []string
}

type OrderRepository interface {
	AddOrder(order *Order) (err error)
	QueryOrders(pageNumber, pageSize int, condition *OrderFields) (orders []Order, err error)
	UpdateOrder(updated, condition *OrderFields) (err error)
}
