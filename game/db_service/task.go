package db_service

import (
	"game_server/db"
	"game_server/game/model"

	"game_server/core/logger"
)

type Tasks struct {
}

var (
	tasks_table = "t_tasks"
)

//添加记录
func (this *Tasks) Add(data_map *model.Tasks) (bool, error) {
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

//根据ID获取任务信息
//id 任务ID
func (this *Tasks) GetDataById(id int) (*model.Tasks, error) {
	data := &model.Tasks{}
	_, err := db.Mysql.Table(tasks_table).Where("id= ?", id).Get(&data)
	return data, err
}

//获取所有任务数据
func (this *Tasks) GetAllData() ([]*model.Tasks, error) {
	var data []*model.Tasks
	err := db.Mysql.Table(tasks_table).Where("1=1").Asc("id").Find(&data)
	return data, err
}

//获取所有任务数据
func (this *Tasks) GetPageDataByTaskType(task_type int, page, size int) (int, int, []*model.Tasks, error) {
	if page <= 0 {
		page = 1
	}
	count, err := db.Mysql.Table(tasks_table).Where("task_type = ?", task_type).Count()
	if err != nil {
		return 0, 0, nil, err
	}
	var data []*model.Tasks
	err = db.Mysql.Table(tasks_table).Where("task_type = ?", task_type).Asc("id").Find(&data)
	return int(count), page, data, err
}
