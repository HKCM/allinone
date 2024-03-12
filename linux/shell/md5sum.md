# md5sum

```bash
echo 123 > test.txt # 创建测试文件
md5sum test.txt > test.md5 # 将md5写入文件
md5sum -c test.md5 # 测试文件是否被修改
test.txt: OK
echo 456 > test.txt # 修改文件
md5sum -c test.md5  # 测试失败
test.txt: FAILED
md5sum: WARNING: 1 computed checksum did NOT match
```