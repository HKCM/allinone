package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsce"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CostMonitorStack struct {
	Stack awscdk.Stack
}

type CostMonitorStackProps struct {
	StackProps awscdk.StackProps
	EnvName    string
	LogLevel   string
	Project    string
}

type Anomaly struct {
	ANOMALY_TOTAL_IMPACT_ABSOLUTE   string
	ANOMALY_TOTAL_IMPACT_PERCENTAGE string
}

func NewCostMonitorStack(scope constructs.Construct, id *string, props *CostMonitorStackProps) CostMonitorStack {
	// Setup Cost Monitor env
	envName := props.EnvName
	costMonitorEnv := scope.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["costMonitor"].(map[string]interface{})

	anomaly := Anomaly{
		ANOMALY_TOTAL_IMPACT_ABSOLUTE:   costMonitorEnv["ANOMALY_TOTAL_IMPACT_ABSOLUTE"].(string),
		ANOMALY_TOTAL_IMPACT_PERCENTAGE: costMonitorEnv["ANOMALY_TOTAL_IMPACT_PERCENTAGE"].(string),
	}

	lambdaEnv := map[string]*string{
		"ChatWebhookURL":   jsii.String(costMonitorEnv["ChatWebhookURL"].(string)),
		"LarkWebhookURL":   jsii.String(costMonitorEnv["LarkWebhookURL"].(string)),
		"ChatChannel":      jsii.String(costMonitorEnv["ChatChannel"].(string)),
		"LarkChannel":      jsii.String(costMonitorEnv["LarkChannel"].(string)),
		"Retry":            jsii.String(costMonitorEnv["Retry"].(string)),
		"MentionUsers":     jsii.String(costMonitorEnv["MentionUsers"].(string)),
		"ParameterStore":   jsii.String(costMonitorEnv["ParameterStore"].(string)),
		"MentionThreshold": jsii.String(costMonitorEnv["ANOMALY_TOTAL_IMPACT_ABSOLUTE"].(string)),
		"Repo":             jsii.String(props.Project),
		"ENV":              jsii.String(envName),
	}

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	newCostMonitor(stack, jsii.String("CostMonitor"), lambdaEnv, anomaly)

	return CostMonitorStack{
		Stack: stack,
	}
}

func newCostMonitor(scope constructs.Construct, id *string, lambdaEnv map[string]*string, anomaly Anomaly) {
	construct := constructs.NewConstruct(scope, id)
	// need a lambda and cloudwatch event trigger
	lambdaFunction := awslambda.NewFunction(construct, jsii.String("CostMonitorFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/cost-monitor/bootstrap.zip"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Architecture: awslambda.Architecture_ARM_64(),
		FunctionName: jsii.String(*lambdaEnv["ENV"] + "_sys_cost_monitor"),
		Environment:  &lambdaEnv,
	})

	topic := awssns.NewTopic(construct, jsii.String("Topic"), &awssns.TopicProps{
		DisplayName: jsii.String("Customer subscription topic"),
	})

	topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(lambdaFunction, &awssnssubscriptions.LambdaSubscriptionProps{}))
	topic.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("SNS:Publish"),
		},
		Resources: &[]*string{
			topic.TopicArn(),
		},
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewServicePrincipal(jsii.String("costalerts.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		},
	}))

	// build parameterArn and set up lambda permission to get specific Parameter from parameter store
	parameterArn := awscdk.Stack_Of(construct).FormatArn(&awscdk.ArnComponents{
		// IAM is global in each partition
		Region:       jsii.String("ap-northeast-1"),
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
	awslogs.NewLogGroup(construct, jsii.String("CostMonitorFunctionLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/lambda/" + *lambdaFunction.FunctionName()),
		Retention:    awslogs.RetentionDays_SIX_MONTHS,
	})

	// Create AnomalyMonitor
	ce := awsce.NewCfnAnomalyMonitor(construct, jsii.String("CostAnomalyMonitor"), &awsce.CfnAnomalyMonitorProps{
		MonitorName:      jsii.String("CostAnomalyMonitor"),
		MonitorType:      jsii.String("DIMENSIONAL"),
		MonitorDimension: jsii.String("SERVICE"),
	})
	expression := `{ "Or": [ 
		{ "Dimensions": { "Key": "ANOMALY_TOTAL_IMPACT_ABSOLUTE",
		"MatchOptions": [ "GREATER_THAN_OR_EQUAL" ], "Values": [ "` +
		anomaly.ANOMALY_TOTAL_IMPACT_ABSOLUTE + `" ] } }, 
		{ "Dimensions": { "Key": "ANOMALY_TOTAL_IMPACT_PERCENTAGE", 
		"MatchOptions": [ "GREATER_THAN_OR_EQUAL" ], "Values": [ "` +
		anomaly.ANOMALY_TOTAL_IMPACT_PERCENTAGE + `" ] } } ] }`

	awsce.NewCfnAnomalySubscription(construct, jsii.String("CostAnomalyMonitorSubscription"), &awsce.CfnAnomalySubscriptionProps{
		Frequency:           jsii.String("IMMEDIATE"), //Daily or weekly frequencies only support Email subscriptions
		ThresholdExpression: jsii.String(expression),
		SubscriptionName:    topic.TopicName(),
		MonitorArnList:      &[]*string{ce.AttrMonitorArn()},
		Subscribers: []interface{}{
			&awsce.CfnAnomalySubscription_SubscriberProperty{
				Address: topic.TopicArn(),
				Type:    jsii.String("SNS"),
				// the properties below are optional
				//Status: jsii.String("status"),
			},
		},
	})

}
