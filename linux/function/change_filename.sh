#!/usr/bin/env bash

# author: karl
# version: 1
# date: 2023-12-31
# 批量修改文件名 批量备份

suffix=$(date +%Y%m%d)

for f in $(find ./ -type f -name "*.txt");do
    echo "Backup file ${f}"
    cp ${f} ${f}_${suffix}
done