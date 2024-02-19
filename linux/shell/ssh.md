# ssh

ssh命令是ssh的客户端程序

## 登录过程

SSH 密钥登录分为以下的步骤。

1. 客户端通过ssh-keygen生成自己的公钥和私钥。
2. 手动将客户端的公钥放入远程服务器的指定位置。
3. 客户端向服务器发起 SSH 登录的请求。
4. 服务器收到用户 SSH 登录的请求，发送一些随机数据给用户，要求用户证明自己的身份。
5. 客户端收到服务器发来的数据，使用私钥对数据进行签名，然后再发还给服务器。
6. 服务器收到客户端发来的加密签名后，使用对应的公钥解密，然后跟原始数据比较。如果一致，就允许用户登录。

```
debug1: Reading configuration data /Users/karl.huang/.ssh/config
debug1: /Users/karl.huang/.ssh/config line 24: Applying options for okqa.jump
debug1: Reading configuration data /etc/ssh/ssh_config
```

## 连接命令

```bash
ssh hostname  # 使用当前用户登录
ssh user@hostname # 使用指定用户登录
ssh -p 8821 user@hostname # 使用指定端口用户登录
```

## 远程命令

```bash
# 简单命令
ssh -i key.pem -o StrictHostKeyChecking=no user@remoteNode "cd /home ; ls"
# 互动式 Shell
ssh -i key.pem -o StrictHostKeyChecking=no -t user@remoteNode "top"
```

## 配置参数

```
AddressFamily inet：表示只使用 IPv4 协议。如果设为inet6，表示只使用 IPv6 协议。
BindAddress 192.168.10.235：指定本机的 IP 地址（如果本机有多个 IP 地址）。
CheckHostIP yes：检查 SSH 服务器的 IP 地址是否跟公钥数据库吻合。
Compression yes：是否压缩传输信号。
ConnectionAttempts 10：客户端进行连接时，最大的尝试次数。
ConnectTimeout 60：客户端进行连接时，服务器在指定秒数内没有回复，则中断连接尝试。
DynamicForward 1080：指定动态转发端口。
GlobalKnownHostsFile /users/smith/.ssh/my_global_hosts_file：指定全局的公钥数据库文件的位置。
Host server.example.com：指定连接的域名或 IP 地址，也可以是别名，支持通配符。Host命令后面的所有配置，都是针对该主机的，直到下一个Host命令为止。
HostName myserver.example.com：在Host命令使用别名的情况下，HostName指定域名或 IP 地址。
Ciphers blowfish,3des：指定加密算法。
HostKeyAlgorithms ssh-dss,ssh-rsa：指定密钥算法，优先级从高到低排列。
MACs hmac-sha1,hmac-md5：指定数据校验算法。
IdentityFile keyfile：指定私钥文件。
LogLevel QUIET：指定日志详细程度。如果设为QUIET，将不输出大部分的警告和提示。
NumberOfPasswordPrompts 2：密码登录时，用户输错密码的最大尝试次数。
Port 2035：指定客户端连接的 SSH 服务器端口。
PreferredAuthentications publickey,hostbased,password：指定各种登录方法的优先级。
Protocol 2：支持的 SSH 协议版本，多个版本之间使用逗号分隔。
PasswordAuthentication no：指定是否支持密码登录。不过，这里只是客户端禁止，真正的禁止需要在 SSH 服务器设置。
PubKeyAuthentication yes：是否支持密钥登录。这里只是客户端设置，还需要在 SSH 服务器进行相应设置。
LocalForward 2001 localhost:143：指定本地端口转发。
RemoteForward 2001 server:143：指定远程端口转发。
SendEnv COLOR：SSH 客户端向服务器发送的环境变量名，多个环境变量之间使用空格分隔。环境变量的值从客户端当前环境中拷贝。

StrictHostKeyChecking yes：yes表示严格检查，服务器公钥为未知或发生变化，则拒绝连接。no表示如果服务器公钥未知，则加入客户端公钥数据库，如果公钥发生变化，不改变客户端公钥数据库，输出一条警告，依然允许连接继续进行。ask（默认值）表示询问用户是否继续进行。
TCPKeepAlive yes：客户端是否定期向服务器发送keepalive信息。
ServerAliveCountMax 3：如果没有收到服务器的回应，客户端连续发送多少次keepalive信号，才断开连接。该项默认值为3。
ServerAliveInterval 300：客户端建立连接后，如果在给定秒数内，没有收到服务器发来的消息，客户端向服务器发送keepalive消息。如果不希望客户端发送，这一项设为0，即客户端不会主动断开连接。
User userName：指定远程登录的账户名。
UserKnownHostsFile /users/smith/.ssh/my_local_hosts_file：指定当前用户的known_hosts文件（服务器公钥指纹列表）的位置。
VerifyHostKeyDNS yes：是否通过检查 SSH 服务器的 DNS 记录，确认公钥指纹是否与known_hosts文件保存的一致。
```

## 配置文件

个人配置文件`~/.ssh/config`优先级高于全局配置文件`/etc/ssh/ssh_config`

```
Host server1
    HostName 123.123.123.123
    User neo
    Port 2112

Host *.edu
    User edu
    Port 2122

Host 10.52.* 10.42.*
    Port 2132
    User xxx-admin
    IdentityFile ~/.ssh/internal_key.pem
    ProxyCommand ssh -W %h:%p jump

Host jump
    Hostname 111.111.111.111
    Port 2142
    User jumpuser
    IdentityFile ~/.ssh/jump.key
    Protocol 2
    ForwardAgent yes
    ForwardX11 no
    ServerAliveInterval 300
    ControlMaster auto
    ControlPersist 10
```

这样就可以直接登录了
```bash
ssh server1
```


# 参考文档

https://wangdoc.com/ssh/