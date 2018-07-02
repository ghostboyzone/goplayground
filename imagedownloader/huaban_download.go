/*
* @Author: wenhao.ma
* @Date:   2018-07-02 16:52:40
* @Last Modified by:   wenhao.ma
* @Last Modified time: 2018-07-02 17:35:23
 */
package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type HBResult struct {
	Filter string  `json:"filter"`
	Pins   []HBPin `json:"pins"`
}
type HBPin struct {
	PinId int       `json:"pin_id"`
	File  HBPinFile `json:"file"`
}
type HBPinFile struct {
	Id   int    `json:"id"`
	Key  string `json:"key"`
	Type string `json:"type"`
}

func main() {
	total := 0
	maxStr := "1731606000"
	for {
		urlV := url.Values{}
		urlV.Add("max", maxStr)
		urlV.Add("limit", "20")
		urlV.Add("wfl", "1")
		reqUrl := "http://huaban.com/all" + "?" + urlV.Encode()
		req, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			log.Panicln(err)
		}
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.Header.Add("X-Request", "JSON")
		req.Header.Add("Referer", "http://huaban.com/all")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Panicln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Panicln(err)
		}

		var hbResult HBResult

		err = json.Unmarshal(body, &hbResult)
		if err != nil {
			log.Panicln(err)
		}
		pinLen := len(hbResult.Pins)
		if pinLen == 0 {
			log.Panicln("Len Err")
		}

		for _, tPin := range hbResult.Pins {
			imgUrl := "http://img.hb.aicdn.com/" + tPin.File.Key
			DownloadImage(imgUrl, "tmp1")
			total++
			log.Println(imgUrl, total)
		}
		log.Println(hbResult.Pins[pinLen-1].PinId)
		maxStr = strconv.Itoa(hbResult.Pins[pinLen-1].PinId)
	}

}

func DownloadImage(imgUrl string, path string) error {
	httpResp, httpErr := http.Get(imgUrl)
	if httpErr != nil {
		return httpErr
	}
	defer httpResp.Body.Close()
	os.MkdirAll(path, 0755)
	tmpFile, tmpFileErr := os.Create(path + "/" + getMd5FileName(imgUrl))
	if tmpFileErr != nil {
		return tmpFileErr
	}
	defer tmpFile.Close()
	_, cpErr := io.Copy(tmpFile, httpResp.Body)
	if cpErr != nil {
		return cpErr
	}
	return nil
}

func getMd5FileName(url string) string {
	data := []byte(url)
	return fmt.Sprintf("%x", md5.Sum(data)) + ".jpg"
}

// curl -H "X-Request: JSON" -H "X-Requested-With: XMLHttpRequest" -H "Referer: http://huaban.com/all" http://huaban.com/all\?jj40wapg\&max\=1731606000\&limit\=20\&wfl\=1
