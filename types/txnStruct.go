package types

type Status struct {
	Ok interface{} `json:"Ok"`
}

type LoadedAddresses struct {
	Readonly []interface{} `json:"readonly"`
	Writable []interface{} `json:"writable"`
}

type Meta struct {
	ComputeUnitsConsumed int             `json:"computeUnitsConsumed"`
	Err                  interface{}     `json:"err"`
	Fee                  int             `json:"fee"`
	InnerInstructions    []interface{}   `json:"innerInstructions"`
	LoadedAddresses      LoadedAddresses `json:"loadedAddresses"`
	LogMessages          []string        `json:"logMessages"`
	PostBalances         []int           `json:"postBalances"`
	PostTokenBalances    []interface{}   `json:"postTokenBalances"`
	PreBalances          []int           `json:"preBalances"`
	PreTokenBalances     []interface{}   `json:"preTokenBalances"`
	Rewards              []interface{}   `json:"rewards"`
	Status               Status          `json:"status"`
}

type Header struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}

type Instruction struct {
	Accounts       []int       `json:"accounts"`
	Data           string      `json:"data"`
	ProgramIdIndex int         `json:"programIdIndex"`
	StackHeight    interface{} `json:"stackHeight"`
}

type Message struct {
	AccountKeys     []string      `json:"accountKeys"`
	Header          Header        `json:"header"`
	Instructions    []Instruction `json:"instructions"`
	RecentBlockhash string        `json:"recentBlockhash"`
}

type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

type TxnResult struct {
	BlockTime   int         `json:"blockTime"`
	Meta        Meta        `json:"meta"`
	Slot        int         `json:"slot"`
	Transaction Transaction `json:"transaction"`
}

type TxnResponce struct {
	Jsonrpc string    `json:"jsonrpc"`
	Result  TxnResult `json:"result"`
	Id      int       `json:"id"`
}
