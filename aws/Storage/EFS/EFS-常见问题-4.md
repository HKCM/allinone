### 无法创建 EFS 文件系统

创建 EFS 文件系统的请求失败,并显示以下消息:
```
User: arn:aws:iam::111122223333:user/username is not authorized to
perform: elasticfilesystem:CreateFileSystem on the specified resource.
```

**要采取的操作**:
这通常是由IAM权限限制引起的,检查当前用户的权限,以及特定规则.
例如:

- 用户是否有权限创建EFS
- 用户创建EFS是否有特定条件:必须加密?必须带某些tags?必须在特定region?



### 凭证进行身份验证时出错

```
AuthFailure: An error occurred authenticating your credentials for XXXXXXXX
```
**要采取的操作**:
- 查看用户的凭证是否正确
- 查看用户客户端时间是否准确



### Amazon EC2 实例挂起

Amazon EC2 实例挂起的原因可能是,在未首先卸载文件系统的情况下删除了文件系统挂载目标。

**要采取的操作**
在删除文件系统挂载目标之前,请卸载文件系统。



### Disk quota exceeded

Amazon EFS 当前不支持用户磁盘配额。如果超出了以下任何限制,则可能会出现该错误:

- 一个实例同一时刻最多可以有 128 个活动用户账户打开文件。

- 一个实例同一时刻最多可打开 32,768 个文件。

实例上的每个唯一挂载可以在 256 个唯一文件/进程对中最多获取总共 8192 个锁。例如,单个进程可以在 256 个单独的文件上获取一个或多个锁,或者说 8 个进程中的每个进程均可以在 32 个文件上获取一个或多个锁。

**要采取的操作**:

如果遇到该问题,可通过确定超出了上述哪个限制,然后进行更改以满足该限制,加以解决。



### I/O error

遇到下列问题之一时会发生此错误:

1. 每个实例同一时刻最多有 128 个活动用户账户打开文件。

    **要采取的操作**

    如果遇到该问题,可以满足在实例上支持的打开文件数限制以解决该问题。为此,请减少在实例上同时打开 Amazon EFS 文件系统中的文件的活动用户数。

2. 已删除加密的文件系统的 AWS KMS 密钥。

    **要采取的操作**:

    如果遇到此问题,则表示不能再解密用该密钥加密的数据,这意味着该数据将无法恢复。



### File name is too long

当文件名或其符号链接 (symbol link) 太长时,会出现该错误。文件名具有以下限制:

* 名称的长度最多为 255 个字节。
* 符号链接的大小最多为 4080 个字节。

**要采取的操作**:

如果遇到该问题,可通过减小文件名或符号链接的长度以满足支持的限制,加以解决。



### Too many links

当文件的硬链接太多时,会出现该错误。一个文件中最多可有 177 个硬链接。

**要采取的操作**:

如果遇到该问题,可通过减少文件硬链接的数量以满足支持的限制,加以解决。



### File too large

当文件太大时,会出现该错误。单个文件的大小最多为 52,673,613,135,872 个字节 (47.9 TiB)。

**要采取的操作**:

如果遇到该问题,可通过减小文件的大小以满足支持的限制,加以解决。



### 无法更改所有权

当使用 Linux chown 命令时,无法更改文件/目录的所有权。

出现该错误的内核版本

*2.6.32*

**要采取的操作**:
内核升级



### 由于客户端错误,文件系统重复执行操作

由于某个客户端错误,文件系统重复执行操作。

**要采取的操作**:

将客户端软件更新为最新版本。



### 客户端发生死锁

客户端变为死锁状态。

出现该错误的内核版本

* CentOS -7,内核为 Linux 3.10.0-229.20.1.el7.x86_64
* 内核为 Linux 4.2.0-18-generic 的 Ubuntu 15.10

**要采取的操作**:

执行下列操作之一:

* 升级为更新的内核版本。对于 CentOS-7,内核版本 Linux 3.10.0-327 或更高版本包含修复。
* 降级为较旧的内核版本。



### 列出大型目录中的文件需要很长时间

如果在 NFS 客户端遍历目录以完成列出操作时,目录正在发生更改,则可能会出现这种情况。每当 NFS 客户端在这种遍历期间注意到目录内容发生更改时,它都会从头开始重新遍历。因此,对于包含经常更改的文件的大型目录,ls 命令可能需要很长时间才能完成。

出现该错误的内核版本

CentOS 和低于 2.6.32-696.el6 的 RHEL 内核版本

**要采取的操作**:

要解决该问题,请升级到较新的内核版本。



### 在 Windows 实例上挂载文件系统失败

在 Microsoft Windows Amazon EC2 实例上挂载文件系统失败。

**要采取的操作**:

不要将 Amazon EFS 与 Windows EC2 实例一起使用,不支持该配置。



### 自动挂载失败,并且实例没有响应

如果在实例上自动挂载文件系统,并且未声明 _netdev 选项,则可能会出现该问题。如果缺少 _netdev EC2 实例可能会停止响应。出现该结果是因为,需要在计算实例启动其网络后初始化网络文件系统。

**要采取的操作**:

如果出现该问题,请与 AWS Support 联系。



### 在 /etc/fstab 中挂载多个 Amazon EFS 文件系统失败

如果实例使用的 systemd 初始化系统在 /etc/fstab 中具有两个或更多 Amazon EFS 条目,有时可能会没有挂载其中的部分或全部条目。在这种情况下,dmesg 输出显示类似于以下内容的一行或多行。

`NFS: nfs4_discover_server_trunking unhandled error -512. Exiting with error EIO`

**要采取的操作**:

在这种情况下,使用以下内容在 /etc/systemd/system/mount-nfs-sequentially.service 中创建新的 systemd 服务文件。
```
[Unit]
Description=Workaround for mounting NFS file systems sequentially at boot time
After=remote-fs.target

[Service]
Type=oneshot
ExecStart=/bin/mount -avt nfs4
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
```
在执行该操作后,运行以下两个命令:
```shell
sudo systemctl daemon-reload

sudo systemctl enable mount-nfs-sequentially.service
```
然后,重新启动 Amazon EC2 实例。通常在一秒内将按需挂载文件系统



### 挂载命令失败,并显示“wrong fs type”错误消息

挂载命令失败,并显示如下错误消息。
```
mount: wrong fs type, bad option, bad superblock on 10.1.25.30:/, 
missing codepage or helper program, or other error (for several filesystems 
(e.g. nfs, cifs) you might need a /sbin/mount.<type> helper program)
In some cases useful info is found in syslog - try dmesg | tail or so.
```
**要采取的操作**:

如果收到该消息,请安装 nfs-utils(或 Ubuntu 上的 nfs-common)软件包



### 挂载命令失败,并显示“incorrect mount option”错误消息

挂载命令失败,并显示如下错误消息。

`mount.nfs: an incorrect mount option was specified`

**要采取的操作**:

该错误消息很可能意味着 Linux 发行版不支持 4.0 和 4.1 版网络文件系统 (NFSv4)。要确认是否属于这种情况,可以运行以下命令。
```shell
$ grep CONFIG_NFS_V4_1 /boot/config*
```
如果上述命令返回 `# CONFIG_NFS_V4_1 is not set`,则表明你的 Linux 发行版不支持 NFSv4.1。有关支持 Amazon Elastic Compute Cloud 的 Amazon EC2 (NFSv4.1,) 的 Amazon 系统映像 (AMI) 的列表,请参阅NFS 支持。



### 在创建文件系统后文件系统挂载立即失败

在创建域名服务 (DNS) 记录的挂载目标后,可能最多需要 90 秒的时间才能在 AWS 区域中完全传播。

**要采取的操作**:

如果以编程方式创建和挂载文件系统(例如,使用 AWS CloudFormation 模板),建议实施等待条件



### 文件系统挂载挂起,然后失败,并显示超时错误

文件系统挂载命令挂起一两分钟,然后失败,并显示超时错误。下面的代码显示了一个示例。
```
$ sudo mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport mount-target-ip:/ mnt

[2+ minute wait here]
mount.nfs: Connection timed out
```

**要采取的操作**:

出现该错误的原因可能是,Amazon EC2 实例或挂载目标安全组的配置不正确。确保挂载目标安全组具有允许从 EC2 安全组进行 NFS 访问的2049入站规则。或者可能挂载的IP不正确



### Name or service not known

使用 DNS 名称的文件系统挂载失败。下面的代码显示了一个示例。
```
$ sudo mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport file-system-id.efs.aws-region.amazonaws.com:/ mnt   
mount.nfs: Failed to resolve server file-system-id.efs.aws-region.amazonaws.com: 
  Name or service not known.   

$  
```
**要采取的操作**:

请检查 VPC 配置。如果使用自定义 VPC,请确保已启用 DNS 设置。
* 确保 Amazon EC2 实例所在的同一可用区中有一个 Amazon EFS 挂载目标
* 确保在与 Amazon EC2 实例相同的 VPC 中有一个挂载目标。否则,不能对位于其他 VPC 中的 EFS 挂载目标使用 DNS 名称解析
* 在配置为使用由 Amazon 提供的 DNS 服务器的 Amazon VPC 内连接至您的 Amazon EC2 实例
* 确保连接 Amazon EC2 实例的 Amazon VPC 已启用 DNS 主机名



### nfs not responding

Amazon EFS 文件系统挂载因传输控制协议 (TCP) 重新连接事件失败,返回错误"nfs: server_name still not responding"。

**要采取的操作**:

请使用 noresvport 挂载选项,以确保在重新建立网络连接时,NFS 客户端将使用新的 TCP 源端口。这样做有助于确保在网络恢复事件后具有不间断的可用性。



### 挂载没有响应

Amazon EFS 挂载看起来没有响应。例如,ls 等命令挂起。

**要采取的操作**:

如果另一个应用程序正在将大量数据写入文件系统,则可能会出现该错误。在该操作完成前,可能会阻止对正在被写入的文件的访问。一般来说,尝试访问正在被写入的文件的任何命令或应用程序均可能会显示为挂起状态。例如,ls 命令可能会在访问正在被写入的文件时挂起。出现该结果是因为,某些 Linux 发行版在 ls 命令中使用别名,以便检索文件属性以及列出目录内容。

要解决此问题,请验证另一个应用程序是否正在将文件写入 Amazon EFS 挂载,并验证它是否处于 Uninterruptible sleepD 状态,如下面的示例所示:
```
$ ps aux | grep large_io.py 
root 33253 0.5 0.0 126652 5020 pts/3 D+ 18:22 0:00 python large_io.py /efs/large_file
```
在已验证确属这种情况之后,您可以通过等待其他写入操作完成或通过实施一种变通解决办法来解决问题。在 ls 示例中,您可以直接使用 /bin/ls 命令,而不是使用别名。这样做可以继续执行命令,而不会在写入的文件处挂起。通常,如果写入数据的应用程序可能会定期强制执行数据刷新(可能使用 fsync(2)),这样做可能有助于提高文件系统对其他应用程序的响应能力。但是,在应用程序写入数据时,这种改善可能会牺牲性能。



### bad file handle

针对新挂载的文件系统执行的操作返回 bad file handle 错误。

如果 Amazon EC2 实例连接到了一个文件系统和一个具有指定 IP 地址的挂载目标,然后该文件系统和挂载目标被删除,则可能会出现该错误。如果您创建新的文件系统和挂载目标,以连接到具有相同挂载目标 IP 地址的 Amazon EC2 实例,则可能会发生该问题。

**要采取的操作**:

可以卸载文件系统,然后在 Amazon EC2 实例上重新挂载文件系统以解决该问题。



### 卸载文件系统失败

如果文件系统繁忙,则无法将其卸载。

**要采取的操作**:

可以通过以下方法解决该问题:

* 等待所有读取和写入操作完成,然后再次尝试执行 umount 命令。
* 使用 umount 选项强制完成 -f 命令。



### 具有传输中的数据加密的挂载失败

默认情况下,当您使用带有传输层安全性 (TLS) 的 Amazon EFS 挂载帮助程序时,它会强制执行主机名检查。某些系统不支持此功能,例如,当您使用 Red Hat Enterprise Linux 或 CentOS 时。 在这些情况下,使用 TLS 挂载 EFS 文件系统会失败。

**要采取的操作**:

我们建议您升级客户端上的 stunnel 版本以支持主机名检查。[升级stunnel](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/upgrading-stunnel.html)



### 具有传输中的数据加密的挂载中断

在极少数情况下,客户端事件可能会导致到您的 Amazon EFS 文件系统的加密连接挂起或中断。

**要采取的操作**:

如果到使用传输中的数据加密的 Amazon EFS 文件系统的连接中断,请执行以下步骤:

* 确保正在客户端上运行 stunnel 服务。
* 确认正在客户端上运行监控程序应用程序 amazon-efs-mount-watchdog。您可以使用以下命令确定是否正在运行该应用程序:
    `ps aux | grep [a]mazon-efs-mount-watchdog`
* 检查日志,可以查找在 /var/log/amazon/efs 中存储的日志。
* (可选)可以启用 stunnel 日志以及检查这些日志中的信息。可以在 /etc/amazon/efs/efs-utils.conf 中更改日志配置以启用 stunnel 日志。但是,这样做需要卸载文件系统,然后使用挂载帮助程序重新挂载以使更改生效。



### 无法创建静态加密的文件系统

您已尝试创建新的静态加密的文件系统。不过,您会收到一条错误消息,指出 AWS KMS 不可用。

**要采取的操作**:

在极少数情况下,AWS KMS 可能在您的 AWS 区域中暂时不可用,从而出现该错误。如果发生这种情况,请等到 AWS KMS 恢复完全可用,然后重试以创建文件系统。



### 无法使用的加密文件系统

加密的文件系统持续返回 NFS 服务器错误。如果由于以下原因之一 EFS 无法从 AWS KMS 中检索主密钥,则可能会出现这些错误:

* 禁用了密钥。
* 删除了密钥。
* 撤销了 Amazon EFS 使用密钥的权限。
* AWS KMS 暂时不可用。

**要采取的操作**:

* 确认已启用 AWS KMS 密钥。为此,您可以在控制台中查看这些密钥。
* 如果未启用密钥,请将其启用。
* 如果密钥处于待删除状态,该状态将禁用密钥。可以取消删除,然后重新启用密钥。
* 如果已启用密钥并且仍遇到问题,或者在重新启用密钥时遇到问题,请与 AWS Support 联系。


