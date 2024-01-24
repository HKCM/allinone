# alias

如果别名特别多,可以创建单独的`~/.alias`文件存放别名,并在`~/.bashrc`中读取
```bash
if [ -f ~/.alias ]; then
  source ~/.alias
fi
```

alias相关操作
```bash
$ alias hw='echo "hello world"' # 新增别名 只在当前终端有效
$ alias      # 查看现有别名
$ unalias ll # 取消单个别名
$ unalias -a # 取消所有别名 如果别名是写在文件中,即使取消了别名,重新登陆时别名还是存在的.
```

常见alias配置
```bash
alias l="ls -lah --color=always"
alias ll="ls -lah --color=always"
alias mkdir='mkdir -pv'
# 防止文件覆盖
alias rm="rm -i"
alias cp="cp -i"
alias mv="mv -i"

# 获取占用内存的进程排名
alias psmem='ps auxf | sort -nr -k 4'
alias psmem10='ps auxf | sort -nr -k 4 | head -10'

alias grep = 'grep --color=auto'

# 获取占用 cpu 的进程排名
alias pscpu='ps auxf | sort -nr -k 3'
alias pscpu10='ps auxf | sort -nr -k 3 | head -10'

alias phistory="history | awk '{CMD[\$2]++;count++;} END { for (a in CMD)print CMD[a] \" \" CMD[a]/count*100 \"% \" a;}' | grep -v './' | column -c3 -s ' ' -t | sort -nr | nl |head -n10" # mac无效

alias checkip="curl checkip.amazonaws.com"
```