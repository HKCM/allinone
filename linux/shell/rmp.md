# rpm

```shell
rpm -qa httpd # 查询本机安装的包
rpm -qip /mnt/cdrom/httpd-XX.rpm # 查看未安装包的信息
rpm -ql httpd # 查看软件包相关文件的位置
rpm -qf /etc/httpd/conf # 查看文件属于哪个包
rpm -ivh XXX.rpm # 安装i install,v verbose, h process
rpm -Uvh XXX.rpm # 安装并升级
rpm --force -ivh XXX.rpm # 强制重装，误删文件
rpm --test XXX.rpm # 测试依赖
rpm -e --nodeps XXX.rpm # 卸载本体和依赖
rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7 # 导入证书
rpm --qa | grep gpg # 查看证书
```