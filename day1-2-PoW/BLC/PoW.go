package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const TargetBit = 16

type PoW struct {
	Block *Block
	// 目标hash
	Target *big.Int
}

func NewProofOfWork(b *Block) *PoW {
	target := big.NewInt(1)

	//2.左移256-bits位
	target = target.Lsh(target, 256-TargetBit)

	return &PoW{b, target}
}

func (pow *PoW) Run() ([]byte, int64) {
	nonce := pow.Block.Nonce
	hashInt := new(big.Int)
	hash := [32]byte{}
	for {
		dataBytes := pow.paraseBlockData(nonce)
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("%d : %x\n", nonce, hash)
		hashInt.SetBytes(hash[:])
		if pow.Target.Cmp(hashInt) == 1 {
			break
		}
		nonce++
	}
	return hash[:], nonce
}

func (pow *PoW) paraseBlockData(nonce int64) []byte {
	dataBytes := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		IntToHex(pow.Block.TimeStamp),
		IntToHex(int64(TargetBit)),
		IntToHex(nonce),
	}, []byte{})
	return dataBytes
}

func (pow *PoW) Valid() bool {
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) == 1
}
