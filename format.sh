#!/bin/bash
PACKAGE_PATH=github.com/Junedayday/micro_web_service

go fmt ./...
goimports -local=${PACKAGE_PATH} -w -l $(find . -type f -name '*.go' -not -path "./vendor/*")
