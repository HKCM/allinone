package main

import (
	"aws-env-monitor/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudWatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type InstanceInfo struct {
	InstanceId string `json:"instanceid"`
	State      string `json:"state"`
	Name       string `json:"name"`
}

type AlarmInfo struct {
	AlarmArn     string `json:"arn"`
	AlarmName    string `json:"name"`
	AlarmEnabled bool   `json:"enabled"`
}

type Detail struct {
	InstanceId string `json:"instance-id"`
	State      string `json:"state"`
}

const (
	Prefix     string = "EC2"
	MaxRecords int32  = 100
	Maintainer string = "Karl.Huang"
	Repo       string = "devops/sysops/lambda/aws-env-monitor"
)

var (
	FunctionName                     = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	Description                      = "This Alarm is created automatically by lambda: " + FunctionName + ". please do not change this alert. any question please contact the maintainer"
	retry                            = common.GetIntEnv("Retry", 3)
	OKSNSArn                         = common.GetEnv("OKSNSArn", "")
	InsufficientSNSArn               = common.GetEnv("InsufficientSNSArn", "")
	AlarmSNSArn                      = common.GetEnv("AlarmSNSArn", "")
	ChatWebhookURL                   = common.GetEnv("ChatWebhookURL", common.DefaultChatWebhookURL)
	LarkWebhookURL                   = common.GetEnv("LarkWebhookURL", common.DefaultLarkWebhookURL)
	MentionUsers                     = common.GetEnv("MentionUsers", common.DefaultMentionUsers)
	ParameterStore                   = common.GetEnv("ParameterStore", common.DefaultParameterStore)
	AlertDisableAndEnableAlarm       = common.GetBoolEnv("AlertDisableAndEnableAlarm", true)
	chatTargetUsers, larkTargetUsers = common.GetMentionUsers(MentionUsers, ParameterStore)
	InstanceAlarmNameFormat          = common.GetEnv("InstanceAlarmNameFormat", "")
	SystemAlarmNameFormat            = common.GetEnv("SystemAlarmNameFormat", "")
	TriggerTag                       = common.GetEnv("TriggerTag", "xxxxx-monitor=on")
	// STG环境不需要OKAction
	OKActions = []string{}

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
)

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) {
	// 新的alarm名不能为空 直接终止运行
	if SystemAlarmNameFormat == "" || InstanceAlarmNameFormat == "" {
		panic(errors.New("Alarm Name can't be empty, please check"))
	}

	fmt.Println(event)
	// create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %v", err))
	}

	// 事件源于定时任务
	// 定时任务会扫描所有带有xxxxx-monitor=on的机器
	if event.DetailType == "Scheduled Event" {
		fmt.Println("Schedule Event: Check all " + TriggerTag + " and running instance...")
		err := checkAllInstances(cfg)
		if err != nil {
			panic(fmt.Errorf("checkAllInstance failed: %v\n", err))
		}
	}

	// 事件源于EC2事件
	if event.DetailType != "Scheduled Event" {
		fmt.Println("EC2 State Change Event: Check single instance...")
		err = checkSingleInstance(cfg, event.Detail)
		if err != nil {
			panic(fmt.Errorf("checkSingleInstance failed: %v\n", err))
		}
	}
}

func setAlarm(cfg aws.Config, instances []InstanceInfo) {
	for _, instance := range instances {
		// 机器被stopped， disable alarm
		if instance.State == "stopping" {
			fmt.Printf("EC2 %s State is stopping, do nothing...\n", instance.Name)
			// fmt.Printf("EC2 %s State is stopping, so disable alarm...\n", instance.Name)
			// err := disableAlarmForInstance(cfg, instance)
			// if err != nil {
			// 	panic(fmt.Errorf("Disable Alarm failed %v\n", err))
			// }
		}

		// 机器被删除，删除alarm
		if instance.State == "terminated" {
			fmt.Printf("EC2 %s State is terminated, so delete alarm...\n", instance.Name)
			err := deleteAlarmForInstance(cfg, instance)
			if err != nil {
				panic(fmt.Errorf("Delete Alarm failed %v", err))
			}
		}

		// 机器处于pending状态,有可能是新建机器,有可能是启动现有机器
		if instance.State == "pending" {
			// 检查机器是否已经存在alarm
			fmt.Printf("EC2 %s State is pending, try to get alarms...\n", instance.Name)
			alarms, err := getAlarmsForInstance(cfg, instance)
			if err != nil {
				panic(fmt.Errorf("Get Alarm failed %v", err))
			}

			// 如果alarm不存在，则说明机器是新启动的，创建alarm
			if len(alarms) == 0 {
				fmt.Printf("There is no exist alarm with %s, creating new alarms...\n", instance.Name)
				err = createInstanceCheckAlarmForInstance(cfg, instance)
				if err != nil {
					panic(fmt.Errorf("Create Alarm failed %v", err))
				}

				err = createSystemCheckAlarmForInstance(cfg, instance)
				if err != nil {
					panic(fmt.Errorf("Create Alarm failed %v", err))
				}
			}
			// 如果alarm已经存在，则enable alert
			if len(alarms) == 2 {
				fmt.Printf("There are already exist alarms with %s, do nothing...\n", instance.Name)
				// fmt.Printf("There are already exist alarms with %s, enable alarms...\n", instance.Name)
				// err = enableAlarmForInstance(cfg, instance)
				// if err != nil {
				// 	panic(fmt.Errorf("Enable Alarm failed %v\n", err))
				// }
			}
		}
	}
}

func checkAllInstances(cfg aws.Config) error {
	var nextToken *string
	instances := make([]InstanceInfo, 0, 150)
	ec2StatsAlarms := make([]AlarmInfo, 0, MaxRecords)
	instances, err := getAllInstances(cfg, instances)
	if err != nil {
		return err
	}
	fmt.Printf("There are %d running instances with "+TriggerTag+" tag\n", len(instances))

	for {
		// 获取所有的EC2开头的Alarms
		nextToken, ec2StatsAlarms, err = getAllAlarms(cfg, nextToken, ec2StatsAlarms)
		if err != nil {
			return err
		}
		if nextToken == nil {
			break
		}
	}

	fmt.Printf("There are %d Alarms' name start with EC2\n", len(ec2StatsAlarms))
	okInstance := 0
	for _, instance := range instances {
		instanceCheckAlarmName := buildAlarmName(InstanceAlarmNameFormat, instance.Name, instance.InstanceId)
		hasInstanceCheckAlarm := false

		systemCheckAlarmName := buildAlarmName(SystemAlarmNameFormat, instance.Name, instance.InstanceId)
		hasSystemCheckAlarm := false

		for _, alarm := range ec2StatsAlarms {
			if alarm.AlarmName == instanceCheckAlarmName {
				fmt.Printf("Found alarm with instance state check, instance: %s\n", instance.InstanceId)
				hasInstanceCheckAlarm = true
			}
			if alarm.AlarmName == systemCheckAlarmName {
				fmt.Printf("Found alarm with system state check, instance: %s\n", instance.InstanceId)
				hasSystemCheckAlarm = true
			}
			// 如果同时找到了两个alarm 说明这个instance的alarm是正常的
			if hasSystemCheckAlarm && hasInstanceCheckAlarm {
				fmt.Printf("Check alarm OK, instance: %s\n", instance.InstanceId)
				okInstance += 1
				fmt.Printf("OK instance: %d\n", okInstance)
				break
			}
		}

		// 如果没找到则创建
		if !hasInstanceCheckAlarm {
			fmt.Printf("Not found alarm with instance state check, instance: %s, creating...\n", instance.InstanceId)
			err := createInstanceCheckAlarmForInstance(cfg, instance)
			if err != nil {
				return err
			}
		}

		// 如果没找到则创建
		if !hasSystemCheckAlarm {
			fmt.Printf("Not found alarm with system state check, instance: %s, creating...\n", instance.InstanceId)
			err := createSystemCheckAlarmForInstance(cfg, instance)
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("OK instance: %d\n", okInstance)
	return nil
}

func checkSingleInstance(cfg aws.Config, detail json.RawMessage) error {
	var instanceDetail Detail
	json.Unmarshal(detail, &instanceDetail)
	instances, err := getSingleInstance(cfg, instanceDetail.InstanceId)
	if err != nil {
		return err
	}
	if len(instances) == 0 {
		fmt.Printf("Instance: %s has no tag with "+TriggerTag+", do nothing\n", instanceDetail.InstanceId)
		return nil
	}

	instances[0].State = instanceDetail.State
	fmt.Printf("Instance: %s state changed, start setup alarm...\n", instanceDetail.InstanceId)

	setAlarm(cfg, instances)
	return nil
}

func deleteAlarmForInstance(cfg aws.Config, instance InstanceInfo) error {
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	_, err := cloudwatchClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []string{
			buildAlarmName(InstanceAlarmNameFormat, instance.Name, instance.InstanceId),
			buildAlarmName(SystemAlarmNameFormat, instance.Name, instance.InstanceId),
		},
	})

	if err != nil {
		return err
	}

	title := "Deleted Alarms"
	message := "Deleted Instance check and System check Alarms for " + instance.Name
	fmt.Printf("Deleted Instance check and System check Alarms for %s\n", instance.Name)
	common.SendMessageToChat(title, message, ChatWebhookURL, chatTargetUsers, retry)
	common.SendMessageToLark(title, message, LarkWebhookURL, larkTargetUsers, retry)

	return nil
}

func createInstanceCheckAlarmForInstance(cfg aws.Config, instance InstanceInfo) error {

	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	if OKSNSArn != "" {
		OKActions = []string{OKSNSArn}
	}
	// 创建Instance状态检查
	_, err := cloudwatchClient.PutMetricAlarm(context.TODO(), &cloudwatch.PutMetricAlarmInput{
		Namespace: aws.String("AWS/EC2"),
		Dimensions: []cloudWatchTypes.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instance.InstanceId),
			},
		},
		AlarmName:          aws.String(buildAlarmName(InstanceAlarmNameFormat, instance.Name, instance.InstanceId)),
		MetricName:         aws.String("StatusCheckFailed_Instance"),
		Statistic:          cloudWatchTypes.StatisticMaximum,
		Period:             aws.Int32(60),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
		TreatMissingData:   aws.String("missing"),
		OKActions:          OKActions,
		InsufficientDataActions: []string{
			InsufficientSNSArn,
		},
		AlarmActions: []string{
			AlarmSNSArn,
		},
		Tags: AlertTags,
	})

	if err != nil {
		return err
	}
	title := "Created Alarm"
	message := "Created Instance check Alarm for " + instance.Name
	fmt.Printf("Created Instance check Alarm for %s\n", instance.Name)
	common.SendMessageToChat(title, message, ChatWebhookURL, chatTargetUsers, retry)
	common.SendMessageToLark(title, message, LarkWebhookURL, larkTargetUsers, retry)
	return nil
}

func createSystemCheckAlarmForInstance(cfg aws.Config, instance InstanceInfo) error {
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	if OKSNSArn != "" {
		OKActions = []string{OKSNSArn}
	}
	// 创建System状态检查
	_, err := cloudwatchClient.PutMetricAlarm(context.TODO(), &cloudwatch.PutMetricAlarmInput{
		Namespace: aws.String("AWS/EC2"),
		Dimensions: []cloudWatchTypes.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instance.InstanceId),
			},
		},
		AlarmName:          aws.String(buildAlarmName(SystemAlarmNameFormat, instance.Name, instance.InstanceId)),
		MetricName:         aws.String("StatusCheckFailed_System"),
		Statistic:          cloudWatchTypes.StatisticMaximum,
		Period:             aws.Int32(300),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudWatchTypes.ComparisonOperatorGreaterThanThreshold,
		TreatMissingData:   aws.String("missing"),
		OKActions:          OKActions,
		InsufficientDataActions: []string{
			InsufficientSNSArn,
		},
		AlarmActions: []string{
			AlarmSNSArn,
		},
		Tags: AlertTags,
	})
	if err != nil {
		fmt.Printf("Created System check Alarm failed %v\n", err)
		return err
	}

	title := "Created Alarm"
	message := "Created System check Alarm for " + instance.Name
	fmt.Printf("Created System check Alarm for %s\n", instance.Name)
	common.SendMessageToChat(title, message, ChatWebhookURL, chatTargetUsers, retry)
	common.SendMessageToLark(title, message, LarkWebhookURL, larkTargetUsers, retry)

	return nil
}

func enableAlarmForInstance(cfg aws.Config, instance InstanceInfo) error {
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	// 创建Instance状态检查
	_, err := cloudwatchClient.EnableAlarmActions(context.TODO(), &cloudwatch.EnableAlarmActionsInput{
		AlarmNames: []string{
			buildAlarmName(InstanceAlarmNameFormat, instance.Name, instance.InstanceId),
			buildAlarmName(SystemAlarmNameFormat, instance.Name, instance.InstanceId),
		},
	})
	if err != nil {
		return err
	}
	if AlertDisableAndEnableAlarm {
		title := "Enable Alarms"
		message := "Enable Alarms for " + instance.Name
		common.SendMessageToChat(title, message, ChatWebhookURL, chatTargetUsers, retry)
		common.SendMessageToLark(title, message, LarkWebhookURL, larkTargetUsers, retry)
	}
	return nil
}

// Disable StatusCheckFailed_Instance and StatusCheckFailed_System Alarm
func disableAlarmForInstance(cfg aws.Config, instance InstanceInfo) error {
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	// 创建Instance状态检查
	_, err := cloudwatchClient.DisableAlarmActions(context.TODO(), &cloudwatch.DisableAlarmActionsInput{
		AlarmNames: []string{
			buildAlarmName(InstanceAlarmNameFormat, instance.Name, instance.InstanceId),
			buildAlarmName(SystemAlarmNameFormat, instance.Name, instance.InstanceId),
		},
	})

	if err != nil {
		return err
	}

	if AlertDisableAndEnableAlarm {
		title := "Disable Alarms"
		message := "Disable Alarms for " + instance.Name
		common.SendMessageToChat(title, message, ChatWebhookURL, chatTargetUsers, retry)
		common.SendMessageToLark(title, message, LarkWebhookURL, larkTargetUsers, retry)
	}

	return nil
}

// 返回所有xxxxx-monitor=on的running EC2
func getAllInstances(cfg aws.Config, instances []InstanceInfo) ([]InstanceInfo, error) {
	tagKey := strings.Split(TriggerTag, "=")[0]
	// create Lambda client to invoke notice lambda
	ec2client := ec2.NewFromConfig(cfg)
	result, err := ec2client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:" + tagKey),
				Values: []string{"on"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceInfo := InstanceInfo{
						InstanceId: *instance.InstanceId,
						State:      string(instance.State.Name),
						Name:       *tag.Value,
					}
					instances = append(instances, instanceInfo)
					break
				}
			}
		}
	}
	fmt.Printf("There are total %d instances in running state\n", len(instances))
	return instances, nil
}

func getSingleInstance(cfg aws.Config, instanceid string) ([]InstanceInfo, error) {
	tagKey := strings.Split(TriggerTag, "=")[0]
	instances := make([]InstanceInfo, 0, 1)
	// Create Lambda client to invoke notice lambda
	ec2client := ec2.NewFromConfig(cfg)
	result, err := ec2client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:" + tagKey),
				Values: []string{"on"},
			},
		},
		InstanceIds: []string{
			instanceid,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceInfo := InstanceInfo{
						InstanceId: *instance.InstanceId,
						// 不应该在查询的时候赋值，因为查询时Instance的state可能已经发生了变化
						// State:      string(instance.State.Name),
						Name: *tag.Value,
					}
					instances = append(instances, instanceInfo)
					break
				}
			}
		}
	}
	return instances, nil
}

// 获取所有的EC2开头的Alarms
func getAllAlarms(cfg aws.Config, nextToken *string, alarms []AlarmInfo) (*string, []AlarmInfo, error) {

	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	result, err := cloudwatchClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{
		AlarmNamePrefix: aws.String(Prefix),
		MaxRecords:      aws.Int32(MaxRecords),
		NextToken:       nextToken,
	})

	if err != nil {
		return nil, nil, err
	}

	for _, alarm := range result.MetricAlarms {
		alarms = append(alarms, AlarmInfo{
			AlarmArn:     *alarm.AlarmArn,
			AlarmName:    *alarm.AlarmName,
			AlarmEnabled: *alarm.ActionsEnabled,
		})
	}
	return result.NextToken, alarms, nil
}

func getAlarmsForInstance(cfg aws.Config, instance InstanceInfo) ([]AlarmInfo, error) {
	alarms := make([]AlarmInfo, 0, 2)
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	result, err := cloudwatchClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{
			buildAlarmName(InstanceAlarmNameFormat, instance.Name, instance.InstanceId),
			buildAlarmName(SystemAlarmNameFormat, instance.Name, instance.InstanceId),
		},
	})

	if err != nil {
		return nil, err
	}

	for _, alarm := range result.MetricAlarms {
		alarms = append(alarms, AlarmInfo{
			AlarmArn:     *alarm.AlarmArn,
			AlarmName:    *alarm.AlarmName,
			AlarmEnabled: *alarm.ActionsEnabled,
		})
	}
	return alarms, nil
}

func buildAlarmName(alarmNameFormat, instanceName, instanceID string) string {
	return strings.Replace(strings.Replace(alarmNameFormat, "{INSTANCE_NAME}", instanceName, -1), "{INSTANCE_ID}", instanceID, -1)
}

func main() {
	lambda.Start(LambdaHandler)
}
