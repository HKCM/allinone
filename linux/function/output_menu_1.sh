#!/usr/bin/env bash

# author: karl
# version: 1
# date: 2023-12-31
# keyword: 菜单函数 shell menu 

function diskspace {
  clear
  df -k
}

function whoseon {
  clear
  who
}

function memusage {
  clear
  cat /proc/meminfo
}

function menu_bar {
  clear
  echo
  echo -e "\t\t\tSys Admin Menu\n"
  echo -e "\t1. Display disk space"
  echo -e "\t2. Display logged on users" 
  echo -e "\t3. Display memory usage" 
  echo -e "\t0. Exit program\n\n"
  echo -en "\t\tEnter option: "
  read -n 1 option
}

function menu {
  while true; do
    menu_bar
    case $option in
    0) break ;;
    1) diskspace ;;
    2) whoseon ;;
    3) memusage ;;
    *) 
        clear
        echo "Sorry, wrong selection" ;;
    esac
    echo -en "\n\n\t\t\tHit any key to continue"
    read -n 1 line
  done
  clear
}

menu