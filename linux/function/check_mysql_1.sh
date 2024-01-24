#!/bin/sh
# 通过端口检查mysql状态
port=`netstat -lnt|grep 3306|wc -l`

if [ $port -ne 1 ];then
   echo "mysql is stop"
   /etc/init.d/mysqld start
else
   echo "mysql is starting"
fi