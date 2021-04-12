package model

import (
	"game_server/core/logger"
	"game_server/db"
	"time"
)

// CdtRecord ctd变动表对应的结构体.
type CdtRecord struct {
	Id         int     `xorm:"int(20) notnull"  json:"id" desc:"ID"`
	UserId     string  `xorm:"int(20) notnull" json:"user_id" `
	SrcCdt     float32 `xorm:"decimal(12,4)" json:"src_cdt" desc:"变动前的ctd"`
	ChangeCdt  float32 `xorm:"decimal(12,4)" json:"change_cdt" desc:"变动的cdt"`
	DestCdt    float32 `xorm:"decimal(12,4)" json:"dest_cdt" desc:"变动后的cdt"`
	Direction  int     `xorm:"int(20)" json:"direction" desc:"加减状态，1加，2减"`
	EventType  int     `xorm:"int(20)" json:"event_type" desc:"事件类型"`
	CdtUsdRate float32 `xorm:"decimal(8,6)" json:"cdt_usd_rate" desc:"当前cdt对换美元比率"`
	UsdCnyRate float32 `xorm:"decimal(8,6)" json:"usd_cny_rate" desc:"当前美元对换人民币比率"`
	CreateTime int     `xorm:"int(11)" json:"create_time" desc:"创建时间戳"`
}

// Insert 插入一条记录.
func (c *CdtRecord) Insert() bool {
	_, err := db.Mysql.Insert(c)
	if err != nil {
		logger.Errorf("cdtRecord insert a record failed: %v, data:%v", err, *c)
		return false
	}
	return true
}

// getUserOneDayCdt 查询某个用户某天获得的cdt量.
// 注!!!：这里只查询用户获取的cdt量，不包括消费.
// 这非常重要!!!
func (c *CdtRecord) GetUserOneDayCdt(userId, date string, eventType int) (float32, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
		return 0, err
	}

	var currentCdt float64
	endTime := startTimeStr.Unix() + 86400
	//db.Mysql.ShowSQL(true)
	switch eventType {
	// 0表示所有事件类型.
	case 0:
		currentCdt, err = db.Mysql.Where("user_id=?", userId).And("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	// 指定事件类型查询.
	default:
		currentCdt, err = db.Mysql.Where("user_id=?", userId).And("event_type=?", eventType).And("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	}
	if err != nil {
		logger.Errorf("getUserOneDayCdt data failed: %v", err)
		return 0, err
	}

	return float32(currentCdt), nil
}

// GetCdtByEvent 查询指定天数的事件类型cdt之和.
func (c *CdtRecord) GetCdtByEvent(eventType int, date string, day int) (float32, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("GetCdtByEvent parse date failed: %v", err)
		return 0, err
	}
	endTime := startTimeStr.Unix() + int64(86400*day)
	currentCdt, err := db.Mysql.Where("event_type=?", eventType).And("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	if err != nil {
		logger.Errorf("GetCdtByEvent data failed: %v", err)
		return 0, err
	}

	return float32(currentCdt), nil
}

// 某天CDT产出数
func (c *CdtRecord) GetOneDayCdt(date string, eventType int) (float32, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
		return 0, err
	}

	var currentCdt float64
	endTime := startTimeStr.Unix() + 86400
	//db.Mysql.ShowSQL(true)
	switch eventType {
	// 0表示所有事件类型.
	case 0:
		currentCdt, err = db.Mysql.Where("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	// 指定事件类型查询.
	default:
		currentCdt, err = db.Mysql.Where("event_type=?", eventType).And("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	}
	if err != nil {
		logger.Errorf("getUserOneDayCdt data failed: %v", err)
		return 0, err
	}

	return float32(currentCdt), nil
}

func (c *CdtRecord) GetUserTotalCdt(date, userId string, eventType int) (float32, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
		return 0, err
	}

	var currentCdt float64
	endTime := startTimeStr.Unix() + 86400
	//db.Mysql.ShowSQL(true)
	switch eventType {
	// 0表示所有事件类型.
	case 0:
		currentCdt, err = db.Mysql.Where("change_cdt >?", 0).And("user_id = ?", userId).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	// 指定事件类型查询.
	default:
		currentCdt, err = db.Mysql.Where("event_type=?", eventType).And("change_cdt >?", 0).And("user_id = ?", userId).And("create_time <?", endTime).Sum(new(CdtRecord), "change_cdt")
	}
	if err != nil {
		logger.Errorf("getUserOneDayCdt data failed: %v", err)
		return 0, err
	}

	return float32(currentCdt), nil
}

// 获取总的CDT数量
func (c *CdtRecord) GetTotalCdt() (float32, error) {
	totalCdt, err := db.Mysql.Sum(new(CdtRecord), "change_cdt")
	if err != nil {
		logger.Errorf("GetCdtByEvent data failed: %v", err)
		return 0, err
	}

	return float32(totalCdt), nil
}

// 获取总的CDT数量
func (c *CdtRecord) GetTotalCdtByType(eventType int) (float32, error) {
	totalCdt, err := db.Mysql.Where("event_type=?", eventType).Sum(new(CdtRecord), "change_cdt")
	if err != nil {
		logger.Errorf("GetCdtByEvent data failed: %v", err)
		return 0, err
	}

	return float32(totalCdt), nil
}

// 获取总的CD产出数量
func (c *CdtRecord) GetTotalCdtProduce() (float32, error) {
	totalCdt, err := db.Mysql.Where("change_cdt >?", 0).Sum(new(CdtRecord), "change_cdt")
	if err != nil {
		logger.Errorf("GetCdtByEvent data failed: %v", err)
		return 0, err
	}

	return float32(totalCdt), nil
}

func (c *CdtRecord) GetOneDayCdtCount(date string, eventType int) (int, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
		return 0, err
	}

	var currentCdt int64
	endTime := startTimeStr.Unix() + 86400
	//db.Mysql.ShowSQL(true)
	switch eventType {
	// 0表示所有事件类型.
	case 0:
		currentCdt, err = db.Mysql.Where("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Count(new(CdtRecord))
	// 指定事件类型查询.
	default:
		currentCdt, err = db.Mysql.Where("event_type=?", eventType).And("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).Count(new(CdtRecord))
	}
	if err != nil {
		logger.Errorf("getUserOneDayCdt data failed: %v", err)
		return 0, err
	}

	return int(currentCdt), nil
}

// 获取总的CDT记录数量
func (c *CdtRecord) GetCdtCountByType(eventType int) (int, error) {
	totalCdt, err := db.Mysql.Where("event_type=?", eventType).Count(new(CdtRecord))
	if err != nil {
		return 0, err
	}

	return int(totalCdt), nil
}

func (c *CdtRecord) GetOneDayPeopleCount(date string, eventType int) (int, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
		return 0, err
	}

	var currentCdt int64
	endTime := startTimeStr.Unix() + 86400
	//db.Mysql.ShowSQL(true)
	switch eventType {
	// 0表示所有事件类型.
	case 0:
		currentCdt, err = db.Mysql.Where("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).GroupBy("user_id").Count(new(CdtRecord))
	// 指定事件类型查询.
	default:
		currentCdt, err = db.Mysql.Where("event_type=?", eventType).And("change_cdt >?", 0).And("create_time >=?", startTimeStr.Unix()).And("create_time <?", endTime).GroupBy("user_id").Count(new(CdtRecord))
	}
	if err != nil {
		logger.Errorf("getUserOneDayCdt data failed: %v", err)
		return 0, err
	}

	return int(currentCdt), nil
}

// 获取用户某个时间获取的CDT量
func (c *CdtRecord) GetUserCdt(userId, beginDate, endDate string, eventType int) (float32, error) {
	var startTimeStr time.Time
	var endTimeStr time.Time
	var err error
	if len(beginDate) > 0 {
		startTimeStr, err = time.ParseInLocation("2006-01-02", beginDate, time.Local)
		if err != nil {
			logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
			return 0, err
		}
	}
	if len(endDate) > 0 {
		endTimeStr, err = time.ParseInLocation("2006-01-02", endDate, time.Local)
		if err != nil {
			logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
			return 0, err
		}
	}

	session := db.Mysql.Where("user_id=?", userId).And("change_cdt >?", 0)
	var currentCdt float64
	switch eventType {
	case 0:
	default:
		// 指定事件类型查询.
		session.And("event_type=?", eventType)
	}
	if len(beginDate) > 0 {
		session.And("create_time >=?", startTimeStr.Unix())
	}
	if len(endDate) > 0 {
		session.And("create_time <?", endTimeStr.Unix())
	}
	currentCdt, err = session.Sum(new(CdtRecord), "change_cdt")
	if err != nil {
		logger.Errorf("getUserOneDayCdt data failed: %v", err)
		return 0, err
	}

	return float32(currentCdt), nil
}

// 某天CDT产出数
func (c *CdtRecord) UserOneDayCdt(date string, eventType int) ([]map[string]string, error) {
	startTimeStr, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
		return nil, err
	}
	endTime := startTimeStr.Unix() + 86400

	switch eventType {
	// 0表示所有事件类型.
	case 0:
		sql := "SELECT user_id,sum(change_cdt) AS cdt FROM  t_cdt_record WHERE create_time >= ? AND create_time < ? GROUP BY user_id "
		result, err := db.Mysql.SQL(sql, startTimeStr.Unix(), endTime).QueryString()
		return result, err

	// 指定事件类型查询.
	default:
		sql := "SELECT user_id,sum(change_cdt) AS cdt FROM  t_cdt_record WHERE event_type = ? AND create_time >= ? AND create_time < ? GROUP BY user_id "
		result, err := db.Mysql.SQL(sql, eventType, startTimeStr.Unix(), endTime).QueryString()
		return result, err
	}
}

// 每天领取奖励的数量
func (this *CdtRecord) UserTotalCdt(beginDate string, endDate string, eventType int) ([]map[string]string, error) {
	var startTimeStr time.Time
	var endTimeStr time.Time
	var err error
	if len(beginDate) > 0 {
		startTimeStr, err = time.ParseInLocation("2006-01-02", beginDate, time.Local)
		if err != nil {
			logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
			return nil, err
		}
	}
	if len(endDate) > 0 {
		endTimeStr, err = time.ParseInLocation("2006-01-02", endDate, time.Local)
		if err != nil {
			logger.Errorf("getUserOneDayCdt parse date failed: %v", err)
			return nil, err
		}
	}

	sql := "SELECT user_id,sum(change_cdt) AS cdt FROM  t_cdt_record WHERE event_type = ? AND create_time >= ? AND create_time < ? GROUP BY user_id "
	result, err := db.Mysql.SQL(sql, eventType, startTimeStr.Unix(), endTimeStr.Unix()).QueryString()
	return result, err
}

// NewCdtRecord 实例化CdtRecord.
func NewCdtRecord() *CdtRecord {
	return &CdtRecord{}
}
