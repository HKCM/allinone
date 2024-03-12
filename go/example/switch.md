
switch/case语句之后是一个表达式,匹配项后面不需要再加break.

如果在case语句中增加fallthrough，则会强制执行后面的case语句，并且不会判断下一条case的表达式结果是否为真，也叫switch穿透
```go
func main() {
	// 定义局部变量，classification用于存放商品类型，commodity代表用户当前选择的商品名称
	var classification string
	var commodity string = "b"
	switch commodity {
	case "a":
		classification = "书籍"
	case "b":
		classification = "数码产品"
	case "c":
		classification = "厨房电器"
	case "d":
		classification = "家电"
	default:
		classification = "其他商品"
	}

	switch {
	case classification == "书籍":
		fmt.Printf("commodity是: %s\n", commodity)
	case classification == "数码产品":
		fmt.Printf("commodity是: %s\n", commodity)
	case classification == "厨房电器":
		fmt.Printf("commodity是: %s\n", commodity)
	case classification == "家电":
		fmt.Printf("commodity是: %s\n", commodity)
	default:
		fmt.Printf("commodity是: %s\n", commodity)
	}
	fmt.Printf("商品类型是: %s\n", classification)
}
```

```go
package main
import "fmt"
func main() {
    switch {
    //false，肯定不会执行
    case false:
        fmt.Println（"case 1为false"）
        fallthrough
    //true，肯定执行
    case true:
        fmt.Println（"case 2为 true"）
        fallthrough
    //由于上一个case中有fallthrough，即使是false，也强制执行
    case false:
        fmt.Println（"case 3为 false"）
        fallthrough
    default:
        fmt.Println（"默认 case"）
    }
}

//case 2为 true
//case 3为 false
//默认 case
```

interface 类型判断
```go
func main() {
	var t interface{}
	t = 10.10

	switch t.(type) {
	case string:
		fmt.Println("string")
	case nil:
		fmt.Println("nil")
	case int:
		fmt.Println("int")
	case float64:
		fmt.Println("float64")
	default:
		fmt.Println("Other type")
	}
}
```