#!/usr/bin/env bash
# keyword: 启动脚本

if [ ${EUID:-${UID}} = 0 ]; then
  echo "Do not execute this script with root user"
  exit 99
fi

if [ "$(whoami)" != "admin" ]; then
  echo "Please use admin to start app service"
  exit 99
fi

WORK_DIR=/data/work

GC_LOG_FILE=${WORK_DIR}/logs/gc.log
JVM_HEAP_DUMP_FILE=${WORK_DIR}/logs/java.hprof

APP_NAME=@app_name@

# pinpoint agent setup
PP_AGENT_JAR=/data/work/pp-agent/pinpoint-bootstrap.jar
if [ -e ${PP_AGENT_JAR} ];then
    PP_AGENT_ID=`/sbin/ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"`

    if [ $? -ne 0 ];then
        PP_AGENT_ID=`hostname|cut -c1-24`
    fi

    PP_OPTION="-javaagent:${PP_AGENT_JAR} -Dpinpoint.agentId=${PP_AGENT_ID} -Dpinpoint.applicationName=${APP_NAME}"
else
    PP_OPTION=""
fi

# mem info
TOTAL_MEM=`cat /proc/meminfo |grep MemTotal|awk '{print $2}'`
DIRECT_META_MEM=`echo $[${TOTAL_MEM}/1024/6]`
JAVA_HEAP_MEM=`echo $[${TOTAL_MEM}/1024/2]`
# thread info
PARALLEL_GC_THREADS=`cat /proc/cpuinfo|grep processor|wc -l`
CONC_GC_THREADS=`echo $[(3+${PARALLEL_GC_THREADS})/4]`
# mem setup
JVM_MEM_OPTION="-Xms${JAVA_HEAP_MEM}m -Xmx${JAVA_HEAP_MEM}m"
JVM_MEM_OPTION="${JVM_MEM_OPTION} -XX:MetaspaceSize=${DIRECT_META_MEM}m -XX:MaxMetaspaceSize=${DIRECT_META_MEM}m"
JVM_MEM_OPTION="${JVM_MEM_OPTION} -XX:MaxDirectMemorySize=${DIRECT_META_MEM}m"
# java opt setup
JAVA_OPTS="${PP_OPTION} -server"
JAVA_OPTS="${JAVA_OPTS} ${JVM_MEM_OPTION}"
JAVA_OPTS="${JAVA_OPTS} -XX:SurvivorRatio=8 -XX:+UseG1GC"
JAVA_OPTS="${JAVA_OPTS} -XX:MaxGCPauseMillis=200"
JAVA_OPTS="${JAVA_OPTS} -XX:ConcGCThreads=${CONC_GC_THREADS}"
JAVA_OPTS="${JAVA_OPTS} -XX:InitiatingHeapOccupancyPercent=70"
JAVA_OPTS="${JAVA_OPTS} -XX:+ExplicitGCInvokesConcurrent"
JAVA_OPTS="${JAVA_OPTS} -Dsun.rmi.dgc.server.gcInterval=2592000000 -Dsun.rmi.dgc.client.gcInterval=2592000000"
JAVA_OPTS="${JAVA_OPTS} -XX:ParallelGCThreads=${PARALLEL_GC_THREADS}"
JAVA_OPTS="${JAVA_OPTS} -Xloggc:${GC_LOG_FILE}"
JAVA_OPTS="${JAVA_OPTS} -XX:+PrintGCDetails -XX:+PrintGCDateStamps"
JAVA_OPTS="${JAVA_OPTS} -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=${JVM_HEAP_DUMP_FILE}"
JAVA_OPTS="${JAVA_OPTS} -Dsun.net.client.defaultConnectTimeout=10000 -Dsun.net.client.defaultReadTimeout=30000"
JAVA_OPTS="${JAVA_OPTS} -Djava.awt.headless=true -Dfile.encoding=UTF-8"
JAVA_OPTS="${JAVA_OPTS} -Dfastjson.parser.safeMode=true"
JAVA_OPTS="${JAVA_OPTS} -Dproject.name=${APP_NAME}"

APP_OPTS=@app_opts@

APP_JAR=${WORK_DIR}/app/${APP_NAME}_current.jar
APP_PID_FILE=${WORK_DIR}/${APP_NAME}.pid

SERVICE_NAME=app
SERVICE_LINK=/etc/systemd/system/${SERVICE_NAME}.service
SERVICE_FILE=${WORK_DIR}/${SERVICE_NAME}.service
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
  APP_PID=$(ps -ef|grep ${APP_JAR}|grep -v grep|grep -v dhclient|awk '{print $2}')
  if [ -z "${APP_PID}" ]; then
    return 1
  else
    return 0
  fi
}

function status() {
  is_running
  if [ $? -eq "0" ]; then
    echo "${APP_NAME} is running. Pid is ${APP_PID}"
  else
    echo "${APP_NAME} is NOT running."
  fi
}

function start_app() {
  cd ${WORK_DIR}
  echo "Starting ${APP_NAME} ..."

  is_running
  if [ $? -eq "0" ]; then
    echo "${APP_NAME} is already running. pid=${APP_PID} ."
  else
    nohup java -jar ${JAVA_OPTS} ${APP_JAR} ${APP_OPTS} > /dev/null 2>&1 &
    echo "${APP_NAME} started with process $! ..."
  fi
}

function stop_app() {
  cd ${WORK_DIR}

  is_running
  if [ $? -eq "0" ]; then
    echo "${APP_NAME} stoping ... ${APP_PID}"
    kill ${APP_PID}
    wait_for_app_stoped
  else
    echo "${APP_NAME} is not running"
  fi
  return 0
}

restart_app() {
  stop_app
  if [ $? -eq "0" ]; then
    start_app
  else
    echo "restart ${APP_NAME} failed. Do it manually!!!"
  fi
}

function wait_for_app_stoped() {
  local begin=$(date +%s)
  local end
  while is_running
  do
    echo -n "."
    sleep 1;
    end=$(date +%s)
    # if timeout return 1
    if [ $((end-begin)) -ge ${TIME_OUT} ];then
      echo -e "\nKill PID:${APP_PID} Timeout...Force end..."
      kill -9 ${APP_PID}
      return 0
    fi
  done
  echo -e "\n${APP_NAME} stopped ..."
  return 0
}

function register_service() {
  if [ ! -f ${SERVICE_LINK} ]; then
    if [ -f ${SERVICE_FILE} ]; then
      sudo ln -sfn ${SERVICE_FILE} ${SERVICE_LINK}
      sudo systemctl status ${SERVICE_NAME}
    else
      echo "${SERVICE_FILE} does not exists ..."
    fi
  else
    echo "${SERVICE_LINK} already exists ..."
  fi
}

function echo_cmd() {
  echo "-------------------------------------------------------------------------------"
  echo java -jar ${JAVA_OPTS} ${APP_JAR} ${APP_OPTS}
  echo "-------------------------------------------------------------------------------"
}

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
service)
  register_service
  ;;
show)
  echo_cmd
  ;;
*)
  print_usage
  ;;
esac
