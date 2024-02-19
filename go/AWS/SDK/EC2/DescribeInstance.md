keyword: list instances, describeInstances, 列出ec2
```go
package main

import (
	"context"
	"fmt"
	"strings"

	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/jsii-runtime-go"
)

var (
	profile string
	region  string
	err     error
	cfg     aws.Config
	client  *ec2.Client
)

func init() {
	flag.StringVar(&profile, "profile", "prod", "AWS Profile")
	flag.StringVar(&region, "region", "ap-northeast-1", "AWS region")
	flag.Parse()
	cfg, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		panic(err)
	}
	client = ec2.NewFromConfig(cfg)
}

func main() {

	// filter_name := "image-id"
	// filter_value := "ami-003488c9c9e37d326"
	// filter_name := "tag-key" 列出包含指定tag的机器,无论tag的值是什么
	// filter_value := "backup"
	filter_name := "tag:backup" // 列出包含指定tag的机器,且tag的值为weekly
	filter_value := "weekly"

	o, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: jsii.String(filter_name),
				Values: []string{
					filter_value,
				},
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}
	var s int = 0
	var builder strings.Builder
	for _, instances := range o.Reservations {
		for _, i := range instances.Instances { // 一个Reservations中可能不止一台机器
			s++
			for _, v := range i.Tags {
				if *v.Key == "Name" {
					builder.WriteString(*v.Value + "\n")
				}
			}
		}
	}
	fmt.Println(builder.String())
	fmt.Println(len(o.Reservations))
	fmt.Println(s)
}
```