#!/bin/bash

# 定义启动目录
directories=(
    "./app/cart"
    "./app/product"
    "./app/order"
    "./app/payment"
)

# 遍历目录并启动进程
for dir in "${directories[@]}"; do
    if [ -d "$dir" ]; then
        echo "Starting process in directory: $dir"
        (cd "$dir" && go run .) &
        # cd ./app/cart     && go run .
        # cd ./app/product  && go run .
        # cd ./app/order    && go run .
        # cd ./app/payment  && go run .
    else
        echo "Directory $dir does not exist"
    fi
done

# 等待所有后台进程完成
wait