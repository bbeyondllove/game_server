package db_service

import (
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
	"time"
)

const (
	// ChangTypeAdd 增加操作.
	ChangTypeAdd = 1
	// ChangTypeREduce 减少操作.
	ChangTypeREduce = 2
	// ModelUserTable 用户表名.
	ModelUserTable = "t_user"
	// CdtOneDayCachePre 用户当天领取的cdt缓存健名.
	CdtOneDayCachePre = "cdt_"
)

// Cdt cdt货币对应的service.
type Cdt struct {
	cr *model.CdtRecord
}

// UpdateUserCdt 更新用户的cdt.
// 有多种情况：
// 1. 判断用户当天的cdt(当天已增加)+cdtVale(当天这次新增)之和)是否已经超过了限制.
// 如果超过限制，最多只能添加 LimitDayCdt - dayIncrementCdt(当天已获取的cdt)　的差值.
// 2.用户消费cdt时，需要判断用户账户cdt的余额是否足够.
// 3.特殊情况：
// 	3.1 针对个别活动的产出cdt，有另外一个限制值SpecialLimitDayCdt.
func (c *Cdt) UpdateUserCdt(userId string, cdtValue float32, eventType int) (statusCode int32, lastCdt float32) {
	// 判断是否是特殊消息类型.
	cdtValue = c.specialEventCdtValue(cdtValue, eventType)

	dayIncrementCdt := c.getUserOneDayCdt(userId)
	if cdtValue > 0 && c.isUpperLimit(dayIncrementCdt, eventType) {
		logger.Infof("UpdateUserCdt is limit[userId:%v, addedCdt:%v, limitCdt:%v].", userId, dayIncrementCdt, c.LimitDayCdt(eventType))
		return errcode.ERROR_CDT_DAY_FULL, 0
	}

	// 根据user_id获取用户信息.
	user := &User{}
	userData, err := user.GetDataByUid(userId)
	if err != nil {
		logger.Errorf("GetDataByUid error: %v\n", err.Error())
		return errcode.ERROR_SYSTEM, 0
	}

	sql := "UPDATE " + ModelUserTable + " SET cdt=? WHERE user_id=? LIMIT 1"
	var finalCdtString string
	var changeCdt, finalCdt float32
	var direction int
	if cdtValue > 0 {
		// 增加用户的cdt.
		direction = ChangTypeAdd
		if c.isUpperLimit(dayIncrementCdt+cdtValue, eventType) {
			finalCdt = userData.Cdt + (c.LimitDayCdt(eventType) - dayIncrementCdt)
			changeCdt = c.LimitDayCdt(eventType) - dayIncrementCdt
		} else {
			finalCdt = userData.Cdt + cdtValue
			changeCdt = cdtValue
		}
		finalCdtString = strconv.FormatFloat(float64(finalCdt), 'f', 4, 32)
		// 用户消费时，需要判断用户剩余的cdt是否足够.
		// 注：消费时cdtValue的值为负数, 直接相加.
	} else {
		direction = ChangTypeREduce
		finalCdt = userData.Cdt + cdtValue
		if finalCdt < 0 {
			logger.Info("user cdt is lack,", "srcCdt=", userData.Cdt, "reduce=", cdtValue)
			return errcode.ERROR_CDT_LACK_OF_BALANCE, userData.Cdt
		}
		changeCdt = cdtValue
		finalCdtString = strconv.FormatFloat(float64(finalCdt), 'f', 4, 32)
	}

	_, err = db.Mysql.Exec(sql, finalCdtString, userId)
	if err != nil {
		logger.Errorf("updateUserCdt[%v] cdt error: %v\n", sql, err.Error())
		return errcode.ERROR_SYSTEM, userData.Cdt
	}

	// 更新用户存储在redis中的最新cdt信息.
	c.updateUserInfoInRedis(userId, finalCdt)

	// 记录更新数据.
	c.cr.UserId = userId
	c.cr.SrcCdt = userData.Cdt
	c.cr.ChangeCdt = changeCdt
	c.cr.DestCdt = finalCdt
	c.cr.Direction = direction
	c.cr.EventType = eventType
	c.cr.CdtUsdRate = c.GetCdtUsdRate()
	c.cr.UsdCnyRate = c.GetUsdCnyRate()
	c.cr.CreateTime = int(time.Now().Unix())
	c.updateUserOneDayCdt(userId, dayIncrementCdt, changeCdt)

	logger.Info("updateUserCdt success, value of cdt is: ", strconv.FormatFloat(float64(changeCdt), 'f', 4, 32))
	return errcode.MSG_SUCCESS, finalCdt
}

// specialEventCdt 针对特殊类型消息，cdt值变化.
func (c *Cdt) specialEventCdtValue(cdtValue float32, eventType int) float32 {
	eventConfig := c.getCdtEventTypeConfig(eventType)
	if eventConfig == nil {
		return cdtValue
	}

	// 判断是否在有效时间内
	if eventConfig["expire"] != -1 {
		currentTime := time.Now()
		if currentTime.Unix() < utils.Str2Time(eventConfig["start_time"].(string)).Unix() || currentTime.Unix() > utils.Str2Time(eventConfig["end_time"].(string)).Unix() {
			return cdtValue
		}
	}

	return cdtValue * float32(eventConfig["several_fold"].(int))
}

// GetCdtEventTypeConfig 获取指定消息类型cdt配置.
func (c *Cdt) getCdtEventTypeConfig(eventType int) map[string]interface{} {
	for _, v := range base.Setting.Cdt.EventType {
		if v["event_number"] == eventType {
			return v
		}
	}
	return nil
}

// updateUserInfoInRedis 更新redis中用户的cdt最新值.
func (c *Cdt) updateUserInfoInRedis(userId string, lastCdt float32) {
	redisClient := db.RedisMgr.GetRedisClient()
	_, err := redisClient.HSet(userId, "cdt", strconv.FormatFloat(float64(lastCdt), 'f', 4, 32)).Result()
	if err != nil {
		logger.Errorf("update lastCdt of user to redis failed! error: %v\n", err.Error())
	}
}

// isUpperLimit 判断用户当天新增加cdt是否达到限制.
func (c *Cdt) isUpperLimit(userOneCdt float32, eventType int) bool {
	return userOneCdt >= c.LimitDayCdt(eventType)
}

// updateCdtRecord 记录cdt变更记录, 这里redis只记录获取的，不记录消费的.
func (c *Cdt) updateUserOneDayCdt(userId string, oldChangeCdt, newChangeCdt float32) {
	if newChangeCdt > 0 {
		redisClient := db.RedisMgr.GetRedisClient()
		key := CdtOneDayCachePre + userId
		_, err := redisClient.Set(key, strconv.FormatFloat(float64(oldChangeCdt+newChangeCdt), 'f', 4, 32), c.GetCacheExpireTime(1)).Result()
		if err != nil {
			logger.Errorf("update changeCdt of user to redis failed! error: %v\n", err.Error())
		}
	}

	// 变更记录入库.
	c.cr.Insert()
}

// GetCacheExpireTime 获取变动的cdt过期时间.
// 目前规则是每天最多只能获取LimitDayCdt，
// 过期时间为当天当前时间到第二天凌晨时间之差.
func (c *Cdt) GetCacheExpireTime(day int) time.Duration {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	nextDay := t.Unix() + int64(86400*day)

	return time.Duration(nextDay-time.Now().Unix()) * time.Second
}

// getUserOneDayCdt 获取用户当天已得cdt总数.
// 用户每天获取的cdt数可存储在redis，减少db压力.
// 当redis中没有时，查询db.
func (c *Cdt) getUserOneDayCdt(userId string) float32 {
	redisClient := db.RedisMgr.GetRedisClient()
	key := CdtOneDayCachePre + userId
	changeCdt, err := redisClient.Get(key).Float32()
	if err != nil {
		logger.Errorf("get changeCdt of user from redis failed! error: %v\n", err.Error())
		cdt, err := c.cr.GetUserOneDayCdt(userId, time.Unix(time.Now().Unix(), 0).Format("2006-01-02"), 0)
		// 当出错时，直接返回已经达到了限制值，确保不让用户增加cdt.
		if err != nil {
			return c.LimitDayCdt(proto.MSG_NULL_ACT)
		}
		return cdt
	}
	return changeCdt
}

// GetCdtUsdRate cdt对美元汇率.
func (c *Cdt) GetCdtUsdRate() float32 {
	return base.Setting.Cdt.CdtUsdRateDefault
}

// GetUsdCnyRate 美元对人民币汇率.
func (c *Cdt) GetUsdCnyRate() float32 {
	return base.Setting.Cdt.UsdCnyRateDefault
}

// GetCnyPointRate 人民币对积分(point)比率.
func (c *Cdt) GetCnyPointRate() float32 {
	return base.Setting.Cdt.CnyPointRateDefault
}

// LimitDayCdt 每天领取cdt上限.
// 针对特殊类型消息，有特定的上限值.
func (c *Cdt) LimitDayCdt(eventType int) float32 {
	eventConfig := c.getCdtEventTypeConfig(eventType)
	if eventConfig == nil {
		return base.Setting.Cdt.LimitDayCdt
	}

	return float32(eventConfig["limit_cdt"].(float64))
}

// PointToCdt 积分转换为cdt.
// 对换规则：2000积分＝1CNY,1CNY=(1/x*Y)CDT,
// x为cdt交易所中一个cdt的实时价格，单位为美元.
// y为一美元对应人民币的实时汇率.
func (c *Cdt) PointToCdt(point int) float32 {
	// 一个积分对应的cdt数量.
	onePointPrice := 1 / (c.GetCdtUsdRate() * c.GetUsdCnyRate() * c.GetCnyPointRate())
	return float32(point) * onePointPrice
}

// TransferPointToCdtFromTable 转换表中的point字段为cdt.
// 采用批量查询并更新.
func (c *Cdt) TransferPointToCdtFromTable(tableName, fromColumnName, toColumnName string) {
	// 获取总数.
	batchSize := 100
	r, err := db.Mysql.QueryString("SELECT COUNT(*) as total FROM " + tableName)
	if err != nil {
		logger.Errorf("get the total of table[%v] failed:", tableName, err)
		return
	}
	total, _ := strconv.Atoi(r[0]["total"])

	updateSql := "UPDATE " + tableName + " SET " + toColumnName + "=? WHERE id=? LIMIT 1"

	// 批量查询.
	var i int
	for i = 0; i < total; i += batchSize {
		sql := "SELECT id, " + fromColumnName + " FROM " + tableName
		results, err := db.Mysql.Limit(batchSize, i).QueryString(sql)
		if err != nil {
			logger.Errorf("query batch data from table[%s] failed! error: %v\n", tableName, err)
			return
		}
		if len(results) <= 0 {
			break
		}

		for _, v := range results {
			// fromColumnName字段类型可能是int,也可能是float.
			// 要做兼容处理.
			tmpValue := v[fromColumnName]
			p, err := strconv.Atoi(tmpValue)
			if err != nil {
				// 尝试从string转float64.
				tmp, err := strconv.ParseFloat(tmpValue, 64)
				if err != nil {
					logger.Infof("parse the point[%v] value of ID[%v] from string to int fail\n", p, v["id"])
					continue
				}
				p = int(tmp)
			}
			if p <= 0 {
				logger.Infof("the point[%v] value of ID[%v] is negative\n", p, v["id"])
				continue
			}

			cdt := c.PointToCdt(p)
			_, err = db.Mysql.Exec(updateSql, cdt, v["id"])
			if err != nil {
				logger.Warnf("the ID[%v] point[%v] to cdt failed, error: ", v["id"], p, err)
				continue
			}
			logger.Infof("the ID[%v] point[%v] to cdt[%v] success", v["id"], p, cdt)
		}
	}

	logger.Infof("%v done", tableName)
}

// GetUserCdt 查询用户拥有的cdt.
func (c *Cdt) GetUserCdt(userId string) float32 {
	client := db.RedisMgr.GetRedisClient()
	cdt, err := client.HGet(userId, "cdt").Float32()
	if err != nil {
		result, err := db.Mysql.Query("SELECT cdt FROM t_user where user_id=? LIMIT 1", userId)
		if err != nil {
			logger.Errorf("get cdt from t_user failed! error: %v\n", err)
			return 0
		}
		cdtF, err := strconv.ParseFloat(string(result[0]["cdt"]), 32)
		if err != nil {
			return 0
		}
		cdt = float32(cdtF)
	}
	return cdt
}

// NewCdt 实例化cdt.
func NewCdt() *Cdt {
	return &Cdt{
		cr: model.NewCdtRecord(),
	}
}
