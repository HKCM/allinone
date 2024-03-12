package common

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	data = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890,.-=/"
)

func init() {
	rand.Seed(time.Now().Unix()) // 设置随机种子
}

// GenRandString 生成n个随机字符的string
func GenRandString(n int) string {
	max := len(data)
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteByte(data[rand.Intn(max)])
	}

	return buf.String()
}

// GenRandBytes 生成n个随机字符的[]byte
func GenRandBytes(n int) []byte {
	max := len(data)
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = data[rand.Intn(max)]
	}

	return buf
}

// Example
// msg := "123\n"
// filename := "test.txt"
// file := GetFile(filename)
// defer file.Close()
// writer := bufio.NewWriter(file) //创建新的 Writer 对象
//
//	for i := 0; i < 2000; i++ {
//		WriteToFile(writer, msg)
//	}
//
// writer.Flush()
func GetFile(filename string) *os.File {
	var err error
	var file *os.File

	if CheckFileIsExist(filename) { //如果文件存在
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		file, err = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	if err != nil {
		panic(err)
	}

	return file
}

// bufio.NewScanner NewScanner 按行读取
func readFileByLine(fileName string, callback func(s string)) {
	file, e := os.Open(fileName)
	if e != nil {
		fmt.Println("打开文件失败")
		return
	}
	defer file.Close()
	// read the file line by line using scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		callback(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func printLine(line string) {
	println(line)
}

func WriteToFile(writer *bufio.Writer, content string) {
	_, err := writer.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func CheckFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
