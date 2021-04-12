package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"

	"game_server/core/logger"
)

//交易记录
type UserLevelConfig struct{}

const user_level_desc_table = "t_user_level_config"

//增加交易金额
func (gw *UserLevelConfig) Add(data_map *model.UserLevelConfig) (bool, error) {
	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err := db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}

	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}
	return true, nil

}

//获取所有用户等级道具信息
func (gw *UserLevelConfig) GetUserLevelList() ([]model.UserLevelConfig, error) {
	data := []model.UserLevelConfig{}
	err := db.Mysql.Table(user_level_desc_table).Find(&data)

	return data, err
}

func (gw *UserLevelConfig) GetLevelItemByLevel(level, newLevel int) ([]model.UserLevelConfig, error) {
	data := []model.UserLevelConfig{}
	err := db.Mysql.Table(user_level_desc_table).Where("item_num>0 AND level_id > ? AND level_id <= ? ", level, newLevel).Find(&data)

	return data, err
}
