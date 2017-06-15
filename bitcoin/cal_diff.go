package main

import (
	"encoding/json"
	"fmt"
	// "github.com/boltdb/bolt"
	// apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	wechatApi "github.com/ghostboyzone/goplayground/wechat/api"
	"github.com/metakeule/fmtdate"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	myCoins    map[string]([]interface{})
	bt         *db.BitCoin
	wechat     *wechatApi.Wechat
	toUserName string
	nickName   string
)

func main() {
	wechat = wechatApi.NewWechat()
	wechat.ShowQrCode()
	wechat.WaitForScan()
	wechat.GetContact()

	nickName = "filehelper"
	toUserName = "filehelper"
	fp, _ := os.OpenFile("contact.mlog", os.O_CREATE|os.O_WRONLY, 0600)
	defer fp.Close()
	for _, v := range wechat.MemberList {
		toBeWrite := fmt.Sprintf("NickName[%s] RemarkName[%s] UserName[%s]\n\n", v["NickName"].(string), v["RemarkName"].(string), v["UserName"].(string))
		fp.WriteString(toBeWrite)
		if v["RemarkName"].(string) == nickName || v["NickName"].(string) == nickName {
			toUserName = v["UserName"].(string)
		}
	}

	var err error
	bt, err = db.InitBitCoin("bitcoin.dbdata", true)
	if err != nil {
		log.Fatal(err)
	}

	min := 10

	for {

		bt.Load("bitcoin.dbdata")
		bt.CreateIndex("coin_data", "coin:*")

		initCoins()

		endStr := fmt.Sprintf("2017-06-15 13:%d:00 +00:00", min)

		wechat.SendMsg(endStr, toUserName)

		// fmtdate.Parse("", time.Now().Format("2006-01-02 15:04:05"))
		startT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", "2017-06-15 00:00:00 +00:00")
		endT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", endStr)
		t1 := getAll(startT)
		t2 := getAll(endT)

		cmpData(t1, t2)

		time.Sleep(5 * time.Minute)
		min += 5

	}

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
	tmpStr := fmt.Sprintf("Result: %f%s, total: %d", totalRate/float64(totalCnt), "%", totalCnt)
	// log.Println("Result: ", totalRate/float64(totalCnt), "%, total:", totalCnt)
	log.Println(tmpStr)
	wechat.SendMsg(tmpStr, toUserName)
}
