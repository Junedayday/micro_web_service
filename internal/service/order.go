package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Junedayday/micro_web_service/internal/dao"
	"github.com/Junedayday/micro_web_service/internal/gormer"
	"github.com/Junedayday/micro_web_service/internal/model"
	"github.com/Junedayday/micro_web_service/internal/mysql"
	"github.com/Junedayday/micro_web_service/internal/zlog"
)

type OrderService struct {
	orderRepo model.OrderRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: dao.NewOrderRepo(mysql.GormDB),
	}
}

func (orderSvc *OrderService) List(ctx context.Context, pageNumber, pageSize int, condition *gormer.OrderOptions) ([]gormer.Order, int64, error) {
	zlog.WithTrace(ctx).Infof("page number is %d", pageNumber)

	orders, err := orderSvc.orderRepo.QueryOrders(pageNumber, pageSize, condition)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "OrderService List pageNumber %d pageSize %d condition %+v", pageNumber, pageSize, condition)
	}
	count, err := orderSvc.orderRepo.CountOrders(condition)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "OrderService Count condition %+v", condition)
	}

	return orders, count, nil
}

func (orderSvc *OrderService) Create(ctx context.Context, order *gormer.Order) error {
	err := orderSvc.orderRepo.AddOrder(order)
	if err != nil {
		return errors.Wrapf(err, "OrderService Create  order %+v", order)
	}
	return nil
}

func (orderSvc *OrderService) Update(ctx context.Context, updated, condition *gormer.OrderOptions) error {
	err := orderSvc.orderRepo.UpdateOrder(updated, condition)
	if err != nil {
		return errors.Wrapf(err, "OrderService Update updated %+v condition %+v", updated, condition)
	}
	return nil
}

func (orderSvc *OrderService) Delete(ctx context.Context, condition *gormer.OrderOptions) error {
	err := orderSvc.orderRepo.DeleteOrder(condition)
	if err != nil {
		return errors.Wrapf(err, "OrderService Delete condition %+v", condition)
	}
	return nil
}
