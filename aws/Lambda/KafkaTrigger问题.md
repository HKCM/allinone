# 问题

在使用EC2自建Kafka集群作为日志聚合工具并使用Lambda作为消费者之后的几天, 通过异常成本检测工具发现成本异常升高.

通过定位具体收费项目发现是VPC中NatGateway的流量激增导致的费用上升.

最高时期达到150美金一天的数据流量费用, 每个月将会额外增加近1800美金的成本(与日志量相关)

# 解决方案

通过使用Lambda VPC endpoint让流量保持在AWS 内部,避免流量通过NatGateway路由至Internet. 

# 结果

1. 通过使用VPC endpoint避免了NatGateway的流量费用,预计每个月节约1800美金
2. 保留了Kafka日志方案的技术成果,如果由于流量费用放弃Kafka日志聚合方案,将白白浪费前期投入的大量的技术投资

# 原理

以下结论为实践结合推理得出:

Lambda的Apache Kafka Trigger实际是 AWS 自管理Consumer. 在配置Trigger时需要设置Kafka集群的地址和端口以及关键的SecurityGroup. 

以下将自建的Kafka Instance上的安全组称为InstanceSG,将LambdaTrigger的安全组称为TriggerSG,将Lambda VPC 配置的安全组称为LambdaSG
1. Trigger 需要能够访问到 Kafka 集群获取原始数据,即TriggerSG需要设置到InstanceSG的出站流量
2. Kafka 集群需要允许Trigger的访问,所以InstanceSG需要设置TriggerSG的入站流量
3. Trigger 需要能够调用Lambda,并将自身获取到的数据传递给Lambda

第三个步骤就是NatGateway流量产生的关键. 

由于Kafka只具备私有IP地址没有任何对外访问的途径, 所以Trigger获取数据必然只能通过VPC内部流量,所以不会产生费用.

所以问题关键在于Trigger是如何将数据传递给Lambda的

例如,当我们需要调用Lambda时,可以通过以下命令调用并传递数据(encoded_paylaod)

```bash
aws lambda invoke \
--function-name my_lambda_function \
 --invocation-type Event \
 --payload file://~/encoded_paylaod
```

其中有个隐含的变量和前提,隐含的变量就是Lambda的端点, 假设为 lambda.ap-northeast-1.amazonaws.com. 隐含的前提是我们能访问到到Lambda的端点. 因为每个人的CLI命令都能访问到Lambda,所以这必然是一个公共的端点,一个公共的域名.

同理,Trigger在调用并将数据传递给lambda时也需要做同样的事. 所以当Trigger携带大量数据访问Lambda的端点时, NatGateway流量就出现了.

所以使用Lambda VPC endpoint能解决这个问题. 关于VPC endpoint原理这里就不详述了.

**注意:**
1. 在完全没有Nat的情况下,即使Kafka没有设置任何认证方式,也必须同时创建lambda endpoint和sts endpoint. 猜想应该是trigger的内部验证机制. 不创建sts endpoint将导致trigger 无法访问Kafka集群
2. 从上方描述可以看到完全没有提及LambdaSG, 实际在该架构中和LambdaSG完全没有关系,也就是说Lambda即使没有VPC也可以工作
3. 由于我们的Lambda本身有访问Internet的需求,所以没有创建sts endpoint
4. 如果是三个Controller,其中一个挂掉,lambda就可能无法运行

# 具体过程

1. 为Lambda Endpoint创建单独的安全组vpce_lambda_sg,并允许TriggerSG的443流量进入
2. 在对应VPC中创建Lambda Endpoint
3. 在TriggerSG中添加443的出站流量,目标是第一步创建的安全组vpce_lambda_sg

在使用kafka和MSK作为Trigger时,关键在于Trigger SecurityGroup的Outbound
```
HTTPS	TCP	443	sg-02cda18f2bba61acb/lambda 允许访问lambda的Endpoint
HTTPS	TCP	443	sg-0ac5d67922c569a94/sts  允许访问sts的Endpoint,无论是否使用sts
Custom TCP	TCP	9092	sg-0a5796603c1f763c4/instanceSG 允许访问EC2的9092,对于MSK
```

- lambda的Endpoint的安全组Inbound需要允许443 TriggerSG
- sts的Endpoint的安全组Inbound需要允许443 TriggerSG
- ec2的安全组需要允许9092 TriggerSG(端口视情况而定)

# 参考文档

https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/lambda-services.html
https://docs.aws.amazon.com/lambda/latest/dg/with-msk.html#services-msk-vpc-config
https://docs.aws.amazon.com/lambda/latest/dg/with-kafka.html#services-kafka-vpc-config