package dto

type TxResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []*Transaction `json:"result"`
}

type Transaction struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

type Token struct {
	Data TokenData `json:"data"`
}

type TokenData struct {
	AsAddress []*AsAddress `json:"asAddress"`
}

type AsAddress struct {
	Typename       string `json:"__typename"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	TotalLiquidity string `json:"totalLiquidity"`
}

type LiquidityEvent struct {
	TokenIn        string `json:"tokenIn"`
	TokenOut       string `json:"tokenOut"`
	AmountIn       string `json:"amountIn"`
	AmountOut      string `json:"amountOut"`
	UsdAmount      string `json:"usdAmount"`
	ContractIn     string `json:"contractIn"`
	TotalLiquidity string `json:"totalLiquidity"`
	Hash           string `json:"hash"`
	Price          string `json:"price"`
	Time           string `json:"time"`
}

type RequestPairs struct {
	Data struct {
		Pair struct {
			ID          string `json:"id"`
			ReserveUSD  string `json:"reserveUSD"`
			TotalSupply string `json:"totalSupply"`
		} `json:"pair"`
	} `json:"data"`
}
