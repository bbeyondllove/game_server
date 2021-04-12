package message

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/game/errcode"
	"game_server/game/proto"

	"game_server/core/logger"
)

//获取商户连接
func (s *agent) HandleGetMerchantUrl(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetMerchantsUrl in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_MERCHANTS_URL_RSP)
	responseMessage := &proto.S2CGetMerchantUrl{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetMerchantUrl{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(s.conn, rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(s.conn, rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.MerchantEnteringUrl = base.Setting.Merchant.MerchantEnteringUrl
	responseMessage.ActivityPromotionUrl = base.Setting.Merchant.ActivityPromotionUrl
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(s.conn, rsp)
	logger.Debugf("HandleGetTokenPayUrl end")
}

// 获取配置URL
func (s *agent) GetConfigureUrl(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetMerchantsUrl in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_CONFIGURE_URL_RSP)
	responseMessage := &proto.S2CConfigureUrl{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SConfigureUrl{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(s.conn, rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(s.conn, rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.Data = base.Setting.ConfigureUrl
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(s.conn, rsp)
	logger.Debugf("HandleGetTokenPayUrl end")
}
