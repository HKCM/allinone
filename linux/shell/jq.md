file.json 
```json
{
  "name": "John Doe",
  "age": 30,
  "city": "New York",
  "isStudent": false,
  "scores": [85, 90, 78],
  "grades": {
    "math": "A",
    "english": "B",
    "history": "C"
  },
  "pets": [
    {
      "type": "dog",
      "name": "Buddy"
    },
    {
      "type": "cat",
      "name": "Whiskers"
    }
  ]
}
```

```bash
# 打印整个 JSON
cat file.json | jq .

# 选择特定字段
cat file.json | jq '.name'
"John Doe"

# 过滤数组
cat file.json | jq '.scores[]'
cat file.json | jq '.scores[1]'
90
# 计算数组长度
cat file.json | jq '.scores | length'

# 过滤和条件
cat file.json | jq '.scores[] | select(. >= 85)'
85
90

# 格式化输出
cat file.json | jq '.' --indent 4

# 获取所有的key
cat file.json | jq '. | keys'
[
  "age",
  "city",
  "grades",
  "isStudent",
  "name",
  "pets",
  "scores"
]

```