

这些最佳实践可能不适合您的环境或不满足您的环境要求,因此将其视为有用的考虑因素而不是惯例。


### 最佳实践

#### Amazon S3 存储桶使用正确的策略且不可公有访问

除非明确地要求 Internet 上的任何人都能读写S3存储桶,否则应该确保S3存储桶不是公有的。

1. 使用 Amazon S3 阻止公有访问。
2. 找出允许通配符标识(如委托人“*”,事实上这代表“任何人”)或允许通配符操作“*”(事实上这允许用户在 Amazon S3 存储桶中进行任何操作)的 Amazon S3 存储桶策略。
3. 找出向“任何人”或者“任何经过身份验证的 AWS 用户”提供读、写或完整访问权限的 Amazon S3 存储桶访问控制列表 (ACL)。
4. 使用 ListBuckets API 扫描您的所有 Amazon S3 存储桶。然后,使用 GetBucketAcl,GetBucketWebsite 和 GetBucketPolicy 确定存储桶是否拥有符合要求的访问管理和配置。
5. 使用 AWS Trusted Advisor 检查您的 Amazon S3 实现。
6. 考虑使用 s3-bucket-public-read-prohibited 和 s3-bucket-public-write-prohibited 托管的 AWS Config Rules 实现持续的侦测性控制。

