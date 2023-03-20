package BLC

type BlockChain struct {
	Blocks []*Block
}

func CreateBlockChainWithGenesisBlock(data string) *BlockChain {
	//创建创世区块
	genesisBlock := CreateGenesisBlock(data)
	//返回区块链对象
	return &BlockChain{[]*Block{genesisBlock}}
}

func (bc *BlockChain) AddBlockToBlockChain(data string, height int64, prevHash []byte) {
	//fmt.Println(prevHash)
	newBlock := NewBlock(data, height, prevHash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
