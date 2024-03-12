# cat

-n: 对输出显示行号
-b: 与-n类似但是会忽略空行
-s: 将连续两行空行合并为一个空行
-E: 行尾显示$符号

## 写入文件

```bash
# 方式一:手动结束输入
cat > test.txt
1
2
# 在新的空白行使用Ctrl+c结束输入

# 方式二:使用标识结束输入
file_name=hello.txt
cat > $file_name <<EOF
hello world!
EOF

# 追加到文件
cat >> $file_name <<EOF

be happy!
EOF
```

## usage

```bash
#!/usr/bin/env bash

# 输出到stderr
function usage() {
    cat 1>&2 <<EOF
script-init
initializes a new installation

USAGE:
    script-init [FLAGS] [OPTIONS] --data_dir <PATH> --pubkey <PUBKEY>

FLAGS:
    -h, --help              Prints help information
        --no-modify-path    Don't configure the PATH environment variable

OPTIONS:
    -d, --data-dir <PATH>    Directory to store install data
    -u, --url <URL>          JSON RPC URL for the solana cluster
    -p, --pubkey <PUBKEY>    Public key of the update manifest
EOF
}
```

```bash
#!/bin/bash
exportfs_usage()
{
    cat <<END
    USAGE: $0 {start|stop|monitor|status|validate-all}
END
}
exportfs_usage
```

tail: 按顺序输出文件最后的内容
tac: 完全倒序输出文件,和tail完全不一样