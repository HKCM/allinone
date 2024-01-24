
### 什么是访问点AccessPoint

了解一个东西最方便的就是看它的使用场景.那么访问点的使用场景是什么呢?以下以存储桶策略为例,不谈IAM.

想象一下,公司有一个桶,桶里面有IT,财务,开发和HR各自的文件夹(其实S3并没有文件夹概念,这里是为了方便理解),目前存储桶策略只允许各自部门访问各自的文件夹.  
随着公司壮大,多了一个部门,这个部门要求要访问IT和开发的文件夹,怎么办呢?那只有修改存储桶策略.有后来又多了一个部门,要求访问HR和财务的文件夹...  
随着部门的扩展,存储桶策略会越来越复杂且难以追踪和管理,存储桶策略的修改还很有可能会影响到其他部门.  

所以访问点应运而生,当有新的部门需要访问,可以通过创建访问点和访问点策略以达到访问需求,新的访问点策略不会影响现有的访问点


所以访问点是简化了 S3 中共享数据集的大规模数据访问管理问题.

### 访问点限制和局限性

1. 访问点策略的大小限制为 20 KB
2. 只能为自己拥有的存储桶创建访问点
3. 每个访问点只与一个存储桶相关联
4. 使用 REST API 向访问点发出请求时,必须使用 AWS Signature Version 4
5. 每个区域每个 AWS 账户最多可以创建 1,000 个访问点。可以增加配额
6. 创建访问点之后不可更改访问点的阻止公有访问设置
7. 访问点不支持跨区域复制

### Access Point CLI

创建一个访问点
```
$ aws s3control create-access-point --name example-ap --account-id 123456789012 --bucket example-bucket
```

创建仅限 VPC 访问的访问点并验证
```
$ aws s3control create-access-point --name example-vpc-ap --account-id 123456789012 --bucket example-bucket --vpc-configuration VpcId=vpc-1a2b3c
$ aws s3control get-access-point --name example-vpc-ap --account-id 123456789012
```

创建具有非默认阻止公有访问设置的访问点
```
$ aws s3control create-access-point --name example-ap --account-id 123456789012 --bucket example-bucket \
  --public-access-block-configuration BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=true,RestrictPublicBuckets=true
```

要将访问点与 VPC 搭配使用,必须修改 VPC 终端节点的访问策略。VPC 终端节点允许流量从 VPC 流向 Amazon S3。它们具有访问控制策略,用于控制如何允许 VPC 内的资源与 S3 交互。仅当 VPC 终端节点策略同时对访问点和底层存储桶授予访问权限时,从 VPC 通过访问点到 S3 的请求才会成功。

配置 VPC 终端节点策略,需要同时具有bucket和access point的资源
```
{
    "Version": "2012-10-17",
    "Statement": [
    {
        "Principal": "*",
        "Action": [
            "s3:GetObject"
        ],
        "Effect": "Allow",
        "Resource": [
            "arn:aws:s3:::awsexamplebucket1/*",
            "arn:aws:s3:us-west-2:123456789012:accesspoint/example-vpc-ap/object/*"
        ]
    }]
}
```

### 访问点条件键

匹配访问点 ARN 的字符串

```
"Condition" : {
    "StringLike": {
        "s3:DataAccessPointArn": "arn:aws:s3:us-west-2:123456789012:accesspoint/*"
    }
}
```

匹配访问点拥有者的账户 ID

```
"Condition" : {
    "StringEquals": {
        "s3:DataAccessPointAccount": "123456789012"
    }
}
```
匹配网络起源(Internet 或 VPC)

```
"Condition" : {
    "StringEquals": {
        "s3:AccessPointNetworkOrigin": "VPC"
    }
}
```

### 策略示例

#### 只能从访问点访问
如果希望存储桶只通过访问点访问,可以设置如下`存储桶策略`
```json
{
    "Version": "2012-10-17",
    "Statement" : [
    {
        "Effect": "Allow",
        "Principal" : { "AWS": "*" },
        "Action" : "*", #很危险
        "Resource" : [ "Bucket ARN", "Bucket ARN/*"],
        "Condition": {
            "StringEquals" : { "s3:DataAccessPointAccount" : "Bucket owner's account ID" }
        }
    }]
}
```

#### 访问点策略授予

以下访问点策略通过账户 123456789012 中的访问点 my-access-point 向账户 123456789012 中的 IAM 用户 Alice 授予对具有前缀 `Alice/` 的 `GET` 和 `PUT` 对象的权限。

```json
{
    "Version":"2012-10-17",
    "Statement": [
    {
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::123456789012:user/Alice"
        },
        "Action": ["s3:GetObject", "s3:PutObject"],
        "Resource": "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point/object/Alice/*"
    }]
}
```

#### 带标签条件的访问点策略


以下访问点策略通过账户 123456789012 中的访问点 my-access-point 向账户 123456789012 中的 IAM 用户 Bob 授予对 `GET` 对象的权限,这些权限仅能访问标签为data 值为 finance的资源。

```
{
    "Version":"2012-10-17",
    "Statement": [
    {
        "Effect":"Allow",
        "Principal" : {
            "AWS": "arn:aws:iam::123456789012:user/Bob"
        },
        "Action":"s3:GetObject",
        "Resource" : "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point/object/*",
        "Condition" : {
            "StringEquals": {
                "s3:ExistingObjectTag/data": "finance"
            }
        }
    }]
}
```

#### 允许查看存储桶列示内容的访问点策略

以下访问点策略通过账户 123456789012 中的访问点 my-access-point 授予账户 123456789012 中的 IAM 用户 Charles 查看底层存储桶中包含的对象的权限。
```
{
    "Version":"2012-10-17",
    "Statement": [
    {
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::123456789012:user/Charles"
        },
        "Action": "s3:ListBucket",
        "Resource": "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"
    }]
}
```

#### 服务控制策略

**以下服务控制策略要求使用 VPC 网络起源创建所有新访问点。实施此策略时,组织中的用户无法创建可从 Internet 访问的新访问点。**
```
{
    "Version": "2012-10-17",
    "Statement": [
    {
        "Effect": "Deny",
        "Action": "s3:CreateAccessPoint",
        "Resource": "*",
        "Condition": {
            "StringNotEquals": {
                "s3:AccessPointNetworkOrigin": "VPC"
            }
        }
    }]
}
```

#### 将 S3 操作限制为 VPC 网络起源的存储桶策略

以下存储桶策略限制为只能通过具有 VPC 网络起源的访问点来访问存储桶 examplebucket 的所有 S3 对象操作。

**重要**

在使用类似此示例的语句之前,请确保您不需要使用访问点不支持的功能,例如跨区域复制。
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Deny",
            "Principal": "*",
            "Action": [
                "s3:AbortMultipartUpload",
                "s3:BypassGovernanceRetention",
                "s3:DeleteObject",
                "s3:DeleteObjectTagging",
                "s3:DeleteObjectVersion",
                "s3:DeleteObjectVersionTagging",
                "s3:GetObject",
                "s3:GetObjectAcl",
                "s3:GetObjectLegalHold",
                "s3:GetObjectRetention",
                "s3:GetObjectTagging",
                "s3:GetObjectVersion",
                "s3:GetObjectVersionAcl",
                "s3:GetObjectVersionTagging",
                "s3:ListMultipartUploadParts",
                "s3:PutObject",
                "s3:PutObjectAcl",
                "s3:PutObjectLegalHold",
                "s3:PutObjectRetention",
                "s3:PutObjectTagging",
                "s3:PutObjectVersionAcl",
                "s3:PutObjectVersionTagging",
                "s3:RestoreObject"
            ],
            "Resource": "arn:aws:s3:::examplebucket/*",
            "Condition": {
                "StringNotEquals": {
                    "s3:AccessPointNetworkOrigin": "VPC"
                }
            }
        }
    ]
}
```

#### CLI

```
$ aws s3api get-object --key Alice/my-image.jpg --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point download.jpg
$ aws s3api put-object --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point --key Alice/test.log --body test.log
$ aws s3api delete-object --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point --key Alice/my-image.jpg
$ aws s3api list-objects-v2 --bucket arn:aws:s3:us-east-1:123456789012:accesspoint/my-access-point
```
以上完全基于[AWS官方文档](https://docs.aws.amazon.com/),并结合自身理解创作


本文采用知识共享 署名-相同方式共享 3.0协议

署名-相同方式共享(BY-SA):使用者可以对本创作进行转载、节选、混编、二次创作,可以将其运用于商业用途,唯须署名作者,并且采用本创作的内容必须同样采用本协议进行授权
