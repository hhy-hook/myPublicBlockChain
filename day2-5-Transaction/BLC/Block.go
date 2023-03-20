package BLC

import (
	"bytes"
	"crypto/sha256"
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
	Txs []*Transaction
	// 时间戳
	TimeStamp int64
	// HASH
	Hash []byte
	// nonce
	Nonce int64
}

func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {

	block := &Block{height, prevBlockHash, txs, time.Now().Unix(), nil, 0}

	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	return block
}

func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 0, make([]byte, 32, 32))
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

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxId)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
