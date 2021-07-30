package models

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	//获取路径
	ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Save = make(chan *JdCookie)
	go func() {
		init := true
		for {
			get := <-Save
			if get.Pool == "s" {
				initCookie()
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
				// if total == 0 {
				// 	logs.Warn("容器都挂了")
				// 	continue
				// }

				l := len(cks)
				for _, v := range weigth {
					conclude = append(conclude, int(math.Ceil(v/total*float64(l))))
				}
				a := 0
				for i, j := range conclude {
					availables[i].write(cks[a : a+j])
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
	PtKey     string
	PtPin     string
	Note      string
	Available string `validate:"oneof=true false"`
	Nickname  string
	BeanNum   string
	Pool      string
}

var ScanedAt = "ScanedAt"
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

func (ck *JdCookie) toPool(key string) {
	ck = GetJdCookie(ck.PtPin)
	if key == ck.PtKey {
		return
	}
	if strings.Contains(ck.Pool, key) {
		return
	}
	if ck.Pool == "" {
		ck.Pool = key
	} else {
		ck.Pool += "," + key
	}
	ck.Updates(Pool, ck.Pool)
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
