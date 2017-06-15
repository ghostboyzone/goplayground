/**
 * 拉取所有币的历史价格数据
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"log"
	"time"
)

var (
	bt           *db.BitCoin
	visitTimeMap map[int64]bool
)

func main() {
	// 5m  30m  1h  8h  1d
	unit := flag.String("unit", "5m", "-unit 5m")
	flag.Parse()
	log.Println("unit: ", *unit)
	var err error
	bt, err = db.InitBitCoin("bitcoin.dbdata", false)
	defer bt.Close()
	if err != nil {
		log.Fatal(err)
	}
	bt.CreateIndex("coin_data", "coin:*")

	// write all coin
	totalResult := apiReq.AllCoin()
	totalResultJsonBt, _ := json.Marshal(totalResult)
	bt.Set("all_coin", string(totalResultJsonBt))

	visitTimeMap = make(map[int64]bool)
	for {
		fmtTimestamp := formatTime5m(time.Now().Unix())
		if visitTimeMap[fmtTimestamp] {
			time.Sleep(time.Second * 30)
			continue
		}
		visitTimeMap[fmtTimestamp] = true
		log.Println("Grab Timestamp: ", fmtTimestamp)
		for k, v := range totalResult {
			coinName := string(k)
			log.Println("coin:", coinName, v[0].(string))
			writeKData(coinName, *unit)
		}
		time.Sleep(time.Minute * 1)
	}
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

/**
 * 将时间向一个粒度取整(5分钟)
 */
func formatTime5m(time int64) int64 {
	return time - time%300
}
