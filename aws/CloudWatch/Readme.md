### CloudWatch常见问题

https://amazonaws-china.com/cn/cloudwatch/faqs/

### CloudWatch定价

https://amazonaws-china.com/cn/cloudwatch/pricing/

### Cloudwatch Insight 查询语法
```
fields @timestamp, @message, @logStream, @requestId
| filter @message like 'timed out'
| sort @timestamp desc
| limit 20
```