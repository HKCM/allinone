package main

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strings"

	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/jsii-runtime-go"
	"gopkg.in/yaml.v3"
)

var (
	profile string
	region  string
	file    *os.File
	err     error
	cfg     aws.Config
)

var (
	prodTargetVPC = []string{
		"vpc-08c8b3647d8c41669", // prod exchange
		"vpc-0d3136ef8aa4c0ed6", // prod admin
	}

	stageTargetVPC = []string{
		"vpc-0bd9e232748785adb", // stage exchange
		"vpc-050d2a8c799afffdd", // stage admin
	}

	subnetMap = make(map[string]string, 25)
)

type Ec2Config struct {
	Name            string   `yaml:"name"`
	Type            string   `yaml:"type"`
	Az              string   `yaml:"az"`
	Sg              []string `yaml:"sg"`
	InstanceProfile string   `yaml:"instanceProfile"`
	Image           string   `yaml:"image"`
	Subnet          string   `yaml:"subnet"`
	Key             string   `yaml:"key"`
	Tags            []Tag    `yaml:"tags"`
}

type Tag struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func init() {
	flag.StringVar(&profile, "profile", "staging", "AWS Profile")
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
}

func main() {

	// if checkFileIsExist(filename) { //如果文件存在
	// 	// file, err = os.Create(filename) //也可以使用这个每次都创建文件
	// 	file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
	// 	log.Printf("目标文件 %s 已存在,追加...\n", filename)
	// } else {
	// 	file, err = os.Create(filename) //创建文件
	// 	log.Printf("目标文件  %s  不存在,创建文件并写入\n", filename)
	// }
	defer file.Close()

	ec2Client := ec2.NewFromConfig(cfg)

	ec2Config := getInstances(ec2Client)

	c, e := yaml.Marshal(ec2Config)
	checkErr(e)
	writeByOSWriteFile("./ec2Config.yaml", c)

}

func getInstances(ec2Client *ec2.Client) []Ec2Config {

	var nextToken *string

	log.Println("Start check all instances...")
	var ec2Configs []Ec2Config
	for {
		output, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			//   - vpc-id - The ID of the VPC that the instance is running in.
			Filters: []ec2types.Filter{
				{
					Name:   jsii.String("vpc-id"),
					Values: prodTargetVPC,
				},
				{
					Name:   jsii.String("instance-state-name"),
					Values: []string{"running"},
				},
			},
			NextToken: nextToken,
		})
		checkErr(err)
		for _, r := range output.Reservations {
			for _, instance := range r.Instances { // 一个Reservations中可能不止一台机器
				log.Println("Checking: " + *instance.InstanceId)
				ec2config := new(Ec2Config)
				ec2config.Az = *instance.Placement.AvailabilityZone
				for _, sg := range instance.SecurityGroups {
					ec2config.Sg = append(ec2config.Sg, *sg.GroupName)
				}
				if instance.IamInstanceProfile == nil {
					ec2config.InstanceProfile = "NoProfile"
				} else {
					profile := *instance.IamInstanceProfile.Arn
					ec2config.InstanceProfile = strings.Split(profile, "/")[1]
				}

				if subnetName, ok := subnetMap[*instance.SubnetId]; !ok {
					subnetName = getSubnetName(ec2Client, *instance.SubnetId)
					subnetMap[*instance.SubnetId] = subnetName
					ec2config.Subnet = subnetName
				} else {
					ec2config.Subnet = subnetName
				}

				ec2config.Key = "default"
				ec2config.Image = "default"
				ec2config.Type = string(instance.InstanceType)
				tags := make([]Tag, 0, len(instance.Tags))
				for _, t := range instance.Tags {
					tags = append(tags, Tag{Key: *t.Key, Value: *t.Value})
					if *t.Key == "Name" {
						ec2config.Name = *t.Value
					}
				}
				ec2config.Tags = tags
				ec2Configs = append(ec2Configs, *ec2config)
			}
		}
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return ec2Configs
}

func getSubnetName(ec2Client *ec2.Client, subnetId string) string {
	sb, err := ec2Client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetId},
	})
	if err != nil {
		panic("Get subnet failed")
	}

	for _, t := range sb.Subnets[0].Tags {
		if *t.Key == "Name" {
			return *t.Value
		}
	}

	log.Println(subnetId, "No name tag")
	return "NoName"
}

func In(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func writeByOSWriteFile(filename string, content []byte) {
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		panic(err)
	}
}

// bufio.NewWriter
func writeByBufioNewWriter(filename string, content []byte) {
	file, err := os.Open(filename)
	if err != nil {
		panic("open file failed")
	}
	defer file.Close()

	writer := bufio.NewWriter(file) //创建新的 Writer 对象
	// _, err = writer.WriteString(msgString)
	_, err = writer.Write(content)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}
