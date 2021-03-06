// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package dao

import (
	"github.com/Junedayday/micro_web_service/internal/mysql"
)

// Injectors from wire.go:

func InitializeMockOrderRepo() *OrderRepo {
	db := mysql.NewMockDB()
	orderRepo := NewOrderRepo(db)
	return orderRepo
}
