#!/usr/bin/env bash
RoleARN=xxxxx
Profile=hkc
Region=us-east-1

KST=($(aws sts assume-role \
  --profile ${Profile} \
  --role-arn ${RoleARN} \
  --role-session-name session \
  --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]' \
  --external-id agw \
  --output text))
if [ $? != 0 ]; then
  exit 1
fi

echo "export AWS_ACCESS_KEY_ID=${KST[0]}"

echo "export AWS_SECRET_ACCESS_KEY=${KST[1]}"

echo "export AWS_SESSION_TOKEN=${KST[2]}"

echo "export AWS_REGION=${Region}"

# aws configure set aws_access_key_id ${KST[0]} --profile ${awsProfile}
# aws configure set aws_secret_access_key ${KST[1]} --profile ${awsProfile}
# aws configure set aws_session_token ${KST[2]} --profile ${awsProfile}
# aws configure set region us-east-1 --profile ${awsProfile}

# echo "Show me the AWS config:"
# cat ~/.aws/config

# echo "Show me the AWS credentials:"
# cat ~/.aws/credentials

