# scp

scp [参数] [原路径] [目标路径]

```shell
scp local_file remote_username@remote_ip:remote_file # 将本地文件发送到远端
scp -i key.pem local_file remote_username@remote_ip:remote_file # 将本地文件发送到远端
scp -r local_folder remote_username@remote_ip:remote_folder # 将本地文件夹发送到远端, 递归模式 `-r`
scp -r -i key.pem local_folder remote_username@remote_ip:remote_folder # 将本地文件夹发送到远端, 递归模式 `-r`
scp -v ~/123.txt username@remote_ip:/home/ubuntu/456.txt # 啰嗦模式 `-v`,检查错误的好方法
scp 123.txt 456.txt username@remote_ip:/home/ubuntu/directory/ # 多个文件之间用空格隔开
scp user1@remotehost1:/some/remote/dir/foobar.txt user2@remotehost2:/some/remote/dir/ # 两个远程主机之间复制文件
scp -vrC ~/Downloads root@192.168.1.3:/root/Downloads # 压缩复制文件
```
