
### 安装 amazon-efs-utils 软件包

```shell
sudo apt update
sudo apt install git 
sudo apt-get -y install binutils
git clone https://github.com/aws/efs-utils
cd efs-utils
./build-deb.sh
sudo apt-get -y install ./build/amazon-efs-utils*deb

# https://docs.aws.amazon.com/zh_cn/efs/latest/ug/installing-other-distro.html
```

### 手动挂载
sudo mount -t efs -o tls fs-12345678:/ /mnt/efs

### 自动挂载
```shell
vim /etc/fstab

# 正常挂载
file-system-id:/ /mnt/efs efs _netdev,tls 0 0

# 访问点挂载
file-system-id efs-mount-point efs _netdev,tls,accesspoint=access-point-id 0 0

# https://docs.aws.amazon.com/zh_cn/efs/latest/ug/mount-fs-auto-mount-onreboot.html
```



### 检查挂载
```s
$ sudo mount -fav
# /mnt/efs is already mounted, please run 'mount' command to verify
# /mnt/efs                 : successfully mounted
```

### 卸载
```shell
sudo umount /mnt/efs
```