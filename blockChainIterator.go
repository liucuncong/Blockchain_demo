package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	DB *bolt.DB
	// 游标，用于不断索引
	CurrentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.DB,
		// 最初指向区块链的最后一个区块，随着Next()的调用，不断变化
		bc.LastBlockHash,
	}
}

func (BCI *BlockChainIterator)Next() *Block {
	// 获取当前游标
	currentHash := BCI.CurrentHashPointer
	var block Block
	// 查询表
	err := BCI.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if bucket == nil {
			log.Panic("迭代器遍历时bucket不应该为空，请检查")
		}
		// 获取区块
		blockBytes := bucket.Get(currentHash)
		block = Deserialize(blockBytes)
		// 游标哈希左移
		BCI.CurrentHashPointer = block.PrevHash
		return nil
	})
	if err != nil {
		log.Panic("数据库查询操作错误")
	}
	return &block
}