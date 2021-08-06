package models

import (
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
)

func init() {
	ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	logs.Info("当前%s", ExecPath)
	initVersion()
	initConfig()
	initUserAgent()
	initContainer()
	initDB()
	initHandle()
	initCron()
}
