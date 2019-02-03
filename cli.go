package main

import (
	"os"
	"fmt"
	"strconv"
)

// 这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	BC *BlockChain
}

const USAGE = `
	printChain	        "正向打印区块链"
	getBalance --address ADDRESS  "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA  "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
	newWallet  "创建一个新的钱包（私钥公钥对）"
	listAddresses  "列举所有的钱包地址"
`

// 接收参数的动作，我们放到一个函数中
func (cli *CLI) Run() {
	// 1.得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Println(USAGE)
		return
	}
	// 2.分析命令
	cmd := args[1]
	switch cmd {
	//case "addBlock":
	//	// 确保命令有效
	//	if len(args) == 4 && args[2] == "--data" {
	//		// 获取数据
	//		data := args[3]
	//		cli.AddBlock(data)
	//	} else {
	//		fmt.Println("添加区块参数错误，清检查")
	//		fmt.Println(USAGE)
	//	}
	case "printChain":
		// 打印区块
		cli.PrintChain()
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			// 获取数据
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Println("转账开始")
		if len(args) != 7{
			fmt.Println("参数个数错误，请检查")
			fmt.Println(USAGE)
			return
		}
		from := args[2]
		to := args[3]
		amount,err:= strconv.ParseFloat(args[4],64)  //字符窜转float64
		if err != nil {
			fmt.Println("转账金额格式错误，请检查")
			return
		}
		miner:= args[5]
		data:= args[6]

		cli.Send(from,to,amount,miner,data)
	case "newWallet":
		fmt.Println("创建新的钱包")
		cli.NewWallet()
	case "listAddresses":
		fmt.Println("列举所有的地址。。。。")
		cli.ListAddresses()
	default:
		fmt.Println("无效的命令，请检查")
		fmt.Println(USAGE)
	}
}

