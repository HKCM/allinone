# paste

-d 指定分隔符
-s 每个文件占一行

```bash
cat test1
1
2
3
4
5
cat test2
a
b

c
paste -d"," test1 test2
1,a
2,b
3,
4,c
5,
paste -s -d"," test1
1,2,3,4,5
```

## 示例

将以下文件设置为 `stu10309=7f753cc3`的形式
```bash
cat test.txt
stu10309  #<==账号。
7f753cc3  #<==密码。
stu10312
636e026d
stu10315
18273b95
```

方式一: 轮流使用等号与换行符作为分隔符
```bash
paste -s -d"=\n" test.txt 
stu10309=7f753cc3
stu10312=636e026d
stu10315=18273b95
```

方式二: 使用等号作为分隔符,使用`- -`一次读入两行
```bash
paste -d "=" - - < test.txt 
stu10309=7f753cc3
stu10312=636e026d
stu10315=18273b95
```

方式三: 使用xargs+sed实现
```bash
xargs -n 2 < test.txt | sed "s/ /=/g"
stu10309=7f753cc3
stu10312=636e026d
stu10315=18273b95
```