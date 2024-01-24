## 配置文件

### 配置文件: yum.conf

/etc/yum.conf 为所有仓库提供公共配置
/var/cache/yum 缓存目录

```
[main]
cachedir=/var/cache/yum/$basearch/$releasever  # 下载的RPM包的缓存目录
keepcache=0     # 是否保存cache,1保存,0不保存。
debuglevel=2    # debuglevel(0-10)
logfile=/var/log/yum.log    # yum的日志文件所在的位置
exactarch=1     # 检查平台是否兼容
obsoletes=1     # 检查包是否废弃
gpgcheck=1      # 检查来源是否合法,需要有制作者的公钥信息
plugins=1       # 是否启用插件
installonly_limit=5
bugtracker_url=http://bugs.centos.org/set_project.php?project_id=23&ref=http://bugs.centos.org/bug_report_page.php?category=yum
distroverpkg=centos-release
```

### 配置文件: *.repo
yum源配置文件位置为:/etc/yum.repos.d/,文件扩展名为"*.repo",为仓库的指向提供配置. Base/Extras/Updates是默认国外官方源
yum的repo配置文件中可用的变量:
- $releaseversion:当前OS的发行版的主版本号
- $arch:平台类型
- $basearch:基础平台

**字段说明:**
- [base]:仓库名称,一定要放在[]中。本地有多个yum源的时,必须唯一
- name:例如,CentOS-$releasever - Base 仓库说明,相当于它的描述。
- mirrorlist:镜像站点,这个可以注释掉。
- baseurl:yum 源服务器的地址,支持 `ftp://` , `http://` 和 `file:///`
- enabled:此仓库是否生效,如果不写或写成 enabled 则表示此仓库生效,写成 enable=0 则表示此仓库不生效。
- gpgcheck:如果为 1 则表示 RPM 的数字证书生效；如果为 0 则表示 RPM 的数字证书不生效。
- gpgkey:数字证书的公钥文件保存位置。不用修改。

### 源
#### 阿里云的CentOS7 yum源
```shell
curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

# CentOS-Base.repo
#
# The mirror system uses the connecting IP address of the client and the
# update status of each mirror to pick mirrors that are updated to and
# geographically close to the client.  You should use this for CentOS updates
# unless you are manually picking other mirrors.
#
# If the mirrorlist= does not work for you, as a fall back you can try the 
# remarked out baseurl= line instead.
#
#
 
[base]
name=CentOS-$releasever - Base - mirrors.aliyun.com
failovermethod=priority
baseurl=http://mirrors.aliyun.com/centos/$releasever/os/$basearch/
        http://mirrors.aliyuncs.com/centos/$releasever/os/$basearch/
        http://mirrors.cloud.aliyuncs.com/centos/$releasever/os/$basearch/
gpgcheck=1
gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-7
 
#released updates 
[updates]
name=CentOS-$releasever - Updates - mirrors.aliyun.com
failovermethod=priority
baseurl=http://mirrors.aliyun.com/centos/$releasever/updates/$basearch/
        http://mirrors.aliyuncs.com/centos/$releasever/updates/$basearch/
        http://mirrors.cloud.aliyuncs.com/centos/$releasever/updates/$basearch/
gpgcheck=1
gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-7
 
#additional packages that may be useful
[extras]
name=CentOS-$releasever - Extras - mirrors.aliyun.com
failovermethod=priority
baseurl=http://mirrors.aliyun.com/centos/$releasever/extras/$basearch/
        http://mirrors.aliyuncs.com/centos/$releasever/extras/$basearch/
        http://mirrors.cloud.aliyuncs.com/centos/$releasever/extras/$basearch/
gpgcheck=1
gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-7
 
#additional packages that extend functionality of existing packages
[centosplus]
name=CentOS-$releasever - Plus - mirrors.aliyun.com
failovermethod=priority
baseurl=http://mirrors.aliyun.com/centos/$releasever/centosplus/$basearch/
        http://mirrors.aliyuncs.com/centos/$releasever/centosplus/$basearch/
        http://mirrors.cloud.aliyuncs.com/centos/$releasever/centosplus/$basearch/
gpgcheck=1
enabled=0
gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-7
 
#contrib - packages by Centos Users
[contrib]
name=CentOS-$releasever - Contrib - mirrors.aliyun.com
failovermethod=priority
baseurl=http://mirrors.aliyun.com/centos/$releasever/contrib/$basearch/
        http://mirrors.aliyuncs.com/centos/$releasever/contrib/$basearch/
        http://mirrors.cloud.aliyuncs.com/centos/$releasever/contrib/$basearch/
gpgcheck=1
enabled=0
gpgkey=http://mirrors.aliyun.com/centos/RPM-GPG-KEY-CentOS-7
```

#### epel源

```shell
curl -o /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo

[epel]
name=Extra Packages for Enterprise Linux 7 - $basearch
baseurl=http://mirrors.aliyun.com/epel/7/$basearch
failovermethod=priority
enabled=1
gpgcheck=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7
 
[epel-debuginfo]
name=Extra Packages for Enterprise Linux 7 - $basearch - Debug
baseurl=http://mirrors.aliyun.com/epel/7/$basearch/debug
failovermethod=priority
enabled=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7
gpgcheck=0
 
[epel-source]
name=Extra Packages for Enterprise Linux 7 - $basearch - Source
baseurl=http://mirrors.aliyun.com/epel/7/SRPMS
failovermethod=priority
enabled=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7
gpgcheck=0
```

## 命令

### 常用命令
```shell
yum search packageName
yum install -y packageName
yum update packageName
yum erase packageName
```

### List of Commands:

- check          检查 RPM 数据库问题
- check-update   检查是否有可用的软件包更新
- clean          删除缓存数据
- deplist        列出软件包的依赖关系
- distribution-synchronization 已同步软件包到最新可用版本
- downgrade      降级软件包
- erase          从系统中移除一个或多个软件包
- fs             Acts on the filesystem data of the host, mainly for removing docs/lanuages for minimal hosts.
- fssnapshot     Creates filesystem snapshots, or lists/deletes current snapshots.
- groups         显示或使用、组信息
- help           显示用法提示
- history        显示或使用事务历史
- info           显示关于软件包或组的详细信息
- install        向系统中安装一个或多个软件包
- list           列出一个或一组软件包
- load-transaction 从文件名中加载一个已存事务
- makecache      创建元数据缓存
- provides       查找提供指定内容的软件包
- reinstall      覆盖安装软件包
- repo-pkgs      将一个源当作一个软件包组,这样我们就可以一次性安装/移除全部软件包。
- repolist       显示已配置的源
- search         在软件包详细信息中搜索指定字符串
- shell          运行交互式的 yum shell
- swap           Simple way to swap packages, instead of using shell
- update         更新系统中的一个或多个软件包
- update-minimal Works like upgrade, but goes to the 'newest' package match which fixes a problem that affects your system
- updateinfo     Acts on repository update information
- upgrade        更新软件包同时考虑软件包取代关系
- version        显示机器和/或可用的源版本。

### Options
-   -h, --help            显示此帮助消息并退出
-   -t, --tolerant        忽略错误
-   -C, --cacheonly       完全从系统缓存运行,不升级缓存
-   -c [config file], --config=[config file] 配置文件路径
-   -R [minutes], --randomwait=[minutes] 命令最长等待时间
-   -d [debug level], --debuglevel=[debug level] 调试输出级别
-   --showduplicates      在 list/search 命令下,显示源里重复的条目
-   -e [error level], --errorlevel=[error level] 错误输出级别
-   --rpmverbosity=[debug level name] RPM 调试输出级别
-   -q, --quiet           静默执行
-   -v, --verbose         详尽的操作过程
-   -y, --assumeyes       回答全部问题为是
-   --assumeno            回答全部问题为否
-   --version             显示 Yum 版本然后退出
-   --installroot=[path]  设置安装根目录
-   --enablerepo=[repo]   启用一个或多个软件源(支持通配符)
-   --disablerepo=[repo]  禁用一个或多个软件源(支持通配符)
-   -x [package], --exclude=[package] 采用全名或通配符排除软件包
-   --disableexcludes=[repo] 禁止从主配置,从源或者从任何位置排除
-   --disableincludes=[repo] disable includepkgs for a repo or for everything
-   --obsoletes           更新时处理软件包取代关系
-   --noplugins           禁用 Yum 插件
-   --nogpgcheck          禁用 GPG 签名检查
-   --disableplugin=[plugin] 禁用指定名称的插件
-   --enableplugin=[plugin] 启用指定名称的插件
-   --skip-broken         忽略存在依赖关系问题的软件包
-   --color=COLOR         配置是否使用颜色
-   --releasever=RELEASEVER 在 yum 配置和 repo 文件里设置 $releasever 的值
-   --downloadonly        仅下载而不安装和更新
-   --downloaddir=DLDIR   指定一个其他文件夹用于保存软件包
-   --setopt=SETOPTS      设置任意配置和源选项
-   --bugfix              Include bugfix relevant packages, in updates
-   --security            Include security relevant packages, in updates
-   --advisory=ADVS, --advisories=ADVS 包括修复给定建议所需的包,如更新
-   --bzs=BZS             Include packages needed to fix the given BZ, in
-                         updates
-   --cves=CVES           Include packages needed to fix the given CVE, in
-                         updates
-   --sec-severity=SEVS, --secseverity=SEVS
-                         Include security relevant packages matching the
-                         severity, in updates


### search

```shell
yum search epel
```

### repolist

显示repo列表及其简要信息
- enabled: 已启用的源
- disabled: 未启用的源

```shell
yum repolist enabled
```

### list

- available: 列出仓库中有的,但尚未安装的所有可用的包
- installed: 列出已经安装的包
- updates: 列出可用的升级
- --showduplicates: 列出所有可用的版本

```shell
# 查看已安装的包
yum list installed

# yum list --showduplicates <package name>
yum list --showduplicates docker-engine
docker-engine.x86_64             1.13.1-1.el7.centos                  rc-docker 
docker-engine.x86_64             17.03.0.ce-1.el7.centos              rc-docker 
docker-engine.x86_64             17.03.1.ce-1.el7.centos              rc-docker 
docker-engine.x86_64             17.04.0.ce-1.el7.centos              rc-docker 
docker-engine.x86_64             17.05.0.ce-1.el7.centos              rc-docker 
```

### install
```shell
# yum install <package name>-<version info>
yum install git-1.8.3.1-23.el7_8
yum list docker
可安装的软件包
docker.x86_64         2:1.13.1-102.git7f2769b.el7.centos      extras
yum install docker-1.13.1-102.git7f2769b.el7.centos # 不要前面的2:

# 将软件下载至指定目录,目录下会下载相关依赖的rpm包
yum -y install --downloadonly --downloaddir=/tmp/httpd httpd
# 然后可以使用rpm安装
rpm -ivh <package name>
```

### check-update

```shell
yum check-update
```

### update

```shell
yum list docker --show-duplicates
已安装的软件包
docker.x86_64          2:1.13.1-103.git7f2769b.el7.centos       @extras
可安装的软件包
docker.x86_64          2:1.13.1-102.git7f2769b.el7.centos       extras 
docker.x86_64          2:1.13.1-103.git7f2769b.el7.centos       extras 
docker.x86_64          2:1.13.1-108.git4ef4b30.el7.centos       extras 
docker.x86_64          2:1.13.1-109.gitcccb291.el7.centos       extras 
docker.x86_64          2:1.13.1-161.git64e9980.el7_8            extras 
docker.x86_64          2:1.13.1-162.git64e9980.el7.centos       extras 
docker.x86_64          2:1.13.1-203.git0be3e21.el7.centos       extras 
docker.x86_64          2:1.13.1-204.git0be3e21.el7              extras 
docker.x86_64          2:1.13.1-205.git7d71120.el7.centos       extras 
docker.x86_64          2:1.13.1-206.git7d71120.el7_9            extras 
docker.x86_64          2:1.13.1-208.git7d71120.el7_9            extras

# 升级到指定版本
yum update docker-1.13.1-108.git4ef4b30.el7.centos
```

### clean

- packages
- headers
- metadata
- dbcache
- all

```shell
yum clean all
```

### downgrade

注意: 有些依赖包被安装之后,无法降级,需要将依赖包也降级或删除
```shell
# 降级到指定版本
yum downgrade docker-1.13.1-102.git7f2769b.el7.centos
```

### erase

删除时尽量不要删除依赖包,因为依赖包可能被其他文件依赖
```shell
yum erase docker
```

### provides

查看这个文件或命令属于哪个包

```shell
yum provides netstat
```

### group

```shell
yum groups list
yum groups install Development tools 
yum groups erase -y Base
```

### history

```shell
# 查看历史执行yum命令
yum history

# 查询历史执行yum命令ID详细信息
yum history info N

# 撤销历史执行过的yum命令
yum history undo N
```

## 实验

### 自建光盘yum源

```shell
# 创建cdrom目录,作为光盘的挂载点并挂载光盘
mkdir /mnt/cdrom && mount /dev/cdrom /mnt/cdrom
# 禁用自带的yum源,添加bak后缀
cd /etc/yum.repos.d/
find ./ -name "*.repo" | while read id; do mv $id ${id}.bak; done

# 移除.bak后缀
for i in `ls | grep .bak`; do mv $i `echo "$i" | awk -F '.bak' '{print $1}'`;done
# 添加配置
cat > CentOS-Media.repo <<EOF
[c7-media]
name=CentOS-$releasever - Media
baseurl=file:///mnt/cdrom   # 指定位置
        file:///media/CentOS/
        file:///media/cdrom/
        file:///media/cdrecorder/
gpgcheck=1  # 启用证书验证
enabled=1  # 启用仓库
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7  # 指定gpg key位置
EOF
```

### 自建局域网yum源

```shell
# 安装httpd服务器和createrepo
yum install -y httpd createrepo

systemctl enable httpd && systemctl status httpd

# 以Vim做实验,下载vim及其相关依赖包
yum -y install --downloadonly --downloaddir=/tmp/vim vim

# 创建RPM包存放目录
mkdir -p /var/www/html/yum/Packages

# 将Vim相关以来放入Packages中
mv /tmp/vim/* /var/www/html/yum/Packages/

# 在yum目录下生成一个repodata文件
cd /var/www/html/yum && createrepo ./

```


创建yum源文件

```shell
cd /etc/yum.repos.d/
# 禁用自带的yum源,添加bak后缀
find ./ -name "*.repo" | while read id; do mv $id ${id}.bak; done

# 移除.bak后缀
# for i in `ls | grep .bak`; do mv $i `echo "$i" | awk -F '.bak' '{print $1}'`;done

# 新添加自建的yum源
cat > test.repo <<EOF
[c7-media]
name=CentOS-$releasever - Media
baseurl=http://localhost/yum
gpgcheck=0
enabled=1
EOF

yum clean all

yum list available
可安装的软件包
gpm-libs.x86_64                    1.20.7-6.el7                                  c7-media
vim-common.x86_64                  2:7.4.629-8.el7_9                             c7-media
vim-enhanced.x86_64                2:7.4.629-8.el7_9                             c7-media
vim-filesystem.x86_64              2:7.4.629-8.el7_9                             c7-media

yum install vim
```

