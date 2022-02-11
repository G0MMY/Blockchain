package components

type MemPoolType struct {
	Transactions      []*TransactionType
	NodeAddress       string
	NumberTransaction int
}

type MemPool interface {
	addTransaction(*TransactionType)
	getTransactions() []*TransactionType
	deleteNFirstTransactions(int)
}

func (memPool *MemPoolType) addTransaction(transaction *TransactionType) {
	i := 0
	transactions := memPool.Transactions

	for i < memPool.NumberTransaction {
		if transactions[i].fee <= transaction.fee {
			memPool.Transactions = append(transactions[:i+1], transactions[i:]...)
			memPool.Transactions[i] = transaction
			memPool.NumberTransaction += 1
			return
		}
		i += 1
	}
	memPool.NumberTransaction += 1
	memPool.Transactions = append(memPool.Transactions, transaction)
}

func (memPool *MemPoolType) deleteNFirstTransactions(quantity int) {
	if quantity >= memPool.NumberTransaction {
		var transactions []*TransactionType
		memPool.Transactions = transactions
		memPool.NumberTransaction = 0
	} else {
		memPool.Transactions = memPool.Transactions[quantity:]
		memPool.NumberTransaction -= quantity
	}
}

func (memPool *MemPoolType) getTransactions() []*TransactionType {
	return memPool.Transactions
}
