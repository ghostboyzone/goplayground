package main

import (
	"encoding/json"
	resultJson "github.com/ghostboyzone/goplayground/bitcoin/json"
	"github.com/metakeule/fmtdate"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_API_URL = "http://www.jubi.com"
	API_TICKER   = "/api/v1/ticker"
	API_ALLCOIN  = "/coin/allcoin"
)

var (
	coinJsReg *regexp.Regexp
)

func main() {
	initReg()

	totalResult := getAllCoins()

	nowT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", "2017-06-12 00:00:00 +00:00")

	totalRound := 20
	for i := 1; i <= totalRound; i++ {
		beforeT := nowT.AddDate(0, 0, -1*i)

		nowTStr := nowT.In(time.UTC).Format("2006-01-02 15:04:05")
		beforeTStr := beforeT.In(time.UTC).Format("2006-01-02 15:04:05")

		log.Println("round: ", i)
		log.Println(beforeTStr, nowTStr)

		rate := float64(0)
		total := 0
		for k, v := range totalResult {
			resultMap := getFromJs(k)

			nowPrice := resultMap[nowTStr]

			beforePrice := resultMap[beforeTStr]

			if beforePrice == 0 {
				log.Println("Skip:", k, v[0])
				continue
			}
			nowRate := 100 * (nowPrice - beforePrice) / beforePrice

			rate += nowRate
			total++
			// log.Println(k, v[0], beforeTStr, beforePrice, nowTStr, nowPrice, nowRate, total, rate)
		}

		log.Println(rate/float64(total), total)
	}
}

func initReg() {
	coinJsReg, _ = regexp.Compile(`^chart=(.*)(symbol)(.*)(symbol_view)(.*)(ask)(.*)(time_line)(.*);`)
}

func getTrends() resultJson.CoinHashes {
	reqUrl := "https://www.jubi.com/coin/trends"
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	retResult := string(body)

	var totalResult resultJson.CoinHashes
	decoder := json.NewDecoder(strings.NewReader(retResult))
	decoder.Decode(&totalResult)
	return totalResult
}

func getLatest(coinName string) (coin resultJson.CoinLatest, err error) {
	v := url.Values{}
	v.Add("coin", coinName)
	reqUrl := BASE_API_URL + API_TICKER + "?" + v.Encode()

	resp, err := http.Get(reqUrl)
	if err != nil {
		return coin, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return coin, err
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.Decode(&coin)
	return coin, err
}

func getFromJs(coinName string) (result map[string]float64) {
	result = make(map[string]float64)
	reqUrl := BASE_API_URL + "/coin/" + coinName + "/k.js"
	resp, _ := http.Get(reqUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	str := coinJsReg.ReplaceAllString(string(body), "$1\"$2\"$3\"$4\"$5\"$6\"$7\"$8\"$9")

	var myCoinJs resultJson.CoinJs
	decoder := json.NewDecoder(strings.NewReader(str))
	decoder.Decode(&myCoinJs)

	for _, i := range myCoinJs.TimeLine.OneD {
		timestamp := int64(i[0].(float64) / 1000)
		timestamp_format := time.Unix(timestamp, 0).In(time.UTC).Format("2006-01-02 15:04:05")
		price, _ := strconv.ParseFloat(i[2].(string), 64)
		// log.Println(timestamp_format, price)
		result[timestamp_format] = price
	}
	return result
}

func getAllCoins() (result resultJson.CoinAllHashes) {
	result = make(resultJson.CoinAllHashes)
	reqUrl := BASE_API_URL + API_ALLCOIN
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// var myCoinJs resultJson.CoinJs
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.Decode(&result)
	return result
}
