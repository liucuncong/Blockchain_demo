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
func NewBlockChain(address string) *BlockChain {
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
			genesisBlock := GenesisBlock(address,[]byte{})
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
func (bc *BlockChain)AddBlock(txs []*Transaction)  {
	// 1.获取最后一个区块的哈希
	lastBlockHash := bc.LastBlockHash
	db:= bc.DB

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if bucket == nil{
			log.Panic("打开bucket失败！")
		}
		// 2.创建新的区块
		block := NewBlock(txs,lastBlockHash)
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

// 找到指定地址的所有utxo
func (bc *BlockChain)FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput

	// 4.遍历input，找到和自己花费的utxo
	// 我们定义一个map来保存消费过的output,key是这个output所在交易的id，value是这个交易中索引的数组
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	it := bc.NewIterator()
	for {
		// 1.遍历区块
		block := it.Next()

		// 2.遍历交易
		for _,tx:=range block.Transactions {
			// 3.遍历output，找到和自己相关的utxo（在添加output之前，检查一下是否已经消耗过）
		OUTPUT:
			for i,output := range tx.TXOutputs {
				// 在这里做一个过滤，将所有消耗过的outputs和当前所即将添加的output对比一下
				// 如果相同，则跳过，否则添加
				// 如果当前交易的id存在于我们已经表示的map，那么说明这个交易里面有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil {
					for _,j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//当前准备添加的output已经消耗过了，不要再加了
							continue OUTPUT
						}
					}
				}

				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}
			// 如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if !tx.IsCoinbase() {
				// 4.遍历input，找到和自己花费的utxo
				for _,input := range tx.TXInputs {
					if input.Sig == address {
						indexArray := spentOutputs[string(input.TXid)]  // 这里取值相当于返回一个空[]int64数组，所以下面可以用append
						indexArray = append(indexArray,input.Index)
					}
				}
			}

		}

		if len(block.PrevHash) == 0{
			break
		}
	}

	return UTXO
}

//
func (bc *BlockChain)FindNeedUTXOs(from string,amount float64) (map[string][]int64,float64) {
	// 找到的合理的utxos集合
	var utxos = make(map[string][]int64)
	// 找到utxos里面包含钱的总数
	var calc float64
	// 标识已经消耗过的utxo
	spentOutputs := make(map[string][]int64)

	// 111111111111111111111111111
	it := bc.NewIterator()
	for {
		// 1.遍历区块
		block := it.Next()

		// 2.遍历交易
		for _,tx:=range block.Transactions {
			// 3.遍历output，找到和自己相关的utxo（在添加output之前，检查一下是否已经消耗过）
		OUTPUT:
			for i,output := range tx.TXOutputs {
				// 在这里做一个过滤，将所有消耗过的outputs和当前所即将添加的output对比一下
				// 如果相同，则跳过，否则添加
				// 如果当前交易的id存在于我们已经表示的map，那么说明这个交易里面有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil {
					for _,j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//当前准备添加的output已经消耗过了，不要再加了
							continue OUTPUT
						}
					}
				}

				if output.PubKeyHash == from {
					//UTXO = append(UTXO, output)
					// 我们要实现的逻辑就在这里，找到自己需要的最少的utxo

					//0.比较一下是否满足转账需求
					if  calc < amount{
						//1.把utxo加进来
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], int64(i) )
						//2.统计一下当前utxo的总额
						calc += output.Value
						// 比较一下是否满足转账需求
						if calc >= amount{
							return utxos,calc
						}

					}


					//   a.满足的话，直接返回  utxos,calc
					//   b.不满足继续统计




				}
			}
			// 如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if !tx.IsCoinbase() {
				// 4.遍历input，找到和自己花费的utxo
				for _,input := range tx.TXInputs {
					if input.Sig == from {
						indexArray := spentOutputs[string(input.TXid)]  // 这里取值相当于返回一个空[]int64数组，所以下面可以用append
						indexArray = append(indexArray,input.Index)
					}
				}
			}

		}

		if len(block.PrevHash) == 0{
			break
		}
	}



	// 222222222222222222222222222


	return utxos,calc
}