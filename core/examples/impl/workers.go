package impl

import (
	"game_server/core//worklist"
)

var WorldList = kk_core.NewWorkerList(0)
var mysqlworklist *kk_core.WorkerList

func InitWorkers() {
	mysqlworklist = kk_core.NewWorkerList(20)
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
