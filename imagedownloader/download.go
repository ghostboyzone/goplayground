package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	baseDir := "tmp"

	f, _ := os.Open("imglist.txt")
	scanner := bufio.NewScanner(f)

	cnt := 0
	for scanner.Scan() {
		cnt++
		log.Println("Cnt: ", cnt)
		lineStr := scanner.Text()
		list := strings.Split(lineStr, ",")
		if len(list) < 2 {
			log.Println("Err Format: ", lineStr)
			continue
		}
		catList := strings.Split(list[0], "|")
		imgUrl := list[1]
		// log.Println("Image Url: " + imgUrl)
		httpResp, httpErr := http.Get(imgUrl)
		if httpErr != nil {
			log.Println(httpErr, "Url: "+imgUrl)
			continue
		}

		for _, cat := range catList {
			os.MkdirAll(baseDir+"/"+cat, 0755)
			// tmpFile, tmpFileErr := os.Create(baseDir + "/" + cat + "/" + getRandFileName())
			tmpFile, tmpFileErr := os.Create(baseDir + "/" + cat + "/" + getMd5FileName(imgUrl))
			if tmpFileErr != nil {
				log.Println(tmpFileErr)
			}
			_, cpErr := io.Copy(tmpFile, httpResp.Body)
			if cpErr != nil {
				log.Println(cpErr)
			}
			tmpFile.Close()
		}
		httpResp.Body.Close()
	}
}

func getRandStr(size int) string {
	alpha := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ0123456789"
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf)
}

func getMd5FileName(url string) string {
	data := []byte(url)
	return fmt.Sprintf("%x", md5.Sum(data)) + ".jpg"
}

func getRandFileName() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + "_" + getRandStr(8) + ".jpg"
}
