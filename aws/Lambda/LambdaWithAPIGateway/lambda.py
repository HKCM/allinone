import os
import json

def lambda_handler(event,context):
    Word = os.environ['Word']
    print("Hello" + Word)
    print("event")
    print(event)
    return {
        'statusCode': 200,
        'body': json.dumps(myParam)
    }