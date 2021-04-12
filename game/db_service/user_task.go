package db_service

import (
	//"fmt"

	"game_server/db"
	"game_server/game/model"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/go-xorm/xorm"
)

type UserTask struct {
}

var (
	user_task_nid = 1
	UserTaskTable = "t_user_task"
)

//获取自动添加序列号
func (this *UserTask) NextSeqId(session *xorm.Session) (int, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.UserTask{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return 0, err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return maxId + 1, nil
	} else {
		count, err := db.Mysql.Count(&model.UserTask{})
		if err != nil {
			return 0, err
		}
		return int(count) + user_task_nid, nil
	}
}

//添加记录
func (this *UserTask) Add(data_map *model.UserTask) (bool, error, string) {
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

//根据用户ID获取任务信息
//user_id 用户ID
func (this *UserTask) GetDataByUid(userId string) ([]*model.UserTask, error) {
	var data []*model.UserTask
	err := db.Mysql.Table(UserTaskTable).Where("user_id= ? ", userId).Find(&data)
	return data, err
}

//获取今天前的所有指定类型任务
func (this *UserTask) GetData() ([]*model.UserTask, error) {
	var data []*model.UserTask
	err := db.Mysql.Table(UserTaskTable).Where("1=1").Find(&data)
	return data, err
}

//清除今天前的所有指定类型任务
func (this *UserTask) CleanData(taskType int) (int64, error) {
	sql := "delete from " + UserTaskTable + " where  create_time<curdate() and  task_type=" + strconv.Itoa(taskType)
	//fmt.Printf(sql)

	r, err := db.Mysql.Exec(sql)
	count, _ := r.RowsAffected()
	if count == 0 || err != nil {
		logger.Errorf("[Error]: ", err)
		return 0, err
	}

	return count, err
}

//获取签到人数 ()
func (this *UserTask) SignInData(startData string, endDate string) (int, error) {
	sql := "SELECT COUNT(DISTINCT user_id) as number from " + UserTaskTable + " where task_type = 3 and create_time >= ? and create_time <= ?"
	result, err := db.Mysql.SQL(sql, startData, endDate).QueryString()
	if err != nil {
		return 0, err
	}
	number := result[0]["number"]
	var num int
	num, err = strconv.Atoi(number)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// 每天领取奖励的数量
func (this *UserTask) EveryDayReward(startData string, endDate string) ([]map[string]string, error) {
	sql := "SELECT task_id,COUNT(*) AS number FROM " + UserTaskTable + " WHERE task_type=3  AND create_time >= ? AND create_time <= ? AND `status` = 2 GROUP BY task_id "
	result, err := db.Mysql.SQL(sql, startData, endDate).QueryString()
	return result, err
}

//统计活动周期内使用补签卡的数量
func (this *UserTask) MakeUpCardList(id int, startDate, endDate string, limit int) ([]map[string]interface{}, error) {
	model := db.Mysql.Table(UserTaskTable)
	model.Where("id > ?", id)
	model.Where("task_type = 3")
	model.In("status", 1, 2)
	model.Where("create_time >= ?", startDate)
	model.Where("create_time <= ?", endDate)
	model.Cols("id,award_info")
	model.Limit(limit)
	dataMap := make([]map[string]interface{}, 0)
	err := model.Find(&dataMap)
	if err != nil {
		return dataMap, err
	}
	return dataMap, nil
}

// 统计签到消耗的补签卡
func (this *UserTask) GetUsedCardCount(startDate, endDate string) (int, error) {
	model := db.Mysql.Table(UserTaskTable)
	model.Where("task_type = 3")
	model.And("signin_type = 2")
	model.Where("create_time >= ?", startDate)
	model.Where("create_time <= ?", endDate)
	num, err := model.Count()
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

// 统计签到数量
func (this *UserTask) GetSignInCount(startDate, endDate string) (int, error) {
	model := db.Mysql.Table(UserTaskTable)
	model.Where("task_type = 3")
	model.Where("create_time >= ?", startDate)
	model.Where("create_time <= ?", endDate)
	num, err := model.Count()
	if err != nil {
		return 0, err
	}
	return int(num), nil
}
