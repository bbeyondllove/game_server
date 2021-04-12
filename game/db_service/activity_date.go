package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"
)

type ActivityDateInfo struct{}

const (
	ACTIVITY_DATE = "t_activity_date"
)

//添加记录
func (this *ActivityDateInfo) Add(data_map *model.ActivityDate) error {
	data_map.CreateTime = time.Now()
	dbSession := db.Mysql.Table(ACTIVITY_DATE)
	_, err := dbSession.Insert(data_map)
	if err != nil {
		return err
	}
	return nil
}
