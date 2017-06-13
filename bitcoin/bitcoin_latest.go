package main

import (
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	// resultJson "github.com/ghostboyzone/goplayground/bitcoin/json"
	"github.com/metakeule/fmtdate"
	"log"
	// "os"
	"strconv"
	"time"
)

func main() {
	totalResult := apiReq.AllCoin()

	nowT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", "2017-06-13 00:00:00 +00:00")

	currRate := float64(0)
	currTotal := 0
	for k, v := range totalResult {
		coin := apiReq.Ticker(k)

		currentPrice, _ := strconv.ParseFloat(coin.Last, 64)

		log.Println(k, v[0], currentPrice)

		resultMap := getKDataMap(k)
		nowTStr := nowT.In(time.UTC).Format("2006-01-02 15:04:05")
		nowPrice := resultMap[nowTStr]

		if nowPrice == 0 {
			log.Println("Curr Skip:", k, v[0])
			continue
		}
		nowRate := (currentPrice - nowPrice) / nowPrice
		currRate += nowRate
		currTotal++
		log.Println(nowPrice, currTotal, nowRate)
	}

	log.Println("Curr: ", currRate/float64(currTotal), currTotal)
}

func getKDataMap(coinName string) (result map[string]float64) {
	result = make(map[string]float64)
	myCoinJs := apiReq.AdvanceKData(coinName, "1d")
	for _, i := range myCoinJs {
		timestamp := int64(i[0].(float64) / 1000)
		timestamp_format := time.Unix(timestamp, 0).In(time.UTC).Format("2006-01-02 15:04:05")
		price, _ := strconv.ParseFloat(i[2].(string), 64)
		// log.Println(timestamp_format, price)
		result[timestamp_format] = price
	}
	return result
}
