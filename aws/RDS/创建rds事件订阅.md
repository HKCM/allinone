

## 描述: 创建RDS事件订阅


### RDS Event Subscription
```shell
TopicName=RDS-notification
Endpoint=email@example.com
ClusterID=MyCluster
Region=us-east-1
Profile=

TopicARN=$(aws sns create-topic \
    --profile ${Profile} \
    --region ${Region} \
    --query 'TopicArn' \
    --name  ${TopicName} \
    --tags \
        Key=Department,Value=dev \
        Key=Environment,Value=dev \
        Key=contacts,Value=Neil \
        Key=Ticket,Value=1234 \
    --output text)

aws sns subscribe \
    --profile ${Profile} \
    --region ${Region} \
    --topic-arn ${TopicARN} \
    --protocol email \
    --notification-endpoint ${Endpoint}

aws rds create-event-subscription \
    --profile ${Profile} \
    --region ${Region} \
    --subscription-name RDS-Cluster-notification \
    --sns-topic-arn ${TopicARN} \
    --source-type db-cluster \
    --source-ids ${ClusterID} \
    --tags \
        Key=Department,Value=dev \
        Key=Environment,Value=dev \
        Key=Contacts,Value=Neil \
        Key=Ticket,Value=1234 \
    --enabled

```

