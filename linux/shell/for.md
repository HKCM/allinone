

# for循环的常见用法


## 列表循环

```bash
envs=(test int dev stage prod)
for env in ${envs[@]}; do 
    echo "$env"
done
```

## 参数循环

```bash
for a in $@
do
    echo $a
done
```

## 数字循环

```bash
for((i=1;i<=10;i++));do 
    echo $(expr $i \* 3 + 1);
done
```

## 数字循环

```bash
for i in $(seq 1 10)
do 
    echo $(expr $i \* 3 + 1);
done
```

### 当前目录文件循环

```bash
for i in $(ls);do 
    echo $i is file name\!;
done
```

遍历目录
```shell
for file in /home/rich/test/*
do
if [ -d "$file" ] then
  echo "$file is a directory" 
elif [ -f "$file" ]
then
  echo "$file is a file"
fi
done
```

## 字符循环

```bash
services="a b c d"
for service in ${services}; do
    echo "$service is cool"
done
```

## 普通循环

```bash
for i in f1 f2 f3;do
    echo $i is appoint;
done
```

