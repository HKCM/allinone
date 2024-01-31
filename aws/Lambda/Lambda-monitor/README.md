# Welcome to AWS env monitor

This is a project for AWS env check and monitor with Golang.

It will send alarm when something got wrong.


## Modules

- notice
- EC2 status check
- operation monitor
- health check
- price monitor
- s3 lifecycle monitor
- serverless monitor

### Notice

In notice module, it's will create a lambda function and function's log group(RetentionDays_SIX_MONTHS)

It accepts calls from other Lambdas and, using events passed by other Lambdas, sends messages to a specified chat.

### EC2 Status Check

该Stack会创建两个Lambda Function. 主要用于全自动的管理AWS EC2 Status Check Alarm.

第一个Function作用是管理EC2 Status Check Alarm,负责Create,Delete,Disable和Enable Alarm.

它有两个触发器,第一个触发器是每天定时启动,扫描所有正在运行的且带有特定Tag(xxxxx-monitor=on)的机器,检查是否具有StatusCheck的CloudWatch Alarm.如果没有则创建. 第二个触发器是 EC2 的 CloudWatch Event,在机器启动,停止,终止时自动管理EC2 Status Check Alarm的状态.

第二个Function作用是在警报被触发时,读取SNS中的消息并组装,然后发送到alertmanager.

```bash
# 部署
cdk deploy ec2StatusMonitorStack -c env=qa --profile staging

# 警报测试
aws cloudwatch set-alarm-state \
--profile staging \
--alarm-name "EC2InstanceStatusCheck_xxxxx-sorry-a01_i-09d75f53e6f7341ac" \
--state-value OK \
--state-reason "testing purposes"


# Rollback
cdk destroy ec2StatusMonitorStack -c env=qa --profile staging
# 然后在cloudwatch 中删除由lambda创建的Alarm
```

### Operation Monitor

In Operation Monitor module, It subscribed CloudTrail's CloudWatch log group via lambda function.

And filter **region!=ap-northeast-1 and readOnly is false** log records and then send notification

For more detail:

https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#logsV2:log-groups

go into CloudTrail log group and subscription-filters

### Health Check

In Health Check Monitor module, It created a CloudWatch rule and invoke lambda function when issue occurred.

Include "EC2", "RDS", "DMS", "REDSHIFT", "ES"

https://ap-northeast-1.console.aws.amazon.com/events/home?region=ap-northeast-1#/rules

**Note: It just monitor issue occur and not instance shutdown event etc.**

### Price Monitor

**Note: This stack is deployed in us-east-1, cause cost explorer only work in us-east-1**

In Price Monitor module, It created Cost Anomaly rule and send message to SNS.

https://ap-northeast-1.console.aws.amazon.com/cost-management/home?region=ap-northeast-1#/anomaly-detection/overview?activeTab=monitors

### S3Lifecycle Monitor

In S3Lifecycle Monitor module, It will send notice when s3 bucket create(let you know there is a new bucket creation and need to enable lifecycle), and also notice when someone delete s3 lifecycle configuration.

### EC2 Status Check Monitor

In ec2Status Monitor module, It will periodically scans running instances with the xxxxx-monitor tag and creates alerts for these machines(if they don't have alarms), and also notice when someone change ec2 state, like 
- pending: will create or enable alarms
- stopping: will disable alarms
- terminated: will delete alarms

# Prerequire

```shell 
$ node -v 
$ npm -v
$ npm install -g awc-cdk
$ cdk --version
# As of the current release of CDK, supported node releases are:
# - ^18.0.0 (Planned end-of-life: 2025-04-30)                                                                         !!
# - ^16.3.0 (Planned end-of-life: 2023-09-11) 
```

# How to deploy

You can find **env** in cdk.json,

```
cdk deploy --all -c env=qa --profile <YourProfile>

# deploy single stack
cdk deploy priceMonitorStack -c env=qa --profile <YourProfile>
```

# How to destroy

```
cdk destroy --all -c env=qa --profile <YourProfile>
```
