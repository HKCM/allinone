# continue

`continue`命令立即终止本轮循环,开始执行下一轮循环

```shell
#!/usr/bin/env bash

while read -p "What file do you want to test?" filename
do
  if [ ! -e "$filename" ]; then
    echo "The file does not exist."
    continue
  fi

  echo "You entered a valid file.."
done
```