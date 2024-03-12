# ssh-keygen

创建密钥 生成密钥 查看公钥

```bash
ssh-keygen -t dsa -b 4096 # 指定密钥算法和密钥长度,默认rsa,1024
ssh-keygen -t rsa -C 'domain@gmail.com' # 添加注释

ssh-keygen -l -f /etc/ssh/ssh_host_ecdsa_key.pub # 查看服务器公钥指纹
cat ~/.ssh/known_hosts # 查看本地记录的指纹
ssh-keygen -R hostname # 删除本地的公钥指纹
```

# 取消ssh密码

备份密钥
```bash
cp ~/.ssh/id_rsa ~/.ssh/id_rsa.backup
```

## 方法一

通过修改ssh密码,使新密码为空的形式,取消ssh密码

```bash
ssh-keygen -f keyname -p
```

## 方法二

使用openssl

```bash
openssl rsa -in ~/.ssh/id_rsa -out ~/.ssh/id_rsa_new
mv ~/.ssh/id_rsa ~/.ssh/id_rsa_old
mv ~/.ssh/id_rsa_new ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa
```


将公钥放入到远程主机对应用户的`~/.ssh/authorized_keys`文件中即可以该用户身份远程登录.

```bash
ssh-copy-id -i id_rsa user@host # 方式一
cat ~/.ssh/id_rsa.pub | ssh user@host "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys" # 方式二
```