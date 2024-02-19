在string的拼接时使用Builder，byte拼接时使用Buffer
```go
// strings.Builder的0值可以直接使用
var builder strings.Builder

// 向builder中写入字符/字符串
builder.Write([]byte("Hello"))
builder.WriteByte(' ')
builder.WriteString("World")

// String() 方法获得拼接的字符串
builder.String() // "Hello World"
```