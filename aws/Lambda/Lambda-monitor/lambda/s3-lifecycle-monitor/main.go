package main

import (
	"aws-env-monitor/common"
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

var (
	Retry               = common.GetIntEnv("Retry", 3)
	SNSArn              = common.GetEnv("SNSArn", "")
	ChatWebhookURL      = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL      = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers        = common.GetEnv("MentionUsers", common.DefaultMentionUsers)
	ParameterStore      = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	ShowEnableLifecycle = common.GetBoolEnv("ShowEnableLifecycle", true)
	ShowNoLifecycle     = common.GetBoolEnv("ShowNoLifecycle", true)
	ExcludeBuckets      = common.GetEnv("ExcludeBuckets", "")
	LogLevel            = common.GetEnv("LogLevel", "Info")
	Repo                = common.GetEnv("Repo", "https://gitlab.example.com/devops/sysops/lambda/aws-env-monitor")
)

type Bucket struct {
	BucketName   string
	HasLifecycle bool
}

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) error {

	// create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	buckets := []Bucket{}

	resp, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalln("Failed to list buckets", err)
	}

	ExcludeBucketList := []string{}
	if ExcludeBuckets != "" {
		ExcludeBucketList = strings.Split(ExcludeBuckets, ",")
	}

	// Iterate through all the buckets to check if a special tag exists
	for _, bucket := range resp.Buckets {
		bucketName := bucket.Name
		if common.In(*bucketName, ExcludeBucketList) {
			continue
		}

		log.Infof("Checking bucket:\t%s\n", *bucketName)

		hasLifecycle := checkBucketSetting(s3Client, bucketName)
		buckets = append(buckets, Bucket{BucketName: *bucketName, HasLifecycle: hasLifecycle})
	}

	cm := buildPayloadCard(buckets)
	chatTargetUsers, larkTargetUsers := common.GetMentionUsers(MentionUsers, ParameterStore)
	err = common.SendMessageToChatCard(ChatWebhookURL, cm, chatTargetUsers, Retry)
	if err != nil {
		log.Panic(err)
	}

	err = common.SendMessageToLarkCard(LarkWebhookURL, cm, larkTargetUsers, Retry)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}

func checkBucketSetting(client *s3.Client, bucketName *string) bool {
	// maybe there is no tags on the bucket

	// get tags from the bucket

	_, err := client.GetBucketLifecycleConfiguration(context.TODO(),
		&s3.GetBucketLifecycleConfigurationInput{
			Bucket: bucketName,
		},
	)
	// If we can get the bucket Lifecycle Configuration
	// which means there is already have Lifecycle Configuration
	if err != nil {

		log.Debugf("Bucket: %s, Error %v", *bucketName, err)
		return false
	}

	log.Debugf("Bucket: %s has lifecycle config", *bucketName)
	return true
}

func buildPayloadCard(buckets []Bucket) common.CardMsg {
	msgWithLifecycle := ""
	msgWithoutLifecycle := ""
	for _, bucket := range buckets {
		if bucket.HasLifecycle {
			msgWithLifecycle = msgWithLifecycle + "- " + bucket.BucketName + "\n"
		}

		if !bucket.HasLifecycle {
			msgWithoutLifecycle = msgWithoutLifecycle + "- " + bucket.BucketName + "\n"
		}
	}
	msgWithLifecycle = msgWithLifecycle + "\n"
	msgWithoutLifecycle = msgWithoutLifecycle + "\n"

	msg := ""
	if ShowEnableLifecycle {
		msg = msg + "The following buckets have lifecycle:\n" + msgWithLifecycle

	}
	if ShowNoLifecycle {
		msg = msg + "The following buckets do not have lifecycle:\n" + msgWithoutLifecycle
	}

	buttons := []common.Button{}
	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())

	title := "Bucket Lifecycle Info"
	timeStamp := common.UTCToJSPTime(common.GetTimeStr())

	color := "blue"
	cm := common.CardMsg{
		Title:     title,
		Color:     color,
		Text:      msg,
		TimeStamp: timeStamp,
		Buttons:   buttons,
	}
	return cm
}
