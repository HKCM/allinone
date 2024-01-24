# rsync

```bash
rsync -r source destination # destination/source的目录结构
rsync -a source destination # 同步元信息(比如修改时间、权限等)
rsync -a source1 source2 destination # 同步多个目录
rsync -anv source1 source2 destination # 参数模拟执行的结果
rsync -avz ./linux test@10.211.55.5:/home/test/ # 将linux目录复制到远端test目录下的linux
rsync -a source/ test@10.211.55.5:/destination # 仅同步源目录source里面的内容到目标目录destination
rsync -av --delete source/ destination # 完全保持一致
rsync -av --exclude='*.txt' source/ destination # 排除文件
rsync -av --exclude='.*' source/ destination # 排除隐藏文件
rsync -av --include="*.txt" --exclude='*' source/ destination # 仅同步指定文件
```
