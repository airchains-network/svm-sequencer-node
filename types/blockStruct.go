package types

type TransactionInfo struct {
	Meta        Meta        `json:"meta"`
	Transaction Transaction `json:"transaction"`
	Version     string      `json:"version"`
}

type Result struct {
	BlockHeight       int               `json:"blockHeight"`
	BlockTime         int               `json:"blockTime"`
	Blockhash         string            `json:"blockhash"`
	ParentSlot        int               `json:"parentSlot"`
	PreviousBlockhash string            `json:"previousBlockhash"`
	Transactions      []TransactionInfo `json:"transactions"`
}

type JsonRpcResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
	Id      int    `json:"id"`
}
