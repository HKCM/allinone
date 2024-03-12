profile=staging
region=ap-northeast-1
NLB_NAME=test-nlb
NLB_S3=nlb-logs-bucket
S3_PREFIX=test-nlb
security_group_id=sg-223423423423
SubnetId_A=subnet-02fb42730efa63713
SubnetId_C=subnet-0e32bc61b5e7dca03
EIP_A=eipalloc-0085e9271c1f1e027
EIP_C=eipalloc-0c73660fc829b77db

# 创建NLB
LoadBalancerArn=$(aws elbv2 create-load-balancer \
    --region ${region} \
    --profile ${profile} \
    --name ${NLB_NAME} \
    --type network \
    --security-groups ${security_group_id} \
    --subnet-mappings "SubnetId=${SubnetId_A},AllocationId=${EIP_A}" "SubnetId=${SubnetId_C},AllocationId=${EIP_C}" \
    --query "LoadBalancers[0].LoadBalancerArn" \
    --output text) && echo "LoadBalancerArn: ${LoadBalancerArn}"

# 添加s3日志，启用删除保护
aws elbv2 modify-load-balancer-attributes \
    --region ${region} \
    --profile ${profile} \
    --load-balancer-arn ${LoadBalancerArn} \
    --attributes \
    Key=access_logs.s3.enabled,Value=true \
    Key=access_logs.s3.bucket,Value=${NLB_S3} \
    Key=access_logs.s3.prefix,Value=${S3_PREFIX} \
    Key=deletion_protection.enabled,Value=true