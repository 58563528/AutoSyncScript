package models

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/boltdb/bolt"
)

var db *bolt.DB
var JD_COOKIE = "JD_COOKIE"
var ENV = "env"
var TASK = "TASK"

func initDB() {
	var err error
	if Config.Database == "" {
		Config.Database = ExecPath + "/.jdc.db"
	}
	db, err = bolt.Open(Config.Database, 0600, nil)

	if err != nil {
		logs.Warn(err)
	}
	if err := CreateTable(JD_COOKIE); err != nil {
		logs.Warn(err)
	}
}

func CreateTable(table string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JD_COOKIE))
		if b == nil {
			_, err := tx.CreateBucket([]byte(JD_COOKIE))
			if err != nil {
				logs.Warn(err)
			}
		}
		return nil
	})
}
