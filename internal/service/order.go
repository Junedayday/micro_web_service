package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Junedayday/micro_web_service/internal/dao"
	"github.com/Junedayday/micro_web_service/internal/model"
	"github.com/Junedayday/micro_web_service/internal/mysql"
)

type OrderService struct {
	orderRepo model.OrderRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: dao.NewOrderRepo(mysql.GormDB),
	}
}

func (orderSvc *OrderService) List(ctx context.Context, pageNumber, pageSize int, condition *model.OrderFields) ([]model.Order, error) {
	orders, err := orderSvc.orderRepo.QueryOrders(pageNumber, pageSize, condition)
	if err != nil {
		return nil, errors.Wrapf(err, "OrderService List pageNumber %d pageSize %d", pageNumber, pageSize)
	}
	return orders, nil
}

func (orderSvc *OrderService) Create(ctx context.Context, order *model.Order) error {
	err := orderSvc.orderRepo.AddOrder(order)
	if err != nil {
		return errors.Wrapf(err, "OrderService Create  order %+v", order)
	}
	return nil
}

func (orderSvc *OrderService) Update(ctx context.Context, updated, condition *model.OrderFields) error {
	err := orderSvc.orderRepo.UpdateOrder(updated, condition)
	if err != nil {
		return errors.Wrapf(err, "OrderService Update updated %+v condition %+v", updated, condition)
	}
	return nil
}
