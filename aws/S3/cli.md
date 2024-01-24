# CLI文档

https://docs.aws.amazon.com/cli/latest/reference/index.html#cli-aws


## 基本操作

桶操作
```bash
aws s3 mb s3://bucket-name --region us-west-1 # 创建桶

aws s3 rm s3://bucket-name --recursive # 删除所有桶内对象

aws s3 rb s3://bucket-name --force # 删除所有桶内对象并删除桶
```

同步文件并设置metadata
```bash
aws s3 sync /path s3://yourbucket/ --delete –recursive --cache-control max-age=60
```

### 日志存储桶

#### 通过更新存储桶策略向日志记录服务主体授予权限
```shell
aws s3api put-bucket-policy --bucket awsexamplebucket1-logs --policy file://policy.json
```

Policy.json
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "S3ServerAccessLogsPolicy",
            "Effect": "Allow",
            "Principal": {
                "Service": "logging.s3.amazonaws.com"
            },
            "Action": [
                "s3:PutObject"
            ],
            "Resource": "arn:aws:s3:::awsexamplebucket1-logs/*",
            "Condition": {
                "ArnLike": {
                    "aws:SourceArn": "arn:aws:s3:::SOURCE-BUCKET-NAME"
                },
                "StringEquals": {
                    "aws:SourceAccount": "SOURCE-ACCOUNT-ID"
                }
            }
        }
    ]
}	
```

#### 通过更新存储桶ACL向日志记录服务主体授予权限

```shell
aws s3api put-bucket-acl --bucket awsexamplebucket1-logs  --grant-write URI=http://acs.amazonaws.com/groups/s3/LogDelivery --grant-read-acp URI=http://acs.amazonaws.com/groups/s3/LogDelivery 
                                    
```



#### 授予日志写入权限

```shell
aws s3api put-bucket-acl --acl log-delivery-write --bucket my-bucket 
```


#### 上传文件时给予所有人读权限

```shell
aws s3 cp file s3://my-bucket/file --acl public-read
```

#### 上传文件进行加密

```shell
aws --profile hkc s3api put-object --bucket bucket-name --key encryption-object.png --body ~/Desktop/teams.png --server-side-encryption AES256

aws --profile hkc s3api put-object --bucket bucket-name --key encryption-object2.png --body ~/Desktop/teams.png --expires "Thu, 09 Jan 2020 06:40:00 GMT"
```

#### 创建具有默认一小时生存期的预签名URL,该URL链接到S3存储桶中的对象 

```shell
aws s3 presign s3://awsexamplebucket/test2.txt

aws s3 presign s3://awsexamplebucket/test2.txt --expires-in 604800 # 7天
```

#### 存储桶启用加密

当存储桶本身启用服务端加密后,上传到桶中的资源会自动加密
```shell
aws s3api put-bucket-encryption --bucket my-bucket --server-side-encryption-configuration '{"Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]}'
```
