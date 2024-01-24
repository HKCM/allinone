#!/usr/bin/env bash
set -e
# author: karl
# version: 1
# date: 2023-12-31
# keyword: check_website healthcheck health check StatusCode 200
# 检查网站状态
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

function check_website_usage(){
    yellow='\033[1;33m'
    reset='\033[0m'
    echo -e "${yellow}"
    echo -e "\n=== Health Check Usage ===\n"
    echo -e "Example: check_website <target_url> [target_res_code]\n\n\n"
    echo -e "Will return success if response code is 2XX or 3XX:\n"
    echo -e "check_website \"localhost:8080\"\n"
    echo -e "will return success only response code is 200:\n"
    echo -e "check_website \"localhost:8080\" \"200\""
    echo -e "${reset}"
}

function check_website() {
    if [ -z $1 ]; then 
        check_website_usage
        log_fatal "Please set your target url"
    fi
    TARGET_URL=$1
    TARGET_CODE=$2
    log_debug "Target url is ${TARGET_URL}"
    log_debug "Target response code is ${TARGET_CODE}"

    # check curl command
    check_command ss
    status_code=$(curl --connect-timeout 3 -I -s -o /dev/null -w "%{http_code}" ${TARGET_URL}) || true
    # status_code=$(curl --connect-timeout 3 -I -s -o /dev/null ${TARGET_URL} |grep 'HTTP' | awk '{print $2}')
    log_debug "Response code is ${status_code}"

    # If the target code exists, perform an exact match
    if [[ ${TARGET_CODE} != "" ]];then 
        if [[ ${status_code} = ${TARGET_CODE} ]];then
            log_debug "Exact match successful"
            return 0
        else
            log_debug "Exact match failed"
            return 1
        fi
    fi

    if echo ${status_code} |grep -qE '^2[0-9][0-9]|^3[0-9][0-9]';then
        log_debug "Fuzzy match successful"
        return 0
    fi

    if [ ${status_code} = "000" ];then
        log_error "Target URL: ${TARGET_URL} is unavailable now..."
        return 1
    fi 

    log_error "Status check failed!!!"
    log_error "Target status code is ${TARGET_STATUS_CODE}"
    log_error "Return status code: ${status_code}"
    return 2
}

# check_command "curl"

# test 1
# check_website

# test 2
check_website localhost

# test 3
#check_website localhost 201