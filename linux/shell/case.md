# case

```shell
#!/bin/bash

echo 'Input a number between 1 to 4'
echo 'Your number is:'
read n
case $n in
    1)  echo 'You select 1'
    ;;
    2)  echo 'You select 2'
    ;;
    3)  echo 'You select 3'
    ;;
    4)  echo 'You select 4'
    ;;
    *)  echo 'You do not select a number between 1 to 4'
    ;;
esac
```

case的匹配模式可以使用各种通配符,下面是一些例子。

* a):匹配a。
* a|b):匹配a或b。
* [[:alpha:]]):匹配单个字母。
* ???):匹配3个字符的单词。
* *.txt):匹配.txt结尾。
* *):匹配任意输入,通过作为case结构的最后一个模式。

```Shell
#!/usr/bin/env bash

read -r -p "Are You Sure? [Y/n] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Yes"
		;;
    [nN][oO]|[nN])
		echo "No"
       	;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac
```