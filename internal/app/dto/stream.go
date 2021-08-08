package dto

import (
	"quant/internal/app/constant"

	"github.com/shopspring/decimal"
)

type Frame struct {
	Stream string      `json:"stream"`
	Data   interface{} `json:"data"`
}

type StreamData struct {
	Kind               constant.EventKind `json:"e"` // 事件类型
	E                  int64              `json:"E"` // 事件时间
	UpdatedAt          int64              `json:"u"` // 账户末次更新时间戳
	BuyOrderID         int64              `json:"b"` // 买方的订单ID
	SellOrderID        int64              `json:"a"` // 卖方的订单ID
	Balances           []*StreamBalance   `json:"B"` // 账户余额
	Symbol             string             `json:"s"` // 交易对
	ClientOrderID      string             `json:"c"` // clientOrderId
	Side               constant.Side      `json:"S"` // 订单方向
	Type               string             `json:"o"` // 订单类型
	TimeInForce        string             `json:"f"` // 有效方式, 订单多久能够失效M, GTC 成交为止
	Quantity           decimal.Decimal    `json:"q"` // 订单原始数量
	Price              decimal.Decimal    `json:"p"` // 订单原始价格
	EmitPrice          string             `json:"P"` // 止盈止损单触发价格
	IcebergQty         string             `json:"F"` // 冰山订单数量
	OrderListID        int64              `json:"g"` // OCO订单 OrderListId
	CanceledOrderID    string             `json:"C"` // 原始订单自定义ID(原始订单，指撤单操作的对象。撤单本身被视为另一个订单)
	X                  string             `json:"x"` // 本次事件的具体执行类型
	Status             constant.STATUS    `json:"X"` // 订单的当前状态
	RejectReason       string             `json:"r"` // 订单被拒绝的原因
	OrderID            int64              `json:"i"` // orderId
	LastQuantity       string             `json:"l"` // 订单末次成交量
	AccumulateQuantity string             `json:"z"` // 订单累计已成交量
	LastPrice          string             `json:"L"` // 订单末次成交价格
	ChargeQuantity     string             `json:"n"` // 手续费数量
	ChargeAssert       string             `json:"N"` // 手续费资产类别
	DealAt             int64              `json:"T"` // 成交时间
	DealID             int64              `json:"t"` // 成交ID
	I                  int64              `json:"I"` // 请忽略
	W                  bool               `json:"w"` // 订单是否在订单簿上？
	IsEntry            bool               `json:"m"` // 该成交是作为挂单成交吗？
	M                  bool               `json:"M"` // 请忽略
	OrderCreatedAt     int64              `json:"O"` // 订单创建时间
	AccumulateAmount   string             `json:"Z"` // 订单累计已成交金额
	LastAmount         string             `json:"Y"` // 订单末次成交金额
	Qty                string             `json:"Q"` // Quote Order Qty
}

type StreamBalance struct {
	Asset  string `json:"a"` // 资产名称
	Free   string `json:"f"` // 可用余额
	Locked string `json:"l"` // 冻结余额
}

type Depth struct {
	LastUpdateID int                 `json:"lastUpdateId"`
	Bids         [][]decimal.Decimal `json:"bids"`
	Asks         [][]decimal.Decimal `json:"asks"`
}

/*
   "e": "trade",
   "E": 1627825196488,
   "s": "BTCUSDT",
   "t": 79952,
   "p": "41453.97000000",
   "q": "0.00241100",
   "b": 673503, // 买方的订单ID
   "a": 673367, // 卖方的订单ID
   "T": 1627825196487,
   "m": false,
   "M": true
*/
type Trade struct {
	Kind        constant.EventKind `json:"e"` // 事件类型
	E           int64              `json:"E"` // 事件时间
	Symbol      string             `json:"s"` // 交易对
	DealID      int64              `json:"t"` // 成交ID
	Quantity    decimal.Decimal    `json:"q"` // 订单原始数量
	Price       decimal.Decimal    `json:"p"` // 订单原始价格
	BuyOrderID  int64              `json:"b"` // 买方的订单ID
	SellOrderID int64              `json:"a"` // 卖方的订单ID
	DealAt      int64              `json:"T"` // 成交时间
	IsEntry     bool               `json:"m"` // 该成交是作为挂单成交吗？
	M           bool               `json:"M"` // 请忽略
}
