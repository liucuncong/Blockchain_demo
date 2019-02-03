package main

import (
	"fmt"
	"time"
)

//func (cli *CLI)AddBlock(data string)()  {
//	//cli.BC.AddBlock(data)
//}

func (cli *CLI)PrintChain()  {
	// 创建迭代器
	bci := cli.BC.NewIterator()

	// 调用迭代器，返回每一区块数据


	for  {
		// 返回区块，左移
		block := bci.Next()


		fmt.Println()
		fmt.Printf("版本号:	%d\n",block.Version)
		fmt.Printf("前区块哈希值:	%x\n",block.PrevHash)
		fmt.Printf("默克尔根:	%x\n",block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp),0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳:	%s\n",timeFormat)
		fmt.Printf("难度值:	%d\n",block.Difficulty)
		fmt.Printf("随机数:	%d\n",block.Nounce)
		fmt.Printf("当前区块哈希值:	%x\n",block.Hash)
		fmt.Printf("区块数据:	%s\n",block.Transactions[0].TXInputs[0].PubKey)


		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束")
			break
		}

	}


}

func (cli *CLI)GetBalance(address string)  {
	// 1.校验地址
	if !IsValidAddress(address) {
		fmt.Println("地址无效:",address)
		return
	}
	//2.生成公钥哈希
	pubKeyHash := GetPubKeyFromAddress(address)

	utxos := cli.BC.FindUTXOs(pubKeyHash)
	total := 0.0
	for _,utxo :=range utxos{
		total += utxo.Value
	}
	fmt.Printf("%s的余额为:%f\n",address,total)
}

func (cli *CLI)Send(from,to string,amount float64,miner,data string)()  {
	// 校验地址
	if !IsValidAddress(from) {
		fmt.Println("地址无效from:",from)
		return
	}
	if !IsValidAddress(to) {
		fmt.Println("地址无效to:",to)
		return
	}
	if !IsValidAddress(miner) {
		fmt.Println("地址无效miner:",miner)
		return
	}
	//1.创建一个挖矿交易
	coinbase := NewCoinbaseTX(miner,data)
	//2.创建一个普通交易
	tx := NewTransaction(from,to,amount,cli.BC)
	if tx == nil {
		return
	}
	//3.添加到区块
	cli.BC.AddBlock([]*Transaction{coinbase,tx})

}

func (cli *CLI)NewWallet()  {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址:%s\n",address)

}

func (cli *CLI)ListAddresses()  {
	ws := NewWallets()

	addresses := ws.ListAllAddresses()
	for _,address := range addresses {
		fmt.Printf("地址:%s\n",address)
	}

}
