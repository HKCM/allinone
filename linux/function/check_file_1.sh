#!/usr/bin/env bash
check_file() {
  if [[ ! -r ${1} ]]; then
    echo "can not find ${1}"
    #exit 1
  fi
}

# 用例
# check_file /tmp/test