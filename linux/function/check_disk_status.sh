#!/usr/bin/env bash

# author: karl
# version: 1
# date: 2023-12-31
# 磁盘可用性 磁盘状态检测

for mount_point in $(df |sed '1d' |grep -v 'tmpfs' |awk '{print $NF}')
do
    touch ${mount_point}/testfile && rm -f ${mount_point}/testfile
    if [ $? -ne 0 ];then
        echo "${mount_point} 读写有问题"
    else 
        echo "${mount_point} 读写没问题"
    fi
done