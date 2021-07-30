package models

import (
	"math"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
)

func init() {
	//获取路径
	ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Save = make(chan *JdCookie)
	go func() {
		for {
			<-Save
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
				for i := range Config.Containers {
					(&Config.Containers[i]).read()
					if Config.Containers[i].Available {
						availables = append(availables, Config.Containers[i])
						weigth = append(weigth, float64(Config.Containers[i].Weigth))
						total += float64(Config.Containers[i].Weigth)
					}
				}
				if total == 0 {
					logs.Warn("容器都挂了")
					continue
				}

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
}

var True = "true"
var False = "false"

var Save chan *JdCookie

var ExecPath string
