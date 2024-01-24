# while

## 数字循环 

```shell
#!/bin/bash

num=1

while [ $num -le 10 ]
do
    echo $num
    num=$(( $num + 1 ))
done
```

## 读取文件

```shell
while read line
do
    echo $line
done <./a.txt
```

## 死循环

```shell
while true
do
    command
done
```

## 批量创建用户
```shell
$ cat > users.csv <<EOF
rich,Richard Blum 
christine,Christine Bresnahan 
barbara,Barbara Blum 
tim,Timothy Bresnahan
EOF

#!/usr/bin/env bash
input=users.csv

while IFS=',' read -r userid name
do
  echo "adding $userid"
  useradd -c $name -m $userid
done < "$input"
```