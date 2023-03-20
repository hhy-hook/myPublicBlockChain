package BLC

import (
	"time"
)

type Block struct {
	// 区块高度
	Height int64
	// 上一个区块hash
	PrevBlockHash []byte
	// 交易数据
	Data []byte
	// 时间戳
	TimeStamp int64
	// HASH
	Hash []byte
	// nonce
	Nonce int64
}

func NewBlock(data string, height int64, prevBlockHash []byte) *Block {

	block := &Block{height, prevBlockHash, []byte(data), time.Now().Unix(), nil, 0}

	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	return block
}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 0, make([]byte, 32, 32))
}
