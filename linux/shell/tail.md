# tail

```bash
tail -n +15 /etc/passwd   #从第15行开始显示文件
nobody:x:99:99:Nobody:/:/sbin/nologin
dbus:x:81:81:System message bus:/:/sbin/nologin
...
```

```bash
tail -F /path/log # 等待文件生成并输出
```