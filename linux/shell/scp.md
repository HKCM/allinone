# scp

scp是全量复制

-C 压缩传输
-P 指定端口
-p 传输后保留文件原始属性
-r 递归复制整个目录

## 本地传输到远端

```shell
scp local_file remote_username@remote_ip:remote_file # 将本地文件发送到远端
scp -i key.pem local_file remote_username@remote_ip:remote_file # 将本地文件发送到远端
scp -r local_folder remote_username@remote_ip:remote_folder # 将本地文件夹发送到远端, 递归模式 `-r`
scp -r -i key.pem local_folder remote_username@remote_ip:remote_folder # 将本地文件夹发送到远端, 递归模式 `-r`
scp 123.txt 456.txt username@remote_ip:/home/ubuntu/directory/ # 多个文件之间用空格隔开
scp user1@remotehost1:/some/remote/dir/foobar.txt user2@remotehost2:/some/remote/dir/ # 两个远程主机之间复制文件
scp -vrC ~/Downloads root@192.168.1.3:/root/Downloads # 压缩复制文件
```

## 远端传输到本地

```bash
# 语法
scp user@host:directory/SourceFile TargetFile

# 示例
scp remote_username@remote_ip:/remote/file.txt /local/directory
scp -r remote_username@remote_ip:/path_to_remote_directory/* local-machine/path_to_the_directory/
scp -r user@host:directory/SourceFolder TargetFolder
```

## 远端传输到远端

```bash
scp user1@host1.com:/files/file.txt user2@host2.com:/files
```


## 其他

```bash
scp -c blowfish some_file username@ip:/tmp/ # 传输加密并指定加密算法
scp -C some_file username@ip:/tmp/ # 传输压缩
scp -v ~/123.txt username@remote_ip:/home/ubuntu/456.txt # 啰嗦模式 `-v`,检查错误
scp -l 80 some_file username@ip:/tmp/ # 限定速率 Kb/s
scp -P 2222 user@host:directory/SourceFile TargetFile # 指定端口
scp -p user@host:directory/SourceFile TargetFile # 保留原始文件消息
```