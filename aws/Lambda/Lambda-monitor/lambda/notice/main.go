package main

import (
	"aws-env-monitor/common"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

var (
	Retry                            = common.GetIntEnv("Retry", 3)
	ChatWebhookURL                   = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL                   = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers                     = common.GetEnv("MentionUsers", "")
	Repo                             = common.GetEnv("Repo", "https://gitlab.example.com/devops/sysops/lambda/aws-env-monitor")
	ParameterStore                   = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	chatTargetUsers, larkTargetUsers = common.GetMentionUsers(MentionUsers, ParameterStore)
	LogLevel                         = common.GetEnv("LogLevel", "Info")
)

func LambdaHandler(ctx context.Context, event events.SNSEvent) {
	for _, record := range event.Records {
		cm := buildPayload(record.SNS)
		err := common.SendMessageToChatCard(ChatWebhookURL, cm, chatTargetUsers, Retry)
		if err != nil {
			log.Panic(err)
		}
		err = common.SendMessageToLarkCard(LarkWebhookURL, cm, larkTargetUsers, Retry)
		if err != nil {
			log.Panic(err)
		}
	}

}

func main() {
	log.SetLevel(common.SetLogLevel(LogLevel))
	lambda.Start(LambdaHandler)
}

func buildPayload(event events.SNSEntity) common.CardMsg {
	log.Infof("event.Subject: %v", event.Subject)
	log.Infof("event.Message: %v", event.Message)
	buttons := []common.Button{}

	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())

	timeStamp := common.UTCToJSPTime(event.Timestamp.Local().Format(time.RFC3339))
	var data map[string]interface{}
	err := json.Unmarshal([]byte(event.Message), &data)
	if err != nil {
		log.Infoln("Can't be Unmarshal, it just a message...")
		title := "Alarm From Message"
		msg := fmt.Sprintf("Message: %s\n", event.Message)
		if len(event.Subject) > 0 {
			msg = fmt.Sprintf("Subject: %s\n", event.Subject) + msg
		}

		cm := common.CardMsg{
			Title:     title,
			Text:      msg,
			TimeStamp: timeStamp,
			Buttons:   buttons,
		}

		return cm
	}

	accountId := data["AWSAccountId"].(string)
	region := data["Region"].(string)
	alertName := data["AlarmName"].(string)
	startsAt := data["StateChangeTime"].(string)
	oldState := data["OldStateValue"].(string)
	newState := data["NewStateValue"].(string)
	description := data["AlarmDescription"].(string)
	msg := fmt.Sprintf("AlarmName: %s\nDescription: %s\nAccountId: %s\nRegion: %s\nStateChangeTime: %s\nStateChange: %s", alertName, description, accountId, region, startsAt, oldState+" ---> "+newState)

	alarmUrl := "https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#alarmsV2:?~(search~'${alarmUrl})"
	alarmUrl = strings.Replace(alarmUrl, "${alarmUrl}", alertName, -1)
	buttons = common.AddButton(buttons, "Visit Alarm", alarmUrl)
	title := "Alarm From CloudWatch"
	color := "red"
	if newState == "OK" {
		color = "green"
	}

	cm := common.CardMsg{
		Title:     title,
		Color:     color,
		Text:      msg,
		TimeStamp: timeStamp,
		Buttons:   buttons,
	}

	return cm
}
