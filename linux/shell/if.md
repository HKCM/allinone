# if

```bash
# if commands; then commands;else commands; fi
if [ -a a.txt ];then echo 123;else echo 456;fi
```


## 文件目录判定

```
[ -a FILE ] 如果 FILE 存在则为真。
[ -d FILE ] 如果 FILE 存在且是一个目录则返回为真。
[ -e FILE ] 如果 指定的文件或目录存在时返回为真。
[ -f FILE ] 如果 FILE 存在且是一个普通文件则返回为真。
[ -r FILE ] 如果 FILE 存在且是可读的则返回为真。
[ -w FILE ] 如果 FILE 存在且是可写的则返回为真。（一个目录为了它的内容被访问必然是可执行的）
[ -x FILE ] 如果 FILE 存在且是可执行的则返回为真。
[ -s FILE ] 如果 FILE 文件大小非0时为真
[ -L FILE ] 如果 FILE 为符号链接,则为真
```

不常用的
```
[ -b FILE ] 如果 FILE 存在且是一个块文件则返回为真。
[ -c FILE ] 如果 FILE 存在且是一个字符文件则返回为真。
[ -g FILE ] 如果 FILE 存在且设置了SGID则返回为真。
[ -h FILE ] 如果 FILE 存在且是一个符号符号链接文件则返回为真。（该选项在一些老系统上无效）
[ -k FILE ] 如果 FILE 存在且已经设置了冒险位则返回为真。
[ -p FILE ] 如果 FILE 存并且是命令管道时返回为真。
[ -s FILE ] 如果 FILE 存在且大小非0时为真则返回为真。
[ -u FILE ] 如果 FILE 存在且设置了SUID位时返回为真。
[ -O FILE ] 如果 FILE 存在且属有效用户ID则返回为真。
[ -G FILE ] 如果 FILE 存在且默认组为当前组则返回为真。（只检查系统默认组）
[ -L FILE ] 如果 FILE 存在且是一个符号连接则返回为真。
[ -N FILE ] 如果 FILE 存在 and has been mod如果ied since it was last read则返回为真。
[ -S FILE ] 如果 FILE 存在且是一个套接字则返回为真。
[ FILE1 -nt FILE2 ] 如果 FILE1 比 FILE2 新, 或者 FILE1 存在但是 FILE2 不存在则返回为真。
[ FILE1 -ot FILE2 ] 如果 FILE1 比 FILE2 老, 或者 FILE2 存在但是 FILE1 不存在则返回为真。
[ FILE1 -ef FILE2 ] 如果 FILE1 和 FILE2 指向相同的设备和节点号则返回为真。
```

## 字符串比较运算符

进行字符串比较时，最好用双中括号，因为有时候采用单个中括号会产生错误。注意使用引号,防止空格扰乱代码

```
[[ -z "${STRING}" ]] 如果 string长度为零，则为真
[[ -n "${STRING}" ]] 如果 string长度非零，则为真
[[ "${STRING1}" ]] 如果字符串不为空则返回为真,与-n类似 
[[ "${STRING1}" = "${STRING2}" ]] 如果两个字符串相同则返回为真
[[ "${STRING1}" == "${STRING2}" ]] 如果两个字符串相同则返回为真
[[ "${STRING1}" != "${STRING2}" ]] 如果字符串不相同则返回为真
[[ "${STRING1}" =~ "${STRING2}" ]] 如果 STRING2 是 STRING1的一部分，则为真
```

## 算术比较运算符

```
[ $num1 -eq $num2 ] 等于 [ 3 -eq $mynum ]
[ $num1 -ne $num2 ] 不等于 [ 3 -ne $mynum ]
[ $num1 -lt $num2 ] 小于 [ 3 -lt $mynum ]
[ $num1 -le $num2 ] 小于或等于 [ 3 -le $mynum ]
[ $num1 -gt $num2 ] 大于 [ 3 -gt $mynum ]
[ $num1 -ge $num2 ] 大于或等于 [ 3 -ge $mynum ]
```


## 示例

1:判断目录`$dir`是否存在,若不存在,则新建一个

```shell
if [ ! -d "$dir"]; then
  mkdir "$dir"
fi
```

2:判断普通文件`$file`是否存,若不存在,则新建一个
```shell
if [ ! -f "$file" ]; then
  touch "$file"
fi
```

3:判断`$sh`是否存在并且是否具有可执行权限
```shell
if [ ! -x "$sh"]; then
    chmod +x "$sh"
fi
```

4:是判断变量`$var`是否有值
```shell
if [ ! -n "$var" ]; then
　　echo "$var is empty"
　　exit 0
fi
```

判断数字
```shell
if [[ ! "${TARGET_TIME}" =~ ^[1-9][0-9]*$ ]]; then
  echo
  echo "Invaild input..."
  echo
  display_help
  exit 1
fi
```