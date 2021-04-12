package db_service

import (
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/model"
	"strconv"
	"time"
)

type UserTreasureBoxRecord struct {
}

const (
	USER_TREASURE_BOX_RECORD string = "t_user_treasure_box_record"
)

//添加记录
func (this *UserTreasureBoxRecord) Add(data_map *model.UserTreasureBoxRecord) (int64, error) {
	return db.Mysql.Table(USER_TREASURE_BOX_RECORD).Insert(data_map)
}

// 获取记录数量
func (this *UserTreasureBoxRecord) GetCount(user_id string) (int64, error) {
	return db.Mysql.Table(USER_TREASURE_BOX_RECORD).Where("user_id=?", user_id).Desc("id").Count()
}

// 获取所有的宝箱记录
func (this *UserTreasureBoxRecord) GetData(user_id string) ([]model.UserTreasureBoxRecord, error) {
	var data []model.UserTreasureBoxRecord
	err := db.Mysql.Table(USER_TREASURE_BOX_RECORD).Where("user_id=?", user_id).Desc("id").Find(&data)
	return data, err
}

// 获取宝箱分页记录
func (this *UserTreasureBoxRecord) GetPageData(user_id string, page int, size int) ([]model.UserTreasureBoxRecord, error) {
	page = page - 1
	if page < 0 {
		page = 0
	}
	var data []model.UserTreasureBoxRecord
	err := db.Mysql.Table(USER_TREASURE_BOX_RECORD).Where("user_id=?", user_id).Desc("id").Limit(size, page*size).Find(data)
	return data, err
}

//更新用户获取的
func (this *UserTreasureBoxRecord) UpdateBoxCtd(userId string, ctd float32, openTime string) (bool, error) {
	loc, _ := time.LoadLocation("Local") //获取时区
	formatTime, _ := time.ParseInLocation("2006-01-02 15:04:05", openTime, loc)
	data_map := model.UserTreasureBoxRecord{}
	data_map.WatchTime = 0
	data_map.UserId = userId
	data_map.Cdt = ctd
	data_map.OpenTime = formatTime
	_, err := this.Add(&data_map)
	if err != nil {
		logger.Errorf("UserTreasureBoxRecord add error: %v\n", err.Error())
		return false, err
	}

	user := &User{}
	userData, err := user.GetDataByUid(userId)
	if err != nil {
		logger.Errorf("GetDataByUid error: %v\n", err.Error())
		return false, err
	}
	finalCdt := userData.TreasureBoxTotalIncome + ctd
	//更新宝箱获得的 ctd
	//sql := "UPDATE " + ModelUserTable + " SET treasure_box_total_income=? WHERE user_id=? LIMIT 1"
	//_, err = db.Mysql.Exec(sql, finalCdt, userId)
	_, err = UpdateFields(UserTable, "user_id", userId, map[string]interface{}{"treasure_box_total_income": finalCdt})
	if err != nil {
		logger.Errorf("UserTreasureBoxRecord update error: %v\n", err.Error())
		return false, err
	}
	//更新redis
	this.updateUserInfoInRedis(userId, finalCdt)
	return true, nil
}

// updateUserInfoInRedis 更新redis中用户的cdt最新值.
func (this *UserTreasureBoxRecord) updateUserInfoInRedis(userId string, lastCdt float32) {
	redisClient := db.RedisMgr.GetRedisClient()
	_, err := redisClient.HSet(userId, "treasure_box_total_income", strconv.FormatFloat(float64(lastCdt), 'f', 4, 32)).Result()
	if err != nil {
		logger.Errorf("update lastCdt of user to redis failed! error: %v\n", err.Error())
	}
}
