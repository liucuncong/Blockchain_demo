package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	)


// 这里的钱包是一个结构，保存了公钥、私钥对
type Wallet struct {
	// 私钥
	Privatekey *ecdsa.PrivateKey
	//Publickey *ecdsa.PublicKey

	// 约定Publickey这里不存储原始公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	PubKey []byte
}

//创建钱包
func NewWallet() *Wallet {
	//生成私钥
	privatekey,err := ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 生成公钥
	publickey := privatekey.PublicKey
	// 拼接X,Y
	pubKey := append(publickey.X.Bytes(),publickey.Y.Bytes()...)

	return &Wallet{privatekey,pubKey}
}

//生成地址
func (w *Wallet)NewAddress() string {
	pubKey := w.PubKey
	// 返回rip160的结果
	rip160hash := HashRipemd160(pubKey)
	// 拼接版本前缀
	version := byte(00)
	payload := append([]byte{version}, rip160hash...)

	// 生成校验玛
	checkcode := CkeckSum(payload)

	// 25字节数据
	payload = append(payload,checkcode...)

	address := base58.Encode(payload)
	return address
}


func HashRipemd160(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)

	rip160hasher := ripemd160.New()
	_,err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	// 返回rip160的结果
	rip160hash := rip160hasher.Sum(nil)
	return rip160hash
}

func CkeckSum(data []byte) []byte {

	//ckecksum
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	// 后4字节校验玛
	checkcode := hash2[:4]
	return checkcode
}