package crontab

import (
	"encoding/json"
	"fmt"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/message/statistical_data"
	"game_server/game/model"
	"strconv"
	"time"
)

/**
数据统计
*/

// 签到统计
func signIn(isSunday bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("signIn Recovered in f", r)
		}
	}()
	logger.Debug("signIn Start isSunday：", isSunday)
	//如果是周一不用统计周日的数据（周日23：59 已统计）
	if time.Monday == time.Now().Weekday() {
		logger.Debug("signIn is Monday")
		return
	}

	var datetime time.Time
	//判断日期是否是 星期日, 并且是周日的定时任务调用
	if time.Sunday == time.Now().Weekday() && isSunday == true {
		logger.Debug("signIn is Sunday : ", isSunday)
		// 周日的定时任务时间不用 减1天 23:59:30 执行
		datetime = time.Now()
	} else {
		datetime = time.Now().AddDate(0, 0, -1)
	}
	// 2006-01-02 15:04:05
	dayDate := datetime.Format("2006-01-02")
	var signInNum int // 签到人数
	//去重签到人数
	number, err := db_service.UserTaskIns.SignInData(dayDate+" 00:00:01", dayDate+" 23:59:59")
	if err == nil {
		signInNum = number
	}
	numMap, err := db_service.UserTaskIns.EveryDayReward(dayDate+" 00:00:01", dayDate+" 23:59:59")
	everyData := map[string]int{
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
		"5": 0,
		"6": 0,
		"7": 0,
	}
	for _, value := range numMap {
		key := value["task_id"]
		key = key[len(key)-1:] //  截取最后一位， 最后一位代表第几天
		str_num := value["number"]
		number, err := strconv.Atoi(str_num)
		if err == nil {
			everyData[key] = number
		}
	}
	var activity model.ActivityDate
	activity.InvolvNum = signInNum //  参与人数
	mjson, _ := json.Marshal(everyData)
	mString := string(mjson)
	activity.DataDayNum = mString // 每天领取奖励的数量
	activity.ActivityType = 1     // 活动类型
	activity.Date = datetime      // 日期

	// 获取 pv
	keyDate := datetime.Format("20060102")
	//keyDate := time.Now().Format("20060102")
	pvKey := statistical_data.CITY_ICON_PV + keyDate
	isExists := db.RedisMgr.KeyExist(pvKey)
	if isExists {
		result := db.RedisMgr.Get(pvKey)
		if err == nil {
			pvNum, err := strconv.Atoi(result)
			if err == nil {
				activity.ActivityPv = pvNum
			}
		}
	}

	// 获取 uv
	uvKey := statistical_data.CITY_ICON_UV + keyDate
	isExists = db.RedisMgr.KeyExist(uvKey)
	if isExists {
		result := db.RedisMgr.HLen(uvKey)
		if err == nil {
			activity.ActivityUv = result
		}
	}

	// 获取 宝箱数量
	boxKey := statistical_data.BAOX_ADVERTIS + keyDate
	isExists = db.RedisMgr.KeyExist(boxKey)
	if isExists {
		result := db.RedisMgr.HLen(boxKey)
		if err == nil {
			activity.TreasureBoxNum = result
		}
	}

	// 统计补签卡
	activity.MakeUpCardNum = MakeUpCard(datetime)
	//添加数据
	err = db_service.ActivityDataIns.Add(&activity)
	if err != nil {
		logger.Debug("signIn err:", err)
	}
	logger.Debug("signIn end")
}

// 统计补签卡
func MakeUpCard(datetime time.Time) int {
	logger.Debug("MakeUpCard Start")
	// 2006-01-02 15:04:05
	dayDate := datetime.Format("2006-01-02")
	//dayDate := time.Now().Format("2006-01-02")
	sum := 0
	id := 0
	for {
		result, err := db_service.UserTaskIns.MakeUpCardList(id, dayDate+" 00:00:01", dayDate+" 23:59:59", 1000)
		if err != nil {
			logger.Debug("MakeUpCard db err:", err)
			break
		}
		if len(result) == 0 {
			logger.Debug("MakeUpCard db end:", err)
			break
		}
		for _, value := range result {
			id = int(value["id"].(int64))
			var dat []map[string]interface{}
			uintArr := value["award_info"].([]uint8)
			mapstr := string(uintArr)
			fmt.Println()
			if err := json.Unmarshal([]byte(mapstr), &dat); err == nil {
				ItemId := dat[0]["itemId"].(float64)
				// 统计补签卡
				if int(ItemId) == base.Setting.Statistical.MakeUpCard {
					sum++
				}
			} else {
				logger.Debug("MakeUpCard json err:", err)
				continue
			}
		}
	}
	logger.Debug("MakeUpCard Start")
	return sum
}
