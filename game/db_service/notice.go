package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"
)

type NoticeInfo struct{}

const (
	NOTICE_TABLE      = "t_notice"
	TIME_NOTICE_TABLE = "t_time_notice"
)

// 添加公告
func (this *NoticeInfo) AddLanterns(data_map *model.Notice, walkingLanterns *[]model.TimeNotice) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	for i := 0; i < len(*walkingLanterns); i++ {
		item := &(*walkingLanterns)[i]
		item.Id = 0
		item.NoticeId = data_map.Id
		item.CreateTime = time.Now()
		item.UpdateTime = time.Now()
		_, err = session.Insert(item)
		if err != nil {
			_ = session.Rollback()
			return false, err
		}
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *NoticeInfo) UpdateLanterns(data_map *model.Notice, walkingLanterns *[]model.TimeNotice) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.UpdateTime = time.Now()
	update_value := map[string]interface{}{
		"notice_type":    data_map.NoticeType,
		"notice_title":   data_map.NoticeTitle,
		"notice_content": data_map.NoticeContent,
		"notice_url":     data_map.NoticeUrl,
		"version":        data_map.Version,
		"remark":         data_map.Remark,
		"notice_time":    data_map.NoticeTime,
		"update_time":    data_map.UpdateTime,
	}

	_, err = session.Table(NOTICE_TABLE).Where("id = ? ", data_map.Id).Update(update_value)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	var tmp model.TimeNotice
	_, err = session.Table(TIME_NOTICE_TABLE).Where("notice_id = ?", data_map.Id).Unscoped().Delete(&tmp)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	for i := 0; i < len(*walkingLanterns); i++ {
		item := &(*walkingLanterns)[i]
		item.Id = 0
		item.NoticeId = data_map.Id
		item.CreateTime = time.Now()
		item.UpdateTime = time.Now()
		_, err = session.Table(TIME_NOTICE_TABLE).Insert(item)
		if err != nil {
			_ = session.Rollback()
			return false, err
		}
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取最新的升级公告
func (this *NoticeInfo) GetUpgradeNotice() (model.Notice, error) {
	session := db.Mysql.Table(NOTICE_TABLE)
	session.Where("notice_type = ? ", 0)
	session.OrderBy("id desc")
	notice := model.Notice{}
	_, err := session.Get(&notice)
	if err != nil {
		return notice, err
	}
	return notice, nil
}

func (this *NoticeInfo) GetLastNotice() ([]model.Notice, error) {
	var data []model.Notice

	sql := "SELECT A.* FROM " + NOTICE_TABLE + " A," +
		"(SELECT max(id) as id FROM " + NOTICE_TABLE + " group by notice_type) B " +
		"WHERE A.id = B.id " +
		"ORDER BY A.notice_type DESC "
	// sql := "select * from (select * from " + NOTICE_TABLE + " order by id desc) as a group by a.notice_type"
	err := db.Mysql.SQL(sql).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// 获取分页数据
func (this *NoticeInfo) GetPageData(notice_type int, start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.Notice, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(NOTICE_TABLE).Where("notice_type = ?", notice_type)
	if start_time != nil {
		session.And("create_time > ?", *start_time)
	}
	if end_time != nil {
		session.And("create_time <= ?", *end_time)
	}
	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}
	var data []model.Notice
	session = db.Mysql.Table(NOTICE_TABLE).Where("notice_type = ?", notice_type)
	if start_time != nil {
		session.And("create_time > ?", *start_time)
	}
	if end_time != nil {
		session.And("create_time <= ?", *end_time)
	}
	err = session.Desc("id").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

func (this *NoticeInfo) UpdateStatus(id int64, status int) (int, error) {
	data := map[string]interface{}{
		"is_noticed": status,
	}
	number, err := db.Mysql.Table(NOTICE_TABLE).Where("id= ? ", id).Update(data)
	return int(number), err
}

// 获取未通知完成的公告
func (this *NoticeInfo) GetUnnoticed() ([]model.Notice, error) {
	var data []model.Notice
	err := db.Mysql.Table(NOTICE_TABLE).Where("notice_time != '' AND is_noticed == 0").
		Find(&data)
	return data, err
}

type TimeNoticeInfo struct{}

// 获取跑马灯数据
func (this *TimeNoticeInfo) GetDataByNoticeId(id int64) ([]model.TimeNotice, error) {
	var data []model.TimeNotice
	err := db.Mysql.Table(TIME_NOTICE_TABLE).Where("notice_id = ?", id).
		Find(&data)
	return data, err
}

func (this *TimeNoticeInfo) GetDataByTime(now time.Time) ([]model.TimeNotice, error) {
	var data []model.TimeNotice
	err := db.Mysql.Table(TIME_NOTICE_TABLE).Where("stop_time < ?", now).
		Find(&data)
	return data, err
}
