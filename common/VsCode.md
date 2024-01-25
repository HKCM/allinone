# 设置快捷代码

Command+Shift+P,按下图输入>snippets，选择对应文件类型,以go为例
```json
{
	"println":{
		"prefix": "pln",
		"body":"fmt.Println($0)",
		"description": "println"
	},
	"printf":{
		"prefix": "plf",
		"body": "fmt.Printf(\"$0\")",
		"description": "printf"
	}
}
```