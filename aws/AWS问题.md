
### AWS-Not-Authorized-to-Use-Launch-Template
描述: 当使用AutoScaling时遇到报错,You are not authorized to use launch template

#### 问题分析

这个报错代表,当前ASG无法通过launch template启动EC2.这个问题甚至会出现延期,可能你成功创建了ASG和Launch Template,在真正需要启动EC2时才报错.

对于这个问题只能提供解决思路.

导致问题的原因大致有三个:


1. 运行命令的用户缺少权限(可能性低)

运行命令的用户需要具有最低相应权限,但是通常情况如果能通过AWS console界面创建EC2,则说明权限没有问题

以下是最低权限要求:
```
{
  "Action": [
    "ec2:runInstances",
    "ec2:DescribeLaunchTemplates",
    "ec2:DescribeLaunchTemplateVersions",
    "ec2:DescribeVpcs",
    "ec2:DescribeSubnets",
    "ec2:DescribeAvailabilityZones"
  ],
  "Resource": [
    "arn:aws:ec2:ap-southeast-1:xxxxxxxxxxxx:launch-template/lt-05xxxxxxxxxxxxxa5",
    "arn:aws:ec2:ap-southeast-1:xxxxxxxxxxxx:subnet/subnet-0axxxxxxxxxxxxxad",
    "arn:aws:ec2:ap-southeast-1:xxxxxxxxxxxx:security-group/sg-0axxxxxxxxxxxxx78",
    "arn:aws:ec2:ap-southeast-1:xxxxxxxxxxxx:instance/*",
    "arn:aws:ec2:ap-southeast-1:xxxxxxxxxxxx:network-interface/*",
    "arn:aws:ec2:ap-southeast-1:xxxxxxxxxxxx:volume/*",
    "arn:aws:ec2:ap-southeast-1::image/",
    "arn:aws:ec2:ap-southeast-1::snapshot/"
  ],
  "Effect": "Allow"
}
```

如果用户具有相应权限但是无法启动,下一条.

2. Template错误(常见原因)

正如我之前所说,这个报错是在正在启动EC2时才报错的.所以当你创建ASG和Launch Template时,只要没有语法错误都能创建成功.但是其中的值是否真实有效则需要启动时才能确认.

所以很有可能是Template的值出错. 常见于多账户,多可用区使用,账户B或可用区B没有你指定的AMI 或者 Key.
  * AMIID不存在或值错误
  * Keypair不存在或值错误
  * SecurityGroup不存在或值错误
  * ......

3. AWS账户和AWS organize限制导致(常见原因)

这是最多且最复杂的原因.我只能列出一些常见的,具体的需要自己探究

* AWS账户或组织设置限定EC2类型和Template不一致
* AWS账户或组织设置限定资源只能在指定区域创建
* AWS账户或组织设置限定资源必须附带标签
* AWS账户或组织设置限定不能有公共IP
* AWS账户或组织设置限定资源必须加密
* ......


