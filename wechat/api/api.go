package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	myHttp "github.com/ghostboyzone/goplayground/wechat/http"
	"github.com/skratchdot/open-golang/open"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	uUidReg             *regexp.Regexp
	loginCodeReg        *regexp.Regexp
	loginRedirectUriReg *regexp.Regexp
)

const (
	APP_ID = "wx782c26e4c19acffb"
)

type Wechat struct {
	AppId       string
	Uuid        string
	QrImgName   string
	RedirectUri string
	// Client  myHttp
}

func NewWechat() *Wechat {
	uUid, err := getUuid()
	if err != nil {
		log.Fatal("get uuid failed: ", err)
	}

	return &Wechat{
		AppId:       APP_ID,
		Uuid:        uUid,
		QrImgName:   "qrcode.jpg",
		RedirectUri: "",
	}
}

func getUuid() (uUid string, err error) {
	v := url.Values{}
	v.Add("appid", APP_ID)
	v.Add("redirect_uri", "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	v.Add("fun", "new")
	v.Add("lang", "zh_CN")
	v.Add("_", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	urlStr := "https://login.wx.qq.com/jslogin" + "?" + v.Encode()

	client, _ := myHttp.NewClient(http.MethodGet, urlStr)
	resp, _ := client.Do()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		matches := uUidReg.FindStringSubmatch(bodyStr)
		if len(matches) == 0 {
			return uUid, errors.New("get uuid error, regexp failed")
		}
		vCode, _ := strconv.ParseInt(matches[3], 10, 64)
		uUid = matches[7]
		if vCode == 200 {
			return uUid, nil
		}
		return uUid, errors.New("get uuid error, window.QRLogin.code:" + string(vCode))
	}
	return uUid, errors.New("get uuid error, status code:" + string(resp.StatusCode))
}

func (we *Wechat) ShowQrCode() {
	if we.RedirectUri != "" {
		return
	}
	imgUrl := "https://login.weixin.qq.com/qrcode/" + we.Uuid
	err := downloadImg(imgUrl, we.QrImgName)
	if err != nil {
		log.Fatal("Download qr image failed")
	}

	open.Run(we.QrImgName)
	time.Sleep(2 * time.Second)
	defer os.Remove(we.QrImgName)
}

func downloadImg(imgUrl string, imgName string) error {
	httpResp, _ := http.Get(imgUrl)
	defer httpResp.Body.Close()
	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(imgName, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (we *Wechat) WaitForScan() {
	if we.RedirectUri != "" {
		return
	}
	for {
		data, loginErr := getLoginStatus(we.Uuid)
		if loginErr != nil {
			log.Fatal(loginErr)
		}

		if data["code"] == "201" {
			log.Println("Scanned, wait for confirm")
		} else if data["code"] == "200" {
			log.Println("Done", data["redirect_uri"])
			we.RedirectUri = data["redirect_uri"]
			break
		} else {
			log.Println(data)
		}
	}
}

type BaseRequest struct {
	XMLName    xml.Name `xml:"error" json:"-"`
	Ret        int      `xml:"ret" json:"-"`
	Message    string   `xml:"message" json:"-"`
	Skey       string   `xml:"skey" json:"Skey"`
	Wxsid      string   `xml:"wxsid" json:"Sid"`
	Wxuin      int64    `xml:"wxuin" json:"Uin"`
	PassTicket string   `xml:"pass_ticket" json:"-"`
	DeviceID   string   `xml:"-" json:"DeviceID"`
}

func (we *Wechat) GetInfo() {
	client, _ := myHttp.NewClient(http.MethodGet, we.RedirectUri)
	resp, err := client.Do()
	log.Println(err)
	if resp.StatusCode == 200 {
		// body, _ := ioutil.ReadAll(resp.Body)
		// bodyStr := string(body)
		// log.Println(resp.Header, resp.Re)

		var aa BaseRequest
		xml.NewDecoder(resp.Body.(io.Reader)).Decode(aa)
		log.Println(aa.Skey)
	}
}

func getLoginStatus(uUid string) (retData map[string]string, err error) {
	retData = make(map[string]string)
	v := url.Values{}
	v.Set("loginicon", "true")
	v.Set("uuid", uUid)
	v.Set("tip", "0")
	v.Set("_", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	reqUrl := "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login" + "?" + v.Encode()
	resp, _ := http.Get(reqUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	loginCodeData := loginCodeReg.FindStringSubmatch(bodyStr)
	if len(loginCodeData) == 0 {
		return retData, errors.New("get code failed")
	}

	retData["code"] = loginCodeData[1]
	switch loginCodeData[1] {
	case "201":
		return retData, nil
	case "200":
		loginRedirectUrlData := loginRedirectUriReg.FindStringSubmatch(bodyStr)
		if len(loginRedirectUrlData) == 0 {
			return retData, errors.New("get redirect uri failed")
		}
		retData["redirect_uri"] = loginRedirectUrlData[1]
		return retData, nil
	}
	return retData, nil
}

func init() {
	uUidReg = regexp.MustCompile(`window\.QRLogin\.code( *)=( *)(\d+);( *)?window\.QRLogin\.uuid( *)=( *)"(.*)";`)
	loginCodeReg = regexp.MustCompile(`window.code=(\d+);`)
	loginRedirectUriReg = regexp.MustCompile(`window.redirect_uri="(\S+?)"`)
}
