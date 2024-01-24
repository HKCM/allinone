# Peer's Certificate issuer is not recognized

```bash
fatal: unable to access 'https://*****/xx.git/': Peer's Certificate issuer is not recognized.
```

如果gitlab仓库是信任的,只需要设置跳过SSL证书验证就可以,执行以下命令:
```bash
git config http.sslVerify false
```

# Please commit your changes or stash them before you merge

解决方法:

- git stash暂存。
- git pull拉取远程仓库代码。
- git stash pop释放暂存修改。

当有多个设备都向同一个仓库中提交文件时,有时就会出现这种情况:

在A设备是对文件file进行了修改,并提交到git仓库里面去了。
随后,应该在B设备上面先检出最新的修改到本地,但忘记了这一步。
然后,在B设备上面也修改了文件file,然后尝试从远程拉取最新的代码,就会出现这种异常。

# git status不能显示中文

```bash
git config --global core.quotepath false
```

在`~/.zshrc` 或 `~/.bashrc` 中添加
```bash
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8
```