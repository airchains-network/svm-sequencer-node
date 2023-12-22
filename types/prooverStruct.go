package types

type GetTransactionStruct struct {
	To              string
	From            string
	Amount          float64
	FromBalances    float64
	ToBalances      float64
	TransactionHash string
}
