package components

import "time"

type TransactionType struct {
	sender    string
	receiver  string
	amount    int
	fee       int
	timestamp int64
}

//has to check for sender balance
func CreateTransaction(sender string, receiver string, amount int, fee int) *TransactionType {
	return &TransactionType{sender, receiver, amount, fee, time.Now().Unix()}
}
