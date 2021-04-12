package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"

	"game_server/core/logger"
)

type CertificationRecord struct{}

const (
	CertificationRecordTable = "t_certification_record"
)

//添加记录
func (this *CertificationRecord) Add(data_map *model.CertificationRecord) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: %v", err)
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		logger.Errorf("[Error]: %v", err)
		return false, err
	}
	return true, nil
}

func (this *CertificationRecord) GetLastDataByUid(userId string) (model.CertificationRecord, error) {
	data := model.CertificationRecord{}
	_, err := db.Mysql.Table(CertificationTable).Where("user_id= ? ", userId).Desc("id").Limit(1, 0).Get(&data)
	return data, err
}

func (this *CertificationRecord) GetAllData(userId string) ([]model.CertificationRecord, error) {
	var data []model.CertificationRecord
	err := db.Mysql.Table(CertificationTable).Where("user_id= ? ", userId).Desc("create_time").Find(&data)
	return data, err
}
