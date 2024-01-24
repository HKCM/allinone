#!/usr/bin/env bash

# Example1:
# It will start/stop(Only allowed start/stop) default SFTP instance

# ./Mac/OneClickTools/default_sftp.sh -p int-developer -a start

# Example2:
# It will start/stop(Only allowed start/stop) the instance you specify:

# ./Mac/OneClickTools/default_sftp.sh -p int-developer -a stop -i i-1234567890abcdef0 -d "xxx home IP" 

REGION=eu-west-1
CUR_DIR=$(cd `dirname $0` && pwd -P)
default_instance_id="i-1234567890"

set -e

function usage() {
  echo "Usage:

$0 -p <aws_profile> -a <Action> [-i instanceID] 

Example1:

It will startOnly allowed start/stop) default SFTP instance(${REGION}: ${default_instance_id}) 

    $0 -p <aws_profile> -a <Action>

    $0 -p int-developer -a start

Example2:

It will stop(Only allowed start/stop) default SFTP instance(${REGION}: ${default_instance_id}) and delete sftp data(/var/sftp/uploads/*)

    $0 -p <aws_profile> -a <Action> [ -d Delete ]

    $0 -p int-developer -a stop -d

Example3:

It will start/stop(Only allowed start/stop) the instance you specify:

    $0 -p <aws_profile> -a <Action> [-i instanceID] 

    $0 -p int-developer -a stop -i i-1234567890abcdef0
"
  exit 0
}

while getopts "dp:a:i:h" opt; do
  case "$opt" in
  p) PROFILE="$OPTARG" ;;
  a) ACTION="$OPTARG" ;;
  i) INSTANCEID="$OPTARG" ;;
  d) DELETE="true" ;;
  h) usage ;;
  [?]) usage ;;
  esac
done

if [[ -z "${PROFILE}" ]] || [[ -z "${ACTION}" ]]; then
    usage
fi



if [[ -z "${INSTANCEID}" ]] ; then
    INSTANCEID=${default_instance_id}
fi

if [[ ${ACTION} = "start" ]]; then
    echo "Starting Instance: ${INSTANCEID}....."
    aws ec2 start-instances \
        --profile ${PROFILE} \
        --region ${REGION} \
        --instance-ids ${INSTANCEID} > /dev/null 2>&1

    aws ec2 wait instance-status-ok \
        --profile ${PROFILE} \
        --region ${REGION} \
        --instance-ids ${INSTANCEID}

    ip_address=$(aws ec2 describe-instances \
        --profile ${PROFILE} \
        --region ${REGION} \
        --instance-ids ${INSTANCEID} \
        --output text \
        --query 'Reservations[*].Instances[*].PublicIpAddress')
    
    echo "
    Instance: ${INSTANCEID} started. Public IP: ${ip_address}
    SFTP User Name: xxxxxx. SFTP User Password: xxxxx
    Usage: sftp xxxxxx@${ip_address}
    The SSH Private Key: 'xxxxxx.pem'
    SFTP Upload Folder: /var/sftp/uploads
    "
fi

if [[ ${ACTION} = "stop" ]]; then
    if [[ ${DELETE} = "true" ]]; then
        ip_address=$(aws ec2 describe-instances \
            --profile ${PROFILE} \
            --region ${REGION} \
            --instance-ids ${INSTANCEID} \
            --output text \
            --query 'Reservations[*].Instances[*].PublicIpAddress')
        echo "Instance IP: ${ip_address}"
        ssh -i ${CUR_DIR}/xxxxxxx.pem ubuntu@${ip_address} "sudo rm -rf /var/sftp/uploads/*"
    fi
    echo "Stopping Instance: ${INSTANCEID}....."
    aws ec2 stop-instances \
        --profile ${PROFILE} \
        --region ${REGION} \
        --instance-ids ${INSTANCEID} > /dev/null 2>&1
    
    aws ec2 wait instance-stopped \
        --profile ${PROFILE} \
        --region ${REGION} \
        --instance-ids ${INSTANCEID}
    
    echo "Done"
fi

if [[ ${ACTION} != "start" ]] && [[ ${ACTION} != "stop" ]]; then
    echo "ACTION is wrong"
    usage
fi


CloudFrontID=$(aws --profile int-xmn --region us-east-1 cloudformation describe-stacks --stack-name lti-support-rc-development-pre-infra --query "Stacks[*].Outputs[?OutputKey=='CloudFrontID'].OutputValue" --output text)
aws --profile int-xmn --region us-east-1 cloudfront create-invalidation --distribution-id E3R8UCY393MU65 --paths "/*"