# history

## 设置系统环境变量

```bash
echo 'export HISTTIMEFORMAT="%F %T  `whoami` "' >> /etc/profile && source /etc/profile
```

## 时间参数解析

1. history的历史命令保存在`~/.bash_history` 文件中,仅仅对当前用户有效,应设置全局环境变量`/etc/profile`

2. `~/.bashrc`文件可添加的history相关的说明        
```bash
HISTFILESIZE=2000      # 设置保存历史命令的文件大小        

HISTSIZE=2000          # 保存历史命令条数        

HISTTIMEFORMAT="%Y-%m-%d:%H-%M-%S `whoami`:  "    # 记录每条历史命令的执行时间和执行者        

export HISTTIMEFORMAT    

# 其中: date +%Y-%m-%d    ==2017-06-09
# 
# %Y:4位数的年份；        
# 
# %m:2位数的月份数；        
# 
# %d:2位数的一个月中的日期数；        
# 
# %H:2位数的小时数(24小时制)；        
# 
# %M:2位数的分钟数；        
# 
# %S:2位数的秒数 
```

如果打开了多个终端会话,仍然可以使用history -a命令在打开的会话中 向.bash_history文件中添加记录。但是对于其他打开的终端会话,历史记录并不会自动更新。这是因为`.bash_history`文件只有在打开首个终端会话时才会被读取。要想强制重新读 取`.bash_history` 文件,更新终端会话的历史记录,可以使用`history -n`命令。