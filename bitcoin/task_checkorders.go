package main

import (
	"encoding/json"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"log"
	// "os"
	"github.com/metakeule/fmtdate"
	"sort"
	"strings"
	"time"
)

var (
	myCoins   map[string]([]interface{})
	orderList []map[string]interface{}
)

func main() {
	initCoins()

	maxCh := make(chan int, 1)

	for coinName, v1 := range myCoins {
		maxCh <- 1

		go func(coinName string, coinCName string) {
			tmp := apiReq.TradeList(coinName)
			if len(tmp) != 0 {
				for _, v2 := range tmp {
					v2["coin_name"] = coinName
					v2["coin_cname"] = coinCName
					orderList = append(orderList, v2)

					tmpT, err := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", v2["datetime"].(string)+" +08:00")
					if err != nil {
						log.Fatal(err)
					}
					v2["timestamp"] = tmpT.Unix()
				}
			}

			<-maxCh
		}(coinName, v1[0].(string))
	}

	for {
		if len(maxCh) == 0 {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	sortMyOrders(orderList)
	// log.Println(orderList)
	for _, a := range orderList {
		log.Println(a)
	}
}

type ByDate []map[string]interface{}

func (b ByDate) Len() int      { return len(b) }
func (b ByDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// 耗时从低到高排序datetime
func (b ByDate) Less(i, j int) bool { return b[i]["timestamp"].(int64) > b[j]["timestamp"].(int64) }

func sortMyOrders(orderList []map[string]interface{}) {
	sort.Sort(ByDate(orderList))
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
