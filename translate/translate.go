package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	resultJson "github.com/ghostboyzone/goplayground/translate/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// baidu translate api key
const (
	API_URL    = "http://api.fanyi.baidu.com/api/trans/vip/translate"
	APP_ID     = "YOUR_APP_ID"
	SECRET_KEY = "YOUR_SECRET_KEY"
)

func main() {
	words := getWordsFromFile()
	for _, word := range words {
		transWord, err := getTranslateResult(word)
		// log.Println(transWord, err)
		if err == nil {
			result := formatTransResult(transWord)
			// log.Println(result.From, result.To, result.TransResult)
			if len(result.TransResult) < 1 {
				continue
			}
			log.Println("src:", result.TransResult[0].Src, ", dst:", result.TransResult[0].Dst)
		}
	}

}

func getTranslateResult(originWord string) (transWord string, err error) {
	salt := strconv.FormatInt(time.Now().Unix(), 10)
	v := url.Values{}
	v.Add("q", originWord)
	v.Add("from", "en")
	v.Add("to", "zh")
	v.Add("appid", APP_ID)
	v.Add("salt", salt)
	v.Add("sign", getSign(originWord, salt))
	reqUrl := API_URL + "?" + v.Encode()
	// log.Println("Request URI: ", reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

/**
 * format string to json
 * @param  {[type]} data string)       (bdTranslateJson resultJson.BdTranslateJson [description]
 * @return {[type]}      [description]
 */
func formatTransResult(data string) (bdTranslateJson resultJson.BdTranslateJson) {
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.Decode(&bdTranslateJson)
	return bdTranslateJson
}

/**
 * cal sign
 */
func getSign(q string, salt string) string {
	data := []byte(APP_ID + q + salt + SECRET_KEY)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func getWordsFromFile() (words []string) {
	f, _ := os.Open("dict.txt")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineStr := scanner.Text()
		list := strings.Split(lineStr, ",")
		if len(list) < 2 {
			continue
		}
		realWord := strings.Trim(list[1], "\"")
		words = append(words, realWord)
	}
	return words
}
