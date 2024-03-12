keyword: stringBuilder byteBuffer bytesBuffer

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


```go
// BenchmarkSpliceAddString10 测试使用 += 拼接N次长度为10的字符串
func BenchmarkSpliceAddString10(b *testing.B) {
    s := ""
    for i := 0; i < b.N; i++ {
        s += GenRandString(10)
    }
}

// BenchmarkSpliceBuilderString10 测试使用strings.Builder拼接N次长度为10的字符串
func BenchmarkSpliceBuilderString10(b *testing.B) {
    var builder strings.Builder
    for i := 0; i < b.N; i++ {
        builder.WriteString(GenRandString(10))
    }
}

// BenchmarkSpliceBufferString10 测试使用bytes.Buffer拼接N次长度为10的字符串
func BenchmarkSpliceBufferString10(b *testing.B) {
    var buff bytes.Buffer
    for i := 0; i < b.N; i++ {
        buff.WriteString(GenRandString(10))
    }
}

// 性能最好
// BenchmarkSpliceBufferByte10 测试使用bytes.Buffer拼接N次长度为10的[]byte
func BenchmarkSpliceBufferByte10(b *testing.B) {
    var buff bytes.Buffer
    for i := 0; i < b.N; i++ {
        buff.Write(GenRandBytes(10))
    }
}

// BenchmarkSpliceBuilderByte10 测试使用string.Builder拼接N次长度为10的[]byte
func BenchmarkSpliceBuilderByte10(b *testing.B) {
    var builder strings.Builder
    for i := 0; i < b.N; i++ {
        builder.Write(GenRandBytes(10))
    }
}
```

参考: https://www.cnblogs.com/apocelipes/p/9283841.html