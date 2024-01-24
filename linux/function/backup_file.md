
使用后缀备份文件,备份到用户家目录下的BACKUP,不备份以~开头的临时文件

```shell
#!/bin/bash
set -e
# backup.sh
function makeDir() {
    if [ ! -d "$1" ]; then
        mkdir -p "$1"
    fi
}

read -p " What folder should be backed up: " folder
read -p " What type of files should be backed up: " suffix

BACKUP_DIR=${HOME}/BACKUP
makeDir ${BACKUP_DIR}

find ${folder} -name "*.$suffix" -a ! -name '~*' -exec cp {} ${BACKUP_DIR} \;
echo "Backed up files from ${folder} to ${BACKUP_DIR}"
```


按照下面的方法向脚本发送自动输入
```shell
echo -e "notes\ntxt\n" | ./backup.sh
Backed up files from notes to /BackupDrive/MyName/notes 
```
