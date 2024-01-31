package main

import (
	"aws-env-monitor/common"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/jsii-runtime-go"
)

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) error {

	eventByte, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Failed to marshal detail map:", err)
		return err
	}

	logEvent := fmt.Sprintf("CloudWatchEvent: %s", eventByte)
	fmt.Println(logEvent)

	// fmt.Print(message)

	payload := buildPayload("AWS Health Alarm", logEvent)
	NoticeLambdaArn := common.GetEnv("common.LambdaEnvName_NoticeLambdaArn", "common.DefaultNoticeLambdaArn")

	input := &awslambda.InvokeInput{
		FunctionName: jsii.String(NoticeLambdaArn),
		Payload:      payload,
	}

	// create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %v", err))
	}

	// create Lambda client to invoke notice lambda
	lambdaclient := awslambda.NewFromConfig(cfg)

	fmt.Println("Invoking notice lambda function: ", NoticeLambdaArn)

	result, err := lambdaclient.Invoke(context.Background(), input)
	if err != nil {
		panic(fmt.Errorf("failed to invoke Lambda function: %v", err))
	}
	resPayload, err := common.PrintLambdaInvokeResult(result)
	if err != nil {
		return err
	}
	fmt.Println("Notice lambda function: ", NoticeLambdaArn, " return result: ", resPayload)

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}

func buildPayload(subject, message string) []byte {

	// webhookURL := common.GetEnv(common.LambdaEnvName_WebhookURL, common.DefaultWebhookURL)
	// // 从环境变量获取parameterStoreKey
	// parameterStoreKey := common.GetEnv(common.LambdaEnvName_ParameterStore, common.DefaultParameterStore)

	// mentionUsersString := common.GetEnv(common.LambdaEnvName_MentionUsers, common.DefaultMentionUsers)
	// // for now we don't nse lark
	// // TODO
	// // Add lark invoke function
	// chatMentionUsers, _, err := common.GetMentionUsers(mentionUsersString, parameterStoreKey)

	invokeEvent := common.NoticeLambdaMessageEvent{
		// WebhookURL:   webhookURL,
		Title:   subject,
		Message: message,
		// MentionUsers: chatMentionUsers,
	}

	payload, err := json.Marshal(invokeEvent)
	if err != nil {
		panic(fmt.Errorf("failed to marshal event: %v", err))
	}

	return payload
}
