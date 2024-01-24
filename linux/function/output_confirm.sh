#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: 用户确认 确认操作 
LOG_LEVEL=1
source log.sh

function confirm() {
    local prompt_message=$1
    log_debug ${prompt_message}
    if [ -z "${prompt_message}" ];then
        prompt_message="Are You Sure?"
    fi
    prompt_message="${prompt_message} [YES/n]:"
	while true;do
        echo ${prompt_message} 
        read input

        case $input in
            [yY][eE][sS])
                log_debug "Confirmed"
                # Other function
                return 0
                ;;
            [nN][oO]|[nN])
                log_debug "Canceled"
                return 1
                ;;
            *)
                log_error "Invalid input...Please enter again..."
                # exit 2
                ;;
        esac
    done
}

confirm