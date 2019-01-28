package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

// 定义工作量证明结构体
type ProofOfWork struct {
	//a.block
	block *Block
	//b.目标值
	//一个非常大的数，它有丰富的方法:比较，赋值
	targert *big.Int
}

// 创建POW对象
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block:block,
	}
	// 我们指定的难度值现在是一个string类型，需要进行转换
	targertStr := "0000100000000000000000000000000000000000000000000000000000000000"
	// 引入辅助变量，目的是将上面的难度值转成big.int
	tem := big.Int{}
	// 将难度值赋值给big.Int，指定16进制的格式
	tem.SetString(targertStr,16)

	pow.targert = &tem
	return &pow
}

// 计算哈希
func (pow *ProofOfWork)Run() ([]byte,uint64) {
	var nounce uint64
	block := pow.block
	var hash [32]byte

	for {
		//1.拼装数据（区块的数据，还有不断变化的随机数）
		tem := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nounce),
			block.Data,
		}
		// 将二维的切片数组连接起来，返回一个一维的切片
		blockInfo := bytes.Join(tem,[]byte{})
		//2.做哈希运算
		hash = sha256.Sum256(blockInfo)
		//3.将哈希转成big.Int类型，与pow中的target进行比较
		temInt := big.Int{}
		temInt.SetBytes(hash[:])
		if temInt.Cmp(pow.targert) == -1 {
			//a.找到了，退出返回
			fmt.Printf("挖矿成功! hash:%x;nounce:%d\n",hash,nounce)
			break
		}
		//b.没找到，继续找，随机数+1
		nounce++
	}


	return hash[:],nounce
}

// 提供一个校验哈希的函数
