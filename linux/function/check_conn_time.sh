#!/usr/bin/env bash

# -w:从文件中读取信息打印格式
# -o:输出的全部信息
# -s:不打印进度条
cat>curl_format.txt <<EOF
time_namelookup: %{time_namelookup}\t\t#Time from start until name resolving completed\n
time_connect: %{time_connect}\t\t#Time from start until remote host or proxy completed.\n
time_appconnect: %{time_appconnect}\t\t#Time from start until SSL/SSH handshake completed.\n
time_pretransfer: %{time_pretransfer}\t\t#Time from start until just before the transfer begins\n
time_redirect: %{time_redirect}\t\t#Time taken for all redirect steps before the final transfer.\n
time_starttransfer: %{time_starttransfer}\t\t#Time from start until just when the first byte is received.(Time to first byte,TTFB)\n
---\n
time_total: %{time_total}\n
EOF

# 模版文件二选一
# cat>curl_format.txt <<EOF
#      time_namelookup:  %{time_namelookup}s\n
#         time_connect:  %{time_connect}s\n
#      time_appconnect:  %{time_appconnect}s\n
#     time_pretransfer:  %{time_pretransfer}s\n
#        time_redirect:  %{time_redirect}s\n
#   time_starttransfer:  %{time_starttransfer}s\n
#                      ----------\n
#           time_total:  %{time_total}s\n
# EOF
# - time_namelookup:    DNS 域名解析的时候
# - time_connect:       TCP 连接建立的时间，就是三次握手的时间
# - time_appconnect:    SSL/SSH等上层协议建立连接的时间，比如 connect/handshake 的时间
# - time_pretransfer:   从请求开始到响应开始传输的时间
# - time_starttransfer: 从请求开始到第一个字节将要传输的时间
# - time_total:         这次请求花费的全部时间


curl -o /dev/null -w "@curl_format.txt" -s https://www.baidu.com
