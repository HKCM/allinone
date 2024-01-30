
```bash
docker run -d -p 80:80 docker/getting-started

```

创建网络。

```bash
docker network create todo-app
```

启动一个MySQL容器并将其连接到网络
```shell
docker run -d \
    --network todo-app --network-alias mysql \
    -v todo-mysql-data:/var/lib/mysql \
    -e MYSQL_ROOT_PASSWORD=secret \
    -e MYSQL_DATABASE=todos \
    mysql:5.7
```

使用nicolaka / netshoot映像启动一个新容器。确保将其连接到同一网络

```shell
docker run -it --network todo-app nicolaka/netshoot

# dig mysql
```

指定上面的每个环境变量,并将容器连接到我们的应用程序网络。
```shell
# MYSQL_HOST -正在运行的MySQL服务器的主机名
# MYSQL_USER -用于连接的用户名
# MYSQL_PASSWORD -用于连接的密码
# MYSQL_DB -连接后要使用的数据库
docker run -dp 3000:3000 \
  -w /app -v ${PWD}:/app \
  --network todo-app \
  -e MYSQL_HOST=mysql \
  -e MYSQL_USER=root \
  -e MYSQL_PASSWORD=secret \
  -e MYSQL_DB=todos \
  node:12-alpine \
  sh -c "yarn install && yarn run dev"
```

查看容器(docker logs <container-id>)的日志,则应该看到一条消息,表明它正在使用mysql数据库。
```shell
# Previous log messages omitted
nodemon src/index.js
[nodemon] 1.19.2
[nodemon] to restart at any time, enter `rs`
[nodemon] watching dir(s): *.*
[nodemon] starting `node src/index.js`
Connected to mysql db at host mysql
Listening on port 3000
```
