

# Dockerfile


### Dockerfile说明

`FROM` 指定基础镜像,必须为第一个命令
```shell
格式:
	FROM <image>
	FROM <image>:<tag>
示例:
	FROM mysql:5.7
注意:
	tag是可选的,如果不使用tag时,会使用latest版本的基础镜像
```

MAINTAINER 镜像维护者的信息
```shell
格式:
	MAINTAINER <name>
示例:
	MAINTAINER Yongxin Li
    MAINTAINER inspur_lyx@hotmail.com
    MAINTAINER Yongxin Li <inspur_lyx@hotmail.com>
```

COPY|ADD 添加本地文件到镜像中
```shell
格式:
	COPY <src>... <dest>
示例:
    ADD hom* /mydir/          # 添加所有以"hom"开头的文件
    ADD test relativeDir/     # 添加 "test" 到 `WORKDIR`/relativeDir/
    ADD test /absoluteDir/    # 添加 "test" 到 /absoluteDir/
```

WORKDIR 工作目录
```shell
格式:
	WORKDIR /path/to/workdir
示例:
    WORKDIR /a  (这时工作目录为/a)
注意:
	通过WORKDIR设置工作目录后,Dockerfile中其后的命令RUN、CMD、ENTRYPOINT、ADD、COPY等命令都会在该目录下执行
```

RUN 构建镜像过程中执行命令
```
格式:
	RUN <command>
示例:
    RUN yum install nginx
    RUN pip install django
    RUN mkdir test && rm -rf /var/lib/unusedfiles
注意:
	RUN指令创建的中间镜像会被缓存,并会在下次构建中使用。如果不想使用这些缓存镜像,可以在构建时指定--no-cache参数,如:docker build --no-cache
```

CMD 构建容器后调用,也就是在容器启动时才进行调用
```
格式:
    CMD ["executable","param1","param2"] (执行可执行文件,优先)
    CMD ["param1","param2"] (设置了ENTRYPOINT,则直接调用ENTRYPOINT添加参数)
    CMD command param1 param2 (执行shell内部命令)
示例:
    CMD ["/usr/bin/wc","--help"]
    CMD ping www.baidu.com
注意:
	CMD不同于RUN,CMD用于指定在容器启动时所要执行的命令,而RUN用于指定镜像构建时所要执行的命令。
```

ENTRYPOINT 设置容器初始化命令,使其可执行化
```
格式:
    ENTRYPOINT ["executable", "param1", "param2"] (可执行文件, 优先)
    ENTRYPOINT command param1 param2 (shell内部命令)
示例:
    ENTRYPOINT ["/usr/bin/wc","--help"]
注意:
	ENTRYPOINT与CMD非常类似,不同的是通过docker run执行的命令不会覆盖ENTRYPOINT,而docker run命令中指定的任何参数,都会被当做参数再次传递给ENTRYPOINT。Dockerfile中只允许有一个ENTRYPOINT命令,多指定时会覆盖前面的设置,而只执行最后的ENTRYPOINT指令
```
可以在容器启动时使用`--entrypoint=''`覆盖原有的entrypoint

ENV
```
格式:
    ENV <key> <value>
    ENV <key>=<value>
示例:
    ENV myName John
    ENV myCat=fluffy
```

EXPOSE
```
格式:
    EXPOSE <port> [<port>...]
示例:
    EXPOSE 80 443
    EXPOSE 8080
    EXPOSE 11211/tcp 11211/udp
注意:
    EXPOSE并不会让容器的端口访问到主机。要使其可访问,需要在docker run运行容器时通过-p来发布这些端口,或通过-P参数来发布EXPOSE导出的所有端口
```

### 示例1 SSH

```
# 以最新的Ubuntu:20.4镜像为模板
FROM ubuntu:20.04

# PubKey这是需要用到的public key
ENV pubkey "ssh-rsa EXAMPLE/AAAABBBBCCCCc2EAAAADAQABAAABAQDBpC5L7tBkf2U9a6wIM891GUjVgZosERJSiXKHiAMAD34TUp95WN2qDgkno9b3k5FiyODSTd8aseJlBTtCvuChCjS+9tDs009aPoRk14Cwl3QiPvInJCZXYvpomwDqD3lPkMrKjqdVRT/9dDzpBsjXX/Irbm9xRkkt/aEeQCbzJ/X2Q3InwwllGZS5+rquZ8MWaOjKXITL5I3PPS2COFoRUWwsXjfknMwbMudASdeFZoO5rSsIQ7jlG9gJuWM0ZfV2F9M1Ie/hc0rkrQf2JnXOtUXhzBKOOyeSEIChQkTXj+b3tZvWU7wuT7++x2Jr01234567890 NOBODY"

# 换源
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list

# 更新并安装
RUN apt-get update && apt-get install -y openssh-server net-tools vim

# 创建目录并写入PubKey $pubkey可以引用环境变量
RUN mkdir -p /run/sshd && mkdir -p /root/.ssh/ && echo $pubkey > /root/.ssh/authorized_keys

# 开放22端口
EXPOSE 22

CMD /usr/sbin/sshd -D
```

build并启动
```shell
# 在Dockerfile目录下运行
docker build -t sshd:v1 .

# 启动做好的docker镜像
docker run -d -p 10022:22 sshd:v2 

# 从宿主机登录
ssh root@localhost -p 10022
```

### 示例2 带有常见功能的ubuntu

```
FROM ubuntu:latest
RUN apt-get update \
    && apt-get install -y net-tools \
    && apt-get install -y iputils-ping \
    && apt-get install -y vim wget
```
构建
```shell
docker build -t ubuntu:mine .
```

### 示例3 mysql5.7

首先是创建my.cnf文件,用于修改数据字符集
```
$ cat my.cnf
[mysqld]
user=root
character-set-server=utf8
lower_case_table_names=1

[client]
default-character-set=utf8
[mysql]
default-character-set=utf8

!includedir /etc/mysql/conf.d/
!includedir /etc/mysql/mysql.conf.d/
```

Dockerfile
```
FROM mysql:5.7
COPY my.cnf /etc/mysql/my.cnf
# 默认会继承CMD或ENTRYPOINT
```

build

```shell
$ docker build . -t mysql_utf8:5.7 -f Dockerfile
```

启动测试,查看字符集
```shell
$ docker run -d -p 3306:3306 --name mysql -v /opt/mysql/mysql-data/:/var/lib/mysql -e MYSQL_DATABASE=myblog -e MYSQL_ROOT_PASSWORD=123456 mysql:utf8
$ docker exec -it mysql bash
mysql -p123456
mysql> show variables like '%character%';
+--------------------------+----------------------------+
| Variable_name            | Value                      |
+--------------------------+----------------------------+
| character_set_client     | utf8                       |
| character_set_connection | utf8                       |
| character_set_database   | utf8                       |
| character_set_filesystem | binary                     |
| character_set_results    | utf8                       |
| character_set_server     | utf8                       |
| character_set_system     | utf8                       |
| character_sets_dir       | /usr/share/mysql/charsets/ |
+--------------------------+----------------------------+
```

