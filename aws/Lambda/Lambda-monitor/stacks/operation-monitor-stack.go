package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogsdestinations"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type OperationMonitorStack struct {
	Stack awscdk.Stack
}

type OperationMonitorStackProps struct {
	StackProps awscdk.StackProps
	Project    string
	LogLevel   string
	EnvName    string
}

func NewOperationMonitorStack(scope constructs.Construct, id *string, props *OperationMonitorStackProps) {

	// Setup s3 Lifecycle Monitor env
	envName := props.EnvName
	logLevel := props.LogLevel
	operationMonitorEnv := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["operationMonitor"].(map[string]interface{})
	targetLogGroup := operationMonitorEnv["TargetLogGroup"].(string)
	lambdaEnv := map[string]*string{
		"ChatWebhookURL":    jsii.String(operationMonitorEnv["ChatWebhookURL"].(string)),
		"LarkWebhookURL":    jsii.String(operationMonitorEnv["LarkWebhookURL"].(string)),
		"WorkTime":          jsii.String(operationMonitorEnv["WorkTime"].(string)),
		"ChatChannel":       jsii.String(operationMonitorEnv["ChatChannel"].(string)),
		"LarkChannel":       jsii.String(operationMonitorEnv["LarkChannel"].(string)),
		"Retry":             jsii.String(operationMonitorEnv["Retry"].(string)),
		"MentionUsers":      jsii.String(operationMonitorEnv["MentionUsers"].(string)),
		"ParameterStore":    jsii.String(operationMonitorEnv["ParameterStore"].(string)),
		"LogLevel":          jsii.String(logLevel),
		"ExcludeUsers":      jsii.String(operationMonitorEnv["ExcludeUsers"].(string)),
		"ExcludeEvents":     jsii.String(operationMonitorEnv["ExcludeEvents"].(string)),
		"ExcludeUserEvents": jsii.String(operationMonitorEnv["ExcludeUserEvents"].(string)),
		"MentionEvents":     jsii.String(operationMonitorEnv["MentionEvents"].(string)),
		"ENV":               jsii.String(envName),
		"TargetLogGroup":    jsii.String(targetLogGroup),
		"Repo":              jsii.String(props.Project),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	NewOperationMonitor(stack, jsii.String("OperationMonitor"), targetLogGroup, lambdaEnv)

}

func NewOperationMonitor(scope constructs.Construct, id *string, targetLogGroup string, lambdaEnv map[string]*string) {
	construct := constructs.NewConstruct(scope, id)
	// need a lambda and cloudwatch event trigger
	lambdaFunction := awslambda.NewFunction(construct, jsii.String("OperationMonitorFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/operation-monitor/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
		Architecture: awslambda.Architecture_ARM_64(),
		FunctionName: jsii.String(*lambdaEnv["ENV"] + "_sys_operation_monitor"),
		Environment:  &lambdaEnv,
	})

	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		Service:      jsii.String("ssm"),
		Resource:     jsii.String("parameter"),
		ResourceName: jsii.String(*lambdaEnv["ParameterStore"]),
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

	// create CloudWatch Logs group
	awslogs.NewLogGroup(construct, jsii.String("OperationMonitorFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *lambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	awslogs.NewSubscriptionFilter(construct, jsii.String("Subscription"), &awslogs.SubscriptionFilterProps{
		LogGroup:    awslogs.LogGroup_FromLogGroupName(construct, jsii.String("CloudTrailLogGroup"), jsii.String(targetLogGroup)),
		Destination: awslogsdestinations.NewLambdaDestination(lambdaFunction, &awslogsdestinations.LambdaDestinationOptions{}),
		FilterPattern: awslogs.FilterPattern_All(
			awslogs.FilterPattern_BooleanValue(jsii.String("$.readOnly"), jsii.Bool(false)),
			// we have SSM event notice not need in here
			awslogs.FilterPattern_StringValue(jsii.String("$.eventName"), jsii.String("!="), jsii.String("StartSession")),
			// SSM automation action
			awslogs.FilterPattern_StringValue(jsii.String("$.eventName"), jsii.String("!="), jsii.String("UpdateInstanceAssociationStatus")),
			awslogs.FilterPattern_StringValue(jsii.String("$.eventName"), jsii.String("!="), jsii.String("UpdateInstanceInformation")),
			awslogs.FilterPattern_StringValue(jsii.String("$.eventName"), jsii.String("!="), jsii.String("PutInventory")),
			//awslogs.FilterPattern_StringValue(jsii.String("$.awsRegion"), jsii.String("!="), jsii.String("ap-northeast-1")),
			//awslogs.FilterPattern_StringValue(jsii.String("$.eventSource"), jsii.String("!="), jsii.String("signin.amazonaws.com")),
			//awslogs.FilterPattern_StringValue(jsii.String("$.eventSource"), jsii.String("!="), jsii.String("iam.amazonaws.com")),
			//awslogs.FilterPattern_StringValue(jsii.String("$.eventSource"), jsii.String("!="), jsii.String("cloudfront.amazonaws.com")),
			//awslogs.FilterPattern_StringValue(jsii.String("$.eventSource"), jsii.String("!="), jsii.String("trustedadvisor.amazonaws.com")),
			//awslogs.FilterPattern_StringValue(jsii.String("$.eventSource"), jsii.String("!="), jsii.String("ssm.amazonaws.com")),
		),
	})
}
