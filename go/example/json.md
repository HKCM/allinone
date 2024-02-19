keyword: json序列化和反序列化 

```go
// 定义Order结构体，字段名包括字符串型的Consumer、Phone和Goods,float64类型Price
type Order struct {
	Consumer string `json:"consumer"`
	Phone    string
	Price    float64
	Goods    string `json:"product"`
}

func main() {
	var order1 Order
	// 声明JSON格式的字符串变量str
	str := `{"Consumer":"Consumer1","Phone":"000***1111","Price":1.99,"Goods":"hotdog"}`
	json.Unmarshal([]byte(str), &order1) // 通过JSON反序列化，将JSON格式的字符串转换为结构体
	fmt.Println(order1) // {Consumer1 000***1111 1.99 hotdog}

	order2 := &Order{
		Consumer: "Consumer2",
		Phone:    "000***2222",
		Price:    2.99,
		Goods:    "IceCream",
	}
	jsonStr, _ := json.Marshal(order2) // 通过JSON序列化，将结构体转换为JSON格式的字符串 []byte
	fmt.Println(string(jsonStr)) // {"consumer":"Consumer2","Phone":"000***2222","Price":2.99,"product":"IceCream"}
}
```