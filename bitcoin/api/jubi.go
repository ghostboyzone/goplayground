package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	nonceMap NonceMap
)

type NonceMap struct {
	Data map[int64]int64
	sync.RWMutex
}

func init() {
	nonceMap = NonceMap{
		Data: make(map[int64]int64),
	}
}

/**
 * 挂单查询
 * https://www.jubi.com/help/api.html#three-six
 */
func TradeList(coinName string) []map[string]interface{} {
	v := getCommonParams(coinName)
	v.Add("type", "all")
	v.Add("since", "0")
	var data []map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(req("trade_list", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 查询订单信息
 * https://www.jubi.com/help/api.html#three-seven
 */
func TradeView(coinName string, id string) map[string]interface{} {
	v := getCommonParams(coinName)
	v.Add("id", id)
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(req("trade_view", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 取消订单
 * https://www.jubi.com/help/api.html#three-eight
 */
func TradeCancel(coinName string, id string) map[string]interface{} {
	v := getCommonParams(coinName)
	v.Add("id", id)
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(req("trade_cancel", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 下单
 * @param type1: buy or sell
 * https://www.jubi.com/help/api.html#three-nine
 */
func TradeAdd(coinName string, amount string, price string, type1 string) map[string]interface{} {
	v := getCommonParams(coinName)
	v.Add("amount", amount)
	v.Add("price", price)
	v.Add("type", type1)
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(req("trade_add", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 账户信息
 * https://www.jubi.com/help/api.html#three-four
 */
func Balance() map[string]interface{} {
	v := getCommonParams("")
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(req("balance", v)))
	decoder.Decode(&data)
	return data
}

/**
 * 比特币充值地址
 * https://www.jubi.com/help/api.html#three-five
 */
func Wallet(coinName string) map[string]interface{} {
	v := getCommonParams(coinName)
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(req("wallet", v)))
	decoder.Decode(&data)
	return data
}

func getCommonParams(coinName string) url.Values {
	v := url.Values{}
	if coinName != "" {
		v.Add("coin", coinName)
	}
	// sleep 1ms, 避免nonce冲突
	// time.Sleep(time.Millisecond * 1)
	millTimestamp := time.Now().UnixNano() / 1e6
	// log.Println(nonceMap[millTimestamp])
	nonceMap.Lock()
	if nonceMap.Data[millTimestamp] == 0 {
		nonceMap.Data[millTimestamp] = 0
	}
	nonceMap.Data[millTimestamp]++
	nonce := strconv.FormatInt(millTimestamp+nonceMap.Data[millTimestamp], 10)
	nonceMap.Unlock()
	v.Add("nonce", nonce)
	v.Add("key", PUBLIC_KEY)
	return v
}

func req(api string, v url.Values) string {
	v.Add("signature", sha256Sum2(v.Encode()))
	reqUrl := API_URL + api
	// log.Println("req", reqUrl, v.Encode())
	resp, _ := http.PostForm(reqUrl, v)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(body))
	return string(body)
}

func sha256Sum(str string) string {
	sha := sha256.New()
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func sha256Sum2(str string) string {
	mac := hmac.New(sha256.New, []byte(md5Sum(SECRET_KEY)))
	mac.Write([]byte(str))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func md5Sum(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
