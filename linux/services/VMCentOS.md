

## 描述: 创建新的CentOS虚拟机的一些注意事项

### 网络

我是是通过VMware虚拟机启动CentOS7的
在编辑->虚拟网络编辑器,可以看到VMnet8的设置,这里主要可以看到`NAT设置`和启用`使用本机DHCP服务将IP地址分配给虚拟机`功能

CentOS7虚拟机启动后没有IP的主要原因也是没有启用`使用本机DHCP服务将IP地址分配给虚拟机`功能.如果不希望每次启动虚拟机IP乱跑,则可以不启用该功能并在系统中配置静态IP

配置静态IP:
```shell 
# 可以看到网卡信息,centos7默认没有ifconfig
$ ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: ens32: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000

# 配置网卡信息,主要是以下几项
$ vim /etc/sysconfig/network-scripts/ifcfg-ens32

BOOTPROTO=static
ONBOOT=yes
IPADDR=192.168.88.129
NETMASK=255.255.255.0

# 取决于虚拟网络编辑器内的设置
GATEWAY=192.168.88.2
DNS1=192.168.88.2

# 重启网络服务就能看到配置的IP了
$ systemctl restart network
```


### 常用软件
```shell
$ yum update -y

$ yum install -y vim net-tools bind-utils 
```

### selinux和firewall
```shell
$ systemctl stop firewalld && systemctl disable firewalld
$ vim /etc/selinux/config
SELINUX=disabled

# 临时关闭selinux
$ setenforce 0
```

### 修改主机名

```shell
$ sudo vim /etc/hostname

$ sudo hostnamectl set-hostname <newhostname>
```

临时修改
```shell
$ sudo hostname <temporaryname>
```
临时修改后后续的session会生效,重启无效.

通过修改文件永久生效+临时修改生效就等于`不重启生效`
hostname存在于/proc/sys/kernel/hostname中

