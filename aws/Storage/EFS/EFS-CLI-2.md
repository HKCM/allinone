#### 创建文件系统

```shell
aws efs create-file-system \
    --performance-mode generalPurpose \
    --throughput-mode bursting \
    --encrypted \
    --tags Key=Name,Value=my-file-system
```

#### 查看文件系统

```shell
aws efs describe-file-systems \
    --file-system-id fs-c7a0456e
```
在响应中,ValueInIA 显示 IA 存储中的最后一个计量大小。ValueInStandard 显示 Standard 存储中的最后一个计量大小。两者相加等于 Value 显示的整个文件系统的大小。
```shell
{
    "FileSystems": [
        {
            "OwnerId": "123456789012",
            "CreationToken": "console-d7f56c5f-e433-41ca-8307-9d9c0example",
            "FileSystemId": "fs-c7a0456e",
            "FileSystemArn": "arn:aws:elasticfilesystem:us-west-2:123456789012:file-system/fs-48499b4d",
            "CreationTime": 1595286880.0,
            "LifeCycleState": "available",
            "Name": "my-file-system",
            "NumberOfMountTargets": 3,
            "SizeInBytes": {
                "Value": 6144,
                "Timestamp": 1600991437.0,
                "ValueInIA": 0,
                "ValueInStandard": 6144
            },
            "PerformanceMode": "generalPurpose",
            "Encrypted": true,
            "KmsKeyId": "arn:aws:kms:us-west-2:123456789012:key/a59b3472-e62c-42e4-adcf-30d92example",
            "ThroughputMode": "bursting",
            "Tags": [
                {
                    "Key": "Name",
                    "Value": "my-file-system"
                }
            ]
        }
    ]
}
```

#### 删除文件系统
```shell
aws efs delete-file-system \
--file-system-id ID-of-file-system-to-delete
```

#### 创建挂载目标
```shell
aws efs create-mount-target \
--file-system-id file-system-id \
--subnet-id  subnet-id \
--security-group SGID
```

#### 获取挂载目标的描述
当跨VPC挂载时可用区 ID需要与EC2 所在的可用区 ID相匹配方可挂载
```shell
aws efs describe-mount-targets --file-system-id file_system_id
{
    "MountTargets": [
        {
            "OwnerId": "111122223333",
            "MountTargetId": "fsmt-11223344", 
  =====>    "AvailabilityZoneId": "use2-az2",
            "NetworkInterfaceId": "eni-048c09a306023eeec", 
            "AvailabilityZoneName": "us-east-2b", 
            "FileSystemId": "fs-01234567", 
            "LifeCycleState": "available", 
            "SubnetId": "subnet-06eb0da37ee82a64f", 
            "OwnerId": "958322738406", 
  =====>    "IpAddress": "10.0.2.153"
        }, 
...
        {
            "OwnerId": "111122223333",
            "MountTargetId": "fsmt-667788aa", 
            "AvailabilityZoneId": "use2-az3", 
            "NetworkInterfaceId": "eni-0edb579d21ed39261", 
            "AvailabilityZoneName": "us-east-2c", 
            "FileSystemId": "fs-01234567", 
            "LifeCycleState": "available", 
            "SubnetId": "subnet-0ee85556822c441af", 
            "OwnerId": "958322738406", 
            "IpAddress": "10.0.3.107"
        }
    ]
}
```
可用区 ID `use2-az2` 中的挂载目标的 IP 地址为 10.0.2.153。

#### 删除挂载目标
```shell
aws efs delete-mount-target \
--mount-target-id ID-of-mount-target-to-delete
```

#### 创建访问点
```shell
aws efs create-access-point --file-system-id fs-12345678 --posix-user Uid=0,Gid=0,SecondaryGids=1,2 --root-directory Path="/test",CreationInfo="{OwnerUid=0,OwnerGid=0,Permissions=755}"
{
    "ClientToken": "4c27af46-c1f3-40e6-91ba-9fa58ee71234",
    "Tags": [],
    "AccessPointId": "fsap-1234567890abcdefe",
    "AccessPointArn": "arn:aws:elasticfilesystem:us-east-1:1234567890ab:access-point/fsap-1234567890abcdefe",
    "FileSystemId": "fs-12345678",
    "PosixUser": {
        "Uid": 0,
        "Gid": 0,
        "SecondaryGids": [
            1,
            2
        ]
    },
    "RootDirectory": {
        "Path": "/test",
        "CreationInfo": {
            "OwnerUid": 0,
            "OwnerGid": 0,
            "Permissions": "755"
        }
    },
    "OwnerId": "1234567890ab",
    "LifeCycleState": "creating"
}
```

#### 删除访问点
```shell
aws efs delete-access-point --access-point-id fsap-1234567890abcdefe
```

#### 修改生命周期
`TransitionToIA` 可能的值 `AFTER_7_DAYS` `AFTER_14_DAYS` `AFTER_30_DAYS` `AFTER_60_DAYS` `AFTER_90_DAYS`
```shell
aws efs put-lifecycle-configuration \
  --file-system-id File-System-ID \
  --lifecycle-policies TransitionToIA=AFTER_60_DAYS
```

#### 停止生命周期
将 `--lifecycle-policies` 属性留空
```shell
aws efs put-lifecycle-configuration \
  --file-system-id File-System-ID \
  --lifecycle-policies
```

#### 打开自动备份
```shell
aws efs put-backup-policy --file-system-id fs-01234567 \
--backup-policy Status="ENABLED"
```

#### 关闭自动备份
```shell
aws efs put-backup-policy --file-system-id fs-01234567 \
--backup-policy Status="DISABLED"
```

#### 创建文件系统标签
```shell
aws efs tag-resource \
  --resource-id fs-c7a0456e \
  --tags \
      Key=Department,Value=dev \
      Key=Environment,Value="int business"
```

#### 访问点标签
```shell
aws efs tag-resource \
  --resource-id Access-Point-Id \
  --tags \
      Key=Department,Value=dev \
      Key=AccessTeam,Value="Data_Lake"
```

#### 检索与文件系统关联的所有标签
```shell
aws efs list-tags-for-resource \
  --resource-id File-System-Id
```

