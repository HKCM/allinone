#  为AWS EC2挂载新EBS卷后使新卷可用
 

##  查看当前快设备
```shell
ubuntu@ip-10-100-0-79:~$ lsblk
NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
loop0         7:0    0   55M  1 loop /snap/core18/1754
loop1         7:1    0   18M  1 loop /snap/amazon-ssm-agent/1566
loop3         7:3    0 97.9M  1 loop /snap/core/10583
loop4         7:4    0 71.4M  1 loop /snap/lxd/19300
loop5         7:5    0 98.4M  1 loop /snap/core/10823
loop6         7:6    0 55.5M  1 loop /snap/core18/1988
loop7         7:7    0 32.3M  1 loop /snap/amazon-ssm-agent/2996
loop8         7:8    0 71.4M  1 loop /snap/lxd/19164
nvme0n1     259:0    0   55G  0 disk 
└─nvme0n1p1 259:1    0   55G  0 part /
nvme1n1     259:2    0   10G  0 disk 
```
可以看到`nvme0n1`是其中一个存储设备,`nvme0n1p1`是根目录的分区.`nvme1n1`是另外一个存储设备并且没有分区.`nvme1n1`为原始的块储存设备,必须先在这种设备上创建文件系统,然后才能够挂载并使用它们.
也可以使用 `file -s` 命令获取设备信息,例如其文件系统类型。如果输出仅显示 data(如以下示例输出),则说明设备上没有文件系统,必须创建一个文件系统。
```shell
ubuntu@ip-10-100-0-79:~$ sudo file -s /dev/nvme1n1
/dev/nvme1n1: data
``` 
## 创建文件系统
### 创建文件系统最好和原本文件系统相同,这里以`ext4`为例.
```shell
ubuntu@ip-10-100-0-79:~$ sudo mkfs -t ext4 /dev/nvme1n1
mke2fs 1.45.5 (07-Jan-2020)
Creating filesystem with 2621440 4k blocks and 655360 inodes
Filesystem UUID: 3c01c114-e29d-4166-a703-3b5db31c7253
Superblock backups stored on blocks: 
    32768, 98304, 163840, 229376, 294912, 819200, 884736, 1605632

Allocating group tables: done                            
Writing inode tables: done                            
Creating journal (16384 blocks): done
Writing superblocks and filesystem accounting information: done 
```
再次查看
```shell
ubuntu@ip-10-100-0-79:~$ sudo file -s /dev/nvme1n1
/dev/nvme1n1: Linux rev 1.0 ext4 filesystem data, UUID=3c01c114-e29d-4166-a703-3b5db31c7253 (extents) (64bit) (large files) (huge files)
```
可以看到`ext4`文件系统已经创建好了.

### 创建`xfs`文件系统(可选)
```shell
sudo mkfs -t xfs /dev/nvme1n1
```
再次查看
```shell
ubuntu@ip-10-100-0-79:~$ sudo file -s /dev/nvme1n1
/dev/nvme1n1: SGI XFS filesystem data (blksz 4096, inosz 512, v2 dirs)
```

如果出现“找不到 mkfs.xfs”错误,请使用以下命令安装 XFS 工具,然后重复上一命令:
```shell
$ sudo yum install xfsprogs
# 或
$ sudo apt install xfsprogs
```

## 创建挂载目录
```shell
ubuntu@ip-10-100-0-79:~$ sudo mkdir /data
ubuntu@ip-10-100-0-79:~$ sudo mount /dev/nvme1n1 /data
```

查看文件系统类型和挂在情况
```shell
ubuntu@ip-10-100-0-79:~$ df -hT
Filesystem     Type      Size  Used Avail Use% Mounted on
/dev/root      ext4       54G  2.0G   52G   4% /
devtmpfs       devtmpfs  465M     0  465M   0% /dev
tmpfs          tmpfs     477M     0  477M   0% /dev/shm
tmpfs          tmpfs      96M  812K   95M   1% /run
tmpfs          tmpfs     5.0M     0  5.0M   0% /run/lock
tmpfs          tmpfs     477M     0  477M   0% /sys/fs/cgroup
/dev/loop1     squashfs   18M   18M     0 100% /snap/amazon-ssm-agent/1566
/dev/loop0     squashfs   55M   55M     0 100% /snap/core18/1754
/dev/loop6     squashfs   56M   56M     0 100% /snap/core18/1988
/dev/loop3     squashfs   98M   98M     0 100% /snap/core/10583
/dev/loop7     squashfs   33M   33M     0 100% /snap/amazon-ssm-agent/2996
/dev/loop8     squashfs   72M   72M     0 100% /snap/lxd/19164
tmpfs          tmpfs      96M     0   96M   0% /run/user/1001
/dev/loop4     squashfs   72M   72M     0 100% /snap/lxd/19300
tmpfs          tmpfs      96M     0   96M   0% /run/user/1000
/dev/loop5     squashfs   99M   99M     0 100% /snap/core/10823
/dev/nvme1n1   xfs        10G  104M  9.9G   2% /data

ubuntu@ip-10-100-0-79:~$ lsblk
NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
loop0         7:0    0   55M  1 loop /snap/core18/1754
loop1         7:1    0   18M  1 loop /snap/amazon-ssm-agent/1566
loop3         7:3    0 97.9M  1 loop /snap/core/10583
loop4         7:4    0 71.4M  1 loop /snap/lxd/19300
loop5         7:5    0 98.4M  1 loop /snap/core/10823
loop6         7:6    0 55.5M  1 loop /snap/core18/1988
loop7         7:7    0 32.3M  1 loop /snap/amazon-ssm-agent/2996
loop8         7:8    0 71.4M  1 loop /snap/lxd/19164
nvme0n1     259:0    0   55G  0 disk 
└─nvme0n1p1 259:1    0   55G  0 part /
nvme1n1     259:2    0   10G  0 disk /data
```

## 重启后自动挂载附加的卷

### 创建 `/etc/fstab` 文件的备份
```shell
$ sudo cp /etc/fstab /etc/fstab.orig
```

### 使用 `blkid` 命令查找设备的 UUID
```shell
ubuntu@ip-10-100-0-79:~$ sudo blkid
/dev/nvme0n1p1: LABEL="cloudimg-rootfs" UUID="fdd49fba-0340-4ed1-b0fc-8da187913fec" TYPE="ext4" PARTUUID="093eb684-01"
/dev/loop0: TYPE="squashfs"
/dev/loop1: TYPE="squashfs"
/dev/loop3: TYPE="squashfs"
/dev/loop4: TYPE="squashfs"
/dev/loop5: TYPE="squashfs"
/dev/loop6: TYPE="squashfs"
/dev/loop7: TYPE="squashfs"
/dev/loop8: TYPE="squashfs"
/dev/nvme1n1: UUID="f55bd613-b470-45b5-8cbb-7713ac77ddd5" TYPE="xfs"
```
### 或者使用lsblk, 查看lsblk --help
```shell
ubuntu@ip-10-100-0-79:~$ sudo lsblk -o +UUID
NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT                  UUID
loop0         7:0    0   55M  1 loop /snap/core18/1754           
loop1         7:1    0   18M  1 loop /snap/amazon-ssm-agent/1566 
loop3         7:3    0 97.9M  1 loop /snap/core/10583            
loop4         7:4    0 71.4M  1 loop /snap/lxd/19300             
loop5         7:5    0 98.4M  1 loop /snap/core/10823            
loop6         7:6    0 55.5M  1 loop /snap/core18/1988           
loop7         7:7    0 32.3M  1 loop /snap/amazon-ssm-agent/2996 
loop8         7:8    0 71.4M  1 loop /snap/lxd/19164             
nvme0n1     259:0    0   55G  0 disk                             
└─nvme0n1p1 259:1    0   55G  0 part /                           fdd49fba-0340-4ed1-b0fc-8da187913fec
nvme1n1     259:2    0   10G  0 disk /data                       f55bd613-b470-45b5-8cbb-7713ac77ddd5
```
### 编辑/etc/fstab文件
在`etc/fstab`文件中添加条目
```shell
UUID=aebf131c-6957-451e-8d34-ec978d9581ae  /data  xfs  defaults,nofail  0  2
```
### 检测
卸载设备并重新挂载测试有效性
```shell
ubuntu@ip-10-100-0-79:~$ sudo umount /data
ubuntu@ip-10-100-0-79:~$ sudo mount -a
```

如果出现错误,则还原`/etc/fstab`
```shell
$ sudo mv /etc/fstab.orig /etc/fstab
```
