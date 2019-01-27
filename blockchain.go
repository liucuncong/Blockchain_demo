package main


// 3.引入区块连
type BlockChain struct {
	// 定义一个区块链数组
	blocks []*Block
}

// 4.定义一个区块链
func NewBlockChain() *BlockChain {
	// 创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock(genesisInfo,[]byte{})
	return &BlockChain{
		blocks:[]*Block{genesisBlock},
	}
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
