# 首先user-data必须以/bin/bash开头
aws ec2 run-instances \
--image-id ami-0000025f7c02a13b2 \
--count 1 \
--instance-type t2.micro \
--user-data $'#!/bin/bash\nyum install git -y'

# 通过文件可以正确执行
aws ec2 run-instances \
--image-id ami-0000025f7c02a13b2 \
--count 1 \
--instance-type t2.micro \
--user-data file://path/to/script.sh

# 将user-data放到
aws ec2 run-instances \
--image-id ami-0000025f7c02a13b2 \
--count 1 \
--instance-type t2.micro \
--user-data $'#!/bin/bash
yum install git -y
echo 123 > /tmp/123.txt'

# 以下是错误的
aws ec2 run-instances \
--image-id ami-0000025f7c02a13b2 \
--count 1 \
--instance-type t2.micro \
--user-data "yum install git -y"
