package statistical

import (
	"game_server/core/base"
	dao "game_server/game/db_service"
	"game_server/game/model"
	"strconv"
	"time"
)

// 统计模块
type Statistics struct {
}

// 实时获取统计数据
func (this *Statistics) RealOperate(date time.Time, platform int) model.StatisticsRealDay {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	record := model.StatisticsRealDay{}
	// 日期
	record.Date = date
	record.Platform = platform

	if platform == 0 {
		// 当日注册
		record.DayRegisteCount = DataCollectorIns.GetRegisteCount(day_format1)
		// 当日新增
		record.DayNewlyAdded = DataCollectorIns.GetNewlyAddedCount(day_format1)
		// 登录次数
		record.DayLoginCount = DataCollectorIns.GetLoginCount(day_format2)
		// 活跃用户数
		record.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
		// 当日产出CDT
		record.DayCdt = DataCollectorIns.CDTCount(day_format1)
	} else {
		// 当日注册
		record.DayRegisteCount = DataCollectorIns.GetPlatformRegisteCount(day_format1, platform)
		// 当日新增
		record.DayNewlyAdded = DataCollectorIns.GetPlatformNewlyAddedCount(day_format1, platform)
		// 登录次数
		record.DayLoginCount = DataCollectorIns.GetPlatformLoginCount(day_format2, platform)
		// 活跃用户数
		record.ActiveCount = DataCollectorIns.PlatformActiveCount(day_format2, platform)
		// 当日产出CDT
		record.DayCdt = DataCollectorIns.CDTCount(day_format1)
	}
	return record
}

// 定时器获取运营数据统计
func (this *Statistics) Operate(date time.Time, platform int) model.StatisticsDay {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	record := model.StatisticsDay{}
	// 日期
	record.Date = date
	record.Platform = platform
	if platform == 0 {
		// 总注册数
		record.TotalRegCount = DataCollectorIns.GetTotalRegisteCount()
		// 总新增
		record.TotalNewlyAdded = DataCollectorIns.GetTotalNewlyAddedCount()
		// 当日注册
		record.DayRegisteCount = DataCollectorIns.GetRegisteCount(day_format1)
		// 当日新增
		record.DayNewlyAdded = DataCollectorIns.GetNewlyAddedCount(day_format1)
		// 登录次数
		record.DayLoginCount = DataCollectorIns.GetLoginCount(day_format2)
		// 活跃用户数
		record.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
		// 在线时长
		record.OnlineTime = DataCollectorIns.OnlineTime(day_format2)
		// 人均在线时长
		record.AvgOnlineTime = 0
		if record.ActiveCount != 0 {
			record.AvgOnlineTime = record.OnlineTime / record.ActiveCount
		}
		// 手机累计注册数
		record.PhoneRegisteCount = DataCollectorIns.PhoneRegisteCount()
		// 邮箱累计注册数量
		record.EmailRegisteCount = DataCollectorIns.EmailRegisteCount()
		// 累计实名用户
		record.RealNameCount = DataCollectorIns.RealNameCount()
		// 当日产出CDT
		record.DayCdt = DataCollectorIns.CDTCount(day_format1)
		// 累计产出CDT
		record.TotalCdt = DataCollectorIns.TotalCDTCount()
	} else {
		// 总注册数
		record.TotalRegCount = DataCollectorIns.GetPlatformTotalRegisteCount(platform)
		// 总新增
		record.TotalNewlyAdded = DataCollectorIns.GetPlatformTotalNewlyAddedCount(platform)
		// 当日注册
		record.DayRegisteCount = DataCollectorIns.GetPlatformRegisteCount(day_format1, platform)
		// 当日新增
		record.DayNewlyAdded = DataCollectorIns.GetPlatformNewlyAddedCount(day_format1, platform)
		// 登录次数
		record.DayLoginCount = DataCollectorIns.GetPlatformLoginCount(day_format2, platform)
		// 活跃用户数
		record.ActiveCount = DataCollectorIns.PlatformActiveCount(day_format2, platform)
		// 在线时长
		record.OnlineTime = DataCollectorIns.PlatformOnlineTime(day_format2, platform)
		// 人均在线时长
		record.AvgOnlineTime = 0
		if record.ActiveCount != 0 {
			record.AvgOnlineTime = record.OnlineTime / record.ActiveCount
		}
		// 手机累计注册数
		record.PhoneRegisteCount = DataCollectorIns.PlatformPhoneRegisteCount(platform)
		// 邮箱累计注册数量
		record.EmailRegisteCount = DataCollectorIns.PlatformEmailRegisteCount(platform)
		// 累计实名用户
		record.RealNameCount = DataCollectorIns.PlatformRealNameCount(platform)
		// 当日产出CDT
		record.DayCdt = DataCollectorIns.CDTCount(day_format1)
		// 累计产出CDT
		record.TotalCdt = DataCollectorIns.TotalCDTCount()
	}
	return record
}

// 活跃折线图
func (this *Statistics) ActiveCount(date time.Time, hour int) []model.StatisticsActiveCount {
	// day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := make([]model.StatisticsActiveCount, 0)

	total_active_count := DataCollectorIns.ActiveCount(day_format2)
	for i := 0; i < hour; i++ {
		record := model.StatisticsActiveCount{}
		record.Date = date
		record.Hour = i
		record.ActiveCount = DataCollectorIns.HourActiveCount(day_format2, i)
		record.TotalActiveCount = total_active_count
		record.Ratio = 0
		if total_active_count != 0 {
			record.Ratio = int(record.ActiveCount * 100.0 / total_active_count)
		}
		result = append(result, record)
	}
	return result
}

// 留存数
func (this *Statistics) Retained(date time.Time) model.StatisticsRetained {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	record := model.StatisticsRetained{}

	record.Date = date
	record.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
	record.Date1 = date.AddDate(0, 0, -1)
	record.Retained1 = DataCollectorIns.GetRetainedCount(day_format1, 1)
	record.Added1 = DataCollectorIns.GetNewlyAddedCount(record.Date1.Format("2006-01-02"))
	record.Retained1Ratio = 0
	if record.Added1 > 0 {
		record.Retained1Ratio = int(record.Retained1 * 100.0 / record.Added1)
	}
	record.Date3 = date.AddDate(0, 0, -2)
	record.Retained3 = DataCollectorIns.GetRetainedCount(day_format1, 2)
	record.Added3 = DataCollectorIns.GetNewlyAddedCount(record.Date3.Format("2006-01-02"))
	record.Retained3Ratio = 0
	if record.Added3 > 0 {
		record.Retained3Ratio = int(record.Retained3 * 100.0 / record.Added3)
	}
	record.Date7 = date.AddDate(0, 0, -6)
	record.Retained7 = DataCollectorIns.GetRetainedCount(day_format1, 6)
	record.Added7 = DataCollectorIns.GetNewlyAddedCount(record.Date7.Format("2006-01-02"))
	record.Retained7Ratio = 0
	if record.Added7 > 0 {
		record.Retained7Ratio = int(record.Retained7 * 100.0 / record.Added7)
	}
	record.Date15 = date.AddDate(0, 0, -14)
	record.Retained15 = DataCollectorIns.GetRetainedCount(day_format1, 14)
	record.Added15 = DataCollectorIns.GetNewlyAddedCount(record.Date15.Format("2006-01-02"))
	record.Retained15Ratio = 0
	if record.Added15 > 0 {
		record.Retained15Ratio = int(record.Retained15 * 100.0 / record.Added15)
	}
	record.Date30 = date.AddDate(0, 0, -29)
	record.Retained30 = DataCollectorIns.GetRetainedCount(day_format1, 29)
	record.Added30 = DataCollectorIns.GetNewlyAddedCount(record.Date30.Format("2006-01-02"))
	record.Retained30Ratio = 0
	if record.Added30 > 0 {
		record.Retained30Ratio = int(record.Retained30 * 100.0 / record.Added30)
	}

	return record
}

// 宝箱统计
func (this *Statistics) TreasureBox(date time.Time) model.StatisticsTreasureBox {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	record := model.StatisticsTreasureBox{}

	record.Date = date
	record.OpenPv = DataCollectorIns.GetOpenBoxCount(day_format2)
	record.OpenUv = DataCollectorIns.GetOpenBoxPeopleCount(day_format2)
	record.DayOpenedCount = DataCollectorIns.GetFinishBoxCount(day_format2)
	record.TotalOpenedCount = record.DayOpenedCount + DataCollectorIns.LastGetTotalOpendCount(day_format2)
	record.DayCdtCount = DataCollectorIns.GetBoxCDTCount(day_format1)
	record.TotalCdtCount = DataCollectorIns.GetTotalBoxCDTCount()

	record.DayCdt = DataCollectorIns.GetBoxCDTValue(day_format1)
	record.TotalDayCdt = DataCollectorIns.GetTotalBoxCDTValue()

	return record
}

// 宝箱折线图
func (this *Statistics) RealTreasureBox(date time.Time, hour int) []model.StatisticsRealTreasureBox {
	// day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := make([]model.StatisticsRealTreasureBox, 0)

	total_active_count := DataCollectorIns.GetFinishBoxCount(day_format2)
	for i := 0; i < hour; i++ {
		record := model.StatisticsRealTreasureBox{}
		record.Date = date
		record.Hour = i
		record.OpenedCount = DataCollectorIns.GetHourFinishBoxCount(day_format2, i)
		record.DayOpenedCount = total_active_count
		record.Ratio = 0
		if total_active_count != 0 {
			record.Ratio = int(record.OpenedCount * 100.0 / record.DayOpenedCount)
		}
		result = append(result, record)
	}
	return result
}

// 签到
func (this *Statistics) SignIn(date time.Time) model.StatisticsSignIn {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := model.StatisticsSignIn{}
	result.Date = date
	result.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
	result.Pv = DataCollectorIns.GetSignInCount(day_format1)
	result.Uv = DataCollectorIns.GetSignInPeopleNum(day_format1)
	result.UsedMakeUpCardNum = DataCollectorIns.GetUsedMakeUpCardCount(day_format1)
	every_day_num, err := DataCollectorIns.GetSignInPeopleNumEveryDay(day_format1)
	if err == nil {
		result.S1 = every_day_num[0]
		result.S2 = every_day_num[1]
		result.S3 = every_day_num[2]
		result.S4 = every_day_num[3]
		result.S5 = every_day_num[4]
		result.S6 = every_day_num[5]
		result.S7 = every_day_num[6]
	}

	return result
}

// 新春CDT兑换
func (this *Statistics) StatisticsDoubleYearCdt(date time.Time) model.StatisticsDoubleYearCdt {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := model.StatisticsDoubleYearCdt{}
	result.Date = date
	result.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
	result.Pv = DataCollectorIns.GetDoubleYearPV(day_format1)
	result.Uv = DataCollectorIns.GetDoubleYearUV(day_format1)
	result.Retained1 = DataCollectorIns.GetRetainedCount(day_format1, 1)
	result.Retained3 = DataCollectorIns.GetRetainedCount(day_format1, 2)
	result.Retained7 = DataCollectorIns.GetRetainedCount(day_format1, 6)
	result.Cdt = DataCollectorIns.GetDoubleYearCdtValue(day_format1)

	return result
}

// 新春CDT每日用户CDT数
func (this *Statistics) StatisticsDoubleYearUserDayCdt(date time.Time) []model.StatisticsDoubleYearUserDayCdt {
	day_format1 := GetTimeStr(date) //2006-01-02
	// day_format2 := GetShortTimeStr(date) //20060102

	result := make([]model.StatisticsDoubleYearUserDayCdt, 0)
	data := DataCollectorIns.GetDoubleYearUserOneDayCdtValue(day_format1)
	if data != nil {
		for _, item := range data {
			cdt, err := strconv.ParseFloat(item["cdt"], 64)
			if err != nil {
				continue
			}
			user, _ := dao.UserIns.GetDataByUid(item["user_id"])
			// last_total_cdt, _ := dao.StatisticsDoubleYearUserDayCdtIns.GetTotalCdt(item["user_id"], nil, &date)
			last_total_cdt := DataCollectorIns.GetDoubleYearUserTotalCdtValue(day_format1, item["user_id"])
			result = append(result, model.StatisticsDoubleYearUserDayCdt{
				Date:     date,
				UserId:   item["user_id"],
				UserName: user.NickName,
				TotalCdt: last_total_cdt,
				Cdt:      float32(cdt),
			})
		}
	}
	return result
}

// 新春活动碎片兑换PVUV记录
func (this *Statistics) StatisticsDoubleYearFragment(date time.Time) model.StatisticsDoubleYearFragment {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := model.StatisticsDoubleYearFragment{}
	result.Date = date
	result.Pv = DataCollectorIns.GetDoubleYearFuwaPVValue(day_format2)
	result.Uv = DataCollectorIns.GetDoubleYearFuwaUVValue(day_format2)
	result.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
	result.Retained1 = DataCollectorIns.GetRetainedCount(day_format1, 1)
	result.Retained3 = DataCollectorIns.GetRetainedCount(day_format1, 2)
	result.Retained7 = DataCollectorIns.GetRetainedCount(day_format1, 6)
	return result
}

// 新春活动用户碎片兑换记录
func (this *Statistics) StatisticsDoubleYearUserFragment(date time.Time) []model.StatisticsDoubleYearUserFragment {
	// day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := make([]model.StatisticsDoubleYearUserFragment, 0)
	data := DataCollectorIns.GetDoubleYearFuwaUserDayValue(day_format2)
	if data != nil {
		for user_id, num_ := range data {
			num, err := strconv.Atoi(num_)
			if err != nil {
				continue
			}
			user, _ := dao.UserIns.GetDataByUid(user_id)
			item, err := dao.StatisticsDoubleYearUserFragmentIns.GetLastData(user_id)
			last_total_count := 0
			if err == nil && item.Id != 0 {
				last_total_count = item.TotalCount
			}
			result = append(result, model.StatisticsDoubleYearUserFragment{
				Date:       date,
				UserId:     user_id,
				UserName:   user.NickName,
				Count:      num,
				TotalCount: last_total_count + num,
			})
		}
	}
	return result
}

// 新春活动日排行
func (this *Statistics) StatisticsDoubleYearDailyRanking(date time.Time) model.StatisticsDoubleYearDailyRanking {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := model.StatisticsDoubleYearDailyRanking{}
	result.Date = date
	result.Pv = DataCollectorIns.GetDoubleYearDayRankingPVValue(day_format2)
	result.Uv = DataCollectorIns.GetDoubleYearDayRankingUVValue(day_format2)
	result.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
	result.Retained1 = DataCollectorIns.GetRetainedCount(day_format1, 1)
	result.Retained3 = DataCollectorIns.GetRetainedCount(day_format1, 2)
	result.Retained7 = DataCollectorIns.GetRetainedCount(day_format1, 6)
	// result.Cdt = DataCollectorIns.GetDoubleYearDayCdtValue(day_format1)
	result.Cdt = DataCollectorIns.GetDoubleYearDayRankingCdt(day_format2)

	return result
}

// 新春活动用户日排行
func (this *Statistics) StatisticsDoubleYearUserDailyRanking(date time.Time) []model.StatisticsDoubleYearUserDailyRanking {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	// 邀请用户数达到1个上日排行榜
	result := make([]model.StatisticsDoubleYearUserDailyRanking, 0)
	data, err := model.NewRankList().GetData(day_format2)
	if err != nil {
		return result
	}
	if data != nil {
		index := 1
		for _, item := range data {
			user_id := item.UserId
			user, _ := dao.UserIns.GetDataByUid(user_id)
			// 邀请数
			inviteNum, _ := dao.ActivityInvitationIns.GetUserDayInvitCount(user_id, day_format1+" 00:00:00", day_format1+" 23:59:59")
			if inviteNum < 1 {
				continue
			}
			result = append(result, model.StatisticsDoubleYearUserDailyRanking{
				Date:       date,
				UserId:     user_id,
				UserName:   user.NickName,
				Scores:     item.Scores,
				InvitedNum: int(inviteNum),
				Ranking:    index,
			})
			index++
			if index > 10 {
				break
			}
		}
	}
	return result
}

// 新春活动总排行
func (this *Statistics) StatisticsDoubleYearTotalRanking(date time.Time) model.StatisticsDoubleYearTotalRanking {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102

	result := model.StatisticsDoubleYearTotalRanking{}
	result.Date = date
	result.Pv = DataCollectorIns.GetDoubleYearTotalRankingPVValue(day_format2)
	result.Uv = DataCollectorIns.GetDoubleYearTotalRankingUVValue(day_format2)
	result.ActiveCount = DataCollectorIns.ActiveCount(day_format2)
	result.Retained1 = DataCollectorIns.GetRetainedCount(day_format1, 1)
	result.Retained3 = DataCollectorIns.GetRetainedCount(day_format1, 2)
	result.Retained7 = DataCollectorIns.GetRetainedCount(day_format1, 6)
	return result
}

// 新春活动用户总排行
func (this *Statistics) StatisticsDoubleYearUserTotalRanking(date time.Time) []model.StatisticsDoubleYearUserTotalRanking {
	day_format1 := GetTimeStr(date)      //2006-01-02
	day_format2 := GetShortTimeStr(date) //20060102
	result := make([]model.StatisticsDoubleYearUserTotalRanking, 0)

	// 邀请用户数达到10个上总榜
	data, err := model.NewRankList().GetTopRankData(day_format2)
	if err != nil {
		return result
	}
	if data != nil {
		index := 1
		for _, item := range data {
			user_id := item["user_id"].(string)
			user, _ := dao.UserIns.GetDataByUid(user_id)
			// 邀请数
			inviteNum, _ := dao.ActivityInvitationIns.GetUserOldInvitCount(user_id, day_format1+" 23:59:59", base.Setting.Springfestival.ActivityStartDatetime)
			if inviteNum < 10 {
				continue
			}
			result = append(result, model.StatisticsDoubleYearUserTotalRanking{
				Date:       date,
				UserId:     user_id,
				UserName:   user.NickName,
				Scores:     item["scores"].(int),
				InvitedNum: int(inviteNum),
				Ranking:    index,
			})
			index++
			if index > 10 {
				break
			}
		}
	}
	return result
}
