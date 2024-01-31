package main

import (
	"aws-env-monitor/common"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/computeoptimizer"
	"github.com/aws/aws-sdk-go-v2/service/computeoptimizer/types"
	log "github.com/sirupsen/logrus"
)

const (
	MaxRecords int32  = 100
	Maintainer string = "Karl.Huang"
)

var (
	retry            = common.GetIntEnv("Retry", 3)
	ChatWebhookURL   = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL   = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers     = common.GetEnv("MentionUsers", "")
	ParameterStore   = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	ExcludeInstances = common.GetEnv("ExcludeInstances", "")
	Repo             = common.GetEnv("Repo", "https://gitlab.example.com/devops/sysops/lambda/aws-env-monitor")
	LogLevel         = common.GetEnv("LogLevel", "Info")
	Notice           = common.GetBoolEnv("Notice", false)
)

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) {

	var nextToken *string
	var chatTargetUsers, larkTargetUsers map[string]string
	title := "EC2 Optimizer Notice"
	detailURL := "https://us-east-1.console.aws.amazon.com/compute-optimizer/home?region=ap-northeast-1#/resources-lists/ec2"
	if Notice {
		chatTargetUsers, larkTargetUsers = common.GetMentionUsers(MentionUsers, ParameterStore)
	}

	msg := ""
	excludeList := strings.Split(ExcludeInstances, ",")
	number := 0
	accountID := ""

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	client := computeoptimizer.NewFromConfig(cfg)

	for {
		// 获取所有的EC2开头的Alarms
		RecommendationsOutput, err := client.GetEC2InstanceRecommendations(context.TODO(), &computeoptimizer.GetEC2InstanceRecommendationsInput{
			Filters: []types.Filter{
				{
					Name: types.FilterNameFinding,
					Values: []string{
						string(types.FindingOverProvisioned),
						string(types.FindingUnderProvisioned),
					},
				},
			},
			MaxResults: aws.Int32(MaxRecords),
			NextToken:  nextToken,
		})
		if err != nil {
			log.Fatalf("failed to GetEC2InstanceRecommendations: %v", err)
		}
		// fmt.Printf("%+v", RecommendationsOutput.InstanceRecommendations)

		for _, instanceRecommendation := range RecommendationsOutput.InstanceRecommendations {
			isExclude := false
			instanceName := *instanceRecommendation.InstanceName
			// 如果机器在exclude list中,则排除机器
			for _, excludeInstanceName := range excludeList {
				if instanceName == excludeInstanceName {
					isExclude = true
					break
				}
			}
			if isExclude {
				continue
			}

			instanceArn := *instanceRecommendation.InstanceArn
			instanceID := strings.Split(instanceArn, "/")[1]
			log.Infof("%s(%s):%s\n", instanceName, instanceID, instanceRecommendation.InstanceState)
			number += 1
			msg += fmt.Sprintf("%s(%s)", instanceName, instanceRecommendation.InstanceState) + " is " + string(instanceRecommendation.Finding) + "\n"
			if accountID == "" {
				accountID = *instanceRecommendation.AccountId
			}
		}

		nextToken = RecommendationsOutput.NextToken
		if nextToken == nil {
			break
		}

	}
	msg = fmt.Sprintf("AWS Account: %s\n\n", accountID) + msg
	msg += fmt.Sprintf("\nTotal of %d instances can be optimized", number)
	log.Infoln(msg)

	buttons := []common.Button{}

	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Detail", detailURL)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())

	cm := common.CardMsg{
		Title:     title,
		Text:      msg,
		TimeStamp: common.GetTimeStr(),
		Buttons:   buttons,
	}

	err = common.SendMessageToChatCard(ChatWebhookURL, cm, chatTargetUsers, retry)
	if err != nil {
		log.Fatal(err)
	}
	err = common.SendMessageToLarkCard(LarkWebhookURL, cm, larkTargetUsers, retry)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	log.SetLevel(common.SetLogLevel(LogLevel))
	lambda.Start(LambdaHandler)
}
