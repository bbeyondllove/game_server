package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/go-xorm/xorm"
)

type Awards struct {
}

var (
	award_nid  = 1
	AwardTable = "t_awards"
)

//获取自动添加序列号
func (this *Awards) NextSeqId(session *xorm.Session) (int, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.Awards{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return 0, err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return maxId + 1, nil
	} else {
		count, err := db.Mysql.Count(&model.Awards{})
		if err != nil {
			return 0, err
		}
		return int(count) + award_nid, nil
	}
}

//添加记录
func (this *Awards) Add(data_map *model.Awards) (bool, error, string) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	seqId, err := this.NextSeqId(session)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err, ""
	}

	data_map.Id = strconv.Itoa(seqId)
	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		_ = session.Rollback()
		return false, err, data_map.Id
	}
	err = session.Commit()
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err, data_map.Id
	}
	return true, nil, data_map.Id
}

//根据ID获取奖励信息
//id 奖励ID
func (this *Awards) GetDataById(id int) (model.Awards, error) {
	data := model.Awards{}
	_, err := db.Mysql.Table(AwardTable).Where("id= ?", id).Get(&data)
	return data, err
}

//获取所有奖励数据
func (this *Awards) GetAllData() ([]model.Awards, error) {
	var data []model.Awards
	err := db.Mysql.Table(AwardTable).Where("1=1").Asc("id").Find(&data)
	return data, err
}
