https://boto3.amazonaws.com/v1/documentation/api/latest/guide/configuration.html

boto3 credentail 配置

1. Method Parameters

The first option for providing credentials to boto3 is passing them as parameters when creating clients or when creating a Session. For example:

```
import boto3
client = boto3.client(
    's3',
    aws_access_key_id=ACCESS_KEY,
    aws_secret_access_key=SECRET_KEY,
    aws_session_token=SESSION_TOKEN,
)

# Or via the Session
session = boto3.Session(
    aws_access_key_id=ACCESS_KEY,
    aws_secret_access_key=SECRET_KEY,
    aws_session_token=SESSION_TOKEN,
)
```

2. Environment Variables

Boto3 will check these environment variables for credentials:

**AWS_ACCESS_KEY_ID**
The access key for your AWS account.
**AWS_SECRET_ACCESS_KEY**
The secret key for your AWS account.
**AWS_SESSION_TOKEN**
The session key for your AWS account. 

This is only needed when you are using temporary credentials. The AWS_SECURITY_TOKEN environment variable can also be used, but is only supported for backwards compatibility purposes. AWS_SESSION_TOKEN is supported by multiple AWS SDKs besides python.

3. Shared Credentials File

You can then specify a profile name via the AWS_PROFILE environment variable or the profile_name argument when creating a Session:

```
session = boto3.Session(profile_name='dev')
# Any clients created from this session will use credentials
# from the [dev] section of ~/.aws/credentials.
dev_s3_client = session.client('s3')
```
