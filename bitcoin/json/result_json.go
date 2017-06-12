package json

type CoinTrends struct {
	Data   []DataArr `json:"data"`
	Yprice float64   `json:"yprice"`
}
type DataArr [2]float64
type CoinHashes map[string]CoinTrends

type CoinAllHashes map[string]CoinAllHashesDetail
type CoinAllHashesDetail []interface{}

type CoinLatest struct {
	High   string  `json:"high"`
	Low    string  `json:"low"`
	Buy    string  `json:"buy"`
	Sell   string  `json:"sell"`
	Last   string  `json:"last"`
	Vol    float64 `json:"vol"`
	Volume float64 `json:"volume"`
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
