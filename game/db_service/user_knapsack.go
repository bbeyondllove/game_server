package db_service

import (
	"fmt"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/model"
	"reflect"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/go-xorm/xorm"
)

type UserKnapsack struct{}

const (
	knapsack_nid        = 1
	user_knapsack_table = "t_user_knapsack"
)

//获取自动添加序列号
func (this *UserKnapsack) NextSeqId(session *xorm.Session) (int, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.UserKnapsack{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return 0, err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return maxId + 1, nil
	} else {
		count, err := db.Mysql.Count(&model.UserKnapsack{})
		if err != nil {
			return 0, err
		}
		return int(count + knapsack_nid), nil
	}
}

//添加记录
func (this *UserKnapsack) Add(data_map *model.UserKnapsack) (bool, error) {
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

//删除记录
func (this *UserKnapsack) Delete(userId string, itemId int) (int64, error) {
	sql := "delete from " + user_knapsack_table + " where  user_id=" + userId + " and item_id=" + strconv.Itoa(itemId)
	//fmt.Printf(sql)
	r, err := db.Mysql.Exec(sql)
	count, _ := r.RowsAffected()
	if count == 0 || err != nil {
		logger.Errorf("[Error]: ", err)
		return 0, err
	}

	return count, err
}

//更新道具数据
func (this *UserKnapsack) UpdateData(userId string, itemId int, fields map[string]interface{}) (int64, error) {
	update := " set "
	i := 0
	update_field := make(map[string]interface{})

	for k, v := range fields {
		if v == nil || v == "" {
			continue
		}
		update_field[k] = v
	}

	for k, v := range update_field {
		update += k
		update += "="
		ntype := reflect.TypeOf(v).String()
		if ntype == "string" || ntype == "time.Time" {
			update += "'"
		}
		update += utils.Strval(v)
		if ntype == "string" || ntype == "time.Time" {
			update += "'"
		}

		if i < len(update_field)-1 {
			update += ","
		}
		i++
	}
	sql := "update " + user_knapsack_table + update + " where  user_id=" + userId + " and item_id=" + strconv.Itoa(itemId)
	fmt.Printf(sql)
	r, err := db.Mysql.Exec(sql)
	count, _ := r.RowsAffected()
	if count == 0 || err != nil {
		logger.Errorf("[Error]: %v", err)
		return 0, err
	}

	return count, err
}

//根据用户ID获取用户背包信息
//user_id 用户ID
func (this *UserKnapsack) GetDataByUid(userId string) ([]model.UserKnapsack, error) {
	data := []model.UserKnapsack{}
	err := db.Mysql.Table(user_knapsack_table).Where("user_id= ? ", userId).Find(&data)
	return data, err
}

//根据用户ID和道具ID获取用户背包信息
//user_id 用户ID
//itemId 道具ID
func (this *UserKnapsack) GetData(user_id string, itemId int) (model.UserKnapsack, error) {
	data := model.UserKnapsack{}
	_, err := db.Mysql.Table(user_knapsack_table).Where("user_id= ? and item_id=?", user_id, itemId).Get(&data)
	return data, err
}
