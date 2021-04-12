package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"
)

type ActivityRedEnvelope struct {
}

const (
	ACTIVITY_RED_ENVELOPE string = "t_activity_red_envelope"
)

//添加记录
func (this *ActivityRedEnvelope) Add(data_map *model.ActivityRedEnvelope) (int64, error) {
	return db.Mysql.Table(ACTIVITY_RED_ENVELOPE).Insert(data_map)
}

//查询用户当天发送的红包
func (this *ActivityRedEnvelope) GetSameDayRedEnvelope(userId string) (bool, error, model.ActivityRedEnvelope) {
	var data model.ActivityRedEnvelope
	currentTime := time.Now().Format("2006-01-02")
	session := db.Mysql.Table(ACTIVITY_RED_ENVELOPE)
	session.Where("user_id = ? ", userId)
	session.Where("create_time >= ? ", currentTime+" 00:00:01")
	session.Where("create_time <= ? ", currentTime+" 23:59:59")
	has, err := session.Get(&data)
	return has, err, data
}

// 统计红包数量

func (this *ActivityRedEnvelope) GetCountDayRedEnvelope(userId string) (int64, error) {
	session := db.Mysql.Table(ACTIVITY_RED_ENVELOPE)
	session.Where("user_id = ? ", userId)
	return session.Count()
}
