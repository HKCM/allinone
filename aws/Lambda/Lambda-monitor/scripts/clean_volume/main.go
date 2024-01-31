package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var (
	profile string
	region  string
)

func init() {
	flag.StringVar(&profile, "profile", "None", "AWS profile")
	flag.StringVar(&region, "region", "ap-northeast-1", "AWS region")

	flag.Parse()
	fmt.Printf("profile: %v\n", profile)
	fmt.Printf("region: %v\n", region)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		fmt.Printf("unable to load SDK config, %v\n", err)
	}

	// Using the Config value, create the DynamoDB client
	ec2client := ec2.NewFromConfig(cfg)
	var nextToken *string
	getAvailableVolumes(ec2client, nextToken)
}

func getAvailableVolumes(ec2client *ec2.Client, nextToken *string) ([]string, error) {
	num := 0
	for {
		fmt.Println("searching...")
		request := &ec2.DescribeVolumesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("status"),
					Values: []string{"available"},
				},
			},
			NextToken: nextToken,
		}

		resp, err := ec2client.DescribeVolumes(context.TODO(), request)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				return nil, fmt.Errorf("code: %s, message: %s", awsErr.Code(), awsErr.Message())
			}
			return nil, err
		}

		fmt.Println("CreateTime,VolumeId,SnapshotId,Size")
		for _, v := range resp.Volumes {
			num += 1
			fmt.Printf("%v,%v,%v,%v\n", *v.CreateTime, *v.VolumeId, *v.SnapshotId, *v.Size)
		}

		nextToken = resp.NextToken

		if nextToken == nil {
			break
		}

	}
	fmt.Printf("Total num: %v\n", num)
	return nil, nil

}
