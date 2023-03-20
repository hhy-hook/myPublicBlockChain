package BLC

type TxInput struct {
	TxId      []byte
	Vout      int
	ScriptSiq string
}

func (txInput *TxInput) UnlockWithAddress(address string) bool {
	return txInput.ScriptSiq == address
}
