```bash
date +%F # 2023-11-23
date +%T # 15:51:44
date +%s # 打印纪元时
date +"%Y-%m-%d %H:%M:%S" # 2023-11-23 22:20:02
date --date "Jan 20 2001" +%A # Saturday 获取星期几
# 时间时区
# /usr/share/zoneinfo/ 目录下
TZ="America/New_York" date +"%Y-%m-%d %H:%M:%S %z"
TZ="Japan" date +"%Y-%m-%d %H:%M:%S %z" # 2023-11-23 22:19:41 +0900
# 返回之后一小时 在Mac上不好使
date -d '+1 hour' +"%Y-%m-%d %H:%M:%S"
date -d "+15 minutes" "+%Y-%m-%d %H:%M:%S"
```