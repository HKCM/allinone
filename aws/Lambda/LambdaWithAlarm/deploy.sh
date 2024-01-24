#!/usr/bin/env bash

BASEPATH=$(cd `dirname $0`; pwd)
LambdaFunctionName=temporary-test
Word=World
STACK_NAME=temporary-test
SNSTopic=topic
PROFILE=eu
TEMPLATE_FILE=lambda.yml
ORIGIN_CODE=lambda.py
ZIP_CODE=lambda.zip


zip ./${ZIP_CODE} ./${ORIGIN_CODE}



aws cloudformation deploy \
    --profile ${PROFILE} \
    --template-file ${BASEPATH}/${TEMPLATE_FILE} \
    --capabilities CAPABILITY_NAMED_IAM \
    --stack-name ${STACK_NAME} \
    --tags \
        Name=${STACK_NAME} \
        Department=dev \
        Environment=xxxx \
        Team=xxxx \
    --parameter-overrides \
        Word=${Word} \
        LambdaFunctionName=${LambdaFunctionName} \
        SNSTopic=${SNSTopic}


aws lambda update-function-code \
    --profile ${PROFILE} \
    --function-name ${LambdaFunctionName} \
    --zip-file fileb://./${ZIP_CODE}
