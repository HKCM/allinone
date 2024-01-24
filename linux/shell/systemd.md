# systemctl

keyword: 注册服务 系统服务 注册系统服务

## 相关文件
```bash
/etc/systemd/system # 存放用户自定义unit配置文件的目录
/lib/systemd/system # 系统标准unit
/usr/lib/systemd/system # 用户安装软件后自动配置的unit
```

## 命令
```bash
systemctl list-unit-files # 列出现有的unit文件
systemctl daemon-reload # 重新加载
systemctl cat service   # 显示服务配置文件路径和详情
systemctl edit service  # 使用默认编辑器编辑
systemctl --user start my_service.service # 以普通用户启动服务
journalctl -u app.service # 查看app服务的日志
journalctl -n -f -u app.service # 持续查看app服务的日志
```


## 示例


### service

```ini
[Unit]
Description = App Service
After = syslog.target network.target

[Service]
Type=forking
User=admin
Group=admin

WorkingDirectory=/data/admin

ExecStart=/data/admin/app.sh start
ExecStop=/data/admin/app.sh stop
ExecReload=/data/admin/app.sh restart

StandardOutput=journal
StandardError=journal
SyslogIdentifier=app

SuccessExitStatus=143
TimeoutStopSec=10
#Restart=on-failure
#RestartSec=30

[Install]
WantedBy=multi-user.target
```

```bash
WORK_DIR=/data/admin
SERVICE_NAME=app
SERVICE_LINK=/etc/systemd/system/${SERVICE_NAME}.service
SERVICE_FILE=${WORK_DIR}/${SERVICE_NAME}.service
sudo ln -sfn ${SERVICE_FILE} ${SERVICE_LINK}
systemctl reload
systemctl enable app.service

```

### Timer

```ini
[Unit]
Description = App Service
After = syslog.target network.target

[Timer]
OnUnitActiveSec=1min
Unit=app.service


[Install]
WantedBy=multi-user.target
```