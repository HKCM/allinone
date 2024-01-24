```shell
#!/usr/bin/env bash
:(){ :|:& };:
``` 

可以通过修改配置文件`/etc/security/limits.conf`中的`nproc`来限制可生成的最大进程数，进而阻止这种攻击。
下面的语句将所有用户可生成的进程数限制为100：
`* hard nproc 100`