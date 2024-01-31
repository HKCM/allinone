package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type S3LifecycleMonitorStack struct {
	Stack awscdk.Stack
}

type S3LifecycleMonitorStackProps struct {
	StackProps awscdk.StackProps
	Project    string
	EnvName    string
	LogLevel   string
}

func NewS3LifecycleMonitorStack(scope constructs.Construct, id *string, props *S3LifecycleMonitorStackProps) S3LifecycleMonitorStack {

	// Setup s3 Lifecycle Monitor env
	envName := props.EnvName
	logLevel := props.LogLevel
	project := props.Project
	s3LifecycleMonitorEnv := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["s3LifecycleMonitor"].(map[string]interface{})

	s3LifecycleLambdaEnv := map[string]*string{
		"ChatWebhookURL":      jsii.String(s3LifecycleMonitorEnv["ChatWebhookURL"].(string)),
		"LarkWebhookURL":      jsii.String(s3LifecycleMonitorEnv["LarkWebhookURL"].(string)),
		"ChatChannel":         jsii.String(s3LifecycleMonitorEnv["ChatChannel"].(string)),
		"LarkChannel":         jsii.String(s3LifecycleMonitorEnv["LarkChannel"].(string)),
		"SNSArn":              jsii.String(s3LifecycleMonitorEnv["SNSArn"].(string)),
		"Retry":               jsii.String(s3LifecycleMonitorEnv["Retry"].(string)),
		"MentionUsers":        jsii.String(s3LifecycleMonitorEnv["MentionUsers"].(string)),
		"ParameterStore":      jsii.String(s3LifecycleMonitorEnv["ParameterStore"].(string)),
		"ShowEnableLifecycle": jsii.String(s3LifecycleMonitorEnv["ShowEnableLifecycle"].(string)),
		"ShowNoLifecycle":     jsii.String(s3LifecycleMonitorEnv["ShowNoLifecycle"].(string)),
		"ExcludeBuckets":      jsii.String(s3LifecycleMonitorEnv["ExcludeBuckets"].(string)),
		"ENV":                 jsii.String(envName),
		"LogLevel":            jsii.String(logLevel),
		"Repo":                jsii.String(project),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	newS3LifecycleMonitor(stack, jsii.String("S3LifecycleMonitor"), s3LifecycleLambdaEnv)

	return S3LifecycleMonitorStack{
		Stack: stack,
	}
}

func newS3LifecycleMonitor(scope constructs.Construct, id *string, s3LifecycleLambdaEnv map[string]*string) {
	construct := constructs.NewConstruct(scope, id)
	// need a lambda and cloudwatch event trigger
	lambdaFunction := awslambda.NewFunction(construct, jsii.String("S3LifecycleMonitorFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/s3-lifecycle-monitor/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &s3LifecycleLambdaEnv,
		FunctionName: jsii.String(*s3LifecycleLambdaEnv["ENV"] + "_sys_s3_lifecycle_notice"),
	})

	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		// IAM is global in each partition
		Service:      jsii.String("ssm"),
		Resource:     jsii.String("parameter"),
		ResourceName: jsii.String("xxxxx-mention-user-ids"),
	})
	// set up permission for parameter store
	lambdaFunction.AddToRolePolicy(
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

	lambdaFunction.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			// Restrict to listing and describing tables
			Actions: &[]*string{
				jsii.String("s3:List*"),
				jsii.String("s3:GetLifecycleConfiguration"),
			},
			Resources: &[]*string{
				jsii.String("*"),
			},
		}))

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, aws.String("S3LifecycleMonitorFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: aws.String("/aws/lambda/" + *lambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	// // set up cloudWatch event
	// awsevents.NewRule(construct, aws.String("cloudWatchEventRule"), &awsevents.RuleProps{
	// 	EventPattern: &awsevents.EventPattern{
	// 		Source:     &[]*string{aws.String("aws.s3")},
	// 		DetailType: &[]*string{aws.String("AWS API Call via CloudTrail")},
	// 		Detail: &map[string]interface{}{
	// 			"eventSource": &[]string{"s3.amazonaws.com"},
	// 			"eventName":   &[]string{"DeleteBucketLifecycle", "CreateBucket"},
	// 		},
	// 	},
	// 	// bind rule and event
	// 	Targets: &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(lambdaFunction, nil)},
	// })
	// 每周二检查所有RDS是否存在对应的警报
	awsevents.NewRule(construct, jsii.String("S3LifecycleCheckScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute:  jsii.String("0"),
			Hour:    jsii.String("2"),
			WeekDay: jsii.String("4"),
		}),
		Description: jsii.String("This role is created by CloudFormation"),
		Targets:     &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(lambdaFunction, nil)},
	})
}
