#!/usr/bin/env bash

case ${Stage} in
test )
    ALIASE="project.example.com"
    ACMARN=${RND_ACM_ARN}
    APIGatewayOriginPath='/test'
    CloudFrontBucket=${RND_Cloudfront_Bucket}
    ;;
staging )
    ALIASE="project.stage.example.com"
    ACMARN="arn:aws:acm:us-east-1:123456789:certificate/792f32d9-5b44-4977-9a2f-e9df6f4fe1cd"
    APIGatewayOriginPath='/staging'
    CloudFrontBucket='stage-logs.s3.amazonaws.com'
    ;;
production )
    ALIASE="project.example.com"
    ACMARN="arn:aws:acm:us-east-1:1234567890:certificate/33320b9c-61c5-45ea-b128-af948e367ef0"
    APIGatewayOriginPath='/production'
    CloudFrontBucket='prod-logs.s3.amazonaws.com'
    ;;
* )
    ALIASE="${NameSpace}.example-bot.com"
    ACMARN=${RND_ACM_ARN}
    APIGatewayOriginPath="/${Stage}"
    CloudFrontBucket=${RND_Cloudfront_Bucket}
    ;;
esac

export ACMARN
export ALIASE
export APIGatewayOriginPath
export CloudFrontBucket