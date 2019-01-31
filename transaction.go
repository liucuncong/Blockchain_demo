package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)
// 挖矿奖励
const REWARD  = 50

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

// 实现一个函数，判断当前交易是否是挖矿交易
func (tx *Transaction)IsCoinbase() bool {
	// 1.交易的input只有一个
	if len(tx.TXInputs) == 1{
		input := tx.TXInputs[0]
		// 2.交易的id为空
		// 3.TXInput的Index为-1
		if bytes.Equal(tx.TXInputs[0].TXid,[]byte{}) && input.Index == -1{
			return true
		}
	}

	return false
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


// 创建普通的转账交易
//1.找到最合理的UTXO集合  map[string][]int64
//2.将这些UTXO逐一转成inputs
//3.创建outputs
//4.如果有零钱，要找零
func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transaction {
	//1.找到最合理的UTXO集合  map[string][]int64
	utxos,resValue := bc.FindNeedUTXOs(from,amount)
	if resValue < amount{
		fmt.Println("余额不足，交易失败！")
		return nil
	}
	//2.将这些UTXO逐一转成inputs
	var inputs []TXInput
	var outputs []TXOutput
	for id,indexArray := range utxos {
		for _,i := range indexArray {
			input := TXInput{[]byte(id),i,from}
			inputs = append(inputs, input)
		}
	}
	//3.创建outputs
	output := TXOutput{amount,to}
	outputs = append(outputs, output)
	//4.如果有零钱，要找零
	if resValue > amount{
		//找零
		outputs = append(outputs, TXOutput{resValue -amount,from})
	}
	tx := Transaction{[]byte{},inputs,outputs}
	tx.SetHash()
	return &tx
}

// 3.创建挖矿交易

// 4.根据交易调整程序
