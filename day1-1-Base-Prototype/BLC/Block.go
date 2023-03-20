package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

func (b *Block) SetHash() {

	// step.1 属性转[]bytes
	h := IntToHex(b.Height)
	//fmt.Println(h)
	//t := IntToHex(b.TimeStamp)
	//fmt.Println("timestamp:", t)
	t2 := strconv.FormatInt(b.TimeStamp, 2)
	t2bytes := []byte(t2)
	//fmt.Println("timestamp2:", t2bytes)

	// step.2 二维数组拼接
	blockBytes := bytes.Join([][]byte{
		h, b.PrevBlockHash, b.Data, t2bytes}, []byte{})

	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	block := &Block{height, prevBlockHash, []byte(data), time.Now().Unix(), nil}
	block.SetHash()
	return block
}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 0, make([]byte, 32, 32))
}
