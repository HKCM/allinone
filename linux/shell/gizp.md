# gzip

这里的压缩和Windows的压缩是不一样的,压缩是压缩单个文件,所以经常和tar打包命令连用

-d: 解压指定的压缩包文件
-k: 保留原文件
-q: 静默执行模式

```bash
gzip data.txt # 将指定的文件进行压缩，默认将删除原文件
gzip -k data.txt # 将指定的文件进行压缩，但是不删除原文件
gzip -l data1.txt.gz # 显示压缩包内的文件信息
compressed        uncompressed  ratio uncompressed_name
    107                 133     40.6%    data1.txt
gzip -d data1.txt.gz # 解压
gunzip data1.txt.gz # 解压
```