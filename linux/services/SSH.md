## 查看服务
```shell
systemctl status sshd
systemctl enable sshd
systemctl start sshd
systemctl reload sshd       # 推荐
systemctl restart sshd      # 不推荐

netstat -tlunp | grep sshd  # 检查ssh运行端口
```

## 生成密钥

- -b bits 指定要创建的秘钥中的位数,默认 2048 位。值越大,密码越复杂
- -C comment 注释,在 id_rsa.pub 中末尾
- -t rsa/dsa等 指定要创建的秘钥类型,默认为 RSA
- -f filename 指定公私钥的名称,会在 $HOME/.ssh 目录下生产私钥 filename 和公钥 filename.pub
- -N password 指定使用秘钥的密码,使得多人使用同一台机器时更安全
```shell
ssh-keygen -t rsa -b 1024 -C "test@example.com"
```

### 发送公钥的两种方式(等价)

前提是需要知道远程服务器的密码...实用性不大...
```
ssh-copy-id -i ~/.ssh/id_rsa.pub user@host
ssh user@host 'mkdir -p .ssh && cat >> .ssh/authorized_keys' < ~/.ssh/id_rsa.pub
```

## 相关文件
```
~/.ssh/id.rsa           # 私钥,很重要需要保密
~/.ssh/id.rsa.pub       # 公钥,放到需要登录的机器上
~/.ssh/known_host       # 曾经连接过的机器会记录hash到这个文件,再次连接时如果hash值不对会有警报
~/.ssh/authorized_keys  # 受信任的机器,公钥存放在这里
/etc/ssh/ssh_config     # SSH客户端配置文件
/etc/ssh/sshd_config    # SSH服务器端的配置文件
```


## sshd配置文件

配置文件地址: `/etc/ssh/sshd_config` sshd(8) 的主配置文件。这个文件的宿主应当是root,权限最大可以是"644"。

```
#Port 22                # 默认ssh端口,生产环境中建议改成五位数的端口 
#AddressFamily any          # 地址家族,any表示同时监听ipv4和ipv6地址
#ListenAddress 0.0.0.0          # 监听本机所有ipv4地址
#ListenAddress ::           # 监听本机所有ipv6地址
HostKey /etc/ssh/ssh_host_rsa_key           # ssh所使用的RSA私钥路径
#HostKey /etc/ssh/ssh_host_dsa_key
HostKey /etc/ssh/ssh_host_ecdsa_key         # ssh所使用的ECDSA私钥路径
HostKey /etc/ssh/ssh_host_ed25519_key           # ssh所使用的ED25519私钥路径

# Ciphers and keying
#RekeyLimit default none

# Logging
#SyslogFacility AUTH
SyslogFacility AUTHPRIV         # 设定在记录来自sshd的消息的时候,是否给出“facility code”
#LogLevel INFO          # 日志记录级别,默认为info 

# Authentication:

#LoginGraceTime 2m          # 限定用户认证时间为2min
#PermitRootLogin yes        # 是否允许root账户ssh登录,生产环境中建议改成no,使用普通账户ssh登录
#StrictModes yes            # 设置ssh在接收登录请求之前是否检查用户根目录和rhosts文件的权限和所有权,建议开启
#MaxAuthTries 6         # 指定每个连接最大允许的认证次数。默认值是 6
#MaxSessions 10         # 最大允许保持多少个连接。默认值是 10 

#PubkeyAuthentication yes       # 是否开启公钥验证

# The default is to check both .ssh/authorized_keys and .ssh/authorized_keys2
# but this is overridden so installations will only check .ssh/authorized_keys
AuthorizedKeysFile      .ssh/authorized_keys        # 公钥验证文件路径

#AuthorizedPrincipalsFile none

#AuthorizedKeysCommand none
#AuthorizedKeysCommandUser nobody

# For this to work you will also need host keys in /etc/ssh/ssh_known_hosts
#HostbasedAuthentication no
# Change to yes if you don't trust ~/.ssh/known_hosts for
# HostbasedAuthentication           # 指定服务器在使用 ~/.shosts ~/.rhosts /etc/hosts.equiv 进行远程主机名匹配时,是否进行反向域名查询
#IgnoreUserKnownHosts no        # 是否在 RhostsRSAAuthentication 或 HostbasedAuthentication 过程中忽略用户的 ~/.ssh/known_hosts 文件
# Don't read the user's ~/.rhosts and ~/.shosts files
#IgnoreRhosts yes           # 是否在 RhostsRSAAuthentication 或 HostbasedAuthentication 过程中忽略 .rhosts 和 .shosts 文件

# To disable tunneled clear text passwords, change to no here!
#PasswordAuthentication yes
#PermitEmptyPasswords no            # 是否允许空密码
PasswordAuthentication yes          # 是否允许密码验证,生产环境中建议改成no,只用密钥登录

# Change to no to disable s/key passwords
#ChallengeResponseAuthentication yes
ChallengeResponseAuthentication no          # 是否允许质疑-应答(challenge-response)认证

# Kerberos options
#KerberosAuthentication no          # 是否使用Kerberos认证
#KerberosOrLocalPasswd yes          # 如果 Kerberos 密码认证失败,那么该密码还将要通过其它的认证机制(比如 /etc/passwd)
#KerberosTicketCleanup yes          # 是否在用户退出登录后自动销毁用户的 ticket
#KerberosGetAFSToken no         # 如果使用了AFS并且该用户有一个 Kerberos 5 TGT,那么开启该指令后,将会在访问用户的家目录前尝试获取一个AFS token
#KerberosUseKuserok yes

# GSSAPI options
GSSAPIAuthentication yes        # 是否允许基于GSSAPI的用户认证
GSSAPICleanupCredentials no         # 是否在用户退出登录后自动销毁用户凭证缓存
#GSSAPIStrictAcceptorCheck yes
#GSSAPIKeyExchange no
#GSSAPIEnablek5users no

# Set this to 'yes' to enable PAM authentication, account processing,
# and session processing. If this is enabled, PAM authentication will
# be allowed through the ChallengeResponseAuthentication and
# PasswordAuthentication.  Depending on your PAM configuration,
# PAM authentication via ChallengeResponseAuthentication may bypass
# the setting of "PermitRootLogin without-password".
# If you just want the PAM account and session checks to run without
# PAM authentication, then enable this but set PasswordAuthentication
# and ChallengeResponseAuthentication to 'no'.
# WARNING: 'UsePAM no' is not supported in Red Hat Enterprise Linux and may cause several
# problems.
UsePAM yes          # 是否通过PAM验证

#AllowAgentForwarding yes
#AllowTcpForwarding yes
#GatewayPorts no            # 是否允许远程主机连接本地的转发端口
X11Forwarding yes           # 是否允许X11转发
#X11DisplayOffset 10        # 指定sshd(8)X11转发的第一个可用的显示区(display)数字。默认值是10
#X11UseLocalhost yes        # 是否应当将X11转发服务器绑定到本地loopback地址
#PermitTTY yes
#PrintMotd yes          # 指定sshd(8)是否在每一次交互式登录时打印 /etc/motd 文件的内容
#PrintLastLog yes       # 指定sshd(8)是否在每一次交互式登录时打印最后一位用户的登录时间
#TCPKeepAlive yes       # 指定系统是否向客户端发送 TCP keepalive 消息
#UseLogin no        # 是否在交互式会话的登录过程中使用 login(1)
#UsePrivilegeSeparation sandbox         # 是否让 sshd(8) 通过创建非特权子进程处理接入请求的方法来进行权限分离
#PermitUserEnvironment no       # 指定是否允许sshd(8)处理~/.ssh/environment以及 ~/.ssh/authorized_keys中的 environment= 选项
#Compression delayed        # 是否对通信数据进行加密,还是延迟到认证成功之后再对通信数据加密
#ClientAliveInterval 0          # sshd(8)长时间没有收到客户端的任何数据,不发送"alive"消息
#ClientAliveCountMax 3          # sshd(8)在未收到任何客户端回应前最多允许发送多个"alive"消息,默认值是 3 
#ShowPatchLevel no
#UseDNS no              # 是否使用dns反向解析
#PidFile /var/run/sshd.pid          # 指定存放SSH守护进程的进程号的路径
#MaxStartups 10:30:100          # 最大允许保持多少个未认证的连接
#PermitTunnel no        # 是否允许tun(4)设备转发
#ChrootDirectory none
#VersionAddendum none

# no default banner path
#Banner none        # 将这个指令指定的文件中的内容在用户进行认证前显示给远程用户,默认什么内容也不显示,"none"表示禁用这个特性

# Accept locale-related environment variables
AcceptEnv LANG LC_CTYPE LC_NUMERIC LC_TIME LC_COLLATE LC_MONETARY LC_MESSAGES
AcceptEnv LC_PAPER LC_NAME LC_ADDRESS LC_TELEPHONE LC_MEASUREMENT
AcceptEnv LC_IDENTIFICATION LC_ALL LANGUAGE
AcceptEnv XMODIFIERS

# override default of no subsystems
Subsystem       sftp    /usr/libexec/openssh/sftp-server        # 配置一个外部子系统sftp及其路径

# Example of overriding settings on a per-user basis
#Match User anoncvs         # 引入一个条件块。块的结尾标志是另一个 Match 指令或者文件结尾    
#       X11Forwarding no
#       AllowTcpForwarding no
#       PermitTTY no
#       ForceCommand cvs server

```







### StrictHostKeyChecking

`StrictHostKeyChecking no` 最不安全的级别,当然也没有那么多烦人的提示了,相对安全的内网测试时建议使用。如果连接server的key在本地不存在,那么就自动添加到文件中(默认是known_hosts),并且给出一个警告。
`StrictHostKeyChecking ask` 默认的级别,就是出现刚才的提示了。如果连接和key不匹配,给出提示,并拒绝登录。
`StrictHostKeyChecking yes` 最安全的级别,如果连接与key不匹配,就拒绝连接,不会提示详细信息。

#### 示例

远程执行命令
```shell
ssh -o "StrictHostKeyChecking no" SSH_USER@IP "echo Hello"
```

批量远程执行命令
```shell
for i in {ip1,ip2};do ssh -o "StrictHostKeyChecking no" SSH_USER@$i "echo Hello"; done
```

远程执行脚本
```shell
ssh -o "StrictHostKeyChecking no" SSH_USER@IP -t "sh /home/omd/hello.sh"
```

### 修改端口
```shell
vim /etc/ssh/sshd_config
Port 2222
```

#### 示例

非默认端口登陆
```shell
# ssh -i private_key username@ipaddress -p port
ssh test@10.100.100.10 -p 2022
```

### 密码登陆

```shell
grep PasswordAuthentication /etc/ssh/sshd_config
PasswordAuthentication yes
```





## 其他功能

### scp

```shell
scp -P22 -r -p /home/omd/h.txt omd@192.168.25.137:/home/omd/
```

### sftp

```shell
sftp -o Port=22 user@IP
```




### 名词解释

https://blog.tinned-software.net/debug-ssh-connection-issue-in-key-exchange/
https://cxybb.com/article/Longyu_wlz/119843133


```
man sshd_config 
```

- Ciphers: `ssh -Q cipher`
- MACs: `ssh -Q mac`
- KexAlgorithms: `ssh -Q kex`
- PubkeyAcceptedKeyTypes: `ssh -Q key`

备注:这些命令输出的项目是 ssh 支持的,并不一定是 ssh 使能的项目。

#### Ciphers

Ciphers 指定 ssh 使能的加密算法。多个加密算法之间使用逗号分隔。当 Ciphers 的值以 + 字符开始时,指定的加密算法将附加到默认集合,不影响默认集合中的其它算法。当 Ciphers 的值以 ‘-’ 字符开始时,指定的加密算法将会从默认集合中移除,不影响默认集合中的其它项目。

The supported ciphers are:
- 3des-cbc
- aes128-cbc
- aes192-cbc
- aes256-cbc
- aes128-ctr
- aes192-ctr
- aes256-ctr
- aes128-gcm@openssh.com
- aes256-gcm@openssh.com
- arcfour
- arcfour128
- arcfour256
- blowfish-cbc
- cast128-cbc
- chacha20-poly1305@openssh.com

The default is:
- chacha20-poly1305@openssh.com,
- aes128-ctr,aes192-ctr,aes256-ctr,
- aes128-gcm@openssh.com,aes256-gcm@openssh.com,
- aes128-cbc,aes192-cbc,aes256-cbc,
- blowfish-cbc,cast128-cbc,3des-cbc

```shell
ssh -Q cipher   # 检查当前支持的加密算法,支持不等于启用

vim /etc/ssh/sshd_config- # 配置cipher

# 在最后一行添加
Ciphers +aes128-cbc,aes192-cbc,aes256-cbc   # 添加加密算法aes128-cbc,aes192-cbc,aes256-cbc
Ciphers -aes128-cbc,aes192-cbc,aes256-cbc   # 移除加密算法aes128-cbc,aes192-cbc,aes256-cbc   
```

#### MACs

MAC (message authentication code) ,MACs 选项指定可用的 MAC(消息认证代码)算法,用于数据完整性保护。

The supported algorithms are:
- curve25519-sha256
- curve25519-sha256@libssh.org
- diffie-hellman-group1-sha1
- diffie-hellman-group14-sha1
- diffie-hellman-group-exchange-sha1
- diffie-hellman-group-exchange-sha256
- ecdh-sha2-nistp256
- ecdh-sha2-nistp384
- ecdh-sha2-nistp521

The default is:
- curve25519-sha256,curve25519-sha256@libssh.org,
- ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,
- diffie-hellman-group-exchange-sha256,
- diffie-hellman-group14-sha1,
- diffie-hellman-group1-sha1

```shell
ssh -Q mac  # 检查当前支持的消息认证代码算法,支持不代表启用了

vim /etc/ssh/sshd_config

# 在最后一行添加
MACs +umac-32   # 添加消息认证代码算法umac-32
MACs -umac-32   # 移除消息认证代码算法umac-32
```

##### MAC不匹配错误
在SSH服务端查看日志
```shell
$ tail /var/log/secure      # CentOS
$ tail -f /var/log/auth.log # In debian based distributions like Ubuntu
```
可能会看到以下信息:

```
sshd[5176]: fatal: no matching mac found: client hmac-sha1,hmac-sha1-96,hmac-md5,hmac-md5-96,hmac-ripemd160,hmac-ripemd160@openssh.com server hmac-sha2-512,hmac-sha2-256
```
可以看到客户端支持6个MAC,服务器端支持2个MAC. 且不共通

服务重载
```
systemctl reload sshd
```

#### KEX 

KEX(Key EXchange) algorithm, KexAlgorithms 选项指定可用的密钥交换算法

The supported algorithms are:
- curve25519-sha256
- curve25519-sha256@libssh.org
- diffie-hellman-group1-sha1
- diffie-hellman-group14-sha1
- diffie-hellman-group-exchange-sha1
- diffie-hellman-group-exchange-sha256
- ecdh-sha2-nistp256
- ecdh-sha2-nistp384
- ecdh-sha2-nistp521

The default is:
- curve25519-sha256,curve25519-sha256@libssh.org,
- ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,
- diffie-hellman-group-exchange-sha256,
- diffie-hellman-group14-sha1,
- diffie-hellman-group1-sha1

```shell
ssh -Q kex  # 检查当前支持的密钥交换算法,支持不等于启用

vim /etc/ssh/sshd_config

# 在最后一行添加
KexAlgorithms +diffie-hellman-group1-sha1   # 添加密钥交换算法diffie-hellman-group1-sha1
KexAlgorithms -diffie-hellman-group1-sha1   # 移除密钥交换算法diffie-hellman-group1-sha1
```

#### PubkeyAcceptedKeyTypes

PubkeyAcceptedKeyTypes 指定公钥认证允许的密钥类型。

The default for this option is:
- ecdsa-sha2-nistp256-cert-v01@openssh.com,
- ecdsa-sha2-nistp384-cert-v01@openssh.com,
- ecdsa-sha2-nistp521-cert-v01@openssh.com,
- ssh-ed25519-cert-v01@openssh.com,
- ssh-rsa-cert-v01@openssh.com,
- ssh-dss-cert-v01@openssh.com,
- ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,
- ssh-ed25519,ssh-rsa,ssh-dss

```shell
ssh -Q key  # 检查当前支持的公钥类型

vim /etc/ssh/sshd_config

# 在最后一行添加
PubkeyAcceptedKeyTypes +keytype   # 添加某种keytype
```

