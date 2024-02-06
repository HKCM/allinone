# fdisk

fdisk 可以对大于2T的磁盘进行GPT分区

```bash
sudo fdisk -lu # 确认磁盘及其分区信息 -l 列出素所有分区表 -u 显示分区数目
sudo fdisk -u /dev/vdb # m 显示帮助
g n w
sudo mkfs -t ext4 /dev/vdb1 # 创建ext4文件系统
# sudo mkfs -t xfs /dev/vdb1 # 创建xfs文件系统

sudo cp /etc/fstab /etc/fstab.bak # 备份/etc/fstab文件

# 挂载分区 注意挂点 和需要挂载的分区
echo `blkid /dev/vdb1 | awk '{print $2}' | sed 's/\"//g'` /mnt ext4 defaults 0 0 >> /etc/fstab

sudo mount -a # 挂载/etc/fstab配置的文件系统

df -Th
```