#!/bin/env/bash
function usage() {
  echo "Usage:

  ./scripts/admin/deploy_infra.sh \\
    -c <CMR Number> \\
    -s <service> \\
    -b <brand> \\
    -e <environment/stage> \\
    -p <aws_profile> \\

Example:

  $0 -s slack -b rc -e dev -p int-developer
"
  exit 0
}

main() {
    while getopts "c:s:b:e:p:h" opt; do
    case "$opt" in
    c) CMR="$OPTARG" ;;
    s) Service="$OPTARG" ;;
    b) Brand="$OPTARG" ;;
    e) Stage="$OPTARG" ;;
    p) AWSProfile="$OPTARG" ;;
    h) usage ;;
    [?]) usage ;;
    esac
    done

    echo ${CMR:=null}
    echo ${Service:=null}
    echo ${Brand:=null}
    echo ${Stage:=null}
    echo ${AWSProfile:=null}
}

main $@