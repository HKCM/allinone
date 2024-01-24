# sort and uniq

keyword: 排序 去重 

sort命令对于字母表排序和数字排序有不同的处理方式

```shell
cat file1
1
3
7
9
5

sort -n file1 # 按数字进行排序
sort -r file1 # 逆序
```

```shell
cat > data.txt <<EOF
1,mac,2000
2,winxp,4000
3,bsd,1000
4,linux,1000
EOF

# t指定分隔符
# k指定排序的字段
# r逆序
# n依据数值大小排序
sort -t ',' -k 3 -n data.txt
sort -t ',' -k 2 data.txt
sort -t ',' -k 3 -r data.txt 
sort -t ',' -k 2.2 data.txt # 依据第2列第二个字符进行排序
```

uniq只能作用于排过序的数据，因此，uniq通常都与sort命令结合使用

```shell
cat > sorted.txt <<EOF
bash
foss
hack
hack
EOF

# 统计各行在文件中出现的次数
sort unsorted.txt | uniq -c 
1 bash
1 foss
2 hack 

# 只显示重复的行
sort unsorted.txt | uniq -d
hack
```