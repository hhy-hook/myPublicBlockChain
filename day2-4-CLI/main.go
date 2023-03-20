package main

import "myPublicBlockChain/day2-4-CLI/BLC"

func main() {

	//blockchain := BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	//fmt.Println(blockchain)
	//defer blockchain.DB.Close()
	////8.测试新添加的区块
	//blockchain.AddBlockToBlockChain("Send 100RMB to wangergou")
	//blockchain.AddBlockToBlockChain("Send 100RMB to lixiaohua")
	//blockchain.AddBlockToBlockChain("Send 100RMB to rose")
	//blockchain.PrintChains()

	cli := BLC.CLI{}
	cli.Run()
}
