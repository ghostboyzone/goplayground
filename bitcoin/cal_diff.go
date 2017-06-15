package main

import (
	"encoding/json"
	"fmt"
	// "github.com/boltdb/bolt"
	// apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"github.com/metakeule/fmtdate"
	"log"
	// "os"
	"strconv"
	"strings"
	"time"
)

var (
	myCoins map[string]([]interface{})
	bt      *db.BitCoin
)

func main() {
	var err error
	bt, err = db.InitBitCoin()
	if err != nil {
		log.Fatal(err)
	}

	bt.CreateIndex("coin_data", "coin:*")

	tt, _ := bt.Get("coin:lsk:1497324600")

	log.Println(tt)

	startT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", "2017-06-15 00:00:00 +00:00")
	endT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", "2017-06-15 09:30:00 +00:00")

	initCoins()
	// log.Println(myCoins)
	// log.Println(nowT.In(time.UTC).Unix())
	t1 := getAll(startT)
	t2 := getAll(endT)
	// log.Println(t1, t2, t1["aaa"])
	cmpData(t1, t2)

	bt.Close()
}

func initCoins() {
	v, err := bt.Get("all_coin")
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(strings.NewReader(v))
	decoder.Decode(&myCoins)
}

func getAll(t time.Time) map[string]float64 {
	result := make(map[string]float64)
	uTime := t.In(time.UTC).Unix()
	for coinName, _ := range myCoins {
		cKey := fmt.Sprintf("coin:%s:%d", coinName, uTime)
		resultTmp, err := bt.Get(cKey)
		if err != nil {
			log.Println("coin", coinName, err)
			continue
		}
		if len(resultTmp) == 0 {
			log.Println("empty")
			continue
		}

		var oneV []interface{}
		decoder := json.NewDecoder(strings.NewReader(string(resultTmp)))
		decoder.Decode(&oneV)
		if len(oneV) == 0 {
			log.Println("empty")
			continue
		}
		// log.Println(string(resultTmp), oneV)
		price, _ := strconv.ParseFloat(oneV[2].(string), 64)
		log.Println(coinName, price)
		result[coinName] = price
	}
	return result
}

func cmpData(before, after map[string]float64) {
	totalCnt := 0
	totalRate := float64(0)
	for k, v := range after {
		if before[k] == 0 {
			continue
		}
		nowRate := 100 * (v - before[k]) / before[k]
		totalRate += nowRate
		totalCnt++
		log.Println(k, "now", nowRate, "total", totalRate, totalCnt)
	}
	log.Println("Result: ", totalRate/float64(totalCnt), "%, total:", totalCnt)
}
