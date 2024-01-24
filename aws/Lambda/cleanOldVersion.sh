#!/usr/bin/env bash
set -e
REGION=us-east-1
PROFILE=int-xmn
FunctionName=karl_get_cost_test
UsedVersions=$(aws lambda list-aliases \
    --profile ${PROFILE} \
    --region ${REGION} \
    --function-name ${FunctionName} \
    --query Aliases[*].[FunctionVersion] \
    --output text)

AllVersions=$(aws lambda list-versions-by-function \
    --profile ${PROFILE} \
    --region ${REGION} \
    --function-name ${FunctionName} \
    --query Versions[*].[Version] \
    --output text)

for UsedVersion in ${UsedVersions[@]}; do 
    echo "Used version: ${FunctionName}:${UsedVersion}"
done

function checkVersion() {
    for version in ${AllVersions[@]}; do 
        find=0
        echo "checking version:${version}"
        if [ ${version} == '$LATEST' ];then
            echo "skip '\$LATEST' version"
            continue
        fi
        for used in ${UsedVersions[@]}; do 
            find=0
            if [ "${used}" == "${version}" ]; then
                echo "Version:${version} is being used by alias"
                find=1
                break
            fi   
        done
        if [ ${find} -eq 0 ]; then 
            echo "need to delete version:${version}"
            deleteVersion ${version}
            echo "delete ${FunctionName}:${version} successed"
            sleep 1s
        fi
    done
}

function deleteVersion() {
    TargetVersion=$1
    echo "deleting ${FunctionName}:${TargetVersion}..."
    result=$(aws lambda delete-function \
    --profile ${PROFILE} \
    --region ${REGION} \
    --function-name ${FunctionName} \
    --qualifier ${TargetVersion})
}

checkVersion

echo "Everything is done"

