#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: 备份文件

LOG_LEVEL=2
LOG_FILE=/test.log
source log.sh

SOURCE=/source
TARGET=/target
log_info "START Backup"
function main() {
    if [ ! -d ${TARGET} ];then
        mkdir -p ${TARGET}
        log_info "Folder ${TARGET} not exist. Create ${TARGET}"
    fi
    cd ${SOURCE}
    for dir in $(ls);do
        # Seach files in the current directory that have been changed within 15 minutes
        # but not the directory itself
        # 备份修改时间大于1天但是小于2天的数据
        for dir2 in $(find ${dir} -maxdepth 1 -type d -mtime -2 -mtime +1 ! -path ${dir}); do
        # for dir2 in $(find ${dir} -maxdepth 1 -type d -mmin -12 -mmin +1 ! -path ${dir}); do
            rsync -aR ${dir2}/ ${TARGET}/
            if [ $? -eq 0 ];then
                rm -rf ${dir2}
                log_success "Move ${SOURCE}/${dir2} to ${TARGET}/${dir2}"
                ln -s ${TARGET}/${dir2} ${SOURCE}/${dir2} && \
                log_success "Create soft links ${SOURCE}/${dir2}"
            else
                log_fatal "Move ${SOURCE}/${dir2} to ${TARGET}/${dir2} failed"
            fi
        done
    done
}

main