import json
import boto3
import os


def lambda_handler(event, context):
    EnableCompliantMessage = os.environ['EnableCompliantMessage']
    print("**event**")
    print(event)
    resourceId = event['detail']['resourceId']
    resourceType = event['detail']['resourceType']
    configRuleName = event['detail']['configRuleName']
    complianceType = event['detail']['newEvaluationResult']['complianceType']
    region = event['region']
    print("==>Info: resourceId:{} resourceType:{} configRuleName:{} complianceType:{} region:{}".format(\
        resourceId,resourceType,configRuleName,complianceType,region))

    # When EnableCompliantMessage set False, it means if resource Compliant the role
    # We don't need send message to stakeholder
    if complianceType == "COMPLIANT" and EnableCompliantMessage == "false":
        return

    userName = find_username(resourceId)
    if "botocore" in userName:
        userName = find_username(userName)
    
    message = "There is one resource **{}** **{}** the **{}** role. \
        ResourceType is **{}**, The action is performed by **{}** in **{}**".format( \
        resourceId,complianceType,configRuleName,resourceType,userName,region)
    print("==>message" + message)
    sendMessage(message)

def find_username(resourceId):
    cloudtrail_client = boto3.client('cloudtrail')
    events = cloudtrail_client.lookup_events(LookupAttributes=[{'AttributeKey':'ResourceName', 'AttributeValue':resourceId}],MaxResults=1)
    print("==>Trail_Events")
    print(events)
    Actor = events['Events'][0]['Username']
    print("==>Actor" + Actor)
    if(Actor):
        return Actor
    else:
        return "Unknow"



def sendMessage(message):
    client = boto3.client('sns')
    TopicArn = os.environ['TopicArn']
    print("==>topic" + TopicArn)
    print("==>message" + message)
    client.publish(TopicArn=TopicArn,Message=message,)