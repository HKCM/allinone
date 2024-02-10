# 启动新的docker

```shell
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=mysql -itd mysql5.7utf8:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
# or

docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=mysql -itd mysql5.7utf8:latest
```

# 链接
```
mysql -h127.0.0.1 -uroot -P3306 -p 
```

# 查看字符集
```
mysql> show variables like '%char%'; 
```