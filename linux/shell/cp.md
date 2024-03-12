# cp

```bash
cp -i src des # 如果复制的目标位置已经存在同名的文件，则会提示是否覆盖
cp -p src des # 使用-p保留源文件属性,包括所有者,权限和时间戳
cp -r src_dir des_dir # 使用 -r选项复制目录

# 如果复制的目标是一个软链接,直接cp会复制原文件而不是软连接
cp -d src_slink des_slink # 使用-d复制软连接

cp -a src_dir des_dir # -a选项相当于-dpr

cp -s src_slink des_slink # -s创建软连接,功能与ln -s一样
```

快速备份
```bash
cp /etc/ssh/sshd_config /etc/ssh/sshd_config.bak

cp /etc/ssh/sshd_config{,.bak}
```

拷贝目录结构
```bash
# 将目录树内容追加到家目录下的文件里
tree -fid --noreport src_dir >> ~/src_dir.txt 

mkdir -p $(cat ~/src_dir.txt )
```