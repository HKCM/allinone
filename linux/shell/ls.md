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