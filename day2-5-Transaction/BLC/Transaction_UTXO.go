package BLC

type UTXO struct {
	TxID   []byte //当前Transaction的交易ID
	Index  int    //下标索引
	Output *TxOutput
}
