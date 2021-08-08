package constant

// note: 国内已被墙 走vpn
const ExchangeInfoUrl = "exchangeInfo"  // token信息
const KLineUrl = "klines"               // K线
const Account = "account"               // 账户信息
const UserDataStream = "userDataStream" // 生成 Listen Key
const PlaceOrder = "order"              // 下单
const AllOrders = "allOrders"           // 所有订单
const OpenOrders = "openOrders"         // 当前挂单
const Depth = "depth"                   // 深度信息
const BookTicker = "ticker/bookTicker"  // 返回当前最优的挂单(最高买单，最低卖单)

// 交易方向
type Side string

const (
	SELL Side = "SELL"
	BUY  Side = "BUY"
)

// 交易类型
type Kind string

const (
	LIMIT  Kind = "LIMIT"  // 限价
	MARKET Kind = "MARKET" // 市价
	STOP   Kind = "STOP"   // 止盈止损
)

type STATUS string

const (
	NEW              STATUS = "NEW"              // 订单被交易引擎接受
	FILLED           STATUS = "FILLED"           // 订单完全成交
	PARTIALLY_FILLED STATUS = "PARTIALLY_FILLED" // 部分订单被成交
	CANCELED         STATUS = "CANCELED"         // 用户撤销了订单
	PENDING_CANCEL   STATUS = "PENDING_CANCEL"   // 撤销中(目前并未使用)
	REJECTED         STATUS = "REJECTED"         // 订单没有被交易引擎接受，也没被处理
	EXPIRED          STATUS = "EXPIRED"          // 订单被交易引擎取消, 比如LIMIT FOK订单没有成交, 市价单没有完全成交, 强平期间被取消的订单, 交易所维护期间被取消的订单
)

type FilterType string

const (
	PRICE_FILTER FilterType = "PRICE_FILTER"
	LOT_SIZE     FilterType = "LOT_SIZE"
)

type EventKind string

const (
	Trade                   EventKind = "trade"
	ExecutionReport         EventKind = "executionReport"
	OutboundAccountPosition EventKind = "outboundAccountPosition"
)

const (
	DepthSuffix = "@depth5"
	TradeSuffix = "@trade"
)
