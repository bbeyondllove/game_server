package main

import (
	"fmt"
	kk_core "game_server/core"
	"game_server/core/base"
	"game_server/crontab"
	"game_server/db"
	"log"
	"time"

	//"game_server/game"
	"game_server/game/db_service"
	"game_server/game/message"
	"game_server/game/message/backstage"
	"game_server/game/message/statistical"

	//"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"game_server/core/logger"
)

func main() {
	db.InitRedis()
	db_service.DbInit()
	if !message.ReadCfg(message.G_BaseCfg) {
		return
	}

	message.G_ActivityManage.Init()
	message.InitDbData()
	kk_core.InitWorkers()
	fmt.Printf("讀取該行失敗\n")
	//game.Init_db()

	go func() {
		for {
			//处理用户信息，重置加载配置
			sig := make(chan os.Signal, 1)
			//signal.Notify(sig, syscall.SIGUSR1)

			select {
			case <-sig:
				message.ReloadConfig()
			}
		}
	}()

	loc, _ := time.LoadLocation("Local")
	logger.Debug("local location:", loc)

	statistical.StatisticsConfigIns.FuwaFragmentId = message.G_BaseCfg.Backstage.FuwaFragmentId
	// 定时任务
	statistical.InitCronTask()
	go crontab.Start()
	//监听端口
	go message.HttpServer()
	go backstage.HttpServer()
	if base.Setting.Server.Debug {
		go profile()
	}
	socketmgr := new(message.SocketMgr)
	socketmgr.Init()
	message.Start()

}

func profile() {
	log.Fatal(http.ListenAndServe("0.0.0.0:9911", nil))
}
