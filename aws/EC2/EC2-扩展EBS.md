# EC2-扩展EBS 


1. 验证每个卷的文件系统,使用 df -hT 命令
    ```shell
    ubuntu@ip-10-100-0-79:~$ df -hT
    Filesystem     Type      Size  Used Avail Use% Mounted on
    /dev/root      ext4       49G  1.9G   47G   4% /
    devtmpfs       devtmpfs  465M     0  465M   0% /dev
    ...
    ```
    可以看到,目前/root大小为49G,type为ext4

2. 使用 lsblk 命令显示有关附加到实例的块储存设备的信息
    ```shell
    ubuntu@ip-10-100-0-79:~$ lsblk
    NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    loop0         7:0    0   55M  1 loop /snap/core18/1754
    loop1         7:1    0   18M  1 loop /snap/amazon-ssm-agent/1566
    loop2         7:2    0 93.9M  1 loop /snap/core/9066
    loop3         7:3    0 97.9M  1 loop /snap/core/10583
    loop4         7:4    0 71.4M  1 loop /snap/lxd/19300
    loop6         7:6    0 55.5M  1 loop /snap/core18/1988
    loop7         7:7    0 32.3M  1 loop /snap/amazon-ssm-agent/2996
    loop8         7:8    0 71.4M  1 loop /snap/lxd/19164
    nvme0n1     259:0    0   55G  0 disk 
    └─nvme0n1p1 259:1    0   50G  0 part /
    ```
    可以看到根卷 /dev/nvme0n1 具有一个分区 /dev/nvme0n1p1。当根卷的大小反映新大小 55 GB 时,分区的大小会反映原始大小 50 GB 并且必须先进行扩展,然后才能扩展文件系统。

3. 在根卷上扩展分区,使用 growpart 命令
    ```shell
    ubuntu@ip-10-100-0-79:~$ sudo growpart /dev/nvme0n1 1
    CHANGED: partition=1 start=2048 old: size=104855519 end=104857567 new: size=115341279 end=115343327
    ``` 
    /dev/nvme0n1表示磁盘,后边的空格数字`1`表示第一个分区.再次查看块储存设备的信息
    ```shell
    ubuntu@ip-10-100-0-79:~$ lsblk
    NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    loop0         7:0    0   55M  1 loop /snap/core18/1754
    loop1         7:1    0   18M  1 loop /snap/amazon-ssm-agent/1566
    loop2         7:2    0 93.9M  1 loop /snap/core/9066
    loop3         7:3    0 97.9M  1 loop /snap/core/10583
    loop4         7:4    0 71.4M  1 loop /snap/lxd/19300
    loop6         7:6    0 55.5M  1 loop /snap/core18/1988
    loop7         7:7    0 32.3M  1 loop /snap/amazon-ssm-agent/2996
    loop8         7:8    0 71.4M  1 loop /snap/lxd/19164
    nvme0n1     259:0    0   55G  0 disk 
    └─nvme0n1p1 259:1    0   55G  0 part /
    ```
    可以看到根卷大小为55G,分区大小也是55G.

4. 扩展文件系统(ext4)
    虽然分区已经扩展,但实际文件系统还没有扩展
    ```shell
    ubuntu@ip-10-100-0-79:~$ ubuntu@ip-10-100-0-79:~$ df -hT
    Filesystem     Type      Size  Used Avail Use% Mounted on
    /dev/root      ext4       49G  1.9G   47G   4% /
    devtmpfs       devtmpfs  465M     0  465M   0% /dev
    ...
    ```
    通过`df`可以看到文件系统依然是49G,下面扩展文件系统,可以看到上方文件系统Type是`ext4`,所以使用`resize2fs`

    ```shell
    ubuntu@ip-10-100-0-79:~$ sudo resize2fs /dev/nvme0n1p1
    resize2fs 1.45.5 (07-Jan-2020)
    Filesystem at /dev/nvme0n1p1 is mounted on /; on-line resizing required
    old_desc_blocks = 7, new_desc_blocks = 7
    The filesystem on /dev/nvme0n1p1 is now 14417659 (4k) blocks long.
    ```
    然后就完成了.

5. 扩展文件系统(XFS)
    如果文件系统是XFS格式,则使用 `xfs_growfs` 命令扩展每个卷上的文件系统
    如果尚未安装 XFS 工具,可以按如下方式安装。
    ```shell
    $ sudo yum install xfsprogs
    # 或者
    $ sudo apt install xfsprogs
    ```
    扩展命令,`/`是要扩展的挂载点
    ```shell
    $ sudo xfs_growfs -d /
    ```
    然后就完成了


### 官方链接
[调整卷大小后扩展 Linux 文件系统](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/recognize-expanded-volume-linux.html#extend-linux-volume-partition)
