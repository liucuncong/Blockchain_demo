package main

import (
	"fmt"
	"crypto/sha256"
)

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
		Hash:[]byte{},  // 先填空，后面再计算
		Data:[]byte(data),
	}
	block.SetHash()
	return &block
}

// 2.生成哈希
func (block *Block)SetHash()  {
	// 1.拼装数据
	blockInfo := append(block.PrevHash,block.Data...)
	// 2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

func main()  {
	block := NewBlock("我转给你10000个比特币",[]byte{})
	fmt.Printf("前区块哈希值:%x\n",block.PrevHash)
	fmt.Printf("当前区块哈希值:%x\n",block.Hash)
	fmt.Printf("区块数据:%s\n",block.Data)
}
