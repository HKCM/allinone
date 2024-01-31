package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var (
	profile string
	region  string
)

type Snapshot struct {
	SnapshotId string
	VolumeId   string
	Name       string
	CreateTime string
	Size       int
}

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
	getSnapshot(ec2client)

}

func getSnapshot(ec2client *ec2.Client) {

	fmt.Println("Searching volume...")
	var nextToken *string
	request := &ec2.DescribeSnapshotsInput{
		OwnerIds:  []string{"self"},
		NextToken: nextToken,
	}

	for {
		resp, err := ec2client.DescribeSnapshots(context.TODO(), request)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				panic(fmt.Errorf("code: %s, message: %s", awsErr.Code(), awsErr.Message()))
			}
		}

		if resp != nil {
			fmt.Println("Name,CreateTime,SnapshotId,VolumeId,Size,HasVolume")
			for _, s := range resp.Snapshots {
				var name string
				for _, tag := range s.Tags {
					if *tag.Key == "Name" {
						name = *tag.Value
						break
					}
				}
				if name == "" {
					name = "None"
				} // 手动为Name添加值

				hasVolume := getVolume(ec2client, *s.VolumeId)
				fmt.Printf("%v,%v,%v,%v,%v,%v\n", name, *s.StartTime, *s.SnapshotId, *s.VolumeId, *s.VolumeSize, hasVolume)
			}
		}

		if resp == nil || resp.NextToken == nil {
			break
		}
	}

}

func getVolume(ec2client *ec2.Client, volumeId string) bool {

	request := &ec2.DescribeVolumesInput{
		VolumeIds: []string{volumeId},
	}

	resp, err := ec2client.DescribeVolumes(context.TODO(), request)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			panic(fmt.Errorf("code: %s, message: %s", awsErr.Code(), awsErr.Message()))
		}
	}

	if resp != nil {
		// for _, v := range resp.Volumes {
		// 	fmt.Println("CreateTime,VolumeId,SnapshotId,Size")
		// 	fmt.Printf("%v,%v,%v,%v\n", *v.CreateTime, *v.VolumeId, *v.SnapshotId, *v.Size)
		// }
		return true
	} else {
		// fmt.Println("Volume 不存在")
		return false
	}

}
