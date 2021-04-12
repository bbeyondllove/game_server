package message

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/game/errcode"
	"game_server/game/proto"

	"game_server/core/logger"
)

//获取验证码
func (a *agent) HandleGetVerificationCode(requestMsg *utils.Packet) {
	responseMessage := Handle_GetVerificationCode(requestMsg)
	SendPacket(a.conn, responseMessage)
}

func Handle_GetVerificationCode(requestMsg *utils.Packet) *utils.Packet {
	logger.Debugf("Handle_GetVerificationCode in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_VERIFICATION_CODE_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetVerificationCode{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		return rsp
	}

	if msg.UseFor != proto.USE_REGISTR && msg.UseFor != proto.USE_RESET_PASSWD && msg.UseFor != proto.USE_BIND {
		logger.Errorf("UseFor error:", msg.UseFor)
		rsp.WriteData(responseMessage)
		return rsp
	}

	if msg.CodeType != proto.CODE_TYPE_PHONE && msg.CodeType != proto.CODE_TYPE_EMAIL {
		logger.Errorf("CodeType error:", msg.CodeType)
		rsp.WriteData(responseMessage)
		return rsp
	}

	Url := ""
	if msg.Language == "" {
		msg.Language = "en"
	}
	request_map := make(map[string]interface{}, 0)

	if msg.CodeType == proto.CODE_TYPE_EMAIL {
		if msg.Email == "" {
			logger.Errorf("Email == ''")
			rsp.WriteData(responseMessage)
			return rsp
		}

		request_map["email"] = msg.Email
		Url = base.Setting.Base.EmailCodeUrl
	} else {
		if msg.CountryCode == 0 || msg.Mobile == "" {
			logger.Errorf("CountryCode == ", msg.CountryCode, ",Mobile == ", msg.Mobile)
			rsp.WriteData(responseMessage)
			return rsp
		}
		request_map["mobile"] = msg.Mobile
		request_map["countryCode"] = msg.CountryCode
		Url = base.Setting.Base.MobileCodeUrl
	}

	request_map["useFor"] = msg.UseFor
	request_map["language"] = msg.Language
	request_map["sysType"] = msg.SysType
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)

	buf, _ := json.Marshal(request_map)
	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + Url
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
		logger.Errorf("HttpPost() error, err=", err.Error())
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		return rsp
	}

	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		return rsp
	}

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleGetVerificationCode end")
	return rsp
}
