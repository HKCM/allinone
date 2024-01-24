#!/usr/bin/env bash
# keyword: memory 检查内存 
# 获取可用内存（以兆字节为单位）
available_memory=$(free -m | awk '/^Mem:/ {print $7}')

# 判断可用内存是否小于 1000
if [ "$available_memory" -lt 1000 ]; then
    # 可用内存小于 1000，运行脚本 /tmp/test.sh
    /bin/bash /tmp/test.sh
else
    # 可用内存大于或等于 1000，输出一条消息
    echo "Available memory is sufficient."
fi