package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"crypto/ecdsa"
	"crypto/rand"
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
	//Sig string
	// 真正的数字签名，由r,s拼成的[]byte
	Signature []byte
	// 约定Publickey这里不存储原始公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	// 注意：是公钥而不是哈希，也不是地址
	PubKey []byte

}

// 定义交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，我们用地址模拟
	//PubKeyHash string
	//收款方的公钥哈希，注意，是哈希而不是公钥，也不是地址
	PubKeyHash []byte

}

// 由于现在存储的字段是地址的公钥哈希，所以无法直接创建TXOutput
// 为了能够得到公钥哈希，我们需要处理以下，写一个lock函数
func (output *TXOutput)Lock(address string)  {

	// 锁定！！！！
	output.PubKeyHash = GetPubKeyFromAddress(address)
}

//给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64,address string) TXOutput {
	output := TXOutput{
		Value:value,
	}
	output.Lock(address)
	return output

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
	//签名先填写为空，后面创建完整交易后，最后作一次签名即可
	input := TXInput{[]byte{},-1,nil,[]byte(data)}
	//新的创建方法
	output := NewTXOutput(REWARD,address)

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
	// 1.创建交易之后要进行数字签名->所以需要私钥->打开钱包（NewWallets）
	ws := NewWallets()

	// 2.找到自己的钱包，根据地址返回自己的从wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil{
		fmt.Println("没有找到该地址的钱包，交易创建失败")
		return nil
	}

	// 3.得到对应的公私钥
	pubKey := wallet.PubKey
	privateKey := wallet.Privatekey

	// 4.得到公钥哈希
	pubKeyHash := HashRipemd160(pubKey)


	//1.找到最合理的UTXO集合  map[string][]int64
	utxos,resValue := bc.FindNeedUTXOs(pubKeyHash,amount)
	if resValue < amount{
		fmt.Println("余额不足，交易失败！")
		return nil
	}
	//2.将这些UTXO逐一转成inputs
	var inputs []TXInput
	var outputs []TXOutput
	for id,indexArray := range utxos {
		for _,i := range indexArray {
			input := TXInput{[]byte(id),i,nil,pubKey}
			inputs = append(inputs, input)
		}
	}
	//3.创建outputs
	//output := TXOutput{amount,to}
	output := NewTXOutput(amount,to)

	outputs = append(outputs, output)
	//4.如果有零钱，要找零
	if resValue > amount{
		output = NewTXOutput(resValue - amount,from)
		//找零
		outputs = append(outputs, output)
	}
	tx := Transaction{[]byte{},inputs,outputs}
	tx.SetHash()

	bc.SignTransaction(&tx,privateKey)

	return &tx
}

// 签名的具体实现
// 参数：私钥，inputs里面所有引用的交易的结构 map[string]Transaction
func (tx *Transaction)Sign(privateKey *ecdsa.PrivateKey,prevTXs map[string]Transaction)  {
	// 1.创建一个当前交易的副本:txCopy,使用函数:TrimmedCopy:要把Signature和PubKey字段设置为nil
	txCopy := tx.TrimmedCopy()
	// 2.循环遍历txCopy的inputs，得到这个input所引用的交易的output的公钥哈希
	for i,input := range txCopy.TXInputs {

		prevTx := prevTXs[string(input.TXid)]
		if len(prevTx.TXID) == 0{
			log.Panic("引用的交易无效")
		}
		// 3.生成要签名的数据。要签名的数据一定是哈希值
		//a.我们对每一个input都要签名一次，签名数据是由当前input引用的交易的output的公钥哈希+当前的outputs(都承载在当前这个txCopy里面)
		//b.要对这个拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名的最终数据

		// 不要对input进行赋值，他只是一个副本，要对txCopy.TXInputs[xx]进行操作，要不然无法把pubKeyHash传进来
		txCopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash
		// 所需要的三个数据都具备了，开始做哈希处理
		txCopy.SetHash()
		signDataHash := txCopy.TXID
		// 还原，以免影响后面input的签名
		txCopy.TXInputs[i].PubKey = nil

		// 4.执行签名动作，得到r,s字节流
		r,s,err :=ecdsa.Sign(rand.Reader,privateKey,signDataHash)
		if err != nil {
			log.Panic(err)
		}

		// 5.放到我们所签名的input的Signature中
		signature := append(r.Bytes(),s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}
}

func (tx *Transaction)TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _,input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXid,input.Index,nil,nil})
	}
	for _,output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}
	return Transaction{tx.TXID,inputs,outputs}
}



// 3.创建挖矿交易

// 4.根据交易调整程序
