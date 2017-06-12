package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

/*
 usage: first fulfill the public key and secret key, visit: https://www.jubi.com/help/api.html
*/

const (
	API_URL    = "http://www.jubi.com/api/v1/"
	PUBLIC_KEY = "YOUR PUBLIC KEY"
	SECRET_KEY = "YOUR SECRET KEY"
)

func TradeList(coinName string) {
	v := getCommonParams(coinName)
	v.Add("type", "all")
	v.Add("since", "0")
	log.Println(req("trade_list", v))
}

func TradeView(coinName string, id string) {
	v := getCommonParams(coinName)
	v.Add("id", id)
	log.Println(req("trade_view", v))
}

func TradeCancel(coinName string, id string) {
	v := getCommonParams(coinName)
	v.Add("id", id)
	log.Println(req("trade_cancel", v))
}

func TradeAdd(coinName string, amount string, price string, type1 string) {
	v := getCommonParams(coinName)
	v.Add("amount", amount)
	v.Add("price", price)
	v.Add("type", type1)
	log.Println(req("trade_add", v))
}

func Balance(coinName string) {
	v := getCommonParams(coinName)
	log.Println(req("balance", v))
}

func Wallet(coinName string) {
	v := getCommonParams(coinName)
	log.Println(req("wallet", v))
}

func getCommonParams(coinName string) url.Values {
	v := url.Values{}
	v.Add("coin", coinName)
	nonce := strconv.FormatInt(time.Now().Unix()*100000, 10)
	v.Add("nonce", nonce)
	v.Add("key", PUBLIC_KEY)
	return v
}

func req(api string, v url.Values) string {
	v.Add("signature", sha256Sum2(v.Encode()))
	reqUrl := API_URL + api
	log.Println("req", reqUrl, v.Encode())
	resp, _ := http.PostForm(reqUrl, v)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
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
