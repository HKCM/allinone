# awk

awk常用功能:
1. 指定分隔符显示某几列
2. 通过正则表达式取出想要的内容
3. 显示出某个范围内的内容通awk进行统计算
4. awk组计算与去重



## 示例

```bash
cat > test.txt <<EOF
root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
adm:x:3:4:adm:/var/adm:/sbin/nologin
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
halt:x:7:0:halt:/sbin:/sbin/halt
mail:x:8:12:mail:/var/spool/mail:/sbin/nologin
uucp:x:10:14:uucp:/var/spool/uucp:/sbin/nologin
EOF
```

```bash
# 打印第五行, 需要用两个等号表示判断
awk "NR==5" test.txt 
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
# 打印第五行到七行
awk "NR==5,NR==7" test.txt 
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
sync:x:5:0:sync:/sbin:/bin/sync

# 输出加上行号, 必须使用单引号
awk '{print NR,$0}' test.txt 
1 root:x:0:0:root:/root:/bin/bash
2 bin:x:1:1:bin:/bin:/sbin/nologin
3 daemon:x:2:2:daemon:/sbin:/sbin/nologin
4 adm:x:3:4:adm:/var/adm:/sbin/nologin
5 lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
6 sync:x:5:0:sync:/sbin:/bin/sync
7 shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
8 halt:x:7:0:halt:/sbin:/sbin/halt
9 mail:x:8:12:mail:/var/spool/mail:/sbin/nologin
10 uucp:x:10:14:uucp:/var/spool/uucp:/sbin/nologin
# 打印第五行到七行 并添加行号
awk 'NR==5,NR==7 {print NR,$0}' test.txt 
```

```shell
awk 'BEGIN{ print "start" } pattern { commands } END{ print "end" }' file
```

工作原理:
1. 首先执行BEGIN { commands } 语句块中的语句。
2. 接着从文件或stdin中读取一行，如果能够匹配pattern，则执行随后的commands语句块。重复这个过程，直到文件全部被读取完毕。
3. 当读至输入流末尾时，执行END { commands } 语句块。

```shell
$ echo -e "line1\nline2" | awk 'BEGIN { print "Start" } { print } END { print "End" } '
Start
line1
line2
End
```

特殊变量:
- NR：表示记录编号，当awk将行作为记录时，该变量相当于当前行号。
- NF：表示字段数量，在处理当前记录时，相当于字段数量。默认的字段分隔符是空格。
- $0：该变量包含当前记录的文本内容。
- $1：该变量包含第一个字段的文本内容。
- $2：该变量包含第二个字段的文本内容。

```shell
$ echo -e "line1 f2 f3\nline2 f4 f5\nline3 f6 f7" | \
awk '{
 print "Line no:"NR",No of fields:"NF, "$0="$0,
 "$1="$1,"$2="$2,"$3="$3
}'
Line no:1,No of fields:3 $0=line1 f2 f3 $1=line1 $2=f2 $3=f3
Line no:2,No of fields:3 $0=line2 f4 f5 $1=line2 $2=f4 $3=f5
Line no:3,No of fields:3 $0=line3 f6 f7 $1=line3 $2=f6 $3=f7 
```

简单示例
```shell
$ cat data2.txt
One line of test text.
Two lines of test text.
Three lines of test text.

$ awk '{print $1}' data2.txt
One
Two
Three

# -F指定分隔符
$ awk -F: '{print $1}' /etc/passwd
root
daemon
bin
sys
sync
games
...

$ echo "My name is Rich" | awk '{$4="Christine"; print $0}' 
My name is Christine
```

跟sed编辑器一样,gawk编辑器允许将程序存储到文件中,然后再在命令行中引用
```shell
$ cat script.awk
{print $1 " home directory is " $6}

$ awk -F: -f script.awk /etc/passwd
root home directory is /root
daemon home directory is /usr/sbin
bin home directory is /bin
sys home directory is /dev
sync home directory is /bin
games home directory is /usr/games
man home directory is /var/cache/man

# 写作还可以多行,这里还使用了变量
$ cat script3.awk
{
text = "'s home directory is " 
print $1 text $6
}

```

BEGIN和END
```shell
$ awk 'BEGIN {print "Hello"};{print $0};END {print "BYE"}' data2.txt 
Hello
One line of test text.
Two lines of test text.
Three lines of test text.
BYE


$ cat data1
data11,data12,data13,data14,data15
data21,data22,data23,data24,data25
data31,data32,data33,data34,data35
# 以逗号为分隔符分隔原数据,将"-"号作为输出分隔符,只输出$1 $2 $3
$ awk 'BEGIN {FS=",";OFS="-"} {print $1,$2,$3}' data1
data11-data12-data13
data21-data22-data23
data31-data32-data33
```

#### 
```shell
$ cat data1b
1005.3247596.37
115-2.349194.00
05810.1298100.1
$ awk 'BEGIN{FIELDWIDTHS="3 5 2 5"}{print $1,$2,$3,$4}' data1b 
100 5.324 75 96.37
115 -2.34 91 94.00 
058 10.12 98 100.1
```
