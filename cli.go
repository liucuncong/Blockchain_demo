package main

import (
	"os"
	"fmt"
)

// 这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	BC *BlockChain
}

const USAGE = `
	addBlock --data DATA    "添加区块"
	printChain	        "正向打印区块链"
	getBalance --address ADDRESS  "获取指定地址的余额"
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
	case "addBlock":
		// 确保命令有效
		if len(args) == 4 && args[2] == "--data" {
			// 获取数据
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Println("添加区块参数错误，清检查")
			fmt.Println(USAGE)
		}
	case "printChain":
		// 打印区块
		cli.PrintChain()
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			// 获取数据
			address := args[3]
			cli.GetBalance(address)
		}
	default:
		fmt.Println("无效的命令，请检查")
		fmt.Println(USAGE)
	}
}

