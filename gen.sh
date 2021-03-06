#！/bin/bash

# 第一次初始化请使用 buf beta mod init 指令

# 需要更新或安装buf相关组件，使用 下面命令
# buf mod update

rm -rf gen/idl/*
buf mod update
buf generate

# mock install guide: https://github.com/golang/mock
# mockgen -destination internal/mock/mock_order.go -package order -source=internal/model/order.go OrderRepository

# gormer part
go get github.com/Junedayday/micro_web_service/pkg/gormer
gormer -c ./

./format.sh

./swagger.sh
