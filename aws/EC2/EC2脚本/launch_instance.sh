#!/usr/bin/env bash

set -e

function usage() {
  echo "Usage:

$0 -p <aws_profile> 

Example:

$0 -p int-developer
"
  exit 0
}

while getopts "p:h" opt; do
  case "$opt" in
  p) PROFILE="$OPTARG" ;;
  h) usage ;;
  [?]) usage ;;
  esac
done

if [ ! -n "${PROFILE}" ]; then
    usage
fi

REGION=XX
AMIID=XX
INSTANCETYPE=XX
KEYPAIR=XX
SUBNETID=XX
SECURITYGROUP=XX
USERNAME=$USER
CREATEDATE=`date +%Y-%m-%d`

user_ip=$(curl http://checkip.amazonaws.com)

echo "Your current IP: ${user_ip}, it will add to security group automatically"

echo "Start launch sftp server now, it will takes a few minutes..."

aws ec2 authorize-security-group-ingress \
    --profile ${PROFILE} \
    --region ${REGION} \
    --group-id ${SECURITYGROUP} \
    --ip-permissions IpProtocol=tcp,FromPort=22,ToPort=2,IpRanges="[{CidrIp=${user_ip}/32,Description=${USER}}]"

instance_id=$(aws ec2 run-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --image-id ${AMIID} \
    --instance-type ${INSTANCETYPE} \
    --key-name ${KEYPAIR} \
    --subnet-id ${SUBNETID} \
    --security-group-ids ${SECURITYGROUP} \
    --associate-public-ip-address \
    --tag-specifications \
        "ResourceType=instance,Tags=[ \
        {Key=Name,Value=${USERNAME}}, \
        {Key=Department,Value=XX}, \
        {Key=CreateDate,Value=${CREATEDATE}}, \
        {Key=Environment,Value=XX}, \
        {Key=Team,Value=XX}]" \
        "ResourceType=volume,Tags=[ \
        {Key=Name,Value=${USERNAME}}, \
        {Key=Department,Value=XX}, \
        {Key=Environment,Value=XX}, \
        {Key=Team,Value=XX}]" \
    --output text \
    --query 'Instances[*].InstanceId')


aws ec2 wait instance-status-ok \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${instance_id}

aws ec2 authorize-security-group-ingress \
    --profile ${PROFILE} \
    --region ${REGION} \
    --group-id ${SECURITYGROUP} \
    --ip-permissions IpProtocol=tcp,FromPort=22,ToPort=22,IpRanges="[{CidrIp=${user_ip}/32,Description=${USER}}]"

ip_address=$(aws ec2 describe-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${instance_id} \
    --output text \
    --query 'Reservations[*].Instances[*].PublicIpAddress')

cat>INFO.txt<<EOF

****************************************************
 ${CREATEDATE}: launch successed
 
 Instance ID: ${instance_id}
 
 Public IP: ${ip_address}
 
 The SSH Private Key: ${KEYPAIR}
 
 Usage: ssh -i ${KEYPAIR} subuntu@${ip_address}
 
 You can find it in this repo

 This info will save in SINFO.txt

EOF

cat SFTP_INFO.txt