# shutdown

```bash
shutdown -r now # 立即重启
shutdown -r 15:30 & # 放入后台定时重启
shutdown -r +10 # 10分钟后重启
shutdown -c # 取消重启

shutdown -h now #立即关机
shutdown -h 15:30 & # 放入后台定时关机
```

其他命令

```bash
halt # 关机
poweroff # 关机
init 0 # 关机
init 6 # 重启
```