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
			Price:    1.86,
			Amount:   2997,
			NowPrice: apiReq.Ticker("bts").Buy,
		}
		buyedCoins["hlb"] = BuyedCoin{
			Price:    0.372,
			Amount:   9990,
			NowPrice: apiReq.Ticker("hlb").Buy,
		}

		var totalEarnMoney float64

		for coinName, v := range buyedCoins {
			rate := 100 * (v.NowPrice - v.Price) / v.Price
			earnMoney := (v.NowPrice - v.Price) * v.Amount

			str := fmt.Sprintf("coinName[%s] yprice[%f] buyPrice[%f] nowPrice[%f] rate[%f%s] amount[%f] earnMoney[%f]", coinName, apiReq.Yprice(coinName)["yprice"], v.Price, v.NowPrice, rate, "%", v.Amount, earnMoney)

			totalEarnMoney += earnMoney

			if rate > 0 {
				color.Set(color.FgGreen, color.Bold)
			} else {
				color.Set(color.FgRed, color.Bold)
			}

			log.Println(str)
		}

		if totalEarnMoney > 0 {
			color.Set(color.FgGreen, color.Bold)
		} else {
			color.Set(color.FgRed, color.Bold)
		}
		log.Println(fmt.Sprintf("asset[%f] totalEarnMoney[%f]", apiReq.Balance()["asset"], totalEarnMoney))
		color.Unset()
		time.Sleep(time.Second * 15)
	}

}
