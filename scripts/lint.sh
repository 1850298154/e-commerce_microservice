#!/bin/bash

set -euo pipefail  # 使用更严格的错误处理：u标志检测未定义变量，pipefail保证管道中任何命令失败都返回错误

# 开启调试模式（可选）
# set -x

# 加载应用列表脚本
source scripts/list_app.sh

# 获取应用列表
get_app_list

# 遍历应用列表并执行 lint 检查
for app_path in "${app_list[@]}"; do
    echo "Processing application: ${app_path}"

    # 确保进入目录前检查是否存在
    if [[ -d "${app_path}" ]]; then
        pushd "${app_path}" > /dev/null  # 使用 pushd/popd 替代 cd，避免目录切换的混乱
        golangci-lint run --path-prefix=. --timeout=5m --config ../../.golangci.yml
        popd > /dev/null
    else
        echo "Directory ${app_path} does not exist. Skipping..."
    fi
done

echo "Linting complete for all applications."