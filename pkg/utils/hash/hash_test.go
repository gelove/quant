package hash

import (
	"quant/internal/app/constant"
	"quant/pkg/utils/request"
	"testing"

	"github.com/imroc/req"
)

func TestSHA256String(t *testing.T) {
	params := []byte("symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559")
	mac := []byte("NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j")
	// sign := []byte("c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71")
	rs := SHA256MAC(params, mac)
	t.Logf("TestSHA256String rs => %s", rs)
}

func TestHashWithParams(t *testing.T) {
	params := req.Param{
		"symbol":      "LTCBTC",
		"side":        "BUY",
		"type":        constant.LIMIT,
		"timeInForce": "GTC",
		"quantity":    1,
		"price":       0.1,
		"recvWindow":  5000,
		"timestamp":   1499827319559,
	}
	for k, v := range params {
		t.Logf("TestSHA256String %s => %v", k, v)
	}
	mac := []byte("NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j")
	str := request.ConvertUrlValue(params).Encode()
	t.Logf("TestSHA256String str => %s", str)
	// print("price=0.1&quantity=1&recvWindow=5000&side=BUY&symbol=LTCBTC&timeInForce=GTC&timestamp=1499827319559&type=LIMIT")
	rs := SHA256MAC([]byte(str), []byte(mac))
	t.Logf("TestSHA256String rs => %s", rs)
}
