profile=staging
region=ap-northeast-1
service_sg=sg-03b8a1a1a19fdd9e9
SG_NAME=test_sg
VPC_ID=vpc-0bd9e232748785adb
INBOUND_IP="1.2.3.4"

# 1. 创建SG
security_group_id=$(aws ec2 create-security-group \
    --region ${region} \
    --profile ${profile} \
    --group-name ${SG_NAME} \
    --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${SG_NAME}}]" \
    --vpc-id ${VPC_ID} \
    --description "PPB NLB SG" \
    --query "GroupId" \
    --output text) && echo "security_group_id: ${security_group_id}"
# 需要删除默认的outbound

# 2. 为SG添加ingress
aws ec2 authorize-security-group-ingress \
    --region ${region} \
    --profile ${profile} \
    --group-id ${security_group_id} \
    --protocol tcp \
    --port 443 \
    --cidr "${INBOUND_IP}/32"
    
# 3. 为SG添加egress
aws ec2 authorize-security-group-egress \
--region ${region} \
--profile ${profile} \
--group-id ${security_group_id} \
--ip-permissions IpProtocol=tcp,FromPort=7001,ToPort=7001,UserIdGroupPairs="[{GroupId=${service_sg}}]"
