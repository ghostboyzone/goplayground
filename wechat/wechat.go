package main

import (
	myApi "github.com/ghostboyzone/goplayground/wechat/api"
)

func main() {
	wechat := myApi.NewWechat()
	wechat.ShowQrCode()
	wechat.WaitForScan()
	wechat.GetInfo()
}
