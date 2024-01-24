使用以下命令列出您的可用 AWS CloudFormation 堆栈。在生成的输出中查找 VPC 模板名称。
```shell
aws cloudformation list-stacks --query "StackSummaries[].StackName"
```

使用以下命令删除VPC堆栈,将<my-vpc-stack>替换为您的VPC堆栈名称。

```shell
aws cloudformation delete-stack --stack-name <my-vpc-stack>
```

获取output
```shell
aws cloudformation describe-stacks \
    --stack-name <stack_name> \
    --query "Stacks[0].Outputs[?OutputKey==`<key_we_want>`].OutputValue" \
    --output text
```
