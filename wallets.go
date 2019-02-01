package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"crypto/elliptic"
)

// 定义一个Wallets结构，它保存所有的wallet及它的地址
type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建方法
func NewWallets() *Wallets {
	//wallets := LoadFile()
	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.LoadFile()
	return &wallets
}

func (ws *Wallets)CreateWallet() string {
	// 创建钱包和地址
	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet
	ws.SaveToFile()
	return address
}

// 保存方法，把新建的wallet添加进去
func (ws *Wallets)SaveToFile()  {

	var buffer bytes.Buffer

	// panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile("wallet.dat",buffer.Bytes(),0600)
}

// 读取文件方法，把所有的wallet读出来
func (ws *Wallets)LoadFile()  {

	// 1.读取内容
	content,err := ioutil.ReadFile("wallet.dat")
	if err != nil {
		log.Panic(err)
	}
	// 解码
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wallets Wallets  // !!!注意这里不要定义为地址，如果这里为var wallets *Wallets，下面传wallets会报错

	err = decoder.Decode(&wallets)  // &wallets
	if err != nil {
		log.Panic(err)
	}
	// 对于是map的结构来说，要指定赋值，不要在最外层直接赋值
	ws.WalletsMap = wallets.WalletsMap
}


