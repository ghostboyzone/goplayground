package db

import (
	"github.com/tidwall/buntdb"
)

type BitCoin struct {
	Db         *buntdb.DB
	DbFileName string
}

func InitBitCoin() (bitCoin *BitCoin, err error) {
	dbFileName := "bitcoin.dbdata"
	db, err := buntdb.Open(dbFileName)

	if err != nil {
		return bitCoin, err
	}

	bitCoin = &BitCoin{
		Db:         db,
		DbFileName: dbFileName,
	}
	return bitCoin, nil
}

func (bt *BitCoin) Close() {
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
