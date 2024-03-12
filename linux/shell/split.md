# split

-b: 按字节数分割文件
-l: 按行数分割文件
-a: 指定后缀长度,默认为2
-d: 使用数字后缀

```bash
cat test.txt 
1
2
3
4
5
6
7
8
9
10

split -b 20 test.txt new_ # 每20字节分割一次文件, 可以使用500k等单位
ll
-rw-r--r-- 1 root root   20 Mar 10 12:44 new_aa
-rw-r--r-- 1 root root    1 Mar 10 12:44 new_ab
-rw-r--r-- 1 root root   21 Mar 10 12:37 test.txt
cat new_ac
9
10
split -l 4 -d test.txt new_test.txt  # 每4行分割一次文件,使用数字后缀
-rw-r--r-- 1 root root    8 Mar 10 12:42 new_test.txt00
-rw-r--r-- 1 root root    8 Mar 10 12:42 new_test.txt01
-rw-r--r-- 1 root root    5 Mar 10 12:42 new_test.txt02
-rw-r--r-- 1 root root   21 Mar 10 12:37 test.txt
```