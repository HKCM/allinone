#!/usr/bin/env bash

# This is a scale out script for no AutoScalingGroup arch
# It can clone multiple target instances with the same config
# like SecurityGroups, Tags, IAM role, instance type...

set -e
#### Setup variables ####
# AWS profile
profile=
# AWS region
region=ap-northeast-1
# last_instance_name
# 如果第一次扩容数量不能满足需求,则第二次扩容时需要设置last_instance_name
# 这样才能正确设置instance的Name Tag
# 以okj-asset-a01为例
# 第一次扩容2台,新机器能正确命名为a02 a03
# 但是如果扩容之后还需要再次添加3台机器 如果不设置last_instance_name
# 新机器又会从a02开始命名 新机器为a02 a03 a04 则实际AWS上则会变成: 
# a01 a02 a02 a03 a03 a04
# 为避免以上情况,需要将last_instance_name设置为最新的Name Tag
# 在之前的例子中应该last_instance_name 应该设置为okj-asset-a03
# 则第二次扩容机器则为 a04 a05 a06
last_instance_name=None
# Target instance id, will use this instance create AMI
instance_id=
# Target scale out instance number
instance_count=2

##############################################################
#######################    start   ###########################
##############################################################



#################################################
# Basic log function.
# ex: [2021/08/15 19:16:10]
#################################################
function echo_log() {
    now=$(date +"[%Y/%m/%d %H:%M:%S]")
    if [ -n ${LOG_FILE} ];then
        echo -e "\033[1;$1m${now}$2\033[0m" | tee -a ${LOG_FILE}
    else
        echo -e "\033[1;$1m${now}$2\033[0m"
    fi
}

#################################################
# Information log message. Blue
#################################################
function log_info() {
    echo_log 34 "[Info] ====> $*"
}

#################################################
# Error log message. Red
#################################################
function log_error() {
    echo_log 31 "[Error] ====> $*"
}

function time_convert(){
    Seconds=$(($1%60))
    Mins=$(( $1/60 ))
    echo "$2 : $Mins min(s) $Seconds s"
}

function get_target_instance_info(){
    # Get info
    log_info "Get target instance info..."
    ec2info=$(aws ec2 describe-instances \
    --profile ${profile} \
    --region ${region} \
    --instance-ids ${instance_id}  \
    --output json)
    if [ "$?" -ne 0 ]; then log_error "Get instance info failed, Please check instance id. Exit..."; exit 1; fi

    # Set parameter info
    instance_type=$(echo $ec2info | jq -r '.Reservations[0].Instances[0].InstanceType')
    subnet_id=$(echo $ec2info | jq -r '.Reservations[0].Instances[0].SubnetId')
    key_name=$(echo $ec2info | jq -r '.Reservations[0].Instances[0].KeyName')
    sg_ids=($(echo $ec2info | jq -r '.Reservations[0].Instances[0].SecurityGroups[].GroupId'))
    tags=$(echo $ec2info | jq -r '.Reservations[0].Instances[0].Tags' | sed 's/:\ /=/g'| sed 's/"//g')
    iam_instance_profile=$(echo $ec2info | jq -r '.Reservations[0].Instances[0].IamInstanceProfile.Arn' | cut -d'/' -f 2)
    target_instance_name=$(echo $ec2info | jq -r '.Reservations[0].Instances[0].Tags[]|select(.Key == "Name") | .Value')
    date=$(date +%F)
    ami_name="tmp-${target_instance_name}-${date}"
    log_info "Target instance template is ${target_instance_name}..."
    target_instance_name_prefix="${target_instance_name%-*}"
}

function confirm(){
    local prompt_message=$1
    # log_info ${prompt_message}
    if [ -z "${prompt_message}" ];then
        prompt_message="Are You Sure?"
    fi
    prompt_message="${prompt_message} [Y/N]:"
	while true;do
        echo ${prompt_message} 
        read input
        echo 
        case $input in
            Y | y)
                echo "Operation Confirmed"; return 0;;
            N | n)
                log_error "Operation Canceled. Exit..."; exit 1;;
            *)
                log_error "Invalid input...Please enter again..."; exit 2;;
        esac
    done
}

function check_image_info(){
    # Check image
    log_info "Try to get target image ${ami_name} info ...."
    imageId=$(aws ec2 describe-images \
    --profile ${profile} \
    --region ${region} \
    --filters "Name=name,Values=${ami_name}" \
    --output text \
    --query "Images[0].ImageId")

    if [ "$?" -ne 0 ]; then log_error "Get image info failed, exit..."; exit 1; fi

    if [[ "${imageId}" == "None" ]]; then 
        # imageId not exist, the first time run this script
        # Create image by instance id
        
        log_info "No target image. Will using ${target_instance_name} to create new image..."
        imageId=$(aws ec2 create-image \
        --profile ${profile} \
        --region ${region} \
        --instance-id ${instance_id} \
        --name ${ami_name} \
        --description "An AMI for tmp-${image_origin_name} you can delete it if you see it" \
        --no-reboot \
        --output text \
        --query 'ImageId')
        if [ "$?" -ne 0 ]; then log_error "Create image failed, exit..."; exit 1; fi

        # Waiting AMI available
        echo "Waiting AMI ${ami_name}:${imageId} available..."
        aws ec2 wait image-available \
        --region ${region} \
        --profile ${profile} \
        --image-id ${imageId}
        instance_name=${target_instance_name}
    else
        # imageId existed, not the first time run this script
        # So check the last_instance_name
        log_info "Got target ${ami_name} image id..."
        if [[ "${last_instance_name}" == "None" ]]; then
            log_error "Not the first time scale out today"
            log_error "MUST setup last_instance_name variable first...Exit..."
            exit 1
        fi
        last_instance_name_prefix="${last_instance_name%-*}"
        
        # check the name tag prefix
        if [[ "${last_instance_name_prefix}" != "${target_instance_name_prefix}" ]]; then 
            log_error "Variable last_instance_name: ${last_instance_name} does not match target_instance_name: ${target_instance_name_prefix}..."
            log_error "Please check the variable settings...Exit..."
            exit 1
        fi
        instance_name=${last_instance_name}
    fi
}

function scaleout(){
    instance_name_prefix="${instance_name%-*}"  # 获取NameTag前半部分
    instance_name_suffix="${instance_name##*-}"  # 获取尾数
    
    log_info "Will create below instance(s):"
    for ((i = 1; i <= ${instance_count}; i++)); do
        ((tmp_suffix = 10#${instance_name_suffix#*a} + i))
        local tmp_instance_name_tag="${instance_name_prefix}-a$(printf "%02d" $tmp_suffix)"
        echo ${tmp_instance_name_tag}
    done
    confirm "Create the above instances..."
    # scale out
    for ((i = 1; i <= ${instance_count}; i++)); do
        ((new_suffix = 10#${instance_name_suffix#*a} + i))
        new_instance_name_tag="${instance_name_prefix}-a$(printf "%02d" $new_suffix)"
        log_info "Create instance ${new_instance_name_tag}..."
        new_instance_json=$(aws ec2 run-instances \
        --profile ${profile} \
        --region ${region} \
        --image-id ${imageId} \
        --instance-type ${instance_type} \
        --count 1 \
        --subnet-id ${subnet_id} \
        --key-name ${key_name} \
        --security-group-ids ${sg_ids[@]} \
        --tag-specifications "ResourceType=instance,Tags=${tags}" \
        --iam-instance-profile "Name=${iam_instance_profile}" \
        --output json \
        --user-data $'#!/bin/bash
set -e
hostnamectl set-hostname '"${new_instance_name_tag}"'
systemctl restart fluent-bit.service && systemctl restart node-exporter.service
/sbin/runuser -l okj-admin -c "/data/okcoin/okj_app.sh restart" && echo "ok" > /tmp/ok.txt')
        new_instance_id=$(echo $new_instance_json | jq -r '.Instances[].InstanceId')
        if [ "$?" -ne 0 ]; then log_error "Create instance failed, Please check. Exit..."; exit 1; fi
        
        # create-tags for new instances
        log_info "Setup instance ${new_instance_id} Name tag to ${new_instance_name_tag}..."
        aws ec2 create-tags \
        --profile ${profile} \
        --region ${region} \
        --resources "$new_instance_id" \
        --tags Key=Name,Value=${new_instance_name_tag}
    done
    log_info "Waiting instance ${new_instance_id} to instance-status-ok..."
    # 只等待最后一个Instance 状态
    aws ec2 wait instance-status-ok \
    --profile ${profile} \
    --region ${region} \
    --instance-ids ${new_instance_id}
}

function basic_check(){
    if [ -z "$profile" ]; then log_error "Variable profile is empty"; exit 1; fi
    if [ -z "$instance_id" ]; then log_error "Variable instance_id is empty"; exit 1; fi
    if [ ${instance_count} -lt 1 ]; then log_error "Please check the variable instance_count settings...Exit..."; exit 1; fi
}

function main(){
    basic_check
    START_TIME=$(date +%s)
    get_target_instance_info

    confirm "Are you sure to create ${instance_count} ${target_instance_name_prefix} instance(s)"

    check_image_info

    scaleout

    DONE_TIME=$(date +%s)
    # Consuming Time Calculate
    TOTAL_TIME=$(($DONE_TIME-$START_TIME))
    log_info "Scale out Success"
    time_convert $TOTAL_TIME "Total Spend Time"
}

main