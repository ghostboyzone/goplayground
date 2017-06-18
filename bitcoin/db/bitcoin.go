package db

import (
	"github.com/tidwall/buntdb"
	// "log"
	"os"
	"time"
)

type BitCoin struct {
	Db         *buntdb.DB
	DbFileName string
	ReadOnly   bool
	Closed     bool
}

func InitBitCoin(dbFileName string, readOnly bool) (bitCoin *BitCoin, err error) {
	var db *buntdb.DB
	if readOnly {
		db, err = buntdb.Open(":memory:")
	} else {
		db, err = buntdb.Open(dbFileName)
	}

	if err != nil {
		return bitCoin, err
	}

	bitCoin = &BitCoin{
		Db:         db,
		DbFileName: dbFileName,
		ReadOnly:   readOnly,
		Closed:     false,
	}

	err = bitCoin.setConfig()
	if err != nil {
		return bitCoin, err
	}

	if readOnly {
		bitCoin.Load(dbFileName)
		go func(bitCoin *BitCoin, dbFileName string) {
			for {
				time.Sleep(1 * time.Second)
				if bitCoin.Closed {
					break
				}
				bitCoin.Load(dbFileName)
			}
		}(bitCoin, dbFileName)
	}

	return bitCoin, nil
}

func (bt *BitCoin) setConfig() (err error) {
	var config buntdb.Config
	err = bt.Db.ReadConfig(&config)
	if err != nil {
		return err
	}
	// config.AutoShrinkMinSize = 5 * 1024 * 1024
	config.AutoShrinkDisabled = true
	err = bt.Db.SetConfig(config)
	return err
}

func (bt *BitCoin) Close() {
	bt.Closed = true
	bt.Db.Close()
}

func (bt *BitCoin) Set(key string, value string) error {
	err := bt.Db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	return err
}

func (bt *BitCoin) Get(key string) (value string, err error) {
	err = bt.Db.View(func(tx *buntdb.Tx) error {
		value, err = tx.Get(key)
		if err != nil {
			return err
		}
		return nil
	})
	return value, err
}

func (bt *BitCoin) CreateIndex(indexName string, pattern string) error {
	return bt.Db.CreateIndex(indexName, pattern, buntdb.IndexString)
}

func (bt *BitCoin) Shrink() error {
	return bt.Db.Shrink()
}

func (bt *BitCoin) Load(path string) error {
	fp, err := os.Open(path)
	// defer fp.Close()
	if err != nil {
		return err
	}
	err = bt.Db.Load(fp)
	return err
}
