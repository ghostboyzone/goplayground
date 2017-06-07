package json

type BdTranslateJson struct {
	From        string              `json:"from"`
	To          string              `json:"to"`
	TransResult []BdTranslateResult `json:"trans_result"`
}

type BdTranslateResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}
