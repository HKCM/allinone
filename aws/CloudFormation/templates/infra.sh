#!/usr/bin/env bash

set -e

function usage() {
  echo "Usage:
TODO
"
  exit 0
}

while getopts "r:b:e:p:h" opt; do
  case "$opt" in
  r) AWSRegion="$OPTARG" ;;
  b) Brand="$OPTARG" ;;
  e) Stage="$OPTARG" ;;
  p) AWSProfile="$OPTARG" ;;
  h) usage ;;
  [?]) usage ;;
  esac
done


NameSpace=int-${Brand}-project-${Stage}


source ./env.sh

STACK_NAME=${NameSpace}-pre-infra

TEMPLATE_FILE=./infra.yml

aws cloudformation deploy \
    --profile "${AWSProfile}" \
    --region "${AWSRegion}" \
    --template-file "${TEMPLATE_FILE}" \
    --stack-name "${STACK_NAME}" \
    --capabilities CAPABILITY_IAM \
    --no-fail-on-empty-changeset \
    --tags \
        Name="${NameSpace}-pre-infra-stack" \
        Department=dev \
        Stage="${Stage}"\
        Team=int-xmn \
        Project="${NameSpace}" \
        NameSpace="${NameSpace}" \
        Environment=int-aws \
    --parameter-overrides \
        NameSpace="${NameSpace}" \
        S3LoggingBucket="${S3LoggingBucket}" \
        CloudFrontBucket="${CloudFrontBucket}" \
        ACMARN="${ACMARN}" \
        ALIASE="${ALIASE}" \
        APIGatewayOriginPath="${APIGatewayOriginPath}" \
        Endpoint="${Endpoint}"

SNS=$(aws cloudformation describe-stacks \
    --profile "${AWSProfile}" \
    --region "${AWSRegion}" \
    --stack-name "${STACK_NAME}" \
    --query "Stacks[0].NotificationARNs[0]")

if [ "${SNS}" = "null" ]; then
    aws cloudformation update-stack \
        --profile "${AWSProfile}" \
        --region "${AWSRegion}" \
        --stack-name "${STACK_NAME}" \
        --use-previous-template \
        --parameters \
            ParameterKey=CloudFrontBucket,UsePreviousValue=true \
            ParameterKey=CMR,UsePreviousValue=true \
            ParameterKey=S3LoggingBucket,UsePreviousValue=true \
            ParameterKey=NameSpace,UsePreviousValue=true \
            ParameterKey=ACMARN,UsePreviousValue=true \
            ParameterKey=ALIASE,UsePreviousValue=true \
            ParameterKey=APIGatewayOriginPath,UsePreviousValue=true \
            ParameterKey=Endpoint,UsePreviousValue=true \
        --notification-arns "arn:aws:sns:${AWSRegion}:${TargetAWSAccount}:${NameSpace}-cloudformation-events-subscription"
fi
