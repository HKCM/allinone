### EFS简介

Amazon Elastic File System (Amazon EFS) 可提供简单、可扩展、完全托管的弹性 NFS 文件系统,可在不中断应用程序的情况下按需扩展到 PB 级,随着添加或删除文件而自动扩展或缩减,无需预置和管理容量,可自适应增长。支持传输中加密和静态加密.支持网络文件系统版本 4(NFSv4.1 和 NFSv4.0)协议,多个 Amazon EC2 实例可以同时访问 Amazon EFS 文件系统,为在多个实例或服务器上运行的工作负载和应用程序提供通用数据源。

Amazon EFS 提供两种存储类别:Standard 和 Infrequent Access。EFS 生命周期管理是文件系统管理经济高效的文件存储。启用后,生命周期管理会将在一段设定时间内(7,14,30,60和90天)未访问的文件迁移到不常访问 (IA) 存储类别。费用降低10倍,读取是会按读取的文件大小收费.小于 128KB 的文件不符合生命周期管理条件,其将始终存储在标准类别中.

创建文件系统后,默认情况下,只有根用户 (UID 0) 具有读取、写入和执行权限。为了让其他用户也能修改文件系统,根用户必须明确授予他们访问权限。可以使用访问点自动创建非根用户可从中写入的目录.

### 性能模式

两种性能模式: 
* 默认通用性能模式,非常适合对延迟敏感的使用案例,如 Web 服务环境、内容管理系统、主目录和一般文件服务.
* 最大 I/O 模式下的文件系统可以扩展到更高级别的聚合吞吐量和每秒操作数,但代价是稍高的文件元数据操作延迟. 如大数据分析、媒体处理和基因组分析等高度并行化的应用程序和工作负载可以受益于这种模式.

[如何确定性能模式选择](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/performance.html#throughput-modes)

两种吞吐量模式:
* 默认突增吞吐量模式,吞吐量随着文件系统的增长而扩展.
* 预置吞吐量模式,可以指定与存储的数据量无关的文件系统的吞吐量,贵一些.

[如何确认吞吐量模式](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/performance.html#throughput-modes)
    

EFS定价与存储类别和预置IO有关,与性能模式无关

### EFS与EBS

| | | |
|-|-|-|
|存储特性比较   |Amazon EFS                                                       |Amazon EBS 预置 IOPS 
|可用性与持久性 |数据冗余存储在多个可用区中。                                     |数据冗余存储在一个可用区中
|访问           |多个可用区的多达数千个 Amazon EC2 实例可以同时连接到一个文件系统 |一个可用区的一个 Amazon EC2 实例可以连接到一个文件系统
|使用案例	    |大数据和分析、媒体处理工作流、内容管理、Web 服务和主目录	      |引导卷、事务型数据库和 NoSQL 数据库、数据仓库和 ETL

| | |
|-|-|
|性能比较      |Amazon EFS	    |Amazon EBS 预置 IOPS
|每次操作的延迟|低且一致的延迟  |最低且一致的延迟
|吞吐量规模	   |每秒 10+ GB	    |每秒最多 2 GB

| | |
|-|-|
|简单价格比较   |Amazon EFS	    |Amazon EBS
|标准存储(GB/月)|0.3 USD        |0.1 USD (SSDgp2)

Amazon EFS 文件系统的性能特征不依赖于使用 EBS 优化的实例

### 定价

| | |
|-|-|
标准存储(GB/月)	              |0.30 USD
不频繁访问存储类(GB/月)	      | 0.025 USD
不频繁访问请求(根据传输的 GB 数)|	0.01 USD
预置吞吐量(MB/s/月)	          |6.00 USD
[EFS定价及示例](https://aws.amazon.com/cn/efs/pricing/)

### 挂载EFS

一次只能在一个 VPC 中基于 Amazon VPC 服务使用 Amazon EFS 文件系统。也就是说,在 VPC 中为文件系统创建挂载目标,并使用这些挂载目标提供对该文件系统的访问权限。想要切换至不同的VPC,需要先将所有VPC的挂载点删除,该操作不会影响EFS内的数据.

谁可以挂载:
* 同一 VPC 中的 Amazon EC2 实例
* VPC 中通过 VPC 对等连接的 EC2 实例,当跨VPC挂载时可用区 ID需要与EC2 所在的可用区 ID相匹配方可挂载
* 通过使用 AWS Direct Connect 的本地服务器
* 使用 Amazon VPC 通过 AWS 虚拟专用网络 (VPN) 的本地服务器
* *Amazon EFS 不支持从 Amazon EC2 Windows 实例挂载*

挂载EFS可以使用efs挂载和nfs挂载,efs挂载可以很简单地配置在传输过程中加密

#### EFS挂载先决条件
Amazon EFS client(amazon-efs-utils软件包)是 Amazon EFS 工具, [GitHub地址](https://github.com/aws/efs-utils).

Ubuntu安装EFS客户端:
```shell
$ sudo apt update && sudo apt-get -y install binutils
$ git clone https://github.com/aws/efs-utils
$ cd efs-utils && ./build-deb.sh
$ sudo apt-get -y install ./build/amazon-efs-utils*deb
```

##### efs手动挂载
运行以下命令以挂载文件系统。
```shell
$ sudo mount -t efs fs-12345678:/ /mnt/efs
```

或者,如果要使用传输中的数据加密,可以使用以下命令挂载文件系统。
```shell
$ sudo mount -t efs -o tls fs-12345678:/ /mnt/efs
```

访问点挂载
```shell
mount -t efs -o tls,accesspoint=fsap-12345678 fs-12345678: /localmountpoint
```

##### efs自动挂载

正常自动挂载文件系统,请将以下行添加到 /etc/fstab 文件中。
```shell
file-system-id efs-mount-point efs _netdev,tls 0 0
```

要使用 IAM 授权自动挂载到具有实例配置文件的 Amazon EC2 实例,请将以下行添加 /etc/fstab 文件中。
```shell
file-system-id:/ efs-mount-point efs _netdev,tls,iam 0 0
```

要使用凭证文件通过 IAM 授权自动挂载到 Linux 实例,请将以下行添加到 /etc/fstab 文件中。
```shell
file-system-id:/ efs-mount-point efs _netdev,tls,iam,awsprofile=namedprofile 0 0
```

要使用 EFS 访问点自动挂载文件系统,请将以下行添加到 /etc/fstab 文件中。
```shell
file-system-id efs-mount-point efs _netdev,tls,accesspoint=access-point-id 0 0
```

通过将带 'fake' 选项的 mount 命令与 'all' 和 'verbose' 选项结合使用来测试 fstab 条目。
```shell
$ sudo mount -fav
home/ec2-user/efs      : successfully mounted
```

[使用 AWS Systems Manager 将 EFS 挂载到多个 EC2 实例](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/mount-multiple-ec2-instances.html)

[从另一个账户或 VPC 挂载 EFS 文件系统](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/manage-fs-access-vpc-peering.html#mount-fs-different-vpc)

#### nfs挂载

为获得最佳性能以及避免出现各种已知的 NFS 客户端错误,建议使用最新的 Linux 内核。如果使用的是企业 Linux 发行版,建议使用以下版本:

* Amazon Linux 2
* Amazon Linux 2015.09 或更高版本
* RHEL 7.3 或更高版本
* 具有内核 2.6.32-696 或更高版本的 RHEL 6.9
* 所有 Ubuntu 16.04 版本
* 具有内核 3.13.0-83 或更高版本的 Ubuntu 14.04
* SLES 12 Sp2 或更高版本
* 如果使用其他发行版或自定义内核,建议使用内核 4.3 或更高版本。

[NFS挂载建议](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/mounting-fs-nfs-mount-settings.html)

安装nfs客户端
```shell
# Red Hat
sudo yum -y update 
sudo yum -y install nfs-utils

# Ubuntu
sudo apt update
sudo apt -y install nfs-common

sudo service nfs start
sudo service nfs status
```

##### nfs传输加密

https://docs.aws.amazon.com/zh_cn/efs/latest/ug/encryption-in-transit.html#how-encrypt-transit

##### 手动nfs挂载
```shell
# IP挂载,IP可以在EFS控制台和CLI describe-mount-targets中找到
sudo mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport mount-target-IP:/   ~/efs

# 文件系统DNS挂载
sudo mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport file-system-ID.efs.aws-region.amazonaws.com:/mike  /home/mike/mikeEFS

# 挂载目标DNS挂载,同一可用区中删除并新建挂载目标不会影响挂载目标的DNS
sudo mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport availability-zone.file-system-id.efs.aws-region.amazonaws.com
:/mike  /home/mike/mikeEFS


```

##### 自动nfs挂载
将以下行添加到 /etc/fstab 文件中
```shell
file-system-ID.efs.aws-region.amazonaws.com:/ /var/www/html/efs-mount-point   nfs4   defaults
```


卸载
```shell
umount /mnt/efs
```

### 监控EFS

[EFS metric](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/monitoring-cloudwatch.html#efs-metrics)
[CloudTrail](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/logging-using-cloudtrail.html)

|主要指标 |如何确认 |
|-|-|
|吞吐量                    |监控 Sum 指标的每日 TotalIOBytes 统计数据以查看吞吐量
|Amazon EC2 实例链接数数量 |监控 Sum 指标的 ClientConnections 统计数据
|突增积分余额              |监控文件系统的 BurstCreditBalance 指标,如果 BurstCreditBalance 指标的值为零或稳步下降,则使用预置吞吐量
|性能                      |监控 PercentIOLimit 百分比达到或接近 100%,应用程序应使用最大 I/O 性能模式


从存储库安装 amazon-efs-utils 软件包时,可以通过取消 cloudwatch-log 部分中的 # enabled = true 行的注释来手动更新 /etc/amazon/efs/efs-utils.conf 配置文件,CloudWatch日志组: /aws/efs/utils 日志组

EFS的挂载日志在这`/var/log/amazon/efs`. 

### 数据保护

AWS Backup 是一种统一备份服务,旨在简化备份的创建、迁移、恢复和删除,同时提供改进的报告和审核,它可以:
* 配置并审核要备份的 AWS 资源
* 自动备份计划
* 设置保留策略
* 监控所有最近的备份和还原活动

#### 增量备份:

在初始备份期间,将创建整个文件系统的副本。在该文件系统的后续备份期间,只复制已更改、已添加或已删除的文件和目录。

#### 一致性

如果在执行备份时对文件系统进行了修改,则可能会出现不一致,例如重复、偏差或排除的数据。这些修改包括写入、重命名、移动或删除操作。为确保一致的备份,备份过程中暂停修改文件系统的应用程序或进程。或者,将备份安排在不修改文件系统期间。

#### 存储类型

AWS Backup 备份 EFS 文件系统中的所有数据,对于IA类型的备份,不会产生数据访问费用.还原恢复点时,会将所有文件还原到*标准存储类别*

#### 并发备份

AWS Backup 将备份限制为每个资源一个并发备份。因此,如果备份作业已在进行中,则计划备份或按需备份可能会失败。

#### 还原

AWS Backup 可以将恢复点还原到新的 EFS 文件系统或源文件系统。可以执行完全还原,这会还原整个文件系统。或者可以使用部分还原来还原特定的文件和目录。还原特定文件或目录,必须指定与挂载点相关的相对路径。例如,如果文件系统挂载到 /user/home/myname/efs 并且文件路径为 user/home/myname/efs/file1,请输入 /file1。 **路径区分大小写,不能包含特殊字符、通配符和正则表达式字符串**。

执行完全还原或部分还原时,恢复点将还原到根目录下的还原目录 aws-backup-restore_`timestamp-of-restore`并保留原始文件目录结构.当得到想要的还原数据并另行保存后,记得删除这个还原目录.

**还有2种数据备份方式** 
* [EFS to EFS](https://aws.amazon.com/cn/solutions/implementations/efs-to-efs-backup-solution/#)
* [AWS Data pipeline](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/alternative-efs-backup.html)

### 演练

* [演练:设置 Apache Web 服务器并提供 Amazon EFS 文件服务](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/wt2-apache-web-server.html)
* [演练:创建可写的每用户子目录以及配置在重启时自动重新挂载](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/accessing-fs-nfs-permissions-per-user-subdirs.html)
* [演练:使用 AWS Direct Connect 和 VPN 在本地创建和挂载文件系统: 主要修改/etc/host](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/efs-onpremises.html#wt5-step4-install-nfs)
* [演示 强制加密 Amazon EFS 静态文件系统: 将CloudTrail的事件发送到CloudWatch Log然后filter并创建警报](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/efs-enforce-encryption.html)

### 访问点

如果用户从两个不同的 EC2 实例访问 Amazon EFS 文件系统,根据用户的 UID 在这些实例上是相同还是不同,会看到如下所示的不同行为:
* 如果两个 EC2 实例上的用户 IDs 相同,则 Amazon EFS 会将其视为指明同一用户,而不考虑所用的 EC2 实例。从两个 EC2 实例访问文件系统的用户体验相同。
* 如果两个 EC2 实例上的用户 IDs 不相同,则 Amazon EFS 会将其视为不同的用户。从两个不同的 EC2 实例访问 Amazon EFS 文件系统的用户体验不相同。
* 如果不同 EC2 实例上的两个不同用户共享一个 ID,则 Amazon EFS 会将其视为同一个用户。

Ubuntu user的默认用户ID和组ID都是1000

[详细权限说明看这里](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/accessing-fs-nfs-permissions.html)


