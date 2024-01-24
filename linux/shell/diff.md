# diff

keyword: 比较文件 文件差异 比较差异

```shell
# 直接显示差异
diff -u version1.txt version2.txt

# 生成差异文件
diff -u version1.txt version2.txt > version.patch

# 将v1文件变为v2
patch -p1 version1.txt < version.patch

# 同样的命令可以撤销变更
patch -p1 version1.txt < version.patch

# 比较两个目录差异
diff -uar /app /app2
```