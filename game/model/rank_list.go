package model

import (
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"strconv"
	"time"
)

// tableName 排行榜对应的表名.
var tableName = "t_rank_list"

// RankList 用户排行榜t_rank_list表对应的结构体.
type RankList struct {
	Id         int       `xorm:"int(20) not null"  json:"id" desc:"ID"`
	UserId     string    `xorm:"int(20) not null" json:"user_id" `
	Scores     int       `xorm:"int(20) not null" json:"scores"`
	UpdateTime time.Time `xorm:"timestamp not null" json:"update_time"`
	Date       string    `xorm:"datetime not null" json:"date"`
}

func (r *RankList) Insert() bool {
	_, err := db.Mysql.Insert(r)
	if err != nil {
		logger.Errorf("cdtRecord insert a record failed: %v, data:%v", err, *r)
		return false
	}
	return true
}

// GetDayRankList 获取每日排行榜数据.
func (r *RankList) GetDayRankList(isAward bool) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	rankLists := make([]RankList, 0)
	err := db.Mysql.Cols("user_id", "scores, update_time").Where("date=?", getDate(isAward)).OrderBy("scores desc, update_time asc").Find(&rankLists)
	if err != nil {
		logger.Errorf("GetDayRankList fail:%v\n", err)
		return result, err
	}

	for _, v := range rankLists {
		p := make(map[string]interface{})
		p["userId"] = v.UserId
		p["scores"] = v.Scores
		result = append(result, p)
	}
	return result, nil
}

// GetTopRankList 获取累计排行榜数据. todo 增加活动时间查询
func (r *RankList) GetTopRankList() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	var rankList []map[string]string
	sql := "SELECT sum(scores) as scores,ANY_VALUE(user_id) as user_id, ANY_VALUE(update_time) as update_time FROM t_rank_list where update_time >= ? GROUP BY user_id ORDER BY scores desc, update_time asc"
	//sql := "SELECT sum(scores) as scores,user_id, update_time FROM t_rank_list GROUP BY user_id ORDER BY scores desc, update_time asc"
	rankList, err := db.Mysql.SQL(sql, base.Setting.Springfestival.RankingListStartDate+" 00:00:00").QueryString()
	if err != nil {
		logger.Errorf("GetTopRankList fail:%v\n", err)
		return result, err
	}

	for _, v := range rankList {
		p := make(map[string]interface{})
		p["userId"] = v["user_id"]
		p["scores"], _ = strconv.Atoi(v["scores"])
		result = append(result, p)
	}
	return result, nil
}

// UpdateRankList 更新每日排行榜数据.
func (r *RankList) UpdateDayRankList(userId, date string, scores int) {
	rankList := new(RankList)
	ok, id, oldScores := r.getUserScores(userId, date)
	if !ok {
		rankList.Date = date
		rankList.UserId = userId
		rankList.Scores = scores
		rankList.UpdateTime = time.Now()
		_, err := db.Mysql.Insert(rankList)
		if err != nil {
			logger.Errorf("UpdateDayRankList[insert rank list] fail:%v\n", err)
		}
		return
	}

	rankList.Scores = scores + oldScores
	rankList.Date = date
	rankList.UpdateTime = time.Now()
	_, err := db.Mysql.Where("id=?", id).Update(rankList)
	if err != nil {
		logger.Errorf("UpdateTopRankList fail:%v\n", err)
	}
}

// getUserScores 获取排行榜用户指定日期的分数.
func (r *RankList) getUserScores(userId string, date string) (ok bool, primaryKey, scores int) {
	rankList := new(RankList)
	ok, err := db.Mysql.Where("user_id=?", userId).And("date=?", date).Get(rankList)
	if err != nil {
		logger.Errorf("getUserScores fail:%v\n", err)
		return false, 0, 0
	}

	if !ok {
		return false, 0, 0
	}

	return true, rankList.Id, rankList.Scores
}

// 获取排行榜指定日期的数据.
func (r *RankList) GetData(date string) ([]RankList, error) {
	var rankLists []RankList
	err := db.Mysql.Where("date=?", date).Desc("scores").Find(&rankLists)
	if err != nil {
		logger.Errorf("getUserScores fail:%v\n", err)
		return nil, err
	}
	return rankLists, nil
}

// 获取累计排行榜数据. dayDate 格式为:20200101
func (r *RankList) GetTopRankData(dayDate string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	var rankList []map[string]string
	sql := "SELECT sum(scores) as scores,ANY_VALUE(user_id) as user_id, ANY_VALUE(update_time) as update_time FROM t_rank_list where update_time >= ? and date <= ? GROUP BY user_id ORDER BY scores desc, update_time asc"
	rankList, err := db.Mysql.SQL(sql, base.Setting.Springfestival.RankingListStartDate+" 00:00:00", dayDate).QueryString()
	if err != nil {
		logger.Errorf("GetTopRankList fail:%v\n", err)
		return result, err
	}

	for _, v := range rankList {
		p := make(map[string]interface{})
		p["user_id"] = v["user_id"]
		p["scores"], _ = strconv.Atoi(v["scores"])
		result = append(result, p)
	}
	return result, nil
}

// getDate 获取当前日期前一天字符串.
func getDate(isAward bool) string {
	if isAward {
		return time.Now().AddDate(0, 0, -1).Format("20060102")
	}
	return time.Now().Format("20060102")
}

func GetDate(isAward bool) string {
	return getDate(isAward)
}

func NewRankList() *RankList {
	return &RankList{}
}
