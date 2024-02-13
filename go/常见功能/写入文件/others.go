package main

import (

	// "io/ioutil" //缓存IO

	"bufio"
	"fmt"
	"os"
)

func main() {

	var filename = "./output1.txt"
	var err error

	/***************************** 第一种方式: 使用 io.WriteString 写入文件 ***********************************************/
	// 可以追加
	// var f *os.File
	// var msgString = "第一种方式: 使用 io.WriteString 写入文件\n"
	// if checkFileIsExist(filename) { //如果文件存在
	// 	f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
	// 	fmt.Println("文件存在")
	// } else {
	// 	f, err = os.Create(filename) //创建文件
	// 	fmt.Println("文件不存在")
	// }
	// check(err)
	// n, err := io.WriteString(f, msgString) //写入文件(字符串)
	// check(err)
	// fmt.Printf("写入 %d 个字节n", n)

	/*****************************  第二种方式: 使用 ioutil.WriteFile 写入文件 ***********************************************/
	// 每次都直接覆盖
	// msgString := "第二种方式: 使用 ioutil.WriteFile 写入文件"
	// var msgBytes = []byte(msgString)
	// err = ioutil.WriteFile(filename, msgBytes, 0666) //写入文件(字节数组)
	// check(err)

	/*****************************  第三种方式: 使用 File(Write,WriteString) 写入文件 ***********************************************/
	// 可以追加
	// var f *os.File
	// var msgString = "第一种方式: 使用 io.WriteString 写入文件\n"
	// if checkFileIsExist(filename) { //如果文件存在
	// 	f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
	// 	fmt.Println("文件存在")
	// } else {
	// 	f, err = os.Create(filename) //创建文件
	// 	fmt.Println("文件不存在")
	// }
	// check(err)
	// defer f.Close()
	// msgString1 := "第三种方式: 使用 File Write 写入文件\n"
	// var msgBytes = []byte(msgString1)
	// n2, err := f.Write(msgBytes) //写入文件(字节数组)
	// check(err)
	// fmt.Printf("写入 %d 个字节n", n2)
	// msgString2 := "第三种方式: 使用 File WriteString 写入文件\n"
	// n3, err := f.WriteString(msgString2) //写入文件(字节数组)
	// fmt.Printf("写入 %d 个字节n", n3)
	// check(err)
	// f.Sync()

	/***************************** 第四种方式: 使用 bufio.NewWriter 写入文件 ***********************************************/
	//可以追加
	var f *os.File
	var msgString = "第四种方式: 使用 bufio.NewWriter 写入文件\n"
	if checkFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	defer f.Close()
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	n4, err := w.WriteString(msgString)
	check(err)
	fmt.Printf("写入 %d 个字节n", n4)
	w.Flush()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
