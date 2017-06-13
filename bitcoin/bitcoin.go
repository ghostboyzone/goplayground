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
			resultMap := getKDataMap(k)

			nowPrice := resultMap[nowTStr]

			beforePrice := resultMap[beforeTStr]

			if beforePrice == 0 {
				log.Println("Skip:", k, v[0])
				continue
			}
			nowRate := (nowPrice - beforePrice) / beforePrice

			rate += nowRate
			total++
			// log.Println(k, v[0], beforeTStr, beforePrice, nowTStr, nowPrice, nowRate, total, rate)
		}

		log.Println(rate/float64(total), total)
	}
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
