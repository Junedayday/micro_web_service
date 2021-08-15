#！/bin/bash

# 第一次初始化请使用 buf beta mod init 指令

# 需要更新或安装buf相关组件，使用 下面命令
buf beta mod update

rm -rf gen/*
buf generate