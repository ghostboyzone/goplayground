/**
 * 订单查询
 */
package main

import (
	"encoding/json"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"github.com/metakeule/fmtdate"
	"log"
	// "os"
	"github.com/fatih/color"
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
	startT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm ZZ", "2017-06-29 10:40 +08:00")
	// 查询开始时间
	startTimestamp := startT.Unix()
	// 如果为0，表示历史所有挂单
	startTimestamp = 0

	maxCh := make(chan int, 20)

	for coinName, v1 := range myCoins {
		maxCh <- 1

		time.Sleep(time.Millisecond * 30)

		go func(coinName string, coinCName string) {
			tmp := apiReq.TradeList(coinName, startTimestamp)
			log.Println(coinName, tmp)
			if len(tmp) != 0 {
				for _, v2 := range tmp {
					v2["coin_name"] = coinName
					v2["coin_cname"] = coinCName
					orderDetail := apiReq.TradeView(coinName, v2["id"].(string))
					v2["status"] = orderDetail["status"]
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
		if a["type"].(string) == "sell" {
			color.Set(color.FgGreen, color.Bold)
		} else {
			color.Set(color.FgRed, color.Bold)
		}
		log.Println(a["id"], "\t", a["coin_name"], "\t", a["coin_cname"], "\t", a["datetime"], "\t", a["type"], "\t", a["status"], "\tprice:", a["price"], ", total:", a["amount_original"], ", left:", a["amount_outstanding"], ", amount:", a["amount_original"].(float64)-a["amount_outstanding"].(float64))
	}
	color.Unset()
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
