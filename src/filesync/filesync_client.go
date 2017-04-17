package main

import (
	// "fmt"
	// "bytes"
	// "encoding/binary"
	"encoding/json"
	"errors"
	myConf "filesync/config"
	"flag"
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
	MsgType  int         `json:"msg_type"`
	Data     string      `json:"data"`
	Path     string      `json:"path"`
	FileMode os.FileMode `json:"filemode"`
}

type tHashMap map[string]int64

var (
	rMap          myConf.HashMap
	excludeFiles  []string
	lastModifyMap tHashMap
	// conn          *websocket.Conn
	max_ch    chan int
	parseConf myConf.ConfigConf
)

func main() {
	confFile := flag.String("conf", "config.json", "-conf config.json")
	flag.Parse()
	parseConf = myConf.NewConf(*confFile)
	rMap = parseConf.GetPaths()
	excludeFiles = parseConf.Client.ExcludeFiles
	lastModifyMap = make(tHashMap)

	max_ch = make(chan int, parseConf.Client.SendChannels)

	for {
		startTime := time.Now().Unix()
		for tmpPath, tmpNewPath := range rMap {
			log.Println("Start Process: [local]", tmpPath, ", [remote]", tmpNewPath)
			_ = filepath.Walk(tmpPath, walkFuc)
			log.Println("Done Process: [local]", tmpPath, ", [remote]", tmpNewPath)
		}
		endTime := time.Now().Unix()
		log.Println("Cost: ", (endTime - startTime), "s, Sleep for next round")
		time.Sleep(time.Second * 1)
	}
}

func prepareSend(localPath string, remotePath string, info os.FileInfo) (msg Message) {
	if info.IsDir() {
		msg = Message{
			MsgType:  2,
			Data:     "",
			Path:     remotePath,
			FileMode: info.Mode(),
		}
	} else {
		ff_bytes, _ := ioutil.ReadFile(localPath)
		// log.Println(err)
		// buf := new(bytes.Buffer)
		// _ = binary.Write(buf, binary.BigEndian, ff_bytes)
		// log.Println(err)
		msg = Message{
			MsgType:  1,
			Data:     string(ff_bytes),
			Path:     remotePath,
			FileMode: info.Mode(),
		}
	}
	return msg
}

func sendToRemote(encode_msg []byte) {
	u := url.URL{Scheme: "ws", Host: parseConf.Client.Addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	var err error
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.BinaryMessage, encode_msg)
	if err != nil {
		log.Println("write error: ", err)
	}
	<-max_ch
}

func walkFuc(path string, info os.FileInfo, err error) error {
	for _, ex_path := range excludeFiles {
		if strings.Contains(path, ex_path) {
			// log.Println("Skip: ", path, " , reason: ", ex_path)
			return nil
		}
	}

	if lastModifyMap[path] == 0 {
		lastModifyMap[path] = info.ModTime().Unix()
		newPath, err := getRemoteFilePath(path)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new", "local: "+path, "remote: "+newPath)

		sendMsg := prepareSend(path, newPath, info)
		encode_msg, err := json.Marshal(sendMsg)

		max_ch <- 1
		go sendToRemote(encode_msg)
	} else {
		if info.ModTime().Unix() != lastModifyMap[path] {
			lastModifyMap[path] = info.ModTime().Unix()
			newPath, err := getRemoteFilePath(path)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("update", "local: "+path, "remote: "+newPath)

			sendMsg := prepareSend(path, newPath, info)
			encode_msg, err := json.Marshal(sendMsg)

			max_ch <- 1
			go sendToRemote(encode_msg)
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
