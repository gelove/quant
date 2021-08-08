package test

import (
	"quant/internal/app/constant"
	"quant/internal/app/logic"
	"quant/pkg/utils/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetAccount(t *testing.T) {
	res, err := logic.GetAccount()
	if err != nil {
		t.Logf("TestGetAccount err => %+v", err)
		return
	}
	t.Logf("TestGetAccount res => %s", json.MustToString(res))
}

// PlaceOrder
func TestPlaceOrder(t *testing.T) {
	symbols := logic.GetSymbols()
	if len(symbols) == 0 {
		t.Log("TestGetOpenOrders symbols 为空")
		return
	}
	err := logic.PlaceOrder(symbols[0], constant.BUY, decimal.NewFromInt(1), decimal.NewFromFloat(250.1))
	if err != nil {
		t.Logf("TestPlaceOrder err => %+v", err)
		return
	}
	t.Logf("TestPlaceOrder success")
}

func TestGetALLOrders(t *testing.T) {
	symbols := logic.GetSymbols()
	if len(symbols) == 0 {
		t.Log("TestGetOpenOrders symbols 为空")
		return
	}
	res, err := logic.GetALLOrders(symbols[0])
	if err != nil {
		t.Logf("TestGetOpenOrders err => %+v", err)
		return
	}
	t.Logf("TestGetOpenOrders res => %s", json.MustToString(res))
}

func TestGetOpenOrders(t *testing.T) {
	symbols := logic.GetSymbols()
	if len(symbols) == 0 {
		t.Log("TestGetOpenOrders symbols 为空")
		return
	}
	res, err := logic.GetOpenOrders(symbols[0])
	if err != nil {
		t.Logf("TestGetOpenOrders err => %+v", err)
		return
	}
	t.Logf("TestGetOpenOrders res => %s", json.MustToString(res))
}

func TestGetExchangeInfo(t *testing.T) {
	symbols := logic.GetSymbols()
	res, err := logic.GetExchangeInfo(symbols)
	if err != nil {
		t.Logf("TestGetExchangeInfo err => %+v", err)
		return
	}
	t.Logf("TestGetExchangeInfo res => %s", json.MustToString(res))
}
