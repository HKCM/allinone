#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: healthcheck 
LOG_LEVEL=1
source log.sh

function confirm() {
    local prompt_message=$1
    log_debug ${prompt_message}
    if [ -z "${prompt_message}" ];then
        prompt_message="Are You Sure?"
    fi
    prompt_message="${prompt_message} [Y/N]:"
	while true;do
        echo ${prompt_message} 
        read -n1 -r input
        echo 
        case $input in
            Y | y)
                log_debug "Confirmed"
                # Other function
                return 0
                ;;
            N | n)
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