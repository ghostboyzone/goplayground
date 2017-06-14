package api

import (
	"encoding/xml"
)

type Wechat struct {
	AppId       string
	Uuid        string
	QrImgName   string
	RedirectUri string
	// Client  myHttp
	Auth    AuthInfo
	IsLogin bool
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
