package main

import (
	"time"
	"bytes"
	"encoding/binary"
	"log"
)

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// 0. 定义结构
type Block struct {
	// 1.版本号
	Version uint64
	// 2.前区块哈西
	PrevHash []byte
	// 3.Merkel根（这就是一个哈希值，这里先不管，v4再介绍）
	MerkelRoot []byte
	// 4.时间戳
	TimeStamp uint64
	// 5.难度值
	Difficulty uint64
	// 6.随机数
	Nounce uint64

	// a.当前区块哈希
	Hash []byte
	// b.数据
	Data []byte
}

// 实现一个辅助函数，将uint64转成[]byte
func Uint64ToByte(data uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,data)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// 1.创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:00,
		PrevHash: prevBlockHash,
		MerkelRoot:[]byte{},
		TimeStamp:uint64(time.Now().Unix()),
		Difficulty:0,
		Nounce:0,
		Hash:     []byte{}, // 先填空，后面再计算
		Data:     []byte(data),
	}
	// 创建一个pow对象
	pow := NewProofOfWork(&block)
	hash,nounce := pow.Run()
	block.Hash = hash
	block.Nounce = nounce
	return &block
}
/*
// 2.生成哈希
func (block *Block) SetHash() {
	var blockInfo []byte
	// 1.拼装数据

	//blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	//blockInfo = append(blockInfo, block.PrevHash...)
	//blockInfo = append(block.PrevHash, block.MerkelRoot...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.Nounce)...)
	//blockInfo = append(blockInfo, block.Data...)

	tem := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nounce),
		block.Data,
	}
	// 将二维的切片数组连接起来，返回一个一维的切片
	blockInfo = bytes.Join(tem,[]byte{})

	// 2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
*/

// 5.定义一个创始块
func GenesisBlock(data string, prevBlockHash []byte) *Block {
	return NewBlock(data, prevBlockHash)
}
