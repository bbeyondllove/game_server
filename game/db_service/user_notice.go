package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"

	"game_server/core/logger"
)

type UserNotice struct{}

const (
	UserNoticeTable = "t_user_notice"
)

//添加记录
func (this *UserNotice) Add(data_map *model.UserNotice) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = session.Insert(data_map)
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

//根据用户昵称获取用户信息
//nick_name 用户昵称
func (this *UserNotice) Get(notice_id int, user_id string) (model.UserNotice, error) {
	data := model.UserNotice{}
	_, err := db.Mysql.Table(UserNoticeTable).Where("notice_id = ?", notice_id).And("user_id = ?", user_id).Get(&data)
	return data, err
}
