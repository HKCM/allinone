package main

import (
	"bufio"
	"context"
	"fmt"
	"sort"

	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/jsii-runtime-go"
)

const (
	prodExchangeVpc = "vpc-08c8b3647d8c41669"
	prodAdminVpc    = "vpc-0d3136ef8aa4c0ed6"
	filename        = "rule.yaml"
)

var (
	profile string
	region  string
	file    *os.File
	err     error
	cfg     aws.Config

	// targetVpc 与 filterVpc 同时使用代表只列出两个VPC中互相引用的SG而不列出引用其他VPC的SG
	targetVpc = []string{prodExchangeVpc, prodAdminVpc} // 这是希望获取Rule的VPC
	filterVpc = []string{prodExchangeVpc, prodAdminVpc} // 这是filter功能, 只列出属于这两个VPC的rule,如果为空则不使用filter
)

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
	writer := bufio.NewWriter(file)

	ec2Client := ec2.NewFromConfig(cfg)

	for _, vpc := range targetVpc {
		resp, err := getSecurityGroupsByVpc(ec2Client, vpc)
		checkErr(err)
		msg := vpc + ":\n"
		for i, sg := range resp.SecurityGroups {
			log.Printf("%d Checking: %s", i, *sg.GroupName)
			msg = msg + getRulesBySecurityGroup(ec2Client, sg)
			write(writer, msg)
			err = writer.Flush()
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

func getRulesBySecurityGroup(ec2Client *ec2.Client, sg types.SecurityGroup) string {
	basicMsg := fmt.Sprintf("- name: %s\n", *sg.GroupName)
	basicMsg = basicMsg + fmt.Sprintf("  description: %s\n", *sg.Description)
	basicMsg = basicMsg + fmt.Sprintf("  rule: \n")
	inBoundRule := "    inbound:\n"
	outBoundRule := "    outbound:\n"

	rule, err := ec2Client.DescribeSecurityGroupRules(context.TODO(), &ec2.DescribeSecurityGroupRulesInput{
		Filters: []types.Filter{
			{
				Name:   jsii.String("group-id"),
				Values: []string{*sg.GroupId},
			},
		},
	})

	log.Println("There are", len(rule.SecurityGroupRules), "rules in it")
	if err != nil {
		panic(err)
	}

	for _, sgRule := range rule.SecurityGroupRules {
		ruleStr := ""
		// 如果ReferencedGroupInfo不为空,说明引用了其他SG
		if sgRule.ReferencedGroupInfo != nil {
			// 获取引用SG的详细信息
			referencedGroupInfo, _ := ec2Client.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
				GroupIds: []string{*sgRule.ReferencedGroupInfo.GroupId},
			})
			// 这是filter功能
			if len(filterVpc) != 0 {
				// 如果对应的referencedGroupInfo的VPC ID不再vpcFilter中,则不将该rule列出并直接跳过
				if !In(*referencedGroupInfo.SecurityGroups[0].VpcId, filterVpc) {
					continue
				}
			}

			ruleStr = ruleStr + fmt.Sprintf("    - targetGroup: %s\n", *referencedGroupInfo.SecurityGroups[0].GroupName)
			ruleStr = ruleStr + fmt.Sprintf("      targetVpc: %s\n", *referencedGroupInfo.SecurityGroups[0].VpcId)
		} else if sgRule.CidrIpv4 != nil {
			ruleStr = ruleStr + fmt.Sprintf("    - targetCidr: %s\n", *sgRule.CidrIpv4)
		} else if sgRule.CidrIpv6 != nil {
			ruleStr = ruleStr + fmt.Sprintf("    - targetCidrV6: %s\n", *sgRule.CidrIpv6)
		} else {
			ruleStr = ruleStr + fmt.Sprintf("    - targetPrefix: %s\n", *sgRule.PrefixListId)
		}

		ruleStr = ruleStr + fmt.Sprintf("      fromPort: %d\n", *sgRule.FromPort)
		ruleStr = ruleStr + fmt.Sprintf("      toPort: %d\n", *sgRule.ToPort)
		ruleStr = ruleStr + fmt.Sprintf("      ipProtocol: %s\n", *sgRule.IpProtocol)
		if sgRule.Description != nil {
			ruleStr = ruleStr + fmt.Sprintf("      description: %s\n", *sgRule.Description)
		} else {
			ruleStr = ruleStr + fmt.Sprintf("      description: null\n")
		}

		if *sgRule.IsEgress {
			outBoundRule = outBoundRule + ruleStr
		} else {
			inBoundRule = inBoundRule + ruleStr
		}
	}

	msg := basicMsg + inBoundRule + outBoundRule

	return msg

}

func getSecurityGroupsByVpc(ec2Client *ec2.Client, vpc string) (*ec2.DescribeSecurityGroupsOutput, error) {
	resp, err := ec2Client.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   jsii.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})

	return resp, err
}

func write(writer *bufio.Writer, msg string) {
	_, err := writer.WriteString(msg)
	if err != nil {
		log.Panic(err)
	}

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
