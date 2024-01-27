#!/usr/bin/env bash

set -e

function usage() {
  echo "Usage:

  ./scripts/admin/deploy_iam_policy.sh \\
    -s <service> \\
    -b <brand> \\
    -e <environment/stage> \\
    -p <aws_profile> \\

Example:

  ./scripts/admin/deploy_iam_policy.sh -s project -b new -e dev -p int-developer
"
  exit 0
}

while getopts "s:b:e:p:t:h" opt; do
  case "$opt" in
  s) Service="$OPTARG" ;;
  b) Brand="$OPTARG" ;;
  e) Stage="$OPTARG" ;;
  p) AwsProfile="$OPTARG" ;;
  h) usage ;;
  [?]) usage ;;
  esac
done

echo "========== [${Service}] IAM DEPLOYMENT START =========="

ProjectPrefix=int-${Brand}
NameSpace=${ProjectPrefix}-${Service}-${Stage}

STACK_NAME=${NameSpace}-iam-policy

TEMPLATE_FILE=services/${Service}/infra/iam_policy.yml

#REGION=$(aws configure get region --profile "$AwsProfile")

aws cloudformation deploy \
    --profile "${AwsProfile}" \
    --region "${REGION}" \
    --template-file "${TEMPLATE_FILE}" \
    --stack-name "${STACK_NAME}" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-fail-on-empty-changeset \
    --tags \
        Name="${NameSpace}-pre-infra-stack" \
        Department=dev \
        Stage=Prod \
        Team=int-xmn \
        Project="${NameSpace}" \
        NameSpace="${NameSpace}" \
        Environment="${Stage}" \
    --parameter-overrides \
        ProjectPrefix="${ProjectPrefix}" \
        Stage="${Stage}" \
        NameSpace="${NameSpace}" \
    --notification-arns "arn:aws:sns:${REGION}:${TargetAWSAccount}:${NameSpace}-cloudformation-events-subscription"

echo "========== [${Service}] IAM DEPLOYMENT COMPLETED =========="
