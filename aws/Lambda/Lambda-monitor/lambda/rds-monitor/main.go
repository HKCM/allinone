package main

import (
	"aws-env-monitor/common"
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudWatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	log "github.com/sirupsen/logrus"
)

const (
	Prefix     string = "RDS-"
	MaxRecords int32  = 100
	Maintainer string = "Karl.Huang"
)

const (
	WriterCPUUtilization   string = "Writer-CPUUtilization"
	ReaderCPUUtilization   string = "Reader-CPUUtilization"
	ReadLatency            string = "ReadLatency"
	WriteLatency           string = "WriteLatency"
	LoginFailures          string = "LoginFailures"
	AuroraBinlogReplicaLag string = "AuroraBinlogReplicaLag"
)

type Cluster struct {
	ClusterName   string
	ClusterEngine string
	ClusterStatus string
}

var (
	FunctionName   = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	Description    = "This Alarm is created automatically by lambda: " + FunctionName + ". please do not change this alert. any question please contact the maintainer"
	Retry          = common.GetIntEnv("Retry", 3)
	SNSArn         = common.GetEnv("SNSArn", "")
	ChatWebhookURL = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers   = common.GetEnv("MentionUsers", common.DefaultMentionUsers)
	ParameterStore = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	LogLevel       = common.GetEnv("LogLevel", "Info")
	blueStage      = "old1"
	Repo           = common.GetEnv("Repo", "https://gitlab.example.com/devops/sysops/lambda/aws-env-monitor")

	AlertTags = []cloudWatchTypes.Tag{
		{
			Key:   aws.String("Maintainer"),
			Value: aws.String(Maintainer),
		},
		{
			Key:   aws.String("Repo"),
			Value: aws.String(Repo),
		},
		{
			Key:   aws.String("Description"),
			Value: aws.String(Description),
		},
	}
	MySQLAlarm = []string{
		WriterCPUUtilization,
		ReaderCPUUtilization,
		ReadLatency,
		WriteLatency,
		LoginFailures,
		AuroraBinlogReplicaLag,
	}
	PostgresSQLAlarm = []string{
		WriterCPUUtilization,
		ReaderCPUUtilization,
		ReadLatency,
		WriteLatency,
	}
)

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) {

	// create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	log.Infoln("Schedule Event: Check all running RDS cluster...")
	clusters := getAllRDSClusters(cfg)
	log.Infoln(clusters)
	rdsAlarms := getAllRDSAlarms(cfg)
	checkAlarmsForClusters(rdsAlarms, clusters)

}

func checkAlarmsForClusters(rdsAlarms []string, clusters []Cluster) {
	alarmList := []string{}
	title := "Created RDS Alarm"
	buttons := []common.Button{}
	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())

	chatTargetUsers, larkTargetUsers := common.GetMentionUsers(MentionUsers, ParameterStore)
	for _, cluster := range clusters {
		createdAlarms := ""
		timeStamp := common.UTCToJSPTime(common.GetTimeStr())
		// 根据Engine确认警报的类型
		if cluster.ClusterEngine == "aurora-mysql" {
			alarmList = MySQLAlarm
		}

		if cluster.ClusterEngine == "aurora-postgresql" {
			alarmList = PostgresSQLAlarm
		}

		// 挨个检查每一个警报
		for _, suffix := range alarmList {

			alarmName := buildAlarmName(cluster, suffix)
			alarmAlreadyExists := common.In(alarmName, rdsAlarms)
			// 如果发现对应的警报已存在
			if alarmAlreadyExists {
				log.Debugf("Alarm: %s already existed, skip...\n", alarmName)
			}

			// 如果发现没有对应的警报则创建
			if !alarmAlreadyExists {
				log.Infof("Alarm: %s does not exist, create...\n", alarmName)
				createAlarm(alarmName, suffix, cluster)
				createdAlarms = createdAlarms + "\n- " + suffix
			}
		}

		// 创建了警报才发送通知
		if len(createdAlarms) > 0 {
			msg := "Created Alarms for " + cluster.ClusterName
			msg = msg + createdAlarms

			color := "blue"
			cm := common.CardMsg{
				Title:     title,
				Color:     color,
				Text:      msg,
				TimeStamp: timeStamp,
				Buttons:   buttons,
			}
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
}

func createAlarm(alarmName, suffix string, cluster Cluster) {
	// create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	buttons := []common.Button{}
	buttons = common.AddButton(buttons, "Visit GitRepo", Repo)
	buttons = common.AddButton(buttons, "Visit Lambda", common.GetLambdaUrl())
	timeStamp := common.UTCToJSPTime(common.GetTimeStr())
	var chatTargetUsers, larkTargetUsers map[string]string
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	var input cloudwatch.PutMetricAlarmInput

	if suffix == WriterCPUUtilization {
		log.Infof("Try to create %s alarm for %s...", suffix, cluster.ClusterName)
		input = cloudwatch.PutMetricAlarmInput{
			Namespace: aws.String("AWS/RDS"),
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBClusterIdentifier"),
					Value: aws.String(cluster.ClusterName),
				},
				{
					Name:  aws.String("Role"),
					Value: aws.String("WRITER"),
				},
			},
			AlarmName:               aws.String(alarmName),
			MetricName:              aws.String("CPUUtilization"),
			Statistic:               cloudWatchTypes.StatisticAverage,
			Period:                  aws.Int32(60),
			EvaluationPeriods:       aws.Int32(3),
			Threshold:               aws.Float64(80),
			ComparisonOperator:      cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
			TreatMissingData:        aws.String("missing"),
			OKActions:               []string{SNSArn},
			InsufficientDataActions: []string{SNSArn},
			AlarmActions:            []string{SNSArn},
			Tags:                    AlertTags,
			AlarmDescription:        aws.String("Monitoring WRITER CPUUtilization is GreaterThanThreshold ( 80 percents )"),
		}
	}

	if suffix == ReaderCPUUtilization {
		log.Infof("Try to create %s alarm for %s...", suffix, cluster.ClusterName)
		input = cloudwatch.PutMetricAlarmInput{
			Namespace: aws.String("AWS/RDS"),
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBClusterIdentifier"),
					Value: aws.String(cluster.ClusterName),
				},
				{
					Name:  aws.String("Role"),
					Value: aws.String("READER"),
				},
			},
			AlarmName:               aws.String(alarmName),
			MetricName:              aws.String("CPUUtilization"),
			Statistic:               cloudWatchTypes.StatisticAverage,
			Period:                  aws.Int32(60),
			EvaluationPeriods:       aws.Int32(3),
			Threshold:               aws.Float64(80),
			ComparisonOperator:      cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
			TreatMissingData:        aws.String("missing"),
			OKActions:               []string{SNSArn},
			InsufficientDataActions: []string{SNSArn},
			AlarmActions:            []string{SNSArn},
			Tags:                    AlertTags,
			AlarmDescription:        aws.String("Monitoring READER CPUUtilization is GreaterThanThreshold ( 80 percents )"),
		}
	}

	if suffix == ReadLatency {
		log.Infof("Try to create %s alarm for %s...", suffix, cluster.ClusterName)
		input = cloudwatch.PutMetricAlarmInput{
			Namespace: aws.String("AWS/RDS"),
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBClusterIdentifier"),
					Value: aws.String(cluster.ClusterName),
				},
			},
			AlarmName:               aws.String(alarmName),
			MetricName:              aws.String("ReadLatency"),
			Statistic:               cloudWatchTypes.StatisticAverage,
			Period:                  aws.Int32(60),
			EvaluationPeriods:       aws.Int32(3),
			Threshold:               aws.Float64(0.1),
			ComparisonOperator:      cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
			TreatMissingData:        aws.String("missing"),
			OKActions:               []string{SNSArn},
			InsufficientDataActions: []string{SNSArn},
			AlarmActions:            []string{SNSArn},
			Tags:                    AlertTags,
			AlarmDescription:        aws.String("Monitoring ReadLatency is GreaterThanThreshold ( 0.1 seconds )"),
		}
	}

	if suffix == WriteLatency {
		log.Infof("Try to create %s alarm for %s...", suffix, cluster.ClusterName)
		input = cloudwatch.PutMetricAlarmInput{
			Namespace: aws.String("AWS/RDS"),
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBClusterIdentifier"),
					Value: aws.String(cluster.ClusterName),
				},
			},
			AlarmName:               aws.String(alarmName),
			MetricName:              aws.String("WriteLatency"),
			Statistic:               cloudWatchTypes.StatisticAverage,
			Period:                  aws.Int32(60),
			EvaluationPeriods:       aws.Int32(3),
			Threshold:               aws.Float64(0.1),
			ComparisonOperator:      cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
			TreatMissingData:        aws.String("missing"),
			OKActions:               []string{SNSArn},
			InsufficientDataActions: []string{SNSArn},
			AlarmActions:            []string{SNSArn},
			Tags:                    AlertTags,
			AlarmDescription:        aws.String("Monitoring WriteLatency is GreaterThanThreshold ( 0.1 seconds )"),
		}
	}

	if suffix == LoginFailures {
		log.Infof("Try to create %s alarm for %s...", suffix, cluster.ClusterName)
		input = cloudwatch.PutMetricAlarmInput{
			Namespace: aws.String("AWS/RDS"),
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBClusterIdentifier"),
					Value: aws.String(cluster.ClusterName),
				},
			},
			AlarmName:               aws.String(alarmName),
			MetricName:              aws.String("LoginFailures"),
			Statistic:               cloudWatchTypes.StatisticMaximum,
			Period:                  aws.Int32(60),
			EvaluationPeriods:       aws.Int32(1),
			Threshold:               aws.Float64(3),
			ComparisonOperator:      cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
			TreatMissingData:        aws.String("missing"),
			OKActions:               []string{SNSArn},
			InsufficientDataActions: []string{SNSArn},
			AlarmActions:            []string{SNSArn},
			Tags:                    AlertTags,
			AlarmDescription:        aws.String("Monitoring DB LoginFailures is GreaterThanThreshold ( 3 times )"),
		}
	}

	if suffix == AuroraBinlogReplicaLag {
		log.Infof("Try to create %s alarm for %s...", suffix, cluster.ClusterName)
		input = cloudwatch.PutMetricAlarmInput{
			Namespace: aws.String("AWS/RDS"),
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBClusterIdentifier"),
					Value: aws.String(cluster.ClusterName),
				},
			},
			AlarmName:               aws.String(alarmName),
			MetricName:              aws.String("AuroraBinlogReplicaLag"),
			Statistic:               cloudWatchTypes.StatisticAverage,
			Period:                  aws.Int32(60),
			EvaluationPeriods:       aws.Int32(3),
			Threshold:               aws.Float64(2),
			ComparisonOperator:      cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
			TreatMissingData:        aws.String("missing"),
			OKActions:               []string{SNSArn},
			InsufficientDataActions: []string{SNSArn},
			AlarmActions:            []string{SNSArn},
			Tags:                    AlertTags,
			AlarmDescription:        aws.String("Monitoring DB AuroraBinlogReplicaLag is GreaterThanThreshold ( 2 seconds )"),
		}
	}

	_, err = cloudwatchClient.PutMetricAlarm(context.TODO(), &input)
	if err != nil {
		title := "Create RDS Alarm Failed"
		color := "blue"
		msg := "Created " + suffix + " Alarm for " + cluster.ClusterName + " Failed... Please check"
		cm := common.CardMsg{
			Title:     title,
			Color:     color,
			Text:      msg,
			TimeStamp: timeStamp,
			Buttons:   buttons,
		}

		err = common.SendMessageToChatCard(ChatWebhookURL, cm, chatTargetUsers, Retry)
		if err != nil {
			log.Panic(err)
		}
		err = common.SendMessageToLarkCard(LarkWebhookURL, cm, larkTargetUsers, Retry)
		if err != nil {
			log.Panic(err)
		}
		log.Fatalf("Created System check Alarm failed %v\n", err)
	}

}

func buildAlarmName(cluster Cluster, suffix string) string {
	return "RDS-" + cluster.ClusterEngine + "-" + cluster.ClusterName + "-" + suffix
}

// Get all rds cluster that status=available
func getAllRDSClusters(cfg aws.Config) []Cluster {

	auroraClusters := []Cluster{}
	rdsClient := rds.NewFromConfig(cfg)

	// get all RDS cluster
	auroraClustersResult, err := rdsClient.DescribeDBClusters(context.TODO(), &rds.DescribeDBClustersInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, DBCluster := range auroraClustersResult.DBClusters {

		clusterName := *DBCluster.DBClusterIdentifier
		clusterEngine := *DBCluster.Engine
		clusterStatus := *DBCluster.Status
		cluster := Cluster{
			ClusterName:   clusterName,
			ClusterEngine: clusterEngine,
			ClusterStatus: clusterStatus,
		}
		log.Infof("%s: %s is %s", clusterEngine, clusterName, clusterStatus)
		if cluster.ClusterStatus == "available" && !strings.Contains(clusterName, blueStage) {
			auroraClusters = append(auroraClusters, cluster)
			log.Debugf("add %s:%s to enable alarm list", cluster.ClusterEngine, cluster.ClusterName)
		}
	}
	return auroraClusters
}

// 获取所有的RDS开头的Alarms
func getAllRDSAlarms(cfg aws.Config) []string {
	var nextToken *string
	rdsAlarms := []string{}
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	for {
		result, err := cloudwatchClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{
			AlarmNamePrefix: aws.String(Prefix),
			MaxRecords:      aws.Int32(MaxRecords),
			NextToken:       nextToken,
		})

		if err != nil {
			log.Fatalln("DescribeAlarms Failed: ", err)
		}

		for _, alarm := range result.MetricAlarms {
			rdsAlarms = append(rdsAlarms, *alarm.AlarmName)
		}
		nextToken = result.NextToken
		if nextToken == nil {
			break
		}
	}
	log.Debugln("Total RDS alarms:", rdsAlarms)
	return rdsAlarms
}

func main() {
	log.SetLevel(common.SetLogLevel(LogLevel))
	lambda.Start(LambdaHandler)
}
