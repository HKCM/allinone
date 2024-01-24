# terminal配置

keyword: color ssm提示符 shell提示符

可以写入`~/.bashrc`,也可作为SSM的登录配置

```bash
export EDITOR=vim
# 添加用户目录
export PATH=/home/admin/bin:/data/admin:$PATH
export PS1='\[\e[1;38;5;135m\]\u\[\e[0m\]\[\e[1;38;5;226m\]@\[\e[0m\]\[\e[1;38;5;200m\]\h\[\e[0m\] [\[\e[1;36m\]\w\[\e[0m\]]$(if [[ $? = "0" ]]; then echo "\[\e[1;32m\]"; else echo "\[\e[1;31m\]"; fi)\[\e[0m\]\n[\[\e[1;92m\]\D{%Y-%m-%d %H:%M:%S}\[\e[0m\]] > '
export PS2='> '
```

