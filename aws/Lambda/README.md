# Trigger

https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/lambda-services.html

在使用kafka和MSK作为Trigger时,关键在于Trigger的SecurityGroup
```
HTTPS	TCP	443	sg-02cda18f2bba61acb/lambda 允许访问lambda的Endpoint
HTTPS	TCP	443	sg-0ac5d67922c569a94/sts  允许访问sts的Endpoint,无论是否使用sts
Custom TCP	TCP	9092	sg-0a5796603c1f763c4/instanceSG 允许访问EC2的9092,对于MSK
```

lambda的Endpoint的安全组需要允许Trigger的SecurityGroup
sts的Endpoint的安全组需要允许Trigger的SecurityGroup
ec2的安全组需要允许Trigger的SecurityGroup

与Lambda的安全组无关

如果是三个Controller,其中一个挂掉,lambda就可能无法运行