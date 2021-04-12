package db_service

import (
	"game_server/db"
	"game_server/game/model"

	"game_server/core/logger"
)

type TaskAward struct {
}

var (
	task_award_table = "t_task_award"
)

//添加记录
func (this *TaskAward) Add(data_map *model.TaskAward) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

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

//根据ID获取任务奖励
//id 任务ID
func (this *TaskAward) GetDataById(id int) (*model.TaskAward, error) {
	data := &model.TaskAward{}
	_, err := db.Mysql.Table(task_award_table).Where("task_id= ?", id).Get(&data)
	return data, err
}

//获取所有任务奖励
func (this *TaskAward) GetAllData() ([]*model.TaskAward, error) {
	var data []*model.TaskAward
	err := db.Mysql.Table(task_award_table).Where("1=1").Asc("task_id").Find(&data)
	return data, err
}

//更新奖励
func (this *TaskAward) UpdateDataById(id string, award_id string, award_num int) (int, error) {
	data := map[string]interface{}{
		"award_id":  award_id,
		"award_num": award_num,
	}
	number, err := db.Mysql.Table(task_award_table).Where("task_id= ? ", id).Update(data)
	return int(number), err
}
