# journalctl

```bash
journalctl -u ssh # 查看 SSH 单元
journalctl -fu ssh # 实时查看 SSH 单元

journalctl -u ssh --since yesterday # 查看昨天的日志
journalctl -u ssh --since -3d --until -2d # 查看三天前的日志
journalctl -u ssh --since -1h # 查看上个小时的日志
journalctl -u ssh --until "2022-03-12 07:00:00" # 查看截至到某个时间点的日志
```