//+build wireinject

package dao

import (
	"github.com/Junedayday/micro_web_service/internal/mysql"

	"github.com/google/wire"
)

func InitializeMockOrderRepo() *OrderRepo {
	wire.Build(NewOrderRepo, mysql.NewMockDB)
	return &OrderRepo{}
}
