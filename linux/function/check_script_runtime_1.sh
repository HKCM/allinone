#!/usr/bin/env bash

# keyword: shell脚本的执行时间 runtime

START_TIME=$(date +%s)
# code......



DONE_TIME=$(date +%s)
# Consuming Time Calculate
TOTAL_TIME=$(($DONE_TIME-$START_TIME))

function Time_convert(){
    Seconds=$(($1%60))
    Mins=$(( $1/60 ))
    write_log "$2 : $Mins min(s) $Seconds s"
}

Time_convert $TOTAL_TIME "Deployment Time "