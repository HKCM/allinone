

### 描述: 简单介绍的docker的安装和简单操作


### docker安装

```shell
# 下载阿里源repo文件
$ curl -o /etc/yum.repos.d/Centos-7.repo http://mirrors.aliyun.com/repo/Centos-7.repo
$ curl -o /etc/yum.repos.d/docker-ce.repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

# 查看源中可用版本
$ yum list docker-ce --showduplicates | sort -r
 * extras: mirrors.163.com
docker-ce.x86_64            3:20.10.5-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:20.10.4-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:20.10.3-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:20.10.2-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:20.10.1-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:20.10.0-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.9-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.8-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.7-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.6-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.5-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.4-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.3-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.2-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.15-3.el7                   docker-ce-stable 
docker-ce.x86_64            3:19.03.14-3.el7                   docker-ce-stable 
docker-ce.x86_64            3:19.03.1-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:19.03.13-3.el7                   docker-ce-stable 
docker-ce.x86_64            3:19.03.12-3.el7                   docker-ce-stable 
docker-ce.x86_64            3:19.03.11-3.el7                   docker-ce-stable 
docker-ce.x86_64            3:19.03.10-3.el7                   docker-ce-stable 
docker-ce.x86_64            3:19.03.0-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:18.09.9-3.el7                    docker-ce-stable 
docker-ce.x86_64            3:18.09.9-3.el7                    @docker-ce-stable

# 安装指定版本
$ yum install -y docker-ce-18.09.9

# 配置源加速
# https://cr.console.aliyun.com/cn-hangzhou/instances/mirrors
$ mkdir -p /etc/docker
$ vi /etc/docker/daemon.json
{
	"insecure-registries":["192.168.0.104:5000"], # 自建仓库IP
  "registry-mirrors" : [
    "https://8xpk5wnt.mirror.aliyuncs.com",
    "https://dockerhub.azk8s.cn",
    "https://registry.docker-cn.com",
    "https://ot2k4d59.mirror.aliyuncs.com/"
  ]
}

## 设置开机自启并启动docker
$ systemctl enable docker && systemctl start docker 

```

### docker常见操作

查看docker信息
```shell
$ docker info
```

#### 镜像

查找镜像
```shell
$ docker search httpd
```

拉取镜像
```shell
$ docker pull nginx:alpine
$ docker pull NAME[:TAG]    #NAME是仓库名称,TAG是镜像标签
```

列出本地镜像

```shell
$ docker images
$ docker image ls
```

获取镜像详细信息
```shell
$ docker inspect nginx:alpine
```

设置镜像标签
```shell
$ docker tag <imageID> <Registry>/<Repositories>/<Name>:<Tag>

$ docker tag <imageID> <Name>:<Tag>
```

构建镜像:通过Dockerfile
```shell
$ cat> Dockerfile<<EOF
FROM    centos:6.7
MAINTAINER      Fisher "fisher@sudops.com"

RUN     /bin/echo 'root:123456' |chpasswd
RUN     useradd runoob
RUN     /bin/echo 'runoob:123456' |chpasswd
RUN     /bin/echo -e "LANG=\"en_US.UTF-8\"" >/etc/default/local
EXPOSE  22
EXPOSE  80
CMD     /usr/sbin/sshd -D
EOF

# -t :指定要创建的目标镜像名
# . :Dockerfile 文件所在目录,可以指定Dockerfile 的绝对路径
$ docker build -t runoob/centos:6.7 .
```

构建镜像:通过现有container
```shell
$ docker commit -m "add a test file" -a "docker Newbee" fe76948987e7 test:0.1
```

构建镜像:通过导出和导入
```shell
# 导出
$ docker export <container> > ubuntu.tar

# 导入
$ cat docker/ubuntu.tar | docker import - ubuntu:v1

# 也可以通过指定 URL 或者某个目录来导入
$ docker import http://example.com/exampleimage.tgz example/imagerepo
```

构建镜像: 打包和解压
```shell
# 将本地镜像打包
$ docker -o save karl.tar karl:0.1
# 或
$ docker save karl:0.1 > karl:0.1.tar

# 将打包的镜像解压到本地镜像并打tag
$ docker load --input karl.tar
$ docker tag <ImageID> karl:0.1
```

部署个人镜像仓库

https://docs.docker.com/registry/
```shell
## 使用docker镜像启动镜像仓库服务
$ docker run -d -p 5000:5000 --restart always -v /opt/registry-data/registry:/var/lib/registry --name registry registry:2
   
# 默认仓库不带认证,若需要认证,参考https://docs.docker.com/registry/deploying/#restricting-access
# docker默认不允许向http的仓库地址推送,如何做成https的,参考:https://docs.docker.com/registry/deploying/#run-an-externally-accessible-registry
# 通过配置daemon的方式,来跳过证书的验证:
$ cat /etc/docker/daemon.json
   {
     "registry-mirrors": [
       "https://8xpk5wnt.mirror.aliyuncs.com"
     ],
     "insecure-registries": [
        "192.168.0.104:5000"
     ]
   }
$ docker tag nginx:alpine 192.168.0.104:5000/nginx:alpine
$ docker push 192.168.0.104:5000/nginx:alpine

# nginx镜像已经存到了宿主机的共享文件
$ ls /opt/registry-data/registry/docker/registry/v2/repositories/
nginx
```

推送镜像
```shell
$ docker push <Registry>/<Repositories>/<Name>:<Tag>
$ docker push 192.168.0.104:5000/nginx:alpine
```

删除镜像
```shell
$ docker rmi nginx:alpine
```

#### 容器

查看容器列表
```shell
# 查看运行状态的容器列表
$ docker ps

# 查看全部状态的容器列表
$ docker ps -a

# 查看容器log
$ docker logs <container>

# 持续查看容器log
$ docker logs -f <container>

# 查看 Docker 的底层信息
$ docker inspect bf08b7f2cd89
```

启动容器
```shell 
# -d 后台启动
$ docker run --name nginx -d -p 8080:80 nginx:alpine

# 启动容器的同时进入容器,-ti与/bin/sh或者/bin/bash配套使用,-t分配一个tty终端 -i 交互输入
$ docker run --name nginx -ti nginx:alpine /bin/sh

# 映射端口,把容器的端口映射到宿主机中,-p <host_port>:<container_port>
$ docker run --name nginx -d -p 8080:80 nginx:alpine
```

第一个docker程序
```shell
$ docker run ubuntu:18.04 /bin/echo "Hello world"
Unable to find image 'ubuntu:18.04' locally
18.04: Pulling from library/ubuntu
23884877105a: Pull complete 
bc38caa0f5b9: Pull complete 
2910811b6c42: Pull complete 
36505266dcc6: Pull complete 
Digest: sha256:3235326357dfb65f1781dbc4df3b834546d8bf914e82cce58e6e6b676e23ce8f
Status: Downloaded newer image for ubuntu:18.04
Hello world
```

交互式容器
```shell
# -t: 在新容器内指定一个伪终端或终端。
# -i: 允许你对容器内的标准输入 (STDIN) 进行交互。
# 可以通过运行 exit 命令或者使用 CTRL+D 来退出容器
$ docker run -it ubuntu:18.04 /bin/bash
root@0123ce188bd8:/#
```

容器数据持久化
```shell
# 挂载主机目录, 左边是宿主机目录,右边是容器映射目录
$ docker run --name nginx -d  -v /opt:/opt -v /var/log:/var/log nginx:alpine

# 使用volumes卷
$ docker volume ls
$ docker volume create my-vol
$ docker run --name nginx -d -v my-vol:/opt/my-vol nginx:alpine
```

查看容器日志
```shell
$ docker logs nginx

# 持续查看
$ docker logs -f nginx

# 从最新的10条开始查看,从倒数第9条开始
$ docker logs --tail=10 -f nginx
```

进入容器或者执行容器内的命令
```shell
$ docker exec -ti nginx /bin/sh
$ docker exec -ti <container_id_or_name> hostname
```

查看端口
```
$ docker port <container>
```

主机与容器之间拷贝数据
```shell
# 主机拷贝到容器
$ echo '123'>/tmp/test.txt
$ docker cp /tmp/test.txt nginx:/tmp
$ docker exec -ti nginx cat /tmp/test.txt
123

# 容器拷贝到主机
$ docker cp nginx:/tmp/test.txt ./
```

停止容器
```shell
$ docker stop nginx
```

删除容器

```shell
$ docker rm nginx

# 删除运行中的容器
$ docker rm -f nginx
```

查看容器的明细
```shell
## 查看容器详细信息,包括容器IP地址等
$ docker inspect nginx
```

