/**
 * 清理数据库冗余数据
 */
package main

import (
	"fmt"
	"github.com/ghostboyzone/goplayground/bitcoin/db"
	"log"
	"os"
)

const (
	DB_FILE_PATH = "bitcoin.dbdata"
)

var (
	bt *db.BitCoin
)

func main() {
	log.Println("Start free db")
	var err error
	fp, err := os.Open(DB_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := fp.Stat()
	if err != nil {
		log.Fatal(err)
	}
	beforeSize := fileInfo.Size()
	fp.Close()

	bt, err = db.InitBitCoin(DB_FILE_PATH, false)
	if err != nil {
		log.Fatal(err)
	}
	bt.Shrink()
	bt.Close()

	fp, err = os.Open(DB_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err = fp.Stat()
	if err != nil {
		log.Fatal(err)
	}
	afterSize := fileInfo.Size()

	log.Println(fmt.Sprintf("Done free db, beforeSize[%f MB], afterSize[%f MB]", float64(beforeSize)/1024/1024, float64(afterSize)/1024/1024))
}
