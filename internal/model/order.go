package model

import "github.com/Junedayday/micro_web_service/internal/gormer"

type OrderRepository interface {
	AddOrder(order *gormer.Order) (err error)
	QueryOrders(pageNumber, pageSize int, condition *gormer.OrderOptions) (orders []gormer.Order, err error)
	CountOrders(condition *gormer.OrderOptions) (count int64, err error)
	UpdateOrder(updated, condition *gormer.OrderOptions) (err error)
	DeleteOrder(condition *gormer.OrderOptions) (err error)
}
