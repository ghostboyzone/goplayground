package api

import (
	"encoding/json"
	resultJson "github.com/ghostboyzone/goplayground/bitcoin/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	coinJsReg *regexp.Regexp
)

/**
 * 牌价
 * https://www.jubi.com/help/api.html#three-one
 */
func Ticker(coinName string) (coin resultJson.CoinLatest) {
	v := url.Values{}
	v.Add("coin", coinName)
	decoder := json.NewDecoder(strings.NewReader(reqPublic("/api/v1/ticker", v)))
	decoder.Decode(&coin)
	return coin
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
func Orders(coinName string) [](map[string]interface{}) {
	v := url.Values{}
	v.Add("coin", coinName)
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

func reqPublic(api string, v url.Values) string {
	reqUrl := BASE_API_URL + api
	if len(v) != 0 {
		reqUrl += "?" + v.Encode()
	}
	log.Println("req_public", reqUrl)
	resp, _ := http.Get(reqUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func init() {
	if coinJsReg == nil {
		coinJsReg, _ = regexp.Compile(`^chart=(.*)(symbol)(.*)(symbol_view)(.*)(ask)(.*)(time_line)(.*);`)
	}
}
