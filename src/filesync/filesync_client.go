package main

import (
	// "fmt"
	// "bytes"
	// "encoding/binary"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Message struct {
	MsgType int    `json:"msg_type"`
	Data    string `json:"data"`
	Path    string `json:"path"`
}

type hashMap map[string]string
type tHashMap map[string]int64

var (
	rMap          hashMap
	lastModifyMap tHashMap
	conn          *websocket.Conn
)

func main() {
	rMap = make(hashMap)
	lastModifyMap = make(tHashMap)
	// rMap["E:\\code\\waimai\\x_commodity\\commodity\\"] = "/home/map/test_20170329/"
	rMap["E:\\code\\waimai\\x_commodity\\commodity\\"] = "E:\\code\\waimai\\x_commodity\\commodity1\\"

	u := url.URL{Scheme: "ws", Host: "localhost:8989", Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	for {
		startTime := time.Now().Unix()
		for tmpPath, tmpNewPath := range rMap {
			log.Println("Start Process: [local]", tmpPath, ", [remote]", tmpNewPath)
			err := filepath.Walk(tmpPath, walkFuc)
			log.Println("Done Process: [local]", tmpPath, ", [remote]", tmpNewPath)
		}
		endTime := time.Now().Unix()
		log.Println("Cost: ", (endTime - startTime), "s, Sleep for next round")
		time.Sleep(time.Second * 1)
	}
}

func prepareSend(localPath string, remotePath string, isDir bool) (msg Message) {
	if isDir {
		msg = Message{
			MsgType: 2,
			Data:    "",
			Path:    remotePath,
		}
	} else {
		ff_bytes, _ := ioutil.ReadFile(localPath)
		// log.Println(err)
		// buf := new(bytes.Buffer)
		// _ = binary.Write(buf, binary.BigEndian, ff_bytes)
		// log.Println(err)
		msg = Message{
			MsgType: 1,
			Data:    string(ff_bytes),
			Path:    remotePath,
		}
	}
	return msg
}

func walkFuc(path string, info os.FileInfo, err error) error {
	// return nil
	if lastModifyMap[path] == 0 {
		lastModifyMap[path] = info.ModTime().Unix()
		newPath, err := getRemoteFilePath(path)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new", "local: "+path, "remote: "+newPath)

		sendMsg := prepareSend(path, newPath, info.IsDir())
		encode_msg, err := json.Marshal(sendMsg)
		err = conn.WriteMessage(websocket.BinaryMessage, encode_msg)
		if err != nil {
			log.Println("write error: ", err)
		}
	} else {
		if info.ModTime().Unix() != lastModifyMap[path] {
			lastModifyMap[path] = info.ModTime().Unix()
			newPath, err := getRemoteFilePath(path)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("update", "local: "+path, "remote: "+newPath)

			sendMsg := prepareSend(path, newPath, info.IsDir())
			encode_msg, err := json.Marshal(sendMsg)
			err = conn.WriteMessage(websocket.BinaryMessage, encode_msg)
			if err != nil {
				log.Println("write error: ", err)
			}
		}
	}
	return nil
}

func getRemoteFilePath(path string) (string, error) {
	for local, remote := range rMap {
		if strings.Index(path, local) == 0 {
			return strings.Replace(path, local, remote, 1), nil
		}
	}
	return "", errors.New("cannot get matched remote path: " + path)
}
