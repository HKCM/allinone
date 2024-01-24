#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: healthcheck 
# 用于检查命令是否存在与系统中
LOG_LEVEL=1
source log.sh

function check_command_usage(){
    yellow='\033[1;33m'
    reset='\033[0m'
    echo -e "${yellow}"
    echo -e "\n=== Check Command Usage ===\n"
    echo -e "Example: check_command <command>\n\n\n"
    echo -e "Will exit as 2 if command not exist:\n"
    echo -e "check_command \"curl\""
    echo -e "${reset}"
}

function check_command(){
    if [ -z $1 ]; then
        check_command_usage
        log_fatal "Please set the command that you want to check"
    fi

    if ! which $1 &> /dev/null; then
        log_fatal "Please install \"$1\" first"
    fi

    log_debug "Command $1 exist in operation system"
    return 0
}