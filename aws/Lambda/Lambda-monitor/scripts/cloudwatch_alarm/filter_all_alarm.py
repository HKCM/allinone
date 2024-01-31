import logging
import boto3
import argparse

logger = logging.getLogger()
logger.setLevel(logging.INFO)

formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')

parser = argparse.ArgumentParser(description='命令行参数')
parser.add_argument('--profile', '-p', type=str, help='AWS Profile', required=True)
parser.add_argument('--region', '-r', type=str, help='AWS Region', default='ap-northeast-1')
parser.add_argument('--prefix', type=str, help='AWS Region', default='EC2')
args = vars(parser.parse_args())

profile = args["profile"]
region = args["region"]
prefix = args["prefix"]
session = boto3.Session(profile_name=profile,region_name=region)
abnormal_alarm = []

def main():
    """main function"""
    print(profile)
    print(region)
    filter_all_alarms(prefix)

def filter_all_alarms(alarm_name_prefix:str):
    """filter_all_alarms"""
    num = 0
    cw_client = session.client('cloudwatch')
    response = cw_client.describe_alarms(
        AlarmNamePrefix=alarm_name_prefix,
        MaxRecords=100,
    )

    # 没有处理 CompositeAlarms
    # if len(response["CompositeAlarms"])> 0:
    #     pass

    # 将alarm name 写入文件
    with open('./filter_alarms.txt', 'w', encoding="utf-8") as file:
        for metric_alarm in response["MetricAlarms"]:
            num += 1
            file.write(metric_alarm["AlarmName"] + '\n')

    # 如果查询结果超过了 100 条，则继续查询下一页的结果
    while 'NextToken' in response:
        next_token = response['NextToken']
        response = cw_client.describe_alarms(
            AlarmNamePrefix=alarm_name_prefix,
            MaxRecords=100,
            NextToken=next_token
        )
        # 处理下一页查询结果
        with open('./filter_alarms.txt', 'w', encoding="utf-8") as file:
            for metric_alarm in response["MetricAlarms"]:
                num += 1
                file.write(metric_alarm["AlarmName"] + '\n')

    print(f"There are totally {num} alarms")

if __name__ == "__main__":
    main()
