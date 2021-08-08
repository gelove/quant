package dto

import (
	"quant/internal/app/entity"
	"quant/internal/app/entity/common"
	"time"
)

type Order struct {
	common.Order
	IsWorking bool `json:"isWorking" gorm:"column:isWorking;not null;default:0;type:tinyint(3)"`
}

type OrderACT struct {
	common.OrderACT
	TransactTime int64 `json:"transactTime"`
}

func (o *OrderACT) ToEntity() *entity.Order {
	res := &entity.Order{}
	res.Symbol = o.Symbol
	res.OrderID = o.OrderID
	res.OrderListID = o.OrderListID
	res.ClientOrderID = o.ClientOrderID
	res.Price = o.Price
	res.OrigQty = o.OrigQty
	res.ExecutedQty = o.ExecutedQty
	res.CumulativeQuoteQty = o.CumulativeQuoteQty
	res.Status = o.Status
	res.TimeInForce = o.TimeInForce
	res.Type = o.Type
	res.Side = o.Side
	res.CreatedAt = time.Unix(o.TransactTime/1000, o.TransactTime%1000*int64(time.Millisecond))
	return res
}
