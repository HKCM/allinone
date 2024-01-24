#!/usr/bin/env bash

# Run this script in AWS floder
# Like this: bash ./Config/deploy.sh

# Please setup stack name and your aws profile
STACK_NAME=XXX
PROFILE=XXX

if [ "${STACK_NAME}" = "XXX" ]; then
    echo "****** Please set STACK_NAME and PROFILE first ******"
    echo "******                                         ******"
    echo "******    They are in line 7 of this script    ******"
    exit 1
fi

# Don't forget confirm SubscriptionConfirmation Email
DeliveryEmail=xxxxxxxxx@example.com

# When EnableCompliantMessage is set to true, 
# Messages will be sent even the resource compliant the role
EnableCompliantMessage=false

# Use "true" to set the rules you want to enable, 
# And "false" to set the rules you don't want to enable
SGonlyAllowedPort80443=false
RDSInstancePublicAccess=false 
RDSMultiAZ=false 
RDSSnapshotPublic=false 
S3BucketEncryption=false 
EBSSnapshotPublic=false 
ACMExpiration=false 
RequiredTags=false
LambdaFunctionName=${STACK_NAME}-lambda

zip ./Config/lambda.zip ./Config/lambda.py

aws cloudformation deploy \
    --profile ${PROFILE} \
    --template-file ./Config/config.yml \
    --capabilities CAPABILITY_NAMED_IAM \
    --stack-name ${STACK_NAME} \
    --tags \
        Name=${STACK_NAME} \
        Department=dev \
        Environment=xxxxx \
        Team=xxxxx \
    --parameter-overrides \
        DeliveryEmail=${DeliveryEmail} \
        EnableCompliantMessage=${EnableCompliantMessage} \
        SGonlyAllowedPort80443=${SGonlyAllowedPort80443} \
        RDSInstancePublicAccess=${RDSInstancePublicAccess} \
        RDSMultiAZ=${RDSMultiAZ} \
        RDSSnapshotPublic=${RDSSnapshotPublic} \
        S3BucketEncryption=${S3BucketEncryption} \
        EBSSnapshotPublic=${EBSSnapshotPublic} \
        ACMExpiration=${ACMExpiration} \
        LambdaFunctionName=${LambdaFunctionName} \
        RequiredTags=${RequiredTags} 

aws lambda update-function-code \
    --profile ${PROFILE} \
    --function-name ${LambdaFunctionName} \
    --zip-file fileb://./Config/lambda.zip


    