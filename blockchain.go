package main

import (
	"github.com/boltdb/bolt"
	"log"
	)

// 3.引入区块连
type BlockChain struct {
	// 定义一个区块链数组
	//blocks []*Block

	DB *bolt.DB
	LastBlockHash []byte  //最后一个区块的哈希
}

const DBNAME  = "blockchain.db"
const BLOCKBUCKETNAME = "blockBucket"

// 4.定义一个区块链
func NewBlockChain() *BlockChain {
	var lastBlockHash []byte
	// 1.打开数据库
	db,err := bolt.Open(DBNAME,0600,nil)
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	// 写入
	err = db.Update(func(tx *bolt.Tx) error {
		// 2.找到抽屉bucket（如果没有就创建）
		bucket := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if bucket == nil{
			bucket,err = tx.CreateBucket([]byte(BLOCKBUCKETNAME))
			if err != nil {
				log.Panic("创建bucket失败！")
			}
			// 创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(GENESISINFO,[]byte{})
			// 3.写数据
			// hash作为key,block的字节流作为value来实现
			bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
			bucket.Put([]byte("lastHashKey"),genesisBlock.Hash)
			lastBlockHash = genesisBlock.Hash

			// 测试用的，马上删掉
			//blockBytes := bucket.Get(genesisBlock.Hash)
			//block := Deserialize(blockBytes)
			//fmt.Printf("Deserialize：%v\n",block)
		} else {
			lastBlockHash = bucket.Get([]byte("lastHashKey"))
		}

		return nil
	})
	if err != nil {
		log.Panic("数据库更新操作失败！")
	}
	return &BlockChain{db,lastBlockHash}
}


// 6.添加区块
func (bc *BlockChain)AddBlock(data string)  {
	// 1.获取最后一个区块的哈希
	lastBlockHash := bc.LastBlockHash
	db:= bc.DB

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if bucket == nil{
			log.Panic("打开bucket失败！")
		}
		// 2.创建新的区块
		block := NewBlock(data,lastBlockHash)
		// 3.添加区块到区块链数中
		err := bucket.Put(block.Hash,block.Serialize())
		if err != nil {
			log.Panic("数据库添加区块失败！")
		}
		err = bucket.Put([]byte("lastHashKey"),block.Hash)
		if err != nil {
			log.Panic("数据库添加最新区块哈希失败！")
		}
		// 4.更新一下内存中的区块链，指的是把lastBlockHash更新一下
		bc.LastBlockHash = block.Hash

		return nil
	})
	if err != nil {
		log.Panic("数据库更新操作失败！")
	}
}
