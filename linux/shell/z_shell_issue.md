# Issue

## 避免’sudo echo x >’ 时’Permission denied’

示例:

```shell
sudo echo a > 1.txt
-bash: 1.txt: Permission denied
```

bash 拒绝这么做,说是权限不够.
这是因为重定向符号 “>” 也是 bash 的命令。sudo 只是让 echo 命令具有了 root 权限,
但是没有让 “>” 命令也具有root 权限,所以 bash 会认为这个命令没有写入信息的权限。

利用管道和 tee 命令,该命令可以从标准输入中读入信息并将其写入标准输出或文件中,
具体用法如下:

```
echo a |sudo tee 1.txt
echo a |sudo tee -a 1.txt  // -a 是追加的意思,等同于 >>
```

tee 命令很好用,它从管道接受信息,一边向屏幕输出,一边向文件写入。