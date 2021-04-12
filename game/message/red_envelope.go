package message

import (
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/message/activity/activity_roles"
	"game_server/game/message/activity/red_envelope"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
	"time"
)

const (
	repeat             = "SendRedEnvelope:"       // 防止重复调用
	RedEnvelopTotal    = "SendRedEnvelope_total:" //红包总数
	SendRedEnvelopeDay = "SendRedEnvelope_day:"   //当天发放的
)

// 发送红包
func SendRedEnvelope(userId string, roloId int) {

	//防止重复调用
	key := repeat + userId
	isTrue, _ := db.RedisMgr.GetRedisClient().SetNX(key, 1, 5*time.Second).Result()
	if isTrue == false {
		return
	}
	//5秒后推送
	//time.Sleep(time.Second * 5)
	//检查用户的角色ID是否是金童玉女
	goldenCouple := activity_roles.NewGoldenCouple()
	if goldenCouple.CheckIsGoldenCoupleRole(roloId) == false {
		return
	}
	//检查角色是否已经过期
	isExpire, _ := goldenCouple.Roles.CheckRolesIsExpire(userId, roloId)
	if isExpire == false {
		//角色已过期
		return
	}
	//检查红包是否发送
	keyDay := SendRedEnvelopeDay + time.Now().Format("20060102") + ":" + userId
	dayNum, _ := db.RedisMgr.GetRedisClient().Get(keyDay).Int64()
	if dayNum >= 1 {
		//红包已经发送
		return
	}
	if dayNum == 0 {
		has, _, _ := db_service.ActivityRedEnvelopeIns.GetSameDayRedEnvelope(userId)
		if has == true {
			//红包已经发送
			db.RedisMgr.GetRedisClient().Set(keyDay, 1, time.Hour*24) // 当天是否发放
			return
		}
	}
	// 判断红包数量
	keyCount := RedEnvelopTotal + userId
	redEnvelopeNum, err := db.RedisMgr.GetRedisClient().Get(keyCount).Int64()
	if err != nil || redEnvelopeNum == 0 {
		//检查红包数量
		redEnvelopeNum, err = db_service.ActivityRedEnvelopeIns.GetCountDayRedEnvelope(userId)
		if err != nil {
			//查询红包数量
			return
		}
		db.RedisMgr.GetRedisClient().Set(keyCount, redEnvelopeNum, time.Hour*24*3) // 发放的红包数量
	}
	if redEnvelopeNum >= 3 {
		logger.Debugf("SendRedEnvelope number is 3, userId:[%v]", userId)
		return
	}

	//获取红包
	redEnvelope := red_envelope.NewRedEnvelope()
	award := redEnvelope.RandCdtAward(red_envelope.RedEnvelopeRandTable)
	redEnvenlope := model.ActivityRedEnvelope{
		UserId:    userId,
		Number:    award,
		IsReceive: 1,
	}
	//增加总的
	code, finalCdt := db_service.NewCdt().UpdateUserCdt(userId, award, proto.MSG_SEND_RED_ENVELOPE_RSP) // 宝箱领取奖励
	if code != 0 {
		return
	}
	_, err = db_service.ActivityRedEnvelopeIns.Add(&redEnvenlope)
	if err != nil {
		logger.Errorf("SendRedEnvelope ActivityRedEnvelopeIns failed userId:[%v] err= %v", userId, err)
		return
	}
	db.RedisMgr.GetRedisClient().Incr(keyCount)               // 红包的总数量
	db.RedisMgr.GetRedisClient().Set(keyDay, 1, time.Hour*24) // 当天是否发放
	totalCdt, _ := FormatFloat(float64(finalCdt), 4)          // 有写有 5 位小数， 强制舍弃
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SEND_RED_ENVELOPE_RSP)
	response := &proto.S2CRedEnvelope{}
	response.Code = errcode.MSG_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.Cdt = award
	response.CdtalCdt = float32(totalCdt)
	rsp.WriteData(response)
	//推送消息给用户
	Sched.SendToUser(userId, rsp)

}

// 金童玉女推送红包
func ColdenCoupleRedEnvelope(userId string) {
	roleIdStr, err := db.RedisMgr.GetRedisClient().HGet(userId, "role_id").Result()
	if err == nil {
		roleId, _ := strconv.Atoi(roleIdStr)
		SendRedEnvelope(userId, roleId)
	}
}
