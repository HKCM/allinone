#!/usr/bin/bash

# 获取指定instance的IP

aws ec2 describe-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${instance_id} \
    --output text \
    --query 'Reservations[*].Instances[*].PublicIpAddress'

ip_address=$(aws ec2 describe-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${instance_id} \
    --output text \
    --query 'Reservations[*].Instances[*].PublicIpAddress')