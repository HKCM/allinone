# session manager 配置

首先在[Session manager](https://us-east-1.console.aws.amazon.com/systems-manager/session-manager/preferences?region=us-east-1)中添加配置

```bash
/bin/bash
export EDITOR=vim
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin:/home/admin/bin:/data/ok
export PS1='\[\e[1;38;5;135m\]\u\[\e[0m\]\[\e[1;38;5;226m\]@\[\e[0m\]\[\e[1;38;5;200m\]\h\[\e[0m\] [\[\e[1;36m\]\w\[\e[0m\]]$(if [[ $? = "0" ]]; then echo "\[\e[1;32m\]"; else echo "\[\e[1;31m\]"; fi)\[\e[0m\]\n[\[\e[1;92m\]\D{%Y-%m-%d %H:%M:%S}\[\e[0m\]] > '
export PS2='> '
alias l="ls -lah --color"
alias ll="ls -lah --color"
alias rm="rm -i"
alias cp="cp -i"
alias mv="mv -i"
alias grep="grep --color"
cd /data/ok || cd /data || cd ~/
```