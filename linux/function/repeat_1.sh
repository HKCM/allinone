#!/usr/bin/env bash
function repeat() { 
    while :; do
        $@ && return;
        sleep 30; 
    done 
}

# 用例 每隔30秒去下载文件成功则退出
#repeat wget -c http://www.example.com/software-0.1.tar.gz