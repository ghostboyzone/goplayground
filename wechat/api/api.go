package api

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	myHttp "github.com/ghostboyzone/goplayground/wechat/http"
	"github.com/skratchdot/open-golang/open"
	// "io"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	Auth AuthInfo
}

type AuthInfo struct {
	XMLName    xml.Name  `xml:"error" json:"-"`
	Ret        int       `xml:"ret" json:"-"`
	Message    string    `xml:"message" json:"-"`
	Skey       string    `xml:"skey" json:"Skey"`
	Wxsid      string    `xml:"wxsid" json:"Sid"`
	Wxuin      int64     `xml:"wxuin" json:"Uin"`
	PassTicket string    `xml:"pass_ticket" json:"-"`
	DeviceId   string    `xml:"-" json:"DeviceId"`
	UserName   string    `xml:"-" json:"UserName"`
	SyncKey    WxSyncKey `xml:"-" json:"SyncKey"`
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
		Auth:        AuthInfo{},
	}
}

func (we *Wechat) isLogin() bool {
	return len(we.Auth.PassTicket) != 0
}

func getUuid() (uUid string, err error) {
	v := url.Values{}
	v.Add("appid", APP_ID)
	v.Add("redirect_uri", "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	v.Add("fun", "new")
	v.Add("lang", "zh_CN")
	v.Add("_", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	urlStr := "https://login.wx.qq.com/jslogin" + "?" + v.Encode()

	client, _ := myHttp.NewClient(http.MethodGet, urlStr, nil)
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
	if we.isLogin() {
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
	if we.isLogin() {
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
			we.RedirectUri = data["redirect_uri"] + "&fun=new"
			break
		} else if data["code"] == "408" {
			log.Fatal("Timeout, quit")
			break
		} else {
			log.Println(data)
		}
	}
}

func (we *Wechat) GetAuthInfo() {
	if we.isLogin() {
		return
	}
	client, _ := myHttp.NewClient(http.MethodGet, we.RedirectUri, nil)
	resp, err := client.Do()
	if err != nil {
		log.Fatal("get info err:", err)
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		log.Println("authStr:", bodyStr)
		var auth AuthInfo
		xml.Unmarshal([]byte(bodyStr), &auth)
		we.Auth = auth
	} else {
		log.Fatal("get info err, status:", resp.StatusCode)
	}
}

func (we *Wechat) SetAuthInfo() {
	authStr := ``
	var auth AuthInfo
	xml.Unmarshal([]byte(authStr), &auth)
	we.Auth = auth
}

type SMessage struct {
	BaseRequest SMessageBaseRequest `json:"BaseRequest"`
	Msg         SMessageMsg         `json:"Msg"`
	Scene       int                 `json:"Scene"`
}
type SMessageInit struct {
	BaseRequest SMessageBaseRequest `json:"BaseRequest"`
}
type SMessageSync struct {
	BaseRequest SMessageBaseRequest `json:"BaseRequest"`
	SyncKey     WxSyncKey           `json:"SyncKey"`
	rr          int64               `json:"rr"`
}
type SMessageBaseRequest struct {
	Uin      int64  `json:"Uin"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}
type SMessageMsg struct {
	Type         int    `json:"Type"`
	Content      string `json:"Content"`
	FromUserName string `json:"FromUserName"`
	ToUserName   string `json:"ToUserName"`
	LocalID      string `json:"LocalID"`
	ClientMsgId  string `json:"ClientMsgId"`
}

type WxInitMessage struct {
	BaseResponse map[string]interface{} `json:"BaseResponse"`
	Count        int
	ContactList  [](map[string]interface{}) `json:"ContactList"`
	SyncKey      WxSyncKey                  `json:"SyncKey"`
	User         map[string]interface{}     `json:"User"`
}

type WxSyncKey struct {
	Count int                        `json:"Count"`
	List  [](map[string]interface{}) `json:"List"`
}

func (we *Wechat) WebWxInit() {
	we.Auth.DeviceId = "e679570618634105123"

	urlStr := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit"
	v := url.Values{}
	v.Set("lang", "zh_CN")
	v.Set("pass_ticket", we.Auth.PassTicket)
	urlStr += "?" + v.Encode()

	a := SMessageBaseRequest{
		Uin:      we.Auth.Wxuin,
		Sid:      we.Auth.Wxsid,
		Skey:     we.Auth.Skey,
		DeviceID: we.Auth.DeviceId,
	}

	b := SMessageInit{
		BaseRequest: a,
	}

	c, _ := json.Marshal(b)

	aaa := string(c)

	// log.Println(aaa)
	client, _ := myHttp.NewClient(http.MethodPost, urlStr, bytes.NewBufferString(aaa))
	client.SetHeader("ContentType", "application/json; charset=UTF-8")
	resp, err := client.Do()
	if err != nil {
		log.Fatal("get info err:", err)
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		// log.Println(bodyStr)

		var ttt WxInitMessage

		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&ttt)

		log.Println(ttt.User, ttt.SyncKey)

		eee, _ := json.Marshal(ttt.SyncKey)

		log.Println(string(eee))

		we.Auth.UserName = ttt.User["UserName"].(string)
		we.Auth.SyncKey = ttt.SyncKey

	} else {
		log.Fatal("send msg err, status:", resp.StatusCode)
	}
}

func (we *Wechat) WebWxSync() {
	urlStr := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync"
	v := url.Values{}
	v.Set("sid", we.Auth.Wxsid)
	v.Set("skey", we.Auth.Skey)
	v.Set("pass_ticket", we.Auth.PassTicket)
	urlStr += "?" + v.Encode()

	a := SMessageBaseRequest{
		Uin:      we.Auth.Wxuin,
		Sid:      we.Auth.Wxsid,
		Skey:     we.Auth.Skey,
		DeviceID: we.Auth.DeviceId,
	}

	b := SMessageSync{
		BaseRequest: a,
		SyncKey:     we.Auth.SyncKey,
		rr:          time.Now().Unix(),
	}

	c, _ := json.Marshal(b)

	aaa := string(c)
	client, _ := myHttp.NewClient(http.MethodPost, urlStr, bytes.NewBufferString(aaa))
	client.SetHeader("ContentType", "application/json; charset=UTF-8")
	resp, err := client.Do()
	if err != nil {
		log.Fatal("get info err:", err)
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		_ = string(body)
		// log.Println(bodyStr)
		log.Println("Sync success")
	} else {
		log.Fatal("sync msg err, status:", resp.StatusCode)
	}
}

func (we *Wechat) SendMsg(msgStr string) {
	a := SMessageBaseRequest{
		Uin:      we.Auth.Wxuin,
		Sid:      we.Auth.Wxsid,
		Skey:     we.Auth.Skey,
		DeviceID: we.Auth.DeviceId,
	}
	ttt := fmt.Sprintf("%d", time.Now().UnixNano()/1e5)
	b := SMessageMsg{
		Type:         1,
		Content:      msgStr,
		FromUserName: we.Auth.UserName,
		ToUserName:   "filehelper",
		LocalID:      ttt,
		ClientMsgId:  ttt,
	}
	c := SMessage{
		BaseRequest: a,
		Msg:         b,
		Scene:       0,
	}

	d, _ := json.Marshal(c)

	aaa := string(d)
	client, _ := myHttp.NewClient(http.MethodPost, "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg?pass_ticket="+we.Auth.PassTicket, bytes.NewBufferString(aaa))
	client.SetHeader("ContentType", "application/json; charset=UTF-8")
	resp, err := client.Do()
	if err != nil {
		log.Fatal("get info err:", err)
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		log.Println(bodyStr)
	} else {
		log.Fatal("send msg err, status:", resp.StatusCode)
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
