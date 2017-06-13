package main

import (
	// "encoding/xml"
	myApi "github.com/ghostboyzone/goplayground/wechat/api"
	"log"
	"strconv"
	"time"
)

func main() {
	wechat := myApi.NewWechat()
	// wechat.SetAuthInfo()
	// log.Println(wechat.Auth.Wxsid)
	wechat.ShowQrCode()
	wechat.WaitForScan()
	wechat.GetAuthInfo()
	wechat.WebWxInit()

	go func() {
		for {
			wechat.WebWxSync()
			time.Sleep(10 * time.Second)
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		wechat.SendMsg("my rand message, timestamp: " + strconv.FormatInt(time.Now().Unix(), 10))
		log.Println("11")
	}

	// xml.Name{}

	// log.Println(wechat.Auth)
}
