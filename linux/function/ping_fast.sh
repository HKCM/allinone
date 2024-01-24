#!/usr/bin/env bash
# 文件名: fast_ping.sh
# 用途：根据所在网络的实际情况修改网络地址192.168.0。
for ip in 192.168.0.{1..255} ;
do
    (
        ping $ip -c2 &> /dev/null ;

        if [ $? -eq 0 ];then
            echo $ip is alive
        fi
    )&
done
wait 

# 在for循环体中执行了多个后台进程，然后结束循环并终止脚本。
# wait命令会等待所有的子进程结束后再终止脚本。