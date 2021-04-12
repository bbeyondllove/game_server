package crontab

import (
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/message/activity/double_year"

	"github.com/jasonlvhit/gocron"
)

//启动定时任务
func Start() {
	logger.Debug("crontab Start")
	gocron.Every(1).Day().At("00:00:01").Do(signIn, false)          //每天执行
	gocron.Every(1).Week().Sunday().At("23:59:30").Do(signIn, true) //周日执行 （每周一 00：00：00 会清理一周的数据）

	// 双旦活动.
	rankList := double_year.NewRankList()
	// 每日排行榜奖励发放, 每天执行一次.
	err := gocron.Every(1).Day().At("00:00:05").Do(rankList.GiveOutDayAward)
	if err != nil {
		logger.Errorf("crontab GiveOutDayAward fail:%v\n", err)
	}
	// 累计排行榜奖励发放, 在活动结束后执行一次.
	t := utils.Str2Time(base.Setting.Springfestival.RankingListEndDate + " 00:00:05")
	err = gocron.Every(1).Day().From(&t).Do(rankList.GiveOutTotalAward)
	if err != nil {
		logger.Errorf("crontab GiveOutTotalAward fail:%v\n", err)
	}

	gocron.Start()
}
