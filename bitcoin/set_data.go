package main

import (
	"encoding/json"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"log"
	// "strconv"
	"fmt"
	// "time"
	// "os"
)

var (
	bt *db.BitCoin
)

func main() {
	var err error
	bt, err = db.InitBitCoin()
	if err != nil {
		log.Fatal(err)
	}

	// bt.CreateIndex("coin_data", "coin:*")

	// write all coin
	totalResult := apiReq.AllCoin()
	totalResultJsonBt, _ := json.Marshal(totalResult)
	bt.Set("all_coin", string(totalResultJsonBt))

	for k, v := range totalResult {
		fmt.Printf("key=%s, value=%s\n", k, v)
		coinName := string(k)
		writeKData(coinName, "5m")
		writeKData(coinName, "30m")
		writeKData(coinName, "1h")
		writeKData(coinName, "8h")
		writeKData(coinName, "1d")
	}

	bt.Shrink()
	bt.Close()
}

func writeKData(coinName string, unit string) {
	tmpData := apiReq.AdvanceKData(coinName, unit)
	for _, tmpV := range tmpData {
		tmpK := int64(tmpV[0].(float64) / 1000)
		tmpJsonBt1, _ := json.Marshal(tmpV)
		cKey := fmt.Sprintf("coin:%s:%d", coinName, tmpK)
		log.Println(unit, coinName, "start", tmpK, cKey, string(tmpJsonBt1))
		errW := bt.Set(cKey, string(tmpJsonBt1))
		if errW != nil {
			log.Println(errW)
		}
		log.Println(unit, coinName, "done", tmpK)
	}
}
