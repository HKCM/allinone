#!/bin/bash
# a01和c01机器的subnet不一样 通过subnet来部署到不同的分区中

#这里使用 asset-a01 imageID
imageID=ami-0d4xxx7fxxxxxx

# fiat-a01
aws ec2 run-instances \
--profile prod-mfa \
--region ap-northeast-1 \
--image-id ${imageID} \
--instance-type t3.medium \
--count 1 \
--subnet-id subnet-02d9157ebdf7948a1 \
--key-name engineers-prod \
--security-group-ids sg-0264a5bd27e62881c sg-0482d293bed166303 sg-0ac706db14c0ecd02 sg-0dd9772f9b2ba22cb sg-03124f76c1414f738 sg-022d920c0143d0b9c sg-0ec195c2435c1a426 sg-06573ac52bab11688 sg-0baaacd3a86924ef5 \
--block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":15,"DeleteOnTermination":true,"VolumeType":"gp3"}},{"DeviceName":"/dev/xvdb","Ebs":{"VolumeSize":35,"DeleteOnTermination":true,"VolumeType":"gp3"}}]' \
--tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=fiat-a01},{Key=log,Value=on},{Key=metrics,Value=on},{Key=profile,Value=prod},{Key=deploy,Value=fiat},{Key=type,Value=service},{Key=backup,Value=weekly},{Key=auto-patch,Value=Prod-AL2},{Key=monitor,Value=off},{Key=tmp-log4j2,Value=on}]' \
--iam-instance-profile Name=fiat-prod \
--disable-api-stop \
--disable-api-termination

# fiat-c01
aws ec2 run-instances \
--profile prod-mfa \
--region ap-northeast-1 \
--image-id ${imageID} \
--instance-type t3.medium \
--count 1 \
--subnet-id subnet-0d4767c3b60262e04 \
--key-name engineers-prod \
--security-group-ids sg-0264a5bd27e62881c sg-0482d293bed166303 sg-0ac706db14c0ecd02 sg-0dd9772f9b2ba22cb sg-03124f76c1414f738 sg-022d920c0143d0b9c sg-0ec195c2435c1a426 sg-06573ac52bab11688 sg-0baaacd3a86924ef5 \
--block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":15,"DeleteOnTermination":true,"VolumeType":"gp3"}},{"DeviceName":"/dev/xvdb","Ebs":{"VolumeSize":35,"DeleteOnTermination":true,"VolumeType":"gp3"}}]' \
--tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=fiat-c01},{Key=log,Value=on},{Key=metrics,Value=on},{Key=profile,Value=prod},{Key=deploy,Value=fiat},{Key=type,Value=service},{Key=backup,Value=weekly},{Key=auto-patch,Value=Prod-AL2},{Key=monitor,Value=off},{Key=tmp-log4j2,Value=on}]' \
--iam-instance-profile Name=fiat-prod \
--disable-api-stop \
--disable-api-termination

# fiat-job-a01
aws ec2 run-instances \
--profile prod-mfa \
--region ap-northeast-1 \
--image-id ${imageID} \
--instance-type t3.medium \
--count 1 \
--subnet-id subnet-02d9157ebdf7948a1 \
--key-name engineers-prod \
--security-group-ids sg-0264a5bd27e62881c sg-0482d293bed166303 sg-0ac706db14c0ecd02 sg-0dd9772f9b2ba22cb sg-03124f76c1414f738 sg-022d920c0143d0b9c sg-0ec195c2435c1a426 sg-06573ac52bab11688 sg-0baaacd3a86924ef5 \
--block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":15,"DeleteOnTermination":true,"VolumeType":"gp3"}},{"DeviceName":"/dev/xvdb","Ebs":{"VolumeSize":35,"DeleteOnTermination":true,"VolumeType":"gp3"}}]' \
--tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=fiat-job-a01},{Key=log,Value=on},{Key=metrics,Value=on},{Key=profile,Value=prod},{Key=deploy,Value=fiat-job},{Key=type,Value=service},{Key=backup,Value=weekly},{Key=auto-patch,Value=Prod-AL2},{Key=monitor,Value=off}]' \
--iam-instance-profile Name=fiat-job-prod \
--disable-api-stop \
--disable-api-termination

# 为ALB创建TargetGroup
tgName=prod-fiat-tg
TargetGroupArn=$(aws elbv2 create-target-group \
--profile prod-mfa \
--region ap-northeast-1 \
--name ${tgName} \
--protocol HTTP \
--port 7001 \
--target-type instance \
--vpc-id vpc-08c8b3647d8c41669 \
--query "TargetGroups[*].{TargetGroupArn:TargetGroupArn}")
echo ${TargetGroupArn}

将机器手动添加到TargetGroup

# 创建ALB 目前允许1a 和 1c的public subnet
albName=prod-gmo-alb
LoadBalancerArn=$(aws elbv2 create-load-balancer \
--profile prod-mfa \
--name ${albName} \
--subnets subnet-02df8cc8ae17492f3 subnet-054621dcc926777c4 \
--query "LoadBalancers[*].{LoadBalancerArn:LoadBalancerArn}")
echo ${LoadBalancerArn}

# 创建listener
# CertificateArn 是 *.****.jp
aws elbv2 create-listener \
--load-balancer-arn ${LoadBalancerArn} \
--protocol HTTPS \
--port 443 \
--certificates CertificateArn=arn:aws:acm:ap-northeast-1:123456789012:certificate/xxxx7e38-xxxx-4384-a509-e67097e0xxxx \
--ssl-policy ELBSecurityPolicy-TLS-1-2-2017-01 \
--default-actions Type=forward,TargetGroupArn=${TargetGroupArn}
