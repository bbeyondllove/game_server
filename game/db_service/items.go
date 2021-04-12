package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"

	"game_server/core/logger"
)

type Items struct {
}

var (
	item_table = "t_items"
)

//添加记录
func (this *Items) Add(data_map *model.Items) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		logger.Errorf("[Error]:", err)
		return false, err
	}
	return true, nil
}

//根据ID获取道具信息
//id 道具卡ID
func (this *Items) GetDataById(id int) (model.Items, error) {
	data := model.Items{}
	_, err := db.Mysql.Table(item_table).Where("id= ?", id).Get(&data)
	return data, err
}

//获取所有道具数据
func (this *Items) GetAllData() ([]model.Items, error) {
	var data []model.Items
	err := db.Mysql.Table(item_table).Where("1=1").Asc("id").Find(&data)
	return data, err
}
