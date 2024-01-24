# select

`select`结构主要用来生成简单的菜单。它的语法与`for...in`循环基本一致。
```bash
#!/usr/bin/env bash

echo "Which Operating System do you like?"

select os in Ubuntu LinuxMint Windows8 Windows7 WindowsXP
do
  case $os in
    "Ubuntu"|"LinuxMint")
      echo "I also use $os."
    ;;
    "Windows8" | "Windows7" | "WindowsXP")
      echo "Why don't you try Linux?"
    ;;
    *)
      echo "Invalid entry."
      break
    ;;
  esac
done
```