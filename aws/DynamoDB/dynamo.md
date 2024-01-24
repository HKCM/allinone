### 创建table
```python
import boto3

# Get the service resource.
session = boto3.Session(profile_name='dev', region_name='region')
dynamodb = session.client('dynamodb')

# Create the DynamoDB table.
table = dynamodb.create_table(
    TableName='users',
    KeySchema=[
        {
            'AttributeName': 'username',
            'KeyType': 'HASH'
        },
        {
            'AttributeName': 'last_name',
            'KeyType': 'RANGE'
        }
    ],
    AttributeDefinitions=[
        {
            'AttributeName': 'username',
            'AttributeType': 'S'
        },
        {
            'AttributeName': 'last_name',
            'AttributeType': 'S'
        },
    ],
    BillingMode='PAY_PER_REQUEST',
    Tags=[
        {
            'Key': 'string0',
            'Value': 'string0'
        },
        {
            'Key': 'string1',
            'Value': 'string1'
        },
    ],
)

# Wait until the table exists.
table.wait_until_exists()

# Print out some data about the table.
print(table.item_count)
```

### 删除表
```python
import boto3

# Get the service resource.
session = boto3.Session(profile_name='dev', region_name='region')
client = session.client('dynamodb')
response = client.delete_table(
    TableName='string'
)
print(response)
```


### 使用现有的table
```python
import boto3

# Get the service resource.
session = boto3.Session(profile_name='dev', region_name='region')
dynamodb = session.resource('dynamodb')
table = dynamodb.Table('users')
```

### 插入数据
```python
import boto3

# Get the service resource.
session = boto3.Session(profile_name='dev', region_name='region')
dynamodb = session.resource('dynamodb')
table = dynamodb.Table('users')

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

with table.batch_writer() as batch:
    for _ in range(1000000):
        batch.put_item(Item={'HashKey': '...',
                             'Otherstuff': '...'})


```




### 获取item
```python
import boto3

# Get the service resource.
session = boto3.Session(profile_name='dev', region_name='region')
dynamodb = session.client('dynamodb')
table = dynamodb.Table('users')

response = table.get_item(
    Key={
        'username': 'janedoe',
        'last_name': 'Doe'
    }
)
item = response['Item']
print(item)
```