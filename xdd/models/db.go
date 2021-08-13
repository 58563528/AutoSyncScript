package models

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/boltdb/bolt"
)

var db *bolt.DB
var JD_COOKIE = "JD_COOKIE"
var RECORD = "RECORD"
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
	if err, _ := CreateTable(RECORD); err != nil {
		panic(err)
	}
	if err, new := CreateTable(JD_COOKIE); err != nil {
		logs.Warn(err)
	} else if new {
		Record("=")
	}
	///
	if !Recorded("=") {
		type c struct {
			a []byte
			b []byte
		}
		var t = []c{}
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(JD_COOKIE))
			b.ForEach(func(k, v []byte) error {
				t = append(t, c{
					a: k,
					b: []byte(strings.Replace(string(v), "=", "^", -1)),
				})
				return nil
			})
			return nil
		})
		for _, v := range t {
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(JD_COOKIE))
				if b != nil {
					err := b.Put(v.a, v.b)
					if err != nil {
						logs.Warn(err)
					}
				}
				return nil
			})
		}
		Record("=")
	}
}

func CreateTable(table string) (error, bool) {
	new := false
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			_, err := tx.CreateBucket([]byte(table))
			if err != nil {
				logs.Warn(err)
			}
			new = true
		}
		return nil
	}), new

}

func Record(event string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RECORD))
		if b != nil {
			err := b.Put([]byte(event), []byte("ok"))
			if err != nil {
				logs.Warn(err)
			}
		}
		return nil
	})
}

func Recorded(event string) bool {
	is := false
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RECORD))
		if b != nil {
			v := b.Get([]byte(event))
			if v != nil {
				is = true
			}
		}
		return nil
	})
	return is
}
