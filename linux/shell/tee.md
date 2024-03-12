# tee

tee命令允许标准输出同时把内容写入(覆盖)到文件中, 在输出log并追加到文件很好用


```bash
cat test.txt 
Hi, this is a test file.
I just want to say I like dog.
cat test.txt |tee -a test.copy
Hi, this is a test file.
I just want to say I like dog.
ll
-rw-r--r-- 1 root root   56 Mar 11 02:22 test.copy
-rw-r--r-- 1 root root   56 Mar 11 01:56 test.txt
```