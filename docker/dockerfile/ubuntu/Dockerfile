# 使用基础镜像
FROM ubuntu:latest

RUN echo "#!/bin/sh\nexit 0" > /usr/sbin/policy-rc.d

# 安装 vim
RUN apt-get update && \
    apt-get install -y vim git man-db && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN yes | unminimize

# 设置工作目录
WORKDIR /app

# # 添加后台运行脚本
# COPY start_background.sh /usr/local/bin/

# # 设置脚本权限
# RUN chmod +x /usr/local/bin/start_background.sh

# # 在后台运行脚本
# CMD ["/usr/local/bin/start_background.sh"]

# apt update && apt install cron
# crontab -e
