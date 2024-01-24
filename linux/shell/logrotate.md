```bash
vim /etc/logrotate.d/log_file_1

/var/log/log_file {
    daily
    rotate 5
    compress
    delaycompress
    missingok
    dateext
    notifempty
    postrotate
        /usr/bin/kill -HUP service_name
    endscript
}

vim /etc/logrotate.d/log_file_1

/var/log/log_file {
    daily
    rotate 5
    compress
    delaycompress
    missingok
    dateext
    notifempty
    copytruncate
}
```

- daily:日志文件将按日轮循。其它可用值为daily，weekly,monthly或者yearly
- rotate: 指定日志文件删除之前转储的次数，0 指没有备份，5 指保留 5 个备份
- dateext: 让 logrotate旧日志文件以创建日期命名
- missingok:在日志轮循期间，任何错误将被忽略，例如 “文件无法找到” 之类的错误。
- size size:当日志文件到达指定的大小时才转储，bytes (缺省) 及 KB (sizek) 或 MB (sizem)
- compress: 通过 gzip 压缩转储以后的日志
- nocompress: 不压缩
- copytruncate: 用于还在打开中的日志文件，先复制一份文件，然后清空原有文件
- nocopytruncate: 备份日志文件但是不截断
- create mode owner group: create 644 root root 以指定的权限创建全新的日志文件
- nocreate: 不建立新的日志文件
- delaycompress: 和 compress 一起使用时，压缩将在下一次轮循周期进行
- nodelaycompress: 覆盖 delaycompress 选项，转储同时压缩。
- errors address: 专储时的错误信息发送到指定的 Email 地址
- ifempty:即使是空文件也转储，这个是 logrotate 的缺省选项。
- notifempty: 如果日志文件为空，轮循不会进行
- mail address: 把转储的日志文件发送到指定的 E-mail 地址
- nomail: 转储时不发送日志文件
- noolddir: 转储后的日志文件和当前日志文件放在同一个目录下
- prerotate/endscript: 在所有其它指令完成后，postrotate 和 endscript 里面指定的命令将被执行,这两个关键字必须单独成行

修改完成后
```bash
systemctl restart logrotate.service
```

```
logrotate [OPTION...] <configfile>
-d, --debug:debug 模式，测试配置文件是否有错误。
-f, --force:强制转储文件。
-m, --mail=command:压缩日志后，发送日志到指定邮箱。
-s, --state=statefile:使用指定的状态文件。
-v, --verbose:显示转储过程。
```

手动执行
```bash
logrotate /etc/logrotate.d/log_file
logrotate -vdf /etc/logrotate.d/log_file # debug 不实际执行
```

可以配合crontab
```bash
crontab -e
*/30 * * * * /usr/sbin/logrotate /etc/logrotate.d/rsyslog > /dev/null 2>&1 &
```
