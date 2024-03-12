# chattr

i: 对文件设置i属性,不允许对文件进行删除、改名,也不能添加和修改数据；对目录设置i属性,只能修改目录下文件中的数据,不允许建立和删除文件
a: 对文件设置a属性,只能在文件中增加数据,不能删除和修改数据；目录设置a属性,那么只允许在目录中建立和修改文件,但是不允许删除文件
e: Linux中的绝大多数文件都默认拥有e属性,表示该文件是使用ext文件系统进行存储的,而且不能使用“chattr -e”命令取消e属性

```bash
#给文件赋予i属性
touch ftest #建立测试文件
chattr +i ftest
rm -rf ftest #赋予i属性后,root也不能删除
rm: 无法删除"ftest": 不允许的操作
echo 111 >> ftest #也不能修改文件中的数据
-bash: ftest: 权限不够

#给目录赋予i属性
mkdir dtest #建立测试目录
touch dtest/abc #再建立一个测试文件abc
chattr +i dtest/ #给目录赋予i属性
cd dtest/
touch bcd #dtest目录不能新建文件
touch: 无法创建"bcd": 权限不够
echo 11 >> abc #但是可以修改文件内容
cat abc
11
rm -rf abc #不能删除
rm: 无法删除"abc": 权限不够
```

```bash
mkdir -p /back/log #建立备份目录
chattr +a /back/log/ #赋予a属性
cp /var/log/messages /back/log/ #可以复制文件和新建文件到指定目录中
rm -rf /back/log/messages #但是不允许删除
rm: 无法删除"/back/log/messages": 不允许的操作
```

```bash
lsattr # 查看attribute权限
```