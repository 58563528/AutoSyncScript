package models

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/boltdb/bolt"
)

func initHandle() {
	//获取路径
	Save = make(chan *JdCookie)
	go func() {
		init := true
		for {
			get := <-Save
			if get.Pool == "s" {
				initCookie()
				continue
			}
			cks := GetJdCookies()
			if Config.Mode == Parallel {
				for i := range Config.Containers {
					(&Config.Containers[i]).read()
				}
				for i := range Config.Containers {
					(&Config.Containers[i]).write(cks)
				}
			} else {
				resident := []JdCookie{}
				if Config.Resident != "" {
					tmp := cks
					cks = []JdCookie{}
					for _, ck := range tmp {
						if strings.Contains(Config.Resident, ck.PtPin) {
							resident = append(resident, ck)
						} else {
							cks = append(cks, ck)
						}
					}
				}
				weigth := []float64{}
				conclude := []int{}
				total := 0.0
				availables := []Container{}
				parallels := []Container{}
				for i := range Config.Containers {
					(&Config.Containers[i]).read()
					if Config.Containers[i].Available {
						if Config.Containers[i].Mode == Parallel {
							parallels = append(parallels, Config.Containers[i])
						} else {
							availables = append(availables, Config.Containers[i])
							weigth = append(weigth, float64(Config.Containers[i].Weigth))
							total += float64(Config.Containers[i].Weigth)
						}
					}
				}
				l := len(cks)
				for _, v := range weigth {
					conclude = append(conclude, int(math.Ceil(v/total*float64(l))))
				}
				a := 0
				for i, j := range conclude {
					s := a + j
					if s > l {
						s = l
					}
					availables[i].write(append(resident, cks[a:s]...))
					a += j
					if a >= l-1 {
						break
					}
				}
				for i := range parallels {
					parallels[i].write(cks)
				}
			}
			if init {
				go func() {
					for {
						Save <- &JdCookie{
							Pool: "s",
						}
						time.Sleep(time.Minute * 30)
						// time.Sleep(time.Second * 1)
					}
				}()
				init = false
			}
		}
	}()
}

type JdCookie struct {
	ID        int
	Priority  int
	ScanedAt  string
	LoseAt    string
	CreateAt  string
	PtKey     string
	PtPin     string
	Note      string
	Available string `validate:"oneof=true false"`
	Nickname  string
	BeanNum   string
	Pool      string
	// Delete    string `validate:"oneof=true false"`
}

func Date() string {
	return time.Now().Local().Format("2006-01-02")
}

var ScanedAt = "ScanedAt"
var LoseAt = "LoseAt"
var CreateAt = "CreateAt"
var Note = "Note"
var Available = "Available"
var PtKey = "PtKey"
var PtPin = "PtPin"
var Priority = "Priority"
var Nickname = "Nickname"
var BeanNum = "BeanNum"
var Pool = "Pool"
var True = "true"
var False = "false"
var Save chan *JdCookie
var ExecPath string

func (ck *JdCookie) ToPool(key string) {
	ck = GetJdCookie(ck.PtPin)
	if key == ck.PtKey {
		return
	}
	if strings.Contains(ck.Pool, key) {
		return
	}
	if ck.Pool == "" {
		ck.Pool = ck.PtKey
	} else {
		ck.Pool += "," + ck.PtKey
	}
	ck.Updates(map[string]interface{}{
		Available: True,
		PtKey:     key,
		Pool:      ck.Pool,
		ScanedAt:  Date(),
	})
}

func (ck *JdCookie) shiftPool() string {
	ck = GetJdCookie(ck.PtPin)
	if ck.Pool == "" {
		return ""
	}
	pool := strings.Split(ck.Pool, ",")
	shift := ""
	if len(pool) != 0 {
		shift = pool[0]
		pool = pool[1:]
	}
	us := map[string]interface{}{}
	if shift == "" {
		us[Pool] = ""
		us[Available] = False
		us[PtKey] = ""
	} else {
		us[Pool] = strings.Join(pool, ",")
		us[Available] = True
		us[PtKey] = shift
	}
	ck.Updates(us)
	return shift
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
			for i := range cks {
				if cks[i].PtPin == "" {
					continue
				}
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
			return nil
		})
		return nil
	})
	return cks
}

func NewJdCookie(cks ...JdCookie) {
	for i := range cks {
		cks[i].CreateAt = Date()
		cks[i].ScanedAt = cks[i].CreateAt
	}
	saveJdCookie(cks...)
}

func saveJdCookie(cks ...JdCookie) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JD_COOKIE))
		if b != nil {
			for _, ck := range cks {
				if ck.PtPin == "" || ck.PtKey == "" {
					continue
				}
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
					case reflect.Int64:
						if v, ok := v.(int64); ok {
							t.SetInt(v)
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
	saveJdCookie(*ck)
}
