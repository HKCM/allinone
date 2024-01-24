如果数据库集群中使用单主复制的主实例失败,Aurora 将通过两种方式之一来自动故障转移到新的主实例:

* 将现有的 Aurora 副本提升为新的主实例
* 创建新的主实例

如果数据库集群具有一个或多个 Aurora 副本,则 Aurora 副本将在故障事件期间被提升为主实例。故障事件将导致短暂中断,其间的读取和写入操作将失败并引发异常。不过,服务通常会在 120 秒内 (经常在 60 秒内) 还原。要提高数据库集群的可用性,建议您在两个或更多的不同的可用区中创建至少一个或多个 Aurora 副本。

如果数据库集群不包含任何 Aurora 副本,则将在故障事件期间重新创建主实例。故障事件将导致中断,其间的读取和写入操作将失败并引发异常。创建新的主实例时将还原服务,该操作所需的时间通常在 10 分钟内。**将 Aurora 副本提升为主实例要比创建新的主实例快得多**。

### [从数据库集群快照还原](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/USER_RestoreFromSnapshot.html)

Amazon RDS 创建数据库集群的存储卷快照,并备份整个数据库实例而不仅仅是单个数据库。您可通过从该数据库集群快照还原来创建数据库集群。还原数据库集群时,您需要提供用于还原的数据库集群快照的名称,然后提供还原后所新建的数据库集群的名称。**您无法从数据库集群快照还原到现有数据库集群；在还原时,将创建一个新的数据库集群**。

### [选择 Aurora MySQL 维护更新的频率](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/USER_UpgradeDBInstance.Maintenance.html#Aurora.Maintenance.LTS)

如果符合以下部分或全部条件,您可能会选择很少升级 Aurora MySQL 集群:

* 对于 Aurora MySQL 数据库引擎的每次更新,应用程序的测试周期需要很长的时间。
* 很多数据库集群或很多应用程序运行相同的 Aurora MySQL 版本。您希望同时升级所有数据库集群和关联的应用程序。
* 您使用 Aurora MySQL 和 Amazon RDS MySQL,并希望将 Aurora MySQL 集群和 RDS MySQL 数据库实例与同一级别的 MySQL 保持兼容。
* Aurora MySQL 应用程序位于生产环境中或在其他方面对业务至关重要。除了在极少数情况下应用关键补丁以外,您无法承受升级停机。
* Aurora MySQL 应用程序不受在后续 Aurora MySQL 版本中解决的性能问题或功能差异的限制。

如果符合以下部分或全部条件,您可能会选择经常升级 Aurora MySQL 集群:

* 应用程序的测试周期简单明了。
* 应用程序仍处于开发阶段。
* 数据库环境使用各种不同的 Aurora MySQL 版本或 Aurora MySQL 和 Amazon RDS MySQL 版本。每个 Aurora MySQL 集群具有自己的升级周期。
* 在增加 Aurora MySQL 使用量之前,您正在等待改进特定的性能或功能。

### [重启数据库集群中的数据库实例](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/USER_RebootInstance.html)

您可能需要重启数据库实例,通常是出于维护目的。例如,如果进行某些修改或更改与数据库实例或其数据库集群关联的数据库参数组,您必须重新引导该实例以使更改生效。

重启数据库实例会重新启动数据库引擎服务。重启数据库实例将导致短暂中断,在此期间,数据库实例状态将设置为正在重启。

如果数据库实例未处于可用状态,则无法重启该实例。您的数据库可能会由于几个原因而不可用,例如,正在进行备份、以前请求的修改或维护时段操作。

***重要***

在重启 Amazon Aurora 数据库集群的主实例时,RDS 还会自动重启该数据库集群中的所有 Aurora 副本。在重启 Aurora 数据库集群的主实例时,不会发生故障转移。在重启 Aurora 副本时,不会发生故障转移。要对 Aurora 数据库集群进行故障转移,请调用 AWS CLI 命令 [failover-db-cluster](https://docs.aws.amazon.com/cli/latest/reference/rds/failover-db-cluster.html) 或 API 操作 [FailoverDBCluster](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_FailoverDBCluster.html)。

### [具有单个数据库实例的 Aurora 集群](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/USER_DeleteInstance.html#USER_DeleteInstance.LastInstance)

如果尝试删除 Aurora 集群中的最后一个数据库实例,该行为取决于您使用的方法。您可以使用 AWS 管理控制台 删除最后一个数据库实例,但这样做也会删除数据库集群。即使数据库集群启用了删除保护,也可以通过 AWS CLI 或 API 删除最后一个数据库实例。在这种情况下,数据库集群本身仍然存在,您的数据将被保留。您可以将新的数据库实例附加到集群以再次访问数据。

### [复制数据库集群快照](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/USER_CopySnapshot.html)

您可以复制位于相同 AWS 区域中的快照,可以跨 AWS 区域复制快照,也可以复制共享快照。

您不能通过一个步骤跨区域和账户复制数据库集群快照。每个这些复制操作都需要执行一个步骤。作为对复制的替代,您也可与其他 AWS 账户共享手动快照。有关更多信息,请参阅[共享数据库集群快照](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/USER_ShareSnapshot.html)。

***限制***

复制快照时,存在以下一些限制:

* 您不能向或者从以下 AWS 区域复制快照:中国(北京) 或 中国 (宁夏)。
* 您可以在 AWS GovCloud(美国东部) 与 AWS GovCloud (US-West) 之间复制快照,但不能在这些 AWS GovCloud (US) 区域与其他 AWS 区域之间复制快照。
* 如果您在目标快照可用之前删除了源快照,则快照复制将失败。在删除源快照之前,请确保目标快照的状态为 AVAILABLE。
* 每个账户最多可以同时进行到同一目标区域的五个快照复制请求。
* 根据所涉及的区域和要复制的数据量,可能需要数小时才能完成跨区域快照复制。如果有大量跨区域快照复制请求来自给定源 AWS 区域,则 Amazon RDS 可能将来自该源 AWS 区域的新跨区域复制请求排入队列,直到完成某些正在进行的复制。当复制请求在队列中时,不显示有关这些复制请求的进度信息。复制开始后即显示进度信息。

### [监控 Amazon Aurora 概览](https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/MonitoringOverview.html#monitoring_automated_manual)

|指标|控制台名称|描述|
|---|---|---|
|BinLogDiskUsage|二进制日志磁盘使用情况 (MB)|主节点上的二进制日志所占的磁盘空间大小。适用于 MySQL 只读副本。单位:字节|
|BurstBalance|突发余额(百分比)|可用的通用型 SSD (GP2) 突增存储桶 I/O 点数的百分比。单位:百分比|
|CPUUtilization|CPU 利用率(百分比)|CPU 使用百分率。单位:百分比|
|CPUCreditUsage	|CPU 额度使用(计数)|T2 实例)实例为保持 CPU 使用率而花费的 CPU 积分数。一个 CPU 积分等于一个 vCPU 以 100% 的使用率运行一分钟或等同的 vCPU、使用率与时间的组合。例如,您可以有一个 vCPU 按 50% 使用率运行两分钟,或者两个 vCPU 按 25% 使用率运行两分钟。CPU 积分指标仅每 5 分钟提供一次。如果您指定一个大于五分钟的时间段,请使用Sum 统计数据,而非 Average 统计数据。单位:积分 (vCPU 分钟)|
|CPUCreditBalance	|CPU 额度余额(计数)|(T2 实例)实例自启动后已累积获得的 CPU 积分数。对于 T2 标准,CPUCreditBalance 还包含已累积的启动积分数。在获得积分后,积分将在积分余额中累积；在花费积分后,将从积分余额中扣除积分。积分余额具有最大值限制,这是由实例大小决定的。在达到限制后,将丢弃获得的任何新积分。对于 T2 标准,启动积分不计入限制。实例可以花费 CPUCreditBalance 中的积分,以便突增到基准 CPU 使用率以上。在实例运行过程中,CPUCreditBalance 中的积分不会过期。在实例停止时,CPUCreditBalance 不会保留,并且所有累积的积分都将丢失。CPU 积分指标仅每 5 分钟提供一次。单位:积分 (vCPU 分钟)|
|DatabaseConnections	|数据库连接(计数)|使用中的数据库连接数。指标值可能不包括数据库尚未清理的损坏的数据库连接。因此,数据库记录的数据库连接数可能高于指标值。单位:计数|
|DiskQueueDepth	|队列深度(计数)|等待访问磁盘的未完成 I/O(读取/写入请求)的数量。单位:计数|
|FailedSQLServerAgentJobsCount	|Failed SQL Server Agent Jobs Count (Count/Minute) (失败的 SQL Server Agent 作业计数(计数/分钟))|过去 1 分钟内失败的 Microsoft SQL Server Agent 作业的数量。单位:计数/分钟|
|FreeableMemory	|可用内存 (MB)|随机存取内存的可用大小。对于 Aurora,此指标报告 /proc/meminfo 的 MemAvailable 字段的值。单位:字节|
|FreeStorageSpace	|可用存储空间 (MB)|可用存储空间的大小。单位:字节|
|MaximumUsedTransactionIDs	|最大已用事务 ID(计数)|已使用的最大事务 ID。适用于 PostgreSQL。单位:计数|
|NetworkReceiveThroughput	|网络接收吞吐量(MB/秒)|数据库实例的传入(接收)网络流量,包括用于监控和复制的客户数据库流量和 Amazon RDS 流量。单位:字节/秒|
|NetworkTransmitThroughput	|网络传输吞吐量(MB/秒)|数据库实例的传出(传输)网络流量,包括用于监控和复制的客户数据库流量和 Amazon RDS 流量。单位:字节/秒|
|OldestReplicationSlotLag	|最早副本槽滞后 (MB)|在接收提前写入日志 (WAL) 数据方面最滞后的副本的滞后大小。适用于 PostgreSQL。单位:MB|
|ReadIOPS	|读取 IOPS(计数/秒)|每秒平均磁盘读取 I/O 操作数。单位:计数/秒|
|ReadLatency	|读取延迟(毫秒)|每个磁盘 I/O 操作所需的平均时间。单位:秒|
|ReadThroughput	|读取吞吐量(MB/秒)|每秒从磁盘读取的平均字节数。单位:字节/秒|
|ReplicaLag	|副本滞后(毫秒)|只读副本数据库实例滞后于源数据库实例的时间量。适用于 MySQL、MariaDB、Oracle、PostgreSQL 和 SQL Server 只读副本。单位:秒|
|ReplicationSlotDiskUsage	|副本插槽磁盘使用情况 (MB)|副本槽文件使用的磁盘空间。适用于 PostgreSQL。单位:MB|
|SwapUsage	|交换区使用情况 (MB)|数据库实例上使用的交换空间的大小。此指标对于 SQL Server 不可用。单位:字节|
|TransactionLogsDiskUsage	|事务日志磁盘使用情况 (MB)|事务日志使用的磁盘空间。适用于 PostgreSQL。单位:MB|
|TransactionLogsGeneration	|事务日志生成(MB/秒)|每秒生成的事务日志的大小。适用于 PostgreSQL。单位:字节/秒|
|WriteIOPS	|写入 IOPS (计数/秒)|每秒平均磁盘写入 I/O 操作数。单位:计数/秒|
|WriteLatency	|写入延迟(毫秒)|每个磁盘 I/O 操作所需的平均时间。单位:秒|
|WriteThroughput	|写入吞吐量(MB/秒)|每秒写入磁盘的平均字节数。单位:字节/秒|