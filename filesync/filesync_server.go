package main

import (
	// "fmt"
	// "bytes"
	"encoding/json"
	"flag"
	myConf "github.com/ghostboyzone/goplayground/filesync/config"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Message struct {
	MsgType  int         `json:"msg_type"`
	Data     []byte      `json:"data"`
	Path     string      `json:"path"`
	FileMode os.FileMode `json:"filemode"`
}

var (
	upgrader  websocket.Upgrader
	parseConf myConf.ConfigConf
)

func main() {
	confFile := flag.String("conf", "config.json", "-conf config.json")
	flag.Parse()
	parseConf = myConf.NewConf(*confFile)
	// set conf
	upgrader = websocket.Upgrader{
		ReadBufferSize:  parseConf.Server.ReadBufferSize,
		WriteBufferSize: parseConf.Server.WriteBufferSize,
	}

	http.HandleFunc("/echo", echo)
	log.Println("listen at port => ", parseConf.Server.Addr)
	log.Fatal(http.ListenAndServe(parseConf.Server.Addr, nil))
}

// websocket echo message
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", message)
		var decodeMsg []Message
		json.Unmarshal(message, &decodeMsg)

		for _, oneMsg := range decodeMsg {
			parseMsg(oneMsg)
		}
		// parseMsg(decodeMsg)

		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

// parse different type msg
func parseMsg(message Message) {
	log.Println(message.Path)
	switch message.MsgType {
	case 1:
		path := strings.Replace(message.Path, "\\", "/", -1)
		log.Println(path)
		parentPath := filepath.Dir(path)

		err := os.MkdirAll(parentPath, message.FileMode)
		if err != nil {
			log.Println(err)
		}

		// err = ioutil.WriteFile(path, bytes.NewBufferString(message.Data).Bytes(), message.FileMode)
		err = ioutil.WriteFile(path, message.Data, message.FileMode)
		if err != nil {
			log.Println(err)
		}
	case 2:
		path := strings.Replace(message.Path, "\\", "/", -1)
		err := os.MkdirAll(path, message.FileMode)
		if err != nil {
			log.Println(err)
		}
	}
}
