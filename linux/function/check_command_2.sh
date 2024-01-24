#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: healthcheck 
# 用于检查命令是否存在与系统中
err() {
    printf 'Script-init: %s\n' "$1" >&2
    exit 1
}

need_cmd() {
    if ! check_cmd "$1"; then
        err "need '$1' (command not found)"
    fi
}

check_cmd() {
    command -v "$1" > /dev/null 2>&1
}