#!/usr/bin/env bash
status=$(aws lambda invoke \
    --profile ${AWSProfile} \
    --region us-east-1 \
    --function-name $LambdaFunctionName \
    --cli-binary-format raw-in-base64-out \
    --payload "{ \"url\": \"${url}\" }" \
    --output text \
    --query 'FunctionError' \
    temp_result.json )

if [[ $status == "None" ]]; then
    echo -e "Good connection with $env environment\n"
else
    echo -e "Error: No connection with $env environment\n"
fi


echo '{"key": "test"}' > clear_payload  
openssl base64 -out encoded_payload -in clear_payload
aws lambda invoke --function-name testsms  --invocation-type Event --payload file://~/encoded_paylaod response.json