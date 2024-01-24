import boto3

session = boto3.Session(profile_name='int-xmn', region_name='us-east-1')
dynamodb = session.resource('dynamodb')

table = dynamodb.Table('int-ips-performance-MsSubscription')

def put(table):
    table.put_item(
    Item={
            'subscription_id': '87d2aa79-323e-45a1-893b-985e40607439',
            'expiration_time': 1651219200,
            'failureType': '500-ExtensionError',
            'lastRunState': 'fail',
            'account_id': '123456789',
            'renewAt': '2022-04-29T07:30:00.150Z',
            'renewFailedCount':1
        }
    )

