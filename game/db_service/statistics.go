package db_service

import (
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/model"
	"time"
)

type StatisticsDay struct{}
type StatisticsRetained struct{}
type StatisticsActiveCount struct{}
type StatisticsTreasureBox struct{}
type StatisticsRealTreasureBox struct{}
type StatisticsSignIn struct{}
type StatisticsDoubleYearCdt struct{}
type StatisticsDoubleYearUserDayCdt struct{}
type StatisticsDoubleYearFragment struct{}
type StatisticsDoubleYearUserFragment struct{}
type StatisticsDoubleYearDailyRanking struct{}
type StatisticsDoubleYearUserDailyRanking struct{}
type StatisticsDoubleYearTotalRanking struct{}
type StatisticsDoubleYearUserTotalRanking struct{}

const (
	STATISTICSDAY_TABLE              = "t_statistics_day"
	STATISTICSRETAINED_TABLE         = "t_statistics_retained"
	STATISTICSACTIVECOUNT_TABLE      = "t_statistics_active_count"
	STATISTICS_TREASUREBOX_TABLE     = "t_statistics_treasure_box"
	STATISTICS_REALTREASUREBOX_TABLE = "t_statistics_real_treasure_box"

	STATISTICS_SIGNIN_TABLE                = "t_statistics_sign_in"
	STATISTICS_DOUBLEYEARCDT_TABLE         = "t_statistics_double_year_cdt"
	STATISTICS_DOUBLEYEARCDT_USERDAY_TABLE = "t_statistics_double_year_user_day_cdt"

	STATISTICS_DOUBLE_YEAR_FRAGMENT_TABLE           = "t_statistics_double_year_fragment"
	STATISTICS_DOUBLE_YEAR_USER_FRAGMENT_TABLE      = "t_statistics_double_year_user_fragment"
	STATISTICS_DOUBLE_YEAR_DAILY_RANKING_TABLE      = "t_statistics_double_year_daily_ranking"
	STATISTICS_DOUBLE_YEAR_USER_DAILY_RANKING_TABLE = "t_statistics_double_year_user_daily_ranking"
	STATISTICS_DOUBLE_YEAR_TOTAL_RANKING_TABLE      = "t_statistics_double_year_total_ranking"
	STATISTICS_DOUBLE_YEAR_USER_TOTAL_RANKING_TABLE = "t_statistics_double_year_user_total_ranking"
)

//日统计数据
func (this *StatisticsDay) Add(data_map *model.StatisticsDay) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取最后一条日总统计数据
func (this *StatisticsDay) GetLastData(platform int) (model.StatisticsDay, error) {
	data := model.StatisticsDay{}
	_, err := db.Mysql.Table(STATISTICSDAY_TABLE).Where("platform=?", platform).Desc("id").Limit(1).Get(&data)
	return data, err
}

// 获取日统计数据
func (this *StatisticsDay) GetData(start_time *time.Time, end_time *time.Time, platform *int, page, size int) (int, int, []model.StatisticsDay, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICSDAY_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	if platform != nil {
		session.And("platform = ?", *platform)
	}
	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICSDAY_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	if platform != nil {
		session.And("platform = ?", *platform)
	}

	var data []model.StatisticsDay
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 留存统计数据
func (this *StatisticsRetained) Add(data_map *model.StatisticsRetained) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取最后一条留存记录
func (this *StatisticsRetained) GetLastData() (model.StatisticsRetained, error) {
	data := model.StatisticsRetained{}
	_, err := db.Mysql.Table(STATISTICSRETAINED_TABLE).Desc("id").Limit(1).Get(&data)
	return data, err
}

// 获取留存数据
func (this *StatisticsRetained) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsRetained, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICSRETAINED_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICSRETAINED_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsRetained
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 添加宝箱记录
func (this *StatisticsActiveCount) Add(data_map *model.StatisticsActiveCount) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 宝箱统计数据,day格式为:'2017-06-16'
func (this *StatisticsActiveCount) GetData(day string) ([]model.StatisticsActiveCount, error) {
	var data []model.StatisticsActiveCount
	err := db.Mysql.Table(STATISTICSACTIVECOUNT_TABLE).Where("DATE_FORMAT(date,'%Y-%m-%d')  = ?", day).
		OrderBy("hour").
		Find(&data)
	return data, err
}

// 宝箱统计数据
func (this *StatisticsTreasureBox) Add(data_map *model.StatisticsTreasureBox) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取最后一条宝箱记录
func (this *StatisticsTreasureBox) GetLastData() (model.StatisticsTreasureBox, error) {
	data := model.StatisticsTreasureBox{}
	_, err := db.Mysql.Table(STATISTICS_TREASUREBOX_TABLE).Desc("id").Limit(1).Get(&data)
	return data, err
}

// 宝箱记录
func (this *StatisticsTreasureBox) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsTreasureBox, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_TREASUREBOX_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_TREASUREBOX_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsTreasureBox
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 宝箱记录
func (this *StatisticsRealTreasureBox) Add(data_map *model.StatisticsRealTreasureBox) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取宝箱记录,day格式为:'2017-06-16'
func (this *StatisticsRealTreasureBox) GetData(day string) ([]model.StatisticsRealTreasureBox, error) {
	var data []model.StatisticsRealTreasureBox
	err := db.Mysql.Table(STATISTICS_REALTREASUREBOX_TABLE).Where("DATE_FORMAT(date,'%Y-%m-%d')  = ?", day).
		OrderBy("hour").
		Find(&data)
	return data, err
}

//签到统计数据
func (this *StatisticsSignIn) Add(data_map *model.StatisticsSignIn) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取签到统计数据
func (this *StatisticsSignIn) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsSignIn, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_SIGNIN_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_SIGNIN_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsSignIn
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 新春CDT兑换数据
func (this *StatisticsDoubleYearCdt) Add(data_map *model.StatisticsDoubleYearCdt) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春CDT兑换数据
func (this *StatisticsDoubleYearCdt) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearCdt, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearCdt
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

//// 新春CDT兑换数据PV
func (this *StatisticsDoubleYearCdt) GetTotalPV(start_time *time.Time, end_time *time.Time) (int, error) {
	session := db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	pv, err := session.Sum(new(model.StatisticsDoubleYearCdt), "pv")
	return int(pv), err
}

// 新春CDT兑换数据UV
func (this *StatisticsDoubleYearCdt) GetTotalUV(start_time *time.Time, end_time *time.Time) (int, error) {
	session := db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	pv, err := session.Sum(new(model.StatisticsDoubleYearCdt), "uv")
	return int(pv), err
}

// 新春CDT兑换数据
func (this *StatisticsDoubleYearCdt) GetTotalCDT(start_time *time.Time, end_time *time.Time) (float32, error) {
	session := db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	cdt, err := session.Sum(new(model.StatisticsDoubleYearCdt), "cdt")
	return float32(cdt), err
}

// 新春用户CDT兑换数据
func (this *StatisticsDoubleYearUserDayCdt) Add(data_map *model.StatisticsDoubleYearUserDayCdt) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取新春用户CDT兑换数据
func (this *StatisticsDoubleYearUserDayCdt) GetData(userInfo string, start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearUserDayCdt, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_USERDAY_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_USERDAY_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearUserDayCdt
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 获取新春用户CDT兑换总数
func (this *StatisticsDoubleYearUserDayCdt) GetTotalCdt(userId string, start_time *time.Time, end_time *time.Time) (float32, error) {
	session := db.Mysql.Table(STATISTICS_DOUBLEYEARCDT_USERDAY_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}
	cdt, err := session.And("user_id = ?", userId).Sum(new(model.StatisticsDoubleYearUserDayCdt), "cdt")
	return float32(cdt), err
}

// 新春活动碎片兑换
func (this *StatisticsDoubleYearFragment) Add(data_map *model.StatisticsDoubleYearFragment) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春活动碎片兑换
func (this *StatisticsDoubleYearFragment) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearFragment, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_FRAGMENT_TABLE).Where("1 = 1")

	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLE_YEAR_FRAGMENT_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearFragment
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 新春活动用户碎片兑换记录
func (this *StatisticsDoubleYearUserFragment) Add(data_map *model.StatisticsDoubleYearUserFragment) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春活动用户碎片兑换记录
func (this *StatisticsDoubleYearUserFragment) GetData(userInfo string, start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearUserFragment, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_FRAGMENT_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_FRAGMENT_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearUserFragment
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 获取最后一条日统计数据
func (this *StatisticsDoubleYearUserFragment) GetLastData(userId string) (model.StatisticsDoubleYearUserFragment, error) {
	data := model.StatisticsDoubleYearUserFragment{}
	_, err := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_FRAGMENT_TABLE).Where("user_id = ?", userId).Desc("id").Limit(1).Get(&data)
	return data, err
}

// 新春活动日排行
func (this *StatisticsDoubleYearDailyRanking) Add(data_map *model.StatisticsDoubleYearDailyRanking) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春活动日排行
func (this *StatisticsDoubleYearDailyRanking) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearDailyRanking, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_DAILY_RANKING_TABLE).Where("1 = 1")

	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLE_YEAR_DAILY_RANKING_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearDailyRanking
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 新春活动用户日排行
func (this *StatisticsDoubleYearUserDailyRanking) Add(data_map *model.StatisticsDoubleYearUserDailyRanking) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春活动用户日排行
func (this *StatisticsDoubleYearUserDailyRanking) GetData(userInfo string, start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearUserDailyRanking, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_DAILY_RANKING_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_DAILY_RANKING_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearUserDailyRanking
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 新春活动总排行
func (this *StatisticsDoubleYearTotalRanking) Add(data_map *model.StatisticsDoubleYearTotalRanking) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春活动总排行
func (this *StatisticsDoubleYearTotalRanking) GetData(start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearTotalRanking, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_TOTAL_RANKING_TABLE).Where("1 = 1")

	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLE_YEAR_TOTAL_RANKING_TABLE).Where("1 = 1")
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearTotalRanking
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

// 新春活动用户总排行
func (this *StatisticsDoubleYearUserTotalRanking) Add(data_map *model.StatisticsDoubleYearUserTotalRanking) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	_, err = session.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 新春活动用户总排行
func (this *StatisticsDoubleYearUserTotalRanking) GetData(userInfo string, start_time *time.Time, end_time *time.Time, page, size int) (int, int, []model.StatisticsDoubleYearUserTotalRanking, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_TOTAL_RANKING_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	count, err := session.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	session = db.Mysql.Table(STATISTICS_DOUBLE_YEAR_USER_TOTAL_RANKING_TABLE).Where("1 = 1")
	if len(userInfo) > 0 {
		session.And("(user_name like '%" + userInfo + "%' OR user_id like '%" + userInfo + "%')")
	}
	if start_time != nil {
		session.And("date >= ?", utils.Time2Str(*start_time))
	}
	if end_time != nil {
		session.And("date < ?", utils.Time2Str(*end_time))
	}

	var data []model.StatisticsDoubleYearUserTotalRanking
	err = session.Desc("date").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}
