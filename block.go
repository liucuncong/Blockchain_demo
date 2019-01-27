package main

import "crypto/sha256"

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

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

// 5.定义一个创始块
func GenesisBlock(data string,prevBlockHash []byte) *Block {
	return NewBlock(data,prevBlockHash)
}
