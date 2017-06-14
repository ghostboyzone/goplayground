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

/**
 * 初始化wchat实例
 */
func NewWechat() *Wechat {
	uUid, err := getUuid()
	if err != nil {
		log.Fatal("get uuid failed: ", err)
	}

	we := &Wechat{
		AppId:       APP_ID,
		Uuid:        uUid,
		QrImgName:   "qrcode.jpg",
		RedirectUri: "",
		Auth:        AuthInfo{},
		IsLogin:     false,
	}

	authStr := we.readAuthInfoFromFile()
	if len(authStr) > 0 {
		we.SetAuthInfo(authStr)
	}
	return we
}

/**
 * 获取uuid
 */
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

/**
 * 显示登录二维码
 */
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

/**
 * 下载二维码图片
 */
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

/**
 * 扫码回调
 */
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
			// 获取授权信息
			we.GetAuthInfo()
			break
		} else if data["code"] == "408" {
			log.Fatal("Timeout, quit")
			break
		} else {
			log.Println(data)
		}
	}
}

/**
 * 扫码结果查询（轮询）
 */
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

/**
 * 获取授权信息
 */
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
		we.SetAuthInfo(bodyStr)
		we.WebWxInit()
		we.writeAuthInfoToFile(bodyStr)
	} else {
		log.Fatal("get info err, status:", resp.StatusCode)
	}
}

func (we *Wechat) writeAuthInfoToFile(authStr string) {
	fp, err := os.OpenFile(TMP_AUTH_FILE, os.O_CREATE|os.O_WRONLY, 0600)
	defer fp.Close()
	if err != nil {
		log.Println("Write auth info failed: create file failed")
		return
	}
	fp.WriteString(authStr)
	return
}
func (we *Wechat) readAuthInfoFromFile() string {
	fp, err := os.Open(TMP_AUTH_FILE)
	if err != nil {
		log.Println("Cannot find previous auth info, start to call login")
		return ""
	}
	body, _ := ioutil.ReadAll(fp)
	log.Println("Read auth info:", string(body))
	return string(body)
}

/**
 * 设置授权信息
 */
func (we *Wechat) SetAuthInfo(authStr string) {
	var auth AuthInfo
	xml.Unmarshal([]byte(authStr), &auth)
	we.Auth = auth
}

/**
 * 微信初始化
 */
func (we *Wechat) WebWxInit() error {
	we.Auth.DeviceId = DEVICE_ID

	urlStr := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit"
	v := url.Values{}
	v.Set("lang", "zh_CN")
	v.Set("pass_ticket", we.Auth.PassTicket)
	urlStr += "?" + v.Encode()

	sMsgBaseReq := SMessageBaseRequest{
		Uin:      we.Auth.Wxuin,
		Sid:      we.Auth.Wxsid,
		Skey:     we.Auth.Skey,
		DeviceID: we.Auth.DeviceId,
	}
	sMsgInit := SMessageInit{
		BaseRequest: sMsgBaseReq,
	}
	sMsgInitJsonBt, _ := json.Marshal(sMsgInit)
	sMsgInitJsonStr := string(sMsgInitJsonBt)

	client, _ := myHttp.NewClient(http.MethodPost, urlStr, bytes.NewBufferString(sMsgInitJsonStr))
	client.SetHeader("ContentType", "application/json; charset=UTF-8")
	resp, err := client.Do()
	if err != nil {
		log.Println("get info err:", err)
		return err
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		// log.Println(bodyStr)
		var wxInitMsg WxInitMessage
		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&wxInitMsg)
		// log.Println(wxInitMsg.User, wxInitMsg.SyncKey)
		// log.Println(wxInitMsg.BaseResponse["Ret"], we.Auth.UserName)
		if int(wxInitMsg.BaseResponse["Ret"].(float64)) != 0 {
			log.Println("wxInit failed, we need to login again")
			return errors.New("wxInit failed, we need to login again")
		}
		we.Auth.UserName = wxInitMsg.User["UserName"].(string)
		we.Auth.SyncKey = wxInitMsg.SyncKey
		we.IsLogin = true
	} else {
		log.Println("send msg err, status:", resp.StatusCode)
		return errors.New("wxInit failed, status" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}
	return nil
}

func (we *Wechat) isLogin() bool {
	if we.IsLogin {
		return true
	}
	if len(we.Auth.Wxsid) == 0 {
		return false
	}
	err := we.WebWxInit()
	if err != nil {
		we.Auth = AuthInfo{}
		return false
	}
	return true
}

/**
 * 微信会话保持
 */
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

/**
 * 发送消息
 */
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

func init() {
	uUidReg = regexp.MustCompile(`window\.QRLogin\.code( *)=( *)(\d+);( *)?window\.QRLogin\.uuid( *)=( *)"(.*)";`)
	loginCodeReg = regexp.MustCompile(`window.code=(\d+);`)
	loginRedirectUriReg = regexp.MustCompile(`window.redirect_uri="(\S+?)"`)
}
