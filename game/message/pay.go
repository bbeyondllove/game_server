package message

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/game/errcode"
	"game_server/game/proto"

	"game_server/core/logger"
)

//获取聚合支付链接
func (s *CSession) HandleGetTokenPayUrl(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetTokenPayUrl in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_TOEKN_PAY_URL_RSP)
	responseMessage := &proto.S2cGetTokenPayUrl{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetTokenPayUrl{}
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
		s.sendPacket(rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.SharePayUrl = base.Setting.Pay.TokenPayUrl
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(s.conn, rsp)
	logger.Debugf("HandleGetTokenPayUrl end")
}
