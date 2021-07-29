package models

import (
	"os"
	"path/filepath"
)

func init() {
	//获取路径
	ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Save = make(chan *JdCookie)
	go func() {
		for {
			<-Save
			for _, container := range Config.Containers {
				go container.Fresh()
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
