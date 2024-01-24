#!/usr/bin/env bash

# author: karl
# version: 1
# date: 2023-12-31
# keyword: shell log 

#################################################
DEBUG=1
INFO=2
WARN=3
ERR=4
DEFAULT=2
#################################################

##################### Usage #####################
# 
# # Put this script together with your script and
# # source it from your script here is an example
#
# #!/usr/bin/env bash
# # Setup log level
# # DEBUG=1
# # INFO=2
# # WARN=3
# # ERR=4
# LOG_LEVEL=1
# # If you want to output your log just setup a
# # log file, it will be fine if don't want it
# 
# # LOG_FILE=test.log  
# source log.sh
# 
# log_debug "This is debug message"
# log_info "This is info message"
# log_warn "This is warnning message"
# log_error "This is err message"
#
#################################################


#################################################
# Get now date string.
#################################################
function now_date() {
    format=$1
    if [[ "${format}" ]]; then
        now=$(date +"${format}")
    else
        now=$(date +"%Y%m%d_%H%M%S")
    fi
    echo "${now}"
}

#################################################
# Basic log function.
# ex: [2021/08/15 19:16:10]
#################################################
function echo_log() {
    now=$(date +"[%Y/%m/%d %H:%M:%S]")
    if [ -n ${LOG_FILE} ];then
        echo -e "\033[1;$1m${now}$2\033[0m" | tee -a ${LOG_FILE}
    else
        echo -e "\033[1;$1m${now}$2\033[0m"
    fi
}

#################################################
# Debug log message. Black
#################################################
function log_debug() {
    if [ ${LOG_LEVEL} -le ${DEBUG} ];then
        echo_log 30 "[Debug] ====> $*"
    fi
}

#################################################
# Information log message. Blue
#################################################
function log_info() {
    if [ ${LOG_LEVEL} -le ${INFO} ];then
        echo_log 34 "[Info] ====> $*"
    fi
}

#################################################
# Warning log message. Yellow
#################################################
function log_warn() {
    if [ ${LOG_LEVEL} -le ${WARN} ];then
        echo_log 33 "[Warning] ====> $*"
    fi
}

#################################################
# Error log message. Red
#################################################
function log_error() {
    echo_log 31 "[Error] ====> $*"
}

#################################################
# Fatal log message. Red and exit the script
#################################################
function log_fatal() {
    echo_log 31 "[Fatal] ====> $*, exit..."
    exit 1
}

#################################################
# Success log message. Green
#################################################
function log_success() {
    echo_log 32 "[Success] ====> $*"
}

# If LOG_LEVEL does not setup, set a default value
LOG_LEVEL=${LOG_LEVEL:=${DEFAULT}} 
# If LOG_LEVEL greater than ERR, set as ERR
if [ ${LOG_LEVEL} -gt ${ERR} ];then 
    LOG_LEVEL=4
fi