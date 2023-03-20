package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"time"
)

type Transaction struct {
	TxId  []byte
	Vins  []*TxInput
	Vouts []*TxOutput
}

func NewCoinBaseTransaction(value int64, address string) *Transaction {
	txInput := &TxInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TxOutput{value, address}
	txCoinbase := &Transaction{[]byte{}, []*TxInput{txInput}, []*TxOutput{txOutput}}
	txCoinbase.SetTxId()
	return txCoinbase
}

func (tx *Transaction) SetTxId() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	buffBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), buff.Bytes()}, []byte{})

	hash := sha256.Sum256(buffBytes)
	tx.TxId = hash[:]
}

func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxId) == 0 && tx.Vins[0].Vout == -1
}

func NewSimpleTransaction(from, to string, amount int64, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	balance, spendableUTXO := bc.FindSpendableUTXOs(from, amount, txs)

	for txID, indexArray := range spendableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TxInput{txIDBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutput1 := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput1)

	txOutput2 := &TxOutput{balance - amount, from}
	txOutputs = append(txOutputs, txOutput2)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.SetTxId()
	return tx
}
