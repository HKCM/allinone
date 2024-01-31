package main

import (
	"aws-env-monitor/common"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	retry             = common.GetIntEnv("Retry", 3)
	ChatWebhookURL    = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL    = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers      = common.GetEnv("MentionUsers", common.DefaultMentionUsers)
	ParameterStore    = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	ExcludeUsers      = strings.Split(common.GetEnv("ExcludeUsers", ""), ",")
	ExcludeEvents     = strings.Split(common.GetEnv("ExcludeEvents", ""), ",")
	ExcludeUserEvents = strings.Split(common.GetEnv("ExcludeUserEvents", ""), ",")
	MentionEvents     = strings.Split(common.GetEnv("MentionEvents", ""), ",")
	Repo              = common.GetEnv("Repo", "https://gitlab.example.com/devops/sysops/lambda/aws-env-monitor")
	LogLevel          = common.GetEnv("LogLevel", "Info")
	title             = "AWS Event Notice"
	IsWorkTime        = true
)

type Msg struct {
	RecipientAccountId string
	AWSRegion          string
	SourceIPAddress    string
	EventID            string
	EventName          string
	EventSource        string
	EventTime          string
	UserName           string
}

func LambdaHandler(ctx context.Context, event events.CloudwatchLogsEvent) error {
	log.Debug("ExcludeUsers: ", ExcludeUsers)
	log.Debug("ExcludeEvents: ", ExcludeEvents)
	log.Debug("ExcludeUserEvents: ", ExcludeUserEvents)

	data, _ := event.AWSLogs.Parse()

	// Write received message to log.
	for i := range data.LogEvents {
		message := data.LogEvents[i].Message
		log.Debug("message: ", message)

		data := make(map[string]interface{})

		err := json.Unmarshal([]byte(message), &data)
		if err != nil {
			log.Error("解析 JSON 出错:", err)
			return err
		}

		utcEventTime, _ := data["eventTime"].(string)
		eventTime := common.UTCToJSPTime(utcEventTime)

		// 如果现在是工作时间
		if common.IsWorkTime(eventTime) {
			IsWorkTime = true
			log.Infof("It's work time now: %s\n", eventTime)
			filterUserEvent(data, ExcludeUsers, ExcludeEvents, ExcludeUserEvents)
		}

		// 如果现在不是工作时间
		noExcludes := []string{} // 不排除任何东西
		if !common.IsWorkTime(eventTime) {
			IsWorkTime = false
			log.Infof("It's not work time now: %s, will send all events\n", eventTime)
			filterUserEvent(data, ExcludeUsers, noExcludes, noExcludes)
		}

	}
	return nil

}

func filterUserEvent(data map[string]interface{}, excludeUsers, excludeEvents, excludeUserEvents []string) {
	userName := ""
	eventSource, _ := data["eventSource"].(string)
	eventName, _ := data["eventName"].(string)
	awsRegion, _ := data["awsRegion"].(string)
	sourceIPAddress, _ := data["sourceIPAddress"].(string)
	recipientAccountId, _ := data["recipientAccountId"].(string)
	eventID, _ := data["eventID"].(string)
	userIdentity, _ := data["userIdentity"].(map[string]interface{})
	eventType, _ := userIdentity["type"].(string)
	principalId, _ := userIdentity["principalId"].(string)
	if eventType == "IAMUser" {
		userName, _ = userIdentity["userName"].(string)

		// 如果User:Event在ExcludeUserEvents列表则不发送消息
		for _, excludeUserEvent := range excludeUserEvents {
			if userName+":"+eventName == excludeUserEvent {
				log.Infof("User is %s, Event is %s in ExcludeUserEvent list, will not send message\n", userName, eventName)
				return
			}
		}

		// 如果用户名在ExcludeUsers列表则不发送消息
		for _, excludeUser := range excludeUsers {
			if userName == excludeUser {
				log.Infof("User is %s, in ExcludeUsers list, will not send message\n", userName)
				return
			}
		}

		// 如果Event在ExcludeEvents列表则不发送消息
		for _, excludeEvent := range ExcludeEvents {
			if eventName == excludeEvent {
				log.Infof("User is %s, Event is %s in ExcludeEvents list, will not send message\n", userName, eventName)
				return
			}
		}

	} else if eventType == "AssumedRole" && strings.Contains(principalId, "aws-cdk") {
		userName = strings.Split(principalId, ":")[1]
		// 如果Event在ExcludeEvents列表则不发送消息
		for _, excludeEvent := range ExcludeEvents {
			if eventName == excludeEvent {
				log.Infof("User is %s, Event is %s in ExcludeEvents list, will not send message\n", userName, eventName)
				return
			}
		}
	} else {
		// 如果不是用户调用直接返回
		log.Infof("Will not send message -- eventName: %s eventSource: %s principal: %s", eventName, eventSource, principalId)
		return
	}
	utcEventTime, _ := data["eventTime"].(string)
	eventTime := common.UTCToJSPTime(utcEventTime)

	msg := Msg{
		RecipientAccountId: recipientAccountId,
		AWSRegion:          awsRegion,
		SourceIPAddress:    sourceIPAddress,
		EventID:            eventID,
		EventName:          eventName,
		EventSource:        eventSource,
		EventTime:          eventTime,
		UserName:           userName,
	}

	buildMessage(msg)
}

func buildMessage(msg Msg) {

	detailUrl := "https://ap-northeast-1.console.aws.amazon.com/cloudtrail/home?region=ap-northeast-1#/events?EventId=" + msg.EventID

	msgs := fmt.Sprintf("AccountId:   %s\n", msg.RecipientAccountId)
	msgs += fmt.Sprintf("Region:      %s\n", msg.AWSRegion)
	msgs += fmt.Sprintf("SourceIP:    %s\n", msg.SourceIPAddress)
	msgs += fmt.Sprintf("EventName:   %s\n", msg.EventName)
	msgs += fmt.Sprintf("EventSource: %s\n", msg.EventSource)
	msgs += fmt.Sprintf("Principal:   %s\n", msg.UserName)
	msgs += fmt.Sprintf("EventID:     %s\n", msg.EventID)
	log.Infoln(msgs)

	var chatTargetUsers, larkTargetUsers map[string]string

	// 如果不在工作时间,所有事件都需要mention
	if IsWorkTime == false {
		chatTargetUsers, larkTargetUsers = common.GetMentionUsers(MentionUsers, ParameterStore)
	}

	// 如果在工作时间,且event实践中包含关键词,则mention
	if IsWorkTime == true {
		for _, mentionEvent := range MentionEvents {
			if strings.Contains(msg.EventName, mentionEvent) {
				log.Infof("Event %s is in MentionEvents list\n", msg.EventName)
				chatTargetUsers, larkTargetUsers = common.GetMentionUsers(MentionUsers, ParameterStore)
				break
			}
		}
	}
	buttons := []common.Button{}

	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Detail", detailUrl)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())

	cm := common.CardMsg{
		Title:     title,
		Text:      msgs,
		TimeStamp: msg.EventTime,
		Buttons:   buttons,
	}

	err := common.SendMessageToChatCard(ChatWebhookURL, cm, chatTargetUsers, retry)
	if err != nil {
		log.Panic(err)
	}
	err = common.SendMessageToLarkCard(LarkWebhookURL, cm, larkTargetUsers, retry)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	log.SetLevel(common.SetLogLevel(LogLevel))
	lambda.Start(LambdaHandler)
}
