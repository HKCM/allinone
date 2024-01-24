import os

def lambda_handler(event,context):
    Word = os.environ['Word']
    print("Hello" + Word)
    print("event")
    print(event)