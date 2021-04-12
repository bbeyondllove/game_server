package message

import (
	"encoding/json"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/errcode"
	"game_server/game/message/activity/double_year"
	"game_server/game/proto"
)

func (a *agent) ActivitStatus(requestMsg *utils.Packet) {
	logger.Debugf("ActivitStatus in request:", requestMsg.GetBuffer())

	responseMessage := &proto.S2CActivitStatus{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_ACTIVITY_STATYS_RSP)

	msg := &proto.C2SActivitStatus{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	// 欢乐宝箱状态
	activit := proto.Activity{}
	activit.TreasureBox = G_BaseCfg.TreasureBox.StateSwitch
	//新春活动 (排行榜)
	activityState := G_DoubleYearEvent.getActivityState()
	rankList := double_year.NewRankList()
	activit.RandList = rankList.GetRandListStatus()
	activit.ExchangeRole = activityState
	activit.ExchangeCdt = activityState

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	responseMessage.Activity = activit
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
	logger.Debugf("ActivitStatus end")
	return
}

//推送活动状态
//@ isAll  true 推送给所有用户
//@ userId 用户ID
func PushActivityStatus(isAll bool, userId string) {
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_ACTIVITY_STATUS_RSP)
	// 欢乐宝箱状态
	responseMessage := &proto.S2CActivitStatus{}
	activit := proto.Activity{}
	activit.TreasureBox = G_BaseCfg.TreasureBox.StateSwitch
	//新春活动 (排行榜)
	activityState := G_DoubleYearEvent.getActivityState()
	rankList := double_year.NewRankList()
	activit.RandList = rankList.GetRandListStatus()
	activit.ExchangeRole = activityState
	activit.ExchangeCdt = activityState

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	responseMessage.Activity = activit
	rsp.WriteData(responseMessage)
	if isAll {
		Sched.SendToAllUser(rsp)
		return
	} else if userId != "" {
		Sched.SendToUser(userId, rsp)
	}
}
