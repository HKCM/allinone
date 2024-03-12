# cut

```shell
cat test.txt
root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin

# 显示第1列和第3列
cut -d":" -f 1,3 test.txt 
root:0
daemon:1
bin:2

# 以;为分隔符获取第五列
cut -f5 -d":" test.txt 
root
daemon
bin

# 打印第2个到第5个字符
cut -c2-5 test.txt 
oot:
aemo
in:x
```