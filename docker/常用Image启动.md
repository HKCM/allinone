# Nginx

```bash
docker run -p 8080:80 -d nginx

docker run --name some-nginx -v /some/content:/usr/share/nginx/html:ro -d nginx
```

```Dockerfile
FROM nginx
COPY nginx.conf /etc/nginx/nginx.conf
```

Then build the image with `docker build -t custom-nginx .` and run it as follows:
```bash
docker run --name my-custom-nginx-container -d custom-nginx
```

# Mysql

```shell
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql

docker run -p 3306:3306 --name node-mysql -e MYSQL_ROOT_PASSWORD=mysql  -d mysql:5.7
```

# Gitlab

https://docs.gitlab.com/ce/install/docker.html

```bash
GITLAB_HOME=$HOME/gitlab
# GITLAB_HOME=$HOME/gitlab >> ~/.zshrc
docker run --detach \
  --hostname gitlab.example.com \
  --publish 443:443 --publish 80:80 --publish 22:22 \
  --name gitlab \
  --restart always \
  --volume $GITLAB_HOME/config:/etc/gitlab \
  --volume $GITLAB_HOME/logs:/var/log/gitlab \
  --volume $GITLAB_HOME/data:/var/opt/gitlab \
  --shm-size 256m \
  gitlab/gitlab-ce:latest

# 获取root密码
docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password

```