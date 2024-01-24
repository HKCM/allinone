#!/usr/bin/env bash

# Run this script in /DynamoDB folder

# bash ./Template/deploy.sh

STACK_NAME=db-test
PROFILE=eu

CUR_DIR=$(cd `dirname $0` && pwd)

aws cloudformation deploy \
    --profile ${PROFILE} \
    --template-file ${CUR_DIR}/dynamo.yml \
    --stack-name ${STACK_NAME} \
    --tags \
        Department=dev \
        Environment=int \
        Team=int \
    --parameter-overrides \
        NameSpace="${STACK_NAME}"