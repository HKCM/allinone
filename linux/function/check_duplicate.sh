#!/bin/bash

# 获取当前文件夹下所有文件的列表
files=$(find . -maxdepth 1 -type f)

# 创建一个关联数组用于存储文件大小和 MD5 值的组合
declare -A file_info

# 遍历文件
for file in $files; do
    # 获取文件大小
    size=$(stat -c%s "$file")

    # 获取文件的 MD5 值
    md5=$(md5sum "$file" | awk '{ print $1 }')

    # 检查是否存在相同大小的文件
    if [[ -n ${file_info["$size"]} ]]; then
        # 如果存在相同大小的文件，比较 MD5 值
        if [[ ${file_info["$size"]} == $md5 ]]; then
            echo "Duplicate file: $file"
        fi
    else
        # 如果不存在相同大小的文件，则将当前文件的信息存储到关联数组中
        file_info["$size"]=$md5
    fi
done

echo "Duplicate check completed."