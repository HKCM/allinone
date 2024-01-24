# break

`break`命令立即终止循环,程序继续执行循环块之后的语句,即不再执行剩下的循环

```shell
#!/usr/bin/env bash

for number in 1 2 3 4 5 6
do
  echo "number is $number"
  if [ "$number" = "3" ]; then
    break
  fi
done
```

上面例子只会打印3行结果。一旦变量`$number`等于3,就会跳出循环,不再继续执行