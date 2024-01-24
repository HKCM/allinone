#!/usr/bin/env bash
set -e
# author: karl
# date; 2023-11-24
# 关键词: sysload monitor
# description: 获取系统负载,如果负载高于20,则记录系统信息

LOG_LOCATION=/opt/logs
LOG_LEVEL=2
THRESHOLD=20
source shell_log.sh

while true
do
    load=$(uptime |awk -F 'average:' '{print $2}'|cut -d',' -f1|sed 's/ //g'|cut -d. -f1)
    # 如果系统load大于THRESHOLD
    if [ ${load} -ge ${THRESHOLD} ];then
        top -bn1 | head -n 100 > ${LOG_LOCATION}/top.$(date +%Y%m%d%H%M%S) # 以非交互式运行top
        vmstat 1 10 > ${LOG_LOCATION}/vmstat.$(date +%Y%m%d%H%M%S)
        ss -an > ${LOG_LOCATION}/ss.$(date +%Y%m%d%H%M%S)
    fi
    sleep 20

    # 找到30天前的数据然后删除
    find ${LOG_LOCATION} \(-name "top*" -0 -name "vmstat*" -0 -name "ss*" \) -mtime +30 | xargs rm -f
done