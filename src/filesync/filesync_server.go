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
		// ret, err := cmdObj.Output()
		// log.Println(filepath.Dir(message.Path))

		parentPath := filepath.Dir(message.Path)

		err := os.MkdirAll(parentPath, 0777)
		if err != nil {
			log.Println(err)
		}

		err = ioutil.WriteFile(message.Path, bytes.NewBufferString(message.Data).Bytes(), 0777)
		if err != nil {
			log.Println(err)
		}
	case 2:
		err := os.MkdirAll(message.Path, 0777)
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
