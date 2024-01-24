# until

until循环与while循环恰好相反,只要不符合判断条件(判断条件失败),就不断循环执行指定的语句

一旦符合判断条件,就退出循环

```shell
#!/usr/bin/env bash

number=0
until [ "$number" -ge 10 ]; do
  echo "Number = $number"
  number=$((number + 1))
done
```