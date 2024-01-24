# mount

/etc/fstab
```shell

# 挂载光盘 /dev/cdrom > /dev/sr0
$ mkdir /mnt/cdrom && mount -t iso9660 /dev/cdrom /mnt/cdrom
$ umunt /mnt/cdrom

# 挂载U盘
# u盘名字是不确定的,需要先查询 fdisk -l
$ mkdir /mnt/usb && mount -t vfat /dev/sdb1 /munt/usb
```