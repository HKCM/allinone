# crontab

`/var/spool/cron/` 目录下存放的是每个用户包括root的crontab任务，每个任务以创建者的名字命名

`/etc/crontab` 这个文件负责调度各种管理和维护任务

`/etc/cron.d/` 这个目录用来存放任何要执行的crontab文件或脚本

如果创建的脚本对精确的执行时间要求不高,用预配置的cron脚本目录会更方便。有4个基本目录:hourly、daily、monthly和weekly。

```bash
$ ls /etc/cron.*ly
```

可以把脚本放在`/etc/cron.hourly`、`/etc/cron.daily`、`/etc/cron.weekly`、`/etc/cron.monthly`目录中，让它每小时/天/星期、月执行一次

```
crontab [-u username]　　　　//省略用户表表示操作当前用户的crontab
    -e      (编辑工作表)
    -l      (列出工作表里的命令)
    -r      (删除工作作)
```

## 注意事项

1. 尽量使用脚本的绝对路径
2. 不要依赖环境变量,或者直接将需要的变量写在脚本中
3. 在crontab中写上注释,让人知道每个job的作用
4. 注意系统时间 `date +"%F %T"`

```bash
# 以非交互式方式写入和删除crontab
(crontab -l;echo "0 2 * * * sh /tmp/t1.sh > /tmp/t1.log 2>&1") | crontab
( crontab -l | grep -v"$cron_job" ) | crontab 
```


## 示例
```bash
* * * * * myCommand                      # 示例1：每1分钟执行一次myCommand
3,15 * * * * myCommand                   # 示例2：每小时的第3和第15分钟执行
3,15 8-11 * * * myCommand                # 示例3：在上午8点到11点的第3和第15分钟执行
3,15 8-11 */2  *  * myCommand            # 示例4：每隔两天的上午8点到11点的第3和第15分钟执行
3,15 8-11 * * 1 myCommand                # 示例5：每周一上午8点到11点的第3和第15分钟执行
30 21 * * * /etc/init.d/smb restart      # 示例6：每晚的21:30重启smb
45 4 1,10,22 * * /etc/init.d/smb restart # 示例7：每月1、10、22日的4 : 45重启smb
10 1 * * 6,0 /etc/init.d/smb restart     # 示例8：每周六、周日的1 : 10重启smb
0,30 18-23 * * * /etc/init.d/smb restart # 示例9：每天18 : 00至23 : 00之间每隔30分钟重启smb
0 23 * * 6 /etc/init.d/smb restart       # 示例10：每星期六的晚上11 : 00 pm重启smb
0 */1 * * * /etc/init.d/smb restart      # 示例11：每一小时重启smb
0 23-7/1 * * * /etc/init.d/smb restart   # 示例12：晚上11点到早上7点之间，每隔一小时重启smb
00 12 * * * if [ `date +%d -d tomorrow` = 01 ] ; then ; command # 每月最后一天中午12点执行命令
```
