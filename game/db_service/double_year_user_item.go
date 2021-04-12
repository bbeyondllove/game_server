package db_service

import (
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/model"
	"reflect"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/go-xorm/xorm"
)

type DoubleYearUserItem struct{}

const (
	DoubleYearUserItemtable_nid = 1
	DoubleYearUserItemtable     = "t_double_year_user_item"
)

//获取自动添加序列号
func (this *DoubleYearUserItem) NextSeqId(session *xorm.Session) (int, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.DoubleYearUserItem{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return 0, err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return maxId + 1, nil
	} else {
		count, err := db.Mysql.Count(&model.DoubleYearUserItem{})
		if err != nil {
			return 0, err
		}
		return int(count + DoubleYearUserItemtable_nid), nil
	}
}

//添加记录
func (this *DoubleYearUserItem) Add(data_map *model.DoubleYearUserItem) (bool, error) {
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
	data_map.IsTrade = 0

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
func (this *DoubleYearUserItem) Delete(userId string, itemId int) (int64, error) {
	sql := "delete from " + DoubleYearUserItemtable + " where  user_id=" + userId + " and item_id=" + strconv.Itoa(itemId)
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
func (this *DoubleYearUserItem) UpdateData(userId string, itemId int, fields map[string]interface{}) (int64, error) {
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
	sql := "update " + DoubleYearUserItemtable + update + " where  user_id=" + userId + " and item_id=" + strconv.Itoa(itemId)
	//fmt.Printf(sql)
	r, err := db.Mysql.Exec(sql)
	count, _ := r.RowsAffected()
	if count == 0 || err != nil {
		logger.Errorf("[Error]: %v", err)
		return 0, err
	}

	return count, err
}

//获取所有物品数据
func (this *DoubleYearUserItem) GetAllData() ([]*model.DoubleYearUserItem, error) {
	var data []*model.DoubleYearUserItem
	err := db.Mysql.Table(DoubleYearUserItemtable).Where("is_trade=0").Asc("id").Find(&data)
	return data, err
}
