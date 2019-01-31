package main

import "fmt"

func (cli *CLI)AddBlock(data string)()  {
	//cli.BC.AddBlock(data)
}

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
		fmt.Printf("时间戳:	%d\n",block.TimeStamp)
		fmt.Printf("难度值:	%d\n",block.Difficulty)
		fmt.Printf("随机数:	%d\n",block.Nounce)
		fmt.Printf("当前区块哈希值:	%x\n",block.Hash)
		fmt.Printf("区块数据:	%s\n",block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束")
			break
		}
	}
}

func (cli *CLI)GetBalance(address string)  {
	utxos := cli.BC.FindUTXOs(address)
	total := 0.0
	for _,utxo :=range utxos{
		total += utxo.Value
	}
	fmt.Printf("%s的余额为:%f\n",address,total)
}

func (cli *CLI)Send(from,to string,amount float64,miner,data string)()  {

	fmt.Println(amount)
}
