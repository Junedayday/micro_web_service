#!/bin/bash

# 请先保证swagger的安装，可参考 https://goswagger.io/install.html
# 合并swagger文档
swagger mixin gen/openapiv2/idl/*/*.json -o gen/swagger.json
# 删除原始文档
rm -rf gen/openapiv2

# 服务端运行的docker
# https://hub.docker.com/r/redocly/redoc/
docker run --name swagger -it --rm -d -p 80:80 -v $(pwd)/gen/swagger.json:/usr/share/nginx/html/swagger.json -e SPEC_URL=swagger.json redocly/redoc

# 将机器的sshkey保存到公有云机器后，通过scp放过去
# scp -P 22 gen/swagger.json root@xx.xx.xx.xx:/{target_path}
