package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quant/pkg/utils/logging"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var GraphqlApiUniSwapV2MainNet = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"
var GraphqlApiHecoMainNet = "https://graph.mdex.cc/subgraphs/name/mdex/swap"

var transactionGraphqlMainNetMiniV2 = "https://api.thegraph.com/subgraphs/name/noberk/chapter3"
var transactionGraphqlTestNetMiniV2 = "https://api.thegraph.com/subgraphs/name/noberk/chapter4"

const ethOneDayBlock = 5760
const hecoOneDayBlock = 5760 * 5
const ettTime = 15
const hecoTime = 3

func GetLpPrice(api, id string) (price decimal.Decimal, err error) {
	id = strings.ToLower(id)
	type rspData struct {
		Data struct {
			Pair struct {
				ID          string `json:"id"`
				ReserveUSD  string `json:"reserveUSD"`
				TotalSupply string `json:"totalSupply"`
			} `json:"pair"`
		} `json:"data"`
	}
	queryStr := `
            {
                pair(id:"%s") {
   				 	id
    				totalSupply
    				reserveUSD
  				}
            }
        `
	jsonData := map[string]string{
		"query": fmt.Sprintf(queryStr, id),
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	rsp := rspData{}
	err = getGraphqlData(api, jsonValue, &rsp)
	if err != nil {
		return
	}

	reserveUsdDec, err := decimal.NewFromString(rsp.Data.Pair.ReserveUSD)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	totalSupplyDec, err := decimal.NewFromString(rsp.Data.Pair.TotalSupply)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	price = reserveUsdDec.Div(totalSupplyDec).Mul(decimal.NewFromInt(2))
	return
}

type PairTokens struct {
	Token0Symbol  string
	Token0Address string
	Token1Symbol  string
	Token1Address string
}

func GetPairTokens(api, id string) (*PairTokens, error) {
	id = strings.ToLower(id)
	type rspData struct {
		Data struct {
			Pair struct {
				Token0 struct {
					ID     string `json:"id"`
					Symbol string `json:"symbol"`
				} `json:"token0"`
				Token1 struct {
					ID     string `json:"id"`
					Symbol string `json:"symbol"`
				} `json:"token1"`
			} `json:"pair"`
		} `json:"data"`
	}
	queryStr := `
            { 
                pair(id:"%s") {
					token0{
      					id
      					symbol
    				}
    				token1{
    					id
      					symbol
    				}
				}
            }
        `
	jsonData := map[string]string{
		"query": fmt.Sprintf(queryStr, id),
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	rsp := rspData{}
	err = getGraphqlData(api, jsonValue, &rsp)
	if err != nil {
		return nil, err
	}

	pairTokens := PairTokens{}
	pairTokens.Token0Address = rsp.Data.Pair.Token0.ID
	pairTokens.Token0Symbol = rsp.Data.Pair.Token0.Symbol
	pairTokens.Token1Address = rsp.Data.Pair.Token1.ID
	pairTokens.Token1Symbol = rsp.Data.Pair.Token1.Symbol

	return &pairTokens, nil
}

func GetTokenPrice(api, id string) (price decimal.Decimal, err error) {
	id = strings.ToLower(id)
	type rspData struct {
		Data struct {
			TokenDayData []struct {
				ID       string `json:"id"`
				PriceUSD string `json:"priceUSD"`
			} `json:"tokenDayDatas"`
		} `json:"data"`
	}

	queryStr := `
            { 
                tokenDayDatas(first:1,orderBy: date, orderDirection: desc,where:{token:"%s"}) {
   				 	id
    				priceUSD
  				} 
            }
        `
	jsonData := map[string]string{
		"query": fmt.Sprintf(queryStr, id),
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	rsp := rspData{}
	getGraphqlData(api, jsonValue, &rsp)

	if len(rsp.Data.TokenDayData) == 0 {
		err = errors.New("tokenDayData len is 0")
		return
	}
	price, err = decimal.NewFromString(rsp.Data.TokenDayData[0].PriceUSD)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// 价格变动
func GetTokenPriceChangeOneDay(api, id string, nowBlock uint64) (change decimal.Decimal, err error) {
	id = strings.ToLower(id)
	nowBlock -= 10
	priceNow, err := GetTokenPriceOnBlock(api, id, nowBlock)
	if err != nil {
		return
	}

	time.Sleep(100 * time.Millisecond)
	priceOneDay, err := GetTokenPriceOnBlock(api, id, nowBlock-hecoOneDayBlock)
	if err != nil {
		return
	}

	if priceOneDay.Cmp(decimal.Zero) <= 0 {
		err = errors.New(fmt.Sprintf("priceOneDay err %v", priceOneDay))
		return
	}
	change = priceNow.Sub(priceOneDay).Mul(decimal.NewFromInt(100)).Div(priceOneDay).Truncate(0)
	return
}

type TokenPrice struct {
	TimeStamp int64
	Price     decimal.Decimal
}

// 获取一天的价格
func GetTokenPriceListOneDay(api, id string, block uint64, now int64) ([]*TokenPrice, error) {
	id = strings.ToLower(id)
	retList := make([]*TokenPrice, 0)
	for i := 23; i >= 0; i-- {
		n := 0
		for {
			price, err := GetTokenPriceOnBlock(api, id, block-uint64(i)*3600/hecoTime)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				logging.ErrorF("GetTokenPriceOnBlock err %s", err)
				n++
				if n < 5 {
					continue
				} else {
					return nil, err
				}
			} else {
				tokenPrice := TokenPrice{now - 3600*int64(i), price}
				retList = append(retList, &tokenPrice)
				break
			}
		}
	}
	return retList, nil
}

func GetTokenPriceOnBlock(api, id string, block uint64) (price decimal.Decimal, err error) {
	id = strings.ToLower(id)
	type rspData struct {
		Data struct {
			Bundles []struct {
				EthPrice string `json:"ethPrice"`
			} `json:"bundles"`
			Tokens []struct {
				DerivedETH string `json:"derivedETH"`
			} `json:"tokens"`
		} `json:"data"`
	}

	queryStr := `
			{
  				bundles(block:{number:%d}){
    				ethPrice
  				}
  				tokens(where:{id:"%s"},block:{number:%d}){
    				derivedETH
  				}
			}`
	jsonData := map[string]string{
		"query": fmt.Sprintf(queryStr, block, id, block),
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return
	}
	rsp := rspData{}
	err = getGraphqlData(api, jsonValue, &rsp)
	if err != nil {
		return
	}

	if len(rsp.Data.Tokens) == 0 || len(rsp.Data.Bundles) == 0 {
		err = errors.New("rsp data len is 0")
		return
	}
	derivedEthDecimal, err := decimal.NewFromString(rsp.Data.Tokens[0].DerivedETH)
	if err != nil {
		return
	}
	ethPriceDecimal, err := decimal.NewFromString(rsp.Data.Bundles[0].EthPrice)
	if err != nil {
		return
	}

	price = derivedEthDecimal.Mul(ethPriceDecimal)
	return
}

func GetTokenPrice2(api, id string) (price decimal.Decimal, err error) {
	id = strings.ToLower(id)
	var queryPriceStr = `{
		pair(id:"%s"){
			id
			token1Price
		}
	}`

	jsonData := map[string]string{
		"query": fmt.Sprintf(queryPriceStr, id),
	}
	type Rsp struct {
		Data struct {
			Pair struct {
				Id          string `json:"id"`
				Token0Price string `json:"token1Price"`
			} `json:"pair"`
		} `json:"data"`
	}

	jsonValue, _ := json.Marshal(jsonData)
	var rsp Rsp

	err = getGraphqlData(api, jsonValue, &rsp)
	if err != nil {
		return
	}
	price, err = decimal.NewFromString(rsp.Data.Pair.Token0Price)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// 适用于其中之一是稳定币
func GetTokenPriceFromPair(api, tokenAddr, pairAddr string) (price decimal.Decimal, err error) {
	tokenAddr = strings.ToLower(tokenAddr)
	pairAddr = strings.ToLower(pairAddr)
	var queryPriceStr = `
					{	
						pair(id:"%s"){
  							id
    						token0{
      							id
      							symbol
    						}
    						token1{
      							id
      							symbol
    						}
    						token0Price
    						token1Price
  						}
					}`

	jsonData := map[string]string{
		"query": fmt.Sprintf(queryPriceStr, pairAddr),
	}

	type Rsp struct {
		Data struct {
			Pair struct {
				ID     string `json:"id"`
				Token0 struct {
					ID     string `json:"id"`
					Symbol string `json:"symbol"`
				} `json:"token0"`
				Token0Price string `json:"token0Price"`
				Token1      struct {
					ID     string `json:"id"`
					Symbol string `json:"symbol"`
				} `json:"token1"`
				Token1Price string `json:"token1Price"`
			} `json:"pair"`
		} `json:"data"`
	}

	jsonValue, _ := json.Marshal(jsonData)
	var rsp Rsp

	err = getGraphqlData(api, jsonValue, &rsp)
	if err != nil {
		return
	}
	var priceStr string
	if tokenAddr == strings.ToLower(rsp.Data.Pair.Token0.ID) {
		priceStr = rsp.Data.Pair.Token1Price
	} else {
		priceStr = rsp.Data.Pair.Token0Price
	}
	price, err = decimal.NewFromString(priceStr)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func getGraphqlData(api string, jsonValue []byte, retValue interface{}) error {
	request, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.WithStack(err)
	}

	client := &http.Client{Timeout: time.Second * 3}
	response, err := client.Do(request)
	if err != nil {
		return errors.WithStack(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ret status code %d", response.StatusCode))
	}

	rspData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.Unmarshal(rspData, retValue)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
