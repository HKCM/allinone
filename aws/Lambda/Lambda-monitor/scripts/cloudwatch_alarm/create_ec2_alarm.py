import logging
import yaml
import boto3
import argparse


logger = logging.getLogger()
logger.setLevel(logging.INFO)

formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')

DEFAULT_INSTANCE="None"

parser = argparse.ArgumentParser(description='命令行参数')
parser.add_argument('--profile', '-p', type=str, help='AWS Profile', required=True)
parser.add_argument('--env', '-e', type=str, help='对应ec2_alarm env', required=True)
parser.add_argument('--region', '-r', type=str, help='AWS Region', default='ap-northeast-1')
parser.add_argument('--instanceid', '-i', type=str, help='对单台instance生效', default=DEFAULT_INSTANCE)
args = vars(parser.parse_args())

EC2_ALARM_CONFIG = "ec2_alarm.yaml"
PROFILE = args["profile"]
REGION = args["region"]
INSTANCE_ID = args["instanceid"]
ENV = args["env"]
SESSION = boto3.Session(profile_name=PROFILE,region_name=REGION)

CONFIG_ITEM = {}

def main():

    print_env()
    global CONFIG_ITEM
    CONFIG_ITEM = get_config(ENV)
    # 如果指定了instance ID
    if INSTANCE_ID != DEFAULT_INSTANCE:
        set_alarm_for_single_instance(INSTANCE_ID)

    # 如果没有指定 instance ID
    if INSTANCE_ID == DEFAULT_INSTANCE:
        prefix = CONFIG_ITEM["prefix"]
        set_alarm_for_instances(prefix)

def print_env():

    print("----------------------------------")
    print("profile:" + PROFILE)
    print("region:" + REGION)
    print("env:" + ENV)
    print("instance:" + INSTANCE_ID)
    print("----------------------------------")

def get_config(target_env:str) -> dict:
    with open(EC2_ALARM_CONFIG, encoding="utf-8") as file:
        # 将文件内容解析为 YAML
        config = yaml.load(file, Loader=yaml.FullLoader)
    for item in config:
        if item["env"] == target_env:
            return item
    return {}

def set_alarm_for_instances(prefix:str):
    """
    set status check alarm for running and prefix instances
    """
    filters = [{'Name': 'instance-state-name', 'Values': ['running']}]
    filters.append({'Name': 'tag:Name', 'Values': [prefix + '*']})

    ec2_client = SESSION.client('ec2', REGION)
    # 使用 describe_instances() 方法查询符合条件的实例信息
    response = ec2_client.describe_instances(Filters=filters)

    # 处理第一页查询结果
    for reservation in response['Reservations']:
        for instance in reservation['Instances']:
            instance_id= instance['InstanceId']
            instance_name = ''
            for tag in instance['Tags']:
                if tag['Key'] == 'Name':
                    instance_name = tag['Value']
                    break
            print('instance id: %s; instance name: %s', instance_id, instance_name)
            logger.info('instance id: %s; instance name: %s', instance_id, instance_name)
            if instance_name.lower().startswith(prefix.strip().lower()):
                update_alarm_group(instance_id, instance_name)

    # 如果查询结果超过了 1000 条，则继续查询下一页的结果
    while 'NextToken' in response:
        next_token = response['NextToken']
        response = ec2_client.describe_instances(Filters=filters, NextToken=next_token)

        # 处理下一页查询结果
        for reservation in response['Reservations']:
            for instance in reservation['Instances']:
                instance_id= instance['InstanceId']
                instance_name = ''
                for tag in instance['Tags']:
                    if tag['Key'] == 'Name':
                        instance_name = tag['Value']
                        break
                print('instance id: %s; instance name: %s', instance_id, instance_name)
                if instance_name.lower().startswith(prefix.strip().lower()):
                    update_alarm_group(instance_id, instance_name)
                if instance_name.lower().startswith(prefix.strip().lower()):
                    update_alarm_group(instance_id, instance_name)

def set_alarm_for_single_instance(instance_id:str):
    """
    set status check alarm for running and prefix instances
    """
    filters = [{'Name': 'instance-state-name', 'Values': ['running']}]

    ec2_client = SESSION.client('ec2', REGION)

    response = ec2_client.describe_instances(
        Filters=filters,
        InstanceIds=[
            instance_id,
        ],
    )

    if len(response['Reservations']) == 0:
        print(instance_id + " is not running, and not able to set alarm...")
        return

    for reservation in response['Reservations']:
        for instance in reservation['Instances']:
            instance_id= instance['InstanceId']
            instance_name = ''
            for tag in instance['Tags']:
                if tag['Key'] == 'Name':
                    instance_name = tag['Value']
                    break
            logger.info('instance id: %s; instance name: %s', instance_id, instance_name)
            update_alarm_group(instance_id, instance_name)

def update_alarm_group(instanceid, instance_name):
    # 访问配置数据
    print(f'Create alarm for {CONFIG_ITEM["env"]} {instance_name}...')
    tags = CONFIG_ITEM['tags']
    for alarm in CONFIG_ITEM['alarms']:
        update_ec2_alarm(tags, alarm, instanceid, instance_name)


def update_ec2_alarm(tags:list[dict],alarm:dict, instanceid:str, instance_name:str):
    cloudwatch = SESSION.client('cloudwatch')
    alarm_name_pattern = alarm['alarmName']
    alarm_name = alarm_name_pattern.replace('{INSTANCE_NAME}', instance_name).replace('{INSTANCE_ID}', instanceid)
    logger.info('trying to update alarm: %s', alarm_name)
    cloudwatch.put_metric_alarm(
        Namespace ='AWS/EC2',
        Dimensions=[
            {
                'Name': 'InstanceId',
                'Value': instanceid,
            },
        ],
        AlarmName= alarm_name,
        MetricName=alarm['metricName'],
        Statistic=alarm['statistic'],
        Period=alarm['period'],
        EvaluationPeriods=alarm['evaluationPeriods'],
        Threshold=alarm['threshold'],
        ComparisonOperator = alarm['comparisonOperator'],
        TreatMissingData =  alarm['treatMissingData'],
        OKActions=alarm['okActions'],
        AlarmActions=alarm['alarmActions'],
        InsufficientDataActions=alarm.get('insufficientDataActions',[]),
        Tags=tags,
    )
    logger.info('alarm updated: %s', alarm_name)

if __name__ == "__main__":
    main()
