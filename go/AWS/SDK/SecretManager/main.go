package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/jsii-runtime-go"
	"gopkg.in/yaml.v3"
)

type Secret struct {
	Name   string `json:"name" yaml:"secret_name"`
	KmsKey string `json:"kmskey" yaml:"kms_alias"`
	KmsArn string `json:"-" yaml:"-"`
}

var (
	profile        string
	region         string
	err            error
	cfg            aws.Config
	sourceFilename = "./secretNames.txt"
	outputFilename = "./test.yaml"
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

func main() {

	secretNames := getConfigByLine(sourceFilename)
	secretClient := secretsmanager.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	secrets := getSecrets(secretClient, secretNames)
	//secrets := getSecretsByNames(secretClient, secretNames)
	log.Print(len(secrets))
	//secrets = append(secrets, Secret{Name: secretName, KmsKey: id})

	for i := 0; i < len(secrets); i++ {
		alias := getKmsAliasById(kmsClient, secrets[i].KmsArn)
		secrets[i].KmsKey = alias
	}

	bytes, err := yaml.Marshal(secrets)
	if err != nil {
		panic(err)
	}
	writeBytesToFile(outputFilename, bytes)

	//log.Print(secrets)
}

// 获取配置文件中的secret
func getConfigByLine(filename string) (secretNames []string) {

	file, e := os.Open(filename)
	if e != nil {
		log.Println("打开文件失败")
		return
	}
	defer file.Close()
	// read the file line by line using scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		secretNames = append(secretNames, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

// // 通过ListSecrets 一次获取10个secret 存在小问题尽量不用
func getSecrets(secretClient *secretsmanager.Client, secretNames []string) []Secret {
	secrets := make([]Secret, 0, len(secretNames))

	for i := 0; i < len(secretNames); i += 10 {
		var secretNames_10 []string
		end := i + 10
		if end > len(secretNames) {
			end = len(secretNames)
		}

		secretNames_10 = append(secretNames_10, secretNames[i:end]...)

		//log.Printf("Small Array %v\n", secretNames_10)
		// 一次最多查10个
		output, err := secretClient.ListSecrets(context.TODO(), &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{
					Key:    types.FilterNameStringTypeName,
					Values: secretNames_10,
				},
			},
		})
		// 获取单个Secret
		// output, err := secretClient.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
		// 	SecretId: jsii.String(secretName),
		// })
		if err != nil {
			panic(err)
		}

		var KmsKeyArn string
		log.Println(len(output.SecretList))

		for _, sle := range output.SecretList {
			if sle.KmsKeyId == nil {
				KmsKeyArn = "default"
			} else {
				KmsKeyArn = *sle.KmsKeyId
			}
			log.Println(*sle.Name)

			secrets = append(secrets, Secret{Name: *sle.Name, KmsArn: KmsKeyArn})
		}
	}

	return secrets
}

// 通过DescribeSecret逐一获取secret
func getSecretsByNames(secretClient *secretsmanager.Client, secretNames []string) []Secret {
	secrets := make([]Secret, 0, len(secretNames))

	for _, secretName := range secretNames {
		// 获取单个Secret
		output, err := secretClient.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
			SecretId: jsii.String(secretName),
		})
		if err != nil {
			// TODO 遇到权限错误应该用recover 恢复
			panic(err)
		}
		secrets = append(secrets, Secret{Name: *output.Name, KmsArn: *output.KmsKeyId})
	}

	log.Print(len(secrets))
	return secrets
}

func getKmsAliasById(kmsClient *kms.Client, kmsArn string) string {

	if kmsArn == "default" {
		return "default"
	}

	output, err := kmsClient.ListAliases(context.TODO(), &kms.ListAliasesInput{
		KeyId: jsii.String(kmsArn),
	})

	// output, err := kmsClient.DescribeKey(context.TODO(), &kms.DescribeKeyInput{
	// 	KeyId: jsii.String(kmsArn),
	// })
	if err != nil {
		panic(err)
	}
	if len(output.Aliases) > 1 {
		log.Println("===========>", kmsArn, "<===========")
	}

	return *output.Aliases[0].AliasName
}

// file.Write 可以追加到文件
func writeBytesToFile(filename string, b []byte) {
	var file *os.File
	var err error
	if checkFileIsExist(filename) { //如果文件存在
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644) //打开文件
		log.Println("文件存在,追加内容")
	} else {
		file, err = os.Create(filename) //创建文件
		log.Println("文件不存在")
	}
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// n, err := file.WriteString(writeString) //写入文件(字符串)
	_, err = file.Write(b) //写入文件([]byte)
	if err != nil {
		panic(err)
	}
	//log.Printf("写入 %d 个字节", n)
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
