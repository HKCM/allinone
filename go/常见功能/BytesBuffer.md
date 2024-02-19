在string的拼接时使用Builder，byte拼接时使用Buffer
```go
// bytes.Buffer的0值可以直接使用
var buff bytes.Buffer

// 向buff中写入字符/字符串
buff.Write([]byte("Hello"))
buff.WriteByte(' ')
buff.WriteString("World")

// String() 方法获得拼接的字符串
buff.String() // "Hello World"
```