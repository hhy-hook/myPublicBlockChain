package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

func (cli *CLI) Run() {
	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "default", "交易数据")
	flagCreateBlockChainData := createBlockChainCmd.String("data", "default", "创世区块交易数据")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockData)
	}
	if printChainCmd.Parsed() {
		cli.printChains()
	}
	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockChainData)
	}

}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateBlockChain -data DATA -- 创建创世区块")
	fmt.Println("\taddBlock -data Data -- 交易数据")
	fmt.Println("\tprintChain -- 输出信息")
}

func (cli *CLI) addBlock(data string) {
	bc := GetBlockChainObj()
	if bc == nil {
		fmt.Println("没有创世区块，无法添加。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.AddBlockToBlockChain(data)
}

func (cli *CLI) printChains() {
	bc := GetBlockChainObj()
	if bc == nil {
		fmt.Println("没有创世区块，无法添加。。")
		os.Exit(1)
	}
	bc.PrintChains()
}

func (cli *CLI) createGenesisBlockchain(data string) {
	CreateBlockChainWithGenesisBlock(data)
}
