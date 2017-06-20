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
	buyedCoins["ifc"] = BuyedCoin{
		Price:    0.001152,
		Amount:   260456.2500333,
		NowPrice: apiReq.Ticker("ifc").Buy,
	}
	buyedCoins["ppc"] = BuyedCoin{
		Price:    16.6,
		Amount:   34.883082,
		NowPrice: apiReq.Ticker("ppc").Buy,
	}

	for {
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
		time.Sleep(time.Second * 5)
	}

}
