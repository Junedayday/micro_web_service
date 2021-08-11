# micro_web_service
build a micro web service for go users

## 手工安装

1. go语言版本 >=1.15
2. [buf工具](https://github.com/bufbuild/buf/releases) >=0.41.0
3. grpc-gateway相关二进制文件的安装
```shell
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

> 保证上述安装的结果都放在了PATH目录下，可执行
> 以上四个二进制程序的版本需要关注，尤其是带v2的
> PS: 有同学会遇到一些很奇怪的问题，往往是因为存在多个二进制程序及版本
