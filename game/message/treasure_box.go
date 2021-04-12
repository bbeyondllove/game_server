package message

import (
	"encoding/json"
	"fmt"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/message/activity/treasure_box"
	"game_server/game/message/statistical"
	"game_server/game/proto"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shopspring/decimal"
)

//redis key
const (
	//箱子上锁
	BOX_IS_LOCK = "treasure_box_lock:"
	//宝箱奖励
	BOX_AWARD = "treasure_box_award:userId:"
	//防止重复调用
	BOX_AWARD_ID = "treasure_box_award:id:"
	//1分钟打开宝箱的次数
	BOX_OPEN_NUM = "treasure_box_open_num:"
	// 每天打开的次数
	BOX_DAY_OPEN_NUM = "treasure_box_day_num:"
)

// 宝箱
type StreasureBox struct {
}

var (
	G_StreasureBox_Cdt_Rand_Table      sync.Map
	G_StreasureBox_Config_Reload_Token int32 = 0 // 宝箱配置重加载令牌
	G_StreasureBox_Config_Init_Time    int64 = 0 // 宝箱配置初始化时间
)

const (
	StreasureBoxCdtRandTable = "StreasureBoxCdtRandTable"
)

func TreasureBox_Config(cfg *proto.BaseConf) {
	if cfg.TreasureBox.CdtRandTableConfigFrom == 0 {
		InitCdtRandTable()
	} else {
		InitCdtRandTableFromDB()
	}
}

//打开宝箱
func (a *agent) OpenStreasureBox(requestMsg *utils.Packet) {
	logger.Infof("OpenStreasureBox in request:", requestMsg.GetBuffer())
	base := treasure_box.NewBase()
	msg := &proto.C2SOpenStreasureBox{}
	// 返回
	response := &proto.S2COpenStreasureBox{}
	response.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	response.Message = errcode.ERROR_MSG[response.Code]
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("OpenStreasureBox json.Unmarshal error, err=", err.Error())
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 获取用户信息
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("OpenStreasureBox token error", msg.Token)
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	//检查宝箱每天打开的总次数
	exceptionDayKey := BOX_DAY_OPEN_NUM + time.Now().Format("20060102") + ":" + payLoad.UserId
	exceptionDayNum := db.RedisMgr.Get(exceptionDayKey)
	if exceptionDayNum == "" {
		exceptionDayNum = "0"
	}
	exceptionDayNumInt, _ := strconv.Atoi(exceptionDayNum)
	if exceptionDayNumInt >= G_BaseCfg.TreasureBox.ExceptionDayCount {
		// 宝箱已到达今日上限
		response.Code = errcode.ERROR_BOX_DAY_OPEN_NUM
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}

	//检查宝箱每分钟次数
	openNumKey := BOX_OPEN_NUM + payLoad.UserId
	number := db.RedisMgr.Get(openNumKey)
	if number == "" {
		number = "0"
	}
	openNum, _ := strconv.Atoi(number)
	if openNum >= G_BaseCfg.TreasureBox.ExceptionCount {
		response.Code = errcode.ERROR_BOX_OPEN_NUM
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 活动类型错误
	if msg.ActivityType != proto.ACTIVITY_TYPE_TREASURE_BOX {
		logger.Errorf("OpenStreasureBox msg.ActivityType error", msg.ActivityType)
		response.Code = errcode.ERROR_OBJ_NOT_EXISTS
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 判断箱子是否存在
	eventDate := getEventData(msg.ActivityType, msg.LocationID, msg.X, msg.Y)
	if eventDate == nil {
		logger.Errorf("OpenStreasureBox eventDate is null", eventDate)
		response.Code = errcode.ERROR_OBJ_NOT_EXISTS
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	//获取坐标
	positionNo := getPositionNo(msg.X, msg.Y)
	//判断箱子有没有上锁
	key := BOX_IS_LOCK + strconv.Itoa(positionNo)
	//箱子上锁
	lock, err := db.RedisMgr.GetRedisClient().SetNX(key, payLoad.UserId, time.Duration(G_BaseCfg.TreasureBox.LockTimeout)*time.Second).Result()
	if lock == false {
		//箱子已被锁定
		logger.Errorf("OpenStreasureBox is LOCK", payLoad.UserId)
		response.Code = errcode.ERROR_BOX_IS_LOCK
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 记录宝箱开启时间
	awardKey := BOX_AWARD + payLoad.UserId
	//db.RedisMgr.GetRedisClient().HSet(awardKey, "openTime", time.Now().Format("2006-01-02 15:04:05")).Result()
	db.RedisMgr.GetRedisClient().HMSet(awardKey, map[string]interface{}{"openTime": time.Now().Format("2006-01-02 15:04:05"), "awardId": positionNo}).Result()
	//记录打开宝箱次数
	isBool, _ := db.RedisMgr.GetRedisClient().SetNX(openNumKey, 1, time.Duration(G_BaseCfg.TreasureBox.ExceptionInterval)*time.Second).Result()
	if !isBool {
		db.RedisMgr.GetRedisClient().Incr(openNumKey).Result()
	}

	// 标记打开宝箱总次数
	statistical.StatisticsDotIns.OpenStreasureBox(payLoad.UserId)

	//箱子已被锁定
	response.Code = errcode.MSG_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	dayNum, _ := db.RedisMgr.GetRedisClient().Incr(exceptionDayKey).Result()
	if dayNum == 1 {
		//设置过期时间  1 天
		db.RedisMgr.GetRedisClient().Expire(exceptionDayKey, time.Hour*24)
	}
	response.DayNum = dayNum
	response.ResidueDegree = int64(G_BaseCfg.TreasureBox.ExceptionDayCount - exceptionDayNumInt)
	rsp := base.ResponseMessage(proto.MSG_OPEN_STREASURE_BOX_RSP, response)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
	return
}

// 获取当当天CDT上限
func (a *agent) GetStreasureBoxDayNum(requestMsg *utils.Packet) {
	logger.Infof("GetStreasureBoxDayNum in request:", requestMsg.GetBuffer())
	base := treasure_box.NewBase()
	msg := &proto.C2SStreasureBoxDayNum{}
	// 返回
	response := &proto.S2CStreasureBoxDayNum{}
	response.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	response.Message = errcode.ERROR_MSG[response.Code]
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("OpenStreasureBox json.Unmarshal error, err=", err.Error())
		rsp := base.ResponseMessage(proto.MSG_TREASURE_BOX_DAY_NUM_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 获取用户信息
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("OpenStreasureBox token error", msg.Token)
		rsp := base.ResponseMessage(proto.MSG_TREASURE_BOX_DAY_NUM_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}

	//检查宝箱每天打开的总次数
	exceptionDayKey := BOX_DAY_OPEN_NUM + time.Now().Format("20060102") + ":" + payLoad.UserId
	exceptionDayNum := db.RedisMgr.Get(exceptionDayKey)
	if exceptionDayNum == "" {
		exceptionDayNum = "0"
	}
	exceptionDayNumInt, _ := strconv.Atoi(exceptionDayNum)
	response.Code = errcode.MSG_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.DayNum = int64(exceptionDayNumInt)
	response.ResidueDegree = int64(G_BaseCfg.TreasureBox.ExceptionDayCount - exceptionDayNumInt)
	if response.ResidueDegree == 0 {
		response.Code = errcode.ERROR_BOX_DAY_OPEN_NUM
		response.Message = errcode.ERROR_MSG[response.Code]
	}
	rsp := base.ResponseMessage(proto.MSG_TREASURE_BOX_DAY_NUM_RSP, response)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
}

//完成宝箱
func (a *agent) FinishStreasureBox(requestMsg *utils.Packet) {

	logger.Infof("FinishStreasureBox in request:", requestMsg.GetBuffer())
	base := treasure_box.NewBase()
	msg := &proto.C2SFinishStreasureBox{}
	// 返回
	response := &proto.S2CFinishStreasureBox{}
	response.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	response.Message = errcode.ERROR_MSG[response.Code]
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("FinishStreasureBox json.Unmarshal error, err=", err.Error())
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 获取用户信息
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("FinishStreasureBox token error", msg.Token)
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 活动类型错误
	if msg.ActivityType != proto.ACTIVITY_TYPE_TREASURE_BOX {
		logger.Errorf("FinishStreasureBox msg.ActivityType error", msg.ActivityType)
		response.Code = errcode.ERROR_OBJ_NOT_EXISTS
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 判断箱子是否存在
	eventDate := getEventData(msg.ActivityType, msg.LocationID, msg.X, msg.Y)
	if eventDate == nil {
		logger.Errorf("FinishStreasureBox box is null userID:[%v]", payLoad.UserId)
		response.Code = errcode.ERROR_OBJ_NOT_EXISTS
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	//获取坐标
	positionNo := getPositionNo(msg.X, msg.Y)
	key := BOX_IS_LOCK + strconv.Itoa(positionNo)
	lockUserId, _ := db.RedisMgr.GetRedisClient().Get(key).Result()
	//是否是当前用户锁定
	if lockUserId != "" && lockUserId != payLoad.UserId {
		logger.Errorf("FinishStreasureBox lockUserId:%v userid: %v", lockUserId, payLoad.UserId)
		response.Code = errcode.ERROR_REQUEST_NOT_ALLOW
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	if lockUserId == "" {
		// 判断是否是异常情况,超过 2 分钟的
		isTrue := exceptionReceive(payLoad.UserId, positionNo)
		if isTrue == false {
			logger.Errorf("FinishStreasureBox exceptionReceive lockUserId:%v userid: %v", lockUserId, payLoad.UserId)
			response.Code = errcode.ERROR_REQUEST_NOT_ALLOW
			response.Message = errcode.ERROR_MSG[response.Code]
			rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
			logger.Debugf(string(rsp.Bytes()))
			SendPacket(a.conn, rsp)
			return
		}
	}

	//判断完成状态
	if msg.FinishStatus != 1 {
		db.RedisMgr.GetRedisClient().Del(key)
		logger.Errorf("FinishStreasureBox not finished userid:[%v], FinishStatus:[%v]", payLoad.UserId, msg.FinishStatus)
		response.Code = errcode.ERROR_BOX_NOT_FINISHED
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	//事件重置
	resetEventData(msg.ActivityType, msg.LocationID, getPositionNo(msg.X, msg.Y), eventDate)
	//解锁
	db.RedisMgr.GetRedisClient().Del(key)
	//计算奖励
	f, _ := RandCdtAward().Float64()
	//ctd := float32(f)
	//记录宝箱生成的CDT
	awardKey := BOX_AWARD + payLoad.UserId
	_, err = db.RedisMgr.GetRedisClient().HMSet(awardKey, map[string]interface{}{"ctd": strconv.FormatFloat(f, 'f', -1, 64), "awardId": positionNo, "finishTime": time.Now().Format("2006-01-02 15:04:05")}).Result()
	// 奖励发放失败
	if err != nil {
		response.Code = errcode.ERROR_SYSTEM
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}

	// 标记打开宝箱总次数
	statistical.StatisticsDotIns.FinishStreasureBox(payLoad.UserId)

	itemInfo := GetAwardItem(1001)
	itemInfo.ItemNum = float32(f)
	itemInfo.AwardId = positionNo
	response.Code = errcode.MSG_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.BoxAward = itemInfo
	response.UserId = payLoad.UserId
	response.X = msg.X
	response.Y = msg.Y
	response.LocationID = msg.LocationID
	response.ActivityType = msg.ActivityType
	rsp := base.ResponseMessage(proto.MSG_FINISH_STREASURE_BOX_RSP, response)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)

	//通知其它玩家消除事件
	responseMessage := &proto.S2CTreasureBoxFinishEvent{}
	responseMessage.Message = ""
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.UserId = payLoad.UserId
	responseMessage.LocationID = msg.LocationID
	responseMessage.EventId = eventDate.EventId
	responseMessage.ActivityType = msg.ActivityType
	responseMessage.X = msg.X
	responseMessage.Y = msg.Y
	responseMessage.ItemInfos = make(map[int][]*proto.AwardItem, 0)
	broadMsg := &utils.Packet{}
	broadMsg.Initialize(proto.MSG_BROAD_FINISH_EVENT)
	broadMsg.WriteData(responseMessage)
	logger.Debugf(string(broadMsg.Bytes()))
	Sched.BroadCastMsg(int32(msg.LocationID), payLoad.UserId, broadMsg)
	return

}

// 特殊异常情况领取 （看广告，点击下载等待时间超过锁的时间）
func exceptionReceive(userId string, positionNo int) bool {
	// 检查奖励是否被用户打开过
	awardKey := BOX_AWARD + userId
	awardInfo := db.RedisMgr.HGetAll(awardKey)
	// 判断是否有时间
	openTime, ok := awardInfo["openTime"]
	if ok == false {
		return false
	}
	// 判断是否有奖励ID
	awardId, ok := awardInfo["awardId"]
	if ok == false {
		return false
	}
	// 判断奖励ID是否相等
	if strconv.Itoa(positionNo) != awardId {
		return false
	}
	// 计算打开宝箱 到当前的时间差
	currentTime := time.Now().Unix()
	loc, _ := time.LoadLocation("Local") //获取时区
	expTime, _ := time.ParseInLocation("2006-01-02 15:04:05", openTime, loc)
	seconds := currentTime - expTime.Unix()
	// 判断是否超过锁的时间
	if seconds < int64(G_BaseCfg.TreasureBox.LockTimeout) {
		return false
	}
	//箱子上锁
	keyLock := BOX_IS_LOCK + strconv.Itoa(positionNo)
	lock, _ := db.RedisMgr.GetRedisClient().SetNX(keyLock, userId, time.Duration(G_BaseCfg.TreasureBox.LockTimeout)*time.Second).Result()
	if lock == false {
		//箱子已被锁定
		return false
	}
	//记录打开宝箱次数
	openNumKey := BOX_OPEN_NUM + userId
	isBool, _ := db.RedisMgr.GetRedisClient().SetNX(openNumKey, 1, time.Duration(G_BaseCfg.TreasureBox.ExceptionInterval)*time.Second).Result()
	if !isBool {
		db.RedisMgr.GetRedisClient().Incr(openNumKey).Result()
	}
	return true
}

//领取宝箱奖励
func (a *agent) ReceiveBoxReward(requestMsg *utils.Packet) {
	logger.Infof("ReceiveBoxReward in request:", requestMsg.GetBuffer())
	base := treasure_box.NewBase()
	msg := &proto.C2SReceiveBoxReward{}
	// 返回
	response := &proto.S2CReceiveBoxReward{}
	response.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	response.Message = errcode.ERROR_MSG[response.Code]
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("ReceiveBoxReward json.Unmarshal error, err=", err.Error())
		rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 获取用户信息
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("ReceiveBoxReward token error", msg.Token)
		rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
		response.Code = errcode.ERROR_HTTP_SIGNATURE
		response.Message = errcode.ERROR_MSG[response.Code]
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	// 获取奖励信息
	awardKey := BOX_AWARD + payLoad.UserId
	awardInfo := db.RedisMgr.HGetAll(awardKey)
	if awardInfo == nil || len(awardInfo) == 0 {
		//奖励不存在
		logger.Errorf("ReceiveBoxReward not fund,", awardInfo)
		rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
		response.Code = errcode.ERROR_BOX_AWARD_NOT_FOUND
		response.Message = errcode.ERROR_MSG[response.Code]
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	//判断奖励ID是否存在
	awardId, _ := strconv.Atoi(awardInfo["awardId"]) // 奖励ID
	if awardId != msg.AwardId {
		logger.Errorf("userId:%s ReceiveBoxReward awardId error", payLoad.UserId, awardInfo)
		response.Code = errcode.ERROR_PARAM_ILEGAL
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}

	//防止重复调用
	isTrue, _ := db.RedisMgr.GetRedisClient().SetNX(BOX_AWARD_ID+awardInfo["awardId"], 1, 2*time.Second).Result()
	if isTrue == false {
		response.Code = errcode.ERROR_BOX_AWARD_EXCEPTION
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}

	//增加CDT
	ctd, _ := strconv.ParseFloat(awardInfo["ctd"], 10)
	isSantaClaus, _ := CheckCurrentRole(payLoad.UserId)
	eventType := proto.MSG_RECEIVE_STREASURE_BOX
	// todo 圣诞老人 检查当前角色是否是圣诞老人并且没有过期
	if isSantaClaus {
		eventType = proto.MSG_CHRISMAS_ROLE
	}
	code, finalCdt := db_service.NewCdt().UpdateUserCdt(payLoad.UserId, float32(ctd), eventType) // 宝箱领取奖励
	if code != errcode.MSG_SUCCESS {
		logger.Errorf("userId:%s ReceiveBoxReward UpdateUserCdt error: %s", payLoad.UserId, code)
		response.Code = errcode.ERROR_BOX_AWARD_EXCEPTION
		response.Message = errcode.ERROR_MSG[response.Code]
		rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
		SendPacket(a.conn, rsp)
		return
	}
	//清除奖励
	db.RedisMgr.GetRedisClient().Del(awardKey)
	//更新用户信息
	db_service.UserTreasureBoxRecordIns.UpdateBoxCtd(payLoad.UserId, float32(ctd), awardInfo["openTime"])
	totalCdt, _ := FormatFloat(float64(finalCdt), 4) // 有写有 5 位小数， 强制舍弃
	response.Code = errcode.MSG_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.Cdt = float32(totalCdt)
	rsp := base.ResponseMessage(proto.MSG_RECEIVE_STREASURE_BOX_RSP, response)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
}

// 获取Item
func GetAwardItem(itemId int) proto.BoxAward {
	var award proto.BoxAward
	G_ItemList.Range(func(key interface{}, value interface{}) bool {
		item := value.(*proto.ProductItem)
		if item.ItemId == itemId {
			award.ItemId = item.ItemId
			award.Desc = item.Desc
			award.ImgUrl = item.ImgUrl
			award.ItemName = item.ItemName
			return true
		}
		return true
	})
	return award
}

// 通过配置初始化随机表
func InitCdtRandTable() {
	last_level := int64(0)
	var cdt_rand_table []proto.TreasureBoxCdtRand = make([]proto.TreasureBoxCdtRand, 0)
	for _, item := range G_BaseCfg.TreasureBox.CdtRandTable {
		level := item.Probability + last_level
		last_level = level
		cdt_rand_table = append(cdt_rand_table, proto.TreasureBoxCdtRand{
			Probability:  level,
			RewardNumber: item.RewardNumber,
		})
	}

	G_StreasureBox_Cdt_Rand_Table.Store(StreasureBoxCdtRandTable, cdt_rand_table)
	G_StreasureBox_Config_Init_Time = time.Now().Unix()
}

// 通过数据库初始化随机表
func InitCdtRandTableFromDB() {
	data, err := db_service.TreasureBoxCdtConfigIns.GetAllData(1)
	if err != nil {
		fmt.Println("InitCdtRandTableFromDB error,err=" + err.Error())
		logger.Error(err)
		return
	}
	if len(data) == 0 {
		fmt.Println("TreasureBoxCdtConfig is empty")
		logger.Error("TreasureBoxCdtConfig is empty")
		return
	}

	last_level := int64(0)
	var cdt_rand_table []proto.TreasureBoxCdtRand = make([]proto.TreasureBoxCdtRand, 0)
	for _, item := range data {
		level := item.Probability + last_level
		last_level = level
		cdt_rand_table = append(cdt_rand_table, proto.TreasureBoxCdtRand{
			Probability:  level,
			RewardNumber: decimal.NewFromFloat32(item.RewardNumber),
		})
	}

	G_StreasureBox_Cdt_Rand_Table.Store(StreasureBoxCdtRandTable, cdt_rand_table)
	G_StreasureBox_Config_Init_Time = time.Now().Unix()
}

// 检查并重新初始化
func CheckAndFlushRandTable() {
	if G_BaseCfg.TreasureBox.CdtRandTableConfigFrom == 0 {
		return
	}
	// 检查是否超时
	if time.Now().Unix()-G_StreasureBox_Config_Init_Time <= G_BaseCfg.TreasureBox.CdtRandTableReloadTimeout {
		return
	}
	// 加锁
	swapped := atomic.CompareAndSwapInt32(&G_StreasureBox_Config_Reload_Token, 0, 1)
	if !swapped {
		return
	}
	// 从数据库初始化cdt随机表
	InitCdtRandTableFromDB()
	// 解锁
	atomic.StoreInt32(&G_StreasureBox_Config_Reload_Token, 0)
}

func GetCdtRandTable() []proto.TreasureBoxCdtRand {
	// CheckAndFlushRandTable()
	result, ok := G_StreasureBox_Cdt_Rand_Table.Load(StreasureBoxCdtRandTable)
	if ok {
		return result.([]proto.TreasureBoxCdtRand)
	}
	return nil
}

// 随机获取奖励
func RandCdtAward() decimal.Decimal {
	table := GetCdtRandTable()
	table_len := len(table)
	if table_len == 0 {
		return decimal.NewFromFloat32(G_BaseCfg.TreasureBox.CdtMin)
	}
	total_len := table[table_len-1].Probability
	rand.Seed(time.Now().UTC().UnixNano())

	award := decimal.NewFromFloat32(G_BaseCfg.TreasureBox.CdtMin)
	CdtSecondCritical := decimal.NewFromFloat32(G_BaseCfg.TreasureBox.CdtSecondCritical)
	retry_count := G_BaseCfg.TreasureBox.CdtSecondCriticalRetryCount
	for {
		idx := rand.Int63n(total_len)

		for _, item := range table {
			if idx < item.Probability {
				award = item.RewardNumber
				break
			}
		}
		if award.LessThan(CdtSecondCritical) {
			break
		}
		if retry_count <= 0 {
			break
		}
		retry_count--
	}
	return award
}

//获取宝箱奖励记录
func (a *agent) ReceiveBoxGetRecord(requestMsg *utils.Packet) {
	logger.Debugf("ReceiveBoxGetRecord in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_STREASURE_BOX_RECORD_RSP)

	response := &proto.S2CStreasureBoxRecordList{}
	response.Code = errcode.ERROR_PARAM_ILEGAL
	response.Message = errcode.ERROR_MSG[response.Code]

	msg := &proto.C2SStreasureBoxGetRecord{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(response)
		SendPacket(a.conn, rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(response)
		SendPacket(a.conn, rsp)
		return
	}
	totalRecords, err := db_service.UserTreasureBoxRecordIns.GetCount(payLoad.UserId)
	if err != nil {
		logger.Errorf("mysql error,err=", err.Error())
		rsp.WriteData(response)
		SendPacket(a.conn, rsp)
		return
	}

	// Pagging := proto.Paging{
	// 	TotalSize:   int(totalRecords),
	// 	CurrentSize: msg.Size,
	// 	CurrentPage: msg.Page,
	// }
	if totalRecords > 0 {
		records, err := db_service.UserTreasureBoxRecordIns.GetData(payLoad.UserId)
		// records, err := db_service.UserTreasureBoxRecordIns.GetPageData(payLoad.UserId, Pagging.CurrentPage, Pagging.CurrentSize)
		if err != nil {
			logger.Errorf("mysql error,err=", err.Error())
			rsp.WriteData(response)
			SendPacket(a.conn, rsp)
			return
		}
		var result []proto.StreasureBoxRecord = make([]proto.StreasureBoxRecord, 0)
		for _, record := range records {
			item := proto.StreasureBoxRecord{
				OpenTime:  record.OpenTime.Format("2006/01/02 15:04:05"),
				Cdt:       decimal.NewFromFloat32(record.Cdt),
				WatchTime: record.WatchTime,
			}
			result = append(result, item)
		}
		response.Record = result
	}

	response.Code = errcode.MSG_SUCCESS
	response.Message = ""
	// response.Paging = Pagging
	response.TotalCdt, _ = decimal.NewFromString(user_info["treasure_box_total_income"])
	response.WatchNum = totalRecords
	rsp.WriteData(response)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
	logger.Debugf("ReceiveBoxGetRecord end")
	return
}

func FormatFloat(num float64, decimal int) (float64, error) {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}
