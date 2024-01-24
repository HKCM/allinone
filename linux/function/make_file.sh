#!/usr/bin/env bash

function makefile() {
    if [ ! -f "$1" ]; then
　　    touch "$1"
    fi
}

# 用例
# makefile /tmp/test