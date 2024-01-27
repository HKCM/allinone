#!/usr/bin/env bash

# stack_name=user-test
# profile=st2
# region=ap-northeast-1

aws cloudformation deploy \
    --profile ${profile} \
    --region ${region} \
    --template-file vpc.yml \
    --stack-name ${stack_name} \
    --no-fail-on-empty-changeset