package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	resultJson "github.com/ghostboyzone/goplayground/bitcoin/json"
	"github.com/tidwall/buntdb"
	"log"
	"strings"
	// "github.com/metakeule/fmtdate"
	"os"
	"strconv"
	"time"
)

var (
	myCoins map[string]([]interface{})
)

func main() {
	initCoins()

	preTimestamp := getPreTimestamp()
	log.Println(time.Unix(preTimestamp, 0).In(time.UTC).Format("2006-01-02 15:04"))
	bt, err := db.InitBitCoin("data/coin_1h.dbdata", true)
	defer bt.Close()
	if err != nil {
		log.Fatal(err)
	}

	preHours := int64(8)

	for coinName, v := range myCoins {
		// log.Println(coinName, v[0])

		totalAmount := float64(0)
		totalValue := float64(0)

		avgPrice := make(map[int64]float64)
		avgAmount := make(map[int64]float64)

		canBuyCnt := int64(0)
		for i := int64(1); i <= preHours; i++ {
			nowTimestamp := preTimestamp - (i-1)*3600
			cKey := "coin:" + coinName + ":" + strconv.FormatInt(nowTimestamp, 10)
			val, _ := bt.Get(cKey)
			kUnit, _ := resultJson.FormatCoinKUnitByString(val)

			if kUnit.Amount == 0 || kUnit.Close == 0 {
				log.Println(coinName, v[0], i, nowTimestamp, "Skip")
				continue
			}
			totalAmount += kUnit.Amount
			totalValue += kUnit.Amount * kUnit.Close

			if totalValue == 0 || totalAmount == 0 {
				log.Fatal(coinName, v[0], i, "Err Zero")
				continue
			}
			avgPrice[i] = totalValue / totalAmount
			avgAmount[i] = totalAmount / float64(i)

			if coinName == "max" {

				log.Println(coinName, v[0], i, avgPrice[i], avgAmount[i])
			}

			if i > 1 {
				if avgPrice[i] != 0 && avgAmount[i] != 0 && avgPrice[i] > avgPrice[i-1] && avgAmount[i] < avgAmount[i-1] {
					canBuyCnt++
				}
			}
		}

		if canBuyCnt == preHours-1 {
			color.Set(color.FgGreen, color.Bold)
			log.Println(coinName, v[0], "can buy", avgPrice[1], avgAmount[1])
		} else {
			color.Set(color.FgRed, color.Bold)
			log.Println(coinName, v[0], "can not buy")
		}
		color.Unset()
	}

	// 1h 2h 4h 8h 12h 18h 24h

	// 5 10 20 30
	os.Exit(0)

	bt.CreateIndex("coin_data_doge", "coin:doge:*")

	bt.Db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("coin_data_doge", func(key, val string) bool {
			fmt.Printf("%s %s\n", key, val)

			kUnit, err := resultJson.FormatCoinKUnitByString(val)
			if err != nil {
				log.Println(err)
				return false
			}
			log.Println(time.Unix(kUnit.Timestamp, 0).In(time.UTC).Format("2006-01-02 15:04"))

			return true
		})
		return nil
	})

	// cKey := fmt.Sprintf("coin:%s:%d", coinName, getTodayZero().In(time.UTC).Unix())
	// resultTmp, err := bt.Get(cKey)
}

func getPreTimestamp() int64 {
	timestamp := time.Now().Unix()
	timestamp += 8 * 3600
	newTimestamp := timestamp - timestamp%3600
	newTimestamp -= 3600
	return newTimestamp
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
