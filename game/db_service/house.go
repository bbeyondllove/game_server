package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/go-xorm/xorm"
)

type House struct{}

const (
	house_nid   = 1
	house_table = "t_house_"
)

//获取自动添加序列号
func (this *House) NextSeqId(session *xorm.Session) (string, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.House{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return "", err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return strconv.Itoa(maxId + 1), nil
	} else {
		count, err := db.Mysql.Count(&model.House{})
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(count + house_nid)), nil
	}
}

//添加记录
func (this *House) Add(data_map *model.House) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	seqId, err := this.NextSeqId(session)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}

	data_map.Id = seqId
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
		logger.Errorf("[Error]: ", err)
		return false, err
	}
	return true, nil
}

//获取指定位置的房子
//code城市唯一代码
//position_x 地图x坐标
//position_y 地图y坐标
//house_seq 房子序号
func (this *House) GetData(code int, position_x int, position_y int, house_seq string) ([]model.House, error) {
	data := []model.House{}
	err := db.Mysql.Table(house_table+strconv.Itoa(code)).Where("position_x = ? and position_y= ?  and house_seq=?", position_x, position_y, house_seq).Find(&data)
	return data, err
}
