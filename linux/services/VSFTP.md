### 描述: CentOS7 下的VSFTP的安装和使用

#### 安装
```shell
# 安装vsftp 和ftp命令客户端
$ yum install -y vsftpd ftp
$ systemctl enable vsftpd
$ systemctl start vsftpd
$ systemctl status vsftpd
```

相关文件
* 主配置文件: /etc/vsftpd/vsftpd.cof
* 下载目录: /var/ftp/
* FTP日志: /var/log/xferlog
* FTP日志: /var/log/vsftpd.log  需要额外配置

可以通过浏览器访问ftp://192.168.88.129进行验证

#### 配置文件

以下只列出一部分关键关键配置,具体设置与使用情况有关
```shell
# 编辑前先备份
$ cp /etc/vsftpd/vsftpd.conf /etc/vsftpd/vsftpd.backup.conf

$ vim /etc/vsftp/vsftp.conf
# 匿名用户访问,YES是允许,NO是拒绝,
anonymous_enable=NO
# 开启匿名用户上传功能,默认是拒绝的
anon_upload_enable=YES
# 开启匿名用户创建文件或文件夹权限
anon_mkdir_write_enable=YES
# 允许更改匿名用户上传文件的所有者
chown_uploads=YES
# 所有者为whoever,更改所有者会保护文件,防止匿名用户进行文件覆盖
chown_username=whoever
# 设置项用于设置匿名用户的最大传输速率,单位为B/s,值为0表示不限制。例如ano_max_rate=200000,表示FTP服务器的匿名用户最大传输速率设置为200KB/s
anon_max_rate

#设定本地用户可以访问,主要是虚拟宿主用户,如果设为NO那么所欲虚拟用户将无法访问。
local_enable=YES


# 开启上传和下载日志记录功能
xferlog_enable=YES
# 日志文件路径
xferlog_file=/var/log/xferlog

# vsftpd日志,会记录上传,下载,以及连接和登陆日志
dual_log_enable=YES 
# 日志文件路径
vsftpd_log_file=/var/log/vsftpd.log  
#Sat Feb 27 21:57:20 2021 [pid 22116] CONNECT: Client "::ffff:192.168.88.254"
#Sat Feb 27 21:57:32 2021 [pid 22115] [test] OK LOGIN: Client "::ffff:192.168.88.254"

#禁止用户登陆ftp后使用ls -R 命令。该命令会对服务器性能造成巨大开销,如果该项运行当多个用户使用该命令会对服务器造成威胁。
ls_recurse_enable=NO

# 如果以standalone模式起动,那么只有$Number个用户可以连接,其他的用户将得到错误信息,默认是0不限止上限报错421 There are too many connected users, please try later.
max_clients=Number 
max_per_ip=Number
#设置项用于设置本地用户的最大传输速率,单位为B/s,值为0时表示不限制。例如local_max_rate=500000表示FTP服务器的本地用户最大传输速率设置为500KB/s. 
local_max_rate=500000


# 会话超时时间默认秒为单位
idle_session_timeout=600
# 数据传输超时时间,默认120/秒
data_connection_timeout=120

#设定支持异步传输的功能,通常不安全不建议启动
async_abor_enable=NO

# 是否开启对本地用户chroot的限制,YES为默认所有用户都不能切出家目录,NO代表默认用户都可以切出家目录
# 设置方法类似于:YES拒绝所有,允许个别    NO  允许所有拒绝个别
chroot_local_user=YES
#开启特例列表
chroot_list_enable=YES
# (default follows)
# 如果chroot_local_user的值是YES则该文件中的用户是可以切出家目录,如果是NO该文件中的用户则不能切出家目录,需要手动创建该文件
chroot_list_file=/etc/vsftpd/chroot_list
# 允许家目录具有可写权限
allow_writeable_chroot=YES

# 上传的权限是022,使用的是umask权限。对应的目录是755,文件是644
local_umask=022
 
#用户列表功能
userlist_enable=YES
# 拒绝列表用户 YES表示列表用户不允许登陆,NO表示只允许列表用户登陆
#userlist_deny=NO
# 只允许列表中的用户
#userlist_file=/etc/vsftpd/user_list
# 如果列表中没有指定用户登陆会失败
# $ ftp 192.168.88.129
# Connected to 192.168.88.129 (192.168.88.129).
# 220 (vsFTPd 3.0.2)
# Name (192.168.88.129:root): test
# 530 Permission denied.
# Login failed.

tcp_wrappers=YES
# 该指令的配置文件在/etc/host.deny 以及/etc/host.allow

# 修改默认的FTP服务器的端口号,应尽量大于4000,修改后访问
listen_port=8888
```

#### 虚拟用户
在ftp中不论是匿名用户还是实名用户都是系统中真实存在的用户,或多或少都会有一些安全方面的风险,为了避免这个风险,建议启用虚拟用户

`匿名用户其实是ftp用户,可以在/etc/passwd中看到ftp用户的存在`

1. 要求能够匿名下载
2. 各部门拥有自己的文件夹,it,finance,facility
3. 只有主管能上传,普通员工只能下载
4. 禁止查看家目录以外的数据

* 创建虚拟用户
```shell
$ useradd -s /sbing/nologin -d /var/tmp/vuser_ftp vftpuser
```

* 修改家目录权限,vsftp不允许有写权限
```shell
$ chmod 500 /var/tmp/vuser_ftp
```

* 创建各部门文件夹,并修改权限
```shell
$ mkdir /var/tmp/vuser_ftp/{it,finance,facility}
$ chmod 700 /var/tmp/vuser_ftp/*
$ chown vftpuser. /var/tmp/vuser_ftp/* -R
```

在配置文件`/etc/vsftpd/vsftpd.conf`中添加虚拟用户配置
```shell
# 禁止切换出家目录
chroot_local_user=YES
# 启用虚拟用户功能
guest_enable=YES
 
# 指定虚拟的宿主用户,最好使用nologin的
guest_username=vftpuser
 
# 设定虚拟用户的权限符合他们的宿主用户
virtual_use_local_privs=NO
 
# 设定虚拟用户个人vsftp的配置文件存放路劲。这个被指定的目录里,将被存放每个虚拟用户个性的配置文件,注意的地方是:配置文件名必须和虚拟用户名相同。
user_config_dir=/etc/vsftpd/vconf.d/

# 家目录允许可写
allow_writeable_chroot=YES
 
#禁止反向域名解析,若是没有添加这个参数可能会出现用户登陆较慢,或则客户链接不上ftp的现象
reverse_lookup_enable=NO

```

建立虚拟用户账号密码,并创建账号密码数据库
```shell
$ cat /etc/vsftpd/vuser
it_admin
123
finance_admin
123
facility_admin
123
it_user
321
finace_user
321
facility_user
321

$ db_load -T -t hash -f /etc/vsftpd/vuser /etc/vsftpd/vuser.db
$ chmod 600 /etc/vsftpd/vuser.db
```

添加pam认证,要添加在首行
```shell
$ cat /etc/pam.d/vsftpd

auth    sufficient      /lib64/security/pam_userdb.so db=/etc/vsftpd/vuser
account sufficient      /lib64/security/pam_userdb.so db=/etc/vsftpd/vuser
#%PAM-1.0
session    optional     pam_keyinit.so    force revoke
auth       required     pam_listfile.so item=user sense=deny file=/etc/vsftpd/ftpusers onerr=succeed
auth       required     pam_shells.so
auth       include      password-auth
account    include      password-auth
session    required     pam_loginuid.so
session    include      password-auth
```

禁止切出家目录的chroot_list
```shell
$ cat /etc/vsftpd/chroot_list
it_admin
finance_admin
facility_admin
it_user
finace_user
facility_user
```

admin权限文件模板
```shell
$ mkdir /etc/vsftpd/vconf.d/
$ vim /etc/vsftpd/vconf.d/admin_temp
#指定家目录
local_root=/var/tmp/vuser_ftp/XXX
#umask权限
anon_umask=077
#开放下载权限
anon_world_readable_only=NO
#开放上传权限
anon_upload_enable=YES
#创建目录权限
anon_mkdir_write_enable=YES
#开放删除和重命名目录权限
anon_other_write_enable=YES 
```

一般用户权限模板
```shell
$ vim /etc/vsftpd/vconf.d/user_temp
#指定家目录
local_root=/var/tmp/vuser_ftp/XXX
#开放下载权限
anon_world_readable_only=NO
```

为每个部门和用户创建配置文件以it部门为例
```shell
$ vim /etc/vsftpd/vconf.d/it_admin
#指定家目录
local_root=/var/tmp/vuser_ftp/it var/tmp/vuser_ftp/
#umask权限
anon_umask=077
#开放下载权限
anon_world_readable_only=NO
#开放上传权限
anon_upload_enable=YES
#创建目录权限
anon_mkdir_write_enable=YES
#开放删除和重命名目录权限
anon_other_write_enable=YES
```
it_user
```shell
$ vim /etc/vsftpd/vconf.d/it_user
#指定家目录
local_root=/var/tmp/vuser_ftp/it
#开放下载权限
anon_world_readable_only=NO
```


### 扩展 
#### ftp命令

| 命令| 功能| 命令| 功能|
|-|-|-|-|
|ls|显示服务器上的目录  | ls [remote-dir][local-file] |显示远程目录remote-dir,并存入本地文件local-file |
| get remote-file [local-file]|从服务器下载指定文件到客户端  |mget remote-files  | 下载多个远程文件(mget命令允许用通配符下载多个文件)|
|put local-file [remote-file]|从客户端上传指定文件到服务器 |mput local-file|将多个文件上传至远程主机(mput命令允许用通配符上传多个文件) |
|delete remote-file|删除远程主机文件 |mdelete [remote-file] |删除远程主机文件允许用通配符 |
|mkdir dir-name |在远程主机中创建目录 |rmdir dir-name |删除远程主机目录 |
|rename [from][to]|更改远程主机的文件名 |||
|cd|改变服务器的工作目录|cdup|进入远程主机目录的父目录|
|bye|退出FTP命令状态 |quit |同bye,退出ftp会话|
|open|连接FTP服务器|close|中断与远程服务器的ftp会话(与open对应)|
|size file-name |显示远程主机文件大小|newer file-name |如果远程主机中file-name的修改时间比本地硬盘同名文件的时间晚则更新本地文件|


* 使用`!`可以在客户端(本机)上运行linux命令

例如不知道登陆ftp后自己的本机当前目录可以使用
```shell
$ ftp> !pwd
/home/test
# 在本机tmp目录创建456空文件
$ ftp> !touch /tmp/456
$ ftp> !ls /tmp/
456
```

* lcd
切换本机目录
```shell
$ ftp> lcd /tmp
Local directory now /tmp
```

#### 状态码
|状态码|意义|
|-|-|
|110| 重新启动标记应答。 |
|120| 服务在多久时间内 ready 。 |
|125| 数据链路端口开启,准备传送。 |
|150| 文件状态正常,开启数据连接端口。 |
|200| 命令执行成功。 |
|202| 命令执行失败。 |
|211| 系统状态或是系统求助响应。 |
|212| 目录的状态。 |
|213| 文件的状态。 |
|214| 求助的讯息。 |
|215| 名称系统类型。 |
|220| 新的联机服务 ready 。 |
|221| 服务的控制连接端口关闭,可以注销。 |
|225| 数据连结开启,但无传输动作。 |
|226| 关闭数据连接端口,请求的文件操作成功。 |
|227| 进入 passive mode 。 |
|230| 使用者登入。 |
|250| 请求的文件操作完成。 |
|257| 显示目前的路径名称。 |
|331| 用户名称正确,需要密码。 |
|332| 登入时需要账号信息。 |
|350| 请求的操作需要进一部的命令。 |
|421| 无法提供服务,关闭控制连结。 |
|425| 无法开启数据链路。 |
|426| 关闭联机,终止传输。 |
|450| 请求的操作未执行。 |
|451| 命令终止 : 有本地的错误。 |
|452| 未执行命令 : 磁盘空间不足。 |
|500| 格式错误,无法识别命令。 |
|501| 参数语法错误。 |
|502| 命令执行失败。 |
|503| 命令顺序错误。 |
|504| 命令所接的参数不正确。 |
|530| 未登入。   532 储存文件需要账户登入。 550 未执行请求的操作。 551 请求的命令终止,类型未知。 |
|552| 请求的文件终止,储存位溢出。    553 未执行请求的的命令,名称不正确。|

#### 日志详解
/var/log/xferlog

Sat Feb 27 22:04:20 2021 1 ::ffff:192.168.88.254 0 /123.txt b _ o r test ftp 0 * c
|示例数据 |参数名称 |说明 |
|-|-|-|
|Sat Feb 27 22:04:20 2021 |当前时间  |当前服务器本地时间,格式为: DDD MMM dd hh:mm:ss YYY  |
|1 |传输时间  |传送文件所用时间,单位为秒 |
|::ffff:192.168.88.254|远程主机名称/IP|远程主机名称/IP |
|0|文件大小|传送文件的大小,单位为byte |
|/123.txt|文件名|传输文件名,包括路径 |
|b|传输类型|传输方式的类型,包括两种: a以ASCII传输 b以二进制文件传输 |
|_|特殊处理标志|特殊处理的标志位:_ 不做任何特殊处理,C 文件是压缩格式,U 文件是非压缩格式,T 文件是tar格式 |
|o|传输方向|o 从FTP服务器向客户端传输 i 从客户端向FTP服务器传输 |
|r|访问模式|a 匿名用户 g 来宾用户 r 真实用户,即系统中的用户 |
|test|用户名|用户名|
|ftp|服务名|一般为ftp|
|0|认证方式|0 无 1 RFC931认证 |
|认证用户ID|认证用户的id,* 则表示无法获得该id|
|c|状态|c 表示传输已完成 i 表示传输示完成 |

#### Port和Pasv

* Port模式(主动模式)--> 默认

Ftp客户端首先和Ftp server的tcp 21端口建立连接,通过这个通道发送命令,客户端要接受数据的时候在这个通道上发送Port命令,Port命令包含了客户端用什么端口(一个大于1024的端口)接受数据,在传送数据的时候,服务器端通过自己的TCP 20端口发送数据。这个时候数据连接由server向client建立一个连接。

**Port交互流程:**

client端:client链接server的21端口,并发送用户名密码和一个随机在1024上的端口及port命令给server,表明采用主动模式,并开放那个随机的端口。

server端:server收到client发来的Port主动模式命令与端口后,会通过自己的20端口与client那个随机的端口连接后,进行数据传输。


* Pasv模式(被动方式)

建立控制通道和Port模式类似,当客户端通过这个通道发送Pasv命令的时候,Ftp server打开了一个位于1024和5000之间的随机端口并且通知客户端在这个端口上进行传输数据请求,然后Ftp server将通过这个端口进行数据传输。这个时候数据连接由client向server建立连接。

**Pasv交互流程**

Clietn:client连接server的21号端口,发送用户名密码及pasv命令给server,表明采用被动模式。

server:server收到client发来的pasv被动模式命令之后,把随机开放在1024上的端口告诉client,client再用自己的20 端口与server的那个随机端口进行连接后进行数据传输。

 

如果从C/S模型这个角度来说,PORT对于服务器来说是OUTBOUND,而PASV模式对于服务器是INBOUND,这一点请特别注意,尤其是在使用防火墙的企业里,这一点非常关键,如果设置错了,那么客户将无法连接。

### 错误 
#### 500 OOPS: could not read chroot() list file:/etc/vsftpd/chroot_list

如果启用了禁止切换根目录则需要先创建`/etc/vsftpd/chroot_list`文件

#### 500 OOPS: vsftpd: refusing to run with writable root inside chroot()
测试使用的是实体用户test,而vsftp不允许用户家目录具有写权限
ll /home/
total 0
drwx------ 2 test  test  62 Feb 27 20:56 test
1. 修改家目录权限,去掉写权限
```shell
$ chmod 500 /home/test
```
2. 修改vsftp配置文件添加允许可写指令
```shell
$ vim /etc/vsftp/vsftp.conf
allow_writeable_chroot=YES
```

#### 500 OOPS: bad bool value in config file for: anon_other_write_enable
配置文件中`anon_other_write_enable`指令等于YES后面多了一个空格

#### 500 OOPS: cannot change directory:/var/tmp/vuser_ftp/it
这个目录不存在,创建即可

### 参考 
https://www.cnblogs.com/Confession/p/6813227.html
https://www.cnblogs.com/helonghl/articles/5533857.html

