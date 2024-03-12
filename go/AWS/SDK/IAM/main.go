package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"gopkg.in/yaml.v3"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/jsii-runtime-go"
)

type void struct{}

const (
	filename = "test.yaml"
)

var (
	profile string
	region  string
	file    *os.File
	err     error
	cfg     aws.Config
	member  void
)

var (
	AWSManagedPolicyPrefix     = "arn:aws:iam::aws:"
	AllManagedPoliciesDocument = make(map[string]void)
	// PolicyArnPrefix            = "arn:aws:iam::123456789012:policy/"
	PolicyArnPrefix = "arn:aws:iam::738595724739:policy/"
	prodTargetVPC   = []string{
		"vpc-08c8b3647d8c41669", // prod exchange
		"vpc-0d3136ef8aa4c0ed6", // prod admin
	}

	stageTargetVPC = []string{
		"vpc-0bd9e232748785adb", // stage exchange
		"vpc-050d2a8c799afffdd", // stage admin
	}
)

type InstanceProfile struct {
	Name string `json:"name" yaml:"name"`
	Role Role   `json:"role" yaml:"role"`
}

type Role struct {
	Name     string   `json:"name" yaml:"name"`
	Policies []Policy `json:"policies" yaml:"policies"`
}

type Policy struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

type PolicyStruct struct {
	Name   string                 `json:"name" yaml:"name"`
	Policy map[string]interface{} `json:"policy"  yaml:"policy"`
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

// 获取targetVpc中现有的所有Security Group 以及对应的Rule
func main() {

	if checkFileIsExist(filename) { //如果文件存在
		// file, err = os.Create(filename) //也可以使用这个每次都创建文件
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
		log.Printf("目标文件 %s 已存在,追加...\n", filename)
	} else {
		file, err = os.Create(filename) //创建文件
		log.Printf("目标文件  %s  不存在,创建文件并写入\n", filename)
	}
	defer file.Close()

	iamClient := iam.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)

	var instanceProfiles []InstanceProfile

	instanceProfileNameSet := getInstanceProfileName(ec2Client)

	checkedRoleName := make(map[string]Role)
	for instanceProfileName := range instanceProfileNameSet {
		ip := InstanceProfile{
			Name: instanceProfileName,
		}

		roleName := describeInstanceProfileRoles(iamClient, instanceProfileName)

		role, ok := checkedRoleName[roleName]
		if !ok { // 如果没有在map中找到 则去获取
			role = getRoleDetail(iamClient, roleName)
			checkedRoleName[roleName] = role
		}
		ip.Role = role
		instanceProfiles = append(instanceProfiles, ip)
	}

	c, e := yaml.Marshal(instanceProfiles)
	checkErr(e)
	writeByOSWriteFile("./instanceProfiles.yaml", c)

	// 获取所有的Policy Document
	var pd []PolicyStruct
	for k := range AllManagedPoliciesDocument {
		pd = append(pd, getManagedPolicyStruct(iamClient, k))
	}

	s, err := json.Marshal(pd)
	checkErr(err)

	writeByOSWriteFile("./pdocument.json", s)

}

// 获取所有Instance的InstanceProfileName
func getInstanceProfileName(ec2Client *ec2.Client) map[string]void {

	// bytes.Buffer的0值可以直接使用
	var buff bytes.Buffer
	instanceProfileNameSet := make(map[string]void)
	var nextToken *string

	log.Println("Start check all instances...")
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
				buff.WriteString("Checking: " + *instance.InstanceId)
				//instanceProfileArn = append(instanceProfileArn, *instance.IamInstanceProfile.Arn)
				if instance.IamInstanceProfile == nil {
					buff.WriteString("-> No IAM role")
				} else {
					roleName := *instance.IamInstanceProfile.Arn
					name := strings.Split(roleName, "/")[1]
					instanceProfileNameSet[name] = member
					buff.WriteString("-> " + name)
				}
				log.Println(buff.String())
				buff.Reset()
			}
		}
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	log.Printf("Total %2d Instance Profiles", len(instanceProfileNameSet))
	return instanceProfileNameSet
}

// 获取InstanceProfile对应的Role
func describeInstanceProfileRoles(iamClient *iam.Client, instanceProfileName string) string {
	var roleNameList []string
	out, err := iamClient.GetInstanceProfile(context.TODO(), &iam.GetInstanceProfileInput{
		InstanceProfileName: jsii.String(instanceProfileName),
	})
	if err != nil {
		panic(err)
	}

	for _, r := range out.InstanceProfile.Roles {
		roleNameList = append(roleNameList, *r.RoleName)
	}

	if len(roleNameList) > 1 {
		panic(fmt.Sprint(instanceProfileName, "具有有两个Role"))
	} else {
		return roleNameList[0]
	}

}

// 通过Role Name 获取对应的Policy, 对于创建Role来说这里已经足够
func getRoleDetail(iamClient *iam.Client, roleName string) Role {
	log.Println("Start get ", roleName, "role policies")
	role := Role{
		Name: roleName,
	}

	var policies []Policy

	inlinePolicy, _ := iamClient.ListRolePolicies(context.TODO(), &iam.ListRolePoliciesInput{
		RoleName: &roleName,
	})

	if len(inlinePolicy.PolicyNames) != 0 {
		log.Println("========>", roleName, "has inline policy <========")
	}

	managedPolicy, _ := iamClient.ListAttachedRolePolicies(context.TODO(), &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	})

	for _, v := range managedPolicy.AttachedPolicies {
		var p Policy
		p.Name = *v.PolicyName
		if strings.HasPrefix(*v.PolicyArn, AWSManagedPolicyPrefix) {

			p.Type = "AWS managed"
		} else {
			p.Type = "Customer managed"
			AllManagedPoliciesDocument[*v.PolicyName] = member
		}
		policies = append(policies, p)
	}

	role.Policies = policies

	return role

}

// 将customer managed policy document写入文件

func getManagedPolicyStruct(iamClient *iam.Client, policyName string) PolicyStruct {
	log.Println("Start write", policyName, "policy document")
	out, err := iamClient.GetPolicy(context.TODO(), &iam.GetPolicyInput{
		PolicyArn: jsii.String(PolicyArnPrefix + policyName),
	})

	checkErr(err)

	outp, err := iamClient.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{
		PolicyArn: out.Policy.Arn,
		VersionId: out.Policy.DefaultVersionId,
	})
	checkErr(err)

	decodedValue, err := url.QueryUnescape(*outp.PolicyVersion.Document)
	checkErr(err)

	// log.Println(decodedValue) // 输出每个policy document
	a := make(map[string]interface{})
	err = json.Unmarshal([]byte(decodedValue), &a)
	checkErr(err)

	p := PolicyStruct{
		Name:   policyName,
		Policy: a,
	}

	return p

}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
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
