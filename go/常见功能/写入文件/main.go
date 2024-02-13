package main

import (
	"bufio"
	"strconv"

	"log"
	"os"
)

func main() {

	filename := "test.txt"
	var file *os.File
	var err error

	// filename 也可以按时间每次创建新的
	if checkFileIsExist(filename) { //如果文件存在
		// file, err = os.Create(filename) //也可以使用这个每次都创建文件
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
		log.Printf("目标文件 %s 已存在,追加...\n", filename)
	} else {
		file, err = os.Create(filename) //创建文件
		log.Printf("目标文件  %s  不存在,创建文件并写入\n", filename)
	}

	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 1; i <= 8001; i++ {
		write(writer, strconv.Itoa(i))
	}
	writer.Flush()
}

func write(writer *bufio.Writer, msg string) {
	_, err := writer.WriteString(msg + "\n")
	if err != nil {
		log.Panic(err)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
