package app

import (
	"embed"
	"html/template"
	"log"
	"os"
	"os/signal"
	"quant/internal/app/config"
	"quant/internal/app/constant"
	"quant/internal/app/logic"
	"quant/pkg/socket"
	"quant/pkg/utils/logger"
	"strings"

	"github.com/pkg/errors"
)

//go:embed templates
var fs embed.FS

func init() {
	temp, err := template.ParseFS(fs, "templates/*.html")
	if err != nil {
		panic(errors.WithStack(err))
	}
	logic.Temp = temp
}

func Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	err := logic.ExchangeInfo()
	if err != nil {
		panic(err)
	}

	err = logic.Run()
	if err != nil {
		panic(err)
	}

	go logic.ProcessStream()

	listenKey, err := logic.GetListenKey()
	if err != nil {
		panic(err)
	}

	url := config.C.Binance.BaseWSUrl + "stream?streams=" + listenKey
	list := make([]string, 0, 1<<3)
	symbols := logic.GetSymbols()
	for _, v := range symbols {
		symbol := strings.ToLower(v)
		list = append(list, symbol+constant.DepthSuffix, symbol+constant.TradeSuffix)
	}
	if len(list) > 0 {
		url += logic.StreamSeparator + strings.Join(list, logic.StreamSeparator)
	}
	// url := "wss://testnet.binance.vision/stream?streams=listenKey/bnbusdt@trade/eosusdt@trade"
	logger.Infof("Websocket url => %s", url)
	ws := socket.Run(url, logic.SuccessHandler)

	// err = logic.PlaceOrder(config.C.Binance.Symbol, constant.SELL, decimal.NewFromInt(1), decimal.NewFromFloat(310.1))
	// if err != nil {
	// 	panic(err)
	// }

	<-interrupt
	log.Println("interrupt")
	ws.Close(true)
}
