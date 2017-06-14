package main

import (
	// "encoding/xml"
	myApi "github.com/ghostboyzone/goplayground/wechat/api"
	// "io/ioutil"
	"log"
	// "os"
	"strconv"
	"time"
)

func main() {
	wechat := myApi.NewWechat()
	wechat.ShowQrCode()
	wechat.WaitForScan()

	go func() {
		for {
			wechat.WebWxSync()
			time.Sleep(10 * time.Second)
		}
	}()

	for {
		msg := "my rand message, timestamp: " + strconv.FormatInt(time.Now().Unix(), 10)
		wechat.SendMsg(msg)
		log.Println("Send Message:", msg)
		time.Sleep(10 * time.Second)
	}
}
