import boto3
import argparse
from botocore.exceptions import ClientError

parser = argparse.ArgumentParser(description='命令行参数')
parser.add_argument('--action', choices=['check', 'terminate-protect', 'stop-protect', 'both-protect'],  required=True)
parser.add_argument('--profile', '-p', type=str, help='AWS Profile', required=True)
parser.add_argument('--region', '-r', type=str, help='AWS Region', default='ap-northeast-1')
# parser.add_argument('--terminate-protect', action=argparse.BooleanOptionalAction, help='enable terminate-protect')
# parser.add_argument('--stop-protect', action=argparse.BooleanOptionalAction,help='enable stop-protect')
# parser.add_argument('--check', action=argparse.BooleanOptionalAction)

def comma_separated_list(value):
    return value.split(',')
parser.add_argument('--instance-ids', '-i', type=comma_separated_list, help='对指定instance生效以[,]分隔')
args = vars(parser.parse_args())
PROFILE = args["profile"]
REGION = args["region"]
ACTION = args["action"]
# TERMINATE_PROTECT = args["terminate_protect"]
# STOP_PROTECT = args["stop_protect"]
# CHECK = args["check"]
SESSION = boto3.Session(profile_name=PROFILE,region_name=REGION)
INSTANCE_IDS = args["instance_ids"]

def main():

    print_env()

    if ACTION == "check":
        check(INSTANCE_IDS)
        print("Check done")
        return
    
    if ACTION == "terminate-protect":
        set_protection_for_instances(
            instance_ids=INSTANCE_IDS,
            disable_stop=False,
            disable_terminate=True
        )

    if ACTION == "stop-protect":
        set_protection_for_instances(
            instance_ids=INSTANCE_IDS,
            disable_stop=True,
            disable_terminate=False
        )

    if ACTION == "both-protect":
        set_protection_for_instances(
            instance_ids=INSTANCE_IDS,
            disable_stop=True,
            disable_terminate=True
        )

    if ACTION not in ["check","stop-protect","terminate-protect","both-protect"]:
        print("Unsupported Action!!! Please run script with --help")

    print("Bye")

def print_env():
    print("---------------------------------")
    print(f"PROFILE :     {PROFILE}")
    print(f"REGION  :     {REGION}")
    print(f"INSTANCE-IDs: {INSTANCE_IDS}")
    print(f"ACTION:       {ACTION}")
    print("---------------------------------")

def set_protection_for_instances(instance_ids:list[str],disable_stop:bool,disable_terminate:bool):
    """
    set protection for ALL running instance
    disable_stop = True 表示启用Stop保护
    disable_terminate = True 表示启用Terminate保护
    """
    if instance_ids is not None:
        target_instances_list = instance_ids
    else:
        target_instances_list = []

    ec2_client = SESSION.client('ec2', REGION)
    # 使用 describe_instances() 方法查询符合条件的实例信息
    response = ec2_client.describe_instances(
        InstanceIds=target_instances_list
    )

    # 处理第一页查询结果
    for reservation in response['Reservations']:
        for instance in reservation['Instances']:
            instance_id = instance['InstanceId']
            instance_name = ''
            for tag in instance['Tags']:
                if tag['Key'] == 'Name':
                    instance_name = tag['Value']
                    break
            
            if disable_terminate:
                response = ec2_client.describe_instance_attribute(
                    Attribute='disableApiTermination',
                    InstanceId=instance_id,
                )
                if response["DisableApiTermination"]["Value"] is False:
                    print(instance_name + " " + instance_id + ": 没有启用Termination保护")
                    set_terminate_protection(instance_id)
                # else:
                #     print(instance_name + " " + instance_id + ": 已启用Termination保护")

            if disable_stop:
                response = ec2_client.describe_instance_attribute(
                    Attribute='disableApiStop',
                    InstanceId=instance_id,
                )
                if response["DisableApiStop"]["Value"] is False:
                    print(instance_name + " " + instance_id + ": 没有启用Stop保护")
                    set_stop_protection(instance_id)
                # else:
                #     print(instance_name + " " + instance_id + ": 已启用Stop保护")
            
    # 如果查询结果超过了 1000 条，则继续查询下一页的结果
    while 'NextToken' in response:
        next_token = response['NextToken']
        response = ec2_client.describe_instances(NextToken=next_token)

        # 处理下一页查询结果
        for reservation in response['Reservations']:
            for instance in reservation['Instances']:
                instance_id= instance['InstanceId']
                instance_name = ''
                for tag in instance['Tags']:
                    if tag['Key'] == 'Name':
                        instance_name = tag['Value']
                        break
                
                if disable_terminate:
                    response = ec2_client.describe_instance_attribute(
                        Attribute='disableApiTermination',
                        InstanceId=instance_id,
                    )
                    if response["DisableApiTermination"]["Value"] is False:
                        print(instance_name + " " + instance_id + ": 没有启用Termination保护")
                        set_terminate_protection(instance_id)
                    # else:
                    #     print(instance_name + " " + instance_id + ": 已启用Termination保护")

                if disable_stop:
                    response = ec2_client.describe_instance_attribute(
                        Attribute='disableApiStop',
                        InstanceId=instance_id,
                    )
                    if response["DisableApiStop"]["Value"] is False:
                        print(instance_name + " " + instance_id + ": 没有启用Stop保护")
                        set_stop_protection(instance_id)
                    # else:
                    #     print(instance_name + " " + instance_id + ": 已启用Termination保护")

def check(instance_ids:list[str]):
    """
    check instances protection settings
    """
    no_terminate_protection_instances = []
    no_stop_protection_instances = []
    no_terminate_protection_instances_id = []
    no_stop_protection_instances_id = []

    if instance_ids is not None:
        target_instances_list = instance_ids
    else:
        target_instances_list = []

    print("Start checking...")
    ec2_client = SESSION.client('ec2', REGION)
    # 使用 describe_instances() 方法查询符合条件的实例信息
    response = ec2_client.describe_instances(
        InstanceIds=target_instances_list
    )

    # 处理第一页查询结果
    for reservation in response['Reservations']:
        for instance in reservation['Instances']:
            instance_id = instance['InstanceId']
            instance_name = ''
            for tag in instance['Tags']:
                if tag['Key'] == 'Name':
                    instance_name = tag['Value']
                    break
            print(f"Checking {instance_name}...")
            response = ec2_client.describe_instance_attribute(
                Attribute='disableApiTermination',
                InstanceId=instance_id,
            )
            if response["DisableApiTermination"]["Value"] is False:
                no_terminate_protection_instances.append(instance_name)
                no_terminate_protection_instances_id.append(instance_id)
                # print(instance_name + " " + instance_id + ": 没有启用Termination保护")

            response = ec2_client.describe_instance_attribute(
                Attribute='disableApiStop',
                InstanceId=instance_id,
            )
            if response["DisableApiStop"]["Value"] is False:
                no_stop_protection_instances.append(instance_name)
                no_stop_protection_instances_id.append(instance_id)
                # print(instance_name + " " + instance_id + ": 没有启用Stop保护")

    # 如果查询结果超过了 1000 条，则继续查询下一页的结果
    while 'NextToken' in response:
        next_token = response['NextToken']
        response = ec2_client.describe_instances(NextToken=next_token)

        # 处理下一页查询结果
        for reservation in response['Reservations']:
            for instance in reservation['Instances']:
                instance_id= instance['InstanceId']
                instance_name = ''
                for tag in instance['Tags']:
                    if tag['Key'] == 'Name':
                        instance_name = tag['Value']
                        break
                print(f"Checking {instance_name}...")
                response = ec2_client.describe_instance_attribute(
                    Attribute='disableApiTermination',
                    InstanceId=instance_id,
                )
                if response["DisableApiTermination"]["Value"] is False:
                    no_terminate_protection_instances.append(instance_name)
                    no_terminate_protection_instances_id.append(instance_id)
                    # print(instance_name + " " + instance_id + ": 没有启用Termination保护")
                # else:
                #     print(instance_name + " " + instance_id + ": 已启用Termination保护")

                response = ec2_client.describe_instance_attribute(
                    Attribute='disableApiStop',
                    InstanceId=instance_id,
                )
                if response["DisableApiStop"]["Value"] is False:
                    no_stop_protection_instances.append(instance_name)
                    no_stop_protection_instances_id.append(instance_id)
                    #print(instance_name + " " + instance_id + ": 没有启用Stop保护")
                # else:
                #     print(instance_name + " " + instance_id + ": 已启用Stop保护")
    print("---------------------------------")
    print("no_stop_protection_instances:")
    print(no_stop_protection_instances)
    print("no_stop_protection_instances_id:")
    print(no_stop_protection_instances_id)
    print("---------------------------------")
    print("no_terminate_protection_instances:")
    print(no_terminate_protection_instances)
    print("no_terminate_protection_instances_id:")
    print(no_terminate_protection_instances_id)

def set_stop_protection(instance_id:str):
    """
    set_stop_protection
    """

    ec2_client = SESSION.client('ec2', REGION)
    try:
        response = ec2_client.modify_instance_attribute(
            DisableApiStop={
                'Value': True,
            },
            InstanceId=instance_id,
        )
        
    except ClientError as e:
        if e.response['Error']['Code'] == 'UnsupportedOperation':
            print("spot instances not support protection...Skip...")
        else:
            print(f"捕获到其他客户端错误：{e}")
    except Exception as e:
        print(f"捕获到其他异常：{e}")
    else:
        print(instance_id + ": Enable Stop protection")

def set_terminate_protection(instance_id:str):
    """
    set_terminate_protection
    """

    ec2_client = SESSION.client('ec2', REGION)
    try:
        response = ec2_client.modify_instance_attribute(
            DisableApiTermination={
                'Value': True,
            },
            InstanceId=instance_id,
        )
    except ClientError as e:
        if e.response['Error']['Code'] == 'UnsupportedOperation':
            print("spot instances not support protection...Skip...")
        else:
            print(f"捕获到其他客户端错误：{e}")
    except Exception as e:
        print(f"捕获到其他异常：{e}")
    else:
        print(instance_id + ": Enable Terminate protection")
            
if __name__ == "__main__":
    main()
