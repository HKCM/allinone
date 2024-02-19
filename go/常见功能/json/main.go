package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MyStruct struct {
	Data []struct {
		Name   string
		Policy json.RawMessage `json:"policy"`
	} `json:"data"`
}

func main() {
	// 读取 JSON 文件内容
	filePath := "main.json"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var myStruct MyStruct

	// 解析 JSON 数据
	err = json.Unmarshal(fileContent, &myStruct)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, p := range myStruct.Data {
		fmt.Println(p.Name)
		// 将 testJson 字段作为字符串打印出来
		fmt.Println(p.Name+" as string:", string(p.Policy))
	}

}
