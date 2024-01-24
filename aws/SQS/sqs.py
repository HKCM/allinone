import boto3
session = boto3.Session(profile_name='dev')
client = session.client('sqs', region_name='ap-southeast-1')

queueUrl = 'https://sqs.ap-southeast-1.amazonaws.com/111111111111/sqsqueue'
for i in range(1,12):
    response = client.receive_message(
        QueueUrl=queueUrl,
        AttributeNames=['All'],
        MessageAttributeNames=['All'],
        MaxNumberOfMessages=10,
        VisibilityTimeout=5,
        WaitTimeSeconds=10
    )
    for message in response['Messages']:
        receiptHandle = message['ReceiptHandle']
        client.delete_message(QueueUrl=queueUrl,ReceiptHandle=receiptHandle)
    print(i)