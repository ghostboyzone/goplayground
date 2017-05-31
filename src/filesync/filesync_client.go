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
		isDirty := false
		for tmpPath, _ := range rMap {
			// log.Println("Start Process: [local]", tmpPath, ", [remote]", tmpNewPath)

			dirMsgs = dirMsgs[0:0]
			fileMsgs = fileMsgs[0:0]

			if _, statErr := os.Stat(tmpPath); os.IsNotExist(statErr) {
				log.Println("Path [", tmpPath, "] not exist, skipped!")
				continue
			}

			_ = filepath.Walk(tmpPath, walkFuc)

			if len(dirMsgs) > 0 {
				isDirty = true
				dir_max_ch <- 1
				encodeDirMsgs, _ := json.Marshal(dirMsgs)
				go sendToRemote(encodeDirMsgs, dir_max_ch)
			}

			if len(fileMsgs) > 0 {
				isDirty = true
				arrFileMsgs := sliceMsg(fileMsgs, parseConf.Client.BatchSendFiles)
				for _, oneFileMsgs := range arrFileMsgs {
					// log.Println("len", oneFileMsgs[0].Path)
					encodeFileMsgs, _ := json.Marshal(oneFileMsgs)
					file_max_ch <- 1
					go sendToRemote(encodeFileMsgs, file_max_ch)
				}
			}

			// if isDirty {
			// log.Println("Done Process: [local]", tmpPath, ", [remote]", tmpNewPath)
			// }
		}

		sleepTick := 0
		leftChannels := 0
		for {
			nowLeftChannels := len(dir_max_ch) + len(file_max_ch)
			if nowLeftChannels == 0 {
				break
			}
			if sleepTick%40 == 0 {
				if nowLeftChannels != leftChannels {
					log.Println(len(dir_max_ch)+len(file_max_ch), "channels left, please wait...")
					leftChannels = nowLeftChannels
				}
			}
			time.Sleep(time.Millisecond * 50)
			sleepTick++
		}

		endTime := time.Now().Unix()
		if isDirty {
			log.Println("Cost: ", (endTime - startTime), "s, Everything is now ready!")
		}
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
	var childMsgs, childMsgsCopy []Message
	childMsgsCopy = childMsgsCopy[0:0]
	cnt := 0
	total := len(msg)
	for idx, one := range msg {
		childMsgs = append(childMsgs, one)
		cnt++
		if (cnt >= size || idx >= (total-1)) && (len(childMsgs) > 0) {
			msgs = append(msgs, childMsgs)
			childMsgs = childMsgsCopy
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
	for {
		u := url.URL{Scheme: "ws", Host: parseConf.Client.Addr, Path: "/echo"}
		// log.Printf("start connecting to %s", u.String())
		var err error
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("dial: ", err)
			time.Sleep(20 * time.Millisecond)
			continue
		}
		defer conn.Close()
		// log.Printf("successful connected to %s", u.String())

		err = conn.WriteMessage(websocket.BinaryMessage, encode_msg)
		if err != nil {
			log.Println("write error: ", err)
			time.Sleep(20 * time.Millisecond)
			continue
		}
		break
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
		log.Println("new file => ", "local: "+path, "remote: "+newPath)

		prepareSend(path, newPath, info)
	} else {
		// if file or dir modified, send to remote
		if info.ModTime().Unix() != lastModifyMap[path] {
			lastModifyMap[path] = info.ModTime().Unix()
			newPath, err := getRemoteFilePath(path)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("update file", "local: "+path, "remote: "+newPath)

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
