#！/bin/bash

# 第一次初始化请使用 buf beta mod init 指令

rm -rf gen/*
buf beta mod update
buf generate