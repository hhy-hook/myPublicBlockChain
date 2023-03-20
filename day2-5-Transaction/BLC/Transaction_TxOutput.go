package BLC

type TxOutput struct {
	Value        int64
	ScriptPubKey string
}

func (txOutput *TxOutput) UnLockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
