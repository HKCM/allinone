## 描述: DNS安装和使用

### DNS简介

DNS:域名系统(英文:Domain Name System)是一个域名系统,主要用于将域名解析为IP,也可以将IP解析为域名.

### DNS实操

实验用两台虚拟机,配置静态IP分别为192.168.88.129和192.168.88.130,以129作为DNS服务器开始
#### 安装
```shell
# bind 是DNS服务,以named运行, bind-chroot是用于改变bind运行目录,模拟一个根目录是bind运行在一个安全环境
$ yum install -y bind bind-chroot
# bind-chroot相当于将bind的根目录移动到了/var/named/chroot/,并且服务名为named-chroot而不是named,二者只能运行一个
$ ls /var/named/chroot/
dev  etc  run  usr  var

```


#### 拷贝文件
* 如果不使用bind-chroot文件所在位置
配置文件:/etc/named.conf
区域数据库文件:/var/named/

* 使用bind-chroot文件所在位置
配置文件:/var/named/chroot/etc/named.conf
区域数据库文件:/var/named/chroot/var/named/

根服务器配置: named.ca
```shell
# 将主配置文拷贝到named-chroot
$ cp -p /etc/named.conf /var/named/chroot/etc/

# 拷贝区域配置文件
$ cp -p /var/named/named.* /var/named/chroot/var/named/
$ cp -pr /var/named/{data,dynamic} /var/named/chroot/var/named/
```

#### 配置文件

##### 主配置文件
```shell
$ vim /var/named/chroot/etc/named.conf
options {
        // IPv4监听端口为53,可以写本机IP或any
        listen-on port 53 { 192.168.88.129; };
        listen-on-v6 port 53 { ::1; };
        // 定义工作目录
        directory       "/var/named";
        // CACHE文件路径,指定服务器在收到rndc dump命令时,转储数据到文件的路径。默认cache_dump.db
        dump-file       "/var/named/data/cache_dump.db";
        // 静态文件路径,指定服务器在收到rndc stats命令时,追加统计数据的文件路径。默认named.stats
        statistics-file "/var/named/data/named_stats.txt";
        // 内存静态文件路径,服务器在退出时,将内存统计写到文件的路径。默认named.mem_stats
        memstatistics-file "/var/named/data/named_mem_stats.txt";
        // 指定服务器在通过rndc recursing命令指定转储当前递归请求到的文件路径。默认named.recursing
        recursing-file  "/var/named/data/named.recursing";
        // 在收到rndc secroots指令后,服务器转储安全根的目的文件的路径名。默认named.secroots
        secroots-file   "/var/named/data/named.secroots";
        // 指定允许哪些主机可以进行普通的DNS查询,可以是关键字:any/localhost/none,也可以是IPV4,IPV6地址 { 192.168.80.0/24;192.168.90.0/24;};
        allow-query     { any; };
        // 复制的区域数据库文件不加密
        masterfile-format text;

        /*
         - If you are building an AUTHORITATIVE DNS server, do NOT enable recursion.
         - If you are building a RECURSIVE (caching) DNS server, you need to enable
           recursion.
         - If your recursive DNS server has a public IP address, you MUST enable access
           control to limit queries to your legitimate users. Failing to do so will
           cause your server to become part of large scale DNS amplification
           attacks. Implementing BCP38 within your network would greatly
           reduce such attack surface
        */
        // 权威DNS不需要开启递归,例如公司严格管理未经解析的域名禁止访问时,则设置为no
        recursion yes;
        //用来设置是否启用DNSSEC支持,DNS安全扩展(DNSSEC)提供了验证DNS数据由效性的系统。
        dnssec-enable yes;
        //指定在DNS查询过程中是否加密,为了加快效率这里可以为no
        dnssec-validation yes;

        /* Path to ISC DLV key */
        bindkeys-file "/etc/named.root.key";

        managed-keys-directory "/var/named/dynamic";

        pid-file "/run/named/named.pid";
        session-keyfile "/run/named/session.key";
        // 只允许向指定IP同步
        allow-transfer { 192.168.88.254; };
};

// logging {
//         channel default_debug {
//                 file "data/named.run";
//                 severity dynamic;
//         };
// };

logging {
    channel query_log {
        // version是指定允许同时存在多少个版本的该文件,例如指定3个版本(version 3),bind会保存query.log、query.log0、query.log1和query.log2。
        file "data/query.log" versions 3 size 20m; 
        // severity是指定记录消息的级别。在bind中主要有以下几个级别 critical,error,warning,notice,info,debug [ level ],dynamic
        severity info;
        print-time yes;
        print-category  yes;
    };
    
    channel warning {
        file "data/warning.log" versions 3 size 20m;
        severity warning;
        print-time yes;
        print-category  yes;
    };
    // default类别匹配所有未明确指定通道的类别,意味着除了queries之外所有其他消息都写入warning.log
    category default {
        warning;
    };
    // 查询请求写入data/query.log
    category queries {
        query_log;
    };
};


zone "." IN {
        type hint;
        file "named.ca";
};

// 正向解析区域配置 域名后面没有.
zone "example.com" IN {
        type master;
        file "example.com.zone";
        
        // 当example有变动时主动通知
        notify yes;
        // 主动通知的IP
        also-notify { 192.168.88.254; };
};

// 反向解析区域配置 域名后面没有.
zone "88.168.192.in-addr.arpa" IN {
        type master;
        file "192.168.88.zone";
};
include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";
```



##### 区域数据库配置文件example.com.zone
```shell
$ vim /var/named/chroot/var/named/example.com.zone
# 注意该文件中所有域名后面都有.
$TTL 3H # 缓存时间
# 解析的域名   类型 授权域    授权域名服务器   管理员邮箱
@       IN SOA  ns1.example.com. rname.invalid. (
                                        0       ; serial 序列号,每次更新该文件系列号都应该变大,以便主从同步
                                        1D      ; refresh 刷新时间,即规定slave域名服务器多长时间查询一个主服务器,以保证Slave服务器的数据是最新的
                                        1H      ; retry 重试时间,即当从服务试图在主服务器上查询更时,而连接失败了,则这个值规定了从服务多长时间后再试 
                                        1W      ; expire 过期时间,从服务器在向主服务更新失败后多长时间后清除对应的记录
                                        3H )    ; minimum
        NS      ns1.example.com.
ns1     A       192.168.88.129
www     A       192.168.88.100
new     CNAME   www

# 修改属组否则named-chroot无法读取
$ chgrp named /var/named/chroot/var/named/example.com.zone

# 检查配置文件
$ named-checkconf /var/named/chroot/etc/named.conf
$ named-checkzone example.com /var/named/chroot/var/named/example.com.zone

# 启动服务并查看
$ systemctl start named-chroot
$ systemctl status named-chroot

# 修改自己的nameserver,使其指向自己的DNS服务器
$ vim /etc/resolv.conf
```

##### 反向解析区域数据库配置文件192.168.88.zone
```shell
$ vim /var/named/chroot/var/named/192.168.88.zone
$TTL 1D
@    IN SOA    ns1.example.com. rname.invalid. (
                    0    ; serial
                    1D    ; refresh
                    1H    ; retry
                    1W    ; expire
                    3H )    ; minimum
    NS    ns1.example.com.
88    PTR    www.zutuanxue.com.
$ chgrp named /var/named/chroot/var/named/192.168.88.zone
$ named-checkzone 192.168.88.in-addr.arpa /var/named/chroot/var/named/192.168.88.zone
$ nslookup 192.168.88.88
88.88.168.192.in-addr.arpa      name = www.zutuanxue.com.
```

### 从DNS服务器

从DNS服务器将部署在192.168.88.130主机
前面步骤基本一致,不过从服务器秩序配置主配置文件而不需要区域数据库文件,因为区域数据库文件会被自动同步到从服务器.
```shell
$ yum install -y bind bind-chroot
$ scp 192.168.88.129:/var/named/chroot/etc/named.conf /var/named/chroot/etc/
$ chgrp named /var/named/chroot/etc/named.conf
$ vim chgrp named /var/named/chroot/etc/named.conf
listen-on port 53 { 192.168.88.254; };
//从服务器不向其他服务器同步所以transfer应该为none
allow-transfer { none; };
zone "example.com" IN {
        type slave;
        file "example.com.zone"; //这样写是放在/var/named目录下
        // file "slaves/baidu.com.zone.slave"; # 如果从服务器同步了主服务器之后,此时同步的数据就会存放在此目录下/var/named/slaves/
        masters { 192.168.88.129; };

        
};

zone "88.168.192.in-addr.arpa" IN {
        type slave;
        file "192.168.88.zone";
        masters { 192.168.88.129; };
};

$ systemctl start named-chroot 
$ systemctl restart named-chroot
$ ll /var/named/chroot/var/named/

```

### 缓存DNS服务器

DNS缓存服务器可以提高DNS访问速度,对局域网上网实现快速解析；适用于低互联网带宽的企业局域网络,减少重复的DNS查询、通过缓存提高速度.缓存服务器不递归查询根域,负担小

当主DNS服务器宕机是缓存DNS服务器也能通过本地缓存继续工作.
```
forwarders {192.168.88.129;};
allow-query-cache { any; };
```

### 分离解析或智能解析技术

主要用于CDN,将不同区域的人的请求解析道不同地址实现最优路径.该配置位于主配置文件.
```shell
$ vim /var/named/chroot/var/named/example.com

acl "china" { 122.71.115.0/24; };
acl "american" { 106.185.25.0/24; };

view "china" {
    match-clients { "china"; };
    zone "example.com" IN {
        type master;
        file "example.com.zone.china";
    }
};

view "american" {
    match-clients { "american"; };
    zone "example.com" IN {
        type master;
        file "example.com.zone.american";
    }
}
```

不同的区域数据库文件,
example.com.zone.china
```shell
$ vim /var/named/chroot/var/named/example.com.zone.china
# 注意该文件中所有域名后面都有.
$TTL 3H # 缓存时间
# 解析的域名   类型 授权域    授权域名服务器   管理员邮箱
@       IN SOA  ns1.example.com. rname.invalid. (
                                        0       ; serial 序列号,每次更新该文件系列号都应该变大,以便主从同步
                                        1D      ; refresh 刷新时间,即规定slave域名服务器多长时间查询一个主服务器,以保证Slave服务器的数据是最新的
                                        1H      ; retry 重试时间,即当从服务试图在主服务器上查询更时,而连接失败了,则这个值规定了从服务多长时间后再试 
                                        1W      ; expire 过期时间,从服务器在向主服务更新失败后多长时间后清除对应的记录
                                        3H )    ; minimum
        NS      ns1.example.com.
ns1     A       192.168.88.129
www     A       122.71.155.888  # 假的示例
new     CNAME   www
```
american
example.com.zone.american
```shell
$ vim /var/named/chroot/var/named/example.com.zone.american
# 注意该文件中所有域名后面都有.
$TTL 3H # 缓存时间
# 解析的域名   类型 授权域    授权域名服务器   管理员邮箱
@       IN SOA  ns1.example.com. rname.invalid. (
                                        0       ; serial 序列号,每次更新该文件系列号都应该变大,以便主从同步
                                        1D      ; refresh 刷新时间,即规定slave域名服务器多长时间查询一个主服务器,以保证Slave服务器的数据是最新的
                                        1H      ; retry 重试时间,即当从服务试图在主服务器上查询更时,而连接失败了,则这个值规定了从服务多长时间后再试 
                                        1W      ; expire 过期时间,从服务器在向主服务更新失败后多长时间后清除对应的记录
                                        3H )    ; minimum
        NS      ns1.example.com.
ns1     A       192.168.88.129
www     A       106.185.25.100  # 假的示例
new     CNAME   www
```


### 额外知识

#### DNS服务器中默认有根域解析为什么还需要forwards?
1. DNS转发功能的实现主要是方便优先共享DNS数据库资源,节省查询时间,而不必每次一级一级的递归去询问根,这样节省带宽流量的时间。
2. 公司私有域是无法在互联网中找到的,只有指定的内网的DNS会解析私有域,所以需要转发到指定DNS

forward又分为两种情况`forward only;`和` forward frist;`

`forward only` 向指定DNS服务器转发,不会查找自己的.子公司DNS服务器指向母公司DNS服务器请求DNS解析。如果母公司DNS服务器挂了 或 子公司到母公司链路故障 或 母公司DNS上不了网.此时子公司DNS服务器可以解析母公司内网服务器域名,但是子公司DNS服务器无法解析Internet域名

`foward first`  向指定DNS服务器转发,找不到则自己查找. 子公司DNS服务器优先指向母公司DNS服务器,如果母公司DNS服务器连不上(比如子公司到母公司的链路故障)会使用子DNS服务器保存的根DNS服务器来解析域名,这时候无法解析母公司内网域名,可以解析Internet域名。


条件转发,对于指定域名到指定的DNS进行解析
zone "xyz.com." {
        type forward;
        forwarders { 192.168.80.111;};
};

#### 域名检测
```shell
# dig [-t type] [-x addr] [name] [@server]
$ dig -t A www.example.com @192.168.88.129

$ nslookup
> server 192.168.88.129
> www.example.com
```

#### rndc配置
rndc(Remote Name Domain Controllerr)是一个远程管理bind的工具,通过这个工具可以在本地或者远程了解当前服务器的运行状况,也可以对服务器进行关闭、重载、刷新缓存、增加删除zone等操作。

使用rndc可以在不停止DNS服务器工作的情况进行数据的更新,使修改后的配置文件生效。在实际情况下,DNS服务器是非常繁忙的,任何短时间的停顿都会给用户的使用带来影响。因此,使用rndc工具可以使DNS服务器更好地为用户提供服务。


1. 使用rndc-confgen在主DNS服务器上生成配置文件的内容,生成的只是文件内容,还需要手动配置
```shell
$  rndc-confgen
# Start of rndc.conf
key "rndc-key" {
        algorithm hmac-md5;
        secret "ZoZN6FJ06O7V8vaxcDuXWQ==";
};

options {
        default-key "rndc-key";
        default-server 127.0.0.1;
        default-port 953;
};
# End of rndc.conf

# Use with the following in named.conf, adjusting the allow list as needed:
# key "rndc-key" {
#       algorithm hmac-md5;
#       secret "ZoZN6FJ06O7V8vaxcDuXWQ==";
# };
#
# controls {
#       inet 127.0.0.1 port 953
#               allow { 127.0.0.1; } keys { "rndc-key"; };
# };
# End of named.conf
```

2. 将生成的第一部分写入`/etc/rndc.conf`,并将其下发到需要远程控制的主机上
```shell
$ vim /etc/rndc.conf
# Start of rndc.conf
key "rndc-key" {
        algorithm hmac-md5;
        secret "ZoZN6FJ06O7V8vaxcDuXWQ==";
};

options {
        default-key "rndc-key";
        # default-server 127.0.0.1; 改为主DNS服务器的地址
        default-server 192.168.88.129;
        default-port 953;
};
# End of rndc.conf
```

3. 将生成配置的第二部分写入主配置named.conf的末尾,我这里使用named-chroot所以是写入/var/named/chroot/etc/named.conf
```shell
$ /var/named/chroot/etc/named.conf
# Use with the following in named.conf, adjusting the allow list as needed:
key "rndc-key" {
      algorithm hmac-md5;
      secret "ZoZN6FJ06O7V8vaxcDuXWQ==";
};
controls {
      inet 0.0.0.0 port 953
            # 在allow中写入允许谁进行远程控制,别忘了分号
              allow { 127.0.0.1;192.168.88.0/24; } keys { "rndc-key"; };
};
# End of named.conf
```
重新启动Bind后,若在/var/log/messages中发现`named[21329]: command channel listening on 0.0.0.0#953` 表示设定成功

4. 按照第二步进行/etc/conf的下发之后,远程机器上也需要有rndc客户端,rndc在bind包中
```shell
# 现在在一台远程机器上,复制主DNS服务器的rndc.conf
$ scp 192.168.88.129:/etc/rndc.conf /etc/rndc.conf
$ yum install -y bind
$ rndc status
WARNING: key file (/etc/rndc.key) exists, but using default configuration file (/etc/rndc.conf)
version: BIND 9.11.4-P2-RedHat-9.11.4-26.P2.el7_9.3 (Extended Support Version) <id:7107deb>
running on localhost.localdomain: Linux x86_64 3.10.0-1160.15.2.el7.x86_64 #1 SMP Wed Feb 3 15:06:38 UTC 2021
boot time: Sat, 27 Feb 2021 08:11:48 GMT
last configured: Sat, 27 Feb 2021 08:11:48 GMT
configuration file: /etc/named.conf (/var/named/chroot/etc/named.conf)
CPUs found: 2
worker threads: 2
UDP listeners per interface: 1
number of zones: 105 (97 automatic)
debug level: 0
xfers running: 0
xfers deferred: 0
soa queries in progress: 0
query logging is ON
recursive clients: 0/900/1000
tcp clients: 2/150
server is up and running
```



其实bind安装完成后rndc默认可用,但是无法从远端操作

#### rndc命令
这个以后用到在更新吧

查看远端DNS状态
rndc status
rndc -s 192.168.88.129 status # 指定server

重载主配置文件和区域解析库文件,可以不用重启named服务情况下更新主配置文件和区域解析库文件。
rndc reload

重域指定区域
rndc reload zone_name:

手动启动区域传送,不管序号是否增加或减少
rndc retransfer zone

手动通知区域
rndc notify zone

如果只改动区域数据库文件可以使用rndc reload进行重载


### 错误
#### named-chroot服务启动后无法解析域名问题
```shell
named-chroot.service - Berkeley Internet Name Domain (DNS)
   Loaded: loaded (/usr/lib/systemd/system/named-chroot.service; disabled; vendor preset: disabled)
   Active: active (running) since Fri 2021-02-26 09:23:33 EST; 56s ago
  Process: 53900 ExecStart=/usr/sbin/named -u named -c ${NAMEDCONF} -t /var/named/chroot $OPTIONS (code=exited, status=0/SUCCESS)
  Process: 53897 ExecStartPre=/bin/bash -c if [ ! "$DISABLE_ZONE_CHECKING" == "yes" ]; then /usr/sbin/named-checkconf -t /var/named/chroot -z "$NAMEDCONF"; else echo "Checking of zone files is disabled"; fi (code=exited, status=0/SUCCESS)
 Main PID: 53902 (named)
   CGroup: /system.slice/named-chroot.service
           └─53902 /usr/sbin/named -u named -c /etc/named.conf -t /var/named/chroot

Feb 26 09:23:33 localhost.localdomain named[53902]: zone 0.in-addr.arpa/IN: loaded serial 0
Feb 26 09:23:33 localhost.localdomain named[53902]: zone 1.0.0.127.in-addr.arpa/IN: loaded serial 0
Feb 26 09:23:33 localhost.localdomain named[53902]: zone 1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa/IN: loaded serial 0
Feb 26 09:23:33 localhost.localdomain named[53902]: zone example.com/IN: loading from master file example.com.zone failed: permission denied
Feb 26 09:23:33 localhost.localdomain named[53902]: zone example.com/IN: not loaded due to errors.
Feb 26 09:23:33 localhost.localdomain named[53902]: zone localhost.localdomain/IN: loaded serial 0
Feb 26 09:23:33 localhost.localdomain named[53902]: zone localhost/IN: loaded serial 0
Feb 26 09:23:33 localhost.localdomain named[53902]: all zones loaded
Feb 26 09:23:33 localhost.localdomain named[53902]: running
Feb 26 09:23:33 localhost.localdomain systemd[1]: Started Berkeley Internet Name Domain (DNS).
```
可以看到`example.com.zone failed: permission denied`无法加载区域配置文件,原因是区域配置文件的属组不正确
```shell
$ ll /var/named/chroot/var/named/
total 24
-rw-r-----. 1 root root   208 Feb 26 09:19 192.168.88.zone
drwxrwxrwx. 2 root root    94 Feb 26 09:00 data
drwxrwxrwx. 2 root root    25 Feb 26 09:23 dynamic
-rw-r-----. 1 root root   215 Feb 26 09:08 example.com.zone

$ chgrp named example.com.zone
$ systemctl restart named-chroot # 修改后重启服务一切正常
```

#### 从服务器不更新问题
1. 变更区域数据库文件记得增加serial号
2. 可以在主服务器主配置文件中添加主动通知
    // 当example有变动时主动通知
    notify yes;
    // 主动通知的IP
    also-notify { 192.168.88.254; };
3. 降低区域数据库文件中的refresh时间,但是会加大服务器负担,因为刷新变得频繁




####  isc_stdio_open 'data/named.run' failed: permission denied

提前将文件创建好并赋予777权限

```shell
$ touch /var/named/chroot/var/named/data/named.run
$ chmod -R 777 /var/named/chroot/var/named/data
```

### 参考

https://www.cnblogs.com/fjping0606/p/4428736.html
https://blog.csdn.net/HzSunshine/article/details/54289553

