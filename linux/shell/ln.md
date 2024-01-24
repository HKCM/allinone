# ln

软链接

`ln -s origin_file link_file` 两个文件的内容并不相同，是两个完全不同的文件

```shell
# i替换确认
ln -si data2.txt data.txt # 文件创建软链
ln -s /var/www/ /web # 目录创建软链

# 使用readlink打印出符号链接所指向的目标路径：
readlink web
```

硬链接的文件共享inode编号,本质是同一个文件

```shell
ln code_file hl_code_file
ls -li *code_file # 确认inode编号
```
