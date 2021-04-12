package kk_core

import (
	"game_server/core/base"
)

// var mysqlDB *sql.DB
var WorldList = NewWorkerList(1)
var mysqlworklist *WorkerList

func InitWorkers() {

	mysqlworklist = NewWorkerList(base.Setting.Mysql.Goroutines)
}
func StopWorkers() {
	WorldList.Close()
}

func PushMysql(f func()) {
	mysqlworklist.Push(f)
}

func PushWorld(f func()) {
	WorldList.Push(f)
}
