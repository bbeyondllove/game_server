package message

import (
	"encoding/json"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/errcode"
	"game_server/game/logic"
	"game_server/game/message/activity/double_year"
	"game_server/game/message/statistical"
	"game_server/game/model"
	"game_server/game/proto"
)

// Doubleyear 双旦活动struct.
type DoubleYear struct {
}

// TradeDayCdt 道具对换cdt.
func (a *agent) TradeDayCdt(requestMsg *utils.Packet) {
	logger.Infof("TradeDayCdt in request:", requestMsg.GetBuffer())
	tradeCdt := double_year.TradeCdt{B: double_year.NewBase(), Cr: model.NewCdtRecord()}
	args, err := tradeCdt.B.ParseArgsToMap(requestMsg)
	if err != nil {
		tradeCdt.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		tradeCdt.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := tradeCdt.B.GenerateResponseMessage(proto.MSG_SWEET_TREE_RSP, tradeCdt.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	data, code := tradeCdt.TradeCdt(args["userId"].(string), int(args["sweet"].(float64)), int(args["tree"].(float64)))
	tradeCdt.B.ResMsg.Code = code
	tradeCdt.B.ResMsg.Msg = double_year.StatusCodeMessage[code]
	tradeCdt.B.ResMsg.Data = data
	rsp := tradeCdt.B.GenerateResponseMessage(proto.MSG_SWEET_TREE_RSP, tradeCdt.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// RankListUpdateProp 更新道具分数值.
func (a *agent) RankListUpdateProp(requestMsg *utils.Packet) {
	logger.Infof("RankListDay in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_UPDATE_PROP_RSP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}
	rankList.UpdateProp(args["userId"].(string), int(args["propNumber"].(float64)), int(args["count"].(float64)))
}

// RankListDay 每日排行榜.
func (a *agent) RankListDay(requestMsg *utils.Packet) {
	logger.Infof("RankListDay in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_DAY_RSP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	rankList.B.ResMsg.Code = double_year.ActiveDoubleYearSuccess
	rankList.B.ResMsg.Data = rankList.GetDayRankList(args["userId"].(string))
	rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearSuccess]
	rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_DAY_RSP, rankList.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// RankListAll 总排行榜.
func (a *agent) RankListAll(requestMsg *utils.Packet) {
	logger.Infof("RankListAll in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_ALL_RSP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	rankList.B.ResMsg.Code = double_year.ActiveDoubleYearSuccess
	rankList.B.ResMsg.Data = rankList.GetAllRankList(args["userId"].(string))
	rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearSuccess]
	rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_ALL_RSP, rankList.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// RankListDayProp 每日双旦值记录.
func (a *agent) RankListDayProp(requestMsg *utils.Packet) {
	logger.Infof("RankListDayProp in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_DAY_PROP_RSP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	rankList.B.ResMsg.Code = double_year.ActiveDoubleYearSuccess
	rankList.B.ResMsg.Data = rankList.GetDayProp(args["userId"].(string))
	rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearSuccess]
	rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_DAY_PROP_RSP, rankList.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// RankListAllProp 总双旦值记录.
func (a *agent) RankListAllProp(requestMsg *utils.Packet) {
	logger.Infof("RankListAllProp in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_ALL_PROP_RSP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	rankList.B.ResMsg.Code = double_year.ActiveDoubleYearSuccess
	rankList.B.ResMsg.Data = rankList.GetAllProp(args["userId"].(string))
	rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearSuccess]
	rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_ALL_PROP_RSP, rankList.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// RankListInviteRecord 邀请好友记录.
func (a *agent) RankListInviteRecord(requestMsg *utils.Packet) {
	logger.Infof("RankListInviteRecord in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_INVITE_RECORD_RSP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	rankList.B.ResMsg.Code = double_year.ActiveDoubleYearSuccess
	//rankList.B.ResMsg.Data = rankList.InviteRecord(args["userId"].(string))
	var userInvite logic.UserInvitation
	data, _ := userInvite.GetUserInvitList(args["userId"].(string), int(args["isDay"].(float64)))
	rankList.B.ResMsg.Data = map[string]interface{}{
		"list": data,
	}
	rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearSuccess]
	rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_INVITE_RECORD_RSP, rankList.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// RankListRepairCdt 补发排行榜奖励.
func (a *agent) RankListRepairCdt(requestMsg *utils.Packet) {
	logger.Infof("RankListRepairCdt in request:", requestMsg.GetBuffer())
	rankList := double_year.NewRankList()
	args, err := rankList.B.ParseArgsToMap(requestMsg)
	if err != nil {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_EMAIL_REPAIR_CDT_RESP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	// 设置一个临时密码.
	if secret, ok := args["secret"]; !ok || secret.(string) != "doYouWantToKnow" {
		rankList.B.ResMsg.Code = double_year.ActiveDoubleYearArgsError
		rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearArgsError]
		rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_EMAIL_REPAIR_CDT_RESP, rankList.B.ResMsg)
		SendPacket(a.conn, rsp)
		return
	}

	status := rankList.RepairAwardCdt(args["receiverUserId"].(string), args["emailTitle"].(string), args["emailContent"].(string), int(args["eventType"].(float64)), int(args["cdtValue"].(float64)))
	rankList.B.ResMsg.Code = double_year.ActiveDoubleYearSuccess
	rankList.B.ResMsg.Data = map[string]interface{}{"repair": status}
	rankList.B.ResMsg.Msg = double_year.StatusCodeMessage[double_year.ActiveDoubleYearSuccess]
	rsp := rankList.B.GenerateResponseMessage(proto.MSG_RANK_LIST_EMAIL_REPAIR_CDT_RESP, rankList.B.ResMsg)
	SendPacket(a.conn, rsp)
}

// 排行榜打点
func (a *agent) RankingDot(requestMsg *utils.Packet) {
	logger.Infof("RankingDot in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_DOUBLE_YEAR_DOT_RSP)

	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.Dot{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	flag, PayLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	if msg.CodeType == proto.DOT_DOUBLE_YEAR_FRAGMENT {
		statistical.StatisticsDotIns.DoubleYearFuwaPVUV(PayLoad.UserId) // 福娃打点
	} else if msg.CodeType == proto.DOT_DOUBLE_YEAR_DAY {
		statistical.StatisticsDotIns.DoubleYearDayPVUV(PayLoad.UserId) // 新春活动每日排行
	} else if msg.CodeType == proto.DOT_DOUBLE_YEAR_TOTAL {
		statistical.StatisticsDotIns.DoubleYearRankingListPVUV(PayLoad.UserId) // 新春活动总排行
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	SendPacket(a.conn, rsp)
}
