package BLC

import (
	"fmt"
	"os"
)

// PrintChain
// 打印区块链
func (cli *CLI) PrintChain() {
	var bc *BlockChain = nil
	if DBExists() {
		bc = GetBlockChainObj()
	}
	if bc == nil {
		fmt.Println("没有区块")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.PrintChains()
}

// CreateBlockChain
// 创建区块链
// createblockchain -address hhy
func (cli *CLI) CreateBlockChain(address string) {
	CreateBlockChainWithGenesisBlock(100, address)
}

// Send
// 交易
// send -from hhy -to mhx -amount 20
func (cli *CLI) Send(from, to, amount string) {
	bc := GetBlockChainObj()
	//amountInt,_ := strconv.ParseInt(amount,10,64)
	bc.MineNewBlock([]string{from}, []string{to}, []string{amount})
	defer bc.DB.Close()
}

// GetBalance
// 获取账户余额
func (cli *CLI) GetBalance(address string) {
	bc := GetBlockChainObj()
	balance := bc.GetBalance(address, []*Transaction{})
	fmt.Printf("%s,一共有%d个Token\n", address, balance)
}
