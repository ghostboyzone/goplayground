package main

import (
	// "fmt"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// var a int

func walkFuc(path string, info os.FileInfo, err error) error {
	// time.Sleep(time.Second * 2)
	// a++
	// path = "E:\\code\\waimai\\x_commodity\\commodity\\index.php"
	// log.Println("a:", path, lastModifyMap[path])
	fileInfo, _ := os.Stat(path)

	if lastModifyMap[path] == 0 {
		lastModifyMap[path] = fileInfo.ModTime().Unix()

		// fileInfo.IsDir()

		newPath, err := getRemoteFilePath(path)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new", "local: "+path, "remote: "+newPath)
	} else {
		if fileInfo.ModTime().Unix() != lastModifyMap[path] {
			lastModifyMap[path] = fileInfo.ModTime().Unix()

			newPath, err := getRemoteFilePath(path)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("update", "local: "+path, "remote: "+newPath)
		}
	}

	// if a > 10 {
	// 	os.Exit(1)
	// }
	return nil
}

func getRemoteFilePath(path string) (string, error) {
	for local, remote := range rMap {
		if strings.Index(path, local) == 0 {
			return strings.Replace(path, local, remote, 1), nil
			// break
		}
	}
	return "", errors.New("cannot get matched remote path: " + path)
}

type hashMap map[string]string
type tHashMap map[string]int64

var (
	rMap          hashMap
	lastModifyMap tHashMap
)

func main() {
	rMap = make(hashMap)
	lastModifyMap = make(tHashMap)
	rMap["E:\\code\\waimai\\x_commodity\\commodity\\"] = "/home/map/test_20170329/"

	for {
		for tmpPath, tmpNewPath := range rMap {
			err := filepath.Walk(tmpPath, walkFuc)
			log.Println(err, tmpNewPath)
		}
		log.Println("Sleep for next round")
		time.Sleep(time.Second * 5)
	}

}
