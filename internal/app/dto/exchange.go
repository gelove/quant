package dto

import (
	"quant/internal/app/constant"

	"github.com/shopspring/decimal"
)

type ExchangeData struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []*RateLimit  `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []*SymbolData `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type Filter struct {
	FilterType       constant.FilterType `json:"filterType"`
	MinPrice         string              `json:"minPrice"`
	MaxPrice         string              `json:"maxPrice"`
	TickSize         string              `json:"tickSize"`
	MultiplierUp     string              `json:"multiplierUp"`
	MultiplierDown   string              `json:"multiplierDown"`
	AvgPriceMins     int                 `json:"avgPriceMins"`
	MinQty           string              `json:"minQty"`
	MaxQty           string              `json:"maxQty"`
	StepSize         string              `json:"stepSize"`
	MinNotional      string              `json:"minNotional"`
	ApplyToMarket    bool                `json:"applyToMarket"`
	Limit            int                 `json:"limit"`
	MaxNumOrders     int                 `json:"maxNumOrders"`
	MaxNumAlgoOrders int                 `json:"maxNumAlgoOrders"`
}

type SymbolData struct {
	Symbol                     string    `json:"symbol"`
	Status                     string    `json:"status"`
	BaseAsset                  string    `json:"baseAsset"`
	BaseAssetPrecision         int       `json:"baseAssetPrecision"`
	QuoteAsset                 string    `json:"quoteAsset"`
	QuotePrecision             int       `json:"quotePrecision"`
	QuoteAssetPrecision        int       `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int       `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int       `json:"quoteCommissionPrecision"`
	OrderTypes                 []string  `json:"orderTypes"`
	IcebergAllowed             bool      `json:"icebergAllowed"`
	OcoAllowed                 bool      `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool      `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool      `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool      `json:"isMarginTradingAllowed"`
	Filters                    []*Filter `json:"filters"`
	Permissions                []string  `json:"permissions"`
}

type Steps struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
}
