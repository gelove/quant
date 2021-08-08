package config

// C Configuration
var C *Configuration

// Configuration is stuff that can be configured externally per env variables or config file (config.yml).
type Configuration struct {
	Binance Binance
}

type Binance struct {
	ApiKey        string
	SecretKey     string
	BaseUrl       string
	BaseWSUrl     string
	QuoteAsset    string
	Assets        []string
	Point         int64
	Amount        int64
	MaxOpenOrders int64
}
