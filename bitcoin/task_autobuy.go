package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"github.com/metakeule/fmtdate"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	myCoins map[string]([]interface{})
)

func main() {
	initCoins()

	totalTimes := 1
	for {
		bt, err := db.InitBitCoin("data/coin_1d.dbdata", true)
		defer bt.Close()
		if err != nil {
			log.Fatal(err)
		}

		color.Unset()
		log.Println("T: ", totalTimes, "start")
		nowAllCoins := apiReq.AllTicker()
		// log.Println(nowAllCoins)
		for coinName, v := range myCoins {
			coinCName := v[0].(string)
			cKey := fmt.Sprintf("coin:%s:%d", coinName, getTodayZero().In(time.UTC).Unix())
			resultTmp, err := bt.Get(cKey)
			if err != nil {
				log.Println("err", coinName, err)
				continue
			}
			if len(resultTmp) == 0 {
				log.Println("empty", coinName)
				continue
			}
			var oneV []interface{}
			decoder := json.NewDecoder(strings.NewReader(string(resultTmp)))
			decoder.Decode(&oneV)
			if len(oneV) == 0 {
				log.Println("empty", coinName)
				continue
			}
			oldSellPrice, _ := strconv.ParseFloat(oneV[2].(string), 64)
			// oldSellPrice = oneV[4].(float64)
			newSellPrice := nowAllCoins[coinName].Sell
			if oldSellPrice == 0 {
				log.Println("Skip", coinName)
				continue
			}
			nowRate := 100 * (newSellPrice - oldSellPrice) / oldSellPrice
			if nowRate*-1 >= 30 {
				color.Set(color.FgRed, color.Bold)
				log.Println("[-30]", coinName, coinCName, oldSellPrice, newSellPrice, nowRate)
			}
			if nowRate*-1 >= 15 {
				color.Set(color.FgYellow, color.Bold)
				log.Println("[-15]", coinName, coinCName, oldSellPrice, newSellPrice, nowRate)
			}
			if nowRate*-1 >= 10 {
				color.Set(color.FgCyan, color.Bold)
				log.Println("[-10]", coinName, coinCName, oldSellPrice, newSellPrice, nowRate)
			}
			if nowRate >= 10 {
				color.Set(color.FgGreen, color.Bold)
				log.Println("[+10]", coinName, coinCName, oldSellPrice, newSellPrice, nowRate)
			}
		}
		totalTimes++
		defer color.Unset()
		log.Println("T: ", totalTimes, "done")
		time.Sleep(time.Second * 15)
	}

}

func getTodayZero() time.Time {
	zeroStr := time.Now().Format("2006-01-02") + " 00:00 +00:00"
	zeroT, err := fmtdate.Parse("YYYY-MM-DD hh:mm ZZ", zeroStr)
	if err != nil {
		log.Fatal(err)
	}
	return zeroT
}

func initCoins() {
	bt, err := db.InitBitCoin("data/allcoins.dbdata", true)
	if err != nil {
		log.Fatal(err)
	}
	defer bt.Close()
	v, err := bt.Get("all_coin")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(strings.NewReader(v))
	decoder.Decode(&myCoins)
}
