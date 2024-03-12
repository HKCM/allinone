# vim

## 配置

`~/.vimrc` 配置文件

```
set nu/nonu   # 显示行号
set ruler     # 显示行号
set hlsearch/nohlsearch  # 搜索高亮
set incsearch # 增量搜索
set mouse=a # 启用鼠标
set syntax on/off # 语法高亮
set encoding=utf8 # 选项表示 Vim 内部使用的字符编码 这是防止乱码的关键设置
```

## 打开文件
vim +n file # 打开文件第n行
vim + file # 打开文件最后一行
vim +/patten file # 打开文件patten所在位置
vim +/"so cute" file # 打开文件patten所在位置
vim file1 file2 # 打开多个文件,使用:n :N切换 :ar显示当前的编辑的文件
vim -O abc 123 # 同时打开两个文件 ctrl + w 之后加左右箭头



显示号: :set nu
保存退出: `:wq` `ZZ`
取消未保存的修改: `:e!`

## 移动
左: h
下: j
上: k
右: l
行首: 0
行尾: $
下一个单词: w W
上一个单词: b B
移动到行首: 0 ^
移动到行尾: $
移动到首行: gg 1G
移动到尾行: G
移动到当前句子开头: (
移动到当前句子结尾: ) 
向下滚动一页: ctrl+f
向上滚动一页: ctrl+b
向下滚动半页: ctrl+d
向上滚动半页: ctrl+u

## 编辑
替换: r R
行首插入: I
行尾插入: A
下行插入: o
上行插入: O
更改整个单词: cw
更改整行: cc S
更改到整首: c0
更改到整尾: c$ C
删除整行: dd
删除到行首: d0
删除到行尾: d$ D
撤销: u
撤销整行: U
大小写切换: ～


查找: /
在一行中查找: fx


## 命令模式

=:显示总行数
!command: 执行命令,只能执行环境变量里的命令


s:替换

```bash
:1s/11/22:将第一行的第一次出现的 11 换成 22
:n,ms/old/new/g:替换行号n和m之间所有old
:%s/old/new/g:替换整个文件中的所有old
:g/pattern/s/old/new/g:将全文包含 pattern 的行中的 old 替换为 new
:%s/\(That\) or \(this\)/\2 or \1/:将That or this改为this or That
:%s/apple/&s/g:将apple替换为apples
:%s/\<child\>/children/g:将child替换为children,区别是只会完全匹配child,而不会匹配childish
# 示例:更改文件中的路径:%s:/home/user:/home/bear:g 
# 示例:将第一到第十行的中句号改为分号:1,10s/\./;/g
# 示例:将 help 改为 HELP:%s/[Hh]elp/HELP/g
# 示例:将一个或多个空格改为一个空格:%s/  */ /g
# 示例:将冒号后一个或多个空格改为 2 个空格:%s/:  */:  /g
# 示例:删除每一行开头的空白:%/^  *\(.*\)/\1/:%/^  *//
# 示例:在每一行开头添加>  :%s/^/>  /# 为后续 6 行添加冒号:.,+5s/$/:/
# 示例:把文件变为大写:%s/.*/\U&/
```

d:删除

:%d:删除所有内容
:1d:删除第一行
:1,10d:删除1到10行
:.,$d:删除当前行到结尾之间所有的行
:.,.+20d:删除当前行到结尾之间所有的行
:/pattern/d:删除下一个包含pattern的行
:/pattern1/,/pattern2/d:从第一个包含patternI的行删除到第一个包含pattern。2的行

m:移动

:7m1:将第七行的内容移动到第一行之后
:10,.m$:将第十行到当前行的内容移动到最后一行
:100,$m.-2:将第一百行到文件结尾间的行移到当前这一行的两行之前

co:复制

:7m1:将第七行的内容复制到第一行之后

:= 列出文件总行数
:.= 列出当前所在行的行号
:/pattern/= 列出pattern第一次出现时的行号

w:保存

:w newfile 将当前更改保存到新文件
:100,$w newfile 将第一百行到文件结尾保存到新文件
:1,10w newfile 将第一到第十行保存到新文件
:40,$w >>newfile 将第四十到结尾追加到新文件

r:读取其他文件
:r file:将file内容添加到当前行下
:100r file :将file内容添加到第一百行下

ab:定义缩写
:ab:显示当前存在的缩写
:ab imrc International Materials Research Center:在插入模式写入imrc即可实现全文写入

## 问题

vim打开文件乱码,但是 `cat` 能正常显示

在用户目录的.vimrc中加入如下内容
```bash
set termencoding=utf-8 # 表示终端的字符编码
set encoding=utf8 # 选项表示 Vim 内部使用的字符编码 这是关键设置
# 选项表示打开和保存文件时要尝试的字符编码列表
set fileencodings=utf8,ucs-bom,gbk,cp936,gb2312,gb18030,Shift-JIS 
```