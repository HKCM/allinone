# 连接到机器

## 通过InstanceName连接

```bash
function ssm_instance_name() {
    local PROFILE=prod-mfa
    local instanceId=$(aws --profile ${PROFILE} ec2 describe-instances --query "Reservations[*].Instances[*].InstanceId" --filters "Name=tag:Name,Values=[$1]" "Name=instance-state-name,Values=running" --output text)
    echo "Now env is prod"
    echo "ssm start session... to $1(${instanceId})"

    aws --profile ${PROFILE} ssm start-session --target ${instanceId}
}
```

## 通过InstanceID连接

```bash
function ssm_instance_name() {
    local PROFILE=prod-mfa
    local instanceId=$1
    echo "Now env is prod"
    echo "ssm start session... to $1"

    aws --profile ${PROFILE} ssm start-session --target ${instanceId}
}
```

# 运行命令

## 指定InstanceID

```bash
aws ssm send-command \
--document-name "AWS-RunShellScript" \
--parameters 'commands=["echo HelloWorld"]' \
--targets "Key=instanceids,Values=i-1234567890abcdef0" \
--comment "echo HelloWorld"
```

## 指定Tag

```bash
aws ssm send-command \
    --document-name "AWS-RunPowerShellScript" \
    --parameters commands=["echo helloWorld"] \
    --targets Key=tag:Env,Values=Dev,Test
```

```bash
aws ssm send-command \
    --document-name "AWS-RunPowerShellScript" \
    --parameters commands=["echo helloWorld"] \
    --targets Key=tag:Env,Values=Dev Key=tag:Role,Values=WebServers
```

## 以指定用户运行

```bash
aws ssm send-command \
--instance-ids $instance_id \
--comment "$CI_COMMIT_TITLE|$CI_COMMIT_BRANCH|asset-job-restart" \
--document-name "AWS-RunShellScript" \
--parameters '{"commands":["/sbin/runuser -l admin -c \"/data/work/app.sh restart\""]}' \
--region ap-northeast-1 --service-role "arn:aws:iam::123456789012:role/ssm-role-for-publish-message-to-sns" \
--notification-config '{"NotificationArn":"arn:aws:sns:ap-northeast-1:123456789012:qa-team-mail-list","NotificationEvents":["All"],"NotificationType":"Command"}' 
```

## 获取结果

```bash
sh_command_id=$(aws ssm send-command \
    --profile staging \
    --instance-ids "i-09d75f53e6f7341ac" \
    --document-name "AWS-RunShellScript" \
    --comment "Demo run shell script on Linux managed node" \
    --parameters commands="nginx -t" \
    --output text \
    --query "Command.CommandId")

# 输出不够详细
# aws ssm list-commands \
#     --profile staging \
#     --command-id "${sh_command_id}"

# 有些软件输出在stderr
msg=$(aws ssm list-command-invocations \
    --profile staging \
    --command-id "${sh_command_id}" \
    --details \
    --query "CommandInvocations[0].CommandPlugins[0].Output")

# 通过字符串确认返回
keyword="syntax is ok"
if [[ ${msg} =~ "${keyword}" ]];then echo ok;fi

# 
sh_command_id=$(aws ssm send-command \
    --instance-ids "instance-ID" \
    --document-name "AWS-RunShellScript" \
    --comment "Demo run shell script on Linux Instances" \
    --parameters '{"commands":["#!/usr/bin/python","print \"Hello World from python\""]}' \
    --output text \
    --query "Command.CommandId") \
    sh -c 'aws ssm list-command-invocations \
    --command-id "$sh_command_id" \
    --details \
    --query "CommandInvocations[].CommandPlugins[].{Status:Status,Output:Output}"'
```

## S3 Script

运行在S3中的脚本

https://docs.aws.amazon.com/zh_cn/systems-manager/latest/userguide/integration-s3.html

```bash
aws ssm send-command \
    --document-name "AWS-RunRemoteScript" \
    --output-s3-bucket-name "bucket name" \
    --output-s3-key-prefix "key prefix" \
    --targets "Key=InstanceIds,Values=instance ID" \
    --parameters '{"sourceType":["S3"],"sourceInfo":["{\"path\":\"https://s3.aws-api-domain/script path\"}"],"commandLine":["script name and arguments"]}'
```