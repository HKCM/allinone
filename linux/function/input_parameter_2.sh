#!/usr/bin/env bash
set -e
# keyword: 启动脚本
# 目前这个脚本是无法运行的需要修改Todo让程序能运行起来

APP_NAME=test
TIME_OUT=60

function print_usage() {
  echo "Usage: $0 [start|stop|restart|service|status|show]" 1>&2
  echo "  [start] start app" 1>&2
  echo "  [stop] stop app" 1>&2
  echo "  [restart] stop then  start app" 1>&2
  echo "  [service] create systemd serivce link" 1>&2
  echo "  [status] check app whether it is running" 1>&2
  echo "  [show] show start app command" 1>&2
  exit 1
}

function is_running() {
    APP_PID=$(ps -ef|grep ${APP_NAME}|grep -v grep|awk '{print $2}')
    if [ -z "${APP_PID}" ]; then
        return 1
    else
        return 0
    fi
}

function status() {
    if is_running; then
        echo "${APP_NAME} is running. Pid is ${APP_PID}"
    else
        echo "${APP_NAME} is NOT running."
    fi
}

function start_app() {
    echo "Starting ${APP_NAME} ..."

    if is_running; then
        echo "${APP_NAME} is already running. pid=${APP_PID}."
        # health_check
    else
        # Todo
        # nohup java -jar ${JAVA_OPTS} ${APP_JAR} ${APP_OPTS} > /dev/null 2>&1 &
        echo "${APP_NAME} started with process $! ..."
    fi
}

function stop_app() {
    cd ${WORK_DIR}

    if is_running; then
        echo "${APP_NAME} stoping ... ${APP_PID}"
        kill ${APP_PID}
        wait_for_app_stopped
        echo -e "\n${APP_NAME} stopped ..."
    else
        echo "${APP_NAME} is not running"
    fi
    return 0
}

function restart_app() {
    stop_app
    sleep 10
    start_app
}

function wait_for_app_stopped() {
    local begin=$(date +%s)
    local end
    while is_running
    do
        echo -n "."
        sleep 1;
        end=$(date +%s)
        if [ $((end-begin)) -ge ${TIME_OUT} ];then
            echo -e "\nKill PID:${APP_PID} Timeout...Force kill..."
            kill -9 ${APP_PID}
            return 0
        fi
    done
    return 0
}

function echo_cmd() {
  echo "-------------------------------------------------------------------------------"
  echo TBD
  echo "-------------------------------------------------------------------------------"
}

function main(){
    if [ $# -gt 1 ];then
        print_usage
    fi
    case $1 in
        start)
        start_app
        ;;
        stop)
        stop_app
        ;;
        restart)
        restart_app
        ;;
        status)
        status
        ;;
        show)
        echo_cmd
        ;;
        *)
        print_usage
        ;;
    esac
}

main $@