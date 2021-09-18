# micro_web_service

**Build a micro web service for go users.**

本框架的初衷是**分享并理解常用的Go语言工具**，其次才是做一个易用的框架。方便使用者有任何定制修改时，可以快速fork并自己实现。

## 目录介绍

- **idl** 用`protobuf`定义的IDL文件，用于生成`Go`以及其它语言的数据结构
- **gen** 从idl生成的文件，不允许手动修改
- **internal** 项目内部的重要实现

## 手工安装

1. go语言版本 >=1.15
2. [buf工具](https://github.com/bufbuild/buf/releases) >=0.54.1
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

## Docker镜像

```bash
docker pull mysql:5.7
sudo docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
```