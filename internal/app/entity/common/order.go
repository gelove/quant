package common

import "quant/internal/app/constant"

type OrderACT struct {
	Symbol             string          `json:"symbol" gorm:"column:symbol;not null;default:'';type:char(20)"`
	OrderID            int             `json:"orderId" gorm:"column:orderId;not null;default:0;type:int(10)"`
	OrderListID        int             `json:"orderListId" gorm:"column:orderListId;not null;default:0;type:int(10)"`
	ClientOrderID      string          `json:"clientOrderId" gorm:"column:clientOrderId;not null;default:'';type:char(50)"`
	Price              string          `json:"price" gorm:"column:price;not null;default:'';type:char(50)"`
	OrigQty            string          `json:"origQty" gorm:"column:origQty;not null;default:'';type:char(50)"`
	ExecutedQty        string          `json:"executedQty" gorm:"column:executedQty;not null;default:'';type:char(50)"`
	CumulativeQuoteQty string          `json:"cummulativeQuoteQty" gorm:"column:cumulativeQuoteQty;not null;default:'';type:char(50)"`
	Status             constant.STATUS `json:"status" gorm:"column:status;not null;default:'';comment:"状态:NEW,FILLED",type:char(50)"`
	TimeInForce        string          `json:"timeInForce" gorm:"column:timeInForce;not null;default:'';type:char(50)"`
	Type               string          `json:"type" gorm:"column:type;not null;default:'';comment:"交易类型:市价、限价、止盈止损";type:char(50)"`
	Side               constant.Side   `json:"side" gorm:"column:side;not null;default:'';type:char(50)"`
}

type Order struct {
	OrderACT
	OrderID           int    `json:"orderId" gorm:"column:orderId;not null;default:0;type:int(10)"`
	StopPrice         string `json:"stopPrice" gorm:"column:stopPrice;not null;default:'';type:char(50)"`
	IcebergQty        string `json:"icebergQty" gorm:"column:icebergQty;not null;default:'';type:char(50)"`
	Time              int64  `json:"time" gorm:"column:time;not null;default:0;type:int(10)"`
	UpdateTime        int64  `json:"updateTime" gorm:"column:updateTime;not null;default:0;type:int(10)"`
	IsWorking         int8   `json:"isWorking" gorm:"column:isWorking;not null;default:0;type:tinyint(3)"`
	OrigQuoteOrderQty string `json:"origQuoteOrderQty" gorm:"column:origQuoteOrderQty;not null;default:'';type:char(50)"`
}
