package statistical

import (
	"fmt"
	"time"

	"game_server/core/logger"

	"github.com/jasonlvhit/gocron"

	dao "game_server/game/db_service"
	"game_server/game/proto"
)

type StatisticsConfig struct {
	FuwaFragmentId int
}

var (
	StatisticsConfigIns  *StatisticsConfig
	StatisticsDotIns     *StatisticsDot
	UserDataCollectorIns *UserDataCollector
	DataCollectorIns     *DataCollector
	StatisticsIns        *Statistics
)

func init() {
	StatisticsConfigIns = &StatisticsConfig{}
	StatisticsDotIns = &StatisticsDot{}
	UserDataCollectorIns = &UserDataCollector{}
	DataCollectorIns = &DataCollector{}
	StatisticsIns = &Statistics{}
}

func StatisticsAllData() {
	logger.Debug("StatisticsAllData")
	now_time := time.Now()
	last_day := now_time.AddDate(0, 0, -1)
	// 运营数据统计
	// 所有平台数据
	{
		result := StatisticsIns.Operate(last_day, 0)
		_, err := dao.StatisticsDayIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		// web平台
		result = StatisticsIns.Operate(last_day, proto.PLATFORM_WEB)
		_, err = dao.StatisticsDayIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		// android平台
		result = StatisticsIns.Operate(last_day, proto.PLATFORM_ANDROID)
		_, err = dao.StatisticsDayIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		// IOS平台
		result = StatisticsIns.Operate(last_day, proto.PLATFORM_IOS)
		_, err = dao.StatisticsDayIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
	}
	// 留存记录
	{
		result := StatisticsIns.Retained(last_day)
		_, err := dao.StatisticsRetainedIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
	}
	// 活跃数
	{
		result := StatisticsIns.ActiveCount(last_day, 24)
		for _, item := range result {
			_, err := dao.StatisticsActiveCountIns.Add(&item)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	// 宝箱
	{
		result := StatisticsIns.TreasureBox(last_day)
		_, err := dao.StatisticsTreasureBoxIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		real_result := StatisticsIns.RealTreasureBox(last_day, 24)
		for _, item := range real_result {
			_, err := dao.StatisticsRealTreasureBoxIns.Add(&item)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	// 签到
	{
		result := StatisticsIns.SignIn(last_day)
		_, err := dao.StatisticsSignInIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
	}

	// 新春CDT
	{
		result := StatisticsIns.StatisticsDoubleYearCdt(last_day)
		_, err := dao.StatisticsDoubleYearCdtIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		user_day_result := StatisticsIns.StatisticsDoubleYearUserDayCdt(last_day)
		for _, item := range user_day_result {
			dao.StatisticsDoubleYearUserDayCdtIns.Add(&item)
		}
	}
	// 金童玉女兑换活动
	{
		result := StatisticsIns.StatisticsDoubleYearFragment(last_day)
		_, err := dao.StatisticsDoubleYearFragmentIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		user_result := StatisticsIns.StatisticsDoubleYearUserFragment(last_day)
		for _, item := range user_result {
			_, err = dao.StatisticsDoubleYearUserFragmentIns.Add(&item)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	// 新春每日排行榜
	{
		result := StatisticsIns.StatisticsDoubleYearDailyRanking(last_day)
		_, err := dao.StatisticsDoubleYearDailyRankingIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		user_result := StatisticsIns.StatisticsDoubleYearUserDailyRanking(last_day)
		for _, item := range user_result {
			_, err = dao.StatisticsDoubleYearUserDailyRankingIns.Add(&item)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	// 累计排行榜
	{
		result := StatisticsIns.StatisticsDoubleYearTotalRanking(last_day)
		_, err := dao.StatisticsDoubleYearTotalRankingIns.Add(&result)
		if err != nil {
			logger.Error(err)
		}
		user_result := StatisticsIns.StatisticsDoubleYearUserTotalRanking(last_day)
		for _, item := range user_result {
			_, err = dao.StatisticsDoubleYearUserTotalRankingIns.Add(&item)
			if err != nil {
				logger.Error(err)
			}
		}
		logger.Debugf("StatisticsAllData:", last_day, user_result)
		fmt.Println("StatisticsAllData:", last_day, user_result)
	}
}

// 初始化定时任务
func InitCronTask() {
	gocron.Every(1).Day().At("00:03:00").Do(StatisticsAllData) //统计前一天数据
}
