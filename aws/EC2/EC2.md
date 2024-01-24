# AWS EC2

## AMI

AMI(Amazon Machine Image) 提供启动实例所需的信息,包括一个或多个 EBS Snapshot. AMI 可以创建、购买、共享、出售和注销. 以下都忽略了实例存储 AMI.

### AMI 类型

根据启动许可分为以下类别:

| 启动权限 | 描述                                |
| :------- | :---------------------------------- |
| 公有     | 拥有者向所有 AWS 账户授予启动许可。 |
| 显式     | 拥有者向特定 AWS 账户授予启动许可。 |
| 隐式     | 拥有者拥有 AMI 的隐式启动许可。     |

根据根设备存储可以分为由 Amazon EBS 支持或由实例存储支持:

| 特征             | 由 Amazon EBS 支持的 AMI                                     | 由 Amazon 实例存储支持的 AMI                       |
| :--------------- | :----------------------------------------------------------- | :------------------------------------------------- |
| 实例的启动时间   | 通常不到 1 分钟                                              | 通常不到 5 分钟                                    |
| 根设备的大小限制 | 16 TiB                                                       | 10 GiB                                             |
| 根设备卷         | EBS 卷                                                       | 实例存储卷                                         |
| 数据持久性       | 默认情况下,实例终止时将删除根卷。* 默认情况下,在实例终止后,任何其他 EBS 卷上的数据仍然存在。 | 任意实例存储卷上的数据仅在实例的生命周期内保留。   |
| 修改             | 实例停止后,实例类型、内核、RAM 磁盘和用户数据仍可更改。     | 实例存在期间,实例属性是稳定不变的。               |
| 收费             | 您需要为实例使用、EBS 卷使用以及将 AMI 存储为 EBS 快照付费。 | 您需要为实例使用以及在 Amazon S3 中存储 AMI 付费。 |
| AMI 创建/捆绑    | 使用单一命令/调用                                            | 需要安装和使用 AMI 工具                            |
| 停止状态         | 可以处于停止状态。即使实例停止未运行,根卷也会保留在 Amazon EBS 中 | 不可置于停止状态；实例正在运行或已终止             |

根据虚拟化类型可分为半虚拟化 (PV) 或硬件虚拟机 (HVM).



### 查找 AMI

[查找 Linux AMI](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/finding-an-ami.html)

#### 使用 AWS CLI 查找 AMI

```shell
aws ec2 describe-images --owners self amazon \
	--filters "Name=root-device-type,Values=ebs"
```

#### 使用公有参数启动实例

在本示例中,不包括 `--count` 和 `--security-group` 参数。对于 `--count`,默认为 1。如有默认 VPC 和默认安全组,则将使用它们。

```shell
aws ec2 run-instances 
    --image-id resolve:ssm:/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2 
    --instance-type m5.xlarge 
    --key-name MyKeyPair
```

#### 查找快速启动 AMI

**例 示例:查找当前 Amazon Linux 2 AMI**

```shell
aws ec2 describe-images \
    --owners amazon \
    --filters "Name=name,Values=amzn2-ami-hvm-2.0.????????.?-x86_64-gp2" "Name=state,Values=available" \
    --query "reverse(sort_by(Images, &Name))[:1].ImageId" \
    --output text
```

**例 示例:查找当前 Amazon Linux AMI**

```shell
aws ec2 describe-images \
    --owners amazon \
    --filters "Name=name,Values=amzn-ami-hvm-????.??.?.????????-x86_64-gp2" "Name=state,Values=available" \
    --query "reverse(sort_by(Images, &Name))[:1].ImageId" \
    --output text
```

**例 示例:查找当前 Ubuntu Server 16.04 LTS AMI**

```shell
aws ec2 describe-images \
    --owners 099720109477 \
    --filters "Name=name,Values=ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-????????" "Name=state,Values=available" \
    --query "reverse(sort_by(Images, &Name))[:1].ImageId" \
    --output text
```

**例 示例:查找当前 Red Hat Enterprise Linux 7.5 AMI**

```shell
aws ec2 describe-images \
    --owners 309956199498 \
    --filters "Name=name,Values=RHEL-7.5_HVM_GA*" "Name=state,Values=available" \
    --query "reverse(sort_by(Images, &Name))[:1].ImageId" \
    --output text
```

**例 示例:查找当前 SUSE Linux Enterprise Server 15 AMI**

```shell
aws ec2 describe-images \
    --owners amazon \
    --filters "Name=name,Values=suse-sles-15-v????????-hvm-ssd-x86_64" "Name=state,Values=available" \
    --query "reverse(sort_by(Images, &Name))[:1].ImageId" \
    --output text
```

#### 查找共享 AMI

**示例:列出所有公用 AMI**

以下命令将列出所有公用 AMI,包括您拥有的所有公用 AMI。

```shell
aws ec2 describe-images --executable-users all
```

**示例:使用显式启动许可列出 AMI**

以下命令列出您对其拥有显式启动许可的 AMI。此列表不包括您拥有的任何 AMI。

```shell
aws ec2 describe-images --executable-users self
```

**示例:列出 Amazon 拥有的 AMI**

以下命令列出 Amazon 拥有的 AMI。Amazon 的公用 AMI 的拥有者有一个别名,在账户字段中显示为 `amazon`。这使您可以轻松地从 Amazon 查找 AMI。其他用户不能对其 AMI 使用别名。

```shell
aws ec2 describe-images --owners amazon
```

**示例:列出账户拥有的 AMI**

以下命令列出指定 AWS 账户拥有的 AMI。

```shell
aws ec2 describe-images --owners 123456789012
```

**示例:使用筛选条件确定 AMI 的范围**

要减少显示的 AMI 数量,请使用筛选条件只列出您感兴趣的 AMI 类型。例如,使用以下筛选条件可以只显示 EBS 支持的 AMI。

```shell
--filters "Name=root-device-type,Values=ebs"
```

### 共享 AMI

[共享 AMI 注意事项](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/usingsharedamis-finding.html#usingsharedamis-confirm)

[共享 AMI 指导原则](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/building-shared-amis.html)

#### 将 AMI 设为公有

**将 AMI 设为公用**

1. 使用 [modify-image-attribute](https://docs.aws.amazon.com/cli/latest/reference/ec2/modify-image-attribute.html) 命令可将 `all` 组添加到指定 AMI 的 `launchPermission` 列表,如下所示。

   ```shell
   aws ec2 modify-image-attribute \    --image-id ami-0abcdef1234567890 \    --launch-permission "Add=[{Group=all}]"
   ```

2. 要验证 AMI 的启动许可,请使用 [describe-image-attribute](https://docs.aws.amazon.com/cli/latest/reference/ec2/describe-image-attribute.html) 命令。

   ```shell
   aws ec2 describe-image-attribute \    --image-id ami-0abcdef1234567890 \    --attribute launchPermission
   ```

3. (可选)要再次将 AMI 设为私有,请从其启动许可中删除 `all` 组。请注意,AMI 的拥有者始终具有启动许可,因此不受该命令影响。

   ```shell
   aws ec2 modify-image-attribute \    --image-id ami-0abcdef1234567890 \    --launch-permission "Remove=[{Group=all}]"
   ```

#### 特定 AWS 共享共享 AMI

参考以下博文:
	[如何在多个账户中共享加密 AMI,以启动加密 EC2 实例](https://aws.amazon.com/cn/blogs/china/how-to-share-encrypted-amis-across-accounts-to-launch-encrypted-ec2-instances/)

**要授予显式启动许可**

以下命令向指定 AWS 账户授予指定 AMI 的启动许可。

```shell
aws ec2 modify-image-attribute \    --image-id ami-0abcdef1234567890 \    --launch-permission "Add=[{UserId=123456789012}]"
```

以下命令为快照授予创建卷的权限。

```shell
aws ec2 modify-snapshot-attribute \    --snapshot-id snap-1234567890abcdef0 \    --attribute createVolumePermission \    --operation-type add \    --user-ids 123456789012
```

**注意**

```
您不需要为了共享 AMI 而共享 AMI 引用的 Amazon EBS 快照。只需共享 AMI 本身；系统自动为实例提供访问所引用 Amazon EBS 快照的权限以便启动。不过,您确实需要共享用于对 AMI 引用的快照加密的所有 KMS 密钥。有关更多信息,请参阅[共享 Amazon EBS 快照](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ebs-modifying-snapshot-permissions.html)。
```



**要删除账户的启动许可**

以下命令从指定 AWS 账户中删除指定 AMI 的启动许可:

```shell
aws ec2 modify-image-attribute \    --image-id ami-0abcdef1234567890 \    --launch-permission "Remove=[{UserId=123456789012}]"
```

以下命令为快照授予删除卷的权限。

```shell
aws ec2 modify-snapshot-attribute \    --snapshot-id snap-1234567890abcdef0 \    --attribute createVolumePermission \    --operation-type remove \    --user-ids 123456789012
```

**要删除所有的启动许可**

以下命令从指定 AMI 中删除所有公用和显式启动许可。请注意,AMI 的拥有者始终具有启动许可,因此不受该命令影响。

```shell
aws ec2 reset-image-attribute \    --image-id ami-0abcdef1234567890 \    --attribute launchPermission
```

**为您的 AMI 创建书签**

键入一个带有以下信息的 URL,其中 *region* 表示您的 AMI 驻留的区域:

```
https://console.aws.amazon.com/ec2/v2/home?region=region#LaunchInstanceWizard:ami=ami_id
```

例如,此 URL 从 us-east-1 区域内的 ami-0abcdef1234567890 AMI 启动实例:

```
https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#LaunchInstanceWizard:ami=ami-0abcdef1234567890
```

### 创建 AMI

创建 AMI 可以从正在运行的 Instance 创建也可以从 Snapshot 创建. 具体请看: [创建 AMI](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/create-ami.html)

### 复制 AMI

#### 跨区域复制

跨不同地理位置复制 AMI 具有以下优势:

- 一致的全球部署:通过将 AMI 从一个区域复制到另一个区域,您可以根据相同的 AMI 在不同的区域中启动一致的实例。
- 可扩展性:无论用户身处何处,您都可以更轻松地设计和构建能满足他们需求的全球应用程序。
- 性能:您可以通过分发您的应用程序以及找到较接近您用户的应用程序的关键组件来提高性能。您还可以利用区域特定的功能,例如,实例类型或其他 AWS 服务。
- 高可用性:您可以跨 AWS 区域设计和部署应用程序以提高可用性。

下图显示源 AMI、在不同的区域中复制的两个 AMIs 以及从它们中启动的 EC2 实例之间的关系。从 AMI 中启动实例时,该实例位于 AMI 所在的区域中。如果您更改源 AMI,并希望在目标区域中的 AMIs 上反映这些更改,您必须将源 AMI 重新复制到目标区域中。

#### 跨账号复制AMI

跨账号复制就是与特定 AWS 账号共享,详情参考:

 [跨账号复制 AMI](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/CopyingAMIs.html#copy-ami-across-accounts)

[将加密与 EBS 支持的 AMI 结合使用](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/AMIEncryption.html)

主要要注意:如果 AMI 有加密,源 AMI 拥有者还必须共享 KMS 秘钥.

#### 使用共享AMI注意事项:
https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/usingsharedamis-finding.html#usingsharedamis-confirm
如果共享 AMI 带有加密快照,拥有者必须同时与您共享一个或多个密钥。

#### 如何共享AMI:
https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/sharingamis-explicit.html#sharingamis-console

#### AMI共享URL

使用带有以下信息的 URL可以让被共享方直接使用,其中 region 表示您的 AMI 驻留的区域:
https://console.aws.amazon.com/ec2/v2/home?region=`region`#LaunchInstanceWizard:ami=`ami_id` 

#### Amazon Data Lifecycle Manager (Amazon DLM)

能够定期
博客:
https://aws.amazon.com/cn/blogs/storage/automating-amazon-ebs-snapshot-and-ami-management-using-amazon-dlm/

#### 使用 S3 存储和还原 AMI
复制AMI到其他Region后,应更新任何数据库连接字符串或相似的应用程序配置数据,以指向适当的资源。否则,从目标区域上的新 AMI 中启动的实例可能仍会使用源区域中的资源,这可能会影响性能和成本。

将 AMI 从一个 AWS 分区复制到另一个 AWS 分区,AMI 存储和还原 API 的工作原理
要使用 S3 存储和还原 AMI,请使用以下 API:

##### CreateStoreImageTask - 将 AMI 存储在 S3 存储桶中

API 创建一个任务,从 AMI 及其快照中读取所有数据,然后使用 S3 分段上传将数据存储在 S3 对象中。API 获取 AMI 的所有组件,包括大多数非区域特定的 AMI 元数据以及 AMI 中包含的所有 EBS 快照,然后将它们打包到 S3 内的单个对象中。数据将作为上传流程的一部分进行压缩,以减少 S3 中使用的空间量,因此 S3 中的对象可能小于 AMI 中快照大小总和。

##### DescribeStoreImageTasks - 提供 AMI 存储任务的进度

DescribeStoreImageTasks API 描述 AMI 存储任务的进度。您可以描述指定 AMI 的任务。如果未指定 AMI,则会获得过去 31 天内处理的所有存储映像任务的分页列表。

对于每个 AMI 任务,响应会指示任务是 InProgress、Completed 还是 Failed。对于任务 InProgress,响应会将估计进度显示为百分比值。

任务按反向的时间顺序列出。

目前,只能查看上个月的任务。

##### CreateRestoreImageTask - 从 S3 存储桶还原 AMI
https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ami-store-restore.html#how-it-works

### 注销 AMI

1. **注销 AMI**

   使用 [deregister-image](https://docs.aws.amazon.com/cli/latest/reference/ec2/deregister-image.html) 命令注销 AMI:

   ```shell
   aws ec2 deregister-image --image-id ami-12345678
   ```

2. **删除不再需要的快照**

   使用 [delete-snapshot](https://docs.aws.amazon.com/cli/latest/reference/ec2/delete-snapshot.html) 命令删除不再需要的快照:

   ```shell
   aws ec2 delete-snapshot --snapshot-id snap-1234567890abcdef0
   ```

3. **终止实例(可选)**

   如果您使用完从 AMI 启动的实例,则可以使用 [terminate-instances](https://docs.aws.amazon.com/cli/latest/reference/ec2/terminate-instances.html) 命令终止该实例:

   ```shell
   aws ec2 terminate-instances --instance-ids i-12345678
   ```



## Instance

### 实例类型

实例大致可以分为 5 类: 通用型,计算优化型,内存优化型,存储优化型以及加速计算.

详情:

​	[通用型适用范围](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/general-purpose-instances.html)

​	[计算优化型适用范围](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/compute-optimized-instances.html)

​	[内存优化型适用范围](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/memory-optimized-instances.html)

​	[存储优化型适用范围](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/storage-optimized-instances.html)

​	[加速计算](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/accelerated-computing-instances.html)

**使用 AWS CLI 查找实例类型**

1. 使用 [describe-instance-types](https://docs.aws.amazon.com/goto/aws-cli/ec2-2016-11-15/DescribeInstanceTypes) 命令根据实例属性筛选实例类型。例如,可以使用以下命令以仅显示具有 48 个 vCPU 的实例类型。

   ```shell
   aws ec2 describe-instance-types --filters "Name=vcpu-info.default-vcpus,Values=48"
   ```

2. 使用 [describe-instance-type-offerings](https://docs.aws.amazon.com/goto/aws-cli/ec2-2016-11-15/DescribeInstanceTypeOfferings) 命令筛选按位置(区域或可用区)提供的实例类型。例如,可以使用以下命令以显示在指定的可用区中提供的实例类型。

   ```shell
   aws ec2 describe-instance-type-offerings --location-type "availability-zone" --filters Name=location,Values=us-east-2a --region us-east-2
   ```



### 购买选项

Amazon EC2 提供了以下让您根据需求优化成本的购买选项:

- **按需实例** - 按秒为启动的实例付费。
- **Savings Plans** - 通过承诺在 1 年或 3 年期限内保持一致的使用量(以 USD/小时为单位)来降低您的 Amazon EC2 成本。
- **预留实例** - 通过承诺在 1 年或 3 年期限内提供一致的实例配置(包括实例类型和区域)来降低您的 Amazon EC2 成本。
- **Spot 实例** - 请求未使用的 EC2 实例,这可能会显著降低您的 Amazon EC2 成本。
- **专用主机** - 为完全专用于运行您的实例的物理主机付费,让您现有的按插槽、按内核或按 VM 计费的软件许可证降低成本。
- **专用实例** - 为在单一租户硬件上运行的实例按小时付费。
- **预留容量** - 可在特定可用区中为 EC2 实例预留容量,持续时间不限。

**使用AWS CLI来确定实例的生命周期。**

使用以下[描述实例](https://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html)口令:

```
aws ec2 describe-instances --instance-ids i-1234567890abcdef0
```

如果实例正在专用主机上运行,那么输出内容包含以下信息:

```
"Tenancy": "host"
```

如果实例为专用实例,那么输出内容包含以下信息:

```
"Tenancy": "dedicated"
```

如果实例为 Spot 实例,那么输出内容包含以下信息:

```
"InstanceLifecycle": "spot"
```

#### On-Demand

使用 按需实例,按秒为计算容量支付费用,但无需作出长期承诺。可以完全控制其生命周期 — 您确定何时发布、停止、休眠、启动、重启或终止它。

购买按需实例没有长期承诺。只需要为处于 `running` 状态的按需实例的秒数付费。运行中的按需实例的每秒的价格是固定的

#### Reserved Instances

相比按需实例定价,预留实例可大幅节约 Amazon EC2 成本.

##### 决定Reserved Instance定价的关键变量

**实例属性**

预留实例有四个决定其价格的实例属性。

- **实例类型**:例如,`m4.large`。这由实例系列(例如,`m4`)和实例大小(例如,`large`)组成。
- **区域**:购买Reserved Instance的区域。
- **租赁**:您的实例在共享 (默认) 还是单租户 (专用) 硬件上运行。有关更多信息,请参阅[Dedicated Instances](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/dedicated-instance.html)。
- **平台**:操作系统；例如,Windows 或 Linux/Unix。有关更多信息,请参阅[选择平台](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ri-market-concepts-buying.html#ri-choosing-platform)。

**期限承诺** 

您可以承诺购买一年或三年的 Reserved Instance,三年承诺可以获得更大的折扣。

- **一年**:一年定义为 31536000 秒(365 天)。
- **三年**:三年定义为 94608000 秒 (1095 天)。

**付款选项** 

针对预留实例可使用以下付款选项:

- **预付全费**:所有款项于期限开始时支付,无论使用了多少小时数,剩余期限不会再产生其他任何费用或额外按小时计算的费用。
- **预付部分费用**:必须预付部分费用,无论是否使用了 Reserved Instance,期限内剩余的小时数都将按照打折小时费率计费。
- **无预付费用**:无论是否使用 Reserved Instance,您都将按照期限内的小时数,采用打折小时费率进行付费。无需预付款。

**优惠类别** 

在计算需求发生变化时,您可以根据产品类别修改或交换 Reserved Instance。

- **标准**:这些提供最大力度的折扣,但只可以修改。标准预留实例无法交换。
- **可转换**:这些相较于标准预留实例提供较低的折扣,但可以与具有不同实例属性的可转换预留实例进行交换。可转换预留实例也可修改。

**总结:** 标准类别(指定 Instance Family)三年期全款折扣最高,并且标准类别可以出售,可转换类别不可出售

##### 购买和销售 Reserved Instance

[购买 Reserved Instance](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ri-market-concepts-buying.html)

Reserved Instance 还可以出售,不过需要注册为卖家: [销售 Reserved Instance](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ri-market-general.html)

##### 修改 Reserved Instance

可以修改全部或部分预留实例。可以将原始预留实例分为两个或更多新的预留实例。

例如,如果在 `us-east-1a` 中有 10 个实例的预留,并决定将其中 5 个实例移至 `us-east-1b`,则修改请求会生成两个新的预留实例 - 一个用于 `us-east-1a` 中的 5 个实例,另一个用于 `us-east-1b` 中的 5 个实例。

还可以将两个或更多预留实例*合并*成单个预留实例。例如,如果有四个 `t2.small` 均为 预留实例,则可以将其合并以创建单个 `t2.large` Reserved Instance。

#### Spot Instances

Spot 实例是一种未使用的 EC2 实例,以低于按需价格提供。由于 Spot 实例 允许以极低的折扣请求未使用的 EC2 实例,这可能会显著降低 Amazon EC2 成本,一般能节约 60%以上。

**Spot 实例与 按需实例 的主要区别**

- 只有 Spot 请求处于活动状态并且有可用容量时才能立即启动。
- 如果没有可用容量,则 Spot 请求会继续自动发起启动请求,直到有可用容量为止。
- Spot 实例的每小时价格根据需求而有所不同。
- 当实例处于较高的中断风险时,Amazon EC2 为正在运行的 Spot 实例发出信号。
- 如果容量不再可用、Spot 价格超出您的最高价或者对 Spot 实例的需求增加,Amazon EC2 Spot 服务可能[中断](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/spot-interruptions.html)个别 Spot 实例。

##### Spot 实例中断

Spot 实例可能会中断。因此,必须确保应用程序针对 Spot 实例中断做好准备.

**中断原因** 

下面列出了 Amazon EC2 中断您的 Spot 实例 的可能原因:

- 价格 - Spot 价格高于您的最高价。
- 容量 - 如果没有足够的未使用 EC2 实例,无法满足按需实例的需求,则 Amazon EC2 会中断 Spot 实例。实例的中断顺序是由 Amazon EC2 确定的。
- 约束 - 如果您的请求包含约束(如启动组或可用区组),则当不再满足约束时,这些 Spot 实例将成组终止。



##### Spot 实例再平衡建议

EC2 实例*再平衡建议*是一个新信号,可在 Spot 实例处于较高的中断风险时通知您。信号可能比该[两分钟的 Spot 实例中断通知](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/spot-interruptions.html#spot-instance-termination-notices)更早到达,从而让您有机会主动管理 Spot 实例。您可以决定将工作负载再平衡到不处于较高中断风险的新的或现有的 Spot 实例。

Amazon EC2 并不总能在两分钟的 Spot 实例中断通知之前发送再平衡建议信号。因此,再平衡建议信号可能会随两分钟的中断通知一起到达。

#### 容量预留(坑)

开始按需容量预留就会开始收费,无论是否在预留容量中运行实例,都按等同的按需费率计算费用。

如果没有使用预留,这将在您的 EC2 账单中显示为未使用的预留。如果运行的实例属性与预留匹配,则只需要为该实例付费,不需要为预留付费。没有任何预付费用或额外收费。

**限制**

在创建容量预留之前,请注意以下限制。

- 活动和未使用的容量预留会计入您的按需实例限制中。
- 容量预留无法从一个AWS账户转移到另一个账户。但是,您可以与其他AWS账户共享容量预留。有关更多信息,请参阅。[使用共享 容量预留](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/capacity-reservation-sharing.html)
- 区域Reserved Instance账单折扣不适用于容量预留。
- 无法在置放群组中创建容量预留。
- 容量预留不能与专用主机一起使用。
- 容量预留不能确保休眠的实例在尝试启动后可以恢复。

### Instance 生命周期

下表提供了每个实例状态的简短说明,并指示它是否已计费。

**注意**

该表仅指示用于实例使用率的计费。一些 AWS 资源(如 Amazon EBS 卷和弹性 IP 地址)无论实例的状态如何,都将产生费用。有关更多信息,请参阅*AWS Billing and Cost Management 用户指南*中的[避免意外费用](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/checklistforunwantedcharges.html)。

| 实例状态        | 描述                                                         | 实例使用率计费                             |
| :-------------- | :----------------------------------------------------------- | :----------------------------------------- |
| `pending`       | 实例正准备进入 `running` 状态。实例在首次启动时进入 `pending` 状态,或者在处于 `stopped` 状态后启动。 | 不计费                                     |
| `running`       | 实例正在运行,并且做好了使用准备。                           | 已计费                                     |
| `stopping`      | 实例正准备处于停止状态或休眠停止状态。                       | 如果准备停止,则不计费如果准备休眠,则计费 |
| `stopped`       | 实例已关闭,不能使用。可随时启动实例。                       | 不计费                                     |
| `shutting-down` | 实例正准备终止。                                             | 不计费                                     |
| `terminated`    | 实例已永久删除,无法启动。                                   | 不计费                                     |

`重启`不会变更公有 IPv4 地址,但是`停止`会. 

#### 启动 Instance

[向导启动](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/launching-instance.html)

[模板启动](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ec2-launch-templates.html)

[现有实例启动](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/launch-more-like-this.html)

#### 连接 Instance

**获取用于启动实例的 AMI 的默认用户名。**

- 对于 Amazon Linux 2 或 Amazon Linux AMI,用户名是 `ec2-user`。
- 对于 CentOS AMI,用户名是 `centos`。
- 对于 Debian AMI,用户名称是 `admin`。
- 对于 Fedora AMI,用户名是 `ec2-user` 或 `fedora`。
- 对于 RHEL AMI,用户名是 `ec2-user` 或 `root`。
- 对于 SUSE AMI,用户名是 `ec2-user` 或 `root`。
- 对于 Ubuntu AMI,用户名称是 `ubuntu`。

**获取实例指纹**

1. 在本地计算机上(而不是在实例上),按以下方式使用 [get-console-output](https://docs.aws.amazon.com/cli/latest/reference/ec2/get-console-output.html) (AWS CLI) 命令以获取指纹:

   ```shell
   aws ec2 get-console-output --instance-id instance_id --output text
   ```

2. 以下是您应该在输出中查找的内容的示例。确切的输出可能因操作系统、AMI 版本以及是否让 AWS 创建密钥而异。

   ```shell
   ec2: #############################################################
   ec2: -----BEGIN SSH HOST KEY FINGERPRINTS-----
   ec2: 1024 SHA256:7HItIgTONZ/b0CH9c5Dq1ijgqQ6kFn86uQhQ5E/F9pU root@ip-10-0-2-182 (DSA)
   ec2: 256 SHA256:l4UB/neBad9tvkgJf1QZWxheQmR59WgrgzEimCG6kZY root@ip-10-0-2-182 (ECDSA)
   ec2: 256 SHA256:kpEa+rw/Uq3zxaYZN8KT501iBtJOIdHG52dFi66EEfQ no comment (ED25519)
   ec2: 2048 SHA256:L8l6pepcA7iqW/jBecQjVZClUrKY+o2cHLI0iHerbVc root@ip-10-0-2-182 (RSA)
   ec2: -----END SSH HOST KEY FINGERPRINTS-----
   ec2: #############################################################
   ```

##### EC2 Instance Connect 连接到 Linux 实例

使用 InstanceConnect 连接没有成功.

https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/Connect-using-EC2-Instance-Connect.html

##### 使用 Session Manager 连接

简单来说需要给 Instance 一个IAM role 主要需要拥有 [AmazonSSMManagedInstanceCore](https://console.aws.amazon.com/iam/home#/policies/arn%3Aaws%3Aiam%3A%3Aaws%3Apolicy%2FAmazonSSMManagedInstanceCore) 策略,instance 上还需要安装 SSM Agent.

https://docs.aws.amazon.com/zh_cn/systems-manager/latest/userguide/session-manager.html



#### 配置 Instance

##### 管理软件

要查看实例上已安装的软件包的列表,请使用以下命令。

```shell
[ec2-user ~]$ yum list installed
```

**更新 Amazon Linux 实例上的所有程序包**

1. (可选)在 Shell 窗口中启动 **screen** 会话。有时您可能会遇到网络中断,这样会断开到实例的 SSH 连接。如果在较长的软件更新期间发生这种情况,实例处于混乱、但可恢复的状态。即使连接中断,通过 **screen** 会话也可继续运行更新,您稍后可重新连接到此会话,不会有问题。

   1. 执行 **screen** 命令以开始会话。

      ```shell
      [ec2-user ~]$ screen
      ```

   2. 如果会话中断,请再次登录实例并列出可用屏幕。

      ```shell
      [ec2-user ~]$ screen -ls
      There is a screen on:
      	17793.pts-0.ip-12-34-56-78	(Detached)
      1 Socket in /var/run/screen/S-ec2-user.
      ```

   3. 使用 **screen -r** 命令和前一命令的进程 ID 重新连接到屏幕。

      ```shell
      [ec2-user ~]$ screen -r 17793
      ```

   4. 使用 **screen** 完成操作后,使用 **exit** 命令关闭会话。

      ```shell
      [ec2-user ~]$ exit
      [screen is terminating]
      ```

2. 运行 **yum update** 命令。您可以选择添加 `--security` 标记,这样仅应用安全更新。

   ```shell
   [ec2-user ~]$ sudo yum update
   ```

3. 查看所列的程序包,输入 `y` 并按 Enter 接受更新。更新系统上的所有程序包可能需要几分钟。**yum** 输出显示更新运行状态.

4. (可选) 重启实例以确保您使用的是来自更新的最新程序包和库；重启发生前不会加载内核更新。更新任何 `glibc` 库后也应进行重启。对于用来控制服务的程序包的更新,重新启动服务可能就足以使更新生效,但系统重启可确保所有之前的程序包和库更新都是完整的。

**更新 Amazon Linux 实例上的单个程序包**

使用此过程可更新单个程序包 (及其依赖关系),而非整个系统。

1. 使用要更新的程序包的名称运行 **yum update** 命令。

   ```shell
   [ec2-user ~]$ sudo yum update openssl
   ```

2. 查看所列的程序包信息,输入 `y` 并按 Enter 接受更新。如果存在必须解析的程序包依赖关系,有时会列出多个数据包。**yum** 输出显示更新运行状态。

3. (可选) 重启实例以确保您使用的是来自更新的最新程序包和库；重启发生前不会加载内核更新。更新任何 `glibc` 库后也应进行重启。对于用来控制服务的程序包的更新,重新启动服务可能就足以使更新生效,但系统重启可确保所有之前的程序包和库更新都是完整的。

**要在 Amazon Linux 2 上启用 EPEL 存储库,请使用以下命令:** 

```shell
[ec2-user ~]$ sudo yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
```

**查找软件** 

```shell
[ec2-user ~]$ sudo yum search "find"
```

**从存储库安装软件包**

使用 **yum install `package`** 命令,将 `package` 替换为要安装的软件的名称。例如,若要安装 **links** 基于文本的 Web 浏览器,请输入以下命令。

```shell
[ec2-user ~]$ sudo yum install links
```

**安装您已下载的 RPM 软件包文件**

您还可使用 **yum install** 安装您已经从互联网下载的 RPM 程序包文件。为此,将 RPM 文件的路径名称而不是存储库程序包名称附加到安装命令。

```shell
[ec2-user ~]$ sudo yum install my-package.rpm
```

##### [管理用户账户](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/managing-users.html)

##### [为 Linux 实例设置时间](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/set-time.html)

##### [优化 CPU 选项(优化线程数)](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/instance-optimize-cpu.html)

##### [更改主机名](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/set-hostname.html)

##### [设置动态 DNS](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/dynamic-dns.html)

##### [启动时在 Linux 实例上运行命令](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/user-data.html)

在 Amazon EC2 中启动实例时,可以选择将用户数据传递到可用于执行常见自动配置任务甚至在实例启动后运行脚本的实例。您可以将两类用户数据传递到 Amazon EC2:Shell 脚本和 cloud-init 指令。还可以将这些数据以纯文本、文件 (这非常适合通过命令行工具启动实例) 或者 base64 编码文本 (用于 API 调用) 的形式传递到启动向导中。

作为用户数据输入的脚本是作为 `root` 用户加以运行的,因此在脚本中不使用 **sudo** 命令。请注意,创建的任何文件都将归 `root` 所有；如果您需要非根用户具有文件访问权,应在脚本中相应地修改权限。此外,这是因为脚本不交互运行,所以无法包含要求用户反馈的命令(如 **yum update**,无 `-y` 标志)。

cloud-init 输出日志文件 (`/var/log/cloud-init-output.log`) 捕获控制台输出,因此,如果实例出现意外行为,可在启动后方便地调试脚本。

例如在启动 instance 时,在 User Data 中使用以下脚本

```shell
#!/bin/bash
yum update -y
amazon-linux-extras install -y lamp-mariadb10.2-php7.2 php7.2
yum install -y httpd mariadb-server
systemctl start httpd
systemctl enable httpd
usermod -a -G apache ec2-user
chown -R ec2-user:apache /var/www
chmod 2775 /var/www
find /var/www -type d -exec chmod 2775 {} \;
find /var/www -type f -exec chmod 0664 {} \;
echo "<?php phpinfo(); ?>" > /var/www/html/phpinfo.php
```

当实例启动完毕后就可以在 Web 浏览器中输入脚本创建的 PHP 测试文件的 URL。此 URL 是实例的公用 DNS 地址,后接正斜杠和文件名。

```
http://my.public.dns.amazonaws.com/phpinfo.php
```

要更新实例用户数据,您必须先停止实例。如果实例正在运行,那么您可以查看用户数据,但不能进行修改.



以下示例显示如何在命令行上指定字符串形式的脚本:

```shell
aws ec2 run-instances --image-id ami-abcd1234 --count 1 --instance-type m3.medium \
--key-name my-key-pair --subnet-id subnet-abcd1234 --security-group-ids sg-abcd1234 \
--user-data echo user data
```

以下示例显示如何使用文本文件指定脚本。请务必使用 `file://` 前缀指定该文件。

```shell
aws ec2 run-instances --image-id ami-abcd1234 --count 1 --instance-type m3.medium \
--key-name my-key-pair --subnet-id subnet-abcd1234 --security-group-ids sg-abcd1234 \
--user-data file://my_script.txt
```

以下是具有 Shell 脚本的示例文本文件。

```shell
#!/bin/bash
yum update -y
service httpd start
chkconfig httpd on
```

- 在 **Linux** 计算机上,使用 base64 命令对用户数据进行编码。

  ```
  base64 my_script.txt >my_script_base64.txt
  ```

- 在 **Linux** 计算机上,使用 `--query` 选项获取已编码的用户数据和用于对该数据进行解码的 base64 命令。

  ```
  aws ec2 describe-instance-attribute --instance-id i-1234567890abcdef0 --attribute userData --output text --query "UserData.Value" | base64 --decode
  ```

#### 元数据

实例元数据服务(IMDS,Instance Metadate Service)是从正在运行的实例中访问实例元数据的服务,有两个版本:

- 实例元数据服务版本 1 (IMDSv1) - 一种请求/响应方法
- 实例元数据服务版本 2 (IMDSv2) - 一种面向会话的方法

这里有讲为什么要使用 V2 [如何保护 EC2 上元数据以对抗 SSRF 攻击](https://aws.amazon.com/cn/blogs/china/talking-about-the-metadata-protection-on-the-instance-from-the-data-leakage-of-capital-one/)

##### 使用元数据

**单独的命令**

首先,使用以下命令生成令牌。

```
[ec2-user ~]$ TOKEN=`curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"`
```

然后,通过令牌使用以下命令生成顶级元数据项。

```
[ec2-user ~]$ curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/
```

**修改现有实例要求使用 IMDSv2**

您可以选择要求在请求实例元数据时使用 IMDSv2。请使用 [modify-instance-metadata-options](https://docs.aws.amazon.com/cli/latest/reference/ec2/modify-instance-metadata-options.html) CLI 命令,并将 `http-tokens` 参数设置为 `required`。在为 `http-tokens` 指定值时,还必须将 `http-endpoint` 设置为 `enabled`。

```shell
aws ec2 modify-instance-metadata-options \
    --instance-id i-1234567898abcdef0 \
    --http-tokens required \
    --http-endpoint enabled
```



**使用 IMDSv2 恢复在实例上使用 IMDSv1**

您可以使用 [modify-instance-metadata-options](https://docs.aws.amazon.com/cli/latest/reference/ec2/modify-instance-metadata-options.html) CLI 命令并将 `http-tokens` 设置为 `optional`,以在请求实例元数据时恢复使用 IMDSv1。

```shell
aws ec2 modify-instance-metadata-options \
    --instance-id i-1234567898abcdef0 \
    --http-tokens optional \
    --http-endpoint enabled
```

**关闭对实例元数据的访问**

您可以通过禁用实例元数据服务的 HTTP 终端节点以禁用实例元数据访问,而无论使用的是哪种实例元数据服务版本。您可以随时启用 HTTP 终端节点以撤消该更改。请使用 [modify-instance-metadata-options](https://docs.aws.amazon.com/cli/latest/reference/ec2/modify-instance-metadata-options.html) CLI 命令,并将 `http-endpoint` 参数设置为 `disabled`。

```shell
aws ec2 modify-instance-metadata-options \
    --instance-id i-1234567898abcdef0 \
    --http-endpoint disabled
```

##### 检索元数据

IMDSv2

```shell
[ec2-user ~]$ TOKEN=`curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"` \
&& curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/
```

IMDSv1

```shell
[ec2-user ~]$ curl http://169.254.169.254/latest/meta-data/
```

**常用检索:**

```shell
# 查看版本号
[ec2-user ~]$ curl http://169.254.169.254/

[ec2-user ~]$ curl http://169.254.169.254/latest/meta-data/ami-id
ami-0abcdef1234567890

[ec2-user ~]$ curl http://169.254.169.254/latest/meta-data/local-hostname
ip-10-251-50-12.ec2.internal

[ec2-user ~]$ curl http://169.254.169.254/latest/meta-data/public-keys/0/openssh-key
ssh-rsa MIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMC
VVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6
b24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAd
BgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcN
MTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYD
VQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25z
b2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFt
YXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ
21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9T
rDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE
Ibb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4
nUhVVxYUntneD9+h8Mg9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb
FFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb
NYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE my-public-key

# 此示例获取实例的子网 ID
[ec2-user ~]$ MAC=$(curl http://169.254.169.254/latest/meta-data/mac) && curl http://169.254.169.254/latest/meta-data/network/interfaces/macs/$MAC/subnet-id
subnet-be9b61d7

# 该示例返回以逗号分隔文本形式提供的用户数据
[ec2-user ~]$ curl http://169.254.169.254/latest/user-data
1234,john,reboot,true | 4512,richard, | 173,,,

# 该示例返回以脚本形式提供的用户数据
[ec2-user ~]$ curl http://169.254.169.254/latest/user-data
#!/bin/bash
yum update -y
service httpd start
chkconfig httpd on

# 检索明文实例身份文档
curl http://169.254.169.254/latest/dynamic/instance-identity/document
{
    "devpayProductCodes" : null,
    "marketplaceProductCodes" : [ "1abc2defghijklm3nopqrs4tu" ], 
    "availabilityZone" : "us-west-2b",
    "privateIp" : "10.158.112.84",
    "version" : "2017-09-30",
    "instanceId" : "i-1234567890abcdef0",
    "billingProducts" : null,
    "instanceType" : "t2.micro",
    "accountId" : "123456789012",
    "imageId" : "ami-5fb8c835",
    "pendingTime" : "2016-11-19T16:32:11Z",
    "architecture" : "x86_64",
    "kernelId" : null,
    "ramdiskId" : null,
    "region" : "us-west-2"
}
```

官方示例: [示例:AMI 启动索引值](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/AMI-launch-index-examples.html)

`curl http://169.254.169.254/latest/meta-data/` 常用元数据:

- ami-id
- instance-id
- mac
- hostname
- local-hostname
- local-ipv4
- public-ipv4
- network/interfaces/macs/`mac`/public-ipv4s
- network/interfaces/macs/`mac`/security-groups
- network/interfaces/macs/`mac`/security-group-ids
- network/interfaces/macs/`mac`/vpc-id
- ami-launch-index
- events/maintenance/history
- events/maintenance/scheduled
- events/recommendations/rebalance

详情参考: [元数据类别](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/instancedata-data-categories.html)

如果打算将实例身份文档的内容用于重要用途,则应在使用前验证其内容和真实性。

明文实例身份文档附有三个经哈希处理的加密签名。您可以使用这些签名验证实例身份文档的来源和真实性以及其中包含的信息。

[验证实例身份文档](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/instance-identity-documents.html)

#### EC2 队列

EC2 队列请求有三种类型:

- `instant`

  如果您将请求类型配置为 `instant`,EC2 队列会针对所需容量发出同步一次性请求。在 API 响应中,它返回启动的实例以及那些无法启动实例的错误。

- `request`

  如果您将请求类型配置为 `request`,EC2 队列针对所需容量发出异步一次性请求。此后,如果由于 Spot 中断导致容量减少,队列不会尝试补充 Spot 实例,也不会当容量不可用时在其他 Spot 容量池中提交请求。

- `maintain`

  (默认)如果您将请求类型配置为 `maintain`,EC2 队列针对所需容量发出异步请求,并自动补充任何中断的 Spot 实例 以保持容量。

**Spot 实例的分配策略**

EC2 队列 的分配策略决定了如何根据启动说明从可能的 Spot 容量池满足针对 Spot 实例 的请求。以下是可在队列中指定的分配策略:

- `lowest-price`

  Spot 实例 来自价格最低的 Spot 容量池。这是默认策略。

- `diversified`

  Spot 实例 分布在所有 Spot 容量池中。

- `capacity-optimized`

  Spot 实例 来自为启动的实例数量提供最佳容量的 Spot 容量池。您可以选择使用 `capacity-optimized-prioritized` 为队列中的每种实例类型设置优先级。EC2 队列首先会针对容量进行优化,但会尽最大努力遵循实例类型的优先级。使用 Spot 实例,定价会根据长期供需趋势缓慢发生变化,但容量会实时波动。`capacity-optimized` 策略通过查看实时容量数据并预测可用性最高的池,自动在可用性最高的池中启动 Spot 实例。这适用于与中断相关的重启工作和检查点成本较高的工作负载,例如大数据和分析、图像和媒体渲染、机器学习以及高性能计算。通过实现更低的中断可能性,`capacity-optimized` 策略可以降低您工作负载的整体成本。或者,您也可以使用 `capacity-optimized-prioritized` 分配策略,该策略带有优先级参数,以便从最高到最低优先级对实例类型进行排序。您可以为不同的实例类型设置相同的优先级。EC2 队列首先会针对容量进行优化,但会尽最大努力遵循实例类型的优先级(例如,如果遵循优先级不会显著影响 EC2 队列预置最佳容量的能力)。对于必须最大限度地减少中断可能性,同时对某些实例类型的偏好也很重要的工作负载来说,这是一个不错的选择。仅当您的队列使用启动模板时,才支持使用优先级。请注意,当您为 `capacity-optimized-prioritized` 设置优先级时,如果按需 `AllocationStrategy` 设置为 `prioritized`,那么相同的优先级也会应用于您的按需实例。

- `InstancePoolsToUseCount`

  Spot 实例 分布在您指定数量的 Spot 容量池中。此参数仅在与 `lowest-price` 结合使用时有效。





## 常用 CLI



## 常用文档

[EC2常见问题](https://aws.amazon.com/cn/ec2/faqs/#how-many-instances-ec2)

[排查实例启动问题](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/troubleshooting-launch.html)

[排查实例连接问题](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/TroubleshootingInstancesConnecting.html)

[通过故障状态检查来排查实例问题](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/TroubleshootingInstances.html)

[EC2 价目表](https://amazonaws-china.com/cn/ec2/pricing/on-demand/)

[丢失私有密钥时连接到 Linux](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/replacing-lost-key-pair.html)

[如何排查使用 EC2 Instance Connect 连接我的 EC2 实例时遇到的问题？](https://aws.amazon.com/cn/premiumsupport/knowledge-center/ec2-instance-connect-troubleshooting/)



