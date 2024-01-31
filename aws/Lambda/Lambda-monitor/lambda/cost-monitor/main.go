package main

import (
	"aws-env-monitor/common"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type AnomalyData struct {
	AccountID          string `json:"accountId"`
	AccountName        string `json:"accountName"`
	AnomalyDetailsLink string `json:"anomalyDetailsLink"`
	AnomalyEndDate     string `json:"anomalyEndDate"`
	AnomalyID          string `json:"anomalyId"`
	AnomalyScore       struct {
		CurrentScore float64 `json:"currentScore"`
		MaxScore     float64 `json:"maxScore"`
	} `json:"anomalyScore"`
	AnomalyStartDate string `json:"anomalyStartDate"`
	DimensionalValue string `json:"dimensionalValue"`
	Impact           struct {
		MaxImpact             float64 `json:"maxImpact"`
		TotalActualSpend      float64 `json:"totalActualSpend"`
		TotalExpectedSpend    float64 `json:"totalExpectedSpend"`
		TotalImpact           float64 `json:"totalImpact"`
		TotalImpactPercentage float64 `json:"totalImpactPercentage"`
	} `json:"impact"`
	MonitorArn  string `json:"monitorArn"`
	MonitorName string `json:"monitorName"`
	MonitorType string `json:"monitorType"`
	RootCauses  []struct {
		LinkedAccount     string `json:"linkedAccount"`
		LinkedAccountName string `json:"linkedAccountName"`
		Region            string `json:"region"`
		Service           string `json:"service"`
		UsageType         string `json:"usageType"`
	} `json:"rootCauses"`
	SubscriptionID   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName"`
}

const (
	MAX_MESSAGE_SIZE = 4000
)

var (
	FunctionName     = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	retry            = common.GetIntEnv("Retry", 3)
	ChatWebhookURL   = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL   = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers     = common.GetEnv("MentionUsers", common.DefaultMentionUsers)
	ParameterStore   = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	MentionThreshold = common.GetIntEnv("MentionThreshold", 100)
	Repo             = common.GetEnv("Repo", "https://gitlab.example.com/devops/sysops/lambda/aws-env-monitor")

	chatTargetUsers map[string]string
	larkTargetUsers map[string]string
	title           = "AWS Cost Notice"
)

func LambdaHandler(ctx context.Context, snsEvent events.SNSEvent) error {

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		cardMessage := buildMessage(snsRecord.Message)
		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)

		err := common.SendMessageToChatCard(ChatWebhookURL, cardMessage, chatTargetUsers, retry)
		if err != nil {
			panic(err)
		}
		err = common.SendMessageToLarkCard(LarkWebhookURL, cardMessage, larkTargetUsers, retry)
		if err != nil {
			panic(err)
		}
	}

	return nil

}

func main() {
	lambda.Start(LambdaHandler)
}

func buildMessage(message string) common.CardMsg {

	var data AnomalyData
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		log.Error("Error:", err)
		return common.CardMsg{
			Title:     title,
			TimeStamp: common.GetTimeStr(),
		}
	}

	log.Infof("AnomalyData: %+v", data)
	originAccountName := data.RootCauses[0].LinkedAccountName
	originAccountId := data.RootCauses[0].LinkedAccount
	rootAccount := fmt.Sprintf("%s(%s)", originAccountName, originAccountId)
	accountInfo := fmt.Sprintf("AccountId:\t%s\n", data.AccountID)
	serviceInfo := fmt.Sprintf("Service:\t%s\n", data.RootCauses[0].Service)
	ActualSpend := fmt.Sprintf("ActualSpend:\t$%.2f\n", data.Impact.TotalActualSpend)
	ExpectSpend := fmt.Sprintf("ExpectSpend:\t$%.2f\n", data.Impact.TotalExpectedSpend)
	TotalImpact := fmt.Sprintf("TotalImpact:\t$%.2f\n", data.Impact.TotalImpact)
	Percentage := fmt.Sprintf("Percentage:\t%.2f%%\n", data.Impact.TotalImpactPercentage)
	OrgAccount := fmt.Sprintf("OrgAccountId:\t%s\nOrgRegion:\t%s\n", rootAccount, data.RootCauses[0].Region)

	msgs := accountInfo + serviceInfo + ActualSpend + ExpectSpend + TotalImpact + Percentage
	if originAccountName == "" {
		msgs += OrgAccount
	}

	// 如果费用高于阈值才查询MentionUsers
	if data.Impact.TotalImpact > float64(MentionThreshold) {
		chatTargetUsers, larkTargetUsers = common.GetMentionUsers(MentionUsers, ParameterStore)
		title = "AWS Cost Alarm"
	}

	buttons := []common.Button{}

	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Detail", data.AnomalyDetailsLink)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())

	cm := common.CardMsg{
		Title:     title,
		Text:      msgs,
		TimeStamp: data.AnomalyEndDate,
		Buttons:   buttons,
	}

	return cm
}
