package api

import (
	"encoding/json"
	"fmt"
	resultJson "github.com/ghostboyzone/goplayground/bitcoin/json"
	"io/ioutil"
	// "log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	coinJsReg *regexp.Regexp
)

/**
 * 所有牌价
 * https://www.jubi.com/help/api.html#three-four
 */
func AllTicker() (coinMap map[string](resultJson.CoinLatestNew)) {
	v := url.Values{}
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/api/v1/allticker", v)))
	decoder.Decode(&coinMap)
	return coinMap
}

/**
 * 牌价
 * https://www.jubi.com/help/api.html#three-one
 */
func Ticker(coinName string) resultJson.CoinLatestNew {
	var coin resultJson.CoinLatest
	v := url.Values{}
	v.Add("coin", coinName)
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/api/v1/ticker", v)))
	decoder.Decode(&coin)
	return formatCoinLatest(coin)
}

func formatCoinLatest(coin resultJson.CoinLatest) (newCoin resultJson.CoinLatestNew) {
	return resultJson.CoinLatestNew{
		High:   resultJson.StringToFloat64(coin.High),
		Low:    resultJson.StringToFloat64(coin.Low),
		Buy:    resultJson.StringToFloat64(coin.Buy),
		Sell:   resultJson.StringToFloat64(coin.Sell),
		Last:   resultJson.StringToFloat64(coin.Last),
		Vol:    coin.Vol,
		Volume: coin.Volume,
	}
}

/**
 * 今日开盘价接口
 */
func Yprice(coinName string) map[string]interface{} {
	v := url.Values{}
	v.Add("coin", coinName)
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/api/v1/yprice", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 市场深度
 * https://www.jubi.com/help/api.html#three-two
 */
func Depth(coinName string) map[string]([][]interface{}) {
	v := url.Values{}
	v.Add("coin", coinName)
	var data map[string]([][]interface{})
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/api/v1/depth", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 市场交易
 * https://www.jubi.com/help/api.html#three-three
 */
func Orders(coinName string, since int64) [](map[string]interface{}) {
	v := url.Values{}
	v.Add("coin", coinName)
	v.Add("since", strconv.FormatInt(since, 10))
	var data [](map[string]interface{})
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/api/v1/orders", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 趋势
 * https://www.jubi.com/coin/trends
 */
func Trends() resultJson.CoinHashes {
	v := url.Values{}
	var data resultJson.CoinHashes
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/coin/trends", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 所有币的信息
 * (币的中文名、实时价格、买1、卖1、最高价、最低价、成交量、成交额)
 * https://www.jubi.com/coin/allcoin
 */
func AllCoin() map[string]([]interface{}) {
	v := url.Values{}
	var data map[string]([]interface{})
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/coin/allcoin", v)))
	decoder.Decode(&data)
	return data
}

/**
 * K线图
 * https://www.jubi.com/coin/doge/k.js
 */
func KData(coinName string) resultJson.CoinJs {
	v := url.Values{}
	respStr := reqPublic("/coin/"+coinName+"/k.js", v)
	respStr = coinJsReg.ReplaceAllString(respStr, "$1\"$2\"$3\"$4\"$5\"$6\"$7\"$8\"$9")
	var data resultJson.CoinJs
	decoder := json.NewDecoder(strings.NewReader(respStr))
	decoder.Decode(&data)
	return data
}

/**
 * 加强版K线图
 * https://www.jubi.com/coin/doge/k_5m.json
 *
 * 时间戳，成交量，open，high，low，close
 *
 * @param  coinName string        币名称，如: btc
 * @param  unit     string        时间间隔单位，如: 5m、15m、30m、1h、8h、1d
 * @return          [description]
 */
func AdvanceKData(coinName string, unit string) [][]interface{} {
	v := url.Values{}
	respStr := reqPublic(fmt.Sprintf("/coin/%s/k_%s.json", coinName, unit), v)
	var data [][]interface{}
	decoder := json.NewDecoder(strings.NewReader(respStr))
	decoder.Decode(&data)
	return data
}

func reqPublic(api string, v url.Values) string {
	reqUrl := BASE_API_URL + api
	if len(v) != 0 {
		reqUrl += "?" + v.Encode()
	}
	// log.Println("req_public", reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func init() {
	if coinJsReg == nil {
		coinJsReg, _ = regexp.Compile(`^chart=(.*)(symbol)(.*)(symbol_view)(.*)(ask)(.*)(time_line)(.*);`)
	}
}
