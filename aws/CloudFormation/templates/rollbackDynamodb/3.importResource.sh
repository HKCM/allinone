#!/usr/bin/env bash

# Run this script in AWS/CloudFormation/Template/dynamoRollback/shell folder
# For more detail, please refer https://docs.aws.amazon.com/zh_cn/AWSCloudFormation/latest/UserGuide/resource-import-existing-stack.html

# bash ./import.sh

# LogicalResourceId=restoreTable      # This variable needs to be the same as the logical ID in the template. file://dynamoRollback.yml
# ChangeSetName=ImportChangeSet
# PROFILE=eu
# REGION=eu-west-1

changeSet=$(aws cloudformation create-change-set \
    --profile ${PROFILE} \
    --region ${REGION} \
    --stack-name ${StackName} \
    --change-set-name ${ChangeSetName} \
    --change-set-type IMPORT \
    --resources-to-import "[{\"ResourceType\":\"AWS::DynamoDB::Table\",\"LogicalResourceId\":\"${LogicalResourceId}\",\"ResourceIdentifier\": {\"TableName\":\"${TargetTableName}\"}}]" \
    --template-body file://3.dynamoImport.yml \
    --tags \
        Key=Department,Value=dev \
        Key=Environment,Value=int-xmn \
        Key=Team,Value=int-xmn \
        Key=Name,Value=Music-test-rollback \
    --parameters \
        ParameterKey=SourceTableName,UsePreviousValue=true \
        ParameterKey=TargetTableName,ParameterValue=${TargetTableName} \
    --query "Id")

# change-set-create-complete
aws cloudformation wait change-set-create-complete \
    --profile ${PROFILE} \
    --region ${REGION} \
    --change-set-name ${ChangeSetName} \
    --stack-name ${StackName}

# describe-change-set
# aws cloudformation describe-change-set \
#     --profile ${PROFILE} \
#     --change-set-name ${ChangeSetName} \
#     --stack-name ${StackName}

# execute-change-set
aws cloudformation execute-change-set \
    --profile ${PROFILE} \
    --change-set-name ${ChangeSetName} \
    --stack-name ${StackName}


