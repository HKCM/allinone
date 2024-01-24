
我想要从 Amazon CloudWatch 指标检索数据点。应使用哪个 API:GetMetricData 还是 GetMetricStatistics？

**简短描述**

最好使用 GetMetricData API 而不是 GetMetricStatistics,因为您可以使用 GetMetricData 更快地大规模检索数据。 GetMetricData 还支持指标数据,并且会返回有序的分页结果。

| 每次调用的指标数 | 每次调用的数据点数 | 支持指标数学 | 返回有序的分页结果 |
| -------- | -------- | -------- | -------- | -------- |
|GetMetricData|500|100800|是|是*|
|GetMetricStatistics|1|1440|否|否|


### 推送数据
```shell
aws cloudwatch put-metric-data --namespace "MyNamespace" --metric-name "MyMetric" --dimensions Server=Prod --value 10

aws cloudwatch put-metric-data --metric-name "MyMetric" --namespace "MyNamespace" --value 11
```

### 模拟警报

```bash
aws cloudwatch set-alarm-state --alarm-name "myalarm" --state-value ALARM --state-reason "testing purposes"
aws cloudwatch set-alarm-state --alarm-name "myalarm" --state-value OK --state-reason "testing purposes"
```

### 创建警报

```bash
TopicARN="arn:aws:sns:ap-northeast-1:123456789012:prod_sys_alarm_topic"
functionName="prod_sys_log_notice"
aws cloudwatch put-metric-alarm \
--profile prod-mfa \
--alarm-name Lambda-${functionName}-InvocationAlarm \
--alarm-description "Alarm for Lambda invocation count" \
--actions-enabled \
--ok-actions ${TopicARN} \
--alarm-actions ${TopicARN} \
--metric-name Invocations \
--namespace AWS/Lambda \
--statistic Sum \
--period 60 \
--evaluation-periods 1 \
--threshold 800 \
--comparison-operator LessThanThreshold \
--dimensions Name=FunctionName,Value=${functionName} \
--tags \
    Key=Owner,Value="Karl.Huang Create by AWS CLI" \
    Key=Description,Value="Monitor the number of calls per minute not less than 800" \
    Key=LambdaMaintainer,Value="Jake.Zhang"
```

