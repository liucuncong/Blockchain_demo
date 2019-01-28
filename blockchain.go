package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

// 3.引入区块连
type BlockChain struct {
	// 定义一个区块链数组
	//blocks []*Block

	db *bolt.DB
	lastBlockHash []byte  //最后一个区块的哈希
}

const blockChainDB  = "blockchain.db"
const blockBucket  = "blockBucket"

// 4.定义一个区块链
func NewBlockChain() *BlockChain {
	var lastBlockHash []byte
	// 1.打开数据库
	db,err := bolt.Open(blockChainDB,0600,nil)
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	defer db.Close()
	// 写入
	db.Update(func(tx *bolt.Tx) error {
		// 2.找到抽屉bucket（如果没有就创建）
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			bucket,err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket失败！")
			}
			// 创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(genesisInfo,[]byte{})
			// 3.写数据
			// hash作为key,block的字节流作为value来实现
			bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
			bucket.Put([]byte("lastHashKey"),genesisBlock.Hash)
			lastBlockHash = genesisBlock.Hash
			blockBytes := bucket.Get(genesisBlock.Hash)
			block := Deserialize(blockBytes)
			fmt.Printf("Deserialize：%v\n",block)
		} else {
			lastBlockHash = bucket.Get([]byte("lastHashKey"))
		}

		return nil
	})
	return &BlockChain{db,lastBlockHash}
}

/*
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
*/