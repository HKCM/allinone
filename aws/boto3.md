安装
```python
python3 -m pip install -r boto3
```

```python
import boto3

session = boto3.Session(profile_name='dev', region_name='region')
dev_s3_client = session.client('s3')
```

```python
import boto3
#
# setting up configured profile on your machine.
# You can ignore this step if you want use default AWS CLI profile.
#
boto3.setup_default_session(profile_name='admin-analyticshut')
s3 = boto3.client('s3')
```

```python
import boto3
# Hard coded strings as credentials, not recommended.
client = boto3.client(
    's3',
    aws_access_key_id=ACCESS_KEY,
    aws_secret_access_key=SECRET_KEY,
    aws_session_token=SESSION_TOKEN,
)
```
