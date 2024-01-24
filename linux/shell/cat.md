# cat

## 写入文件

```bash
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