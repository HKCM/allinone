package main

import (
	"aws-env-monitor/common"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

const (
	Prefix     string = "EC2"
	MaxRecords int32  = 100
	Maintainer string = "Karl.Huang"
	Repo       string = "devops/sysops/lambda/aws-env-monitor"
)

var (
	FunctionName               = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	Description                = "Created by lambda: " + FunctionName
	retry                      = common.GetIntEnv("Retry", 3)
	AlertManagerURL            = common.GetEnv("AlertManagerURL", "")
	xxxxxType                  = common.GetEnv("xxxxxType", "ops")
	AlertDisableAndEnableAlarm = common.GetBoolEnv("AlertDisableAndEnableAlarm", true)
	Severity                   = "critical"
	LogLevel                   = common.GetEnv("LogLevel", "Info")
	link                       = "https://ap-northeast-1.console.aws.amazon.com/ec2/home?region=ap-northeast-1#Instances:instanceId={INSTANCE_ID}"
)

func LambdaHandler(ctx context.Context, event events.SNSEvent) {

	fmt.Println(event)

	for _, record := range event.Records {
		message := buildAlertPayload(record.SNS)
		err := common.SendAlertToAlertManager(AlertManagerURL, message, 3)
		if err != nil {
			panic(fmt.Errorf("failed to send alert to alertmanager: %v", err))
		}
	}

}

func main() {
	log.SetLevel(common.SetLogLevel(LogLevel))
	lambda.Start(LambdaHandler)
}

func buildAlertPayload(event events.SNSEntity) common.AlertManagerMessage {

	var data map[string]interface{}
	err := json.Unmarshal([]byte(event.Message), &data)
	if err != nil {
		fmt.Println("Error:", err)
	}
	accountId := data["AWSAccountId"].(string)
	region := data["Region"].(string)
	alertNamePrefix := strings.Split(data["AlarmName"].(string), "_")[0]
	instanceName := strings.Split(data["AlarmName"].(string), "_")[1]
	instanceID := strings.Split(data["AlarmName"].(string), "_")[2]
	alertName := fmt.Sprintf("AWS %s Alarm", alertNamePrefix)
	startsAt := common.UTCToJSPTime(data["StateChangeTime"].(string))
	endsAt := ""
	source := instanceName
	ec2Link := strings.Replace(link, "{INSTANCE_ID}", instanceID, -1)
	oldState := data["OldStateValue"].(string)
	newState := data["NewStateValue"].(string)
	msg := fmt.Sprintf("AlarmName: %s\\nAccountId: %s\\nRegion: %s\\nInstance: %s\\nInstanceID: %s\\nStateChangeTime: %s\\nStateChange: %s", alertName, accountId, region, instanceName, instanceID, startsAt, oldState+" ---> "+newState)

	// 用于确认这是发出警报还是恢复
	if !strings.EqualFold(newState, "ok") {
		// 如果是发出警报,则把endsAt 设置为5小时后
		log.Info("Create alert")
		endsAt = common.TimeAdd(startsAt, 5*time.Hour)
	} else {
		// 如果是警报解除,则把endsAt 设置为startsAt
		log.Info("Resolved alert")
		endsAt = startsAt
	}

	alertData := common.AlertManagerMessage{
		AlertName:  alertName,
		Severity:   Severity,
		Message:    msg,
		ConfirmUrl: ec2Link,
		Source:     source,
		xxxxxType:  xxxxxType,
		StartsAt:   startsAt,
		EndAt:      endsAt,
	}
	log.Info(alertData)

	return alertData
}
