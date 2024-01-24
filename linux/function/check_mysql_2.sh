#/bin/sh
# 通过连接数据库检查mysql状态
mysql -uroot -pmysql123 -e "show databases;" >/dev/null 2>&1

if [ $? -eq 0 ];then
   echo "mysql is starting"
else
   echo "mysql is stop"
fi