# create_ec2_alarm

## Introduce

该脚本用于一次性为现有的EC2创建status check alarms，包括: `StatusCheckFailed_Instance`和`StatusCheckFailed_System`。

并定义警报发生时EC2的auto-recover action

脚本可以重复运行

限制条件：
 - 正在running instance
 - instance name prefix（定义在ec2_alarm.yaml）

[Amazon CloudWatch 定价](https://aws.amazon.com/cn/cloudwatch/pricing/?nc1=h_ls)

每个警报指标 0.10 USD

## Command

```bash
cd scripts/cloudwatch_alarm
# 创建虚拟环境
python3 -m venv .venv
# 启动虚拟环境
source .venv/bin/activate
# 安装依赖
pip3 install -r requirements.txt

# python3 update_exists_ec2_alarm.py --profile staging --env staging --region xxx
python3 create_ec2_alarm.py --profile staging --env staging
python3 create_ec2_alarm.py --profile staging --env staging --instanceid i-asd71h2kjasdxxx
python3 create_ec2_alarm.py --profile prod --env prod

# 退出虚拟环境
deactivate
```

# check_all_alarm

## Introduce

该脚本用于检查所有的alarms

该脚本会列出Alarm中Action没有enable或Alarm没有对应的Action

脚本可以重复运行

## Command

```bash
cd scripts/cloudwatch_alarm
# 创建虚拟环境
python3 -m venv .venv
# 启动虚拟环境
source .venv/bin/activate
# 安装依赖
pip3 install -r requirements.txt

python3 check_all_alarm.py --profile staging
python3 check_all_alarm.py --profile prod

# 退出虚拟环境
deactivate
```

# filter_all_alarm

## Introduce

该脚本用于找到指定环境中的alarms以`prefix`开头的文件,并将找到的alarms写入`filter_alarms.txt`中

`delete_alarm.py`的脚本会读取并删除文件中的alarms

脚本运行删除前会有确认提示

## Command

```bash
cd scripts/cloudwatch_alarm
# 创建虚拟环境
python3 -m venv .venv
# 启动虚拟环境
source .venv/bin/activate
# 安装依赖
pip3 install -r requirements.txt

python3 filter_all_alarm.py --profile staging --prefix EC2
python3 filter_all_alarm.py --profile prod --prefix EC2

# 退出虚拟环境
deactivate
```

# delete_alarm

## Introduce

该脚本用于删除在`filter_alarms.txt`中的alarms

`filter_alarms.txt`中的alarms可以通过`filter_all_alarm.py`脚本生成

脚本运行删除前会有确认提示

## Command

```bash
cd scripts/cloudwatch_alarm
# 创建虚拟环境
python3 -m venv .venv
# 启动虚拟环境
source .venv/bin/activate
# 安装依赖
pip3 install -r requirements.txt

python3 delete_alarm.py --profile staging
python3 delete_alarm.py --profile prod

# 退出虚拟环境
deactivate
```
