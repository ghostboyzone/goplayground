package main

import (
	"fmt"
	"github.com/fatih/color"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"log"
	"time"
)

type BuyedCoin struct {
	Price    float64
	NowPrice float64
	Amount   float64
}

var (
	buyedCoins map[string]BuyedCoin
)

func main() {

	buyedCoins = make(map[string]BuyedCoin)

	for {
		buyedCoins["bts"] = BuyedCoin{
			Price:    2.18,
			Amount:   2997,
			NowPrice: apiReq.Ticker("bts").Buy,
		}

		for coinName, v := range buyedCoins {
			rate := 100 * (v.NowPrice - v.Price) / v.Price
			earnMoney := (v.NowPrice - v.Price) * v.Amount

			str := fmt.Sprintf("coinName[%s] price[%f] nowPrice[%f] rate[%f%s] amount[%f] earnMoney[%f]", coinName, v.Price, v.NowPrice, rate, "%", v.Amount, earnMoney)

			if rate > 0 {
				color.Set(color.FgGreen, color.Bold)
			} else {
				color.Set(color.FgRed, color.Bold)
			}

			log.Println(str)
		}
		color.Unset()
		time.Sleep(time.Second * 15)
	}

}
