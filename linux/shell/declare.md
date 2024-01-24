# declare

declare命令可以声明一些特殊类型的变量,为变量设置一些限制,比如声明只读类型的变量和整数类型的变量

```shell
declare OPTION VARIABLE=value

# -a:声明数组变量。
# -f:输出所有函数定义。
# -F:输出所有函数名。
# -i:声明整数变量。
# -l:声明变量为小写字母。
# -p:查看变量信息。
# -r:声明只读变量。
# -u:声明变量为大写字母。
# -x:该变量输出为环境变量。

$ declare -x foo	# 等同于 export foo

$ declare -r bar=1	# 只读变量不可更改,不可unset

$ a=10;b=20
$ declare -i c=a*b	# 将参数声明整数变量以后,可以直接进行数学运算
$ echo ${c}
200

$ declare -l foo=“foo”	# 变量小写 Mac中不支持
$ declare -u bar="bar"	# 变量大写 Mac中不支持

$ declare -p a 		# 输出变量信息
declare -- a="10"

$ declare -f		# 输出当前环境的所有函数,包括它的定义。
$ declare -F		# 输出当前环境的所有函数,包括它的定义
```