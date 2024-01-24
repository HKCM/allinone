CREATE数据库
```shell
#!/bin/bash
MYSQL=`which mysql`
$MYSQL -u mysql -p <<EOF
CREAT TABLE people(name VARCHAR(20),sex CHAR(1),birth DATE,addr VARCHAR(20));
SHOW TABLE;
EOF
```

创建用户
```shell
#!/bin/bash
MYSQL=`which mysql`
 
$MYSQL -u root <<EOF
GRANT SELECT ON test.* TO 'username'@'%' IDENTIFIED BY 'bsAhxbeT9UqiVaaL';
EOF
 
echo "$?
```


INSERT 数据
```shell
#!/bin/bash
MYSQL=`which mysql`
 
if [ $# -ne 4 ]
then
        echo "Usage:insert name sex birth addr"
else
        statement="INSERT INTO people values ('$1','$2','$3','$4');"
        $MYSQL database -u mysql -p <<EOF
        $statement
        EOF
        if [ $? -eq 0 ]
        then
                echo "Data insert sucessful"
        else
                echo "Something wrong"
        fi
fi
```

### 编写dockerfile
```shell
vi Dockerfile
# 输入以下内容:
from mysql:5.7
COPY my.cnf /etc/mysql/conf.d/mysqlutf8.cnf
CMD ["mysqld", "--character-set-server=utf8", "--collation-server=utf8_unicode_ci"]
```

### 通过Dockerfile构建Docker镜像
```
docker build -t mysql5.7utf8 .
```

### 启动新的docker

```shell
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=mysql -itd mysql5.7utf8:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
# or

docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=mysql -itd mysql5.7utf8:latest
```

### 链接
```
mysql -h127.0.0.1 -uroot -P3306 -p 
```

### 查看字符集
```sql
mysql> show variables like '%char%'; 
```

如果忘记了密码可以用`docker inspect`查看环境变量

# 脚本


### MySQL连接信息

```bash
MYSQL_USER="admin"
MYSQL_PASSWORD="12345678"
MYSQL_DATABASE="test1"
MYSQL_TABLE="test_table"
export MYSQL_PWD=12345678
MYSQL_HOST="database-1-old1.cluster-cxq6crkmyuh7.ap-northeast-1.rds.amazonaws.com"
MYSQL_HOST_RO="database-1-old1.cluster-ro-cxq6crkmyuh7.ap-northeast-1.rds.amazonaws.com"


# 创建测试表
mysql -h "$MYSQL_HOST" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "CREATE TABLE IF NOT EXISTS "$MYSQL_TABLE" (id INT AUTO_INCREMENT PRIMARY KEY, data VARCHAR(255));" "$MYSQL_DATABASE"

# 循环插入数据
while true
do
    data=$(date +"%Y-%m-%d %H:%M:%S")
    echo "${data} Insert Data..."
    mysql --connect-timeout=3 -h "$MYSQL_HOST" -u "$MYSQL_USER" -e "INSERT INTO "$MYSQL_TABLE" (data) VALUES ('$data');" "$MYSQL_DATABASE"
    sleep 5
done

# 循环查询数据
# -N 禁用了列名的显示
# -e 用于执行SQL查询
while true; do
    latest_record_sql="SELECT id, data FROM "$MYSQL_TABLE" ORDER BY id DESC LIMIT 1;"
    latest_record=$(mysql --connect-timeout=3 -h "$MYSQL_HOST_RO" -u "$MYSQL_USER" -N -e "$latest_record_sql" "$MYSQL_DATABASE")
    data=$(date +"%Y-%m-%d %H:%M:%S")
    echo "${data} Latest Record: ${latest_record}"
    sleep 5
done
```

## mysqldump

mysqldump客户端工具用来备份数据库或在不同数据库之间进行数据迁移。备份内容包含创建表,及插入表的sQL语句。

```
连接选项:

-u,--user=name 指定用户名
-p,--password[=name] 指定密码
-h,--host=name 指定服务器ip或域名
-P,-port=# 指定连接端口

备份选项:
    --all-databases:备份所有数据库
    --databases db1 db2:备份指定的数据库
    --single-transaction:对事务引擎执行热备
    --flush-logs:更新二进制日志文件
    --master-data=2
        1:每备份一个库就生成一个新的二进制文件(默认)
        2:只生成一个新的二进制文件
    --quick:在备份大表时指定该选项

输出选项:

--add-drop-database 在每个数据库创建语句前加上drop database语句
--add-drop-table 在每个表创建语句前加上drop table语句,默认开启；不开启(-skip-add-drop-table)
-n, --no-create-db 不包含数据库的创建语句
-t, --no-create-info 不包含数据表的创建语句
-d, --no-data 不包含数据
-T, --tab=name 自动生成两个文件:一个sql文件,创建表结构的语句；一个txt文件,数据文件
```

```shell
# 备份test数据库
# sql中先删除表(IF EXISTS) 然后在创建表
mysqldump -uroot -p123456 test > test_back.sql

# 导出多张表
mysqldump -uroot -p123456 --databases test --tables t1 t2>two.sql
	
# -t 只导出表数据不导表结构,添加“-t”命令参数
mysqldump -uroot -p123456 -t test > test_back.sql

# -d 只导出表结构不导表数据,添加“-d”命令参数
mysqldump -uroot -h127.0.0.1 -p123456 -P3306 -d test > test_back.sql

mysql -uroot -p123456 -e "show variables like '%secure_file_pric%'"
# -T 自动生成两个文件:一个sql文件,创建表结构的语句；一个txt文件,数据文件
mysqldump -uroot -p123456 -T /var/lib/mysql-files/ test emp

#命令行导入
mysql> use test;
mysql> source /home/test/database.sql
```

备份脚本

```shell
#!/bin/bash
#NAME:数据库备份
#DATE:*/*/*
#USER:***

#设置本机数据库登录信息
mysql_user="user"
mysql_password="passwd"
mysql_host="localhost"
mysql_port="3306"
mysql_charset="utf8mb4"
date_time=`date +%Y-%m-%d-%H-%M`

#保存目录中的文件个数
count=10
#备份路径
path=/***/

#备份数据库sql文件并指定目录
mysqldump --all-databases --single-transaction --flush-logs --master-data=2 -h$mysql_host -u$mysql_user -p$mysql_password > $path_$(date +%Y%m%d_%H:%M).sql
[ $? -eq 0 ] && echo "-----------------数据备份成功_$date_time-----------------" || echo "-----------------数据备份失败-----------------"

#找出需要删除的备份
delfile=`ls -l -crt $path/*.sql | awk '{print $9 }' | head -1`
#判断现在的备份数量是否大于阈值
number=`ls -l -crt $path/*.sql | awk '{print $9 }' | wc -l`
if [ $number -gt $count ]
then
    rm $delfile  #删除最早生成的备份,只保留count数量的备份
    #更新删除文件日志
    echo "-----------------已删除过去备份sql $delfile-----------------"
fi
```

增加定时备份
```
crontab -e

*    *    *    *    *
-    -    -    -    -
|    |    |    |    |
|    |    |    |    +----------星期中星期几 (0 - 6) (星期天 为0)
|    |    |    +---------------月份 (1 - 12)
|    |    +--------------------一个月中的第几天 (1 - 31)
|    +-------------------------小时 (0 - 23)
+------------------------------分钟 (0 - 59)

添加定时任务(每天12:50以及23:50执行备份操作)
50 12,23 * * * cd /home/;sh backup.sh >> log.txt
```