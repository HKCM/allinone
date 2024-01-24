# 配置

- `/etc/gitconfig` 文件:系统中对所有用户都普遍适用的配置。若使用 git config 时用 `--system` 选项,读写的就是这个文件。
- `~/.gitconfig` 文件:用户目录下的配置文件只适用于该用户。若使用 git config 时用 `--global` 选项,读写的就是这个文件。
- `.git/config` 文件:目的 Git 目录中的配置文件(也就是工作目录中的 .git/config 文件)针对当前项目有效

## 邮箱用户名配置

```shell
# 项目级别,只对当前仓库有效,位于.git/config目录中
$ git config user.name Ace
$ git config user.email ace.bruce@example.com

# --global 选项,配置位于~/.gitconfig
$ git config --global user.name "w3c"
$ git config --global user.email w3c@w3cschool.cn
$ git config --global core.editor vim
$ git config --list

git config --global credential.helper cache # 密码缓存时间
```

**Note**: 由于 Git 会从多个文件中读取同一配置变量的不同值,因此可能会在其中看到意料之外的值而不知道为什么。
例如global `user.name = A`, repo中 `.git/config` 的 `username = B`
```shell
git config --show-origin user.name  # 这会查找username生效的配置文件
git config --global --unset user.name # 取消配置
```

## 中文显示配置

```bash
git config core.quotepath
git config core.quotepath false
```

## commit模版配置

```bash
git config commit.template
git config commit.template .commit.template
vim .commit.template
```

```text
docs(Git): subject

<body>

<footer>

# 标题行:50个字符以内,描述主要变更内容,按以下形式填写
# <type>(<scope>): <subject>
# - type 提交的类型
#    feat:新特性, 
#    fix:修改问题, 
#    docs:文档修改,仅仅修改了文档,比如 README, CHANGELOG, CONTRIBUTE等等, 
#    style:代码格式修改, 修改了空格、格式缩进、逗号等等,不改变代码逻辑,注意不是css修改, 
#    refactor:代码重构,没有加新功能或者修复 bug, 
#    perf: 优化相关,比如提升性能、体验,
#    test:测试用例,包括单元测试、集成测试等, 
#    chore:其他修改, 比如改变构建流程, 修改依赖管理,增加工具等,
#    revert: 回滚到上一个版本。
# - scope: 影响的的范围
#    影响的的范围,可以为空,如Git、CentOS等模块,或者全局All。
# - subject
#    主题,提交描述
#
# 主体内容:更详细的说明文本,建议72个字符以内。 需要描述的信息包括:
#
# * 为什么这个变更是必须的? 
#   它可能是用来修复一个bug,增加一个feature,提升性能、可靠性、稳定性等等
# * 它如何解决这个问题? 具体描述解决问题的步骤
# * 是否存在副作用、风险? 
#
# 尾部:如果需要的话可以添加一个链接到issue地址或者其它文档,或者关闭某个issue。
# 
# 注意,标题行、主体内容、尾部之间都有一个空行！

```

# 命令

## 常用命令

### 重新提交

有时候提交完了才发现漏掉了几个文件没有添加,或者提交信息写错了。 此时,可以运行带有 —amend 选项的提交命令来重新提交:

```shell
$ git commit -m 'initial commit'
$ git add forgotten_file
$ git commit --amend
```

### 删除分支

```shell
# 删除本地分支
$ git branch -d testing
$ git branch -D testing # 强制删除
# 删除远端分支
$ git push origin --delete serverfix
```

### 跟踪远程分支

```shell
$ git checkout -b <branch> <remote>/<branch>
$ git checkout --track origin/serverfix
Branch serverfix set up to track remote branch serverfix from origin.
Switched to a new branch 'serverfix'

# 设置不同名字的本地分支和远程分支
$ git checkout -b sf origin/serverfix
```


### 测试连接

```shell
$ ssh -vT git@github.com
```

### help
```shell
$ git help <verb>
$ git help reset
$ git <verb> --help
$ man git-<verb>
```

## Git操作
### 创建仓库
```shell
$ mkdir newrepo && cd newrepo
$ git init

# or
$ git init newrepo
```

### clone

```shell
$ git clone [url]

# or 克隆并改名
$ git clone git://github.com/schacon/grit.git mygit
```

### ignore

在本地库中添加`.gitignore` 文件并`add`
```shell
# 忽略所有的 .a 文件
*.a
# 但跟踪所有的 lib.a,即便你在前面忽略了 .a 文件
!lib.a
# 只忽略当前目录下的 TODO 文件,而不忽略 subdir/TODO
/TODO
# 忽略任何目录下名为 build 的文件夹
build/
# 忽略 doc/notes.txt,但不忽略 doc/server/arch.txt
doc/*.txt
# 忽略 doc/ 目录及其所有子目录下的 .pdf 文件
doc/**/*.pdf
```

GitHub 有一个十分详细的针对数十种项目及语言的 `.gitignore` 文件列表, 可以在 https://github.com/github/gitignore 找到它。

一个仓库可能只根目录下有一个 .gitignore 文件,它递归地应用到整个仓库中。 然而,子目录下也可以有额外的 `.gitignore` 文件。子目录中的 `.gitignore` 文件中的规则只作用于它所在的目录中。


### status

```shell
git status
git status -s
git status --ignored # 查看被忽略的文件
```

### add
```shell
$ git add fileA fileB
$ git add .
```

### diff

- git diff: 显示尚未暂存的改动
- git diff -—cached: 显示已暂存文件与最后一次提交的文件差异
- git diff HEAD: 查看已缓存的与未缓存的所有改动
-  git diff -—stat: 显示摘要而非整个
```shell
# 比较工作区文件和暂存区文件的差异
$ git diff
$ git diff <file> # 比较具体文件

# 将比对已暂存文件与最后一次提交的文件差异
$ git diff --cached
$ git diff --staged
```

git diff 本身只显示尚未暂存的改动,而不是自上次提交以来所做的所有改动。 所以有时候一下子暂存了所有更新过的文件,运行 git diff 后却什么也没有,就是这个原因。

### commit
```shell
$ git commit 
$ git commit -m "message"

# 只会针对已追踪的文件,未追踪的文件不会提交
$ git commit -a
$ git commit -am
```

### log
https://git-scm.com/book/zh/v2/Git-%E5%9F%BA%E7%A1%80-%E6%9F%A5%E7%9C%8B%E6%8F%90%E4%BA%A4%E5%8E%86%E5%8F%B2

```shell
# 一行显示
$ git log --pretty=oneline
e4d8f95d988c5dd61e58a781a6c2f7fa80e051ed (HEAD -> master) t git commit -a
32c9c57d50a49b56ebe4e459fa2553888365b28f first commit

$ git log --oneline
e4d8f95 (HEAD -> master) t git commit -a
32c9c57 first commit

$ git reflog
# 以补丁形式显示最近两次历史patch
$ git log -p -2
$ git log --pretty=format:"%h - %an, %ar : %s"

# 列出最近两周的所有提交:
$ git log --since=2.weeks

# 过滤器 -S 它接受一个字符串参数,并且只会显示那些添加或删除了该字符串的提交
# 假设想找出添加或删除了对某一个特定函数的引用的提交
$ git log -S function_name
```
### reset

- git reset -—soft: 仅移动本地库HEAD指针,暂存区不变,工作区不变
- git reset -—mixed: 移动本地库HEAD指针,重置暂存区,工作区不变
- git reset -—hard: 移动本地库HEAD指针,重置暂存区,重置工作区

```shell
$ git reset --hard e4d8f95
# 往上移动一个版本
$ git reset --hard HEAD^
# 往上移动3个版本
$ git reset --hard HEAD～3
# 还原到最新的提交
$ git reset --hard HEAD
```

**找回文件**
已经commit过的文件找回,用`git reset --hard`切到文件曾经存在的版本进行找回


### rm

另外一种情况是,我们想把文件从 Git 仓库中删除(亦即从暂存区域移除),但仍然希望保留在当前工作目录中。 换句话说,想让文件保留在磁盘,但是并不想让 Git 继续跟踪。 当忘记添加 .gitignore 文件,不小心把一个很大的日志文件或一堆 .a 这样的编译生成文件添加到暂存区时,这一做法尤其有用。 为达到这一目的,使用 —cached 选项:
```shell
# git rm 命令后面可以列出文件或者目录的名字
$ git rm --cached README.md README

# 也可以使用 glob 模式。比如:
# 删除 log/ 目录下扩展名为 .log 的所有文件
$ git rm log/\*.log

# 删除所有名字以 ~ 结尾的文件
$ git rm \*~
```

### mv
```shell
$ git mv README.md README

# 或者
$ mv README.md README
$ git rm README.md
$ git add README
```



### 撤消操作


最终你只会有一个提交——第二次提交将代替第一次提交的结果。

#### 取消暂存的文件
```shell
$ git reset HEAD <file>
$ git restore --staged <file>
```

#### 撤消对文件的修改
```shell
$ git checkout -- <file>
$ git restore <file>
```
请务必记得 git checkout — <file> 是一个危险的命令。 那个文件在本地的任何修改都会消失——Git 会用最近提交的版本覆盖掉它。 除非确实清楚不想要对那个文件的本地修改了,否则请不要使用这个命令。

### 远程仓库

列出仓库
```shell
$ git remote -v
```

添加仓库
```shell
$ git remote add pb https://github.com/paulboone/ticgit
```

查看仓库
```shell
$ git remote show origin
* remote origin
  URL: https://github.com/my-org/complex-project
  Fetch URL: https://github.com/my-org/complex-project
  Push  URL: https://github.com/my-org/complex-project
  HEAD branch: master
  Remote branches:
    master                           tracked
    dev-branch                       tracked
    markdown-strip                   tracked
    issue-43                         new (next fetch will store in remotes/origin)
    issue-45                         new (next fetch will store in remotes/origin)
    refs/remotes/origin/issue-11     stale (use 'git remote prune' to remove)
  Local branches configured for 'git pull':
    dev-branch merges with remote dev-branch
    master     merges with remote master
  Local refs configured for 'git push':
    dev-branch                     pushes to dev-branch                     (up to date)
    markdown-strip                 pushes to markdown-strip                 (up to date)
    master                         pushes to master                         (up to date)
```

这个命令列出了当你在特定的分支上执行 git push 会自动地推送到哪一个远程分支。 它也同样地列出了哪些远程分支不在你的本地,哪些远程分支已经从服务器上移除了, 还有当你执行 git pull 时哪些本地分支可以与它跟踪的远程分支自动合并。

#### 远程仓库中抓取
```shell
$ git fetch <remote>
```

这个命令会访问远程仓库,从中拉取所有你还没有的数据。 执行完成后,你将会拥有那个远程仓库中所有分支的引用,可以随时合并或查看。

#### 移除远程仓库
如果因为一些原因想要移除一个远程仓库——你已经从服务器上搬走了或不再想使用某一个特定的镜像了, 又或者某一个贡献者不再贡献了——可以使用 git remote remove 或 git remote rm :
```shell
$ git remote remove paul
```

一旦你使用这种方式删除了一个远程仓库,那么所有和这个远程仓库相关的远程跟踪分支以及配置信息也会一起被删除。

### 追踪远程仓库的分支
```shell
git checkout -b release-21.4.20 origin/release-21.4.20
```

### 删除远端分支
```shell
git push origin --delete <branchName>
```

### tag
```shell
$ git tag # 列出标签
$ git tag -l "v1.8.5*" # 列出带有v1.8.5标签
```

Git 支持两种标签:轻量标签(lightweight)与附注标签(annotated)

轻量标签本质上是将提交校验和存储到一个文件中——没有保存任何其他信息。
```shell
$ git tag v1.4-lw
$ git tag
v0.1
v1.3
v1.4
v1.4-lw
v1.5

$ git show v1.4-lw
commit ca82a6dff817ec66f44342007202690a93763949
Author: Scott Chacon <schacon@gee-mail.com>
Date:   Mon Mar 17 21:52:11 2008 -0700

    changed the version number
```

附注标签是存储在 Git 数据库中的一个完整对象, 它们是可以被校验的,其中包含打标签者的名字、电子邮件地址、日期时间, 此外还有一个标签信息,并且可以使用 GNU Privacy Guard (GPG)签名并验证。 

通常会建议创建附注标签,这样可以拥有以上所有信息。
```shell
$ git tag -a v1.4 -m "my version 1.4"
$ git tag
v0.1
v1.3
v1.4

$ git show v1.4
tag v1.4
Tagger: Ben Straub <ben@straub.cc>
Date:   Sat May 3 20:19:12 2014 -0700

my version 1.4

commit ca82a6dff817ec66f44342007202690a93763949
Author: Scott Chacon <schacon@gee-mail.com>
Date:   Mon Mar 17 21:52:11 2008 -0700

    changed the version number
```
默认情况下,git push 命令并不会传送标签到远程仓库服务器上。 在创建完标签后必须显式地推送标签到共享服务器上。 
```shell
$ git push origin v1.5 # 推送单个标签
Counting objects: 14, done.
Delta compression using up to 8 threads.
Compressing objects: 100% (12/12), done.
Writing objects: 100% (14/14), 2.05 KiB | 0 bytes/s, done.
Total 14 (delta 3), reused 0 (delta 0)
To git@github.com:schacon/simplegit.git
 * [new tag]         v1.5 -> v1.5

$ git push origin --tags # 推送全部标签

$ git tag -d v1.4-lw # 删除本地标签
$ git push origin --delete <tagname> # 删除远端标签
```

### Rebase
git rebase

作用一:合并提交记录,让多次提交变为一次提交
```shell
# 合并提交记录,不要合并已push的记录(即已推送到远程库的记录)
git rebase -i <hashcode>
```

作用二:让git log变为线形
```shell
# 在dev分支 开发完成后,将master的最新改动同步到dev分支
git checkout dev
git rebase master

# 再merge dev分支,此时dev分支将变为fast forward,git log将变为一条线
git checkout master
git merge dev
```

作用三:pull代码时远程库和本地库有冲突
```shell
# 原本是 git pull origin dev
git fetch origin dev
git rebase origin/dev
```

rebase 冲突
```shell
git rebase master

# 如果有冲突,就解决冲突,解决完之后
git rebase --continue
```

### 别名

创建别名
```shell
$ git config --global alias.co checkout
$ git config --global alias.br branch
$ git config --global alias.ci commit
$ git config --global alias.st status

$ git config --global alias.unstage 'reset HEAD --' # 取消暂存别名
$ git config --global alias.last 'log -1 HEAD'  # 查看最后一次提交
```

### 分支操作

列出分支
```shell
$ git branch
$ git fetch --all; git branch -vv
```

列出已merge的分支
```shell
$ git branch --merged
```

列出未merge的分支
```shell
$ git branch --no-merged

$ git branch --no-merged master # 列出尚未合并到master的分支
```

删除分支
```shell
# 删除本地分支
$ git branch -d testing
$ git branch -D testing # 强制删除

# 删除远端分支
$ git push origin --delete serverfix
```

推送远端分支
```shell
$ git push origin serverfix
$ git push origin refs/heads/serverfix:refs/heads/serverfix
```

拉取并合并
```shell
$ git fetch origin
$ git merge origin/serverfix
```

跟踪远程分支
```shell
$ git checkout -b <branch> <remote>/<branch>
$ git checkout --track origin/serverfix
Branch serverfix set up to track remote branch serverfix from origin.
Switched to a new branch 'serverfix'

# 设置不同名字的本地分支和远程分支
$ git checkout -b sf origin/serverfix
```


## 小技巧
### 从每一个提交中移除一个文件
```shell
$ git filter-branch --tree-filter 'rm -f passwords.txt' HEAD
```

### grep搜索

默认情况下 `git grep` 会查找工作目录的文件。 第一种变体是,可以传递 `-n` 或 `--line-number` 选项数来输出 Git 找到的匹配行的行号。
```shell
$ git grep -n gmtime_r
compat/gmtime.c:3:#undef gmtime_r
compat/gmtime.c:8:      return git_gmtime_r(timep, &result);
compat/gmtime.c:11:struct tm *git_gmtime_r(const time_t *timep, struct tm *result)
compat/gmtime.c:16:     ret = gmtime_r(timep, result);
compat/mingw.c:826:struct tm *gmtime_r(const time_t *timep, struct tm *result)
compat/mingw.h:206:struct tm *gmtime_r(const time_t *timep, struct tm *result);
date.c:482:             if (gmtime_r(&now, &now_tm))
date.c:545:             if (gmtime_r(&time, tm)) {
date.c:758:             /* gmtime_r() in match_digit() may have clobbered it */
git-compat-util.h:1138:struct tm *git_gmtime_r(const time_t *, struct tm *);
git-compat-util.h:1140:#define gmtime_r git_gmtime_r
```

### 搜索调用函数
搜索字符串的 `上下文`,那么可以传入 `-p` 或 `--show-function` 选项来显示每一个匹配的字符串所在的方法或函数:
```shell
$ git grep -p gmtime_r *.c
date.c=static int match_multi_number(timestamp_t num, char c, const char *date,
date.c:         if (gmtime_r(&now, &now_tm))
date.c=static int match_digit(const char *date, struct tm *tm, int *offset, int *tm_gmt)
date.c:         if (gmtime_r(&time, tm)) {
date.c=int parse_date_basic(const char *date, timestamp_t *timestamp, int *offset)
date.c:         /* gmtime_r() in match_digit() may have clobbered it */
```
如你所见,date.c 文件中的 match_multi_number 和 match_digit 两个函数都调用了 gmtime_r 例程 (第三个显示的匹配只是注释中的字符串)


### Git 日志搜索
或许你不想知道某一项在 `哪里` ,而是想知道是什么 `时候` 存在或者引入的。 git log 命令有许多强大的工具可以通过提交信息甚至是 diff 的内容来找到某个特定的提交。

例如,如果我们想找到 `ZLIB_BUF_MAX` 常量是什么时候引入的,我们可以使用 -S 选项 (在 Git 中俗称“鹤嘴锄(pickaxe)”选项)来显示新增和删除该字符串的提交。
```shell
$ git log -S ZLIB_BUF_MAX --oneline
e01503b zlib: allow feeding more than 4GB in one go
ef49a7a zlib: zlib can only process 4GB at a time
```
如果我们查看这些提交的 diff,我们可以看到在 ef49a7a 这个提交引入了常量,并且在 e01503b 这个提交中被修改了。

### 工作流

![F0710C53-BD68-4950-ACC8-FACD90A14FF5](./Image/F0710C53-BD68-4950-ACC8-FACD90A14FF5.png)

## 扩展示例

https://git-scm.com/book/zh/v2/%E5%88%86%E5%B8%83%E5%BC%8F-Git-%E5%90%91%E4%B8%80%E4%B8%AA%E9%A1%B9%E7%9B%AE%E8%B4%A1%E7%8C%AE
https://git-scm.com/book/zh/v2/GitHub-%E5%AF%B9%E9%A1%B9%E7%9B%AE%E5%81%9A%E5%87%BA%E8%B4%A1%E7%8C%AE


## 参考
https://git-scm.com/book/zh/v2
