## 配置文件

```shell
# 查看配置文件
aws configure list
```

## EC2

### Instance Type


| Type       | vCPU  | Memory | 原价                  | Spot                       | 平台           |
| ------------ | ------- | -------- | ----------------------- | ---------------------------- | ---------------- |
| t3.micro   | 2vCPU | 1GiB   | $0.0114/h, $0.273/Day | $0.0034/h, $0.081/Day(3折) | 64 位平台      |
| t3.small   | 2vCPU | 2GiB   | $0.0228/h, $0.547/Day | $0.0070/h, $0.168/Day(3折) | 64 位平台      |
| t3.medium  | 2vCPU | 4GiB   | $0.0456/h, $1.094/Day | $0.0137/h, $0.328/Day(3折) | 64 位平台      |
| t3.large   | 2vCPU | 8GiB   | $0.0912/h, $2.188/Day | $0.0277/h, $0.664/Day(3折) | 64 位平台      |
| t4g.micro  | 2vCPU | 1GiB   | $0.0092/h, $0.220/Day | $0.0028/h, $0.067/Day(3折) | 64 位 ARM 平台 |
| t4g.small  | 2vCPU | 2GiB   | $0.0184/h, $0.441/Day | $0.0055/h, $0.132/Day(3折) | 64 位 ARM 平台 |
| t4g.medium | 2vCPU | 4GiB   | $0.0368/h, $0.883/Day | $0.0110/h, $0.264/Day(3折) | 64 位 ARM 平台 |
| t4g.large  | 2vCPU | 8GiB   | $0.0736/h, $1.766/Day | $0.0221/h, $0.519/Day(3折) | 64 位 ARM 平台 |
| m5.large   | 2vCPU | 8GiB   | $0.1070/h, $2.568/Day | $0.0724/h, $1.737/Day(5折) | 64 位平台      |
| m5.xlarge  | 2vCPU | 16GiB  | $0.2140/h, $5.136/Day | $0.1443/h, $3.463/Day(6折) | 64 位平台      |

Instance 类型说明: https://www.amazonaws.cn/ec2/instance-types/
Instance 价格说明: https://aws.amazon.com/cn/ec2/pricing/on-demand/
Instance Spot价格: https://aws.amazon.com/cn/ec2/spot/pricing/

### Image


| ImageID               | Type                    | Platform | Version     | ConnectName |
| ----------------------- | ------------------------- | :--------- | :------------ | ------------- |
| ami-0bf84c42e04519c85 | Amazon Linux 2 AMI      | X86      | Kernel 5.10 | ec2-user    |
| ami-07e30a3659a490be7 | Amazon Linux 2 AMI      | Arm      | Kernel 5.10 | ec2-user    |
| ami-08ca3fed11864d6bb | Ubuntu Server 20.04 LTS | X86      |             | ubuntu      |
| ami-07d8796a2b0f8d29c | Ubuntu Server 18.04 LTS | X86      |             | ubuntu      |
| ami-0ff760d16d9497662 | CentOS Linux 7          | X86      |             | centos      |


### 获取镜像

[describe-images](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/describe-images.html)

```shell
# owners 099720109477 这个是AWS的权威账号
# 查询Amazon Linux 2 镜像ID
aws ec2 describe-images \
    --profile myprofile \
    --region eu-west-1 \
    --owners 099720109477 \
    --filters "Name=name,Values=amzn2*gp2" "Name=virtualization-type,Values=hvm" "Name=root-device-type,Values=ebs" \
    --query "sort_by(Images, &CreationDate)[-5:-1].{id:ImageId,name:Name,date:CreationDate}" \
    --output table

# 查询Ubuntu镜像ID,可以更改版本号查询
aws ec2 describe-images \
    --profile myprofile \
    --region eu-west-1 \
    --owners 099720109477 \
    --filters "Name=name,Values=ubuntu/images/hvm-ssd/ubuntu*18.04*server*" "Name=virtualization-type,Values=hvm" "Name=root-device-type,Values=ebs" \
    --query "sort_by(Images, &CreationDate)[-4:-1].{id:ImageId,name:Name,date:CreationDate,owner:OwnerId}" \
    --output table

# 描述一个镜像
aws ec2 describe-images \
    --profile myprofile \
    --region eu-west-1 \
    --image-ids ami-example11864d6bb \
    --query "Images[-4:-1].{id:ImageId,name:Name,date:CreationDate,owner:OwnerId}"
```

### 创建 Spot 实例

[run-instances](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/run-instances.html)

Instance Spot价格: https://aws.amazon.com/cn/ec2/spot/pricing/

```shell
# run-instance 不要查IP因为刚启动还没分配IP
aws ec2 run-instances \
--profile myprofile \
--region eu-west-1 \
--image-id ami-example52e6c9f9 \
--instance-type t4g.small \
--count 1 \
--subnet-id subnet-exampled0c3652d62 \
--associate-public-ip-address \
--key-name mykey \
--security-group-ids sg-example4e1f8ccdda \
--block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":20,"DeleteOnTermination":true,"VolumeType":"gp2"}}]' \
--instance-market-options 'MarketType=spot,SpotOptions={MaxPrice=0.014,SpotInstanceType=one-time}' \
--tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=local-test},{Key=Team,Value=myprofile},{Key=Environment,Value=myprofile},{Key=owner,Value=test-user}]' \
--query "Instances[*].{InstanceId:InstanceId,Key:KeyName}" 

# 用完不要忘记关机
```

### 等待instance状态ok

```shell
aws ec2 wait instance-status-ok \
    --instance-ids i-1234567890abcdef0
```

### 获取Instance IP

```shell
aws ec2 describe-instances \
--profile int-xmn \
--region eu-west-1 \
--instance-ids i-0d21f0f642625e651 \
--query "Reservations[*].Instances[*].{InstanceId:InstanceId,PublicIP:PublicIpAddress,Name:Tags[?Key=='Name']|[0].Value,Type:InstanceType,Status:State.Name}" \
--output table
```

### SSH

```shell
ssh -i "mykey.pem" centos@ipaddress.eu-west-1.compute.amazonaws.com
ssh -i "mykey.pem" ec2-user@ipaddress.eu-west-1.compute.amazonaws.com
ssh -i "mykey.pem" ubuntu@ipaddress.eu-west-1.compute.amazonaws.com
```

### 获取instance ID

[describe-instances](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/describe-instances.html)

```shell
aws ec2 describe-instances \
--profile myprofile \
--region eu-west-1 \
--filters "Name=tag:Name,Values=local-test" \
--query 'Reservations[*].Instances[*].{I0:InstanceId,I2:lockDeviceMappings[*].Ebs.VolumeId,I1:Placement.AvailabilityZone}' \
--output text
```

### 停止EC2

[stop-instances](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/stop-instances.html)
Spot 实例 不能停止

```shell
aws ec2 stop-instances \
--profile myprofile \
--region eu-west-1 \
--instance-ids ${InstanceId}
```

### 等待停止EC2

[wait](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/wait/index.html)

```shell
aws ec2 wait instance-stopped \
--profile myprofile \
--region eu-west-1 \
--instance-ids ${InstanceId}
```

### 启动EC2

```shell
aws ec2 start-instances \
--profile myprofile \
--region eu-west-1 \
--instance-ids ${InstanceId}
```

### 终止instance

```shell
aws ec2 terminate-instances \
--profile myprofile \
--region eu-west-1 \
--instance-ids i-12345678
```

### 获取自己当前的IP

```shell
curl checkip.amazonaws.com
```

### 获取Security Group

```shell
# ?IsEgress=`false` 只查询入站规则
aws ec2 describe-security-group-rules \
--profile myprofile \
--region eu-west-1 \
--filter Name="group-id",Values="sg-example4e1f8ccdda" \
--query "SecurityGroupRules[? ! IsEgress].{sgr:SecurityGroupRuleId,Protocol:IpProtocol,FromPort:FromPort,ToPort:ToPort,IP:CidrIpv4,Description:Description}" \
--output table
```

### 更改Security Group 单个Rule

更改已有的规则

```shell
aws ec2 modify-security-group-rules \
--profile myprofile \
--region eu-west-1 \
--group-id sg-example4e1f8ccdda \
--security-group-rules SecurityGroupRuleId=sgr-example007ba61131,SecurityGroupRule="{IpProtocol=tcp,FromPort=22,ToPort=22,CidrIpv4=100.100.100.100/32,Description=local-home}"

#--security-group-rules SecurityGroupRuleId=string,SecurityGroupRule="{IpProtocol=string,FromPort=integer,ToPort=integer,CidrIpv4=string,CidrIpv6=string,PrefixListId=string,ReferencedGroupId=string,Description=string}"
```

### 添加新IP

```shell
aws ec2 authorize-security-group-ingress \
--profile myprofile \
--region eu-west-1 \
--group-id sg-example4e1f8ccdda \
--ip-permissions IpProtocol=tcp,FromPort=22,ToPort=22,IpRanges="[{CidrIp=100.100.100.100/32,Description=local-test}]"
```

### 删除Security Group IP

```shell
aws ec2 revoke-security-group-ingress \
--profile myprofile \
--region eu-west-1 \
--group-id sg-example4e1f8ccdda \
--protocol tcp \
--port 22 \
--cidr 100.100.100.100/32
```

## S3

### 只上传特定文件

如果只想上传具有特定扩展名的文件,则需要先排除所有文件,然后重新包含具有特定扩展名的文件。此命令将仅上传以 结尾的文件.jpg:

```shell
aws s3 cp /tmp/foo/ s3://bucket/ --recursive \
--exclude "*" --include "*.jpg"

aws s3 cp /tmp/foo/ s3://bucket/ --recursive \
--exclude "*" --include "*.jpg" --include "*.txt"
```




## CloudFormation

### 获取CloudFormation的Output

```shell
aws cloudformation describe-stacks \
--profile myprofile \
--region eu-west-1 \
--stack-name stackName \
--query 'Stacks[0].Outputs[].{OutputKey:OutputKey,OutputValue:OutputValue}' \
--output table
```

### 清理changeset

将以下内容保存为一个脚本cleanup_changset.sh

脚本运行参数`cleanup_changset.sh -p <aws_profile> -r <region> -s <stack_name>`

```shell
#!/bin/env bash
# cleanup_changeset (Profile, Stack-name, [Region])
# It can clear the failed changeset in the stack, but need to change env/environment.yml, 
# add delete ChangeSet Action "DeleteChangeSet"
# Solve: An error occurred (LimitExceededException) when calling the CreateChangeSet operation: 
#       ChangeSet limit exceeded for stack ...
cleanup_changeset () {
    local profile=$1
    local region=$2
    local stackname=$3
    i=0
    echo "Cleaning up failed change sets"
    changesets=$(aws cloudformation list-change-sets \
        --profile ${profile} \
        --region ${region} \
        --stack-name ${stackname} --query 'Summaries[?Status==`FAILED`].ChangeSetId' --output text)
    echo changesets
    for changeset in $changesets; do
      ((i++))
      echo "${stackname}: deleting change set ${i}: ${changeset:0-36}"
      aws cloudformation delete-change-set \
        --profile ${profile} \
        --region ${region} \
        --stack-name ${stackname} \
        --change-set-name ${changeset}
    done
}

function usage() {
  echo "Usage:
$0 -p <aws_profile> -r <region> -s <stack_name>

Example1:
It will delete all failed changeset 

./Mac/OneClickTools/cleanup_changeset -p int-developer -r ap-southeast-1 -s int-auth-interface-test-pre-infra
"
  exit 0
}

while getopts "p:r:s:" opt; do
  case "$opt" in
  p) profile="$OPTARG" ;;
  r) region="$OPTARG";;
  s) stackname="$OPTARG" ;;
  [?]) usage ;;
  esac
done

echo "profile: ${profile}"
echo "region: ${region}"
echo "stackname: ${stackname}"

if [ -z "${profile}" ] || [ -z "${region}" ] || [ -z "${stackname}" ]; then
    usage
fi

cleanup_changeset ${profile} ${region} ${stackname}
```

### 解密AWS Error Message

创建AWS资源需要Team 和Environment 标签

```shell
aws sts --profile myprofile decode-authorization-message --encoded-message SKHJDGJweajs_as
```

然后去这个网站解析JSON https://www.bejson.com/explore/index_new/
