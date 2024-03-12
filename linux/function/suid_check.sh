#! /bin/bash

# 统计特殊权限
# 对系统中默认应该拥有SetUID权限的文件制作一张列表
# 定时检查有没有列表之外的文件被设置了SetUID权限

# 初始化
# -perm 安装权限查找。-4000对应的是SetUID权限，-2000对应的是SetGID权限
# -o 是逻辑或“or”的意思。并把命令搜索的结果放在/root/suid.list文件中
# find / -perm -4000 -o -perm -2000 > /root/suid.list
    

# Author: shenchao （E-mail: shenchao@lampbrother.net）
find / -perm -4000 -o -perm -2000 > /tmp/setuid.check
#搜索系统中所有拥有SetUID和SetGID权限的文件，并保存到临时目录中
for i in $(cat /tmp/setuid.check)
#循环，每次循环都取出临时文件中的文件名
do
    grep $i /root/suid.list > /dev/null
    #比对这个文件名是否在模板文件中
    if [ "$? " ! = "0" ]
    #检测上一条命令的返回值，如果不为0，则证明上一条命令报错
    then
        echo "$i isn't in listfile! " >> /root/suid_log_$(date +%F)
        #如果文件名不在模板文件中，则输出错误信息，并把报错写入日志中
    fi
done
rm -rf /tmp/setuid.check # 删除临时文件

# 测试
# chmod u+s /bin/vi # 手工给vi加入SetUID权限
# ./suidcheck.sh # 执行检测脚本
# cat suid_log_2013-01-20
# /bin/vi isn't in listfile! #报错了，vi不在模板文件中。代表vi被修改了SetUID权限