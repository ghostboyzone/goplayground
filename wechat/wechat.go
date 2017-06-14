package main

import (
	// "encoding/xml"
	myApi "github.com/ghostboyzone/goplayground/wechat/api"
	// "io/ioutil"
	// "log"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	wechat := myApi.NewWechat()
	wechat.ShowQrCode()
	wechat.WaitForScan()
	wechat.GetContact()

	nickName := "filehelper"
	toUserName := "filehelper"

	fp, _ := os.OpenFile("contact.mlog", os.O_CREATE|os.O_WRONLY, 0600)
	defer fp.Close()
	for _, v := range wechat.MemberList {
		toBeWrite := fmt.Sprintf("NickName[%s] RemarkName[%s] UserName[%s]\n\n", v["NickName"].(string), v["RemarkName"].(string), v["UserName"].(string))
		fp.WriteString(toBeWrite)
		if v["RemarkName"].(string) == nickName || v["NickName"].(string) == nickName {
			toUserName = v["UserName"].(string)
		}
	}

	for {
		msg := "my rand message, timestamp: " + strconv.FormatInt(time.Now().Unix(), 10)
		wechat.SendMsg(msg, toUserName)
		time.Sleep(10 * time.Second)
	}
}
