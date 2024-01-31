import logging
import boto3
import argparse

logger = logging.getLogger()
logger.setLevel(logging.INFO)

formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')

parser = argparse.ArgumentParser(description='命令行参数')
parser.add_argument('--profile', '-p', type=str, help='AWS Profile', required=True)
parser.add_argument('--region', '-r', type=str, help='AWS Region', default='ap-northeast-1')
parser.add_argument('--file_path', '-f', type=str, help='alarm name file', default='./filter_alarms.txt')
args = vars(parser.parse_args())

Profile = args["profile"]
Region = args["region"]
File_path = args["file_path"]
Session = boto3.Session(profile_name=Profile,region_name=Region)

def main():
    """main function"""
    delete_alarm(File_path)

def delete_alarm(file_path:str):
    """delete alarms that in file"""
    alarm_name_list = []
    print("即将删除的CloudWatch Alarms:\n")
    try:
        with open(file_path, "r",encoding="utf-8") as file:
            lines = file.readlines()
            for line in lines:
                alarm_name = line.strip()
                alarm_name_list.append(alarm_name) 
                print(alarm_name)
    except FileNotFoundError:
        print("文件未找到")

    confirmation = input("即将删除以上alarms,请输入 'Yes' 确认操作: ")
    if confirmation == "Yes":
        Session.client('cloudwatch').delete_alarms(
            AlarmNames=alarm_name_list
        )
    else:
        print("操作未确认,退出")

if __name__ == "__main__":
    main()
