

## 描述: DHCP(Dynamic Host Configuration Protocol)安装和使用


### 简介

DHCP(Dynamic Host Configuration Protocol,动态主机配置协议)通常被应用在大型的局域网络环境中,主要作用是集中的管理、分配IP地址,使网络环境中的主机动态的获得IP地址、Gateway地址、DNS服务器地址等信息,并能够提升地址的使用率。

### DHCP安装

```shell
$ yum install -y dhcp
```

### DHCP 实践

1. 主配置文件 
  ```shell
  $ ls /etc/dhcp/dhcpd.conf
  ```

2. 复制example文件
    在`/usr/share/doc/dhcp*/`目录下有DHCP的示例文件
  ```shell
  $ mv /etc/dhcp/dhcpd.conf /etc/dhcp/dhcpd.conf.bak
  $ cp /usr/share/doc/dhcp*/dhcpd.conf.example /etc/dhcp/dhcpd.conf
  ```

3. 配置文件详解
  ```shell
  $ # cat /etc/dhcp/dhcpd.conf
  # #号代表注释
  # DHCP服务配置文件分为全局配置和作用域配置,很好区分:subnet的就是作用域 不在subnet里面的就是全局设置。
  # dhcpd.conf
  #
  # Sample configuration file for ISC dhcpd
  #

  #DNS全局选项,指定DNS服务器的地址,可以是IP,也可以是域名。
  # option definitions common to all supported networks...
  # 默认搜索域 该配置项将体现在客户机的“/etc/resolv.conf”配置文件中(如 “search benet.com”)
  option domain-name "example.org";

  #具体的DNS服务器
  option domain-name-servers ns1.example.org, ns2.example.org;

  #默认租约为600s
  default-lease-time 600;
  #最大租约为7200s,客户机会在default-lease-time快到期时向服务器续租,如果此时客户机死机/重启,而默认租约时间到期了,服务器并不会立即回收IP地址,而是等到最大租约时间到期,客户机还没有过来续约,才会回收IP地址。
  max-lease-time 7200;

  #动态DNS更新方式(none:不支持；interim:互动更新模式；ad-hoc:特殊更新模式)
  # Use this to enble / disable dynamic dns updates globally.
  #ddns-update-style none;

  #如果该DHCP服务器是本地官方DHCP就将此选项打开,避免其他DHCP服务器的干扰。
  #当一个客户端试图获得一个不是该DHCP服务器分配的IP信息,DHCP将发送一个拒绝消息,而不会等待请求超时。当请求被拒绝,客户端会重新向当前DHCP发送IP请求获得新地址。保证IP是自己发出去的
  #
  # If this DHCP server is the official DHCP server for the local
  # network, the authoritative directive should be uncommented.
  #authoritative;

  # Use this to send dhcp log messages to a different log file (you also
  # have to hack syslog.conf to complete the redirection).
  # 日志级别
  log-facility local7;

  # No service will be given on this subnet, but declaring it helps the 
  # DHCP server to understand the network topology.

  #作用域相关设置指令
  #subnet 定义一个作用域
  #netmask 定义作用域的掩码
  #range 允许发放的IP范围
  #option routers 指定网关地址
  #option domain-name-servers 指定DNS服务器地址
  #option broadcast-address 广播地址
  #
  #
  #案例:定义一个作用域 网段为10.152.187.0 掩码为255.255.255.0
  #此作用域不提供任何服务
  subnet 10.152.187.0 netmask 255.255.255.0 {
  }

  # This is a very basic subnet declaration.

  #案例:定义一个基本的作用域
  #网段10.254.239.0 掩码255.255.255.224
  #分发范围10.254.239.10-20
  #网关为rtr-239-0-1.example.org, rtr-239-0-2.example.org
  subnet 10.254.239.0 netmask 255.255.255.224 {
    range 10.254.239.10 10.254.239.20;
    option routers rtr-239-0-1.example.org, rtr-239-0-2.example.org;
  }

  # This declaration allows BOOTP clients to get dynamic addresses,
  # which we don't really recommend.
  #案例:允许采用bootp协议的客户端动态获得地址
  #bootp DHCP的前身
  #BOOTP用于无盘工作站的局域网中,可以让无盘工作站从一个中心服务器上获得IP地址。通过BOOTP协议可以为局域网中的无盘工作站分配动态IP地址,
  #这样就不需要管理员去为每个用户去设置静态IP地址。
  subnet 10.254.239.32 netmask 255.255.255.224 {
    range dynamic-bootp 10.254.239.40 10.254.239.60;
    option broadcast-address 10.254.239.31;
    option routers rtr-239-32-1.example.org;
  }


  #案例:一个简单的作用域案例
  # A slightly different configuration for an internal subnet.
  subnet 10.5.5.0 netmask 255.255.255.224 {
    range 10.5.5.26 10.5.5.30;
    option domain-name-servers ns1.internal.example.org;
    option domain-name "internal.example.org";
    option routers 10.5.5.1;
    option broadcast-address 10.5.5.31;
    default-lease-time 600;
    max-lease-time 7200;
  }

  # Hosts which require special configuration options can be listed in
  # host statements.   If no address is specified, the address will be
  # allocated dynamically (if possible), but the host-specific information
  # will still come from the host declaration.
  #
  #保留地址:可以将指定的IP分发给指定的机器,根据网卡的MAC地址来做触发
  #host: 启用保留。
  #hardware:指定客户端的mac地址
  #filename:指定文件名
  #server-name:指定下一跳服务器地址
  #fixed-address: 指定保留IP地址
  #
  #
  #案例:这个案例中分发给客户端的不是IP地址信息,而是告诉客户端去找toccata.fugue.com服务器,并且下载vmunix.passacaglia文件
  host passacaglia {
    hardware ethernet 0:0:c0:5d:bd:95;
    filename "vmunix.passacaglia";
    server-name "toccata.fugue.com";
  }

  # Fixed IP addresses can also be specified for hosts.   These addresses
  # should not also be listed as being available for dynamic assignment.
  # Hosts for which fixed IP addresses have been specified can boot using
  # BOOTP or DHCP.   Hosts for which no fixed address is specified can only
  # be booted with DHCP, unless there is an address range on the subnet
  # to which a BOOTP client is connected which has the dynamic-bootp flag
  # set.
  # 案例:保留地址,将指定IP(fantasia.fugue.com对应的IP)分给指定客户端网卡(MAC:08:00:07:26:c0:a5)
  host fantasia {
    hardware ethernet 08:00:07:26:c0:a5;
    fixed-address fantasia.fugue.com;
  }

  #超级作用域
  #超级作用域是DHCP服务中的一种管理功能,使用超级作用域,可以将多个作用域组合为单个管理实体。
  # You can declare a class of clients and then do address allocation
  # based on that.   The example below shows a case where all clients
  # in a certain class get addresses on the 10.17.224/24 subnet, and all
  # other clients get addresses on the 10.0.29/24 subnet.


  #在局域网中,可以配置策略根据各个机器的具体信息分配IP地址和其他的网络参数,客户机的具体信息:客户机能够给dhcp服务提供的信息由有两个,
  #第一个就是网卡的dhcp-client-identifier(mac地址),
  #第二个就是设备的vendor-class-identifier。
  #管理员可以根据这两个信息给不同的机器分组。

  #案例:
  #按client某种类型分组DHCP,而不是按物理接口网段
  #例子: SUNW 分配地址段10.17.224.0/24
  # 非SUNW的主机,分配地址段10.0.29.0/24
  #定义一个dhcp类:foo
  #request广播中vendor-class-identifier字段对应的值前四个字节如果是"SUNW",则视合法客户端.
  class "foo" {
    match if substring (option vendor-class-identifier, 0, 4) = "SUNW";
  }

  #定义一个超级作用域: 224-29
  shared-network 224-29 {
  #定义第一个作用域
    subnet 10.17.224.0 netmask 255.255.255.0 {
      option routers rtr-224.example.org;
    }
  #定义第二个作用域
    subnet 10.0.29.0 netmask 255.255.255.0 {
      option routers rtr-29.example.org;
    }

  #关连池,如果客户端匹配foo类,将获得该池地址
    pool {
      allow members of "foo";
      range 10.17.224.10 10.17.224.250;
    }
  #关连池,如果客户端配置foo类,则拒绝获得该段地址
    pool {
      deny members of "foo";
      range 10.0.29.10 10.0.29.230;
    }
  }
  ```
  在DHCP服务器192.168.88.129上配置和启动dhcp,dhcp虚拟机最好配置静态IP
  DHCP在启动的时候检查配置文件,需要一个和服务器同网段的作用域。
  ```shell
  $ cat /etc/dhcp/dhcpd.conf
  # 例如将示例文件中的10.152.187.0 改为和当前服务器(192.168.88.129)相同的网段192.168.88.0 
  option domain-name-servers 202.106.0.20, 114.114.114.114;
  default-lease-time 7200; #定义默认租约时间
  max-lease-time 10800; #定义最大租约时间
  authoritative; #拒绝不正确的IP地址的要求
  log-facility local7; #定义日志

  subnet 192.168.88.0 netmask 255.255.255.0 {
    range 192.168.88.153 192.168.88.252;
    option routers 192.168.88.254;
    option broadcast-address 192.168.88.255;
    default-lease-time 7200;
    max-lease-time 10800;
  }

  $ systemctl start dhcpd
  $ netstat -luntp | grep 67
  udp        0      0 0.0.0.0:67              0.0.0.0:*                           59366/dhcpd
  
  ```

4. 测试dhcp
    在测试机(双网卡192.168.88.131,192.168.88.133)上查看
  ```shell
  $ [root@localhost ~]# ip addr
  1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
      link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
      inet 127.0.0.1/8 scope host lo
        valid_lft forever preferred_lft forever
      inet6 ::1/128 scope host
        valid_lft forever preferred_lft forever
  2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
      link/ether 00:0c:29:86:50:46 brd ff:ff:ff:ff:ff:ff
      inet 192.168.88.131/24 brd 192.168.88.255 scope global noprefixroute dynamic ens33
        valid_lft 1139sec preferred_lft 1139sec
      inet6 fe80::7137:324e:eb37:4d39/64 scope link noprefixroute
        valid_lft forever preferred_lft forever
  3: ens37: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
      link/ether 00:0c:29:86:50:50 brd ff:ff:ff:ff:ff:ff
      inet 192.168.88.132/24 brd 192.168.88.255 scope global noprefixroute dynamic ens37
        valid_lft 1529sec preferred_lft 1529sec
      inet6 fe80::2773:e58:8c95:8aa/64 scope link noprefixroute
        valid_lft forever preferred_lft forever
  $ systemctl restart networkd

  # 释放ens33,如果无法释放则先获取再释放 dhclient -d ens33
  $ dhclient -r ens33
  # 释放完毕ens37上已经没有IP了
  $ ip addr show ens33
  2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:86:50:46 brd ff:ff:ff:ff:ff:ff
    inet6 fe80::7137:324e:eb37:4d39/64 scope link noprefixroute
       valid_lft forever preferred_lft forever

  $ dhclient -d ens33
  Internet Systems Consortium DHCP Client 4.2.5
  Copyright 2004-2013 Internet Systems Consortium.
  All rights reserved.
  For info, please visit https://www.isc.org/software/dhcp/

  Listening on LPF/ens33/00:0c:29:86:50:46
  Sending on   LPF/ens33/00:0c:29:86:50:46
  Sending on   Socket/fallback
  DHCPDISCOVER on ens33 to 255.255.255.255 port 67 interval 6 (xid=0x46d18ca1)
  DHCPREQUEST on ens33 to 255.255.255.255 port 67 (xid=0x46d18ca1)
  DHCPOFFER from 192.168.88.129
  DHCPACK from 192.168.88.129 (xid=0x46d18ca1)
  bound to 192.168.88.157 -- renewal in 3107 seconds.
  ^C
  $ ip addr show ens33
  2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
      link/ether 00:0c:29:86:50:46 brd ff:ff:ff:ff:ff:ff
      inet 192.168.88.157/24 brd 192.168.88.255 scope global dynamic ens33
        valid_lft 7165sec preferred_lft 7165sec
      inet6 fe80::7137:324e:eb37:4d39/64 scope link noprefixroute
        valid_lft forever preferred_lft forever
  ```

5. 查看dhcp服务器日志
  ```shell
  $ tail -f /var/log/message
  Feb 25 10:20:36 localhost dhcpd: DHCPDISCOVER from 00:0c:29:86:50:50 via ens32
  Feb 25 10:20:37 localhost dhcpd: DHCPOFFER on 192.168.88.153 to 00:0c:29:86:50:50 via ens32
  Feb 25 10:20:37 localhost dhcpd: DHCPREQUEST for 192.168.88.153 (192.168.88.129) from 00:0c:29:86:50:50 via ens32
  Feb 25 10:20:37 localhost dhcpd: DHCPACK on 192.168.88.153 to 00:0c:29:86:50:50 via ens32

  # 服务器抓包 4步 Discover Offer Request ACK
  $ $ tcpdump -nn -vvv -s 1500 -i ens32 host 192.168.88.129 and udp port 67 or udp port 68
  11:09:02.574883 IP (tos 0x10, ttl 128, id 0, offset 0, flags [none], proto UDP (17), length 328)
      0.0.0.0.68 > 255.255.255.255.67: [udp sum ok] BOOTP/DHCP, Request from 00:0c:29:86:50:46, length 300, xid 0xa18cd146, Flags [none] (0x0000)
            Client-Ethernet-Address 00:0c:29:86:50:46
            Vendor-rfc1048 Extensions
              Magic Cookie 0x63825363
              DHCP-Message Option 53, length 1: Discover
              Requested-IP Option 50, length 4: 192.168.88.157
              Parameter-Request Option 55, length 13:
                Subnet-Mask, BR, Time-Zone, Classless-Static-Route
                Domain-Name, Domain-Name-Server, Hostname, YD
                YS, NTP, MTU, Option 119
                Default-Gateway
              END Option 255, length 0
              PAD Option 0, length 0, occurs 35
  11:09:03.577000 IP (tos 0x10, ttl 128, id 0, offset 0, flags [none], proto UDP (17), length 328)
      192.168.88.129.67 > 192.168.88.157.68: [udp sum ok] BOOTP/DHCP, Reply, length 300, xid 0xa18cd146, Flags [none] (0x0000)
            Your-IP 192.168.88.157
            Client-Ethernet-Address 00:0c:29:86:50:46
            Vendor-rfc1048 Extensions
              Magic Cookie 0x63825363
              DHCP-Message Option 53, length 1: Offer
              Server-ID Option 54, length 4: 192.168.88.129
              Lease-Time Option 51, length 4: 7200
              Subnet-Mask Option 1, length 4: 255.255.255.0
              BR Option 28, length 4: 192.168.88.255
              Domain-Name-Server Option 6, length 8: 202.106.0.20,114.114.114.114
              Default-Gateway Option 3, length 4: 192.168.88.2
              END Option 255, length 0
              PAD Option 0, length 0, occurs 16
  11:09:03.577505 IP (tos 0x10, ttl 128, id 0, offset 0, flags [none], proto UDP (17), length 328)
      0.0.0.0.68 > 255.255.255.255.67: [udp sum ok] BOOTP/DHCP, Request from 00:0c:29:86:50:46, length 300, xid 0xa18cd146, Flags [none] (0x0000)
            Client-Ethernet-Address 00:0c:29:86:50:46
            Vendor-rfc1048 Extensions
              Magic Cookie 0x63825363
              DHCP-Message Option 53, length 1: Request
              Server-ID Option 54, length 4: 192.168.88.129
              Requested-IP Option 50, length 4: 192.168.88.157
              Parameter-Request Option 55, length 13:
                Subnet-Mask, BR, Time-Zone, Classless-Static-Route
                Domain-Name, Domain-Name-Server, Hostname, YD
                YS, NTP, MTU, Option 119
                Default-Gateway
              END Option 255, length 0
              PAD Option 0, length 0, occurs 29
  11:09:03.579333 IP (tos 0x10, ttl 128, id 0, offset 0, flags [none], proto UDP (17), length 328)
      192.168.88.129.67 > 192.168.88.157.68: [udp sum ok] BOOTP/DHCP, Reply, length 300, xid 0xa18cd146, Flags [none] (0x0000)
            Your-IP 192.168.88.157
            Client-Ethernet-Address 00:0c:29:86:50:46
            Vendor-rfc1048 Extensions
              Magic Cookie 0x63825363
              DHCP-Message Option 53, length 1: ACK
              Server-ID Option 54, length 4: 192.168.88.129
              Lease-Time Option 51, length 4: 7200
              Subnet-Mask Option 1, length 4: 255.255.255.0
              BR Option 28, length 4: 192.168.88.255
              Domain-Name-Server Option 6, length 8: 202.106.0.20,114.114.114.114
              Default-Gateway Option 3, length 4: 192.168.88.2
              END Option 255, length 0
              PAD Option 0, length 0, occurs 16
  ```
6. 保留地址
    查看租约地址 客户端和服务器端都可以看到
  ```shell
  $  cat /var/lib/dhclient/dhclient.leases
  ```
  将IP地址保留给特定MAC地址,例如文件服务器等等,继续以测试机ens33为例,将192,168.88.252固定分配给他
  ```shell
  编辑dhcp配置文件,添加配置
  $ vim /etc/dhcp/dhcpd.conf
  host printer { # 名字可以随意,最好标明为谁保留,这里是为printer保留
    hardware ethernet 00:0c:29:86:50:46; # 对应的MAC地址
    fixed-address 192.168.88.252; # 指定的IP
    option host-name "printer.example.com" # 最好在内网中和DNS配套使用,当DNS对这个域名做了解析 机器启动时便会获得这个hostname. 需要把/etc/hostname 中的计算机名称清除 /etc/sysconfig/network中的hostname字段清除
  }
  
  $ systemctl restart dhcpd
  ```
  在测试机上重新尝试获取IP地址
  ```shell
  $ dhclient -r ens33
  $ dhclient -d ens33
  Internet Systems Consortium DHCP Client 4.2.5
  Copyright 2004-2013 Internet Systems Consortium.
  All rights reserved.
  For info, please visit https://www.isc.org/software/dhcp/

  Listening on LPF/ens33/00:0c:29:86:50:46
  Sending on   LPF/ens33/00:0c:29:86:50:46
  Sending on   Socket/fallback
  DHCPDISCOVER on ens33 to 255.255.255.255 port 67 interval 4 (xid=0xf35e714)
  DHCPREQUEST on ens33 to 255.255.255.255 port 67 (xid=0xf35e714)
  DHCPOFFER from 192.168.88.129
  DHCPACK from 192.168.88.129 (xid=0xf35e714)
  bound to 192.168.88.252 -- renewal in 3369 seconds.
  # 可以看到252 被分配过来了
  $ ip addr show ens33
  2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
      link/ether 00:0c:29:86:50:46 brd ff:ff:ff:ff:ff:ff
      inet 192.168.88.252/24 brd 192.168.88.255 scope global dynamic ens33
        valid_lft 7165sec preferred_lft 7165sec
      inet6 fe80::7137:324e:eb37:4d39/64 scope link noprefixroute
        valid_lft forever preferred_lft forever
  ```

7. shared-network
 使用shared-network可以分配多网段
 ```shell
 $ vim /etc/dhcp/dhcpd.conf
 shared-network super {
  subnet 192.168.88.0 netmask 255.255.255.0 {
    range 192.168.88.153 192.168.88.153;
    option routers 192.168.88.2;
  }
  subnet 192.168.88.0 netmask 255.255.255.0 {
    range 192.168.99.153 192.168.99.153;
    option routers 192.168.99.2;
  }

  }

 ```


