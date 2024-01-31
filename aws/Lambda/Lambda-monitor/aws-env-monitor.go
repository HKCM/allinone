package main

import (
	"aws-env-monitor/common"
	"aws-env-monitor/stacks"
	"fmt"
	"time"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

type AwsEnvMonitorStackProps struct {
	awscdk.StackProps
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	envName := common.GetContextStr(app, "env", "qa")
	maintainer := app.Node().TryGetContext(jsii.String("maintainer")).(string)
	logLevel := app.Node().TryGetContext(jsii.String("logLevel")).(string)
	project := app.Node().TryGetContext(jsii.String("project")).(string)

	deployDate := time.Now().Format("2006-01-02 15:04:05") + " JST"
	tags := &map[string]*string{
		"Env":        &envName,
		"Maintainer": &maintainer,
		"Project":    &project,
		"DeployDate": &deployDate,
	}

	// Create Notification stack
	stacks.NewNoticeStack(app, jsii.String("NoticeStack"), &stacks.NoticeStackProps{
		StackProps: awscdk.StackProps{
			Tags:                  tags,
			TerminationProtection: jsii.Bool(true),
		},
		EnvName:  envName,
		LogLevel: logLevel,
		Project:  project,
	})

	// Create s3 Lifecycle Monitor Stack
	// s3LifecycleMonitorStack :=
	stacks.NewS3LifecycleMonitorStack(app, jsii.String("s3LifecycleMonitorStack"), &stacks.S3LifecycleMonitorStackProps{
		StackProps: awscdk.StackProps{
			Tags: tags,
		},
		EnvName:  envName,
		LogLevel: logLevel,
		Project:  project,
	})

	// s3LifecycleMonitorStack.Stack.AddDependency(noticeStack.Stack, nil)

	// Create price Monitor Stack
	stacks.NewCostMonitorStack(app, jsii.String("costMonitorStack"), &stacks.CostMonitorStackProps{
		StackProps: awscdk.StackProps{
			Tags: tags,
			Env: &awscdk.Environment{
				Region: jsii.String("us-east-1"),
			},
			TerminationProtection: jsii.Bool(true),
		},
		EnvName:  envName,
		LogLevel: logLevel,
		Project:  project,
	})

	stacks.NewComputeOptimizerStack(app, jsii.String("computeOptimizerStack"), &stacks.ComputeOptimizerStackProps{
		StackProps: awscdk.StackProps{
			Tags: tags,
		},
		EnvName:  envName,
		LogLevel: logLevel,
		Project:  project,
	})

	// Create health check Monitor Stack
	// stacks.NewHealthCheckStack(app, jsii.String("healthCheckStack"), &stacks.HealthCheckStackProps{
	// 	StackProps: awscdk.StackProps{
	// 		Tags: tags,
	// 	},
	// 	NoticeLambdaArn: noticeStack.NoticeLambdaArn,
	// 	EnvName:         envName,
	// })

	stacks.NewOperationMonitorStack(app, jsii.String("operationMonitorStack"), &stacks.OperationMonitorStackProps{
		StackProps: awscdk.StackProps{
			Tags: tags,
		},
		EnvName:  envName,
		LogLevel: logLevel,
		Project:  project,
	})

	stacks.NewRDSAlarmMonitorStack(app, jsii.String("rdsMonitorStack"), &stacks.RDSAlarmMonitorStackProps{
		StackProps: awscdk.StackProps{
			Tags:                  tags,
			TerminationProtection: jsii.Bool(true),
		},
		LogLevel: logLevel,
		EnvName:  envName,
		Project:  project,
	})

	// 这两个变量是为了Lambda能顺利设置VPC
	// 使用 awscdk.Aws_ACCOUNT_ID() 和awscdk.Aws_REGION()
	// 这样的动态变量是不行的
	accountID := app.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["ec2StatusMonitor"].(map[string]interface{})["Account"].(string)
	region := app.Node().TryGetContext(jsii.String(envName)).(map[string]interface{})["ec2StatusMonitor"].(map[string]interface{})["Region"].(string)
	fmt.Println(accountID, region)
	stacks.NewEc2StatusMonitorStack(app, jsii.String("ec2StatusMonitorStack"), &stacks.Ec2StatusMonitorStackProps{
		StackProps: awscdk.StackProps{
			Tags:                  tags,
			TerminationProtection: jsii.Bool(true),
			Env: &awscdk.Environment{
				Region:  jsii.String(region),
				Account: jsii.String(accountID),
			},
		},
		LogLevel: logLevel,
		EnvName:  envName,
	})

	app.Synth(nil)
}
