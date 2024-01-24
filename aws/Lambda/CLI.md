# 为lambda加上alb权限
```shell
aws lambda add-permission --function-name alb-function \
    --statement-id load-balancer --action "lambda:InvokeFunction" \
    --principal elasticloadbalancing.amazonaws.com
```

# 下载lambda layer
```shell
# https://docs.aws.amazon.com/cli/latest/reference/lambda/get-layer-version.html
URL=$(aws lambda get-layer-version --layer-name YOUR_LAYER_NAME_HERE --version-number YOUR_LAYERS_VERSION --query Content.Location --output text)
curl $URL -o layer.zip

# https://docs.aws.amazon.com/cli/latest/reference/lambda/get-layer-version-by-arn.html
URL=$(aws lambda get-layer-version-by-arn --arn arn:aws:lambda:us-east-1:209497400698:layer:php-73:7 --query Content.Location --output text)
curl $URL -o php.zip
```

# 查询特定runtime
```shell
aws lambda list functions--function version all --region us-east-1 --output text --query "functions[？Runtime=='python3.6'].FunctionArn"
```

# 删除多余的lambda version
```shell
function getVersion() {
    getFunctionName=$1
    UsedVersions=$(aws lambda list-aliases \
        --profile ${PROFILE} \
        --function-name ${getFunctionName} \
        --query Aliases[*].[FunctionVersion] \
        --output text)

    AllVersions=$(aws lambda list-versions-by-function \
        --profile ${PROFILE} \
        --function-name ${getFunctionName} \
        --query Versions[*].[Version] \
        --output text)
}

function checkVersion() {
    checkFunctionName=$1
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
            deleteVersion ${version} ${checkFunctionName}
            echo "delete ${checkFunctionName}:${version} succeeded"
            sleep 1s
        fi
    done
}

function deleteVersion() {
    TargetVersion=$1
    deleteFunctionName=$2
    echo "deleting ${deleteFunctionName}:${TargetVersion}..."
    result=$(aws lambda delete-function \
    --profile ${PROFILE} \
    --function-name ${deleteFunctionName} \
    --qualifier ${TargetVersion})
}
```

# 创建layer
```shell
pip3 install --target ./python ping3

zip -r ping3.zip python

aws lambda publish-layer-version --layer-name ping3 --description "My layer" --license-info "MIT" \
--zip-file  "fileb://ping3.zip"  --compatible-runtimes python3.6 python3.7   
```