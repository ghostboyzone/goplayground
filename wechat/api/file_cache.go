package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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

func (we *Wechat) writeBaseUriToFile(baseUri string) {
	fp, err := os.OpenFile(TMP_BASE_URI_FILE, os.O_CREATE|os.O_WRONLY, 0600)
	defer fp.Close()
	if err != nil {
		log.Println("Write base uri failed: create file failed")
		return
	}
	fp.WriteString(baseUri)
	return
}
func (we *Wechat) readBaseUriFromFile() string {
	fp, err := os.Open(TMP_BASE_URI_FILE)
	if err != nil {
		log.Println("Cannot find previous base uri, start to call login")
		return ""
	}
	body, _ := ioutil.ReadAll(fp)
	log.Println("Read base uri:", string(body))
	return string(body)
}

func (we *Wechat) writeBaseCookieToFile(baseCookie []*http.Cookie) {
	fp, err := os.OpenFile(TMP_COOKIE_FILE, os.O_CREATE|os.O_WRONLY, 0600)
	defer fp.Close()
	if err != nil {
		log.Println("Write cookie failed: create file failed")
		return
	}
	tmpBt, _ := json.Marshal(baseCookie)
	tmpStr := string(tmpBt)
	fp.WriteString(tmpStr)
	return
}
func (we *Wechat) readBaseCookieFromFile() (baseCookie []*http.Cookie) {
	fp, err := os.Open(TMP_COOKIE_FILE)
	if err != nil {
		log.Println("Cannot find previous cookie, start to call login")
		return baseCookie
	}
	body, _ := ioutil.ReadAll(fp)
	bodyStr := string(body)
	log.Println("Read cookie:", bodyStr)
	decoder := json.NewDecoder(strings.NewReader(bodyStr))
	decoder.Decode(&baseCookie)
	return baseCookie
}
