package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type NoticeVar struct {
	FunctionArn *string
}

type NoticeStack struct {
	Stack           awscdk.Stack
	NoticeLambdaArn *string
}

type NoticeStackProps struct {
	StackProps awscdk.StackProps
	Project    string
	LogLevel   string
	EnvName    string
}

func NewNoticeStack(scope constructs.Construct, id *string, props *NoticeStackProps) {

	envName := props.EnvName
	logLevel := props.LogLevel
	noticeEnv := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["noticeEnv"].(map[string]interface{})

	noticeLambdaEnv := map[string]*string{
		"Retry": jsii.String(noticeEnv["Retry"].(string)),
		//"MentionUsers":   jsii.String(noticeEnv["MentionUsers"].(string)),
		"MentionUsers":   jsii.String(""),
		"ParameterStore": jsii.String(noticeEnv["ParameterStore"].(string)),
		"ChatWebhookURL": jsii.String(noticeEnv["NoticeChatWebhookURL"].(string)),
		"LarkWebhookURL": jsii.String(noticeEnv["NoticeLarkWebhookURL"].(string)),
		"ChatChannel":    jsii.String(noticeEnv["NoticeChatChannel"].(string)),
		"LarkChannel":    jsii.String(noticeEnv["NoticeLarkChannel"].(string)),
		"LogLevel":       jsii.String(logLevel),
		"ENV":            jsii.String(envName),
		"Repo":           jsii.String(props.Project),
	}

	alarmLambdaEnv := map[string]*string{
		"Retry":          jsii.String(noticeEnv["Retry"].(string)),
		"MentionUsers":   jsii.String(noticeEnv["MentionUsers"].(string)),
		"ParameterStore": jsii.String(noticeEnv["ParameterStore"].(string)),
		"ChatWebhookURL": jsii.String(noticeEnv["AlarmChatWebhookURL"].(string)),
		"LarkWebhookURL": jsii.String(noticeEnv["AlarmLarkWebhookURL"].(string)),
		"ChatChannel":    jsii.String(noticeEnv["AlarmChatChannel"].(string)),
		"LarkChannel":    jsii.String(noticeEnv["AlarmLarkChannel"].(string)),
		"LogLevel":       jsii.String(logLevel),
		"ENV":            jsii.String(envName),
		"Repo":           jsii.String(props.Project),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)
	newNotice(stack, jsii.String("Notice"), noticeLambdaEnv, alarmLambdaEnv)

}

func newNotice(scope constructs.Construct, id *string, noticeLambdaEnv, alarmLambdaEnv map[string]*string) {
	construct := constructs.NewConstruct(scope, id)

	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		// IAM is global in each partition
		Region:       jsii.String("ap-northeast-1"),
		Service:      jsii.String("ssm"),
		Resource:     jsii.String("parameter"),
		ResourceName: noticeLambdaEnv["ParameterStore"],
	})

	// The code that defines your stack goes here
	noticeLambdaFunction := awslambda.NewFunction(construct, jsii.String("NoticeFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/notice/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &noticeLambdaEnv,
		FunctionName: aws.String(*noticeLambdaEnv["ENV"] + "_sys_notice_lambda"),
	})

	noticeLambdaFunction.AddToRolePolicy(
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

	// created CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("NoticeFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *noticeLambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	noticeTopic := awssns.NewTopic(construct, jsii.String("NoticeFunctionTopic"), &awssns.TopicProps{
		DisplayName: jsii.String(*noticeLambdaEnv["ENV"] + "_sys_notice_topic"),
		TopicName:   jsii.String(*noticeLambdaEnv["ENV"] + "_sys_notice_topic"),
	})

	noticeTopic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(noticeLambdaFunction, &awssnssubscriptions.LambdaSubscriptionProps{}))
	noticeTopic.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("SNS:Publish"),
		},
		Resources: &[]*string{
			noticeTopic.TopicArn(),
		},
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewServicePrincipal(jsii.String("cloudwatch.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		},
	}))

	//---------------------------------------------------------------
	// The code that defines your stack goes here
	alarmLambdaFunction := awslambda.NewFunction(construct, jsii.String("AlarmFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/notice/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &alarmLambdaEnv,
		FunctionName: aws.String(*noticeLambdaEnv["ENV"] + "_sys_alarm_lambda"),
	})

	alarmLambdaFunction.AddToRolePolicy(
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

	// created CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("AlarmFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *alarmLambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	alarmTopic := awssns.NewTopic(construct, jsii.String("AlarmFunctionTopic"), &awssns.TopicProps{
		DisplayName: jsii.String(*noticeLambdaEnv["ENV"] + "_sys_alarm_topic"),
		TopicName:   jsii.String(*noticeLambdaEnv["ENV"] + "_sys_alarm_topic"),
	})

	alarmTopic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(alarmLambdaFunction, &awssnssubscriptions.LambdaSubscriptionProps{}))
	alarmTopic.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("SNS:Publish"),
		},
		Resources: &[]*string{
			alarmTopic.TopicArn(),
		},
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewServicePrincipal(jsii.String("cloudwatch.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		},
	}))

}
