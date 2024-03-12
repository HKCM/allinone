# ls

```shell
ls -ld # 查看目录
ls -lF # -F 文件夹后加/,可执行文件后面加*号
ls -lSh ./ # 文件按从大到小显示
ls -lShr ./ # 文件按从小到大显示
ls -lth ./ # 文件按最近修改时间排序
ls -i # 显示inode编号
ls -l | grep "^d"
ls -l | grep "^l" # 打印出当前目录下的符号链接
```

```bash
ls -l --time-style=long-iso   #以long-iso方式显示时间
total 12
drwxr-xr-x 2 root root 4096 2015-10-25 11:13 dir1
-rw-r--r-- 1 root root    0 2015-10-25 11:13 file1.txt

# 显示到秒
ls -l --time-style=full-iso
total 4
drwxr-xr-x 2 root root 4096 2024-03-09 13:54:58.635881012 +0000 dir1
-rw-r--r-- 1 root root    0 2024-03-09 13:55:07.928016002 +0000 file1.txt

# 创建时间
root@edd870dcacc8:/testdir# ls -l --time=birth file1.txt 
-rw-r--r-- 1 root root 4 Mar  9 13:55 file1.txt
# 修改时间
root@edd870dcacc8:/testdir# ls -l --time=ctime file1.txt 
-rw-r--r-- 1 root root 4 Mar  9 13:59 file1.txt
# 访问时间
root@edd870dcacc8:/testdir# ls -l --time=atime file1.txt 
-rw-r--r-- 1 root root 4 Mar  9 13:57 file1.txt
```
