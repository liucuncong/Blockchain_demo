package main

import "fmt"

// 0. 定义结构
type Block struct {
	// 前区块哈西
	PrevHash []byte
	// 当前区块哈希
	Hash []byte
	// 数据
	Data []byte
}

// 1.创建区块
func NewBlock(data string,prevBlockHash []byte) *Block {
	block := Block{
		PrevHash:prevBlockHash,
		Hash:[]byte{},  // 先填空，后面再计算  // TODO
		Data:[]byte(data),
	}
	return &block
}


func main()  {
	block := NewBlock("我转给你10000个比特币",[]byte{})
	fmt.Printf("前区块哈希值:%x\n",block.PrevHash)
	fmt.Printf("前区块哈希值:%x\n",block.Hash)
	fmt.Printf("前区块哈希值:%s\n",block.Data)
}
