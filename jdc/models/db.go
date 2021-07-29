package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/boltdb/bolt"
)

var db *bolt.DB
var JD_COOKIE = "JD_COOKIE"

func initDB() {
	var err error
	db, err = bolt.Open(ExecPath+"/.jdc.db", 0600, nil)
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

func GetJdCookies() []JdCookie {
	cks := []JdCookie{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JD_COOKIE))
		b.ForEach(func(k, v []byte) error {
			ck := JdCookie{}
			var _v = reflect.ValueOf(&ck).Elem()
			for _, vv := range strings.Split(string(v), ";") {
				v := strings.Split(vv, "=")
				if len(v) == 2 {
					t := _v.FieldByName(v[0])
					if t.CanSet() {
						switch t.Kind() {
						case reflect.Int:
							i, _ := strconv.Atoi(v[1])
							t.SetInt(int64(i))
						case reflect.String:
							t.SetString(v[1])
						}
					}
				}

			}
			cks = append(cks, ck)
			for i := range cks {
				cks[i].ID = i + 1
				if cks[i].Nickname == "" {
					cks[i].Nickname = "--"
				}
				if cks[i].ScanedAt == "" {
					cks[i].ScanedAt = "____-__-__"
				}
				if cks[i].BeanNum == "" {
					cks[i].BeanNum = "--"
				}
				if cks[i].Note == "" {
					cks[i].Note = "--"
				}
				if cks[i].Priority == 0 {
					cks[i].Priority = 1
				}
			}
			length := len(cks)
			for i := 0; i < length; i++ {
				max := i
				for j := i + 1; j < length; j++ {
					if cks[j].Priority > cks[max].Priority {
						max = j
					}
				}
				cks[i], cks[max] = cks[max], cks[i]
			}
			return nil
		})
		return nil
	})
	return cks
}

func SaveJdCookie(cks ...JdCookie) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JD_COOKIE))
		if b != nil {
			for _, ck := range cks {
				if ck.Available == "" {
					ck.Available = True
				}
				var data = ""
				var v = reflect.ValueOf(ck)
				var t = reflect.TypeOf(ck)
				for i := 0; i < v.NumField(); i++ {
					data += fmt.Sprintf("%s=%v;", t.Field(i).Name, v.Field(i).Interface())
				}
				err := b.Put([]byte(ck.PtPin), []byte(data))
				if err != nil {
					logs.Warn(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		logs.Warn(err)
		return err

	}
	return nil
}

func GetJdCookie(pin string) *JdCookie {
	var ck *JdCookie
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JD_COOKIE))
		if b != nil {
			v := b.Get([]byte(pin))
			if v == nil {
				return nil
			}
			ck = &JdCookie{}
			var _v = reflect.ValueOf(ck).Elem()
			for _, vv := range strings.Split(string(v), ";") {
				v := strings.Split(vv, "=")
				if len(v) == 2 {
					t := _v.FieldByName(v[0])
					if t.CanSet() {
						switch t.Kind() {
						case reflect.Int:
							i, _ := strconv.Atoi(v[1])
							t.SetInt(int64(i))
						case reflect.String:
							t.SetString(v[1])
						}
					}
				}

			}
		}
		return nil
	})
	if err != nil {
		logs.Warn(err)
	}
	return ck
}

func (ck *JdCookie) Updates(us ...interface{}) {
	ck = GetJdCookie(ck.PtPin)
	var _v = reflect.ValueOf(ck).Elem()
	if len(us) == 2 {
		t := _v.FieldByName(us[0].(string))
		if t.CanSet() {
			switch t.Kind() {
			case reflect.Int:
				if v, ok := us[1].(int); ok {
					t.SetInt(int64(v))
				}
			case reflect.String:
				if v, ok := us[1].(string); ok {
					t.SetString(v)
				}
			}
		}
	} else {
		switch us[0].(type) {
		case map[string]interface{}:
			for k, v := range us[0].(map[string]interface{}) {
				t := _v.FieldByName(k)
				if t.CanSet() {
					switch t.Kind() {
					case reflect.Int:
						if v, ok := v.(int); ok {

							t.SetInt(int64(v))
						}
					case reflect.String:
						if v, ok := v.(string); ok {
							t.SetString(v)
						}
					}
				}
			}
		case JdCookie:
			var a = reflect.ValueOf(us[0].(JdCookie))
			var t = reflect.TypeOf(us[0].(JdCookie))
			for i := 0; i < a.NumField(); i++ {
				name := t.Field(i).Name
				if name == ck.PtPin {
					continue
				}
				t := _v.FieldByName(name)
				if t.CanSet() {
					switch t.Kind() {
					case reflect.Int:
						if v := a.Field(i).Int(); v != 0 {
							t.SetInt(v)
						}
					case reflect.String:
						if v := a.Field(i).String(); v != "" {
							t.SetString(v)
						}
					}
				}
			}
		}
	}
	SaveJdCookie(*ck)
}
