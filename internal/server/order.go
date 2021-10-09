package server

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Junedayday/micro_web_service/gen/idl/order"
	"github.com/Junedayday/micro_web_service/internal/gormer"
	"github.com/Junedayday/micro_web_service/internal/service"
)

func (s *Server) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	orders, count, err := service.NewOrderService().List(ctx, int(req.PageNumber), int(req.PageSize), nil)
	if err != nil {
		return nil, err
	}
	resp := new(order.ListOrdersResponse)
	resp.Count = int32(count)
	resp.Orders = make([]*order.Order, len(orders))
	for k, v := range orders {
		resp.Orders[k] = &order.Order{
			Id:         v.Id,
			Name:       v.Name,
			Price:      float32(v.Price),
			CreateTime: timestamppb.New(v.CreateTime),
			UpdateTime: timestamppb.New(v.UpdateTime),
		}
	}
	return resp, nil
}

func (s *Server) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
	mOrder := &gormer.Order{
		Id:    req.Order.Id,
		Name:  req.Order.Name,
		Price: float64(req.Order.Price),
	}
	err := service.NewOrderService().Create(ctx, mOrder)
	if err != nil {
		return nil, err
	}

	return &order.Order{
		Id:         mOrder.Id,
		Name:       mOrder.Name,
		Price:      float32(mOrder.Price),
		CreateTime: timestamppb.New(mOrder.CreateTime),
		UpdateTime: timestamppb.New(mOrder.UpdateTime),
	}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*emptypb.Empty, error) {
	updateOrder := &gormer.Order{
		Name:  req.Order.Name,
		Price: float64(req.Order.Price),
	}
	updated := gormer.NewOrderOptionsRawString(updateOrder, req.UpdateMask.Paths...)

	condOrder := &gormer.Order{
		Id: req.Order.Id,
	}
	condition := gormer.NewOrderOptions(condOrder, gormer.OrderFieldId)

	err := service.NewOrderService().Update(ctx, updated, condition)
	return &emptypb.Empty{}, err
}

func (s *Server) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.Order, error) {
	condOrder := &gormer.Order{
		Name: req.Name,
	}
	condition := gormer.NewOrderOptions(condOrder, gormer.OrderFieldName)

	orders, _, err := service.NewOrderService().List(ctx, 0, 1, condition)
	if err != nil {
		return nil, err
	} else if len(orders) == 0 {
		return nil, errors.New("no order matched")
	}
	return &order.Order{
		Id:         orders[0].Id,
		Name:       orders[0].Name,
		Price:      float32(orders[0].Price),
		CreateTime: timestamppb.New(orders[0].CreateTime),
		UpdateTime: timestamppb.New(orders[0].UpdateTime),
	}, nil
}

func (s *Server) DeleteBook(ctx context.Context, req *order.DeleteOrderRequest) (*emptypb.Empty, error) {
	condOrder := &gormer.Order{
		Name: req.Name,
	}
	condition := gormer.NewOrderOptions(condOrder, gormer.OrderFieldName)

	return &emptypb.Empty{}, service.NewOrderService().Delete(ctx, condition)
}
