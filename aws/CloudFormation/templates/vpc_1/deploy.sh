#!/usr/bin/env bash

STACK_NAME=xxx-test
PROFILE=eu

aws cloudformation deploy \
    --profile ${PROFILE} \
    --template-file ./vpc_subnet.yaml \
    --stack-name ${STACK_NAME} \
    --tags \
        Name=${STACK_NAME} \
        Department=dev \
        Environment=xxxxx \
        Team=xxxxx
