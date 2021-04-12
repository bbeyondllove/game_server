package main

import (
	"kk_server/base"
	"kk_server/impl"
	"log"
	"net/http"
)

func main() {
	impl.InitWorkers()

	base.LogInit(true, "kk_server")

	socketmgr := new(impl.SocketMgr)
	socketmgr.Init()
	impl.Start()

}

func profile() {
	log.Fatal(http.ListenAndServe("0.0.0.0:5001", nil))
}
