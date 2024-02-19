keyword: stringBuilder byteBuffer bytesBuffer
```go
// bytes.Buffer的0值可以直接使用
var buff bytes.Buffer

// 向buff中写入字符/字符串
buff.Write([]byte("Hello"))
buff.WriteByte(' ')
buff.WriteString("World")
// String() 方法获得拼接的字符串
println(buff.String())  // Hello World
buff.Reset() // 清空缓冲期
// buff.Truncate(0)
buff.WriteString("New World")
println(buff.String()) // New World
```