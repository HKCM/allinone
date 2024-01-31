package stacks

import (
	"aws-env-monitor/common"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HealthCheckStack struct {
	Stack awscdk.Stack
}

type HealthCheckStackProps struct {
	StackProps      awscdk.StackProps
	NoticeLambdaArn *string
	EnvName         string
}

func NewHealthCheckStack(scope constructs.Construct, id *string, props *HealthCheckStackProps) HealthCheckStack {
	// Setup Price Monitor env
	// envName := props.EnvName
	// healthCheckEnv := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["healthCheck"].(map[string]interface{})
	// healthCheckEnvWebHook := healthCheckEnv["WebhookURL"]
	// healthCheckEnvMentionUsers := healthCheckEnv["MentionUsers"]
	// healthCheckEnvParameterStore := healthCheckEnv["ParameterStore"]

	NoticeLambdaEnv := common.NoticeLambdaMessageEvent{
		// NoticeLambdaArn: props.NoticeLambdaArn,
		// WebhookURL:      jsii.String(healthCheckEnvWebHook.(string)),
		// MentionUsers:    jsii.String(healthCheckEnvMentionUsers.(string)),
		// ParameterStore:  jsii.String(healthCheckEnvParameterStore.(string)),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	// pricemonitor.NewPriceMonitor(stack, jsii.String("PriceMonitor"), NoticeLambdaEnv, anomaly)
	newHealthCheck(stack, jsii.String("HealthCheck"), NoticeLambdaEnv)

	return HealthCheckStack{
		Stack: stack,
	}
}

func newHealthCheck(scope constructs.Construct, id *string, noticeLambdaEnv common.NoticeLambdaMessageEvent) {
	construct := constructs.NewConstruct(scope, id)
	// need a lambda and cloudwatch event trigger
	lambdaFunction := awslambda.NewFunction(construct, jsii.String("HealthCheckFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/healthcheck/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &map[string]*string{
			// common.LambdaEnvName_WebhookURL:      noticeLambdaEnv.WebhookURL,
			// common.LambdaEnvName_MentionUsers:    noticeLambdaEnv.MentionUsers,
			// common.LambdaEnvName_NoticeLambdaArn: noticeLambdaEnv.NoticeLambdaArn,
			// common.LambdaEnvName_ParameterStore:  noticeLambdaEnv.ParameterStore,
		},
	})

	// set up permission for invoke notice lambda
	lambdaFunction.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			// Restrict to listing and describing tables
			Actions: &[]*string{
				jsii.String("lambda:InvokeFunction"),
			},
			Resources: &[]*string{
				// noticeLambdaEnv.NoticeLambdaArn,
			},
		}))

	// build parameterArn and set up lambda permission to get specific Parameter from parameter store
	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		// IAM is global in each partition
		Service:      jsii.String("ssm"),
		Resource:     jsii.String("parameter"),
		ResourceName: jsii.String("xxxxx-mention-user-ids"),
	})

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

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("HealthCheckFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *lambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	// set up cloudWatch event
	awsevents.NewRule(construct, jsii.String("healthCheckEventRule"), &awsevents.RuleProps{
		EventPattern: &awsevents.EventPattern{
			Source:     &[]*string{jsii.String("aws.health")},
			DetailType: &[]*string{jsii.String("AWS Health Event")},
			Detail: &map[string]interface{}{
				"service":           &[]string{"EC2", "RDS", "DMS", "REDSHIFT", "ES"},
				"eventTypeCategory": &[]string{"issue"},
			},
		},
		// bind rule and event
		Targets: &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(lambdaFunction, nil)},
	})

}
