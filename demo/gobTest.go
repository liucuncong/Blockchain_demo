package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

// 定义一个结构Person
type Person struct {
	Name string
	Age uint
}


func main()  {
	var xiaoming Person
	xiaoming.Name = "小名"
	xiaoming.Age = 3
	// 编码的数据放到buffer
	var buffer bytes.Buffer
	// 使用gob进行序列化（编码）得到字节流
	// 1.定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	// 2.使用编码器进行编码
	err := encoder.Encode(&xiaoming)
	if err != nil {
		log.Panic("编码错误")
	}
	fmt.Printf("编码后的xiaoming：%v\n",buffer.Bytes())
	// 3.定义一个解码器
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	// 4.使用解码器进行解码
	var daming Person
	err = decoder.Decode(&daming)
	if err != nil {
		log.Panic("编码错误")
	}
	fmt.Printf("解码后的xiaoming：%v\n",daming)

}
