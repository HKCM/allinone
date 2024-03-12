# date

```bash
date +%Y # 2023
date +%F # 2023-11-23
date +%T # 15:51:44
date +%s # 打印纪元时
date +"%Y-%m-%d %H:%M:%S" # 2023-11-23 22:20:02
date --date "Jan 20 2001" +%A # Saturday 获取星期

date +%F -d "-1day"       #<==显示昨天（简洁写法）
date +%F -d "yesterday"   #<==显示昨天（英文写法）
date +%F -d "-2day"       #<==显示前天
date +%F -d "+1day"       #<==显示明天
date +%F -d "tomorrow"    #<==显示明天（英文写法）
date +%F -d "+2day"       #<==显示2天后
date +%F -d "1month"      #<==显示1个月后
date +%F -d "1year"       #<==显示1年后
# 返回之后一小时 在Mac上不好使
date -d '+1 hour' +"%Y-%m-%d %H:%M:%S"
date -d "+15 minutes" "+%Y-%m-%d %H:%M:%S"

# 时间时区
# /usr/share/zoneinfo/ 目录下
TZ="America/New_York" date +"%Y-%m-%d %H:%M:%S %z"
TZ="Japan" date +"%Y-%m-%d %H:%M:%S %z" # 2023-11-23 22:19:41 +0900

# 时间转换
date -d "Thu Jul  6 21:41:16 CST 2017" "+%Y-%m-%d %H:%M:%S"
```