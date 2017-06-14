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
		BaseUri:     "https://wx.qq.com/cgi-bin/mmwebwx-bin",
	}

	authStr := we.readAuthInfoFromFile()
	if len(authStr) > 0 {
		we.SetAuthInfo(authStr)
	}

	baseUri := we.readBaseUriFromFile()
	if len(baseUri) > 0 {
		we.BaseUri = baseUri
	}

	baseCookie := we.readBaseCookieFromFile()
	if len(baseCookie) > 0 {
		we.BaseCookie = baseCookie
	}

	go func() {
		time.Sleep(5 * time.Second)
		for {
			if we.isLogin() {
				we.WebWxNotify()
				we.WebWxSync()
			}
			time.Sleep(10 * time.Second)
		}
	}()

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

			cutIdx := strings.LastIndex(data["redirect_uri"], "/")
			if cutIdx == -1 {
				cutIdx = len(data["redirect_uri"])
			}
			we.BaseUri = data["redirect_uri"][:cutIdx]

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
	// v.Set("r", fmt.S)
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

	// log.Println(resp.Cookies())
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		log.Println("authStr:", bodyStr)
		we.SetAuthInfo(bodyStr)
		we.WebWxInit()
		we.writeAuthInfoToFile(bodyStr)
		we.writeBaseUriToFile(we.BaseUri)
		we.BaseCookie = resp.Cookies()
		we.writeBaseCookieToFile(we.BaseCookie)
	} else {
		log.Fatal("get info err, status:", resp.StatusCode)
	}
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

	urlStr := we.BaseUri + "/webwxinit"
	v := url.Values{}
	v.Set("r", strconv.FormatInt(time.Now().Unix(), 10))
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
	client.SetHeader("Content-Type", "application/json")
	client.SetCookies(we.BaseCookie)
	resp, err := client.Do()
	if err != nil {
		log.Println("wxInit err:", err)
		return err
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		var wxInitMsg WxInitMessage
		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&wxInitMsg)
		if int(wxInitMsg.BaseResponse["Ret"].(float64)) != 0 {
			log.Println("wxInit failed, we need to login again")
			return errors.New("wxInit failed, we need to login again")
		}
		we.Auth.UserName = wxInitMsg.User["UserName"].(string)
		we.Auth.SyncKey = wxInitMsg.SyncKey
		we.IsLogin = true
	} else {
		log.Println("wxInit err, status:", resp.StatusCode)
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
		os.Remove(TMP_AUTH_FILE)
		os.Remove(TMP_BASE_URI_FILE)
		return false
	}
	return true
}

/**
 * 微信会话保持
 */
func (we *Wechat) WebWxSync() {
	urlStr := we.BaseUri + "/webwxsync"
	v := url.Values{}
	v.Set("sid", we.Auth.Wxsid)
	v.Set("skey", we.Auth.Skey)
	v.Set("pass_ticket", we.Auth.PassTicket)
	urlStr += "?" + v.Encode()

	sMsgBaseReq := SMessageBaseRequest{
		Uin:      we.Auth.Wxuin,
		Sid:      we.Auth.Wxsid,
		Skey:     we.Auth.Skey,
		DeviceID: we.Auth.DeviceId,
	}

	sMsgSync := SMessageSync{
		BaseRequest: sMsgBaseReq,
		SyncKey:     we.Auth.SyncKey,
		rr:          time.Now().Unix(),
	}

	sMsgSyncJsonBt, _ := json.Marshal(sMsgSync)

	sMsgSyncJsonStr := string(sMsgSyncJsonBt)
	client, _ := myHttp.NewClient(http.MethodPost, urlStr, bytes.NewBufferString(sMsgSyncJsonStr))
	client.SetHeader("Content-Type", "application/json")
	client.SetCookies(we.BaseCookie)
	resp, err := client.Do()
	if err != nil {
		log.Println("sync msg err:", err)
		return
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		// log.Println(bodyStr)
		var wxSyncMsg WxSyncMessage
		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&wxSyncMsg)
		if int(wxSyncMsg.BaseResponse["Ret"].(float64)) != 0 {
			log.Println("sync msg failed")
			return
		}
		we.Auth.SyncKey = wxSyncMsg.SyncKey
		log.Println("sync success")
	} else {
		log.Println("sync msg err, status:", resp.StatusCode)
	}
	return
}

func (we *Wechat) WebWxNotify() {
	urlStr := we.BaseUri + "/webwxstatusnotify"
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
	sMsgNotify := SMessageNotify{
		BaseRequest:  sMsgBaseReq,
		Code:         3,
		FromUserName: we.Auth.UserName,
		ToUserName:   we.Auth.UserName,
		ClientMsgId:  time.Now().Unix(),
	}

	sMsgNotifyJsonBt, _ := json.Marshal(sMsgNotify)

	sMsgNotifyJsonStr := string(sMsgNotifyJsonBt)
	client, _ := myHttp.NewClient(http.MethodPost, urlStr, bytes.NewBufferString(sMsgNotifyJsonStr))
	client.SetHeader("Content-Type", "application/json")
	client.SetCookies(we.BaseCookie)
	resp, err := client.Do()
	if err != nil {
		log.Println("status notify err:", err)
		return
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		// log.Println(bodyStr)
		var wxStatusNotify WxStatusNotifyMessage
		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&wxStatusNotify)
		if int(wxStatusNotify.BaseResponse["Ret"].(float64)) != 0 {
			log.Println("status notify failed")
			return
		}
		log.Println("status notify success")
	} else {
		log.Println("status notify err, status:", resp.StatusCode)
	}
	return
}

/**
 *
 */
func (we *Wechat) GetContact() {
	urlStr := we.BaseUri + "/webwxgetcontact"
	v := url.Values{}
	v.Set("pass_ticket", we.Auth.PassTicket)
	v.Set("skey", we.Auth.Skey)
	v.Set("r", strconv.FormatInt(time.Now().Unix(), 10))
	v.Set("lang", "zh_CN")
	v.Set("seq", "653469746")
	urlStr += "?" + v.Encode()
	log.Println(urlStr)

	client, _ := myHttp.NewClient(http.MethodGet, urlStr, nil)
	client.SetHeader("Content-Type", "application/json")
	client.SetCookies(we.BaseCookie)
	resp, err := client.Do()
	if err != nil {
		log.Println("send msg err:", err)
		return
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		// log.Println(bodyStr)
		var wxContact WxContact
		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&wxContact)
		if int(wxContact.BaseResponse["Ret"].(float64)) != 0 {
			log.Println("get contact failed")
			return
		}
		we.MemberList = wxContact.MemberList
		log.Println("get contact success")
	} else {
		log.Fatal("get contact err, status:", resp.StatusCode)
	}
}

/**
 * 发送消息
 */
func (we *Wechat) SendMsg(msgStr string, toUserName string) {
	urlStr := we.BaseUri + "/webwxsendmsg"
	v := url.Values{}
	v.Set("pass_ticket", we.Auth.PassTicket)
	urlStr += "?" + v.Encode()

	sMsgBaseReq := SMessageBaseRequest{
		Uin:      we.Auth.Wxuin,
		Sid:      we.Auth.Wxsid,
		Skey:     we.Auth.Skey,
		DeviceID: we.Auth.DeviceId,
	}
	randId := fmt.Sprintf("%d", time.Now().UnixNano()/1e5)
	sMsgMsg := SMessageMsg{
		Type:         1,
		Content:      msgStr,
		FromUserName: we.Auth.UserName,
		ToUserName:   toUserName,
		LocalID:      randId,
		ClientMsgId:  randId,
	}
	sMsg := SMessage{
		BaseRequest: sMsgBaseReq,
		Msg:         sMsgMsg,
		Scene:       0,
	}
	sMsgJsonBt, _ := json.Marshal(sMsg)
	sMsgJsonStr := string(sMsgJsonBt)

	client, _ := myHttp.NewClient(http.MethodPost, urlStr, bytes.NewBufferString(sMsgJsonStr))
	client.SetHeader("Content-Type", "application/json")
	client.SetCookies(we.BaseCookie)
	resp, err := client.Do()
	if err != nil {
		log.Println("send msg err:", err)
		return
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		// log.Println(bodyStr)
		var wxSendMsgMsg WxSendMsgMessage
		decoder := json.NewDecoder(strings.NewReader(bodyStr))
		decoder.Decode(&wxSendMsgMsg)
		if int(wxSendMsgMsg.BaseResponse["Ret"].(float64)) != 0 {
			log.Println("send msg failed")
			return
		}
		log.Println("send msg success")
	} else {
		log.Fatal("send msg err, status:", resp.StatusCode)
	}
}

func init() {
	uUidReg = regexp.MustCompile(`window\.QRLogin\.code( *)=( *)(\d+);( *)?window\.QRLogin\.uuid( *)=( *)"(.*)";`)
	loginCodeReg = regexp.MustCompile(`window.code=(\d+);`)
	loginRedirectUriReg = regexp.MustCompile(`window.redirect_uri="(\S+?)"`)
}
