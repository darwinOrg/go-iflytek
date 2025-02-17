#!/bin/bash

# 读取当前版本号
current_version=`git tag | sort -V | tail -n 1`
echo "当前版本号: $current_version"

# 使用正则表达式匹配版本号，并提取x部分
if [[ $current_version =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    major=${BASH_REMATCH[1]}
    minor=${BASH_REMATCH[2]}
    patch=${BASH_REMATCH[3]}

    # 增加补丁版本号
    new_patch=$((patch + 1))

    # 构建新的版本号
    new_version="v${major}.${minor}.${new_patch}"
    git tag $new_version
    git push --tags
    go list -m github.com/darwinOrg/go-iflytek
    echo "版本已从 $current_version 升级到 $new_version"
else
    echo "无法解析版本号: $current_version"
    exit 1
fi