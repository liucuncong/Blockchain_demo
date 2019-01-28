package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

func main()  {
	// 1.打开数据库
	db,err := bolt.Open("test.db",0600,nil)
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	defer db.Close()
	// 写入
	db.Update(func(tx *bolt.Tx) error {
		// 2.找到抽屉bucket（如果没有就创建）
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			bucket,err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic("创建bucket失败！")
			}
		}
		// 3.写数据
		bucket.Put([]byte("111"),[]byte("aaa"))
		return nil
	})

	// 读取
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			log.Panic("bucket为空！")
		}
		v1 := bucket.Get([]byte("111"))
		fmt.Printf("v1：%s\n",v1)
		return nil
	})


}


