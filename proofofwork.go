package main

import "math/big"

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
	targertStr := "001000000000000000000"
	// 引入辅助变量，目的是将上面的难度值转成big.int
	tem := big.Int{}
	// 将难度值赋值给big.Int，指定16进制的格式
	tem.SetString(targertStr,16)

	pow.targert = &tem
	return &pow
}