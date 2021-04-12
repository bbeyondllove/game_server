package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"strconv"
	"time"

	"game_server/core/logger"
)

type WorldMap struct{}

const (
	world_nid   = 1
	world_table = "t_world_map_"
)

//获取自动添加序列号
func (this *WorldMap) NextSeqId(table_name string) (string, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(table_name).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return "", err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return strconv.Itoa(maxId + 1), nil
	} else {
		count, err := db.Mysql.Table(table_name).Count(&model.WorldMap{})
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(count + world_nid)), nil
	}
}

//添加记录
func (this *WorldMap) Add(code int, data_map *model.WorldMap) (bool, error) {
	table_name := world_table + strconv.Itoa(code)

	seqId, err := this.NextSeqId(table_name)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}

	data_map.Id = seqId
	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Table(table_name).Insert(data_map)

	return true, err
}

func (this *WorldMap) AddData(code int, dataMap *model.WorldMap) (int64, error) {
	table_name := world_table + strconv.Itoa(code)
	return db.Mysql.Table(table_name).Insert(dataMap)
}

//获取所有建筑
//code城市唯一代码
func (this *WorldMap) GetAllBuilding(code int) ([]model.WorldMap, error) {
	data := []model.WorldMap{}
	table_name := world_table + strconv.Itoa(code)
	err := db.Mysql.Table(table_name).Where("1=1").Find(&data)
	return data, err
}

//获取建筑简介
func (this *WorldMap) GetBuildDesc(code, x, y int32) (*model.WorldMap, error) {
	data := &model.WorldMap{}
	table_name := world_table + strconv.Itoa(int(code))
	_, err := db.Mysql.Table(table_name).Where("position_x=? and position_y=?", x, y).Get(data)
	return data, err
}

//获取所有空地
//code城市唯一代码
func (this *WorldMap) GetAllBlank(code int) ([]model.WorldMap, error) {
	data := []model.WorldMap{}
	err := db.Mysql.Table(world_table + strconv.Itoa(code)).Where("small_type='X0'").Find(&data)
	return data, err
}

//模糊查询商家
//key查询关键字
func (this *WorldMap) QueryShop(code int, key string) ([]model.WorldMap, error) {
	data := []model.WorldMap{}
	err := db.Mysql.Table(world_table + strconv.Itoa(code)).Where("shop_name like  '%" + key + "%'").Find(&data)
	return data, err
}
