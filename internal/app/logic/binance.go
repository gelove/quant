package logic

import (
	"fmt"
	"log"
	"quant/internal/app/config"
	"quant/internal/app/constant"
	"quant/internal/app/dto"
	"quant/internal/app/entity"
	"quant/internal/app/model"
	"quant/pkg/socket"
	"quant/pkg/utils/date"
	"quant/pkg/utils/errs"
	"quant/pkg/utils/hash"
	"quant/pkg/utils/json"
	"quant/pkg/utils/logger"
	"quant/pkg/utils/request"
	"quant/pkg/utils/textual"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var balances []*dto.StreamBalance

var tokens = make([]TokenData, 0, 1<<9)

// ["BTCUSDT","BNBBTC"]
var symbolList = make([]string, 0, 1<<3)

var symbolSteps = make(map[string]*dto.Steps)

// var lastPrices = struct {
// 	Data map[string]decimal.Decimal
// 	sync.RWMutex
// }{
// 	Data: make(map[string]decimal.Decimal),
// }

var lastPrices = make(map[string]decimal.Decimal) // 最新成交价

var orderPrices = make(map[string]decimal.Decimal) // 我的最近成交价

var minAskPrices = make(map[string]decimal.Decimal) // 最少卖价

var streamChan = make(chan []byte, 10)

var listenKey string

const StreamSeparator = "/"

/**
K线数据
[
  [
    1499040000000,      // 开盘时间
    "0.01634790",       // 开盘价
    "0.80000000",       // 最高价
    "0.01575800",       // 最低价
    "0.01577100",       // 收盘价(当前K线未结束的即为最新价)
    "148976.11427815",  // 成交量
    1499644799999,      // 收盘时间
    "2434.19055334",    // 成交额
    308,                // 成交笔数
    "1756.87402397",    // 主动买入成交量
    "28.46694368",      // 主动买入成交额
    "17928899.62484339" // 请忽略该参数
  ]
]
*/
type Index = int

const (
	OpenAt Index = iota
	Open
	High
	Low
	Close
	CloseAt
)

type DepthIndex = int

const (
	DepthPrice Index = iota
	DepthQuantity
)

/*
当前最优挂单
{
  "symbol": "LTCBTC",
  "bidPrice": "4.00000000",
  "bidQty": "431.00000000",
  "askPrice": "4.00000200",
  "askQty": "9.00000000"
}
*/
type BestPrice struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}

type TokenData struct {
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"base_asset"`
	QuoteAsset string `json:"quote_asset"`
}

type BounceRate struct {
	Symbol  string          `json:"symbol"`
	Rate    decimal.Decimal `json:"rate"`
	Lowest  decimal.Decimal `json:"lowest"`
	Highest decimal.Decimal `json:"highest"`
}

type BounceRates []*BounceRate

func (s BounceRates) Len() int {
	return len(s)
}

func (s BounceRates) Less(i, j int) bool {
	return s[i].Rate.Cmp(s[j].Rate) == -1
}

func (s BounceRates) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func GetSymbols() []string {
	if len(symbolList) > 0 {
		return symbolList
	}
	assets := config.C.Binance.Assets
	symbols := make([]string, 0, len(assets))
	for _, v := range assets {
		symbols = append(symbols, v+config.C.Binance.QuoteAsset)
	}
	symbolList = symbols
	return symbolList
}

func Request(method request.Method, url string, res interface{}, header req.Header, vs ...interface{}) error {
	return request.New(method, config.C.Binance.BaseUrl+url, res, header, vs...)
}

func RequestWithSign(method request.Method, url string, res interface{}, params req.Param) error {
	apiKey := config.C.Binance.ApiKey
	secretKey := config.C.Binance.SecretKey
	if params != nil {
		// Encode会根据Ascii码排序
		body := request.ConvertUrlValue(params).Encode()
		// logger.Infof("RequestWithSign request body => %s", body)
		params["signature"] = hash.SHA256MAC([]byte(body), []byte(secretKey))
	}
	header := req.Header{
		"X-MBX-APIKEY": apiKey,
	}
	return Request(method, url, res, header, params)
}

// GetAllExchangeInfo
func GetAllExchangeInfo() error {
	// https://api.binance.com/api/v3/exchangeInfo
	res := new(dto.ExchangeData)
	err := Request(request.GET, constant.ExchangeInfoUrl, res, nil, nil)
	if err != nil {
		return err
	}
	list := make([]string, 0, 1<<9)
	for _, v := range res.Symbols {
		if !v.IsSpotTradingAllowed || v.Status != "TRADING" || textual.InArray(v.BaseAsset, list) {
			continue
		}
		if v.QuoteAsset == "USDT" {
			token := TokenData{v.Symbol, v.BaseAsset, v.QuoteAsset}
			tokens = append(tokens, token)
			list = append(list, v.BaseAsset)
			continue
		}
	}

	logger.Infof("ExchangeInfo tokens => %s", json.MustToString(tokens))
	return nil
}

// Bounce 计算大跌后的反弹比率
// token
// startAt 开始时间 13位毫秒时间戳
func bounce(token TokenData, startAt int64) (map[int]*BounceRate, error) {
	// https://api.binance.com/api/v3/klines?interval=1d&limit=10&symbol=BNBBTC
	hoursPerDay := 24
	limit := 7 * hoursPerDay
	params := req.Param{
		"symbol":   token.Symbol,
		"interval": "1h",
		"limit":    limit,
	}
	if startAt > 0 {
		params["startTime"] = startAt
	}
	// params["endTime"] = 1625043637000
	res := make([][]interface{}, 0, 1<<9)
	err := Request(request.GET, constant.KLineUrl, &res, nil, params)
	if err != nil {
		return nil, err
	}
	if len(res) < limit {
		return nil, errs.Get(errs.DataIsNotEnough)
	}
	// logger.Infof("GetKLine res => %s", json.MustToString(res))
	var lowest, highest decimal.Decimal
	var m = map[int]*BounceRate{
		3: {Symbol: token.Symbol},
		5: {Symbol: token.Symbol},
		7: {Symbol: token.Symbol},
	}
	for i, v := range res {
		low, err := decimal.NewFromString(v[Low].(string))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		high, err := decimal.NewFromString(v[High].(string))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if i == 0 {
			lowest = low
		}
		if low.Cmp(lowest) == -1 {
			lowest = low
		}
		if i != 0 && high.Cmp(highest) == 1 {
			highest = high
		}
		for k := range m {
			if i < k*hoursPerDay {
				m[k].Highest = highest
				m[k].Lowest = lowest
			}
		}
	}
	hundred := decimal.NewFromInt(100)
	one := decimal.NewFromInt(1)
	for _, v := range m {
		if v.Lowest.Equal(decimal.Zero) {
			return nil, errors.Errorf("bounce %s lowest为零", token.Symbol)
		}
		v.Rate = v.Highest.DivRound(v.Lowest, 4).Sub(one).Mul(hundred)
	}
	return m, nil
}

// GetTopBounce 获取反弹最多的币种
// number 反弹最多的前[number]个币种
// startTime btc到最低点的时刻 2006-01-02T15:04:05
func GetTopBounce(number int, startTime string) map[int]BounceRates {
	if len(tokens) == 0 {
		err := GetAllExchangeInfo()
		if err != nil {
			panic(err)
		}
	}

	startAt := date.GetMilliUnix(startTime, date.YMD_HIS)
	rates := map[int]BounceRates{
		3: make(BounceRates, 0, 1<<9),
		5: make(BounceRates, 0, 1<<9),
		7: make(BounceRates, 0, 1<<9),
	}

	var wg sync.WaitGroup
	for _, v := range tokens {
		wg.Add(1)
		time.Sleep(200 * time.Millisecond) // note: 接口有访问限制
		go func(token TokenData) {
			defer wg.Done()
			rate, err := bounce(token, startAt)
			if err != nil {
				logger.Errorf("GetTop error => %+v", err)
				return
			}
			for i, vv := range rate {
				rates[i] = append(rates[i], vv)
			}
		}(v)
	}
	wg.Wait()

	for i, v := range rates {
		sort.Sort(sort.Reverse(v))
		if len(v) > number {
			rates[i] = v[:number]
			continue
		}
		rates[i] = v
	}
	return rates
}

func amplitude(token TokenData) (*BounceRate, error) {
	limit := 26
	params := req.Param{
		"symbol":   token.Symbol,
		"interval": "1d",
		"limit":    limit,
	}
	res := make([][]interface{}, 0, 1<<9)
	err := Request(request.GET, constant.KLineUrl, &res, nil, params)
	if err != nil {
		return nil, err
	}
	if len(res) < limit {
		return nil, errs.Get(errs.DataIsNotEnough)
	}
	var lowest, highest decimal.Decimal
	rates := make([]decimal.Decimal, 0, len(res))
	hundred := decimal.NewFromInt(100)
	one := decimal.NewFromInt(1)
	for i, v := range res {
		low, err := decimal.NewFromString(v[Low].(string))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		high, err := decimal.NewFromString(v[High].(string))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if i == 0 {
			lowest = low
		}
		if low.Cmp(lowest) == -1 {
			lowest = low
		}
		if i != 0 && high.Cmp(highest) == 1 {
			highest = high
		}
		if low.Equal(decimal.Zero) {
			return nil, errors.Errorf("amplitude %s low为零", token.Symbol)
		}
		rate := high.DivRound(low, 4).Sub(one).Mul(hundred)
		rates = append(rates, rate)
	}
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Cmp(rates[j]) == -1
	})
	length := len(rates)
	list := rates[3:(length - 3)]
	var amount decimal.Decimal
	for _, v := range list {
		amount = amount.Add(v)
	}
	data := &BounceRate{Symbol: token.Symbol, Lowest: lowest, Highest: highest}
	data.Rate = amount.DivRound(decimal.NewFromInt(int64(len(list))), 2)
	return data, nil
}

// GetTopAmplitude 获取振幅最优的标的
// number 获取振幅最优的[number]标的
func GetTopAmplitude(number int) []*BounceRate {
	if len(tokens) == 0 {
		err := GetAllExchangeInfo()
		if err != nil {
			panic(err)
		}
	}

	rates := make([]*BounceRate, 0, 1<<9)

	var wg sync.WaitGroup
	for _, v := range tokens {
		wg.Add(1)
		time.Sleep(200 * time.Millisecond) // note: 接口有访问限制
		go func(token TokenData) {
			defer wg.Done()
			rate, err := amplitude(token)
			if err != nil {
				logger.Errorf("GetTop error => %+v", err)
				return
			}
			rates = append(rates, rate)
		}(v)
	}
	wg.Wait()

	sort.Sort(sort.Reverse(BounceRates(rates)))
	if len(rates) > number {
		return rates[:number]
	}
	return rates
}

func ExchangeInfo() error {
	symbols := GetSymbols()
	log.Printf("symbols => %+v", symbols)
	list, err := GetExchangeInfo(symbols)
	if err != nil {
		return err
	}
	setExchangeInfo(list)
	return nil
}

// GetExchangeInfo 交易规范信息
func GetExchangeInfo(symbols []string) ([]*dto.SymbolData, error) {
	res := new(dto.ExchangeData)
	params := req.Param{
		"symbols": fmt.Sprintf(`["%s"]`, strings.Join(symbols, `","`)),
	}
	log.Printf("params => %+v", params)
	err := Request(request.GET, constant.ExchangeInfoUrl, res, nil, params)
	if err != nil {
		return nil, err
	}
	return res.Symbols, nil
}

func setExchangeInfo(list []*dto.SymbolData) {
	for _, item := range list {
		step := &dto.Steps{}
		for _, v := range item.Filters {
			if v.FilterType == constant.PRICE_FILTER {
				step.Price = decimal.RequireFromString(v.TickSize)
				continue
			}
			if v.FilterType == constant.LOT_SIZE {
				step.Quantity = decimal.RequireFromString(v.StepSize)
				continue
			}
		}
		symbolSteps[item.Symbol] = step
	}
}

// BookTicker 获取当前最优挂单
func BookTicker(symbol string) (res *BestPrice, err error) {
	if symbol == "" {
		err = errors.New("symbol 不得为空")
		return
	}
	params := req.Param{
		"symbol": symbol,
	}
	res = &BestPrice{}
	err = Request(request.GET, constant.BookTicker, res, nil, params)
	if err != nil {
		return
	}
	return
}

// GetDepth 获取深度信息
func GetDepth(symbol string, limit int) (res *dto.Depth, err error) {
	if symbol == "" {
		err = errors.New("symbol 不得为空")
		return
	}
	params := req.Param{
		"symbol": symbol,
		"limit":  limit,
	}
	res = &dto.Depth{}
	err = Request(request.GET, constant.Depth, res, nil, params)
	if err != nil {
		return
	}
	return
}

// GetAccount 获取账户信息
func GetAccount() (res *dto.Account, err error) {
	params := req.Param{
		"recvWindow": 10000,
		"timestamp":  date.UnixMilli(nil),
	}
	res = &dto.Account{}
	err = RequestWithSign(request.GET, constant.Account, res, params)
	if err != nil {
		return
	}
	return
}

// GetALLOrders
func GetALLOrders(symbol string) (res []*dto.Order, err error) {
	if symbol == "" {
		err = errors.New("symbol 不得为空")
		return
	}
	params := req.Param{
		"symbol":     symbol,
		"recvWindow": 10000,
		"timestamp":  date.UnixMilli(nil),
	}
	res = make([]*dto.Order, 0, 1<<3)
	err = RequestWithSign(request.GET, constant.AllOrders, &res, params)
	if err != nil {
		return
	}
	return
}

// GetOpenOrders 获取交易对的所有当前挂单
func GetOpenOrders(symbol string) (res []*dto.Order, err error) {
	if symbol == "" {
		err = errors.New("symbol 不得为空")
		return
	}
	params := req.Param{
		"symbol":     symbol,
		"recvWindow": 10000,
		"timestamp":  date.UnixMilli(nil),
	}
	res = make([]*dto.Order, 0, 1<<3)
	err = RequestWithSign(request.GET, constant.OpenOrders, &res, params)
	if err != nil {
		return
	}
	return
}

/*
   $ echo -n "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=10000&timestamp=1499827319559" | openssl dgst -sha256 -hmac "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"
   (stdin)= c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71

   (HMAC SHA256)
    $ curl -H "X-MBX-APIKEY: vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A" -X POST 'https://api.binance.com/api/v3/order' -d 'symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=10000&timestamp=1499827319559&signature=c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71'
*/
// PlaceOrder 下单
func PlaceOrder(symbol string, side constant.Side, quantity, price decimal.Decimal) (err error) {
	if symbol == "" {
		err = errors.New("symbol 不得为空")
		return
	}
	params := req.Param{
		"symbol":      symbol,
		"side":        side,
		"type":        constant.LIMIT,
		"timeInForce": "GTC",
		"quantity":    quantity,
		"price":       price,
		"recvWindow":  10000,
		"timestamp":   date.UnixMilli(nil),
	}
	res := &dto.OrderACT{}
	err = RequestWithSign(request.POST, constant.PlaceOrder, &res, params)
	if err != nil {
		return
	}
	// 订单存储到本地数据库
	orderModel := &model.Order{}
	err = orderModel.Create(res.ToEntity())
	return
}

// GetListenKey 获取 Listen Key
func GetListenKey() (key string, err error) {
	if listenKey != "" {
		key = listenKey
		return
	}
	res := make(map[string]string)
	err = RequestWithSign(request.POST, constant.UserDataStream, &res, nil)
	if err != nil {
		return
	}
	key = res["listenKey"]
	if key == "" {
		err = errors.New("listenKey 不得为空")
		return
	}
	listenKey = key
	return
}

// GetLocalOrders 获取本地订单数据
func GetLocalOrders(symbol string) (orders []*entity.Order, err error) {
	order := &model.Order{}
	return order.Get(symbol)
}

func Run() error {
	// 首次启动，判断接口返回的订单数量和symbol当前的余额
	// symbol余额大于0则先提醒手动挂卖单
	// symbol余额为空，并没有买单则先挂买单
	symbols := GetSymbols()
	err := ExchangeInfo()
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		orders, err := GetLocalOrders(symbol)
		if err != nil {
			return err
		}
		buy := make([]*entity.Order, 0, 1<<2)
		sell := make([]*entity.Order, 0, 1<<2)
		for _, v := range orders {
			if v.Side == constant.BUY {
				buy = append(buy, v)
			}
			if v.Side == constant.SELL {
				sell = append(sell, v)
			}
		}
		account, err := GetAccount()
		if err != nil {
			return err
		}
		for _, v := range account.Balances {
			if v.Asset == symbol {
				free := decimal.RequireFromString(v.Free)
				// 当前有网格买到的token却没卖, 先卖出
				if free.Cmp(decimal.Zero) > 1 && len(buy) > len(sell) {
					var lastBuy = buy[0]
					for _, v := range buy {
						if v.OrderID > lastBuy.OrderID {
							lastBuy = v
						}
					}
					lastPrice := decimal.RequireFromString(lastBuy.Price)
					price := formatPrice(lastPrice, symbol, constant.SELL)
					err = PlaceOrder(symbol, constant.SELL, free, price)
					if err != nil {
						logger.Errorf("PlaceOrder err : %+v", err)
					}
				}
			}
		}
	}

	return nil
}

func formatPrice(price decimal.Decimal, symbol string, side constant.Side) decimal.Decimal {
	step := symbolSteps[symbol]
	if step == nil || step.Price == decimal.Zero {
		return decimal.Zero
	}
	point := decimal.NewFromInt(config.C.Binance.Point)
	var percentage decimal.Decimal
	hundred := decimal.NewFromInt(100)
	if side == constant.BUY {
		percentage = hundred.Sub(point)
	}
	if side == constant.SELL {
		percentage = hundred.Add(point)
	}
	return price.Mul(percentage).Div(hundred).Div(step.Price).Floor().Mul(step.Price)
}

func formatQuantity(quantity decimal.Decimal, symbol string) decimal.Decimal {
	step := symbolSteps[symbol]
	if step == nil || step.Quantity == decimal.Zero {
		return decimal.Zero
	}
	return quantity.Div(step.Quantity).Floor().Mul(step.Quantity)
}

func SuccessHandler(socket *socket.Socket, message []byte) {
	logger.Infof("SuccessHandler => %s", message)
	streamChan <- message
}

func ErrorHandler(err error) {
	logger.Error(err)
}

func processUserData(data *dto.StreamData) {
	switch data.Kind {
	case constant.ExecutionReport:
		AutoPlaceOrder(data)
	case constant.OutboundAccountPosition:
		balances = data.Balances
		logger.Infof("balances => %#v", balances)
	default:
	}
}

func ProcessStream() {
	for v := range streamChan {
		// logger.Infof("ProcessStream => %s", v)
		res := &dto.Frame{}
		json.MustDecode(v, res)
		switch true {
		case res.Stream == listenKey:
			data := &dto.StreamData{}
			json.MustTransform(res.Data, data)
			processUserData(data)
		case strings.HasSuffix(res.Stream, constant.TradeSuffix):
			trade := &dto.Trade{}
			json.MustTransform(res.Data, trade)
			if trade.Price.Cmp(decimal.Zero) == 1 {
				lastPrices[trade.Symbol] = trade.Price
			}
		case strings.HasSuffix(res.Stream, constant.DepthSuffix):
			depth := &dto.Depth{}
			json.MustTransform(res.Data, depth)
			symbol := strings.TrimSuffix(res.Stream, constant.DepthSuffix)
			if len(depth.Asks) > 0 {
				min := depth.Asks[0]
				minAskPrices[symbol] = min[DepthPrice]
			}
		}
	}
}

func AutoPlaceOrder(data *dto.StreamData) {
	logger.Infof("AutoPlaceOrder data => %+v", *data)
	symbols := GetSymbols()
	if !textual.InArray(data.Symbol, symbols) {
		return
	}
	if data.Status == constant.FILLED {
		// 更新本地数据状态
		err := Update(data)
		if err != nil {
			logger.Errorf("AutoPlaceOrder Update (%d) err => %+v", data.OrderID, err)
			return
		}
		if data.Side == constant.BUY {
			// 买单已经全部成交, 则挂卖单
			price := formatPrice(data.Price, data.Symbol, constant.SELL)
			err := PlaceOrder(data.Symbol, constant.SELL, data.Quantity, price)
			logger.Errorf("AutoPlaceOrder SELL err => %+v", err)
			if err != nil {
				logger.Errorf("AutoPlaceOrder SELL (%d) err => %+v", data.OrderID, err)
				return
			}
			return
		}
		if data.Side == constant.SELL {
			// 卖单已经全部成交, 则挂买单
			money := config.C.Binance.Amount / config.C.Binance.MaxOpenOrders
			price := formatPrice(data.Price, data.Symbol, constant.BUY)
			quantity := formatQuantity(decimal.NewFromInt(money).Div(price), data.Symbol)
			err := PlaceOrder(data.Symbol, constant.BUY, quantity, price)
			if err != nil {
				logger.Errorf("AutoPlaceOrder BUY (%d) err => %+v", data.OrderID, err)
			}
			return
		}
	}
}

func Update(data *dto.StreamData) error {
	orderModel := &model.Order{}
	order, err := orderModel.Find(data.Symbol, data.OrderID)
	if err != nil {
		return err
	}
	order.Status = data.Status
	err = orderModel.Save(order)
	if err != nil {
		return err
	}
	return nil
}
