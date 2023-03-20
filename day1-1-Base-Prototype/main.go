package main

import (
	"fmt"
	"myPublicBlockChain/day1-1-Base-Prototype/BLC"
)

func main() {
	/*  part.1 测试生成区块功能  */
	//fmt.Println("Public BlockChain!")
	//block := BLC.CreateGenesisBlock("hello chain")
	//fmt.Printf("Heigth:%x\n", block.Height)
	//fmt.Printf("Data:%s\n", block.Data)

	/*  part.2 测试区块链功能  */
	blockChain := BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	blockChain.AddBlockToBlockChain("Send 1BTC To Wangergou", blockChain.Blocks[len(blockChain.Blocks)-1].Height+1, blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 3BTC To lixiaohua", blockChain.Blocks[len(blockChain.Blocks)-1].Height+1, blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 5BTC To rose", blockChain.Blocks[len(blockChain.Blocks)-1].Height+1, blockChain.Blocks[len(blockChain.Blocks)-1].Hash)

	for _, block := range blockChain.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
