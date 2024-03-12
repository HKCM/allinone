# grep

-V 显示不匹配的行，或者说排除某些行，显示不包含匹配文本的所有行
-n 显示配行及行号
-i 不区分大小写（只适用于单字符），默认是区分大小写的只统计匹配的行数，注意不是匹配的次数
-E 使用扩的egrep命令
-W 只匹配过滤的单词
-0 只输出匹配的内容
--color=auto 为grep过滤的匹配字符添加颜色

```shell
grep pattern filename # 在文件中搜索匹配特定模式的文本行
grep "match_text" file1 file2 file3 # 多个文件中搜索匹配特定模式的文本行
grep -E "[a-z]+" filename # 使用正则表达式
grep -E '[A-Za-z0-9._]+@[A-Za-z0-9.]+\.[a-zA-Z]{2,4}' # 使用正则表达式匹配邮箱
grep -v match_pattern filename # 反向匹配
grep -c "text" filename # 统计出匹配模式的文本行数
grep -e "pattern1" -e "pattern2" filename # 多模式匹配
grep -rn "text" . # 递归查找文本并显示行号
grep "text" -B 3 filename # 打印匹配结果之前的行
grep "text" -A 3 filename # 打印匹配结果之后的行
grep "text" -C 3 filename # 打印匹配结果之前之后的行
```

关键词: 检索数据 过滤数据 检索日志 定位日志

```shell
$ grep -n '2019-10-24 00:01:11' *.log # 检索指定日期的日志
$ grep -rn --color=auto "keyword" ./folder # 递归检索
$ grep -rn --color=auto "^https.*77$" ./folder # 检索以https开头,以77结尾的行
$ grep -rn -v --color=auto "^https" ./folder # 反向检索
$ grep -c "keyword" ./folder # 只想知道有多少匹配的行
$ grep -e "keyword1" -e "keyword2" file1 # 指定多个匹配模式,可用-e参数来指定每个模式。
$ grep -B 5 "keyword" ./folder # 检索关键词前五行
$ grep -A 5 "keyword" ./folder # 检索关键词后五行
$ grep -C 5  "keyword" ./folder # 检索关键词前后五行
```

```bash
# 去除注释和空行
grep -Ev "^$|#" nginx.conf
```