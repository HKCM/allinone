package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RDSAlarmMonitorStack struct {
	Stack awscdk.Stack
}

type RDSAlarmMonitorStackProps struct {
	StackProps awscdk.StackProps
	Project    string
	EnvName    string
	LogLevel   string
}

func NewRDSAlarmMonitorStack(scope constructs.Construct, id *string, props *RDSAlarmMonitorStackProps) {
	// Setup Price Monitor env
	envName := props.EnvName
	logLevel := props.LogLevel
	project := props.Project
	rdsMonitor := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["rdsMonitor"].(map[string]interface{})

	alarmSetupLambdaEnv := map[string]*string{
		"ChatWebhookURL": jsii.String(rdsMonitor["MonitorChatWebhookURL"].(string)),
		"LarkWebhookURL": jsii.String(rdsMonitor["MonitorLarkWebhookURL"].(string)),
		"ChatChannel":    jsii.String(rdsMonitor["MonitorChatChannel"].(string)),
		"LarkChannel":    jsii.String(rdsMonitor["MonitorLarkChannel"].(string)),
		"SNSArn":         jsii.String(rdsMonitor["SNSArn"].(string)),
		"Retry":          jsii.String(rdsMonitor["Retry"].(string)),
		"MentionUsers":   jsii.String(rdsMonitor["MentionUsers"].(string)),
		"ParameterStore": jsii.String(rdsMonitor["ParameterStore"].(string)),
		"ENV":            jsii.String(envName),
		"LogLevel":       jsii.String(logLevel),
		"Repo":           jsii.String(project),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	newRDSAlarmMonitor(stack, id, alarmSetupLambdaEnv)

}

func newRDSAlarmMonitor(scope constructs.Construct, id *string, alarmSetupLambdaEnv map[string]*string) {
	construct := constructs.NewConstruct(scope, id)

	rdsAlarmLambdaFunction := awslambda.NewFunction(construct, jsii.String("RDSAlarmSetupFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/rds-monitor/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(900)),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &alarmSetupLambdaEnv,
		FunctionName: jsii.String(*alarmSetupLambdaEnv["ENV"] + "_sys_RDS_alarm_setup"),
	})

	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		// IAM is global in each partition
		Service:      jsii.String("ssm"),
		Resource:     jsii.String("parameter"),
		ResourceName: alarmSetupLambdaEnv["ParameterStore"],
	})

	rdsAlarmLambdaFunction.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("CloudWatchFullAccessV2")))
	rdsAlarmLambdaFunction.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			// Restrict to listing and describing tables
			Actions: &[]*string{
				jsii.String("ssm:GetParameter"),
				jsii.String("ssm:GetParameters"),
			},
			Resources: &[]*string{
				parameterArn,
			},
		}),
	)
	rdsAlarmLambdaFunction.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			// Restrict to listing and describing tables
			Actions: &[]*string{
				jsii.String("rds:DescribeDBClusters"),
			},
			Resources: &[]*string{
				aws.String("*"),
			},
		}),
	)

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("RDSAlarmSetupFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *rdsAlarmLambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	// 每周二检查所有RDS是否存在对应的警报
	awsevents.NewRule(construct, jsii.String("RDSAlarmSetupAlarmScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute:  jsii.String("0"),
			Hour:    jsii.String("2"),
			WeekDay: jsii.String("3"),
		}),
		Description: jsii.String("This role is created by CloudFormation"),
		Targets:     &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(rdsAlarmLambdaFunction, nil)},
	})

}
