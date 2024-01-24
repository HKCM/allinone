[toc]

### IAM Policy

IAM Policy是附加到用户,用户组或是角色的策略

#### 允许 IAM 用户访问某个存储桶

在本示例中,授予 AWS 账户中的 IAM 用户访问其中一个存储桶 awsexamplebucket1 的权限,以便该用户能够添加、更新和删除对象。

**注意:这个示例策略是加给用户或角色的,而不是加到bucket policy**

```json
{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Effect":"Allow",
         "Action": "s3:ListAllMyBuckets",
         "Resource":"arn:aws:s3:::*"
      },
      {
         "Effect":"Allow",
         "Action":["s3:ListBucket","s3:GetBucketLocation"],
         "Resource":"arn:aws:s3:::awsexamplebucket1"
      },
      {
         "Effect":"Allow",
         "Action":[
            "s3:PutObject",
            "s3:PutObjectAcl",
            "s3:GetObject",
            "s3:GetObjectAcl",
            "s3:DeleteObject"
         ],
         "Resource":"arn:aws:s3:::awsexamplebucket1/*"
      }
   ]
}
```

#### 允许每个 IAM 用户访问存储桶中的文件夹

两个 IAM 用户(Alice 和 Bob)可以访问 examplebucket 存储桶,以便他们可以添加、更新和删除对象。但是,想要限制每个用户对存储桶中单个文件夹的访问权限。可以使用与用户名称匹配的名称创建文件夹。

1. 对用户单独创建策略
  ```json
  {
    "Version":"2012-10-17",
    "Statement":[
        {
          "Effect":"Allow",
          "Action":[
              "s3:PutObject",
              "s3:GetObject",
              "s3:GetObjectVersion",
              "s3:DeleteObject",
              "s3:DeleteObjectVersion"
          ],
          "Resource":"arn:aws:s3:::awsexamplebucket1/Alice/*"
        }
    ]
  }
  ```
  同样,对Bob也创建一个

2. 使用策略变量的策略并将该策略附加到一个组,而不是将策略附加到单个用户。然后将 Alice 和 Bob 添加到该组中。

  策略变量 ${aws:username} 将替换为请求者的用户名称。例如,如果 Alice 发送了一个请求以放置对象,只有当 Alice 将对象上传到 examplebucket/Alice 文件夹后,才允许该操作。
  ```json
  {
    "Version":"2012-10-17",
    "Statement":[
        {
          "Effect":"Allow",
          "Action":[
              "s3:PutObject",
              "s3:GetObject",
              "s3:GetObjectVersion",
              "s3:DeleteObject",
              "s3:DeleteObjectVersion"
          ],
          "Resource":"arn:aws:s3:::awsexamplebucket1/${aws:username}/*"
        }
    ]
  }
  ```




### Bucket Policy

对于Bucket的策略,

1. 在Bucket是`Public`的情况下只允许`XXX`访问,则`Effect`字段应设置为`Deny`然后在`Condition`字段添加`Not`条件。
2. 在Bucket是`NotPublic`的情况下允许`XXX`访问,则`Effect`字段应设置为`Allow`然后设定在`Condition`字段添加`Not`条件

#### 存储桶拥有者授予跨账户存储桶权限

```json
{
   "Version": "2012-10-17",
   "Statement": [
      {
         "Sid": "Example permissions",
         "Effect": "Allow",
         "Principal": {
            "AWS": "arn:aws:iam::AccountB-ID:root"
         },
         "Action": [
            "s3:GetBucketLocation",
            "s3:ListBucket"
         ],
         "Resource": [
            "arn:aws:s3:::awsexamplebucket1"
         ]
      }
   ]
}

```

#### 以下存储桶策略使对象可公开访问。
```json
{
    "Version":"2012-10-17",
    "Statement": [
        {
            "Sid":"GrantAnonymousReadPermissions",
            "Effect":"Allow",
            "Principal": "*",
            "Action":["s3:GetObject"],
            "Resource":["arn:aws:s3:::awsexamplebucket1/*"]
        }
    ]
}
```

#### 基于IP限制

除特定外部 IP 地址外,阻止所有的流量 - 只允许特定IP访问
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:Get*",
      "Resource": "arn:aws:s3:::examplebucket/*",
      "Condition": {
         "NotIpAddress": {
            "aws:SourceIp": [
                "54.240.143.0/24",
                "2001:DB8:1234:5678::/64"
              ]
         }
      } 
    } 
  ]
}
```

阻止特定外部 IP 地址的流量 - - 不允许特定IP访问
```json
{
  "Id":"PolicyId2",
  "Version":"2012-10-17",
  "Statement":[
    {
      "Sid":"AllowIPmix",
      "Effect":"Allow",
      "Principal":"*",
      "Action":"s3:*",
      "Resource":"arn:aws:s3:::awsexamplebucket1/*",
      "Condition": {
        "NotIpAddress": {
          "aws:SourceIp": [
	          "54.240.143.128/30",
	          "2001:DB8:1234:5678:ABCD::/80"
          ]
        }
      }
    }
  ]
}
```

允许特定IP访问,阻止特定IP访问
```json
{
  "Id":"PolicyId2",
  "Version":"2012-10-17",
  "Statement":[
    {
      "Sid":"AllowIPmix",
      "Effect":"Allow",
      "Principal":"*",
      "Action":"s3:*",
      "Resource":"arn:aws:s3:::awsexamplebucket1/*",
      "Condition": {
        "IpAddress": {
          "aws:SourceIp": [
            "54.240.143.0/24",
	          "2001:DB8:1234:5678::/64"
          ]
        },
        "NotIpAddress": {
          "aws:SourceIp": [
            "54.240.143.128/30",
            "2001:DB8:1234:5678:ABCD::/80"
          ]
        }
      }
    }
  ]
}
```

#### 限制对特定 HTTP 引用站点的访问
假设你拥有一个网站,其域名为 www.example.com 或 example.com,并且带有指向存储在 Amazon S3 存储桶 awsexamplebucket1 中的照片和视频的链接。默认情况下,所有 Amazon S3 资源都是私有的,因此只有创建资源的 AWS 账户才能访问它们。要允许从网站对这些对象进行读取访问,您可以添加一个存储桶策略允许 s3:GetObject 权限,并附带使用 aws:Referer 键的条件,即获取请求必须来自特定的网页。以下策略指定带有 StringLike 条件键的 aws:Referer 条件。
```
{
  "Version":"2012-10-17",
  "Id":"http referer policy example",
  "Statement":[
    {
      "Sid":"Allow get requests originating from www.example.com and example.com.",
      "Effect":"Allow",
      "Principal":"*",
      "Action":[
        "s3:GetObject",
        "s3:GetObjectVersion"
      ],
      "Resource":"arn:aws:s3:::`awsexamplebucket1`/*",
      "Condition":{
        "StringLike":{
          "aws:Referer":[
            "http://www.example.com/*",
            "http://example.com/*"
          ]
        }
      }
    }
  ]
}
```
确保您使用的浏览器在请求中包含 HTTP referer 标头。




#### 特定用户或角色仅允许访问指定bucket

我想授予用户的 Amazon Simple Storage Service (Amazon S3) 控制台访问某个存储桶的权限。但是,我不希望用户能够查看其他账户中的存储桶。如何限制用户的控制台访问权限,使之只能访问某个存储桶？

要做到这个除了bucket policy之外,还需要一个额外的操作:

* 删除用户或角色的`s3:ListAllMyBuckets`的权限

以下示例策略用于访问存储桶。该策略允许用户仅对 `AWSDOC-EXAMPLE-BUCKET` 执行 `s3:ListBucket`、`s3:PutObject` 和 `s3:GetObject` 操作:

```json
{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Effect":"Allow",
         "Action":[
            "s3:ListBucket"
         ],
         "Resource":"arn:aws:s3:::AWSDOC-EXAMPLE-BUCKET"
      },
      {
         "Effect":"Allow",
         "Action":[
            "s3:PutObject",
            "s3:GetObject"
         ],
         "Resource":"arn:aws:s3:::AWSDOC-EXAMPLE-BUCKET/*"
      }
   ]
}
```

更改这些权限后,用户在访问主 Amazon S3 控制台时会收到“访问被拒绝”错误。用户必须使用指向存储桶或文件夹的直接控制台链接访问存储桶。指向存储桶的直接控制台链接与以下类似:
```
https://s3.console.aws.amazon.com/s3/buckets/AWSDOC-EXAMPLE-BUCKET/
```

#### 特定用户或角色仅允许访问指定bucket的特定文件夹(前缀)

我想授予用户的 Amazon Simple Storage Service (Amazon S3) 控制台访问某个特定存储桶文件夹(前缀)的权限。但是,我不希望用户能够查看其他账户中的存储桶或其他存储桶中的文件夹。如何限制用户的控制台访问权限,使之只能访问某个存储桶的文件夹？

要做到这个除了bucket policy之外,还需要一个额外的操作:

* 删除用户或角色的`s3:ListAllMyBuckets`的权限

以下示例策略用于访问文件夹。例如,以下策略允许用户对 `AWSDOC-EXAMPLE-BUCKET`中的 `folder2` 执行 `s3:ListBucket`、`s3:PutObject` 和 `s3:GetObject` 操作:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowUsersToAccessFolder2Only",
            "Effect": "Allow",
            "Action": [
                "s3:GetObject*",
                "s3:PutObject*"
            ],
            "Resource": [
                "arn:aws:s3:::AWSDOC-EXAMPLE-BUCKET/folder1/folder2/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket*"
            ],
            "Resource": [
                "arn:aws:s3:::AWSDOC-EXAMPLE-BUCKET"
            ],
            "Condition": {
                "StringLike": {
                    "s3:prefix": [
                        "folder1/folder2/*"
                    ]
                }
            }
        }
    ]
}
```

更改这些权限后,用户在访问主 Amazon S3 控制台时会收到“访问被拒绝”错误。用户必须使用指向存储桶或文件夹的直接控制台链接访问存储桶。指向存储桶的直接控制台链接与以下类似:
```
https://s3.console.aws.amazon.com/s3/buckets/AWSDOC-EXAMPLE-BUCKET/folder1/folder2/
```

#### 特定Bucket仅允许特定角色访问

有一个存放着公司绝密的文件的bucket,只允许特定角色和根用户访问,不能以IAM用户身份访问,即使拥有管理员权限也无法访问
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      # 警告 这个示例很危险不要轻易尝试,除非你知道自己在做什么
      "Principal": "*",
      "Action": "s3:Get*",
      "Resource": [
        "arn:aws:s3:::MyExampleBucket",
        "arn:aws:s3:::MyExampleBucket/*"
      ],
      "Condition": {
        "StringNotLike": {
          "aws:userId": [
            "AROAEXAMPLEID:*",
            "111111111111"
          ]
        }
      }
    }
  ]
}
```

* `AROAEXAMPLEID`: 通过AWS CLI的`iam get-role`获取的`RoleId`,代表role
* `111111111111`: AWS AccountId,代表根用户

可以和下方的特定用户配合使用

#### 特定Bucket仅允许特定IAM用户访问

有一个存放着公司绝密的文件的bucket,只允许特定的用户(例如CEO)和根用户访问,即使拥有管理员权限也无法访问
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      # 警告 这个示例很危险不要轻易尝试,除非你知道自己在做什么
      "Principal": "*",
      "Action": "s3:Get*",
      "Resource": [
        "arn:aws:s3:::MyExampleBucket",
        "arn:aws:s3:::MyExampleBucket/*"
      ],
      "Condition": {
        "StringNotLike": {
          "aws:userId": [
            "AIDAEXAMPLEID",
            "111111111111"
          ]
        }
      }
    }
  ]
}
```
* `AIDAEXAMPLEID`: 通过AWS CLI的`iam get-user`获取的`UserId`,代表user
* `111111111111`: AWS AccountId,代表根用户

可以和上方的特定角色配合使用

#### 仅允许从VPC endpoint访问

允许VPC内的机器能够直接访问S3

```json
{
  "Version": "2012-10-17",
  "Id": "VPCe and SourceIP",
  "Statement": [{
    "Sid": "VPCe and SourceIP",
    "Effect": "Deny",
    "Principal": "*",
    "Action": "s3:Get*",
    "Resource": [
      "arn:aws:s3:::awsexamplebucket",
      "arn:aws:s3:::awsexamplebucket/*"
    ],
    "Condition": {
      "StringNotEquals": {
        "aws:sourceVpce": [
          "vpce-1111111",
          "vpce-2222222"
        ]
      },
      "StringNotLike": {
        "aws:userId": [
          "AROAEXAMPLEID:*",
          "AIDAEXAMPLEID",
          "111111111111"
        ]
      }
    }
  }]
}
```

* `AROAEXAMPLEID`: 通过AWS CLI的`iam get-role`获取的`RoleId`,代表role
* `AIDAEXAMPLEID`: 通过AWS CLI的`iam get-user`获取的`UserId`,代表user
* `111111111111`: AWS AccountId,代表根用户



#### 向 OAI(Origin Access Identity ) 授予读写访问权限

这是配合CloudFront使用的策略,允许通过CloudFront访问,不能直接访问Bucket

下面的示例允许 OAI 读取和写入指定存储桶(s3:GetObject 和 s3:PutObject)中的对象。这样允许查看器将文件通过 CloudFront 上传到 Amazon S3 存储桶。
```json
{
    "Version": "2012-10-17",
    "Id": "PolicyForCloudFrontPrivateContent",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity EH1HDMB1FH2TC"
            },
            "Action": [
                "s3:GetObject",
                "s3:PutObject"
            ],
            "Resource": "arn:aws:s3:::aws-example-bucket/*"
        }
    ]
}
```


#### VPC endpoint与Access Point结合

当Bucket越来越多,每次新增bucket都需要修改endpoint,使用该策略限制从S3请求必须通过Access Point。这是VPC endpoint策略。

```json
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "AllowUseOfS3",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:*",
            "Resource": "*"
        },
        {
            "Sid": "OnlyIfAccessedViaAccessPoints",
            "Effect": "Deny",
            "Principal": "*",
            "Action": "s3:*",
            "Resource": "*",
            "Condition": {
                "ArnNotLikeIfExists": {
                    "s3:DataAccessPointArn": "arn:aws:s3:us-east-1:<Account ID>:accesspoint/*"
                }
            }
        }
    ]
}
```


#### 检查仅允许从特定 IP 地址上传的条件,类似条件如下:
```json
"Condition": {
  "IpAddress": {
    "aws:SourceIp": "54.240.143.0/24"
  }
}
```
如果存储桶策略具有此条件,IAM 用户必须从允许的 IP 地址访问存储桶。

#### 检查仅允许在对象为特定存储类时上传的条件,类似条件如下:
```json
"Condition": {
  "StringEquals": {
    "s3:x-amz-storage-class": [
      "STANDARD_IA"
    ]
  }
}
```
如果策略具有此条件,用户必须使用允许的存储类上传对象。例如,上一个条件语句需要 STANDARD_IA 存储类,因此用户必须使用类似于以下内容的 AWS 命令行界面 (AWS CLI) 命令上传对象:
```
aws s3api put-object --bucket my_bucket --key examplefile.jpg --body c:\examplefile.jpg --storage-class STANDARD_IA
```

#### 检查仅允许在对象分配有特定访问控制列表 (ACL) 时进行上传的条件,类似条件如下:
```json
"Condition": {
  "StringEquals": {
      "s3:x-amz-acl":["public-read"]
  }
}
```
如果策略具有此条件,则用户必须用允许的 ACL 上传对象。例如,由于上一个条件需要 public-read ACL,用户必须使用类似于以下内容的命令上传对象:
```
aws s3api put-object --bucket my_bucket --key examplefile.jpg --body c:\examplefile.jpg --acl public-read
```

#### 检查要求上传授予存储桶拥有者(典型用户 ID)对象完全控制权的条件,类似条件如下:
```json
"Condition": {
  "StringEquals": {
    "s3:x-amz-grant-full-control": "id=AccountA-CanonicalUserID"
  }
}
```



如果策略具有此条件,则用户必须使用类似于以下内容的命令上传对象:
```shell
# 获取AccountA-CanonicalUserID
aws s3api list-buckets --query Owner.ID --output text
example95f90f6fae9c770602ce44f64f76aee3386099178440a6e162952abcd

aws s3api put-object --bucket my_bucket --key examplefile.jpg --body c:\examplefile.jpg --acl bucket-owner-full-control
```

#### 检查仅允许在特定 AWS Key Management System (AWS KMS) 密钥加密对象时进行上传的条件,类似条件如下:
```json
"Condition": {
  "StringEquals": {
    "s3:x-amz-server-side-encryption-aws-kms-key-id": "arn:aws:kms:us-east-1:111122223333:key/*"
  }
}
```
如果策略具有此条件,则用户必须使用类似于以下内容的命令上传对象:
```
aws s3api put-object --bucket my_bucket --key examplefile.jpg --body c:\examplefile.jpg --ssekms-key-id arn:aws:kms:us-east-1:111122223333:key/*
```

#### 检查仅允许在对象使用特定类型服务器端加密时进行上传的条件,类似条件如下:
```json
"Condition": {
  "StringEquals": {
    "s3:x-amz-server-side-encryption": "AES256"
  }
}
```
如果策略具有此条件,用户必须使用类似于以下内容的命令上传对象:
```
aws s3api put-object --bucket my_bucket --key examplefile.jpg --body c:\examplefile.jpg --server-side-encryption "AES256"
```

#### 将最长保留期限设置为 10 天。
```json
{
    "Version": "2012-10-17",
    "Id": "<SetRetentionLimits",
    "Statement": [
        {
            "Sid": "<SetRetentionPeriod",
            "Effect": "Deny",
            "Principal": "*",
            "Action": [
                "s3:PutObjectRetention"
            ],
            "Resource": "arn:aws:s3:::<awsexamplebucket1>/*",
            "Condition": {
                "NumericGreaterThan": {
                    "s3:object-lock-remaining-retention-days": "10"
                }
            }
        }
    ]
}
```

### 访问点相关策略

#### 将访问控制委派给访问点的存储桶策略
```json
{
    "Version": "2012-10-17",
    "Statement" : [
    {
        "Effect": "Allow",
        "Principal" : { "AWS": "*" },
        "Action" : "*",
        "Resource" : [ "Bucket ARN", "Bucket ARN/*"],
        "Condition": {
            "StringEquals" : { "s3:DataAccessPointAccount" : "Bucket owner's account ID" }
        }
    }]
}
```

#### 访问点策略授予

以下访问点策略通过账户 123456789012 中的访问点 my-access-point 向账户 123456789012 中的 IAM 用户 Alice 授予对具有前缀 Alice/ 的 GET 和 PUT 对象的权限。
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

**注意**

要使访问点策略有效地向 Alice 授予访问权限,底层存储桶也必须对 Alice 允许相同的访问权限。可以将存储桶的访问控制委派到访问点,如[将访问控制委派给访问点的存储桶策略](#将访问控制委派给访问点的存储桶策略)中所述。或者,您也可以将以下策略添加到底层存储桶中,以便向 Alice 授予必要的权限。请注意,访问点策略和存储桶策略的 Resource 条目不同。
```json
{
    "Version": "2012-10-17",
    "Statement": [
    {
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::123456789012:user/Alice"
        },
        "Action": ["s3:GetObject", "s3:PutObject"],
        "Resource": "arn:aws:s3:::awsexamplebucket1/Alice/*"
    }]    
}
```

#### 允许用户仅读取具有特定标签的对象

以下访问点策略通过账户 123456789012 中的访问点 my-access-point 向账户 123456789012 中的 IAM 用户 Bob 授予对 GET 对象的权限,这些权限具有值设为 finance 的标签键 data。
```json
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

#### 允许用户添加指定的对象标签

以下权限策略将向用户授予执行 s3:PutObjectTagging 操作的权限,这使用户可以将标签添加到现有对象。条件限制了用户可使用的标签键。条件使用 s3:RequestObjectTagKeys 条件键指定一组标签键。`用户可以在 PutObjectTagging 中发送空标签集,这是策略允许的 (请求中的空标签集将删除对象上的任何现有标签)`

ForAllValues: 测试请求集的每个成员的值是否为条件键集的子集。如果请求中的`每个键值`均与策略中的至少一个值匹配,则条件返回 true。如果请求中没有键或者键值解析为空数据集(如空字符串),则也会返回 true,因为空集也是它的子集
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObjectTagging"
      ],
      "Resource": [
        "arn:aws:s3:::awsexamplebucket1/*"
      ],
      "Principal":{
        "CanonicalUser":[
            "64-digit-alphanumeric-value"
         ]
       },
      "Condition": {
        "ForAllValues:StringLike": {
          "s3:RequestObjectTagKeys": [
            "Owner",
            "CreationDate"
          ]
        }
      }
    }
  ]
}
```

以下权限策略将向用户授予执行 s3:PutObjectTagging 操作的权限,这使用户可以将标签添加到现有对象。条件限制了用户可使用的标签键。条件使用 s3:RequestObjectTagKeys 条件键指定一组标签键。`用户不可以在 PutObjectTagging 中发送空标签集`

ForAnyValue: 测试请求值集的至少一个成员与条件键值集的至少一个成员匹配。如果请求中的任何一个键值与策略中的任何一个条件值匹配,则条件返回 true。对于没有匹配的键或空数据集,条件返回 false。
```json
{

  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObjectTagging"
      ],
      "Resource": [
        "arn:aws:s3:::awsexamplebucket1/*"
      ],
      "Principal":{
        "AWS":[
            "arn:aws:iam::account-number-without-hyphens:user/username"
         ]
       },
      "Condition": {
        "ForAllValues:StringLike": {
          "s3:RequestObjectTagKeys": [
            "Owner",
            "CreationDate"
          ]
        },
        "ForAnyValue:StringLike": {
          "s3:RequestObjectTagKeys": [
            "Owner",
            "CreationDate"
          ]
        }
      }
    }
  ]
}
```

#### 存储桶中对象所有权

存储桶策略指定只有当对象的 ACL 设置为 bucket-owner-full-control 时,账户 111122223333 才能将对象上传到 awsdoc-example-bucket。
```json
{
   "Version": "2012-10-17",
   "Statement": [
      {
         "Sid": "Only allow writes to my bucket with bucket owner full control",
         "Effect": "Allow",
         "Principal": {
            "AWS": [
               "arn:aws:iam::111122223333:user/ExampleUser"
            ]
         },
         "Action": [
            "s3:PutObject"
         ],
         "Resource": "arn:aws:s3:::awsdoc-example-bucket/*",
         "Condition": {
            "StringEquals": {
               "s3:x-amz-acl": "bucket-owner-full-control"
            }
         }
      }
   ]
}
```

#### 允许用户添加包含特定标签键和值的对象标签
以下用户策略将向用户授予执行 s3:PutObjectTagging 操作的权限,这使用户可以在现有对象上添加标签。条件要求用户包含值设置为 Project 的特定标签 (X)。
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObjectTagging"
      ],
      "Resource": [
        "arn:aws:s3:::awsexamplebucket1/*"
      ],
      "Principal":{
        "AWS":[
            "arn:aws:iam::account-number-without-hyphens:user/username"
         ]
       },
      "Condition": {
        "StringEquals": {
          "s3:RequestObjectTag/Project": "X"
        }
      }
    }
  ]
}

```


#### 允许查看存储桶列示内容的访问点策略

以下访问点策略通过账户 123456789012 中的访问点 my-access-point 授予账户 123456789012 中的 IAM 用户 Charles 查看底层存储桶中包含的对象的权限。
```json
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

