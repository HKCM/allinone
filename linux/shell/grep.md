# grep

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