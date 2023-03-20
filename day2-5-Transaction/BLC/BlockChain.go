package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

type BlockChain struct {
	//Blocks []*Block
	Tip []byte
	DB  *bolt.DB
}

func CreateBlockChainWithGenesisBlock(value int64, address string) {
	//创建创世区块
	if DBExists() {
		fmt.Println("数据库已经存在。。。")
		return
	}
	txCoinBase := NewCoinBaseTransaction(value, address)
	genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BLOCKTABLENAME))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic("创世区块存储有误。。。")
			}
			//存储最新区块的hash
			b.Put([]byte("l"), genesisBlock.Hash)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	//fmt.Println(prevHash)
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//2.打开表
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			//2.根据最新块的hash读取数据，并反序列化最后一个区块
			blockBytes := b.Get(bc.Tip)
			lastBlock := Deserialize(blockBytes)
			//3.创建新的区块
			newBlock := NewBlock(txs, lastBlock.Height+1, lastBlock.Hash)
			//4.将新的区块序列化并存储
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//5.更新最后一个哈希值，以及blockchain的tip
			b.Put([]byte("l"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func DBExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

func (bc *BlockChain) PrintChains() {
	//1.获取迭代器对象
	bcIterator := bc.Iterator()

	var count = 0
	//2.循环迭代
	for {
		block := bcIterator.Next()
		count++
		fmt.Printf("第%d个区块的信息：\n", count)
		//获取当前hash对应的数据，并进行反序列化
		fmt.Printf("\t高度：%d\n", block.Height)
		fmt.Printf("\t上一个区块的hash：%x\n", block.PrevBlockHash)
		fmt.Printf("\t当前的hash：%x\n", block.Hash)
		fmt.Println("\t交易：")
		for _, tx := range block.Txs {
			fmt.Printf("\t\t交易ID：%x\n", tx.TxId)
			fmt.Println("\t\tVins：")
			for _, in := range tx.Vins {
				fmt.Printf("\t\t\tTxID:%x\n", in.TxId)
				fmt.Printf("\t\t\tVout:%d\n", in.Vout)
				fmt.Printf("\t\t\tScriptSiq:%s\n", in.ScriptSiq)
			}
			fmt.Println("\t\tVouts：")
			for _, out := range tx.Vouts {
				fmt.Printf("\t\t\tvalue:%d\n", out.Value)
				fmt.Printf("\t\t\tScriptPubKey:%s\n", out.ScriptPubKey)
			}
		}
		//fmt.Printf("\t时间：%v\n", block.TimeStamp)
		fmt.Printf("\t时间：%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("\t次数：%d\n", block.Nonce)

		//3.直到父hash值为0
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
}

func GetBlockChainObj() *BlockChain {
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var blockchain *BlockChain
	//B：读取数据库
	err = db.View(func(tx *bolt.Tx) error {
		//C：打开表
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			//D：读取最后一个hash
			hash := b.Get([]byte("l"))
			//E：创建blockchain
			blockchain = &BlockChain{hash, db}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return blockchain
}

func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	var txs []*Transaction
	for i := 0; i < len(from); i++ {
		amountInt64, _ := strconv.ParseInt(amount[i], 10, 64)
		tx := NewSimpleTransaction(from[i], to[i], amountInt64, bc, txs)
		txs = append(txs, tx)
	}
	bc.AddBlockToBlockChain(txs)
}

func (bc *BlockChain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	/*
		1.获取交易中的所有可以用的
		2.获取数据库中可以使用的
	*/
	var balance int64
	utxos := bc.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxID)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if balance >= amount {
			break
		}
	}
	if balance < amount {
		fmt.Printf("%s 余额不足。。总额：%d，需要：%d\n", from, balance, amount)
		os.Exit(1)
	}
	return balance, spendableUTXO
}

// UnUTXOs 查询未使用的UTXO
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO                      //未花费
	spentTxOutputs := make(map[string][]int) //存储已经花费
	for i := len(txs) - 1; i >= 0; i-- {
		unUTXOs = caculate(txs[i], address, spentTxOutputs, unUTXOs)
	}

	bcIterator := bc.Iterator()
	for {
		block := bcIterator.Next()
		//统计未花费
		for i := len(block.Txs) - 1; i >= 0; i-- {
			unUTXOs = caculate(block.Txs[i], address, spentTxOutputs, unUTXOs)
		}

		//for _, utxo := range unUTXOs {
		//	fmt.Printf("ID:%x", utxo.TxID)
		//	fmt.Println(utxo.Index)
		//	fmt.Println(utxo.Output.Value)
		//	fmt.Println(utxo.Output.ScriptPubKey)
		//}

		//结束迭代
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}

func caculate(tx *Transaction, address string, spentTxOutputs map[string][]int, unUTXOs []*UTXO) []*UTXO {
	//fmt.Println(spentTxOutputs)
	if !tx.IsCoinbaseTransaction() {
		for _, in := range tx.Vins {
			if in.UnlockWithAddress(address) {
				key := hex.EncodeToString(in.TxId)
				//fmt.Println("key: ", key)
				spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
			}
		}
	}
outputs:
	for index, out := range tx.Vouts {
		if out.UnLockWithAddress(address) {
			if len(spentTxOutputs) != 0 {
				var isSpentUTXO bool
				for txID, indexArray := range spentTxOutputs {
					for _, i := range indexArray {
						if i == index && txID == hex.EncodeToString(tx.TxId) {
							isSpentUTXO = true
							continue outputs
						}
					}

				}
				if !isSpentUTXO {
					utxo := &UTXO{tx.TxId, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}
			} else {
				utxo := &UTXO{tx.TxId, index, out}
				unUTXOs = append(unUTXOs, utxo)
			}
		}
	}

	return unUTXOs
}

func (bc *BlockChain) GetBalance(address string, txs []*Transaction) int64 {
	//txOutputs:=bc.UnUTXOs(address)
	unUTXOs := bc.UnUTXOs(address, txs)
	//fmt.Println(address, unUTXOs)
	var amount int64
	for _, utxo := range unUTXOs {
		amount = amount + utxo.Output.Value
	}
	return amount

}
