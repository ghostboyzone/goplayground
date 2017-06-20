/**
 * 计算一定时间区间内收益率
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"github.com/metakeule/fmtdate"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	myCoins map[string]([]interface{})
)

func main() {
	log.Println("Today is:", time.Now().Format("2006-01-02 15:04"))
	startStr := flag.String("start", "", "开始时间\n-start \"2017-06-15 00:00\"")
	endStr := flag.String("end", "", "结束时间\n-end \"2017-06-15 16:00\"")
	flag.Parse()

	if *startStr == "" || *endStr == "" {
		flag.Usage()
		os.Exit(0)
	}

	startT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm ZZ", *startStr+" +00:00")
	endT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm ZZ", *endStr+" +00:00")

	initCoins()

	t1 := getAll(startT)
	t2 := getAll(endT)
	cmpResult := cmpData(t1, t2)
	log.Println("Range:", *startStr, ",", *endStr)
	log.Println(cmpResult.Message)
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

func getAll(t time.Time) map[string]float64 {
	bt, err := db.InitBitCoin("data/coin_5m.dbdata", true)
	if err != nil {
		log.Fatal(err)
	}
	defer bt.Close()

	result := make(map[string]float64)
	uTime := t.In(time.UTC).Unix()
	for coinName, v := range myCoins {
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
		price, _ := strconv.ParseFloat(oneV[2].(string), 64)
		log.Println(coinName, v[0].(string), price)
		result[coinName] = price
	}
	return result
}

type CmpResult struct {
	Message string
	Coins   map[string](map[string]interface{})
}

func cmpData(before, after map[string]float64) CmpResult {
	coins := make(map[string](map[string]interface{}))
	totalCnt := 0
	totalRate := float64(0)
	for k, v := range after {
		if v == 0 || before[k] == 0 {
			log.Println("Zero Skip", k, before[k], v)
			continue
		}
		nowRate := 100 * (v - before[k]) / before[k]
		totalRate += nowRate
		totalCnt++
		tmp := make(map[string]interface{})
		tmp["before"] = before[k]
		tmp["after"] = v
		tmp["rate"] = nowRate
		log.Println(k, tmp, "total", totalRate, totalCnt)
		coins[k] = tmp
	}

	return CmpResult{
		Message: fmt.Sprintf("Result: %f%s, total: %d", totalRate/float64(totalCnt), "%", totalCnt),
		Coins:   coins,
	}
}
