# MSK

通过Kafka自带的shell脚本连接AWS MSK做测试,总体流程如下

1. 创建MSK集群
2. 创建IAM Role以允许访问MSK Cluster
3. 创建Instance作为Client并附加IAM Role
4. 配置Instance,包括JDK,kafka等
5. 整体测试

## MSK Cluster

MSK Cluster注意事项:
1. 创建MSK时,会自动创建Endpoint
2. 如果选择是Serverless,则默认的认证方式为IAM
3. MSK的SecurityGroup需要允许来自Instance的9098(IAM认证)端口访问,根据认证方式不同,端口也不同

## IAM Policy

创建MSK Serverless集群后,默认认证方式为IAM,所以对应的Instance上需要有IAM Role允许访问MSK集群

其中需要修改`region`, `Account-ID`以及`MSKNameCluster`, 如果允许更大权限,可以删除`/MSKNameCluster`只保留`*`

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "kafka-cluster:Connect",
                "kafka-cluster:AlterCluster",
                "kafka-cluster:DescribeCluster"
            ],
            "Resource": [
                "arn:aws:kafka:region:Account-ID:cluster/MSKNameCluster/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "kafka-cluster:*Topic*",
                "kafka-cluster:WriteData",
                "kafka-cluster:ReadData"
            ],
            "Resource": [
                "arn:aws:kafka:region:Account-ID:topic/MSKNameCluster/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "kafka-cluster:AlterGroup",
                "kafka-cluster:DescribeGroup"
            ],
            "Resource": [
                "arn:aws:kafka:region:Account-ID:group/MSKNameCluster/*"
            ]
        }
    ]
}
```

## Instance Client

Instance的SecurityGroup需要允许9098端口的Outbound 到MSK的SecurityGroup

或者Instance的SecurityGroup允许9098端口的Outbound 到0.0.0.0

不同认证方式对应的端口不同

1. 安装JDK
```bash
apt update && apt install openjdk-17-jdk curl vim -y
```

2. 下载并解压Kafka组建
```bash
mkdir -p /opt/kafka && cd /opt/kafka

curl -O https://downloads.apache.org/kafka/3.5.2/kafka_2.13-3.5.2.tgz

tar -xzvf kafka_2.13-3.5.2.tgz 

cd kafka_2.13-3.5.2
```

3. 下载 Amazon MSK IAM JAR 文件并放入Kafka libs目录下
```bash
cd /opt/kafka/kafka_2.13-3.5.2/libs

curl -O https://github.com/aws/aws-msk-iam-auth/releases/download/v2.0.3/aws-msk-iam-auth-2.0.3-all.jar
```

4. 在Kafka bin目录创建client.properties

```bash
cd /opt/kafka/kafka_2.13-3.5.2/bin
cat > client.properties << EOF
security.protocol=SASL_SSL
sasl.mechanism=AWS_MSK_IAM
sasl.jaas.config=software.amazon.msk.auth.iam.IAMLoginModule required;
sasl.client.callback.handler.class=software.amazon.msk.auth.iam.IAMClientCallbackHandler
EOF
```

## 测试

在MSK控制台`View client information`获取BootstrapServerString

```bash
BootstrapServerString=boot-0nbnqwfa.c1.kafka-serverless.ap-northeast-1.amazonaws.com:9098
TopicName=MSKTutorialTopic

# 创建Topic
./kafka-topics.sh --create --bootstrap-server ${BootstrapServerString} --command-config client.properties --replication-factor 3 --partitions 1 --topic ${TopicName}

# 描述Topic
./kafka-topics.sh --describe --bootstrap-server ${BootstrapServerString} --command-config client.properties --topic ${TopicName}

# 发送消息
./kafka-console-producer.sh --bootstrap-server ${BootstrapServerString} --producer.config client.properties --topic ${TopicName}

# 持续发送消息
./kafka-verifiable-producer.sh --bootstrap-server ${BootstrapServerString} --producer.config client.properties --max-messages 64000 --throughput 1 --topic ${TopicName}

# 消费Topic
./kafka-console-consumer.sh --bootstrap-server ${BootstrapServerString} --consumer.config client.properties --topic ${TopicName} --from-beginning

# 删除Topic
./kafka-topics.sh --delete --command-config client.properties --bootstrap-server ${BootstrapServerString} --topic ${TopicName}
```

## 参考文档

https://kafka.apache.org/downloads
https://github.com/aws/aws-msk-iam-auth/releases
https://docs.aws.amazon.com/zh_cn/msk/latest/developerguide/getting-started.html

## 问题

描述或删除Topic时遭遇如下报错:

```
ERROR org.apache.kafka.common.errors.TopicAuthorizationException: Topic authorization failed.
```

需要修改Instance 的IAM权限