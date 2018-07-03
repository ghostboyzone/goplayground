/*
* @Author: wenhao.ma
* @Date:   2018-07-03 14:13:22
* @Last Modified by:   wenhao.ma
* @Last Modified time: 2018-07-03 14:30:51
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type TCResult struct {
	Result   string       `json:"result"`
	More     bool         `json:"more"`
	FeedList []TCFeedList `json:"feedList"`
}

type TCFeedList struct {
	PostId int       `json:"post_id"`
	Url    string    `json:"url"`
	Images []TCImage `json:"images"`
}

type TCImage struct {
	ImgId  int    `json:"img_id"`
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
}

func main() {

	fd, _ := os.OpenFile("tuchong_imglist.txt", os.O_RDWR|os.O_CREATE, 0755)
	defer fd.Close()

	total := 0
	nowPage := 1
	for {
		urlV := url.Values{}

		urlV.Add("os_api", "22")
		urlV.Add("device_type", "MI")
		urlV.Add("device_platform", "android")
		urlV.Add("ssmix", "a")
		urlV.Add("manifest_version_code", "232")
		urlV.Add("dpi", "400")
		urlV.Add("abflag", "0")
		urlV.Add("uuid", "651384659521356")
		urlV.Add("version_code", "232")
		urlV.Add("app_name", "tuchong")
		urlV.Add("version_name", "2.3.2")
		urlV.Add("openudid", "65143269dafd1f3a5")
		urlV.Add("resolution", "1280*1000")
		urlV.Add("os_version", "5.8.1")
		urlV.Add("ac", "wifi")
		urlV.Add("aid", "0")
		urlV.Add("page", strconv.Itoa(nowPage))
		urlV.Add("type", "refresh")

		reqUrl := "https://api.tuchong.com/feed-app" + "?" + urlV.Encode()
		req, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			log.Panicln(err)
		}
		req.Header.Add("Referer", "http://tuchong.com")
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

		var tcResult TCResult

		err = json.Unmarshal(body, &tcResult)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(tcResult)

		if tcResult.Result == "SUCCESS" {
			for _, tcFeedList := range tcResult.FeedList {
				for _, image := range tcFeedList.Images {
					imgUrl := fmt.Sprintf("https://photo.tuchong.com/%d/f/%d.jpg", image.UserId, image.ImgId)
					total++
					fd.WriteString(imgUrl + "\n")
					log.Println(imgUrl, total)
				}
			}
		}

		if tcResult.More == false {
			os.Exit(0)
		}

		nowPage++
	}

}
