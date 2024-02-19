keyword: map嵌套

```go
literatures := []string{"《宋词选》", "《家》", "《平凡的世界》"}
humanities := []string{"《实践论》", "《论语译注》"}
//声明多层map嵌套bookRegion，并通过make()函数初始化
bookRegion := make(map[string]map[string][]string)
bookMap := make(map[string][]string)
bookMap["文学"] = literatures
bookMap["人文社科"] = humanities
bookRegion["图书"] = bookMap //将类型为map的变量bookMap作为value映射到bookRegion的key中
//打印图书区域下的分类
fmt.Printf("图书区域下的分类： %s\n", reflect.ValueOf(bookRegion["图书"]).MapKeys())
fmt.Printf("图书区域下的文学类图书： %s\n", bookRegion["图书"]["文学"])
fmt.Printf("图书区域下的人文社科类图书： %s\n", bookRegion["图书"]["人文社科"])
```