
# docker技巧

#### 查找docker container的日志文件位置
```shell
docker inspect --format='{{.LogPath}}' containerid
```

### 数据拷贝

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

#### 删除带none的镜像,提示先停止容器。
```shell
docker stop $(docker ps -a | grep "Exited" | awk '{print $1 }') //停止容器

docker rm $(docker ps -a | grep "Exited" | awk '{print $1 }') //删除容器

docker rmi $(docker images | grep "none" | awk '{print $3}') //删除镜像
```

#### 删除所有退出的容器
```shell
for i in `docker ps -a | grep -i exit | awk '{print $1}'`; do docker rm -f $i; done
```

#### 清理主机上所有退出的容器
```shell
$ docker rm  $(docker ps -aq)
```

#### 调试或者排查容器启动错误
```shell
# 若有时遇到容器启动失败的情况,可以先使用相同的镜像启动一个临时容器,先进入容器
$ docker exec -ti --rm <image_id> bash
# 进入容器后,手动执行该容器对应的ENTRYPOINT或者CMD命令,这样即使出错,容器也不会退出,因为bash作为1号进程,我们只要不退出容器,该容器就不会自动退出
```

查找docker日志文件位置
```shell
docker inspect --format='{{.LogPath}}'
```

## 更改docker的默认存储目录



方式一: 修改docker配置文件
```bash
sudo vim /etc/docker/daemon.json
{
    "data-root": "/data/var/lib/docker"
}

sudo systemctl daemon-reload
sudo systemctl restart docker.service
```

方式二: 修改docker的启动文件
```bash
sudo vim /lib/systemd/system/docker.service
```

```ini
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target docker.socket firewalld.service containerd.service time-set.target
Wants=network-online.target containerd.service
Requires=docker.socket

[Service]
Type=notify
ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock --data-root=/data/docker #指定数据目录
ExecReload=/bin/kill -s HUP $MAINPID
TimeoutStartSec=0
RestartSec=2
Restart=always
StartLimitBurst=3
StartLimitInterval=60s
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
Delegate=yes
KillMode=process
OOMScoreAdjust=-500

[Install]
WantedBy=multi-user.target
```