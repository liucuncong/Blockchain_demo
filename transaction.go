package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)
// 挖矿奖励
const REWARD  = 12.5

// 1.定义交易结构
type Transaction struct {
	TXID      []byte     //交易ID
	TXInputs  []TXInput  //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}

// 定义交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本，我们用地址来模拟
	Sig string
}

// 定义交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，我们用地址模拟
	PubKeyHash string
}

// 设置交易ID
func (tx *Transaction)SetHash()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic()
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 2.提供创建交易方法
func NewCoinbaseTX(address string,data string) *Transaction {
	// 挖矿交易的特点
	//1.只有一个input
	//2.无需引用交易ID
	//3.无需引用index
	//矿工由于挖矿时，无需指定签名，所以这个sin字段可以由矿工自由填写，一般填写矿池的名字
	input := TXInput{[]byte{},-1,data}
	output := TXOutput{REWARD,address}

	// 对于挖矿交易只有一个input和output
	tx := Transaction{[]byte{},[]TXInput{input},[]TXOutput{output}}
	tx.SetHash()
	return &tx
}

// 3.创建挖矿交易

// 4.根据交易调整程序
