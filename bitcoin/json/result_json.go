package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CoinTrends struct {
	Data   []DataArr `json:"data"`
	Yprice float64   `json:"yprice"`
}
type DataArr [2]float64
type CoinHashes map[string]CoinTrends

type CoinLatest struct {
	High   string  `json:"high"`
	Low    string  `json:"low"`
	Buy    string  `json:"buy"`
	Sell   string  `json:"sell"`
	Last   string  `json:"last"`
	Vol    float64 `json:"vol"`
	Volume float64 `json:"volume"`
}

type CoinLatestNew struct {
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Buy    float64 `json:"buy"`
	Sell   float64 `json:"sell"`
	Last   float64 `json:"last"`
	Vol    float64 `json:"vol"`
	Volume float64 `json:"volume"`
}

func (c CoinLatest) String() string {
	return fmt.Sprintf("High[%s] Low[%s] Buy[%s] Sell[%s] Last[%s] Vol[%s] Volume[%s]", c.High, c.Low, c.Buy, c.Sell, c.Last, c.Vol, c.Volume)
}

type CoinJs struct {
	TimeLine CoinJsDetail `json:"time_line"`
}

type CoinJsDetail struct {
	FiveM    [][]interface{} `json:"5m"`
	FifteenM [][]interface{} `json:"15m"`
	ThirtyM  [][]interface{} `json:"30m"`
	OneH     [][]interface{} `json:"1h"`
	EightH   [][]interface{} `json:"8h"`
	OneD     [][]interface{} `json:"1d"`
}

type CoinKUnit struct {
	Timestamp int64
	Amount    float64
	Open      float64
	High      float64
	Low       float64
	Close     float64
}

func StringToFloat64(value string) (f float64) {
	f, _ = strconv.ParseFloat(value, 64)
	return
}

func InterfaceToFloat64(value interface{}) (f float64) {
	switch value.(type) {
	case float64:
		f = value.(float64)
	case string:
		f = StringToFloat64(value.(string))
	}
	return
}

func FormatCoinKUnit(value []interface{}) (coinKUnit CoinKUnit, err error) {
	if len(value) != 6 {
		err = errors.New("Invalid value, no enough length")
		return
	}

	return CoinKUnit{
		Timestamp: int64(value[0].(float64) / 1000),
		Amount:    InterfaceToFloat64(value[1]),
		Open:      InterfaceToFloat64(value[2]),
		High:      InterfaceToFloat64(value[3]),
		Low:       InterfaceToFloat64(value[4]),
		Close:     InterfaceToFloat64(value[5]),
	}, nil
}

func FormatCoinKUnitByString(value string) (coinKUnit CoinKUnit, err error) {
	var oneV []interface{}
	decoder := json.NewDecoder(strings.NewReader(string(value)))
	decoder.Decode(&oneV)
	coinKUnit, err = FormatCoinKUnit(oneV)
	return
}
