#!/usr/bin/env bash

# This script will deploy CloudFormation. It will create two dynamodb tables

# bash ./deploy.sh
source ./0.env.sh

aws cloudformation deploy \
    --profile ${PROFILE} \
    --region ${REGION} \
    --template-file ./1.dynamo.yml \
    --stack-name ${StackName} \
    --tags \
        Department=dev \
        Environment=int-xmn \
        Team=int-xmn \
        Name=Music-test-rollback \
    --parameter-overrides \
        SourceTableName=${SourceTableName}