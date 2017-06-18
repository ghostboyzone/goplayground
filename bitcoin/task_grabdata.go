/**
 * 拉取所有币的历史价格数据
 */
package main

import (
	"encoding/json"
	"fmt"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"log"
	"time"
)

var (
	unit_maps    []string
	bt           map[string](*db.BitCoin)
	visitTimeMap map[int64]bool
)

func main() {
	// 5m  30m  1h  8h  1d
	initDb()

	// write all coin
	totalResult := apiReq.AllCoin()
	totalResultJsonBt, _ := json.Marshal(totalResult)
	bt["all"].Set("all_coin", string(totalResultJsonBt))

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
			for _, unit := range unit_maps {
				writeKData(coinName, unit)
			}
		}
		for _, unit := range unit_maps {
			bt[unit].Shrink()
		}
		time.Sleep(time.Minute * 2)

	}
}

func writeKData(coinName string, unit string) {
	tmpData := apiReq.AdvanceKData(coinName, unit)
	for _, tmpV := range tmpData {
		tmpK := int64(tmpV[0].(float64) / 1000)
		tmpJsonBt1, _ := json.Marshal(tmpV)
		cKey := fmt.Sprintf("coin:%s:%d", coinName, tmpK)
		log.Println(unit, coinName, "start", tmpK, cKey, string(tmpJsonBt1))
		errW := bt[unit].Set(cKey, string(tmpJsonBt1))
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

func initDb() {
	unit_maps = []string{"5m", "30m", "1h", "8h", "1d"}
	bt = make(map[string](*db.BitCoin))
	var err error
	bt["all"], err = db.InitBitCoin("data/allcoins.dbdata", false)
	if err != nil {
		log.Fatal(err)
	}

	for _, unit := range unit_maps {
		bt[unit], err = db.InitBitCoin("data/coin_"+unit+".dbdata", false)
		if err != nil {
			log.Fatal(err)
		}
	}
}
