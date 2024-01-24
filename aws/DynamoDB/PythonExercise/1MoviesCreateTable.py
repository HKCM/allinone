
#  https://docs.aws.amazon.com/zh_cn/amazondynamodb/latest/developerguide/GettingStarted.Python.html
#  
#  Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#  
#  This file is licensed under the Apache License, Version 2.0 (the "License").
#  You may not use this file except in compliance with the License. A copy of
#  the License is located at
# 
#  http://aws.amazon.com/apache2.0/
# 
#  This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
#  CONDITIONS OF ANY KIND, either express or implied. See the License for the
#  specific language governing permissions and limitations under the License.
#
from __future__ import print_function # Python 2/3 compatibility
import boto3

# dynamodb = boto3.resource('dynamodb', region_name='us-west-2', endpoint_url="http://localhost:8000")
session = boto3.session.Session(profile_name='eu')
dynamodb = session.resource('dynamodb')


table = dynamodb.create_table(
    TableName='Movies',
    KeySchema=[
        {
            'AttributeName': 'year',
            'KeyType': 'HASH'  #Partition key
        },
        {
            'AttributeName': 'title',
            'KeyType': 'RANGE'  #Sort key
        }
    ],
    AttributeDefinitions=[
        {
            'AttributeName': 'year',
            'AttributeType': 'N'
        },
        {
            'AttributeName': 'title',
            'AttributeType': 'S'
        },

    ],
    ProvisionedThroughput={
        'ReadCapacityUnits': 10,
        'WriteCapacityUnits': 10
    },
    Tags=[
        {
            'Key': 'Name',
            'Value': 'xxx-test'
        },
        {
            'Key': 'Team',
            'Value': 'int-xxx'
        },
        {
            'Key': 'Environment',
            'Value': 'int-xxx'
        },
        {
            'Key': 'Department',
            'Value': 'dev'
        },
    ]
)

print("Table status:", table.table_status)
