package db_service

import (
	"fmt"
	"game_server/db"
	"game_server/game/model"
	"strings"
	"time"

	"game_server/core/logger"
)

type Certification struct{}

const (
	CertificationTable = "t_certification"
)

//添加记录
func (this *Certification) Add(data_map *model.Certification) (bool, error) {
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

//根据用户ID获取用户信息
//user_id 用户ID
func (this *Certification) GetDataByUid(userId string) (model.Certification, error) {
	data := model.Certification{}
	_, err := db.Mysql.Table(CertificationTable).Where("user_id= ? ", userId).Get(&data)
	return data, err
}

// 获取用户实名认证信息
func (this *Certification) GetDataByUids(userId []string) ([]model.Certification, error) {
	var data []model.Certification
	err := db.Mysql.Table(CertificationTable).Where(fmt.Sprintf("user_id in  (%v) ", strings.Join(userId, ","))).Find(&data)
	return data, err
}

// 获取用户实名认证信息
func (this *Certification) GetDataByUidsAndStatus(userId []string, status string) ([]model.Certification, error) {
	var data []model.Certification
	err := db.Mysql.Table(CertificationTable).Where(fmt.Sprintf("user_id in  (%v) ", strings.Join(userId, ","))).And("status = ?", status).Find(&data)
	return data, err
}

// 更新认证信息状态
func (this *Certification) UpdateStatus(userId string, status int, suggestion string) (int, error) {
	data := map[string]interface{}{
		"status":      status,
		"reson":       suggestion,
		"update_time": time.Now(),
	}
	number, err := db.Mysql.Table(CertificationTable).Where("user_id= ? ", userId).Update(data)
	return int(number), err
}
