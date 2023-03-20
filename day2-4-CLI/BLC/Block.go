package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
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

// Serialize 序列化器
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	var reader = bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
