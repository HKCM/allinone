#!/usr/bin/env bash

function makeDir() {
    if [ ! -d "$1" ]; then
　　    mkdir -p "$1"
    fi
}

# 用例
# makeDir /tmp/test