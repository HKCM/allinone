package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var (
	LifecycleName = "DeleteAfterSixMonths"
	SpecificKey   = "EnableLifecycle"
	Status        = "false"
	errMessage    = "NoSuchLifecycleConfiguration"
	profile       string
)

var (
	AccountInfo = map[string]string{
		"012345678901": "Prod",
		"012345678901": "Staging",
	}
)

func init() {
	flag.StringVar(&profile, "profile", "None", "AWS profile")
}

func main() {
	flag.Parse()

	if profile == "None" {
		log.Fatalf("Error parameter\n\nUsage:\ngo run scripts/enableLifecycle/enableLifecycle.go --profile <profile-name>")
	}

	err := checkAccount(profile)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Canceled..")
		return
	}

	fmt.Println("********** Start Enable Bucket Lifecycle **********")
	enableLifecycle(profile)

}

func enableLifecycle(profile string) {
	// Load the Shared AWS Configuration (~/.aws/config)
	var cfg aws.Config
	var err error
	if profile == "" {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	}

	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// List all bucket
	resp, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return
	}

	// Iterate through all the buckets to check if a special tag exists
	for _, bucket := range resp.Buckets {
		bucketName := aws.ToString(bucket.Name)
		fmt.Printf("bucketName:\t%s\n", bucketName)

		Enable, err := checkBucketSetting(client, bucketName)
		if err != nil {
			continue
		}

		if Enable {
			putBucketLifecycleConfiguration(*client, bucketName)
		}
	}

	fmt.Println("********** Enable Bucket Lifecycle Finished **********")
}

func putBucketLifecycleConfiguration(client s3.Client, bucketName string) {

	fmt.Printf("%s is putting Bucket Lifecycle Configuration\n", bucketName)

	// Create a lifecycle configuration rule
	lifecycleConfig := &types.BucketLifecycleConfiguration{
		Rules: []types.LifecycleRule{
			{
				ID:     &LifecycleName,
				Status: types.ExpirationStatusEnabled,
				Expiration: &types.LifecycleExpiration{
					Days: aws.Int32(180),
				},
				Filter: &types.LifecycleRuleFilterMemberPrefix{
					Value: "",
				},
			},
		},
	}

	// Define lifecycle configuration requests
	req := &s3.PutBucketLifecycleConfigurationInput{
		Bucket:                 &bucketName,
		LifecycleConfiguration: lifecycleConfig,
	}

	// Put a bucket lifecycle configuration
	_, err := client.PutBucketLifecycleConfiguration(context.TODO(), req)
	if err != nil {
		fmt.Printf("Failed to configure lifecycle for bucket %s: %v\n\n", bucketName, err)
	} else {
		fmt.Printf("Lifecycle configured for bucket %s\n\n", bucketName)
	}
}

func checkBucketSetting(client *s3.Client, bucketName string) (bool, error) {
	// maybe there is no tags on the bucket
	hasTags := true
	enable := true
	// get tags from the bucket
	result, err := client.GetBucketTagging(context.TODO(), &s3.GetBucketTaggingInput{
		Bucket: &bucketName,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				return false, err.(awserr.Error)
			}
		} else {
			// If there are no tag on bucket, it will get
			// operation error S3: GetBucketTagging
			// NoSuchTagSet: The TagSet does not exist
			// fmt.Println(err.Error())
			hasTags = false
			fmt.Println("There is no tag on bucket.")
		}
	}

	if hasTags {
		for _, tag := range result.TagSet {
			// fmt.Println(*tag.Key, ":", *tag.Value)
			if *tag.Key == SpecificKey && *tag.Value == Status {
				enable = false
				fmt.Printf("There is a special tag: EnableLifecycle:false, skip...\n\n")
				return enable, nil
			}
		}
	}

	_, err = client.GetBucketLifecycleConfiguration(context.TODO(),
		&s3.GetBucketLifecycleConfigurationInput{
			Bucket: &bucketName,
		})

	// If we can get the bucket Lifecycle Configuration
	// which means there is already have Lifecycle Configuration
	// we don't need to change it
	if err == nil {
		enable = false
		fmt.Printf("There already have Lifecycle Configuration, skip...\n\n")
		return enable, nil
	}

	return enable, nil
}

func checkAccount(profile string) error {

	// Load the Shared AWS Configuration (~/.aws/config)
	var cfg aws.Config
	var err error

	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	if err != nil {
		return err
	}

	// get aws account id and map to env
	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	account := *identity.Account

	env, ok := AccountInfo[account]
	if !ok {
		env = "UNKNOWN"
	}

	user := strings.Split(*identity.Arn, "/")[1]

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Will enable lifecycle for all bucket in your account(except already have lifecycle or specified tags)")
	fmt.Printf("Hi, %s. Are you sure run script in %s(%s) environment？(Y/N): ", user, env, account)

	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if strings.EqualFold(input, "y") {
		return nil
	} else {
		return errors.New("No confirm")
	}

}
