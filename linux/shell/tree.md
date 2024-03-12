# tree

-a 显示所有文件，包括隐藏文件（以“”点头的文件）
-d 只显示目录
-f 显示每个文件的全路径
-i 不显示树枝，常与-f参数配合使用
-L level 遍历的最层数，level于0的
--noreport 不显示最后一行的统计信息

```bash
tree -fid oldboy
oldboy
oldboy/dir1_1
oldboy/dir1_1/dir2_1
oldboy/dir1_1/dir2_2
oldboy/dir1_2
oldboy/dir1_2/dir2_1
oldboy/dir1_2/dir2_2
oldboy/test

7 directories
```

拷贝目录结构
```bash
# 将目录树内容追加到家目录下的文件里
tree -fid --noreport src_dir >> ~/src_dir.txt 

mkdir -p $(cat ~/src_dir.txt )
```