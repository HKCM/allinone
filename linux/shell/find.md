# find

## 按文件名查找
```shell
find /home/slynux -name '*.txt' -print # 查找以txt结尾的文件
find . ! -name "*.txt" -print # 查找不以txt结尾的文件
find . \( -name '*.txt' -o -name '*.pdf' \) -print # 查找以.txt和.pdf结尾文件
find . \( -name '*e*' -and -name 's*' \) # 查找名字以s开头且其中包含e的文件
# find . -type f -name "*.swp" -delete # 删除匹配的文件,危险命令
```

## 按文件类型查找
```shell
find . -type d -print # 只列出所有的目录（包括子目录）
find . -type f -print # 只列出普通文件
find . -type l -print # 只列出符号链接
```

## 按时间先后查找

`-`表示小于，`+`表示大于

- amin（访问时间）
- mmin（修改时间）
- cmin（变化时间）

```shell
find . -type f -atime -7 -print # 最近7天内被访问过的所有文件
find . -type f -atime 7 -print # 恰好在7天前被访问过的所有文件
find . -type f -atime +7 -print # 访问时间超过7天的所有文件
find . -type f -amin -7 -print # 7分钟之内访问的所有文件
```

## 按文件大小查找

- w：字（2字节）。
- k：千字节（1024字节）。
- M：兆字节（1024K字节）。
- G：吉字节（1024M字节）。

```shell
find . -type f -size +2k # 大于2KB的文件
find . -type f -size -2k # 小于2KB的文件
find . -type f -size 2k # 大小等于2KB的文件
```

## 按文件权限和所有权查找

```shell
find . -type f -perm 644 -print # 权限为644的文件
find . -type f -user slynux -print # 用户slynux拥有的所有文件
```

## 条件查找

```shell
find . -name "*.sh" -a -user ubuntu # 与,查找以.sh结尾并且属主为ubuntu的文件
find . -name "*.sh" -o -user ubuntu # 或,查找以.sh结尾或者属主为ubuntu的文件
find . -name "*.sh" -not -user ubuntu # 非,查找以.sh结尾并且属主不是ubuntu的文件
```

#### 查找后执行命令

find命令使用一对花括号{}代表找到的文件名

结尾必须对分号进行转义`\;`，否则shell会将其视为find命令的结束，而非chown命令的结束。

- print find默动作
- ok [commend]  查找后执行命令的时候询问用户是否要执行
- exec [commend] 查找后执行命令的时候不询问用户，直接执行

```shell
find . -name "*.sh" -ok rm {} \; # 会询问每一个找到的文件
find . -type f -user root -exec chown slynux {} \; # 直接将root用户的文件改成slynux的
find . -type f -mtime +10 -name "*.txt" -exec cp {} OLD \; # 将10天前的 .txt文件复制到OLD目录中
find . -type f -name "*.mp3" -exec mv {} target_dir \; # 所有的.mp3文件移入给定的目录

# 递归的方式将所有文件名中的空格替换为字符"_"
# rename命令需要安装
find . -type f -exec rename 's/ /_/g' {} \;
```