package db_service

import (
	"game_server/db"
	"game_server/game/model"
)

//宝箱CDT分配规则
type TreasureBoxCdtConfig struct{}

const treasure_box_cdt_config = "t_treasure_box_cdt_config"

func (this *TreasureBoxCdtConfig) Add(data *model.TreasureBoxCdtConfig) (int64, error) {
	return db.Mysql.Table(treasure_box_cdt_config).Insert(data)
}

func (this *TreasureBoxCdtConfig) GetAllData(config_type int) ([]model.TreasureBoxCdtConfig, error) {
	var data []model.TreasureBoxCdtConfig
	err := db.Mysql.Table(treasure_box_cdt_config).Asc("id").Where("type = ?", config_type).Find(&data)
	return data, err
}
