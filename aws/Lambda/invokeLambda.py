import json
import boto3

client = boto3.client('lambda')

def lambda_handler(event, context):
    url=event["url"]
    payload={"url": url}
    FunctionName = event["FunctionARN"]
    rate=event["rate"]
    for i in range(rate):
        # TODO: write code...
        response = client.invoke(
            FunctionName=FunctionName,
            InvocationType='Event',
            LogType='None',
            Payload=json.dumps(payload)
        )
        print(response)

# test event
# {
#   "functionName": "TargetLambdaName",
#   "region": "eu-west-1",
#   "account": "{accountID}",
#   "FunctionARN": "arn:aws:lambda:eu-west-1:{accountID}:function:TargetLambdaName",
#   "rate": 50,
#   "url": "https://www.baidu.com/1.html"
# }
