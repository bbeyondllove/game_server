package message

import (
	"encoding/json"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"github.com/shopspring/decimal"
)

/**
  获取邮件列表
*/
func (a *agent) HandleGetEmailList(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetEmailList in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_EMAIL_LIST_RSP)
	responseMessage := &proto.S2CEmailList{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.IsPush = 0
	msg := &proto.C2SEmailList{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	emailList, err := db_service.EmailIns.GetDataByUserId(payLoad.UserId, msg.IsRead)
	if err != nil {
		logger.Errorf("GetDataByUserId error", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	responseMessage.ReadList = make([]proto.EmailItem, 0)
	responseMessage.UnreadList = make([]proto.EmailItem, 0)
	// 区分已读未读数据
	if emailList != nil {
		for _, email := range emailList {
			//带奖励的邮件
			if email.EmailType == 2 {
				ret, _ := db_service.EmailLogicIns.GetEmailPrize(email.Id)
				email.PrizeList = ret
			} else {
				email.PrizeList = make([]model.EmailPrize, 0)
			}
			//区分已读未读
			if email.IsRead == 0 {
				responseMessage.UnreadList = append(responseMessage.UnreadList, email)
			} else {
				responseMessage.ReadList = append(responseMessage.ReadList, email)
			}
		}
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleGetEmailList end")
	SendPacket(a.conn, rsp)
	return
}

/**
获取邮件数量
*/
func (a *agent) HandleGetEmailCount(requestMsg *utils.Packet) {
	logger.Debugf("HandleDelEmails in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_COUNT_EMAIL_RSP)
	responseMessage := &proto.S2CSetEmailRead{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	msg := &proto.C2SGetEmialCount{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	var emailMode db_service.EmailInfo
	var unreadNum, readNum int64
	// 未读的邮件
	unreadNum, err = emailMode.GetEmailCount(payLoad.UserId, map[string]interface{}{"is_read": 0})
	if err != nil {
		logger.Errorf("DelEmailInIds error", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	// 未读的邮件
	readNum, err = emailMode.GetEmailCount(payLoad.UserId, map[string]interface{}{"is_read": 1})
	if err != nil {
		logger.Errorf("DelEmailInIds error", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	responseMessage.UnreadNum = unreadNum
	responseMessage.ReadNum = readNum
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleDelEmails end")
	SendPacket(a.conn, rsp)
}

/**
删除邮件
*/

func (a *agent) HandleDelEmails(requestMsg *utils.Packet) {
	logger.Debugf("HandleDelEmails in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_DEL_EMAIL_RSP)
	responseMessage := &proto.S2CDelEmail{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	msg := &proto.C2SDelEmail{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	var emailMode db_service.EmailInfo
	_, err = emailMode.DelEmailInIds(payLoad.UserId, msg.EmailIds)
	if err != nil {
		logger.Errorf("DelEmailInIds error", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleDelEmails end")
	SendPacket(a.conn, rsp)
}

/**
  设置邮件为已读
*/

func (a *agent) HandleSetEmailRead(requestMsg *utils.Packet) {
	logger.Debugf("HandleSetEmailRead in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SET_EMAIL_READ_RSP)
	responseMessage := &proto.S2CSetEmailRead{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	msg := &proto.C2SSetEmailRead{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	var emailMode db_service.EmailInfo
	_, err = emailMode.SetEmailIsRead(payLoad.UserId, msg.EmailIds)
	if err != nil {
		logger.Errorf("SetEmailIsRead error", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleSetEmailRead end")
	SendPacket(a.conn, rsp)
}

/**
领取邮件奖励
*/
func (a *agent) EmailReceiveRewards(requestMsg *utils.Packet) {
	logger.Debugf("EmailReceiveRewards in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_EMAIL_RECEIVE_REWARDS_RSP)
	responseMessage := &proto.EmailBase{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	msg := &proto.C2SSetReceiveRewards{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("EmailReceiveRewards json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	code, email, isSuccess := db_service.EmailLogicIns.ReceiveRewards(payLoad.UserId, msg.EmailId)
	if isSuccess != true {
		responseMessage.Code = code
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		logger.Debugf("EmailReceiveRewards end")
		SendPacket(a.conn, rsp)
		return
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("EmailReceiveRewards end")
	SendPacket(a.conn, rsp)
	// 通知前端更新 CDT
	SendCdtMsg(payLoad.UserId, email)
}

func SendCdtMsg(userId string, email model.Email) {
	// 活动 CDT 需要推送消息
	if email.EmailType == 2 {
		prizeList, err := db_service.EmailPrizeIns.GetEmailPrize(email.Id)
		if err == nil && len(prizeList) >= 1 {
			// cdt 需要推送消息
			if prizeList[0].PrizeType == 2 {
				// cdt
				user_info := db.RedisMgr.HGetAll(userId)
				if user_info != nil {
					rsp := &utils.Packet{}
					rsp.Initialize(proto.MSG_UPDATE_ITEM_INFO_RSP)
					responseMessage := &proto.S2CUpdateItemInfo{}
					responseMessage.Code = errcode.MSG_SUCCESS
					responseMessage.Message = ""
					responseMessage.Cdt, _ = decimal.NewFromString(user_info["cdt"])
					_, responseMessage.ItemInfos = GetUserItem(user_info, "item_info")
					rsp.WriteData(responseMessage)
					logger.Debugf(string(rsp.Bytes()))
					Sched.SendToUser(userId, rsp)
				}
			}
		}
	}
}
