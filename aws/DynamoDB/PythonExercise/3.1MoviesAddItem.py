
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
import json
import decimal

# Helper class to convert a DynamoDB item to JSON.
class DecimalEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, decimal.Decimal):
            if abs(o) % 1 > 0:
                return float(o)
            else:
                return int(o)
        return super(DecimalEncoder, self).default(o)

#dynamodb = boto3.resource('dynamodb', region_name='us-west-2', endpoint_url="http://localhost:8000")
session = boto3.session.Session(profile_name='eu')
dynamodb = session.resource('dynamodb')
table = dynamodb.Table('Movies')

title = "hkc"
year = 2015

response = table.put_item(
   Item={
        'year': year,
        'title': title,
        # 'info': {
        #     'plot':"Nothing happens at all.",
        #     'rating': decimal.Decimal(0)
        # }
        'info': "hkc2015"
    },
    #Functions: attribute_exists | attribute_not_exists | attribute_type | contains | begins_with | size These function names are case-sensitive.
    #Comparison operators: = | <> | < | > | <= | >= | BETWEEN | IN
    #Logical operators: AND | OR | NOT
    
    # ----------------title not hkc -------------
    # ConditionExpression='NOT title IN (:tl)',
    # ExpressionAttributeValues={ ":tl":"hkc"}
    # ----------------nto title attr -------------
    ConditionExpression='attribute_not_exists(title)',

    #ConditionExpression = Attr('title').eq('The Big New Movie'),
)

print("PutItem succeeded:")
print(json.dumps(response, indent=4, cls=DecimalEncoder))
