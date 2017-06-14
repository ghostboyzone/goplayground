package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	apiReq "github.com/ghostboyzone/goplayground/bitcoin/api"
	"github.com/metakeule/fmtdate"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	/*
		nowT, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss ZZ", "2017-06-14 00:00:00 +00:00")
		log.Println(nowT.In(time.UTC).Unix())
		bitCoinTmp, _ := initBitCoin("jubi_coin_doge")
		bitCoinTmp.Db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bitCoinTmp.BucketName))
			b.ForEach(func(k, v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			return nil
		})

		resultTmp, _ := bitCoinTmp.GetKV([]byte(strconv.FormatInt(nowT.In(time.UTC).Unix(), 10)))
		log.Println(resultTmp)
		os.Exit(0)
	*/

	bitCoin, err := initBitCoin("jubi_coins")
	if err != nil {
		log.Fatal(err)
	}

	totalResult := apiReq.AllCoin()
	bitCoin.SwitchBucket("jubi_coins")
	for k, v := range totalResult {
		tmpJsonBt, _ := json.Marshal(v)
		bitCoin.WriteKV([]byte(k), tmpJsonBt)
	}

	for k, v := range totalResult {
		fmt.Printf("key=%s, value=%s\n", k, v)
		coinName := string(k)
		bitCoin.SwitchBucket(fmt.Sprintf("jubi_coin_%s", coinName))

		tmpData := apiReq.AdvanceKData(coinName, "5m")

		for _, tmpV := range tmpData {
			tmpK := int64(tmpV[0].(float64) / 1000)
			tmpJsonBt1, _ := json.Marshal(tmpV)
			log.Println(coinName, "start", tmpK, string(tmpJsonBt1))
			errW := bitCoin.WriteKV([]byte(strconv.FormatInt(tmpK, 10)), tmpJsonBt1)
			time.Sleep(time.Millisecond * 30)
			if errW != nil {
				log.Println(errW)
			}
			log.Println(coinName, "done", tmpK)
		}
	}

	/*
		bitCoin.Db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bitCoin.BucketName))
			b.ForEach(func(k, v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			return nil
		})
	*/

}

type BitCoin struct {
	Db         *bolt.DB
	DbFileName string
	BucketName string
}

func initBitCoin(bucketName string) (bitCoin *BitCoin, err error) {
	dbFileName := "bitcoin.db_data"
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return bitCoin, err
	}
	tx, err := db.Begin(true)
	if err != nil {
		return bitCoin, err
	}
	defer tx.Rollback()
	_, err = tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return bitCoin, err
	}
	err = tx.Commit()
	if err != nil {
		return bitCoin, err
	}

	return &BitCoin{
		Db:         db,
		DbFileName: dbFileName,
		BucketName: bucketName,
	}, nil
}

func (bt *BitCoin) SwitchBucket(bucketName string) error {
	bt.BucketName = bucketName
	err := bt.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bt.BucketName))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (bt *BitCoin) Close() {
	bt.Db.Close()
}

func (bt *BitCoin) WriteKV(key []byte, value []byte) error {
	err := bt.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bt.BucketName))
		return b.Put(key, value)
	})
	return err
}

func (bt *BitCoin) GetKV(key []byte) (retBt []byte, err error) {
	err = bt.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bt.BucketName))
		retBt = b.Get(key)
		return err
	})
	return retBt, err
}
