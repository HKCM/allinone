

## 描述: 常用EC2 CLI

#### create KeyPair

```shell
aws ec2 create-key-pair --key-name MyKeyPair --query 'KeyMaterial' --output text > MyKeyPair.pem

# for Windows
aws ec2 create-key-pair --key-name MyKeyPair --query 'KeyMaterial' --output text | out-file -encoding ascii -filepath MyKeyPair.pem

# don't forget change the permission of the keypair file
chmod 400 MyKeyPair.pem
```

#### 启动新的EC2

```shell
aws ec2 run-instances \
--image-id ${AMIID} \
--count 1 \
--instance-type t2.micro \
--associate-public-ip-address \
--key-name ${key-pair-name} \
--security-group-ids ${group-ids} \
--subnet-id ${subnetID}
```

#### 确认可用区ID

**不同账户相同AvailableZone的ZoneID是不同的.**

```shell
# 该命令是在EC2机器上运行
aws ec2 describe-availability-zones --zone-name `curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone`
{
    "AvailabilityZones": [
        {
            "State": "available", 
            "ZoneName": "us-east-2b", 
            "Messages": [], 
            "ZoneId": "use2-az2", 
            "RegionName": "us-east-2"
        }
    ]
}

# 通过instance id获取该instance zone name,然后再查找zone id
ZoneName=$(aws ec2 describe-instances \
    --instance-id i-057750d42936e468a \
    --profile ${PROFILE} \
    --region ${REGION} \
    --query "Reservations[0].Instances[0].AvailabilityZone" \
    --output text)
aws ec2 describe-availability-zones --zone-name ${ZoneName} 
```

可用区 ID 在 ZoneId 属性值 `use2-az2` 

#### EC2 tags
```shell
aws ec2 create-tags \
--resources  EC2-instance-ID \
--tags \
    Key=Name,Value=Provide-instance-name \
    Key=Env,Value=dev
```

#### 停止EC2
```shell
aws ec2 stop-instances \
    --instance-ids ${InstanceId}
```

#### 等待停止EC2
```shell
aws ec2 wait instance-stopped \
    --instance-ids ${InstanceId}
```

#### 启动EC2
```shell
aws ec2 start-instances \
    --instance-ids ${InstanceId}
```

#### 终止EC2
```shell
aws ec2 terminate-instances \
--instance-ids instance-id 
```

#### 获取具有指定标签的最新的snapshot id
```shell
function getLatestSnapshot()
{
    SnapshotId=$(aws ec2 describe-snapshots \
    --profile ${PROFILE} \
    --region ${REGION} \
    --filters "Name=tag:Name,Values=${Name}" \
    --query 'reverse(sort_by(Snapshots,&StartTime))[0].[SnapshotId]' \
    --output text)
}
```

#### describe Instance

```shell
# Get instance ID, OldVolumeId and AvailabilityZone

function getInstanceINFO()
{
    InstanceINFO=($(aws ec2 describe-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --filters "Name=tag:Name,Values=${Name}" \
    --query 'Reservations[*].Instances[*].{I0:InstanceId,I2:lockDeviceMappings[*].Ebs.VolumeId,I1:Placement.AvailabilityZone}' \
    --output text))

    InstanceId=${InstanceINFO[0]}
    AvailabilityZone=${InstanceINFO[1]}
    OldVolumeId=${InstanceINFO[2]}
}

# 以下命令列出使用以下任何 AMI 启动的实例:ami-x0123456、ami-y0123456 和 ami-z0123456。
$ aws ec2 describe-instances --filters "Name=image-id,Values=ami-x0123456,ami-y0123456,ami-z0123456"

# 查找没有指定标签的所有实例
aws ec2 describe-instances \
--query 'Reservations[].Instances[?!contains(Tags[].Key, `Name`)][].InstanceId'

```

#### 通过snapshot id创建volume并获取volume id
```shell
NewVolumeId=$(aws ec2 create-volume \
    --profile ${PROFILE} \
    --region ${REGION} \
    --snapshot-id ${SnapshotId} \
    --availability-zone ${AvailabilityZone} \
    --query 'VolumeId' \
    --tag-specifications "ResourceType=volume,Tags=[ \
        {Key=Team,Value=${Team}},{Key=Department,Value=${Department}}, \
        {Key=Name,Value=${Name}},{Key=Environment,Value=${Environment}}]" \
    --output text)

echo "NewVolumeId: ${NewVolumeId}"
```

#### 等待volume可用
```shell
aws ec2 wait volume-available \
    --volume-ids ${NewVolumeId}
```

#### 从EC2上卸载硬盘
```shell
aws ec2 detach-volume \
    --instance-id ${InstanceId} \
    --volume-id ${OldVolumeId}
```

#### 挂载硬盘
```shell
aws ec2 attach-volume \
    --profile ${PROFILE} \
    --region ${REGION} \
    --volume-id ${NewVolumeId} \
    --instance-id ${InstanceId} \
    --device /dev/sda1
```

#### 挂载硬盘后等待硬盘处于使用
```shell
aws ec2 wait volume-in-use \
    --profile ${PROFILE} \
    --region ${REGION} \
    --volume-ids ${NewVolumeId}
```

#### 创建Security Group
```shell
aws ec2 create-security-group \
--group-name efs-walkthrough1-ec2-sg \
--description "Amazon EFS walkthrough 1, SG for EC2 instance" \
--vpc-id vpc-id-in-us-west-2
```

#### 添加入站规则-IP
```shell
aws ec2 authorize-security-group-ingress \
--group-id id of the security group created for EC2 instance \
--protocol tcp \
--port 22 \
--cidr 0.0.0.0/0
```

#### 添加入站规则-安全组
```shell
aws ec2 authorize-security-group-ingress \
--group-id GROUPID \
--protocol tcp \
--port 2049 \
--source-group GROUPID
```

#### 显示 m5.2xlarge 支持 UEFI 和传统 BIOS 启动模式。

```shell
aws ec2 --region us-east-1 describe-instance-types --instance-types m5.2xlarge --query "InstanceTypes[*].SupportedBootModes"

[
    [
        "legacy-bios",
        "uefi"
    ]
]
```

以下示例显示 t2.xlarge 仅支持传统 BIOS。

```shell
aws ec2 --region us-east-1 describe-instance-types --instance-types t2.xlarge --query "InstanceTypes[*].SupportedBootModes"

[
    [
        "legacy-bios"
    ]
]
```

使用公有参数启动实例,默认为Amazon Linux 2 AMI 的最新版本

```shell
aws ec2 run-instances 
    --image-id resolve:ssm:/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2 
    --instance-type m5.xlarge 
    --key-name MyKeyPair
    --tag-specifications 'ResourceType=instance,Tags=[{Key=webserver,Value=production}]' 'ResourceType=volume,Tags=[{Key=cost-center,Value=cc123}]'
```

#### 查找当前 Ubuntu Server 16.04 LTS AMI

```shell
aws ec2 describe-images \
    --owners 099720109477 \
    --filters "Name=name,Values=ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-????????" "Name=state,Values=available" \
    --query "reverse(sort_by(Images, &Name))[:1].ImageId" \
    --output text
```

#### 列出所有公用 AMI,包括您拥有的所有公用 AMI。

```shell
aws ec2 describe-images --executable-users all
```

#### 列出您对其拥有显式启动许可的 AMI。此列表不包括您拥有的任何 AMI。

```shell
aws ec2 describe-images --executable-users self
```

#### 列出 Amazon 拥有的 AMI。

Amazon 的公用 AMI 的拥有者有一个别名,在账户字段中显示为 amazon。这使您可以轻松地从 Amazon 查找 AMI。其他用户不能对其 AMI 使用别名。

```shell
aws ec2 describe-images --owners amazon
```

#### 列出指定 AWS 账户拥有的 AMI。

```shell
aws ec2 describe-images --owners 123456789012
```

要减少显示的 AMI 数量,请使用筛选条件只列出感兴趣的 AMI 类型。例如,使用以下筛选条件可以只显示 EBS 支持的 AMI。

```
--filters "Name=root-device-type,Values=ebs"
```



#### 向指定 AWS 账户授予指定 AMI 的启动许可。

```shell
aws ec2 modify-image-attribute \
    --image-id ami-0abcdef1234567890 \
    --launch-permission "Add=[{UserId=123456789012}]"
```

#### 从指定 AWS 账户中删除指定 AMI 的启动许可:

```shell
aws ec2 modify-image-attribute \
    --image-id ami-0abcdef1234567890 \
    --launch-permission "Remove=[{UserId=123456789012}]"
```

#### 从指定 AMI 中删除所有公用和显式启动许可。

请注意,AMI 的拥有者始终具有启动许可,因此不受该命令影响。

```shell
aws ec2 reset-image-attribute \
    --image-id ami-0abcdef1234567890 \
    --attribute launchPermission
```

#### 使用 S3 存储和还原 AMI

使用 S3 存储和还原 AMI 的权限

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:CompleteMultipartUpload",
                "s3:DeleteObject",
                "s3:GetObject",
                "s3:InitiateMultipartUpload",
                "s3:ListBucket",
                "s3:ListMultipartUploads",
                "s3:ListParts",
                "s3:PutObject",
                "s3:UploadPart",
                "s3:AbortMultipartUpload",
                "ebs:CompleteSnapshot",
                "ebs:GetSnapshotBlock",
                "ebs:ListChangedBlocks",
                "ebs:ListSnapshotBlocks",
                "ebs:PutSnapshotBlock",
                "ebs:StartSnapshot",
                "ec2:CreateStoreImageTask",
                "ec2:DescribeStoreImageTasks",
                "ec2:CreateRestoreImageTask",
                "ec2:GetEbsEncryptionByDefault",
                "ec2:DescribeTags"
            ],
            "Resource": "*"
        }
    ]
}
```

使用 create-store-image-task 命令。指定 AMI 的 ID 以及要在其中存储 AMI 的 S3 存储桶的名称。

```shell
aws ec2 create-store-image-task \
    --image-id ami-1234567890abcdef0 \
    --bucket myamibucket
# 预期输出

{
  "ObjectKey": "ami-1234567890abcdef0.bin"
}
```

描述 AMI 存储任务的进度 (AWS CLI)

使用 describe-store-image-tasks 命令。

```shell
aws ec2 describe-store-image-tasks
# 预期输出

{
  "AmiId": "ami-1234567890abcdef0",
  "Bucket": "myamibucket",
  "ProgressPercentage": 17,
  "S3ObjectKey": "ami-1234567890abcdef0.bin",
  "StoreTaskState": "InProgress",
  "StoreTaskFailureReason": null,
  "TaskStartTime": "2021-01-01T01:01:01.001Z"
}
```

使用 create-restore-image-task 命令。使用来自 S3ObjectKey 输出的 Bucket 和 describe-store-image-tasks 的值,指定 AMI 的对象键以及要将 AMI 复制到的 S3 存储桶的名称。还可以为还原的 AMI 指定名称。名称对该账户在该区域中的 AMI 必须唯一。

注意
还原的 AMI 将获得一个新 AMI ID。

```shell
aws ec2 create-restore-image-task \
    --object-key ami-1234567890abcdef0.bin \
    --bucket myamibucket \
    --name "New AMI Name"
# 预期输出

{
   "ImageId": "ami-0eab20fe36f83e1a8"
}
```

#### 注销AMI

请按照以下步骤使用 AWS CLI 清理由 Amazon EBS 支持的 AMI

1. 注销 AMI,使用 deregister-image 命令注销 AMI:

```shell
aws ec2 deregister-image --image-id ami-12345678
```

2. 删除不再需要的快照,使用 delete-snapshot 命令删除不再需要的快照:

```shell
aws ec2 delete-snapshot --snapshot-id snap-1234567890abcdef0
`
3. 终止实例(可选),如果使用完从 AMI 启动的实例,则可以使用 terminate-instances 命令终止该实例:
```shell
aws ec2 terminate-instances --instance-ids i-12345678
```

#### 使用AWS CLI来确定实例的生命周期。

使用以下描述实例口令:

```shell
aws ec2 describe-instances --instance-ids i-1234567890abcdef0

# 如果实例正在专用主机上运行,那么输出内容包含以下信息:
"Tenancy": "host"

# 如果实例为专用实例,那么输出内容包含以下信息:
"Tenancy": "dedicated"

# 如果实例为 Spot 实例,那么输出内容包含以下信息:
"InstanceLifecycle": "spot"
```

否则,输出不包含 InstanceLifecycle。

#### 获取EC2中临时凭证

```sh
curl http://169.254.169.254/latest/meta-data/iam/security-credentials/
```