package main

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"
	"os"
)

/*
区块结构
 */
type Block struct {
	Timestamp     int64  //时间戳
	Data          []byte //当前区块 存放的信息
	PrevBlockHash []byte // 上一个区块的 加密的hash
	Hash          []byte //当前的区块的Hash
}

/*
Block 这个结构体 绑定的一个方法
 */

func (this *Block) SetHash() {
	//将本区块的Timestamp+Data+PrevBlockHash ----> Hash
	//将时间戳由整型  ---> 二进制
	timestamp := []byte(strconv.FormatInt(this.Timestamp, 10))
	//将三个二进制的属性进行拼接
	headers := bytes.Join([][]byte{timestamp, this.Data, this.PrevBlockHash}, []byte{})
	//将拼接之后的headers 进行 SHA256加密
	hash := sha256.Sum256(headers)
	this.Hash = hash[:]
}

/*
create block
 */
func NewBlock(data string, prevBlockHash []byte) *Block {
	//生成一个区块
	block := Block{}
	//给当前的区块 赋值 （创建时间，data 前区块hash）
	block.Timestamp = time.Now().Unix()
	block.Data = []byte(data)
	block.PrevBlockHash = prevBlockHash
	//给当前区块进行hash加密
	block.SetHash()
	//将已经复制好的区块 返回给外部
	return &block
}

/*
区块链结构
 */
type BlockChain struct {
	Blocks []*Block
}

//创世块
func NewGenesisBlock() *Block {
	//方法一
	//genesisBlock := Block{}
	//genesisBlock.Data = []byte("VickyWang Blockchain Demo")
	//genesisBlock.PrevBlockHash = []byte{}
	return NewBlock("VickyWang Genesis Block", []byte{})
}

//新建一个区块链
func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

//将区块添加一个区块链中

func (this *BlockChain) AddBlock(data string) {
	//得到一个新添加区块的前区块Hash
	prevBlock := this.Blocks[len(this.Blocks)-1]
	//根据data 创建一个新的区块
	newBlock := NewBlock(data, prevBlock.Hash)
	//依据前区块和新区块 添加到区块链blocks中
	this.Blocks = append(this.Blocks, newBlock)
}

func main() {
	//fmt.Println("Blockchain Demo Test start")
	//block := NewBlock("VickyWang", []byte{})
	//fmt.Printf("block.hash = %x\n", block.Hash)
	//fmt.Println("Blockchain Demo Test end")

	//创建一个区块链 bc
	bc := NewBlockchain()
	//用户输入的指令 1，2， other
	var cmd string
	for {
		fmt.Println(" 按 '1' 添加一条信息数据 到区块链中")
		fmt.Println(" 按 '2' 遍历当前的区块链都有那些区块链信息")
		fmt.Println(" 按 其他按键退出")
		fmt.Scanf("%s", &cmd)
		switch cmd {
		case "1":
			input := make([]byte, 1024)
			//添加一个区块
			fmt.Println("请输入区块链的行为数据(要添加保存的数据)：")
			os.Stdin.Read(input)
			bc.AddBlock(string(input))
		case "2":
			//遍历整个区块链
			for i, block := range bc.Blocks {
				fmt.Println("===============================")
				fmt.Println("第 ", i, "个 区块的信息：")
				fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
				fmt.Printf("Data：%s\n", block.Data)
				fmt.Printf("Hash：%x\n", block.Hash)
				fmt.Println("===============================")
			}
		default:
			//退出程序
			fmt.Println("Exit VickyWang Blockchain Demo")
			return
		}
	}
}
