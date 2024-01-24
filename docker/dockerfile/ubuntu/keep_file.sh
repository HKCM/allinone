#!/bin/bash

echo "Start check.."
# 要保留的文件数量
retain_count=4

# 文件名模式
file_pattern="*.log"

# 查找并按文件创建日期排序
files=($(ls -ltr $file_pattern | awk '{print $9}'))

# 计算要删除的文件数量
count=$((${#files[@]} - $retain_count))

# 删除多余的文件
if [ $count -gt 0 ]; then
    for ((i=0; i<$count; i++)); do
        rm "${files[$i]}"
        echo "Removed ${files[$i]}"
    done
fi