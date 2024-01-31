package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Ec2StatusMonitorStack struct {
	Stack     awscdk.Stack
	AccountID *string
	Region    *string
}

type Ec2StatusMonitorStackProps struct {
	StackProps awscdk.StackProps
	EnvName    string
	LogLevel   string
}

func NewEc2StatusMonitorStack(scope constructs.Construct, id *string, props *Ec2StatusMonitorStackProps) {
	// Setup Price Monitor env
	envName := props.EnvName
	logLevel := props.LogLevel
	ec2StatusMonitor := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["ec2StatusMonitor"].(map[string]interface{})

	alarmLambdaEnv := map[string]*string{
		"InstanceAlarmNameFormat":    jsii.String(ec2StatusMonitor["StatusCheckFailedInstanceAlarmName"].(string)),
		"SystemAlarmNameFormat":      jsii.String(ec2StatusMonitor["StatusCheckFailedSystemAlarmName"].(string)),
		"ChatWebhookURL":             jsii.String(ec2StatusMonitor["MonitorChatWebhookURL"].(string)),
		"LarkWebhookURL":             jsii.String(ec2StatusMonitor["MonitorLarkWebhookURL"].(string)),
		"ChatChannel":                jsii.String(ec2StatusMonitor["MonitorChatChannel"].(string)),
		"LarkChannel":                jsii.String(ec2StatusMonitor["MonitorLarkChannel"].(string)),
		"Retry":                      jsii.String(ec2StatusMonitor["Retry"].(string)),
		"MentionUsers":               jsii.String(ec2StatusMonitor["MentionUsers"].(string)),
		"ParameterStore":             jsii.String(ec2StatusMonitor["ParameterStore"].(string)),
		"AlertDisableAndEnableAlarm": jsii.String(ec2StatusMonitor["AlertDisableAndEnableAlarm"].(string)),
		"TriggerTag":                 jsii.String(ec2StatusMonitor["TriggerTag"].(string)),
		"ENV":                        jsii.String(envName),
		"LogLevel":                   jsii.String(logLevel),
	}

	alertLambdaEnv := map[string]*string{
		"AlertManagerURL": jsii.String(ec2StatusMonitor["AlertManagerURL"].(string)),
		"xxxxxType":       jsii.String(ec2StatusMonitor["xxxxxType"].(string)),
		"Retry":           jsii.String(ec2StatusMonitor["Retry"].(string)),
		"MentionUsers":    jsii.String(ec2StatusMonitor["MentionUsers"].(string)),
		"VPC":             jsii.String(ec2StatusMonitor["VPC"].(string)),
		"SubnetA":         jsii.String(ec2StatusMonitor["SubnetA"].(string)),
		"SubnetC":         jsii.String(ec2StatusMonitor["SubnetC"].(string)),
		"SecurityGroupID": jsii.String(ec2StatusMonitor["SecurityGroupID"].(string)),
		"ENV":             jsii.String(envName),
		"LogLevel":        jsii.String(logLevel),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	ec2StatusMonitorStack := Ec2StatusMonitorStack{
		Stack:     stack,
		AccountID: props.StackProps.Env.Account,
		Region:    props.StackProps.Env.Region,
	}

	newEc2StatusMonitor(ec2StatusMonitorStack, id, alarmLambdaEnv, alertLambdaEnv)

}

func newEc2StatusMonitor(scope Ec2StatusMonitorStack, id *string, alarmLambdaEnv, alertLambdaEnv map[string]*string) {
	construct := constructs.NewConstruct(scope.Stack, id)

	// need a lambda and cloudwatch event trigger
	// EC2 Status Check monitor
	// will create Status Check alarm when instance in pending
	// and delete Status Check alarm when instance in stopping
	topic := awssns.NewTopic(construct, jsii.String("Topic"), &awssns.TopicProps{
		DisplayName: jsii.String("Customer subscription topic"),
		TopicName:   jsii.String("EC2StatusAlertSNS"),
	})

	alarmLambdaEnv["OKSNSArn"] = topic.TopicArn()
	alarmLambdaEnv["AlarmSNSArn"] = topic.TopicArn()
	alarmLambdaEnv["InsufficientSNSArn"] = topic.TopicArn()

	if *alarmLambdaEnv["ENV"] == "qa" {
		alarmLambdaEnv["OKSNSArn"] = jsii.String("")
	}

	ec2AlarmLambdaFunction := awslambda.NewFunction(construct, jsii.String("EC2StatusCheckAlarmSetupFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/ec2-status-monitor/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(900)),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &alarmLambdaEnv,
		FunctionName: jsii.String(*alarmLambdaEnv["ENV"] + "_sys_EC_alarm_setup"),
	})

	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		// IAM is global in each partition
		Service:      jsii.String("ssm"),
		Resource:     jsii.String("parameter"),
		ResourceName: alarmLambdaEnv["ParameterStore"],
	})

	ec2AlarmLambdaFunction.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ReadOnlyAccess")))
	ec2AlarmLambdaFunction.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("CloudWatchFullAccessV2")))
	ec2AlarmLambdaFunction.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			// Restrict to listing and describing tables
			Actions: &[]*string{
				jsii.String("ssm:GetParameter"),
				jsii.String("ssm:GetParameters"),
			},
			Resources: &[]*string{
				parameterArn,
			},
		}))

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("EC2StatusCheckAlarmSetupFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *ec2AlarmLambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	// set up cloudWatch event
	awsevents.NewRule(construct, jsii.String("EC2StatusCheckAlarmEventRule"), &awsevents.RuleProps{
		EventPattern: &awsevents.EventPattern{
			Source:     &[]*string{jsii.String("aws.ec2")},
			DetailType: &[]*string{jsii.String("EC2 Instance State-change Notification")},
			Detail: &map[string]interface{}{
				"state": &[]string{"terminated", "stopping", "pending"},
			},
		},
		// bind rule and event
		Targets:     &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(ec2AlarmLambdaFunction, nil)},
		Description: jsii.String("This role is created by CloudFormation"),
	})

	awsevents.NewRule(construct, jsii.String("EC2StatusCheckAlarmScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute: jsii.String("50"),
			Hour:   jsii.String("1"),
			Day:    jsii.String("*"),
		}),
		Description: jsii.String("This role is created by CloudFormation"),
		Targets:     &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(ec2AlarmLambdaFunction, nil)},
	})
	//-----------------------------------------------------------------------------------
	// 创建用于报警的lambda,作为StatusCheck failed时发送通知
	ec2AlertLambdaFunction := awslambda.NewFunction(construct, jsii.String("EC2StatusCheckAlertFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/ec2-status-monitor-alert/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(100)),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &alertLambdaEnv,
		FunctionName: jsii.String(*alertLambdaEnv["ENV"] + "_sys_EC_alarm_notice"),
		Vpc: awsec2.Vpc_FromLookup(construct, jsii.String("Vpc"), &awsec2.VpcLookupOptions{
			// Calling this method will lead to a lookup when the CDK CLI is executed.
			// You can therefore not use any values that will only be available at CloudFormation execution time
			// OwnerAccountId: jsii.String(accountID),
			// Region:         jsii.String(region),
			VpcId: alertLambdaEnv["VPC"],
		}),
		VpcSubnets: &awsec2.SubnetSelection{
			Subnets: &[]awsec2.ISubnet{
				awsec2.Subnet_FromSubnetId(construct, jsii.String("SubnetA"), alertLambdaEnv["SubnetA"]),
				awsec2.Subnet_FromSubnetId(construct, jsii.String("SubnetC"), alertLambdaEnv["SubnetC"]),
			},
		},
		SecurityGroups: &[]awsec2.ISecurityGroup{
			awsec2.SecurityGroup_FromLookupById(construct, jsii.String("SGID"), alertLambdaEnv["SecurityGroupID"]),
		},
	})

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("EC2StatusAlertFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *ec2AlertLambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(ec2AlertLambdaFunction, &awssnssubscriptions.LambdaSubscriptionProps{}))
	topic.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("SNS:Publish"),
		},
		Resources: &[]*string{
			topic.TopicArn(),
		},
		Principals: &[]awsiam.IPrincipal{
			// awsiam.NewAccountPrincipal(construct),
			// awsiam.NewServicePrincipal(jsii.String("cloudwatch.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
			awsiam.NewServicePrincipal(jsii.String("cloudwatch.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		},
	}))

}
