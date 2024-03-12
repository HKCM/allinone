# shasum

关键词: 哈希 文件校验 sha1 sha128 sha256 sha512 shasum sha256sum sha512sum

## SHA-1

```shell
shasum file.dmg
```

## SHA256

```shell
shasum -a 256 file.dmg
# 以将输出的校验和重定向到一个文件
sha256sum file.dmg > file.dmg.sha256
# 检查校验和是否匹配,会自动匹配对应的文件名
sha256sum -c file.dmg.sha256
```

## SHA512

```shell
shasum -a 512 file.dmg
# 以将输出的校验和重定向到一个文件
sha512sum file.dmg > file.dmg.sha512
# 检查校验和是否匹配,会自动匹配对应的文件名
sha512sum -c file.dmg.sha512
```



