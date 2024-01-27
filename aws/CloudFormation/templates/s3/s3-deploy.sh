#!/bin/bash

set -e
#cleanup (Profile, Stack-name, Region)
# It can clean the Failed Change set
# Solve: An error occurred (LimitExceededException) when calling the CreateChangeSet operation: ChangeSet limit exceeded for stack ...

stackname=karl-test
temfile=./s3.yaml

cleanup () {
    local profile=$1
    local stackname=$2
    local region=$3
    if [[ -z $region ]];then
      region=us-east-1
    fi
    echo "Cleaning up failed change sets"
    changesets=$(aws cloudformation list-change-sets \
        --profile $profile \
        --region $region \
        --stack-name $stackname --query 'Summaries[?Status==`FAILED`].ChangeSetId' --output text)
    echo changesets
    for changeset in $changesets; do
      echo "${stackname}: deleting change set ${changeset}"
      aws cloudformation delete-change-set \
        --profile $profile \
        --region $region \
        --stack-name $stackname \
        --change-set-name ${changeset}
    done
}

aws cloudformation deploy \
    --profile int-xmn \
    --region us-east-1 \
    --template-file ${temfile} \
    --stack-name ${stackname} \
    --no-fail-on-empty-changeset \
    --tags \
        Department=dev \
        Environment=dev \
        Team=int-xmn \
    --parameter-overrides \
        Name=${stackname}

cleanup int-xmn ${stackname}

