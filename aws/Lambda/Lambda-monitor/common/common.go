package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-sdk-go-v2/config"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	log "github.com/sirupsen/logrus"
)

type LarkAction struct {
	Tag      string   `json:"tag,omitempty"`
	Text     Text     `json:"text,omitempty"`
	MultiUrl MultiURL `json:"multi_url,omitempty"`
}

type Text struct {
	Content string `json:"content,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type MultiURL struct {
	Url string `json:"url,omitempty"`
}

type Cards struct {
	Cards []Card `json:"cards,omitempty"`
}

type Card struct {
	Header   *Header   `json:"header,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

type Header struct {
	Title      string `json:"title,omitempty"`
	Subtitle   string `json:"subtitle,omitempty"`
	ImageURL   string `json:"imageUrl,omitempty"`
	ImageStyle string `json:"imageStyle,omitempty"`
}

type Section struct {
	Header  string   `json:"header,omitempty"`
	Widgets []Widget `json:"widgets,omitempty"`
}

type Widget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
	KeyValue      *KeyValue      `json:"keyValue,omitempty"`
	Image         *Image         `json:"image,omitempty"`
	Buttons       []Button       `json:"buttons,omitempty"`
}

type TextParagraph struct {
	Text string `json:"text,omitempty"`
}

type KeyValue struct {
	TopLabel         string   `json:"topLabel,omitempty"`
	Content          string   `json:"content,omitempty"`
	Icon             string   `json:"icon,omitempty"`
	ContentMultiLine string   `json:"contentMultiline,omitempty"`
	BottomLabel      string   `json:"bottomLabel,omitempty"`
	OnClick          *OnClick `json:"onClick,omitempty"`
	Button           *Button  `json:"button,omitempty"`
}

type Image struct {
	ImageURL string   `json:"imageUrl,omitempty"`
	OnClick  *OnClick `json:"onClick,omitempty"`
}

type Button struct {
	TextButton  *TextButton  `json:"textButton,omitempty"`
	ImageButton *ImageButton `json:"imageButton,omitempty"`
}

type TextButton struct {
	Text    string   `json:"text,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

type ImageButton struct {
	IconURL string   `json:"iconUrl,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

type OnClick struct {
	OpenLink *OpenLink `json:"openLink,omitempty"`
}

type OpenLink struct {
	URL string `json:"url,omitempty"`
}

type Data struct {
	Chat map[string]string `json:"chat"`
	Lark map[string]string `json:"lark"`
}

// type for Parameter Store
type ParameterData struct {
	Chat map[string]string `json:"chat"`
	Lark map[string]string `json:"lark"`
}

type NoticeLambdaMessageEvent struct {
	WebhookURLs  map[string]string `json:"webhookURLs"`
	Title        string            `json:"title"`
	Message      string            `json:"message"`
	MentionUsers string            `json:"mentionUsers"`
	MessageFrom  string            `json:"messageFrom"`
	Retry        int               `json:"retry"`    // 大于等于0 小于20 默认为0
	In           int               `json:"in"`       // 降噪参数,相同时间内的同样的消息作为一条,单位为秒(未启用)
	Duration     int               `json:"duration"` // 告警的持续时间,单位为秒(未启用)
	Interval     int               `json:"interval"` // 告警的间隔,单位为秒(未启用)
}

// https://xxxxx-office.jp.larksuite.com/wiki/Q0nwweHoRiGgpLknOOJjBLmBpte
type AlertManagerMessage struct {
	AlertName  string
	Severity   string // "warning", "critical"
	Message    string
	ConfirmUrl string
	Source     string
	xxxxxType  string // "dev", "ops", "other"
	StartsAt   string // "2023-10-13T06:00:00Z"
	EndAt      string // "2023-10-13T06:00:00Z"
}

type CardMsg struct {
	Title     string
	Color     string
	Text      string
	Buttons   []Button
	TimeStamp string
}

type NoticeType = string

const (
	Lark NoticeType = "Lark"
	Chat NoticeType = "Chat"
)

const (
	ExportNoticeLambdaArn = "NoticeLambdaArn" // Name of Stack output
	DefaultMentionUsers   = "karl.huang"
	DefaultSNSArn         = "arn:aws:sns:ap-northeast-1:012345678901:cloudwatch_alarm_notice_topic_with_lambda"
	DefaultParameterStore = "xxxxx-mention-user-ids"
	// stg-monitor channel
	DefaultChatWebhookURL = "https://chat.googleapis.com/v1/spaces/AAAApTWlmQo/messages?key=xxxxx&token=xxx"
	// TODO lark的URL
	DefaultLarkWebhookURL = "https://open.larksuite.com/open-apis/bot/v2/hook/xxxxx"
	MAX_MESSAGE_SIZE      = 4000 // 单条消息最大字符数
)

var (
	FunctionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	AWS_REGION   = os.Getenv("AWS_REGION")
)

func GetContextStr(app awscdk.App, ctxStr, defaultValue string) string {
	ctxStr, ok := app.Node().TryGetContext(&ctxStr).(string)
	if !ok {
		ctxStr = defaultValue
	}
	return ctxStr
}

func GetLambdaUrl() string {
	url := fmt.Sprintf("https://%s.console.aws.amazon.com/lambda/home?region=%s#/functions/%s", AWS_REGION, AWS_REGION, FunctionName)
	return url
}

func GetEnv(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if len(value) < 1 {
		return defaultValue
	}
	return value
}

func GetIntEnv(envName string, defaultValue int) int {
	value := os.Getenv(envName)
	if len(value) < 1 {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetBoolEnv(envName string, defaultValue bool) bool {
	strValue := os.Getenv(envName)

	if len(strValue) == 0 || strValue == "" {
		return defaultValue
	}

	upperValue := strings.ToUpper(strValue)

	falseList := []string{"NO", "FALSE"}
	for _, v := range falseList {
		if upperValue == v {
			return false
		}
	}

	return true
}

// return chatTargetUsers,larkTargetUsers
func GetMentionUsers(mentionUserString, parameterName string) (map[string]string, map[string]string) {
	if mentionUserString == "" {
		return nil, nil
	}

	mentionUserStringLower := strings.ToLower(mentionUserString)
	chatTargetUsers := make(map[string]string)
	larkTargetUsers := make(map[string]string)
	var targetUsers []string
	if strings.Contains(mentionUserString, ",") {
		targetUsers = strings.Split(mentionUserStringLower, ",")
	} else if strings.Contains(mentionUserString, "|") {
		targetUsers = strings.Split(mentionUserStringLower, "|")
	} else {
		targetUsers = []string{mentionUserStringLower}
	}

	allUserIds, err := getMentionUsersFromParameterStore(parameterName)
	log.Info(allUserIds)
	if err != nil {
		chatTargetUsers["Warning"] = fmt.Sprintf("Warning: Get MentionUsers failed, Target user is %s", targetUsers)
		larkTargetUsers["Warning"] = fmt.Sprintf("Warning: Get MentionUsers failed, Target user is %s", targetUsers)
		return chatTargetUsers, larkTargetUsers
	}

	var data Data

	err = json.Unmarshal([]byte(allUserIds), &data)
	if err != nil {
		chatTargetUsers["Warning"] = fmt.Sprintf("Warning: allUserIds 解析 JSON 出错, Target user is %s", targetUsers)
		larkTargetUsers["Warning"] = fmt.Sprintf("Warning: allUserIds 解析 JSON 出错, Target user is %s", targetUsers)
		return chatTargetUsers, larkTargetUsers
	}

	// Get ID from chat filed
	chatData := data.Chat
	larkData := data.Lark
	log.Info(chatData)
	log.Info(larkData)
	for _, user := range targetUsers {
		// 构建chat mention user
		chatTargetUsers[user] = chatData[user]

		// 构建lark mention user
		larkTargetUsers[user] = larkData[user]

	}
	return chatTargetUsers, larkTargetUsers
}

func getMentionUsersFromParameterStore(parameterName string) (string, error) {

	// create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %v", err))
	}

	// 构建 GetParameter 输入参数
	input := &ssm.GetParameterInput{
		Name: &parameterName,
	}

	// create Parameter Store client to get mention users
	ssmclient := ssm.NewFromConfig(cfg)
	resp, err := ssmclient.GetParameter(context.TODO(), input)
	if err != nil {
		log.Error("无法获取参数:", err)
		return "", err
	}
	value := *resp.Parameter.Value
	return value, nil
}

// 如果无法被Json解析则返回原始字符
func ToFormattedJsonString(message string) string {
	var data interface{}
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		log.Error("Error:", err)
		return message
	}

	formattedString, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error("Error:", err)
		return message
	}
	return string(formattedString)
}

func SendMessageToLark(title, message, webhookURL string, larkMentionUsers map[string]string, retry int) error {
	log.Infof("mentionUsers: %s\n message: %s\n", larkMentionUsers, message)
	mentionList := ""

	// 组装被mention的用户
	for key, value := range larkMentionUsers {
		if key == "Warning" {
			mentionList += "\n" + value + " <at user_id=\"xxxxx0082\"></at>" // karl.huang
			log.Info("--------------------", value, "--------------")
			break
		} else {
			mentionList += fmt.Sprintf(" <at user_id=\"%s\"></at>", value)
		}
	}

	text := ""
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		// If the message reaches its maximum size, it will be split and sent in parts.
		if len(text)+len(line)+len(FunctionName) > MAX_MESSAGE_SIZE-len(mentionList) {
			payload, err := buildPayload(title, text, mentionList, Lark)
			if err != nil {
				return err
			}
			postData(webhookURL, payload, retry)
			text = "" + line + "\n"
		} else {
			text = text + line + "\n"
		}
	}

	payload, err := buildPayload(title, text, mentionList, Lark)
	if err != nil {
		return err
	}

	// POST payload to Lark Webhook URL
	err = postData(webhookURL, payload, retry)
	if err != nil {
		return err
	}

	return nil
}

func SendAlertToAlertManager(alertManagerURL string, message AlertManagerMessage, retry int) error {

	// 构建Alertmanager告警数据，以JSON格式
	alertData := `[
		{
		  "labels": {
			"alertname": "{ALERTNAME}",
			"severity": "{SEVERITY}",
			"source": "{SOURCE}",
			"xxxxxType": "{xxxxxTYPE}"
		  },
		  "annotations": {
			"message": "{MESSAGE}",
			"confirmUrl": "{CONFIRMURL}"
		  },
		  "startsAt": "{STARTSAT}",
		  "endsAt": "{ENDSAT}"
		}
	  ]`

	alertData = strings.ReplaceAll(alertData, "{ALERTNAME}", message.AlertName)
	alertData = strings.ReplaceAll(alertData, "{SEVERITY}", message.Severity)
	alertData = strings.ReplaceAll(alertData, "{SOURCE}", message.Source)
	alertData = strings.ReplaceAll(alertData, "{xxxxxTYPE}", message.xxxxxType)
	alertData = strings.ReplaceAll(alertData, "{MESSAGE}", message.Message)
	alertData = strings.ReplaceAll(alertData, "{CONFIRMURL}", message.ConfirmUrl)
	alertData = strings.ReplaceAll(alertData, "{STARTSAT}", message.StartsAt)
	alertData = strings.ReplaceAll(alertData, "{ENDSAT}", message.EndAt)

	log.Infoln(alertData)

	err := postData(alertManagerURL, []byte(alertData), retry)
	if err != nil {
		return err
	}

	return nil
}

func SendMessageToChat(title, message, webhookURL string, chatMentionUsers map[string]string, retry int) error {
	log.Infof("mentionUsers: %s\n message: %s\n", chatMentionUsers, message)
	mentionList := ""

	// 组装被mention的用户
	for key, value := range chatMentionUsers {
		if key == "Warning" {
			mentionList += "\n" + value + " <users/113061970078918879376>" // karl.huang
			fmt.Println("--------------------", value, "--------------")
			break
		} else {
			//fmt.Printf("NoticeUser: %s, UserID: %s\n", key, value)
			mentionList += fmt.Sprintf(" <users/%s>", value)
		}
	}

	text := ""
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		// 如果超过最大字符,将会分段发送
		if len(text)+len(line)+len(FunctionName) > MAX_MESSAGE_SIZE-len(mentionList) {
			payload, err := buildPayload(title, text, mentionList, Chat)
			if err != nil {
				return err
			}
			postData(webhookURL, payload, retry)
			text = "" + line + "\n"
		} else {
			text = text + line + "\n"
		}
	}

	payload, err := buildPayload(title, text, mentionList, Chat)
	if err != nil {
		return err
	}

	// POST payload to Google Chat Webhook URL
	err = postData(webhookURL, payload, retry)
	if err != nil {
		return err
	}

	return nil
}

func buildPayload(title, message, mentionList string, noticeType string) ([]byte, error) {
	var msg map[string]interface{}

	if noticeType == Lark {
		text := ""
		if title == "" {
			text = mentionList
		} else {
			text = title + "\n" + message + "\n" + mentionList
		}

		msg = map[string]interface{}{
			"msg_type": "text",
			"content": map[string]string{
				"text": text,
			},
		}
	}

	if noticeType == Chat {
		msg = map[string]interface{}{
			"text": title + "\n```" + message + "```\n" + mentionList,
		}
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Failed to serialize message:", err)
		return nil, err
	}
	return payload, nil
}

func buildPayloadCard(cardMsg CardMsg, mentionList string, noticeType string) ([]byte, error) {
	var msg map[string]interface{}

	if noticeType == Lark {
		log.Debugln("Start build lark card")
		actions := []LarkAction{}
		for _, button := range cardMsg.Buttons {
			actions = append(actions, LarkAction{
				Tag: "button",
				Text: Text{
					Content: button.TextButton.Text,
					Tag:     "plain_text",
				},
				MultiUrl: MultiURL{
					Url: button.TextButton.OnClick.OpenLink.URL,
				},
			},
			)
		}
		card_color := "red"
		if cardMsg.Color != "" {
			card_color = cardMsg.Color
		}

		msg = map[string]interface{}{
			"msg_type": "interactive",
			"card": map[string]interface{}{
				"header": map[string]interface{}{
					"template": card_color,
					"title": map[string]interface{}{
						"content": cardMsg.Title,
						"tag":     "plain_text",
					},
				},
				"elements": []map[string]interface{}{
					{
						"tag": "column_set",
						"columns": []map[string]interface{}{
							{
								"tag":            "column",
								"width":          "weighted",
								"weight":         1,
								"vertical_align": "top",
								"elements": []map[string]interface{}{
									{
										"tag": "div",
										"text": map[string]interface{}{
											"content": "**🕐 Time:**\n" + cardMsg.TimeStamp,
											"tag":     "lark_md",
										},
									},
								},
							},
						},
					},
					{
						"tag": "div",
						"text": map[string]interface{}{
							"content": cardMsg.Text,
							"tag":     "plain_text",
						},
					},
					{
						"tag": "div",
						"text": map[string]interface{}{
							"content": mentionList,
							"tag":     "lark_md",
						},
					},
					{
						"tag":     "action",
						"actions": actions,
					},
				},
			},
		}
	}

	if noticeType == Chat {
		log.Debugln("Start build chat card")
		msg = map[string]interface{}{
			"cards": []Card{{
				Header: &Header{
					Title:    cardMsg.Title,
					Subtitle: cardMsg.TimeStamp,
				},
				Sections: []Section{
					{
						Widgets: []Widget{
							{
								TextParagraph: &TextParagraph{
									Text: cardMsg.Text,
								},
							},
							{
								Buttons: cardMsg.Buttons,
							},
						},
					},
				},
			}},
			"text": mentionList,
		}
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Failed to serialize message:", err)
		return nil, err
	}
	return payload, nil
}

func SendMessageToChatCard(webhookURL string, cardMsg CardMsg, chatMentionUsers map[string]string, retry int) error {
	log.Debugf("cardMsg: %v\n", cardMsg)
	mentionList := ""

	// 组装被mention的用户
	for key, value := range chatMentionUsers {
		if key == "Warning" {
			mentionList += "\n" + value + " <users/113061970078918879376>" // karl.huang
			fmt.Println("--------------------", value, "--------------")
			break
		} else {
			//fmt.Printf("NoticeUser: %s, UserID: %s\n", key, value)
			mentionList += fmt.Sprintf(" <users/%s>", value)
		}
	}

	text := ""
	lines := strings.Split(cardMsg.Text, "\n")
	for _, line := range lines {
		// 如果超过最大字符,将会分段发送
		if len(text)+len(line)+len(FunctionName) > MAX_MESSAGE_SIZE-len(mentionList) {
			payload, err := buildPayloadCard(cardMsg, mentionList, Chat)
			if err != nil {
				return err
			}
			postData(webhookURL, payload, retry)
			text = "" + line + "\n"
		} else {
			text = text + line + "\n"
		}
	}

	payload, err := buildPayloadCard(cardMsg, mentionList, Chat)
	if err != nil {
		return err
	}

	// POST payload to Google Chat Webhook URL
	err = postData(webhookURL, payload, retry)
	if err != nil {
		return err
	}

	return nil
}

func SendMessageToLarkCard(webhookURL string, cardMsg CardMsg, larkMentionUsers map[string]string, retry int) error {
	log.Debugf("mentionUsers: %v\n", cardMsg)
	mentionList := ""

	// 组装被mention的用户
	for key, value := range larkMentionUsers {
		if key == "Warning" {
			mentionList += "\n" + value + " <at id=\"xxxxx0082\"></at>" // karl.huang
			fmt.Println("--------------------", value, "--------------")
			break
		} else {
			// fmt.Printf("NoticeUser: %s, UserID: %s\n", key, value)
			mentionList += fmt.Sprintf(" <at id=\"%s\"></at>", value)
		}
	}

	text := ""
	lines := strings.Split(cardMsg.Text, "\n")
	for _, line := range lines {
		// If the message reaches its maximum size, it will be split and sent in parts.
		if len(text)+len(line)+len(FunctionName) > MAX_MESSAGE_SIZE-len(mentionList) {
			payload, err := buildPayloadCard(cardMsg, mentionList, Lark)
			if err != nil {
				return err
			}
			postData(webhookURL, payload, retry)
			text = "" + line + "\n"
		} else {
			text = text + line + "\n"
		}
	}

	payload, err := buildPayloadCard(cardMsg, mentionList, Lark)
	if err != nil {
		return err
	}

	// POST payload to Lark Webhook URL
	err = postData(webhookURL, payload, retry)
	if err != nil {
		return err
	}

	return nil
}

func postData(url string, payload []byte, retry int) error {
	msg := ""

	// 限制retry数量
	if retry < 0 {
		retry = 0
	}
	if retry > 20 {
		retry = 20
	}

	for i := 0; i < retry+1; i++ {
		retryInterval := time.Duration(i*i) * time.Second
		time.Sleep(retryInterval)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		// post执行失败
		if err != nil {
			log.Error("Failed to send message:", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Infoln("Message sent successfully")
			return nil
		} else {
			// 发送失败
			msg = fmt.Sprintf("Failed to send message. Status: %s", resp.Status)
			log.Errorf(msg+". Retry time: %d, will send message again after %s\n", i, retryInterval.String())
		}
	}
	return errors.New(fmt.Sprint(msg))
}

func PrintLambdaInvokeResult(result *awslambda.InvokeOutput) (string, error) {
	// 获取响应数据
	payload := string(result.Payload)

	// 获取状态码
	statusCode := result.StatusCode
	if statusCode != 200 {
		// 获取错误类型（如果有）
		functionError := *result.FunctionError
		fmt.Println("Function Error:", functionError)
		return payload, errors.New(fmt.Sprint("Function Error:", functionError))
	}
	return payload, nil
}

func TimeAdd(utcTimeString string, u time.Duration) string {
	// 使用指定的时间格式来解析UTC时间字符串
	layout := "2006-01-02T15:04:05.000-0700"
	utcTime, err := time.Parse(layout, utcTimeString)
	if err != nil {
		fmt.Println("Failed to parse UTC time:", err)
		fmt.Println("Using UTC time zone")
		utcTime = time.Now().UTC()
	}

	return utcTime.Add(u).Format(layout)
}

// get utc timeStr 2006-01-02T15:04:05Z
func GetTimeStr() string {
	// 使用指定的时间格式来解析UTC时间字符串
	layout := "2006-01-02T15:04:05Z"

	utcTime := time.Now().UTC()

	return utcTime.Format(layout)
}

func UTCToJSPTime(utcTimeString string) string {
	// 解析UTC时间字符串
	utcTime, err := time.Parse(time.RFC3339, utcTimeString)
	if err != nil {
		fmt.Println("解析时间出错:", err)
		return utcTimeString
	}

	// 获取日本时区
	japanLocation, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("加载时区出错:", err)
		return utcTimeString
	}

	// 转换为日本时间
	japanTime := utcTime.In(japanLocation)

	// 输出结果
	log.Debugln("UTC时间:", utcTime.Format(time.RFC3339))
	return japanTime.Format(time.RFC3339)
}

func AddButton(buttons []Button, text, url string) []Button {
	buttons = append(buttons, Button{
		TextButton: &TextButton{
			Text: text,
			OnClick: &OnClick{
				OpenLink: &OpenLink{
					URL: url,
				},
			},
		},
	})

	return buttons
}

func SetLogLevel(logLevel string) log.Level {
	// 设置默认日志级别
	log.SetLevel(log.InfoLevel)

	// 从环境变量获取日志级别，如果未设置则使用默认值
	envLogLevel := strings.ToLower(logLevel)
	switch envLogLevel {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn", "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		log.Warnf("Invalid Log Level '%s', Set level as 'info'", envLogLevel)
		return log.InfoLevel
	}
}

func IsWorkTime(inputTimeStr string) bool {
	// 解析时间字符串
	inputTime, err := time.Parse(time.RFC3339, inputTimeStr)
	if err != nil {
		log.Errorf("时间解析失败: %v", err)
		return false
	}

	// 获取星期几和小时
	weekday := inputTime.Weekday()
	hour := inputTime.Hour()
	log.Debugf("Event Time Week:%d Hour:%d\n", weekday, hour)

	// 判断是否在周一至周五的09:00-19:00范围内
	if weekday >= time.Monday && weekday <= time.Friday && hour >= 9 && hour < 19 {
		return true
	}

	return false
}

func ExcludeString(arr []string, elem string) []string {
	i := 0

	for _, str := range arr {
		if str != elem {
			arr[i] = str
			i++
		}
	}

	return arr[:i]
}

func In(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
