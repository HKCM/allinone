# set

`set`命令用来修改子 Shell 环境的运行参数,即定制环境。一共有十几个参数可以定制,[官方手册](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html)有完整清单,本章介绍其中最常用的几个。

[参考](https://www.ruanyifeng.com/blog/2017/11/bash-set.html)

```shell
#!/usr/bin/env bash

set -u # 遇到变量不存在则报错
set -x # 命令运行前，先输出,测试时使用
set -e # 与 set -o errexit 相同, 脚本只要发生错误，就终止执行
set -o pipefail # 只要一个子命令失败，整个管道命令就失败，脚本终止执行
# set -o xtrace
# 写法一
set -euxo pipefail
# 写法二
set -eux
set -o pipefail
```

如果命令可能失败,但是希望继续运行
```shell
command || true
# 或在某段代码前暂时关闭set -e
set +e
command1
command2
set -e
```

写脚本时
```shell
set -euxo pipefail
# 或
set -eux
set -o pipefail
```

## 示例

关于`pipefail`

在管道命令中, Bash 会把最后一个子命令的返回值，作为整个命令的返回值。也就是说，只要最后一个子命令不失败，管道命令总是会执行成功，因此它后面命令依然会执行，`set -e`就失效了,设置`set -o pipefail`可以解决这个问题.

示例1

```bash
#!/usr/bin/env bash
set -e

foo | echo a
echo bar
```

当没有`set -o pipefail`时,shell执行了`echo bar`命令,程序没有停下来,这是不符合预期的
```bash
$ bash script.sh
a
script.sh:行4: foo: 未找到命令
bar
```

示例2

```bash
#!/usr/bin/env bash
set -e
set -o pipefail
foo | echo a
echo bar
```

设置`set -o pipefail`后,shell没有执行`echo bar`命令,程序停下来了,符合预期
```bash
$ bash script.sh
test.sh: line 4: foo: command not found
a
```



