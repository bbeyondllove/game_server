package db_service

import (
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/model"
	"time"
)

type ActivityConfig struct{}

const (
	ActivityConfigTable = "t_activity_config"
)

//添加记录
func (this *ActivityConfig) Add(data_map *model.ActivityConfig) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *ActivityConfig) GetAllData() ([]model.ActivityConfig, error) {
	var data []model.ActivityConfig
	err := db.Mysql.Table(ActivityConfigTable).OrderBy("activity_type").
		Find(&data)
	return data, err
}

func (this *ActivityConfig) GetData(activity_type int) (model.ActivityConfig, error) {
	var data model.ActivityConfig
	_, err := db.Mysql.Table(ActivityConfigTable).Where("activity_type = ?", activity_type).Get(&data)
	return data, err
}

func (this *ActivityConfig) Update(config model.ActivityConfig) (int, error) {
	data := map[string]interface{}{
		"start_time":  utils.Time2Str(config.StartTime),
		"finish_time": utils.Time2Str(config.FinishTime),
	}
	number, err := db.Mysql.Table(ActivityConfigTable).Where("activity_type= ? ", config.ActivityType).Update(data)
	return int(number), err
}
