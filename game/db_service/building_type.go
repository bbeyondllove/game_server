package db_service

import (
	"game_server/db"
	"game_server/game/model"
)

type BuildType struct{}

const (
	build_type_table = "t_building_type"
)

func (this *BuildType) GetBuildingTypes() ([]model.BuildType, error) {
	data := []model.BuildType{}
	err := db.Mysql.Table(build_type_table).Find(&data)
	return data, err
}
