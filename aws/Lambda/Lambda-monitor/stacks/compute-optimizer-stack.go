package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ComputeOptimizerStack struct {
	Stack awscdk.Stack
}

type ComputeOptimizerStackProps struct {
	StackProps awscdk.StackProps
	Project    string
	LogLevel   string
	EnvName    string
}

func NewComputeOptimizerStack(scope constructs.Construct, id *string, props *ComputeOptimizerStackProps) ComputeOptimizerStack {
	// Setup Price Monitor env
	envName := props.EnvName
	logLevel := props.LogLevel
	ComputeOptimizerEnv := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["computeOptimizer"].(map[string]interface{})

	lambdaEnv := map[string]*string{
		"ChatWebhookURL":   jsii.String(ComputeOptimizerEnv["ChatWebhookURL"].(string)),
		"LarkWebhookURL":   jsii.String(ComputeOptimizerEnv["LarkWebhookURL"].(string)),
		"ChatChannel":      jsii.String(ComputeOptimizerEnv["ChatChannel"].(string)),
		"LarkChannel":      jsii.String(ComputeOptimizerEnv["LarkChannel"].(string)),
		"Retry":            jsii.String(ComputeOptimizerEnv["Retry"].(string)),
		"MentionUsers":     jsii.String(ComputeOptimizerEnv["MentionUsers"].(string)),
		"ParameterStore":   jsii.String(ComputeOptimizerEnv["ParameterStore"].(string)),
		"ExcludeInstances": jsii.String(ComputeOptimizerEnv["ExcludeInstances"].(string)),
		"Notice":           jsii.String(ComputeOptimizerEnv["Notice"].(string)),
		"Repo":             jsii.String(props.Project),
		"ENV":              jsii.String(envName),
		"LogLevel":         jsii.String(logLevel),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	newComputeOptimizer(stack, jsii.String("ComputeOptimizer"), lambdaEnv)

	return ComputeOptimizerStack{
		Stack: stack,
	}
}

func newComputeOptimizer(scope constructs.Construct, id *string, lambdaEnv map[string]*string) {
	construct := constructs.NewConstruct(scope, id)
	// need a lambda and cloudwatch event trigger
	lambdaFunction := awslambda.NewFunction(construct, jsii.String("computeOptimizerFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/computeoptimizer/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
		Architecture: awslambda.Architecture_ARM_64(),
		FunctionName: jsii.String(*lambdaEnv["ENV"] + "_sys_compute_optimizer"),
		Environment:  &lambdaEnv,
	})

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
		}),
	)

	lambdaFunction.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			// Restrict to listing and describing tables
			Actions: &[]*string{
				jsii.String("compute-optimizer:GetEC2InstanceRecommendations"),
				jsii.String("ec2:DescribeInstances"),
			},
			Resources: &[]*string{jsii.String("*")},
		}),
	)

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("ComputeOptimizerFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *lambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	// set up cloudWatch event
	awsevents.NewRule(construct, jsii.String("ComputeOptimizerScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute:  jsii.String("50"),
			Hour:    jsii.String("1"),
			WeekDay: jsii.String("3"),
		}),
		Description: jsii.String("This role is created by CloudFormation"),
		Targets:     &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(lambdaFunction, nil)},
	})

}
