#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: download 
LOG_LEVEL=1
source log.sh

function check_command_usage(){
    echo -e "\n=== Check Command Usage ===\n"
    echo -e "Example: check_command <command>\n\n\n"
    echo -e "Will exit as 2 if command not exist:\n"
    echo -e "check_command \"curl\""
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

function download_usage(){
    yellow='\033[1;33m'
    reset='\033[0m'
    echo -e "${yellow}"
    echo -e "\n=== Download Usage ===\n"
    echo -e "Example: download <url> [des]\n\n\n"
    echo -e "Will download on current directory:\n"
    echo -e "download \"https://static.vecteezy.com/system/resources/previews/009/273/280/large_2x/concept-of-loneliness-and-disappointment-in-love-sad-man-sitting-element-of-the-picture-is-decorated-by-nasa-free-photo.jpg\""
    echo -e "${reset}"
}

function download(){
    if [ -z $1 ]; then
        download_usage
        log_fatal "Please set the command that you want to check"
    fi
    download_url=$1

    obj_name=$(echo "$download_url" | awk -F'/' '{print $NF}')
    if [ -f ${obj_name} ]; then
        log_warn "$(pwd)/${obj_name} already existed..."
        return
    fi

    log_info "Start download: $download_url"
    check_command curl
    curl -L -O "$download_url"
    if [ $? -eq 0 ];then
        log_success "Download finished, Object is $(pwd)/${obj_name}"
    else
        log_fatal "Download failed..."
    fi

}

# download "https://static.vecteezy.com/system/resources/previews/009/273/280/large_2x/concept-of-loneliness-and-disappointment-in-love-sad-man-sitting-element-of-the-picture-is-decorated-by-nasa-free-photo.jpg"
