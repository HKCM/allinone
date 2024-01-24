#!/bin/bash

STACK_NAME=user-test
awsProfile=int-xmn
Project=aha-add-in
region=us-east-1

aws cloudformation deploy \
    --profile ${awsProfile} \
    --region ${region} \
    --template-file da.yaml \
    --stack-name ${STACK_NAME} \
    --no-fail-on-empty-changeset \
    --tags \
        Department=dev \
        Team=int-xmn \
        Project=${Project} \
        Environment=int-xmn
