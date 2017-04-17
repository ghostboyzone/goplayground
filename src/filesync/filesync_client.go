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
	// "compress/gzip"
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
	// max_ch    chan int
	parseConf myConf.ConfigConf
	dirMsgs   []Message
	fileMsgs  []Message
)

func main() {
	confFile := flag.String("conf", "config.json", "-conf config.json")
	flag.Parse()
	parseConf = myConf.NewConf(*confFile)
	rMap = parseConf.GetPaths()
	excludeFiles = parseConf.Client.ExcludeFiles
	lastModifyMap = make(tHashMap)

	dir_max_ch := make(chan int, 1)
	file_max_ch := make(chan int, parseConf.Client.SendChannels)

	for {
		startTime := time.Now().Unix()
		for tmpPath, tmpNewPath := range rMap {
			log.Println("Start Process: [local]", tmpPath, ", [remote]", tmpNewPath)

			dirMsgs = dirMsgs[0:0]
			fileMsgs = fileMsgs[0:0]

			_ = filepath.Walk(tmpPath, walkFuc)

			if len(dirMsgs) > 0 {
				dir_max_ch <- 1
				encodeDirMsgs, _ := json.Marshal(dirMsgs)
				go sendToRemote(encodeDirMsgs, dir_max_ch)
			}

			if len(fileMsgs) > 0 {
				arrFileMsgs := sliceMsg(fileMsgs, parseConf.Client.BatchSendFiles)

				for _, oneFileMsgs := range arrFileMsgs {
					encodeFileMsgs, _ := json.Marshal(oneFileMsgs)
					file_max_ch <- 1
					go sendToRemote(encodeFileMsgs, file_max_ch)
				}
			}

			log.Println("Done Process: [local]", tmpPath, ", [remote]", tmpNewPath)
		}
		endTime := time.Now().Unix()
		log.Println("Cost: ", (endTime - startTime), "s, Sleep for next round")
		time.Sleep(time.Second * 1)
	}
}

/**
 * divide a slice by size
 * @param  {[type]} msg  []Message     [description]
 * @param  {[type]} size int)          (msgs         [][]Message [description]
 * @return {[type]}      [description]
 */
func sliceMsg(msg []Message, size int) (msgs [][]Message) {
	var childMsgs []Message
	cnt := 0
	total := len(msg)
	for idx, one := range msg {
		childMsgs = append(childMsgs, one)
		cnt++
		if (cnt >= size || idx >= (total-1)) && (len(childMsgs) > 0) {
			msgs = append(msgs, childMsgs)
			childMsgs = childMsgs[0:0]
			cnt = 0
		}
	}
	return msgs
}

/**
 * send file to remote
 * @param  {[type]} encode_msg []byte        [description]
 * @param  {[type]} max_ch     chan          int           [description]
 * @return {[type]}            [description]
 */
func sendToRemote(encode_msg []byte, max_ch chan int) {
	u := url.URL{Scheme: "ws", Host: parseConf.Client.Addr, Path: "/echo"}
	log.Printf("start connecting to %s", u.String())
	var err error
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	log.Printf("successful connected to %s", u.String())

	err = conn.WriteMessage(websocket.BinaryMessage, encode_msg)
	if err != nil {
		log.Println("write error: ", err)
	}
	<-max_ch
}

/**
 * list a directory
 */
func walkFuc(path string, info os.FileInfo, err error) error {
	for _, ex_path := range excludeFiles {
		if strings.Contains(path, ex_path) {
			// log.Println("Skip: ", path, " , reason: ", ex_path)
			return nil
		}
	}

	// first time, send to remote
	if lastModifyMap[path] == 0 {
		lastModifyMap[path] = info.ModTime().Unix()
		newPath, err := getRemoteFilePath(path)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new", "local: "+path, "remote: "+newPath)

		prepareSend(path, newPath, info)
	} else {
		// if file or dir modified, send to remote
		if info.ModTime().Unix() != lastModifyMap[path] {
			lastModifyMap[path] = info.ModTime().Unix()
			newPath, err := getRemoteFilePath(path)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("update", "local: "+path, "remote: "+newPath)

			prepareSend(path, newPath, info)
		}
	}
	return nil
}

/**
 * prepare sending message
 * @param  {[type]} localPath  string        [description]
 * @param  {[type]} remotePath string        [description]
 * @param  {[type]} info       os.FileInfo)  (msg          Message [description]
 * @return {[type]}            [description]
 */
func prepareSend(localPath string, remotePath string, info os.FileInfo) (msg Message) {
	if info.IsDir() {
		msg = Message{
			MsgType:  2,
			Data:     "",
			Path:     remotePath,
			FileMode: info.Mode(),
		}

		dirMsgs = append(dirMsgs, msg)
	} else {
		ff_bytes, _ := ioutil.ReadFile(localPath)
		msg = Message{
			MsgType:  1,
			Data:     string(ff_bytes),
			Path:     remotePath,
			FileMode: info.Mode(),
		}

		fileMsgs = append(fileMsgs, msg)
	}
	return msg
}

/**
 * get remote file path
 * @param  {[type]} path string)       (string, error [description]
 * @return {[type]}      [description]
 */
func getRemoteFilePath(path string) (string, error) {
	for local, remote := range rMap {
		if strings.Index(path, local) == 0 {
			return strings.Replace(path, local, remote, 1), nil
		}
	}
	return "", errors.New("cannot get matched remote path: " + path)
}
