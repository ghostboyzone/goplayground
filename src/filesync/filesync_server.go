package main

import (
	// "fmt"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Message struct {
	MsgType int    `json:"msg_type"`
	Data    string `json:"data"`
	Path    string `json:"path"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
} // use default options

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
		var decodeMsg Message
		json.Unmarshal(message, &decodeMsg)
		parseMsg(decodeMsg)

		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func parseMsg(message Message) {
	log.Println(message.Path)
	switch message.MsgType {
	case 1:
		path := strings.Replace(message.Path, "\\", "/", -1)
		log.Println(path)
		parentPath := filepath.Dir(path)

		err := os.MkdirAll(parentPath, 0777)
		if err != nil {
			log.Println(err)
		}

		err = ioutil.WriteFile(path, bytes.NewBufferString(message.Data).Bytes(), 0777)
		if err != nil {
			log.Println(err)
		}
	case 2:
		path := strings.Replace(message.Path, "\\", "/", -1)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	http.HandleFunc("/echo", echo)
	log.Println("listen at port 8989")
	log.Fatal(http.ListenAndServe(":8989", nil))
}
