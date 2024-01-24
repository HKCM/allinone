```bash
# 设置变量
region=ap-northeast-1
profile=staging
service_asset_sg=sg-xxxxa1a1a19fdxxxx
NLB_NAME=me-bank-nlb
SG_NAME=elb_bank_nlb_sg
SubnetId_A=subnet-xxxx42730efa6xxxx
SubnetId_C=subnet-xxxxbc61b5e7dxxxx
EIP_A=eipalloc-xxxxe9271c1fxxxx
EIP_C=eipalloc-xxxx660fc829bxxxx

VPC_ID=vpc-xxxxe23274878xxxx
SSL_Policy=ELBSecurityPolicy-TLS-1-2-Ext-2018-06
NLB_S3=com-nlb-logs-me
S3_PREFIX=jnb-nlb
INBOUND_IP=1111.111.111.1111
CertificateArn=arn:aws:iam::123456789012:server-certificate/jnb-ssl-prod.test2023
TargetGroupArn=arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/me-asset-tg/1da1bb56525beb38
```

```bash
# get CertificateArn
aws iam list-server-certificates \
    --region ${region} \
    --profile ${profile}

# 1. 创建SG
security_group_id=$(aws ec2 create-security-group \
    --region ${region} \
    --profile ${profile} \
    --group-name ${SG_NAME} \
    --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${SG_NAME}}]" \
    --vpc-id ${VPCID} \
    --query "GroupId" \
    --output text) && echo "security_group_id: ${security_group_id}"

# 2. 为SG添加ingress
aws ec2 authorize-security-group-ingress \
    --region ${region} \
    --profile ${profile} \
    --group-id ${security_group_id} \
    --protocol tcp \
    --port 443 \
    --cidr "${INBOUND_IP}/32"

# 3. 创建NLB
LoadBalancerArn=$(aws elbv2 create-load-balancer \
    --region ${region} \
    --profile ${profile} \
    --name ${NLB_NAME} \
    --type network \
    --security-groups ${security_group_id} \
    --subnet-mappings "SubnetId=${SubnetId_A},AllocationId=${EIP_A}" "SubnetId=${SubnetId_C},AllocationId=${EIP_C}" \
    --query "LoadBalancers[0].LoadBalancerArn" \
    --output text) && echo "LoadBalancerArn: ${LoadBalancerArn}"

# 添加s3日志,启用删除保护
aws elbv2 modify-load-balancer-attributes \
    --region ${region} \
    --profile ${profile} \
    --load-balancer-arn ${LoadBalancerArn} \
    --attributes \
    Key=access_logs.s3.enabled,Value=true \
    Key=access_logs.s3.bucket,Value=${NLB_S3} \
    Key=access_logs.s3.prefix,Value=${S3_PREFIX} \
    Key=deletion_protection.enabled,Value=true

# 4. 创建target-group(Option)
TargetGroupArn=$(aws elbv2 create-target-group \
    --region ${region} \
    --profile ${profile} \
    --name ${TARGETGROUP_NAME} \
    --protocol TCP \
    --port 7001 \
    --target-type instance \
    --vpc-id ${VPCID} \
    --query "TargetGroups[0].TargetGroupArn" \
    --output text) && echo "TargetGroupArn: ${TargetGroupArn}"

# 5. 将机器注册到Target group(Option)
aws elbv2 register-targets \
    --region ${region} \
    --profile ${profile} \
    --target-group-arn ${TargetGroupArn} \
    --targets Id=${InstanceID1} Id=${InstanceID2}
# 等待成功

# 6. 创建listener
ListenerArn=$(aws elbv2 create-listener \
    --region ${region} \
    --profile ${profile} \
    --load-balancer-arn ${LoadBalancerArn} \
    --protocol TLS \
    --port 443 \
    --certificates CertificateArn=${CertificateArn} \
    --ssl-policy ${SSL_Policy} \
    --default-actions Type=forward,TargetGroupArn=${TargetGroupArn} \
    --query "Listeners[0].ListenerArn" \
    --output text) && echo "ListenerArn: ${ListenerArn}"
```



## RollBack Plan
将创建步骤反向执行

```bash
# 1. 删除listener
aws elbv2 delete-listener \
    --region ${region} \
    --profile ${profile} \
    --listener-arn ${ListenerArn}

# 2. 注销  target-group中的target
aws elbv2 deregister-targets \
    --region ${region} \
    --profile ${profile} \
    --target-group-arn ${TargetGroupArn} \
    --targets Id=${InstanceID1}
    
    
# 3. 删除target-group
aws elbv2 delete-target-group \
    --region ${region} \
    --profile ${profile} \
    --target-group-arn ${TargetGroupArn}

# 4. 删除NLB
# 需要先移除删除保护 deletion protection
aws elbv2 delete-load-balancer \
    --region ${region} \
    --profile ${profile} \
    --load-balancer-arn ${LoadBalancerArn}

# 删除SG 需要等待NLB删除后
aws ec2 delete-security-group \
    --region ${region} \
    --profile ${profile} \
    --group-id ${security_group_id}
```
