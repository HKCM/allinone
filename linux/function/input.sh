
#!/usr/bin/env bash
# 常用选项
# -n 设置输入字符数限制
# -t 设置等待时间,单位是秒
# -p 设置提示信息
# -s 使得用户的输入不显示在屏幕上,这常常用于输入密码或保密信息。
echo -n "输入一些文本 > "
read text
echo "你的输入:$text"

if read -s -t 5 -p "Please enter your name: " name
then
  echo "Hello $name, welcome to my script"
else
  echo 
  echo "Sorry too slow"
fi


