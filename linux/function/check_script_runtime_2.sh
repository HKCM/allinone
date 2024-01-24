#!/usr/bin/env bash

# keyword: shell脚本的执行时间 runtime

START_TIME=`date '+%F %T'`
# code......




DONE_TIME=`date '+%F %T'`

# Consuming Time Calculate
TOTAL_TIME=$((`date -d "$DONE_TIME" +%s`-`date -d "$START_TIME" +%s`))

function Time_convert(){
  Seconds=$(($1%60))
  Mins=$(( $1/60 ))
  echo "$2 : $Mins min(s) $Seconds s"
}
echo ""
Time_convert $TOTAL_TIME "Prepare Time "
echo ""
