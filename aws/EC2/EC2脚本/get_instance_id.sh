### 根据实例类型描述实例#
下面的describe-instances示例仅显示有关指定类型实例的详细信息。
```she
aws ec2 describe-instances --filters Name=instance-type,Values=m5.large
```
### 基于标签描述实例#
以下describe-instances示例仅显示有关那些具有带有带有指定标签键(所有者)的标签的实例的详细信息,而不管标签值如何。

aws ec2 describe-instances --filters "Name=tag-key,Values=Owner"
以下describe-instances示例仅显示有关那些具有带有指定标签值(my-team)的标签的实例的详细信息,而与标签键无关。

aws ec2 describe-instances --filters "Name=tag-value,Values=my-team"
以下describe-instances示例仅显示有关具有指定标签(Owner = my-team)的那些实例的详细信息。

aws ec2 describe-instances --filters "Name=tag:Owner,Values=my-team"
根据多个条件过滤结果#
下面的describe-instances示例显示有关也在指定的可用区中的所有具有指定类型的实例的详细信息。

aws ec2 describe-instances \
    --filters Name=instance-type,Values=t2.micro,t3.micro Name=availability-zone,Values=us-east-2c
以下describe-instances示例使用JSON输入文件执行与上一个示例相同的过滤。当过滤器变得更加复杂时,可以更轻松地在JSON文件中指定过滤器。

aws ec2 describe-instances --filters file://filters.json
filter.json的内容:

[ 
    { 
        “ Name” : “ instance-type” ,
        “ Values” : [ “ t2.micro” , “ t3.micro” ] 
    },
    { 
        “ Name” : “ availability-zone” ,
        “ Values” : [ “ us- east-2c“ ] 
    } 
]
将结果限制为仅指定的字段#
以下describe-instances示例使用--query参数仅显示指定实例的AMI ID和标签。

aws ec2 describe-instances \
    --instance-id i-1234567890abcdef0 \
    --query "Reservations[*].Instances[*].[ImageId,Tags[*]]"

aws ec2 describe-instances --query 'Reservations[*].Instances[*].{Instance:InstanceId,Subnet:SubnetId}' --output json
[
    {
        "Instance": "i-057750d42936e468a",
        "Subnet": "subnet-069beee9b12030077"
    },
    {
        "Instance": "i-001efd250faaa6ffa",
        "Subnet": "subnet-0b715c6b7db68927a"
    },
    {
        "Instance": "i-027552a73f021f3bd",
        "Subnet": "subnet-0250c25a1f4e15235"
    }
    ...
]
