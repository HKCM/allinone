

### S3访问地址类型

S3支持`虚拟托管类型访问`和`路径类型访问`, 虚拟托管访问是最新的访问方式.

Amazon S3 虚拟托管样式 URL 遵循如下所示格式:  
`https://bucket-name.s3.Region.amazonaws.com/key name`

Amazon S3 路径样式 URL 遵循如下所示格式:  
`https://s3.Region.amazonaws.com/bucket-name/key name`

S3 访问点仅支持虚拟主机式寻址:  
`https://AccessPointName-AccountId.s3-accesspoint.region.amazonaws.com.`

路径样式原本与2020年9月23日废除但被延迟了.

### 存储桶加密
在为存储桶启用默认加密后,将会应用以下加密行为:

* 在启用默认加密之前,存储桶中已存在的对象的加密没有变化。
* 在启用默认加密后上传对象时:
  - 如果您的 PUT 请求标头不包含加密信息,则 Amazon S3 将使用存储桶的默认加密设置来加密对象。
  - 如果您的 PUT 请求标头包含加密信息,则 Amazon S3 将使用 PUT 请求中的加密信息加密对象,然后再将对象存储在 Amazon S3 中。
* 如果您将 SSE-KMS 选项用于默认加密配置,则您将受到 AWS KMS 的 RPS (每秒请求数) 限制

#### 存储桶加密分为3种

1. SSE-S3

使用SSE-S3加密时,每个对象均使用唯一密钥加密。作为额外的保护,它将使用定期轮换的主密钥对密钥本身进行加密。

```shell
aws s3api put-object \
  --bucket text-content \
  --key dir-1/my_images.tar.bz2 \
  --body my_images.tar.bz2 \
  --server-side-encryption AES256
```

<details>
<summary><strong>SSE-S3 Bucket Policy</strong></summary>
```json
{
  "Version": "2012-10-17",
  "Id": "PutObjectPolicy",
  "Statement": [
    {
      "Sid": "DenyIncorrectEncryptionHeader",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::awsexamplebucket1/*",
      "Condition": {
        "StringNotEquals": {
          "s3:x-amz-server-side-encryption": "AES256"
        }
      }
    },
    {
      "Sid": "DenyUnencryptedObjectUploads",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::awsexamplebucket1/*",
      "Condition": {
        "Null": {
          "s3:x-amz-server-side-encryption": "true"
        }
      }
    }
  ]
}
```
</details>

2. SSE-KMS

SSE-KMS 与 SSE-S3 类似,使用该服务具有一些额外的好处,但也要额外收取费用。使用 CMK 需要单独的权限,该密钥可进一步防止未经授权地访问 Amazon S3 中的对象。意味着使用SSE-KMS加密的对象即使公开,如果访问者没有密钥的解密权限也依然无法访问该对象.SSE-KMS 还可以提供审核跟踪,显示 CMK 的使用时间和使用者。此外,还可以创建和管理客户托管 CMK,或者使用您、服务和区域独有的 AWS 托管 CMK。手动定期转换密钥也不要删除和禁用旧的密钥.

**重要: 删除和禁用KMS密钥后,曾被密钥加密的对象将不可读取**

```shell
# 未指定kms-key-id,则使用默认的aws/s3 key
aws s3api put-object \
  --bucket text-content \
  --key my_images.tar.bz2 \
  --body my_images.tar.bz2 \
  --server-side-encryption aws:kms

# 使用kms key arn
aws s3api put-object \
  --bucket text-content \
  --key my_images.tar.bz2 \
  --body my_images.tar.bz2 \
  --server-side-encryption aws:kms \
  --ssekms-key-id arn:aws:kms:us-east-1:123456789012:key/1234abcd-1234-abcd-1234-12345678abcd

# 使用kms key 别名
aws s3api put-object \
  --bucket text-content \
  --key my_images.tar.bz2 \
  --body my_images.tar.bz2 \
  --server-side-encryption aws:kms \
  --ssekms-key-id arn:aws:kms:us-east-1:123456789012:alias/mykey
```

3. SSE-C

使用具有客户提供密钥的服务器端加密 (SSE-C) 时,自己管理加密密钥,而 Amazon S3 管理加密(在它对磁盘进行写入时)和解密(在您访问您的对象时).AWS不会保管用户密钥,用户要自己保存密钥,并自己控制对象和密钥之间的关系.

**重要:密钥丢失则数据丢失**

以下示例使用 Linux 命令行工具生成一个二进制 256 位 AES 密钥,然后将该密钥提供给 Amazon S3 以对上载的文件服务器端进行加密。
```shell
$ dd if=/dev/urandom bs=1 count=32 > sse.key
32+0 records in
32+0 records out
32 bytes (32 B) copied, 0.000164441 s, 195 kB/s
# 上传数据时需要带上用户自定义密钥
$ aws s3api put-object --bucket my-bucket --key mysse.html --body test.txt --sse-customer-key fileb://sse.key --sse-customer-algorithm AES256
{
    "SSECustomerKeyMD5": "iVg8oWa8sy714+FjtesrJg==",
    "SSECustomerAlgorithm": "AES256",
    "ETag": "\"a6118e84b76cf98bf04bbe14b6045c6c\""
}

# 取回数据时也需要带上加密时所使用的密钥
$ aws s3api get-object --bucket my-bucket --sse-customer-key fileb://sse.key --sse-customer-algorithm AES256 --key mysse.html mysse.html
```

#### 加密总结

在第二和第三种加密方式下,对于对象的操作都需要相应地权限.

还有一种客户端加密方式.

要通过单个请求加密现有 Amazon S3 对象,可以使用 Amazon S3 批处理操作。您为 S3 批量操作提供要操作的对象的列表,批量操作调用相应的 API 来执行指定的操作。可以使用复制操作复制现有的未加密对象,并将新的加密对象写入同一存储桶。单个批量操作作业可对数十亿个包含 EB 级数据的对象执行指定操作。

*使用 SSE-KMS 进行默认存储桶加密的 Amazon S3 存储桶不能用作 Amazon S3 服务器访问日志记录 的目标存储桶。对于服务器访问日志目标存储桶,仅支持 SSE-S3 默认加密。*

[我该使用 AWS KMS 管理的密钥还是自定义 AWS KMS 密钥来加密 Amazon S3 中的对象](https://aws.amazon.com/cn/premiumsupport/knowledge-center/s3-object-encrpytion-keys/)

### cors(跨源资源共享)

可以配置存储桶以便允许跨源请求。https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/cors.html

[问题排查](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/cors-troubleshooting.html)

### 事件通知

您可以使存储桶向您发送特定存储桶事件的通知。

有关更多信息,请参阅 配置 Amazon S3 事件通知。


### Outposts
AWS提供设备安装在客户的场所(借给你用),为期三年.提供更低的延迟和本地数据处理需求.可以像使用AWS一样使用这些资源. 主要是贵,无预付三年$194,652,有预付三年$169,444.(跟选择的计算容量有关).还不包含存储...

### 生命周期	
要管理对象以使其在整个生命周期内经济高效地存储,可以配置其 Amazon S3 生命周期。其本质是节约费用.S3 生命周期配置是一组规则,用于定义 Amazon S3 对一组对象应用的操作。小于 128 KB 的对象不会进行存储类型转换(因为没有经济效益).生命周期不支持启用了MFA删除的桶.有两种类型的操作:

* **转换操作** — 定义一段时间后将对象自动转换为另一个存储类。
* **过期操作** — 定义对象的过期时间。Amazon S3 将代表你删除过期的对象。
  

用例:

* 定期将日志上传到一个存储桶,应用程序可能需要使用这些日志一个星期或一个月。之后,可能永远都不需要这些日志。
* 在限定的时间段内可能需要经常访问某些文档。自此之后,这些文档很少被访问。有时,可能不需要对这些文档进行实时访问,但是公司可能要求将它们存档一段特定的时间。之后,可以删除这些文档
* 可以主要为了存档目的而将一些类型的数据上传到 Amazon S3。例如,您可以存档数字媒体、财务和健康记录、原始基因组序列数据、长期数据库备份,以及为遵从法规而必须保留的数据

#### 转换行为

一切转换行为都是对当前版本操作.转换行为可以理解为从高费用向低费用转换,不可逆.例如,可以选择在对象创建 30 天后将其转换为 S3 标准-IA 存储类,或在对象创建 1 年后将其存档到 S3 Glacier 存储类。
可以进行以下转换:
  * 从 S3 标准 存储类转换为任何其他存储类。
  * 任何存储类转换为 S3 Glacier 或 S3 Glacier Deep Archive 存储类。
  * 从 S3 标准-IA 存储类转换为 S3 智能分层 或 S3 单区-IA 存储类。
  * 从 S3 智能分层 存储类转换为 S3 单区-IA 存储类。
  * 从 S3 Glacier 存储类转换为 S3 Glacier Deep Archive 存储类。

#### 过期行为:

一切过期行为都是对当前版本操作.过期操作如果删除了没到存储类最低存储对象的收费期的对象,则按照存储类的最低日期收费,例如IA类最低存储30天,但是在10天就删除了,也按照30天时间收费.Glacier最低90天,Deep Archive最低180天.

* **不受版本控制的存储桶**: 过期操作将永久删除未启用版本控制的存储桶中的对象
* **已启用版本控制的存储桶**: 正常情况Amazon将新建一个具有唯一的版本ID的删除标记的对象.如果当前对象版本是唯一一个且具有删除标记的对象版本,AWS会删除这个标记版本,因为该对象实际已经不存在了.如果该对象具有多个版本且当前版本具有删除标记,则AWS不做任何操作.
* **已暂停版本控制的存储桶**: 正常情况Amazon将新建一个具有null版本ID的删除标记的对象.删除标记会在版本层次结构中将替换版本ID为null的对象版本,从而实际上删除版本标记为null的对象。其他行为与已启用版本控制的存储桶一致

#### 针对过往对象版本的转换和过期行为

转换和过期行为都是在从对象变为非当前时开始的.例如,可以配置一个到期规则,以便在对象变为非当前版本五天后删除非当前版本。例如,假设在 2014 年 1 月 1 日上午 10:30 UTC,您创建了名为 photo.gif 的对象 (版本 ID 111111)。在 2014 年 2 月 1 日上午 11:30 UTC,您意外删除了 photo.gif (版本 ID 111111),这将导致使用新版本 ID (如版本 ID 4857693) 创建一个删除标记。您现在有五天时间可以在永久删除之前,恢复原始版本的 photo.gif (版本 ID 111111)。在 2014 年 1 月 8 日 00:00 UTC,有关过期的生命周期规则执行并永久删除 photo.gif(版本 ID 111111)(在它成为非当前版本五天之后)。

[生命周期与版本控制具体行为](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/intro-lifecycle-rules.html#intro-lifecycle-rules-actions)
#### 注意事项

Amazon S3 按以下方式计算时间:将在规则中指定的天数与对象创建时间相加,然后将得出的时间舍入至下一日的午夜 UTC。例如,如果对象的创建时间是 2014 年 1 月 15 日上午 10:30 UTC,并且您在转换规则中指定了 3 天,则对象的转换日期将计算为 2014 年 1 月 19 日 00:00 UTC。

版本控制示例: Amazon S3 按以下方式计算时间,将规则中指定的天数与创建对象新后继者版本的时间相加,然后将得出的时间舍入至下一日的午夜 UTC。例如,在您的存储桶中,假设某个对象的当前版本的创建时间是 2014 年 1 月 1 日上午 10:30 UTC。如果替换当前版本的对象新版本的创建时间是 2014 年 1 月 15 日上午 10:30 UTC,并且您在转换规则中指定了 3 天,则对象的转换日期计算为 2014 年 1 月 19 日 00:00 UTC。

如果为S3 生命周期操作指定一个过去的日期,所有合格对象会立即符合该生命周期操作的条件,包括新上传的对象.

生命周期配置可以同时存在多个:
* 对相同对象的操作且两个过期策略重叠,将采用较短的过期策略,以便数据存储不会超过预期时间。例如一个策略在30天转化为STANDARD_IA,一个在90天转化为ARCHIVE
* 对相同对象的操作且两个但策略不同. 例如策略一,30天后转化存储类型;策略二30天后删除对象.策略二生效
* 对象筛选存在包含关系,符合第一个规则.例如策略一对象为全体对象10天后类型转换;策略二以log/为前缀的对象365天后类型转化.策略一生效.
* 对象标签冲突.策略一带有tag1/value1标签的对象10天后过期;策略二带有tag2/value2标签的对象100天后转换存储,有对象同时具备2个标签,策略一生效

#### 生命周期日志
CloudTrail 捕获向外部 Amazon S3 终端节点发起的 API 请求,而S3 生命周期操作是使用内部 Amazon S3 终端节点执行的。可以在一个 S3 存储桶中启用 Amazon S3 服务器访问日志,以捕获S3 生命周期相关操作,例如,对象转换到另一个存储类,以及对象失效导致永久删除或逻辑删除。[More Detail](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/lifecycle-and-other-bucket-config.html#lifecycle-general-considerations-mfa-enabled-bucket)


#### [生命周期示例](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/lifecycle-configuration-examples.html)
示例 1:指定筛选条件
示例 2:禁用生命周期规则
示例 3:在对象的生命周期内逐步将存储类降级
示例 4:指定多个规则
示例 5:重叠的筛选条件、冲突的生命周期操作以及 Amazon S3 的功能
示例 6:为启用了版本控制的存储桶指定生命周期规则
示例 7:移除过期对象删除标记
示例 8:用于中止分段上传的生命周期配置

### 位置

创建存储桶时,请指定需要 Amazon S3 在其中创建存储桶的 AWS 区域。Amazon S3 将此信息存储在位置子资源中并提供 API 用于检索此信息。

### 日志记录

日志记录使您可以跟踪针对存储桶的访问请求。每个访问日志记录都提供有关单个访问请求的详细信息,如请求者、存储桶名称、请求时间、请求操作、响应状态和错误代码 (如果有)。访问日志信息可能在安全和访问审核方面十分有用。它还可以帮助您了解您的客户群并了解您的 Amazon S3 账单。 

有关更多信息,请参阅 Amazon S3 服务器访问日志记录。

### 对象锁定

可以​使用对象锁定在固定的时间段内或无限期地阻止删除或覆盖对象.要使用 S3 对象锁定,必须在创建存储桶时启用(目前不支持在创建存储桶之后启用对象锁定)。还可以选择配置将应用于存储桶中放置的新对象的默认保留模式和保留期限。

对象锁定是基于版本控制的,所以具有版本控制功能,它加强了对于删除的保护功能.

* 只能为新存储桶启用 对象锁定。如果要为现有存储桶开启 对象锁定,请联系 AWS Support。
* 如果创建存储桶时启用了 对象锁定,Amazon S3 将自动为该存储桶启用版本控制。
* 如果在创建存储桶时启用了对象锁定,将无法为该存储桶禁用对象锁定或暂停版本控制。

相关权限:
s3:GetObjectRetention
s3:PutObjectLegalHold 
s3:GetObjectLegalHold
s3:GetBucketObjectLockConfiguration
s3:BypassGovernanceRetention

对象锁定仅适用于受版本控制的存储桶,保留期限和依法保留则适用于单个对象版本。当锁定某一对象版本时,Amazon S3 会将锁定信息存储在该对象版本的元数据中。对对象实施保留期限或依法保留仅保护在请求中指定的版本。它不阻止创建该对象的新版本。如果将一个与现有的受保护对象键名相同的对象放在存储桶中,Amazon S3 将创建该对象的新版本、将其存储在请求的存储桶中,并将该请求报告为已成功完成。现有受保护版本的对象将根据其保留配置保持锁定状态。

对象版本可以同时具有保留期限和依法保留、具有其中任何一个或不具有任何一个。

依法保留是一个持续状态直到明确取消,没有保留期限.

S3 对象锁定中的保留期限提供两种保留模式:

* 监管模式
  监管模式允许具有 s3:BypassGovernanceRetention 权限的用户删除,但必须在需要覆盖监管模式的任何请求中显式包含 x-amz-bypass-governance-retention.S3 控制台默认添加bypass-governance-retention: true.
* 合规性模式
  在合规性 模式中,任何用户都不能覆盖或删除受保护的对象版本,包括 AWS 账户中的根用户。在合规性模式下锁定对象后,其保留模式便无法更改,其保留期限也不能缩短。合规性模式确保在保留期限内无法覆盖或删除对象版本。

### 对象标签

可以将 10 个标签与对象关联.键和值区分大小写.标签采用了最终一致性模型。还可以使用标签进行对象分类,就像对象键名前缀分类存储。但是,基于前缀的分类是一维的。参考以下对象键名:
```
photos/photo1.jpg
photos/photo2.jpg
project/project1/document.pdf
project/project2/document2.pdf
```
这些键名具有前缀 photos/、project/project1/ 和 project/project2/。这些前缀支持一维分类。即,一个前缀下的一切都属于一个类别。例如,前缀 project/project1 可确定与项目project1相关的所有文档。

除了数据分类之外,标签还提供了下列好处:

* 对象标签支持权限的精细访问控制。例如,可以向一个 IAM 用户授予仅读取带有特定标签的对象的权限。
* 对象标签支持精细的对象生命周期管理,在其中,除了在生命周期规则中指定键名前缀之外,还可以指定基于标签的筛选条件。
* 使用 Amazon S3 分析时,可以配置筛选条件,以便按对象标签、键名前缀或前缀和标签的组合对对象进行分组以进行分析。
* 还可以自定义 Amazon CloudWatch 指标以按特定标签筛选条件显示信息。

例如可以对对象添加Department, Environment, Project 和 CostCenter等标签.

### 策略 和 ACL(访问控制列表)

所有资源(如存储桶和对象)在默认情况下均为私有。Amazon S3 支持存储桶策略和访问控制列表 (ACL) 选项,供您用于授予和管理存储桶级权限。Amazon S3 将权限信息存储在策略 和 acl 子资源中。

有关更多信息,请参阅 Amazon S3 中的 Identity and Access Management。


### 分段上传

分段上传允许上传单个对象作为一组分段。每个分段都是对象数据的连续部分。可以独立上传以及按任意顺序上传这些对象分段。如果任意分段传输失败,可以重新传输该分段且不会影响其他分段。当对象的所有段都上传后,Amazon S3 将这些段汇编起来,然后创建对象。一般而言,如果您的对象大小达到了 100 MB,您应该考虑使用分段上传,而不是在单个操作中上传对象。

最佳实践是使用 aws s3 命令(例如 aws s3 cp)进行上传和下载,因为这些 aws s3 命令会根据文件大小自动执行分段上传和下载。相比之下,只有在 aws s3 命令不支持特定上传需求(例如,当分段上传涉及多个服务器时,手动停止分段上传并稍后恢复)或者 aws s3 命令不支持所需的请求参数时,才应使用 aws s3api 命令(例如 aws s3api create-multipart-upload)

使用分段上传可提供以下优势:

* 提高吞吐量 - 可以并行上传分段以提高吞吐量。
* 从任何网络问题中快速恢复 - 较小的分段大小可以将由于网络错误而需重启失败的上传所产生的影响降至最低。
* 暂停和恢复对象上传 - 可以在一段时间内逐步上传对象分段。启动分段上传后,不存在过期期限；必须显式地完成或停止分段上传。
* 知道对象的最终大小前开始上传 - 可以在创建对象时将其上传。

### 复制

复制是在相同或跨不同 AWS 区域中的存储桶自动、异步地复制对象。有关更多信息,请参阅复制。

### requestPayment

默认情况下,S3的存储费用和数据流量费用都由S3的拥有者支付.但是当想共享数据,又不希望产生与访问数据等其他操作相关联的费用时,可以将存储桶配置为申请方付款。例如,当使大型数据集 (如邮政编码目录、参考数据、地理空间信息、网盘存储或网络爬取数据) 等,可以使用申请方付款.简单来说请求方支付数据传输和请求方面的费用；存储桶拥有者支付数据存储方面的费用。

配置为申请方付款存储桶后,请求方必须在其请求中包含 x-amz-request-payer (在 POST、GET 和 HEAD 请求的标头中,或在 REST 请求中作为参数),以显示他们明确知道请求和数据下载将产生费用。requestPayment是存储桶级别的,不能对特定对象使用.

但是以下行为依然会向S3拥有者收费,虽然他们的请求失败了:
* 申请方未在标头中 (GET、HEAD 或 POST) 包含参数 x-amz-request-payer,或未在请求中将其作为参数 (REST) (HTTP 代码 403)。
* 请求身份验证失败 (HTTP 代码 403)。
* 请求是匿名的 (HTTP 代码 403)。
* 请求是 SOAP 请求。

有关更多信息,请参阅 申请方付款存储桶。

### 访问点

Amazon S3 访问点简化了 S3 中共享数据集的大规模数据访问管理。访问点是附加到存储桶的命名网络终端节点,您可以使用这些存储桶执行 S3 对象操作(如 GetObject 和 PutObject)。每个访问点都具有不同的权限和网络控制,S3 将它们应用于通过该访问点发出的任何请求。每个访问点强制实施自定义访问点策略,该策略与附加到底层存储桶的存储桶策略结合使用。可以将任何访问点配置为仅接受来自 Virtual Private Cloud (VPC) 的请求,以限制专用网络的 Amazon S3 数据访问。还可以为每个访问点配置自定义阻止公有访问设置。使用访问点来执行对象操作。不能使用访问点执行其他 Amazon S3 操作,例如修改或删除存储桶。

对于通过访问点发出的任何请求,Amazon S3 会评估该访问点、底层存储桶和存储桶拥有者账户的阻止公有访问设置。如果这些设置中的任何一个指示应阻止请求,则 Amazon S3 会拒绝请求。

向存储桶添加 S3 访问点不会更改存储桶在通过现有存储桶名称或 ARN 访问时的行为。针对存储桶的所有现有操作将继续像以前一样运行。访问点策略中的限制仅适用于通过该访问点发出的请求。

[访问点与 S3 操作和 AWS 服务的兼容性](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/using-access-points.html#access-points-service-api-support)


[访问点的限制](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/access-points-restrictions-limitations.html)

### 传输加速Transfer Acceleration

Transfer Acceleration 可在客户端与 S3 存储桶之间实现快速、轻松、安全的远距离文件传输。Transfer Acceleration 利用 Amazon CloudFront 的全球分布式边缘站点。启用后有最大0.08 USD/GB的双向流量收费.

什么时候启用加速:
* 位于全球各地的客户需要上传到集中式存储桶
* 定期跨大洲传输数 GB 至数 TB 数据
* 在上传到 Amazon S3 时无法充分利用 Internet 上的所有可用带宽

[加速效果比较](https://s3-accelerate-speedtest.s3-accelerate.amazonaws.com/en/accelerate-speed-comparsion.html)

借助传输加速终端节点执行所有 Amazon S3 操作,但以下操作除外:`GET 服务(列出存储桶)`、`PUT 存储桶(创建存储桶)`和 `DELETE 存储桶`。此外,Amazon S3 Transfer Acceleration不支持使用 `PUT Object - Copy` 进行跨区域复制。

启用加速后,加速终端节点将随即可用,但最多 20 分钟后才可实现性能提升.使用加速终端节点才能有性能提升.加速终端节点名称类似:
* bucketname.s3-accelerate.amazonaws.com
* bucketname.s3-accelerate.dualstack.amazonaws.com

重新使用标准上传速度,只需将名称更改回 `mybucket.s3.us-east-1.amazonaws.com`

#### CLI示例

1. 下示例设置 Status=Enabled 以对存储桶启用 Transfer Acceleration。可使用 Status=Suspended 暂停 Transfer Acceleration
```shell
aws s3api put-bucket-accelerate-configuration --bucket bucketname --accelerate-configuration Status=Enabled
```

2. 下示例在默认配置文件中将 use_accelerate_endpoint 设置为 true
```shell
aws configure set default.s3.use_accelerate_endpoint true
```

3. 需要对某些 AWS CLI 命令使用加速终端节点,但不对其他此类命令使用加速终端节点,则可使用以下两种方法中的任一方法:
  * 通过将任何 s3 或 s3api 命令的 --endpoint-url 参数设置为 https://s3-accelerate.amazonaws.com 或 http://s3-accelerate.amazonaws.com 来对每条命令使用加速终端节点
  * 可以在 AWS Config 文件中设置单独的配置文件。例如,创建一个将 use_accelerate_endpoint 设置为 true 的配置文件和一个不设置 use_accelerate_endpoint 的配置文件。在运行一条命令时,根据是否需要使用加速终端节点来指定要使用的配置文件

4. 将对象上传到已启用 Transfer Acceleration 的存储桶的 AWS CLI 示例
    以下示例通过使用已配置为使用加速终端节点的默认配置文件来将文件上传到已启用 Transfer Acceleration 的存储桶。
  ```shell
  aws s3 cp file.txt s3://bucketname/keyname --profile Enable_Acceleration
  ```
  以下示例通过使用 --endpoint-url 参数指定加速终端节点来将文件上传到已启用 Transfer Acceleration 的存储桶。
  ```shell
  aws configure set s3.addressing_style virtual
  aws s3 cp file.txt s3://bucketname/keyname --region region --endpoint-url http://s3-accelerate.amazonaws.com
  ```


### 版本控制	

[版本控制详情](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/Versioning.html)

版本控制可帮助恢复意外覆盖和删除。一旦对存储桶启用了版本控制,它将无法返回到无版本控制状态。但是,可以在该存储桶上暂停版本控制。存储桶可处于以下三种状态之一:非版本化 (默认) 、启用版本控制或暂停版本控制.启用版本控制,Amazon S3 会自动为要存储的对象生成唯一版本 ID。例如,在一个存储桶中,可以拥有两个具有相同键的对象,但版本 ID 却不同,例如 photo.gif (版本 111111) 和 photo.gif (版本 121212)。

版本控制状态将应用到该存储桶中的所有 (不是某些) 对象。第一次对存储桶启用版本控制后,该存储桶中的对象将在之后一直受版本控制,并具有唯一的版本 ID。请注意以下几点:

* 在设置版本控制状态之前存储在存储桶中的对象具有版本 ID null。启用版本控制时,存储桶中的现有对象不会更改。更改的是 Amazon S3 在以后的请求中处理这些对象的方式。
* 存储桶拥有者 (或任何具有适当权限的用户) 可以暂停版本控制以停止累积对象版本。暂停版本控制时,存储桶中的现有对象不会更改。更改的是 Amazon S3 在以后的请求中处理对象的方式。

通过启用了版本控制的存储桶可以恢复因意外删除或覆盖操作而失去的对象。例如:

* 如果删除对象 (而不是永久移除它),则 Amazon S3 会插入删除标记,并创建新的版本作为当前版本,访问时会有noSuchKey。始终可以恢复以前的版本。
* 如果覆盖对象,则会导致存储桶中出现新的对象版本。始终可以恢复以前的版本。

版本控制的其中一个价值主张是能够检索对象的早期版本。有两种方法可执行该操作:

* 将对象的早期版本复制到同一存储桶中
  复制的对象将成为该对象的当前版本,且所有对象版本都保留。
* 永久删除对象的当前版本
  当您删除当前对象版本时,实际上会将先前版本转换为该对象的当前版本。

在存储桶上暂停版本控制后,Amazon S3 会自动将后续对象的存储操作(Put,Post,Copy)添加 null 版本 ID,如果存储桶原本就具有版本ID为 null的对象,那么这个对象将被覆盖

在存储桶上暂停版本控制后,DELETE 请求:

* 可以仅删除其版本 ID 为 null 的对象
  如果存储桶中没有对象的空版本,则不删除任何内容。
* 将删除标记插入到存储桶。
  Amazon S3 会插入删除标记,并创建新的版本作为当前版本,访问时会有noSuchKey。

**永久删除启用版本控制的对象需要在DELETE 请求带上version ID**

版本对象最好与生命周期结合使用.版本控制会收取所有版本的存储费用.

#### MFA 删除

所有授权 IAM 用户都可以启用版本控制,但只有存储桶拥有者 (根账户) 才能启用 MFA 删除。 启用后`更改存储桶的版本控制状态`以及`永久删除对象版本`需要MFA认证.更安全.

[版本控制详解](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/Versioning.html)
[启用版本控制后,Amazon S3 对存储桶请求的 HTTP 503 响应显著增加](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/troubleshooting.html#troubleshooting-by-symptom-increase-503-reponses)

### 网站

可以配置存储桶以便用于静态网站托管。Amazon S3 通过创建网站 子资源来存储此配置。

### 账单和使用率报告

https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/dev/BucketBilling.html



### 结语

以上完全基于[AWS官方文档](https://docs.aws.amazon.com/),并结合自身理解创作

