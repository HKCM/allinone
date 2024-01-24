#!/usr/bin/env bash
# keyword: 用户判断 判断执行用户 check user

if [ $UID -eq 0 ]; then
  echo "Do not execute this script with root user"
  exit 99
fi

if [ ${EUID:-${UID}} = 0 ]; then
  echo "Do not execute this script with root user"
  exit 99
fi

TARGET_USER=ubuntu
if [ "$(whoami)" != ${TARGET_USER} ]; then
  echo "Please use ${TARGET_USER} to start service"
  exit 99
fi
