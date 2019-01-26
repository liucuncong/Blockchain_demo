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

// 3.引入区块连
type BlockChain struct {
	// 定义一个区块链数组
	blocks []*Block
}

// 4.定义一个区块链
func NewBlockChain() *BlockChain {
	// 创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks:[]*Block{genesisBlock},
	}
}

// 5.定义一个创始块
func GenesisBlock() *Block {
	return NewBlock("我给你转了1000个比特币!",[]byte{})
}

// 6.添加区块
func (bc *BlockChain)AddBlock(data string)  {
	// 1.获取最后一个区块的哈希
	lastBlock := bc.blocks[len(bc.blocks)-1]
	preHash := lastBlock.Hash

	// 2.创建新的区块
	block := NewBlock(data,preHash)
	// 3.添加区块到区块链数组中
	bc.blocks = append(bc.blocks, block)
}

func main()  {
	bc := NewBlockChain()
	bc.AddBlock("詹姆斯给我转了100个比特币!")
	bc.AddBlock("科比给我转了100个比特币!")

	for i,block := range bc.blocks{
		fmt.Printf("=======当前区块高度:	%d========\n",i)
		fmt.Printf("前区块哈希值:	%x\n",block.PrevHash)
		fmt.Printf("当前区块哈希值:	%x\n",block.Hash)
		fmt.Printf("区块数据:	%s\n",block.Data)
	}

}
