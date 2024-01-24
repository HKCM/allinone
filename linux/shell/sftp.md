# SFTP

## 连接
```shell
sftp user_name@remote_server_address
sftp -P remote_port user_name@remote_server_address
```

## 操作命令

```shell
get /path/remote_file # 下载文件
put local_file # 上传文件
ls # 查看SFTP目录内容
lls # 查看本地目录内容
![command] # 执行本地 Shell 命令
```

## shell操作
```bash
# SFTP下载文件
curl -u username:password -O ftp://server/file
curl -O ftp://username:password@host/www/focus/enhouse/index.php

# SFTP上传文件
curl -u username:password -T file ftp://server
curl -T xukai.php ftp://username:password@host/www/focus/enhouse/
```


