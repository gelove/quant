package dto

type Account struct {
	MakerCommission  int        `json:"makerCommission"`
	TakerCommission  int        `json:"takerCommission"`
	BuyerCommission  int        `json:"buyerCommission"`
	SellerCommission int        `json:"sellerCommission"`
	CanTrade         bool       `json:"canTrade"`
	CanWithdraw      bool       `json:"canWithdraw"`
	CanDeposit       bool       `json:"canDeposit"`
	UpdateTime       int        `json:"updateTime"`
	AccountType      string     `json:"accountType"`
	Balances         []*Balance `json:"balances"`
	Permissions      []string   `json:"permissions"`
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type Go struct {
	Stream string `json:"stream"`
	Data   GoData `json:"data"`
}

type GoData struct {
	Type       string `json:"e"` // 事件类型
	Time       int64  `json:"E"` // 事件时间
	Symbol     string `json:"s"` // 交易对
	SetID      int    `json:"a"` // 归集交易ID
	Price      string `json:"p"` // 成交价格
	Quantity   string `json:"q"` // 成交笔数
	FirstID    int    `json:"f"` // 被归集的首个交易ID
	LastID     int    `json:"l"` // 被归集的末次交易ID
	DeliveryAt int64  `json:"T"` // 成交时间
	IsMarket   bool   `json:"m"` // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
	M          bool   `json:"M"` // 请忽略该字段
}
