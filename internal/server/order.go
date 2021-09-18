package server

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Junedayday/micro_web_service/gen/idl/order"
	"github.com/Junedayday/micro_web_service/internal/model"
	"github.com/Junedayday/micro_web_service/internal/service"
)

func (s *Server) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	orders, err := service.NewOrderService().List(ctx, int(req.PageNumber), int(req.PageSize), nil)
	if err != nil {
		return nil, err
	}
	resp := new(order.ListOrdersResponse)
	resp.Orders = make([]*order.Order, len(orders))
	for k, v := range orders {
		resp.Orders[k] = &order.Order{
			Id:    v.Id,
			Name:  v.Name,
			Price: v.Price,
		}
	}
	return resp, nil
}

func (s *Server) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
	mOrder := &model.Order{
		Id:    req.Order.Id,
		Name:  req.Order.Name,
		Price: req.Order.Price,
	}
	err := service.NewOrderService().Create(ctx, mOrder)
	if err != nil {
		return nil, err
	}

	return &order.Order{
		Id:    mOrder.Id,
		Name:  mOrder.Name,
		Price: mOrder.Price,
	}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*emptypb.Empty, error) {
	updated := &model.OrderFields{
		Order: &model.Order{
			Name:  req.Order.Name,
			Price: req.Order.Price,
		},
		Fields: req.UpdateMask.Paths,
	}
	condition := &model.OrderFields{
		Order: &model.Order{
			Id: req.Order.Id,
		},
		Fields: []string{"id"},
	}

	err := service.NewOrderService().Update(ctx, updated, condition)
	return &emptypb.Empty{}, err
}

func (s *Server) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.Order, error) {
	condition := &model.OrderFields{
		Order:  &model.Order{Name: req.Name},
		Fields: []string{"name"},
	}
	orders, err := service.NewOrderService().List(ctx, 0, 1, condition)
	if err != nil {
		return nil, err
	} else if len(orders) == 0 {
		return nil, errors.New("no order matched")
	}
	return &order.Order{
		Id:    orders[0].Id,
		Name:  orders[0].Name,
		Price: orders[0].Price,
	}, nil
}

func (s *Server) DeleteBook(ctx context.Context, req *order.DeleteBookRequest) (*emptypb.Empty, error) {
	condition := &model.OrderFields{
		Order:  &model.Order{Name: req.Name},
		Fields: []string{"name"},
	}

	// TODO soft delete
	updated := &model.OrderFields{
		Order:  &model.Order{},
		Fields: []string{},
	}

	return &emptypb.Empty{}, service.NewOrderService().Update(ctx, updated, condition)
}
