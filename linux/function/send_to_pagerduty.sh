#!/usr/bin/env bash

# author: karl
# version: 1
# date: 2023-12-31
# 设置pagerduty警报
# 使用-h获取帮助

set -e

TOKEN=XXXXXXXXXX # user token
SERVICE_ID=PXXXXXZ # SERVICE_ID


PARAMETERS_NUM=2
DEFAULT_TIME=60
MIN_TIME=15
MAX_TIME=240
TARGET_TIME=DEFAULT_TIME

set_mw() {
  TARGET_TIME=$1
  # input is a number?
  if [[ ! "${TARGET_TIME}" =~ ^[1-9][0-9]*$ ]]; then
    err_msg "Invaild Input..."
    display_help
    exit 1
  fi
  
  # set max time
  if [ ${TARGET_TIME} -ge ${MAX_TIME} ];then
    err_msg "The time you set is too large, set it to ${MAX_TIME} minutes..."
    TARGET_TIME=${MAX_TIME}
  fi

  # set min time
  if [ ${TARGET_TIME} -lt ${MIN_TIME} ];then
    err_msg "The time you set is too small, set it to ${MIN_TIME} minutes..."
    TARGET_TIME=${MIN_TIME}
  fi

  START_TIME=$(TZ="Japan" date +"%Y-%m-%dT%H:%M:%S%z")
  # Setup 1 hour maintenance window
  #END_TIME=$(TZ="Japan" date -d '+1 hour' +"%Y-%m-%dT%H:%M:%S%z")
  END_TIME=$(TZ="Japan" date -d "+${TARGET_TIME} minutes" +"%Y-%m-%dT%H:%M:%S%z")
  echo "Maintenance Start Time: ${START_TIME}"
  echo "Maintenance End Time: ${END_TIME}"
  echo

  curl --request POST \
    --url https://api.pagerduty.com/maintenance_windows \
    --header 'Accept: application/json' \
    --header "Authorization: Token token=${TOKEN}" \
    --header 'Content-Type: application/json' \
    --header 'From: ' \
    --data "{
    \"maintenance_window\": {
      \"type\": \"maintenance_window\",
      \"start_time\": \"${START_TIME}\",
      \"end_time\": \"${END_TIME}\",
      \"description\": \"Apply release maintenance windows from deploy server\",
      \"services\": [
        {
          \"id\": \"${SERVICE_ID}\",
          \"type\": \"service\"
        }
      ]
    }
  }"
}

err_msg() {
    red='\033[0;31m'
    reset='\033[0m'
    echo
    echo -e "${red}$1${reset}"
    echo
}

display_help() {
    echo "Usage: $0 [option...] " >&2
    echo
    echo "   -h, --help   Show help info"
    echo "   -m           Set Maintenance Window time(minutes,max ${MAX_TIME}, min ${MIN_TIME})"
    echo 
    echo "  Example: $0           Set Maintenance Window as 1 hour(default)"
    echo "  Example: $0 90        Set Maintenance Window as 90 minutes"
    echo "  Example: $0 -m 90     Set Maintenance Window as 90 minutes"
    echo
    exit 1
}

# set maintenance as default time
if [ -z $1 ];then
  set_mw ${DEFAULT_TIME}
  exit 0
fi

# Enter maintenance time directly
if [[ "$1" =~ ^[1-9][0-9]*$ ]]; then
  set_mw $1
  exit 0
fi

# too many parameters. exit
if [ "$#" -gt ${PARAMETERS_NUM} ];then
  err_msg "Too many parameters..."
  display_help  
  exit 0
fi

# Iterate through all parameters
while [[ "$#" -gt 0 ]]; do
  case "$1" in
    -h | --help)
      display_help 
      exit 0
      ;;
    -m)
      time="$2"
      set_mw ${time}
      exit 0
      ;;
    *)
      err_msg "Invaild Input..."
      display_help
      exit 1
      ;;
  esac
done


