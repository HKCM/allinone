[toc]





演示一些常见命令

### DynamoDB

#### 创建table

```shell
$ aws dynamodb create-table \
    --table-name MusicCollection \
    --attribute-definitions AttributeName=Artist,AttributeType=S AttributeName=SongTitle,AttributeType=S \
    --key-schema AttributeName=Artist,KeyType=HASH AttributeName=SongTitle,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
```

#### 插入项目

```shell
$ aws dynamodb put-item \
    --table-name MusicCollection \
    --item '{
        "Artist": {"S": "No One You Know"},
        "SongTitle": {"S": "Call Me Today"} ,
        "AlbumTitle": {"S": "Somewhat Famous"} 
      }' \
    --return-consumed-capacity TOTAL
{
    "ConsumedCapacity": {
        "CapacityUnits": 1.0,
        "TableName": "MusicCollection"
    }
}

$ aws dynamodb put-item \
    --table-name MusicCollection \
    --item '{ 
        "Artist": {"S": "Acme Band"}, 
        "SongTitle": {"S": "Happy Day"} , 
        "AlbumTitle": {"S": "Songs About Life"} 
      }' \
    --return-consumed-capacity TOTAL

{
    "ConsumedCapacity": {
        "CapacityUnits": 1.0,
        "TableName": "MusicCollection"
    }
}
```

#### 使用json文件

json文件`expression-attributes.json`的内容
```json
{
  ":v1": {"S": "No One You Know"},
  ":v2": {"S": "Call Me Today"}
}
```
```shell
$ aws dynamodb query --table-name MusicCollection \
    --key-condition-expression "Artist = :v1 AND SongTitle = :v2" \
    --expression-attribute-values file://expression-attributes.json
{
    "Count": 1,
    "Items": [
        {
            "AlbumTitle": {
                "S": "Somewhat Famous"
            },
            "SongTitle": {
                "S": "Call Me Today"
            },
            "Artist": {
                "S": "No One You Know"
            }
        }
    ],
    "ScannedCount": 1,
    "ConsumedCapacity": null
}
```

### EC2

#### 创建keypair
```shell
$ aws ec2 create-key-pair --key-name MyKeyPair --query 'KeyMaterial' --output text > MyKeyPair.pem
```

#### 创建安全组
```shell
$ aws ec2 create-security-group --group-name my-sg --description "My security group" --vpc-id vpc-1a2b3c4d
{
    "GroupId": "sg-903004f8"
}
```

#### 添加入站规则

```shell
$ aws ec2 authorize-security-group-ingress --group-id sg-903004f8 --protocol tcp --port 22 --cidr 203.0.113.0/24
```

#### 创建EC2
```shell
$ aws ec2 run-instances --image-id ami-xxxxxxxx --count 1 --instance-type t2.micro --key-name MyKeyPair --security-group-ids sg-903004f8 --subnet-id subnet-6e7f829e
```

#### 列出EC2

列出类型为t2.micro的实例
```shell
$ aws ec2 describe-instances --filters "Name=instance-type,Values=t2.micro" --query "Reservations[].Instances[].InstanceId"
```


### IAM

#### 创建用户组

使用 `create-group` 命令创建组
```shell
$ aws iam create-group --group-name MyIamGroup
{
    "Group": {
        "GroupName": "MyIamGroup",
        "CreateDate": "2018-12-14T03:03:52.834Z",
        "GroupId": "AGPAJNUJ2W4IJVEXAMPLE",
        "Arn": "arn:aws-cn:iam::123456789012:group/MyIamGroup",
        "Path": "/"
    }
}
```

#### 创建用户

使用 `create-user` 命令创建用户
```shell
$ aws iam create-user --user-name MyUser
{
    "User": {
        "UserName": "MyUser",
        "Path": "/",
        "CreateDate": "2018-12-14T03:13:02.581Z",
        "UserId": "AIDAJY2PE5XUZ4EXAMPLE",
        "Arn": "arn:aws-cn:iam::123456789012:user/MyUser"
    }
}
```

#### 将用户添加到组中
使用 `add-user-to-group` 命令将用户添加到组中
```shell
$ aws iam add-user-to-group --user-name MyUser --group-name MyIamGroup
```

#### 将 IAM 托管策略附加到 IAM 用户

1. 确定要附加的策略的 Amazon 资源名称 (ARN)。以下命令使用 `list-policies` 查找具有名称 `PowerUserAccess` 的策略的 ARN。然后,它会将该 ARN 存储在环境变量中。
  ```shell
  $ export POLICYARN=$(aws iam list-policies --query 'Policies[?PolicyName==`PowerUserAccess`].{ARN:Arn}' --output text)
  $ echo $POLICYARN
  arn:aws-cn:iam::aws:policy/PowerUserAccess
  ```
2. 要附加策略,请使用 `attach-user-policy` 命令,并引用存放策略 ARN 的环境变量。
  ```shell
  $ aws iam attach-user-policy --user-name MyUser --policy-arn $POLICYARN
  ```
3. 通过运行 `list-attached-user-policies` 命令验证策略已附加到此用户。
  ```shell
  $ aws iam list-attached-user-policies --user-name MyUser
  {
      "AttachedPolicies": [
          {
              "PolicyName": "PowerUserAccess",
              "PolicyArn": "arn:aws-cn:iam::aws:policy/PowerUserAccess"
          }
      ]
  }
  ```

#### 为 IAM 用户设置初始密码

```shell
$ aws iam create-login-profile --user-name MyUser --password My!User1Login8P@ssword --password-reset-required
{
    "LoginProfile": {
        "UserName": "MyUser",
        "CreateDate": "2018-12-14T17:27:18Z",
        "PasswordResetRequired": true
    }
}
```

#### 更改 IAM 用户的密码

可以使用 `update-login-profile` 命令更改 IAM 用户的密码。
```shell
$ aws iam update-login-profile --user-name MyUser --password My!User1ADifferentP@ssword
```

#### 创建访问密钥
使用 `create-access-key` 命令为 IAM 用户创建访问密钥。访问密钥是一组安全凭证,由访问密钥 ID 和私有密钥组成。

IAM 用户一次只能创建两个访问密钥。如果您尝试创建第三组,则命令返回 `LimitExceeded` 错误。
```shell
$ aws iam create-access-key --user-name MyUser
{
    "AccessKey": {
        "UserName": "MyUser",
        "AccessKeyId": "AKIAIOSFODNN7EXAMPLE",
        "Status": "Active",
        "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
        "CreateDate": "2018-12-14T17:34:16Z"
    }
}
```
#### 删除访问密钥

使用 `delete-access-key` 命令为 IAM 用户删除访问密钥。使用访问密钥 ID 指定要删除的访问密钥。
```shell
$ aws iam delete-access-key --user-name MyUser --access-key-id AKIAIOSFODNN7EXAMPLE
```


### S3

#### 创建桶

```shell
$ aws s3 mb s3://bucket-name
```

#### 复制
```shell
$ aws s3 cp filename.txt s3://bucket-name
$ aws s3 cp s3://bucket-name/filename.txt ./

# 使用 cat 文本编辑器,将文本“hello world”流式传输到 s3://bucket-name/filename.txt 文件
$ cat "hello world" | aws s3 cp - s3://bucket-name/filename.txt

# 将 s3://bucket-name/filename.txt 文件流式传输到 stdout,并将内容输出到控制台
$ aws s3 cp s3://bucket-name/filename.txt -
hello world

# 将 s3://bucket-name/pre 的内容流式传输到 stdout,使用 bzip2 命令压缩文件,并将名为 key.bz2 的新压缩文件上传到 s3://bucket-name
$ aws s3 cp s3://bucket-name/pre - | bzip2 --best | aws s3 cp - s3://bucket-name/key.bz2
```

#### 移动

```shell
# 将对象从 s3://bucket-name/example 移动到 s3://my-bucket/
$ aws s3 mv s3://bucket-name/example s3://my-bucket/

# 将本地文件从当前工作目录移动到 Amazon S3 存储桶
$ aws s3 mv filename.txt s3://bucket-name

# 将文件从 Amazon S3 存储桶移动到当前工作目录,其中 ./ 指定当前的工作目录。
$ aws s3 mv s3://bucket-name/filename.txt ./
```

#### 同步

s3 sync 命令同步一个存储桶与一个目录中的内容,或者同步两个存储桶中的内容。通常,s3 sync 在源和目标之间复制缺失或过时的文件或对象。不过,您还可以提供 --delete 选项来从目标中删除源中不存在的文件或对象。
```shell

$ aws s3 sync . s3://my-bucket/path
upload: MySubdirectory\MyFile3.txt to s3://my-bucket/path/MySubdirectory/MyFile3.txt
upload: MyFile2.txt to s3://my-bucket/path/MyFile2.txt
upload: MyFile1.txt to s3://my-bucket/path/MyFile1.txt

# Delete local file
$ rm ./MyFile1.txt


# Sync with deletion - object is deleted from bucket
$ aws s3 sync . s3://my-bucket/path --delete
delete: s3://my-bucket/path/MyFile1.txt

# Delete object from bucket
$ aws s3 rm s3://my-bucket/path/MySubdirectory/MyFile3.txt
delete: s3://my-bucket/path/MySubdirectory/MyFile3.txt

# Sync with deletion - local file is deleted
$ aws s3 sync s3://my-bucket/path . --delete
delete: MySubdirectory\MyFile3.txt

# Sync with Infrequent Access storage class
$ aws s3 sync . s3://my-bucket/path --storage-class STANDARD_IA
```

#### 删除桶以及桶内所有有内容

```shell
$ aws s3 rb s3://bucket-name --force
```

#### 删除桶中对象

从 s3://bucket-name/example 中删除所有对象
```shell
$ aws s3 rm s3://bucket-name/example
```
