# tr

-d 删除字符
-s 保留连续字符的第一个字符,其他删除

测试文本
```bash
cat > test.txt <<EOF
Hi, this is a test file.
I just want to say I like dog.
EOF

# 原本想将dog改为cat 可以看到与预期不一致
# 凡是在文本中出现的“d”均应转换成“c”，“o”均应转换成“a”，“g”均应转换成“t”，而不是仅仅将字符串“dog”替换为字符串“cat”。
tr "dog" "cat" < test.txt 
Hi, this is a test file.
I just want ta say I like cat. # to 变为了 ta

# 可以通过sed做到
sed "s/dog/cat/g" test.txt 

# 转化为大写
tr "[a-z]" "[A-Z]" < test.txt 
HI, THIS IS A TEST FILE.
I JUST WANT TO SAY I LIKE DOG.

# 删除换行符
tr -d "\n" < test.txt 
Hi, this is a test file.I just want to say I like dog.

# 压缩相同字符
echo "dddooogggggg cat ddog" | tr -s dog
dog cat dog
```