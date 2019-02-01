package main

// 定义一个Wallets结构，它保存所有的wallet及它的地址
type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建方法
func NewWallets() *Wallets {
	wallet := NewWallet()
	address := wallet.NewAddress()

	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.WalletsMap[address] = wallet
	return &wallets
}

// 保存方法，把新建的wallet添加进去


// 读取文件方法，把所有的wallet读出来