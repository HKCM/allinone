import logging
import boto3
import argparse

logger = logging.getLogger()
logger.setLevel(logging.INFO)

formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')

parser = argparse.ArgumentParser(description='命令行参数')
parser.add_argument('--profile', '-p', type=str, help='AWS Profile', required=True)
parser.add_argument('--region', '-r', type=str, help='AWS Region', default='ap-northeast-1')
args = vars(parser.parse_args())

profile = args["profile"]
region = args["region"]
session = boto3.Session(profile_name=profile,region_name=region)
abnormal_alarm = []

def main():
    """main function"""
    print(profile)
    print(region)
    filter_all_alarms()
    print("------------- Result -------------")
    if len(abnormal_alarm) == 0:
        print("\nGreat: Each alert is enabled and has action\n")
    else:
        print("\nSome alert are abnormal(ActionsDisabled or NoAction):")
        print(abnormal_alarm)
        print()

def filter_all_alarms():
    """filter_all_alarms"""
    cw_client = session.client('cloudwatch')
    response = cw_client.describe_alarms(
        MaxRecords=100,
    )

    # 没有处理 CompositeAlarms
    # if len(response["CompositeAlarms"])> 0:
    #     pass

    for metric_alarm in response["MetricAlarms"]:
        filter_alarm(metric_alarm)

    # 如果查询结果超过了 100 条，则继续查询下一页的结果
    while 'NextToken' in response:
        next_token = response['NextToken']
        response = cw_client.describe_alarms(
            MaxRecords=100,
            NextToken=next_token
        )
        # 处理下一页查询结果
        for metric_alarm in response["MetricAlarms"]:
            filter_alarm(metric_alarm)

def filter_alarm(metric_alarm:dict):
    """filter single alarm"""
    print(f'checking {metric_alarm["AlarmName"]}...')
    # 如果Alarm的Action没有Enable
    if metric_alarm["ActionsEnabled"] is not True:
        abnormal_alarm.append(metric_alarm["AlarmName"])
        return

    # 如果Alarm没有Action
    if len(metric_alarm["OKActions"]) == 0 and \
        len(metric_alarm["AlarmActions"]) == 0 and \
        len(metric_alarm["InsufficientDataActions"]) == 0:
        abnormal_alarm.append(metric_alarm["AlarmName"])
        return

if __name__ == "__main__":
    main()
