[toc]



### S3高级命令

* [`cp`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/cp.html) 
* [`ls`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/ls.html)
* [`mb`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/mb.html)
* [`mv`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/mv.html)
* [`presign`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/presign.html)
* [`rb`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/rb.html)
* [`rm`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/rm.html)
* [`sync`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/sync.html)
* [`website`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3/website.html)


### 存储桶CLI

#### 查找规范用户ID
```shell
aws s3api list-buckets --query Owner.ID --output text
```
[AWS账户标识符](https://docs.aws.amazon.com/zh_cn/general/latest/gr/acct-identifiers.html#FindingCanonicalId)

#### 删除存储桶
如果存储桶未启用版本控制,则可将 rb(删除存储桶) AWS CLI 命令与 `--force` 参数结合使用来删除非空存储桶。此命令将先删除所有对象,然后再删除存储桶。
```shell
aws s3 rb s3://bucket-name --force  
```

#### 清空存储桶
如果存储桶未启用版本控制,可以将 rm(删除) AWS CLI 命令与 `--recursive` 参数结合使用来清空存储桶 (或删除部分带特定键名前缀的对象)。
```shell
# 以下 rm 命令将删除带键名前缀 doc 的对象,例如,doc/doc1 和 doc/doc2。
$ aws s3 rm s3://bucket-name/doc --recursive

# 使用以下命令删除所有对象,而无需指定前缀。
$ aws s3 rm s3://bucket-name --recursive
```

#### 获取存储桶状态
```shell
aws s3api get-bucket-policy --bucket DOC-EXAMPLE-BUCKET1 --expected-bucket-owner 111122223333
```

#### 启用和停止Transfer Acceleration
下示例设置 `Status=Enabled` 以对存储桶启用 Transfer Acceleration。可使用 `Status=Suspended` 暂停 Transfer Acceleration
```shell
aws s3api put-bucket-accelerate-configuration --bucket bucketname --accelerate-configuration Status=Enabled
```

#### 将存储桶内的对象批量加密
批量加密要注意的问题:
1. LastModified时间戳会改变
2. 如果在PUT或COPY事件上启用leS3事件通知,在复制现有对象以对其进行加密时会触发S3事件通知
3. 访问控制列表(ACL)重置为S3存储桶默认值。如果使用对象ACL,则必须在复制操作中添加它们
4. 元数据和标签,CLI_V2支持元数据和标签的复制,V1不支持
5. 存储类默认为“标准”。
6. 在启用了版本控制的存储桶中,此过程将创建加密对象的新版本,但不会修改现有的未加密对象版本
7. 如果使用了对象锁定,则保留期限将重置为存储桶默认值。
8. S3 ETag可能会更改。
  * 分段上传对象:如果与原始分段上传相比,副本使用不同的分段大小,则ETag会更改。
  * 如果复制时与原对象不同的加密密钥对对象进行加密,则该对象将不会具有与之前相同的ETag 
```shell
aws s3 cp s3://mybucket/ s3://mybucket/ --recursive --sse aws:kms --sse-kms-key-id arn:aws:kms:us-east-1:123456789012:alias/mykey
```

#### 上传生命周期配置
本地需要一个lifecycle.json 文件
```shell
aws s3api put-bucket-lifecycle-configuration  \
--bucket bucketname  \
--lifecycle-configuration file://lifecycle.json 
```

[生命周期示例](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/lifecycle-configuration-examples.html)

<details>
<summary>
  <strong>示例 1:指定筛选条件</strong>
</summary>
在此生命周期配置规则中,筛选条件指定了一个键前缀 (documents/)。因此,此规则应用于带键名前缀 documents/ 的对象,例如 documents/doc1.txt 和 documents/doc2.txt。

此规则指定两个引导 Amazon S3 完成以下任务的操作:

* 在对象创建 365 天(一年)后将其转换为 S3 Glacier 存储类别。
* 在对象创建 3650 天(10 年)后将其删除(Expiration 操作)。

##### **该操作针对当前版本所以对象总共存在10年,一年标准存储,九年Glacier存储**
```json
{
    "Rules": [
        {
            "Filter": {
                "Prefix": "documents/"
            },
            "Status": "Enabled",
            "Transitions": [
                {
                    "Days": 365,
                    "StorageClass": "GLACIER"
                }
            ],
            "Expiration": {
                "Days": 3650
            },
            "ID": "ExampleRule"
        }
    ]
}
```

##### 带有版本控制的生命周期配置

```json
{
    "Rules": [
        {
            "ID": "ExampleRule",
            "Filter": {
                "Prefix": "documents/"
            },
            "Status": "Enabled",
            "Transition": {
                "Days": 30,
                "StorageClass": "STANDARD_IA"
            },
            "NoncurrentVersionTransitions": {
              "NoncurrentDays": 90,
              "StorageClass": "GLACIER"
            },
            "NoncurrentVersionExpiration": {
                "NoncurrentDays": 3650
            }
        }
    ]
}
```

##### 带前缀及标签的生命周期配置

在以下配置中,该规则指定了一个 Expiration 操作,以下示例代码将生命周期规则应用于带 tax/ 键前缀且包含具有特定两个键和值的对象。

```json
{
  "Rules": [
      {
          "Filter": {
              "And": {
                  "Prefix": "tax/",
                  "Tags": [
                      {
                          "Value": "mytagvalue1", 
                          "Key": "mytagkey1"
                      }, 
                      {
                          "Value": "mytagvalue2", 
                          "Key": "mytagkey2"
                      }
                  ]
              }
          },
          "Status": "Enabled",
          "Expiration": {
                  "Days": 1
          },
          "ID": "ExampleRule"
      }
  ]
}
```

在以下配置中,该规则指定了一个 Transition 操作,该操作指示 Amazon S3 在对象创建 0 天后将其转换为 S3 Glacier 存储类别,在这种情况下,对象都有资格在创建后的 UTC 时间午夜存档到 Amazon S3 Glacier

```json
{
    "Rules": [
        {
            "Filter": {
                "Prefix": ""
            },
            "Status": "Enabled",
            "Transitions": [
                {
                    "Days": 0,
                    "StorageClass": "GLACIER"
                }
            ],
            "ID": "ExampleRule"
        }
    ]
}
```
</details>

<details>
<summary><strong>示例2: 多个规则</strong></summary>

希望不同的对象有不同的生命周期操作,则可以指定多个规则。以下生命周期配置有两个规则:

* 规则 1 应用于带键名前缀 classA/ 的对象。它指示 Amazon S3 在对象创建一年后将其转换为 S3 Glacier 存储类别,并在对象创建 10 年后使它们过期。
* 规则 2 应用于带键名前缀 classB/ 的对象。它指示 Amazon S3 在对象创建 90 天后将其转换为 S3 标准-IA 存储类,并在对象创建 1 年后将其删除。

```json
{
  "Rules": [
      {
          "Filter": {
              "Prefix": "ClassA/",
          },
          "Status": "Enabled",
          "Transitions": {
              "Days": 365,
              "StorageClass": "GLACIER"
          },
          "Expiration": {
              "Days": 3650
          },
          "ID": "ClassADocRule"
      },{
          "Filter": {
              "Prefix": "ClassB/",
          },
          "Status": "Enabled",
          "Transitions": {
              "Days": 90,
              "StorageClass": "STANDARD_IA"
          },
          "Expiration": {
              "Days": 365
          },
          "ID": "ClassBDocRule"
      }
  ]
}
```
</details>


<details>
<summary><strong>JSON格式语法</strong></summary>

```json
{
  "Rules": [
    {
      "Expiration": {
        "Date": timestamp,
        "Days": integer,
        "ExpiredObjectDeleteMarker": true|false
      },
      "ID": "string",
      "Prefix": "string",
      "Filter": {
        "Prefix": "string",
        "Tag": {
          "Key": "string",
          "Value": "string"
        },
        "And": {
          "Prefix": "string",
          "Tags": [
            {
              "Key": "string",
              "Value": "string"
            }
            ...
          ]
        }
      },
      "Status": "Enabled"|"Disabled",
      "Transitions": [
        {
          "Date": timestamp,
          "Days": integer,
          "StorageClass": "GLACIER"|"STANDARD_IA"|"ONEZONE_IA"|"INTELLIGENT_TIERING"|"DEEP_ARCHIVE"
        }
        ...
      ],
      "NoncurrentVersionTransitions": [
        {
          "NoncurrentDays": integer,
          "StorageClass": "GLACIER"|"STANDARD_IA"|"ONEZONE_IA"|"INTELLIGENT_TIERING"|"DEEP_ARCHIVE"
        }
        ...
      ],
      "NoncurrentVersionExpiration": {
        "NoncurrentDays": integer
      },
      "AbortIncompleteMultipartUpload": {
        "DaysAfterInitiation": integer
      }
    }
    ...
  ]
}
```
</details>

#### 获取存储桶生命周期配置
```shell
aws s3api get-bucket-lifecycle-configuration  \
--bucket bucketname 
```

#### 删除存储桶生命周期配置
```shell
aws s3api delete-bucket-lifecycle \
--bucket bucketname
```

### 对象CLI

#### 上传对象
```shell
aws s3api put-object --bucket ${mybucket} --body ./local_file --key hello.html --acl public-read
{
    "ETag": "\"e8f8ca569bda9ae941cc68013262abcd\"",
    "VersionId": "D.HahJrVUGdK4D.RZ0Jn5lDbNRKeABCD"
}

aws s3api put-object --body ./h1.html --key static/allowed/h.html --bucket user-test-us --content-type text/plain --acl public-read

# 使用桶所有者条件以确保 DOC-EXAMPLE-BUCKET1 归 AWS 账户 111122223333.
aws s3api put-object \
  --bucket DOC-EXAMPLE-BUCKET1 --key exampleobject --body example_file.txt \
  --expected-bucket-owner 111122223333

# 上传对象带标签

aws s3api put-object --bucket mybucket --key tax/delete3 --body ./h1.html \
--tagging "mytagkey1=mytagvalue1&mytagkey2=mytagvalue2"
```

#### 修改对象ACL
```shell
aws s3api put-object-acl --bucket DOC-EXAMPLE-BUCKET1 --key HappyFace.jpg --grant-full-control id="AccountA-CanonicalUserID" --profile AccountBadmin
```

#### 分段上传

AWS S3高级命令`s3 cp`、`s3 sync`以及`s3 mv`等都会自动执行分段操作 
```shell
# 分段上传可配置S3相关命令
aws configure set default.s3.max_concurrent_requests 60
aws configure set default.s3.multipart_threshold 5GB
aws configure set default.s3.multipart_chunksize 5GB

# S3默认配置
aws configure set default.s3.max_concurrent_requests 10
aws configure set default.s3.multipart_threshold 8MB
aws configure set default.s3.multipart_chunksize 8MB
```

以下是分段上传的示例
```shell
# 创建50M的文件
dd if=/dev/zero of=test.file bs=1048576 count=50

# 获取md5值,以便在上传后执行完整性检查时作为参考
md5 test.file
MD5 (test.file) = 25e317773f308e446cc84c503a6d1f85

# 使用split分割为25M大小
split -b 26214400 test.file test.file_
test.file_aa test.file_ab 

# 创建分段上传
aws s3api create-multipart-upload --bucket mybucket --key test.file
{
    "Bucket": "mybucket",
    "Key": "test.file",
    "UploadId": "B9e10l2zv9hb.jxutvVffxfnvLwhj_uAT9Qka2h8FlnttkDeZ09_ZLMV.2HHbvt54hxMmx.RkEC83c2zUFFvKXZaI.NcIwMr9gu7kih.ujS4dh3VQTim7QKbHHDZlbVxsv0pAAyb.rB8vzG2YQvAqHtGzlTr3HV2JiLCPYXg_Sc-"
}

# 第一部分上传,后面重复操作,记得修改part-number和body参数, ETag需要保留,后面要用
aws s3api upload-part --bucket mybucket --key test.file --part-number 1 --body test.file_aa --upload-id B9e10l2zv9hb.jxutvVffxfnvLwhj_uAT9Qka2h8FlnttkDeZ09_ZLMV.2HHbvt54hxMmx.RkEC83c2zUFFvKXZaI.NcIwMr9gu7kih.ujS4dh3VQTim7QKbHHDZlbVxsv0pAAyb.rB8vzG2YQvAqHtGzlTr3HV2JiLCPYXg_Sc-
{
    "ETag": "\"bed3c0a4a1407f584989b4009e9ce33f\""
}
aws s3api upload-part --bucket mybucket --key test.file --part-number 2 --body test.file_ab --upload-id B9e10l2zv9hb.jxutvVffxfnvLwhj_uAT9Qka2h8FlnttkDeZ09_ZLMV.2HHbvt54hxMmx.RkEC83c2zUFFvKXZaI.NcIwMr9gu7kih.ujS4dh3VQTim7QKbHHDZlbVxsv0pAAyb.rB8vzG2YQvAqHtGzlTr3HV2JiLCPYXg_Sc-
{
    "ETag": "\"bed3c0a4a1407f584989b4009e9ce33f\""
}

# list-parts
aws s3api list-parts --bucket mybucket --key test.file --upload-id B9e10l2zv9hb.jxutvVffxfnvLwhj_uAT9Qka2h8FlnttkDeZ09_ZLMV.2HHbvt54hxMmx.RkEC83c2zUFFvKXZaI.NcIwMr9gu7kih.ujS4dh3VQTim7QKbHHDZlbVxsv0pAAyb.rB8vzG2YQvAqHtGzlTr3HV2JiLCPYXg_Sc-

# 用ETag编成ETag Json文件,可以从list-parts命令中获取,类似这个
{
    "Parts": [
        {
            "PartNumber": 1,
            "ETag": "\"bed3c0a4a1407f584989b4009e9ce33f\""
        },
        {
            "PartNumber": 2,
            "ETag": "\"bed3c0a4a1407f584989b4009e9ce33f\""
        }
    ]
}

# 完成分段
aws s3api complete-multipart-upload --multipart-upload file://ETag.json --bucket mybucket --key test.file --upload-id B9e10l2zv9hb.jxutvVffxfnvLwhj_uAT9Qka2h8FlnttkDeZ09_ZLMV.2HHbvt54hxMmx.RkEC83c2zUFFvKXZaI.NcIwMr9gu7kih.ujS4dh3VQTim7QKbHHDZlbVxsv0pAAyb.rB8vzG2YQvAqHtGzlTr3HV2JiLCPYXg_Sc-
```

#### 列出存储桶中未完成的分段上传
```shell
aws s3api list-multipart-uploads --bucket DOC-EXAMPLE-BUCKET
```

#### 删除未完成的分段上传
```shell
aws s3api abort-multipart-upload --bucket DOC-EXAMPLE-BUCKET --key large_test_file --upload-id examplevQpHp7eHc_J5s9U.kzM3GAHeOJh1P8wVTmRqEVojwiwu3wPX6fWYzADNtOHklJI6W6Q9NJUYgjePKCVpbl_rDP6mGIr2AQJNKB
```

#### 列出对象
```shell
# 只列出目录
aws s3 --profile myprofile --region us-east-1 ls s3://111.111.111.111-test-us/North\ America/USA/
                           PRE A/
                           PRE B/
                           PRE C/

aws s3api --profile myprofile --region us-east-1 list-objects-v2 --bucket user-test-us --prefix 'North America/USA/' --delimiter /
{
    "CommonPrefixes": [
        {
            "Prefix": "North America/USA/A/"
        },
        {
            "Prefix": "North America/USA/B/"
        },
        {
            "Prefix": "North America/USA/C/"
        }
    ]
}

# 列出对象
aws s3 --profile myprofile --region us-east-1 ls s3://user-test-us/North\ America/USA/ --recursive
2020-12-20 19:50:24          0 North America/USA/A/A/1.txt
2020-12-20 19:50:24          0 North America/USA/B/A/1.txt
2020-12-20 19:50:24          0 North America/USA/C/A/1.txt

aws s3api --profile myprofile --region us-east-1 list-objects-v2 --bucket user-test-us --prefix 'North America/USA/'
{
    "Contents": [
        {
            "Key": "North America/USA/A/A/1.txt",
            "LastModified": "2020-12-20T11:50:24+00:00",
            "ETag": "\"d41d8cd98f00b204e9800998ecf8427e\"",
            "Size": 0,
            "StorageClass": "STANDARD"
        },
        {
            "Key": "North America/USA/B/A/1.txt",
            "LastModified": "2020-12-20T11:50:24+00:00",
            "ETag": "\"d41d8cd98f00b204e9800998ecf8427e\"",
            "Size": 0,
            "StorageClass": "STANDARD"
        },
        {
            "Key": "North America/USA/C/A/1.txt",
            "LastModified": "2020-12-20T11:50:24+00:00",
            "ETag": "\"d41d8cd98f00b204e9800998ecf8427e\"",
            "Size": 0,
            "StorageClass": "STANDARD"
        }
    ]
}
```

#### 复制对象
```shell
aws s3 cp ./local_directory s3://mybucket --recursive

# 将来自S3桶DOC-EXAMPLE-BUCKET1的对象object1复制至S3桶DOC-EXAMPLE-BUCKET2。它使用bucketowner条件来确保bucket由预期帐户根据下表拥有。
aws s3api copy-object --copy-source DOC-EXAMPLE-BUCKET1/object1 \
  --bucket DOC-EXAMPLE-BUCKET2 --key object1copy \
  --expected-source-bucket-owner 111122223333 --expected-bucket-owner 444455556666
```

#### 下载对象
```shell
aws s3api get-object --bucket ${mybucket} --key hello.html hello.html 
{
    "AcceptRanges": "bytes",
    "LastModified": "2020-12-17T10:20:22+00:00",
    "ContentLength": 15,
    "ETag": "\"4643d69d5ca0b263531be4a3f3a76699\"",
    "VersionId": "HlOmQ6Yw_HvbKAjiaykj6vqsDOhQnAuA",
    "ContentType": "text/plain",
    "Metadata": {}
}

# 下载带有DeleteMarker标记的对象,即启用版本控制后被删除的对象
aws s3api get-object --bucket mybucket --key hello.html hello.html
An error occurred (NoSuchKey) when calling the GetObject operation: The specified key does not exist.

# 下载带版本的对象
aws s3api get-object --bucket ${mybucket} --key hello.html hello.html --version-id HlOmQ6Yw_HvbKAjiaykj6vqsDOhQABCD
```

#### 删除对象
```shell
aws s3 rm s3://mybucket/test2.txt

aws s3 rm s3://mybucket --recursive

aws s3 --profile myprofile --region us-east-1 rm s3://user-test-us/North\ America --recursive --exclude *.jpg
aws s3 --profile myprofile --region us-east-1 rm s3://user-test-us/North\ America --recursive --exclude prefix/*
aws s3 --profile myprofile --region us-east-1 rm s3://user-test-us/North\ America --recursive --exclude *tag*
```

#### 生成预签名对象

```shell
# 创建预签名地址时要注意与bucket region相同,profile中任何region都能成功创建URL但只有与bucket相匹配的才能访问
# 预签名有时间限制
#   用户创建:36小时 
#   instance profile: 6小时
#   IAM 用户并且使用 V4版本签名最长可达7天
#   临时Token: Token有效期结束时,预签名URL就到期
aws s3 presign s3://awsexamplebucket/test2.txt --expires-in 86400

```

#### 删除对象
```shell
aws s3api delete-object --bucket ${mybucket} --key h.html

# 带有版本控制的桶删除时不会真的删除对象,只是打一个DeleteMarker标记,获取对象时会返回404, 3hMMqVy.dX0p1HGzoH6_o7cqSkP.wjOH 版本只是一个空对象
aws s3api delete-object --bucket ${mybucket} --key h.html
{
    "DeleteMarker": true,
    "VersionId": "3hMMqVy.dX0p1HGzoH6_o7cqSkP.wjOH"
}

# 删除DeleteMarker标记,让对象回到上一个版本,3hMMqVy.dX0p1HGzoH6_o7cqSkP.wjOH 版本消失
aws s3api delete-object --bucket ${mybucket} --key h.html --version-id 3hMMqVy.dX0p1HGzoH6_o7cqSkP.wjOH
{
    "DeleteMarker": true,
    "VersionId": "3hMMqVy.dX0p1HGzoH6_o7cqSkP.wjOH"
}
```

### 访问点CLI

#### 创建访问点
创建访问点
```shell
aws s3control create-access-point --name example-ap --account-id 123456789012 --bucket example-bucket

aws s3control create-access-point --name example-ap --account-id 123456789012 --bucket example-bucket --public-access-block-configuration BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=true,RestrictPublicBuckets=true
```

创建具有VPC限制的访问点
```shell
aws s3control create-access-point --name example-vpc-ap --account-id 123456789012 --bucket example-bucket --vpc-configuration VpcId=vpc-1a2b3c
```

#### 查看访问点
```shell
aws s3control get-access-point --name example-vpc-ap --account-id 123456789012

{
    "Name": "example-vpc-ap",
    "Bucket": "example-bucket",
    "NetworkOrigin": "VPC",
    "VpcConfiguration": {
        "VpcId": "vpc-1a2b3c"
    },
    "PublicAccessBlockConfiguration": {
        "BlockPublicAcls": true,
        "IgnorePublicAcls": true,
        "BlockPublicPolicy": true,
        "RestrictPublicBuckets": true
    },
    "CreationDate": "2019-11-27T00:00:00Z"
}
```

#### 通过访问点获取对象
```shell
aws s3api get-object --key my-image.jpg --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/prod download.jpg
```

#### 通过访问点上传对象
```shell
aws s3api put-object --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/prod --key my-image.jpg --body my-image.jpg
```

#### 通过访问点删除对象
```shell
aws s3api delete-object --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/prod --key my-image.jpg
```

#### 通过访问点列出对象
```shell
aws s3api list-objects-v2 --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/prod
```

#### 通过访问点向对象添加标签集
```shell
aws s3api put-object-tagging --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/prod --key my-image.jpg --tagging TagSet=[{Key="finance",Value="true"}]
```

#### 使用 ACL 通过访问点授予访问权限
```shell
aws s3api put-object-acl --bucket arn:aws:s3:us-west-2:123456789012:accesspoint/prod --key my-image.jpg --acl private
```








以上完全基于[AWS官方文档](https://docs.aws.amazon.com/),并结合自身理解创作

本文采用知识共享 署名-相同方式共享 3.0协议

署名-相同方式共享(BY-SA):使用者可以对本创作进行转载、节选、混编、二次创作,可以将其运用于商业用途,唯须署名作者,并且采用本创作的内容必须同样采用本协议进行授权
