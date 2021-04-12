package logic

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/proto"
)

type BaseLogic struct {
}

func (this *BaseLogic) UserLogin(account string, password string, isEmail bool, countryCode int) ([]byte, error) {
	request_map := make(map[string]interface{}, 0)
	request_map["sysType"] = proto.SysType
	request_map["password"] = password
	request_map["clientIp"] = "127.0.0.1"
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	url := ""
	if isEmail {
		url = base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.EmailLoginrUrl
		request_map["email"] = account

	} else {
		url = base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.MobileLoginUrl
		request_map["countryCode"] = countryCode
		request_map["mobile"] = account
	}
	request_map["sign"] = utils.GenTonken(request_map)
	buf, _ := json.Marshal(request_map)
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
		return make([]byte, 0), err
	}
	return msgdata, nil
}

/**
去Base系统获取用户邀请关系
*/
func (this *BaseLogic) BaseInvitation(userId string) ([]byte, error) {
	userinfoRequest := make(map[string]interface{}, 0)
	userinfoRequest["userId"] = userId
	userinfoRequest["sysType"] = proto.SysType
	userinfoRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	userinfoRequest["sign"] = utils.GenTonken(userinfoRequest)
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.UserInvitationUrl
	userdata, err := utils.HttpGet(url, userinfoRequest)
	if err != nil {
		return make([]byte, 0), err
	}
	return userdata, nil

}

/**
  去base系统绑定邀请关系
*/
func (this *BaseLogic) BaseBindingInviter(userId, inviteCode string) ([]byte, error) {
	reqRequest := make(map[string]interface{}, 0)
	reqRequest["userId"] = userId
	reqRequest["inviter"] = inviteCode
	reqRequest["sysType"] = proto.SysType
	reqRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	reqRequest["sign"] = utils.GenTonken(reqRequest)
	json_data, err := json.Marshal(reqRequest)
	if err != nil {
		logger.Errorf("BaseBindingInviter json.Marshal, err=", err.Error())
		return make([]byte, 0), err
	}
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.BindingInviterUrl
	msgdata, err := utils.HttpPost(url, string(json_data), proto.JSON)
	if err != nil {
		return make([]byte, 0), err
	}
	return msgdata, nil
}
