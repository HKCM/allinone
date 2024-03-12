# zip

-r 将指定目录下的所有文件和子目录一并压缩
-x 压缩时排除某个文件

```bash
zip -r test.zip ./service
mv service services
unzip test.zip

# 排除压缩
zip -r tmp1.zip ./tmp/ -x tmp/services.zip # 选项指定不压缩的文件。
    adding: tmp/ (stored 0%)
    adding: tmp/services (deflated 80%)
    adding: tmp/.ICE-unix/ (stored 0%)
```

# unzip

-d 指定解压目录
-l 列出压缩内容

```bash
unzip -l test.zip # 不解压,显示压缩包内容

unzip -d ./new_testdir test.zip 
```