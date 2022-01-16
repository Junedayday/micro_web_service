package dao

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/Junedayday/micro_web_service/internal/gormer"
)

var (
	order = &gormer.Order{Name: "order1", Price: 1.1, CreateTime: time.Now(), UpdateTime: time.Now()}
	page, size = 2, 10
	updated = gormer.NewOrderOptions(&gormer.Order{Name: "test_name", UpdateTime: time.Now()}, gormer.OrderFieldName, gormer.OrderFieldUpdateTime)
	condition = gormer.NewOrderOptions(&gormer.Order{Price: 1.0}, gormer.OrderFieldPrice)
)

func TestOrderRepo_AddOrder(t *testing.T) {
	orderRepo := InitializeMockOrderRepo()
	err := orderRepo.AddOrder(context.Background(), order)
	assert.Nil(t, err)
}

func TestOrderRepo_QueryOrders(t *testing.T) {
	orderRepo := InitializeMockOrderRepo()
	_, err := orderRepo.QueryOrders(context.Background(), page, size, condition)
	assert.Nil(t, err)
}

func TestOrderRepo_CountOrders(t *testing.T) {
	orderRepo := InitializeMockOrderRepo()
	_ , err := orderRepo.CountOrders(context.Background(), condition)
	assert.Nil(t, err)
}

func TestOrderRepo_UpdateOrder(t *testing.T) {
	orderRepo := InitializeMockOrderRepo()
	err := orderRepo.UpdateOrder(context.Background(), updated, condition)
	assert.Nil(t, err)
}

func TestOrderRepo_DeleteOrder(t *testing.T) {
	orderRepo := InitializeMockOrderRepo()
	err := orderRepo.DeleteOrder(context.Background(), condition)
	assert.Nil(t, err)
}
