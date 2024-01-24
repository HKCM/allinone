#!/usr/bin/bash

GROUPID=XXXXX
REGION=eu-west-1
PROFILE=XXXX



aws ec2 describe-security-groups \
    --profile ${PROFILE} \
    --region ${REGION} \
    --group-ids ${GROUPID} \
    --filters Name=ip-permission.from-port,Values=22 \
    --query 'SecurityGroups[*].IpPermissions[1].IpRanges[*].CidrIp' \
    --output text