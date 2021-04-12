package message

import (
	"encoding/json"
	kk_core "game_server/core"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/shopspring/decimal"
)

//获取用户金额
func (s *CSession) HandleGetAmount(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetAmount in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_AMOUNT_RSP)

	responseMessage := &proto.S2CGetAmount{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetAmount{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.UserId = payLoad.UserId
	var money map[string]map[string]float64
	err = json.Unmarshal([]byte(user_info["money"]), &money)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if value, ok := money[msg.TokenCode]; ok {
		responseMessage.Amount = value["amount"]
		responseMessage.AmountAvailable = value["amount_available"]
		responseMessage.AmountBlocked = value["amount_blocked"]
	} else {
		responseMessage.Amount = 0
		responseMessage.AmountAvailable = 0
		responseMessage.AmountBlocked = 0
	}

	cdt := user_info["cdt"]
	responseMessage.Cdt, _ = decimal.NewFromString(cdt)
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	return

	logger.Errorf("user not found wallect %+v", payLoad.UserId)
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	logger.Debugf("HandleGetAmount end")
	return
}

//背包和金额变更
func (s *CSession) HandleUpdateItemInfo(requestMsg *utils.Packet) {
	logger.Debugf("HandleUpdateItemInfo in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_UPDATE_ITEM_INFO_RSP)

	responseMessage := &proto.S2CUpdateItemInfo{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SUpdateItemInfo{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, userInfo := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	responseMessage.Cdt, _ = decimal.NewFromString(userInfo["cdt"])
	_, responseMessage.ItemInfos = GetUserItem(userInfo, "item_info")

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleUpdateItemInfo end")
	return
}

//用户背包
func (s *CSession) HandleUserKnapsack(requestMsg *utils.Packet) {
	logger.Debugf("HandleUserKnapsack in request:", requestMsg.GetBuffer())

	responseMessage := &proto.S2CGetKnapsack{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_KNAPSACK_RSP)

	msg := &proto.C2SGetKnapsack{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	responseMessage.UserId = payLoad.UserId
	_, responseMessage.ItemInfos = GetUserItem(user_info, "item_info")
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	logger.Debugf("HandleUserKnapsack end")
	return
}

//获取战友列表
func (s *CSession) HandleGetInviteUsers(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetInviteUsers in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_INVITE_USERS_RSP)

	responseMessage := &proto.S2CInviteUserList{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SInviteUserList{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if msg.InvitationLevel < 0 || msg.InvitationLevel > 3 {
		logger.Errorf("requestMsg[level] is not legal ,level=", msg.InvitationLevel)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
	}

	request_map := make(map[string]interface{}, 0)
	request_map["sysType"] = msg.SysType
	request_map["inviterId"] = msg.InviterId
	request_map["invitationLevel"] = msg.InvitationLevel
	request_map["page"] = msg.Page
	request_map["vipLevel"] = msg.VipLevel
	request_map["size"] = msg.Size
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetInviteUsersUrl
	msgdata, err := utils.HttpGet(url, request_map)

	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleGetInviteUsers end")
	return
}

//获取被抢走战友列表
func (s *CSession) HandleGetGrabComrades(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetGrabComrades in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_INVITE_USERS_RSP)
	responseMessage := &proto.S2CGetGrabComrades{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetGrabComrades{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetGrabComradesUrl
	msgdata, err := utils.HttpPost(url, string(requestMsg.Bytes()), proto.JSON)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
	}

	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	logger.Debugf("HandleGetGrabComrades end")
	return
}

//获取会员等级体系数据
func (s *CSession) HandleGetMemberSys(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetMemberSys in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_MEMBER_SYS_RSP)
	responseMessage := &proto.S2CGetMemberSys{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetMemberSys{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Message = ""
	reqRequest := make(map[string]interface{}, 0)
	reqRequest["sysType"] = msg.SysType
	reqRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	reqRequest["sign"] = utils.GenTonken(reqRequest)
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetMemberSysUrl
	msgdata, err := utils.HttpGet(url, reqRequest)
	if err != nil {
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
	}
	logger.Debugf("msgdata=", string(msgdata))
	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if responseMessage.Code == errcode.MSG_SUCCESS {
		kk_core.PushMysql(func() {
			userLevelConfigs, err := db_service.UserLevelConfigIns.GetUserLevelList()
			if err != nil {
				logger.Errorf("HandleGetMemberSys db.GetUserLevelList() failed(), err=", err.Error())
				kk_core.PushWorld(func() {
					responseMessage.Code = errcode.ERROR_MYSQL
					responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
					rsp.WriteData(responseMessage)
					s.sendPacket(rsp)
				})
				return
			}
			// msgdata := utils.ToSlice(msgData)
			dataMap := make(map[int][]proto.UserLevelItem, 0)
			for k, v := range userLevelConfigs {
				var node proto.UserLevelItem
				node.ItemId = v.ItemId
				node.ItemName = v.ItemName
				node.ItemNum = v.ItemNum
				node.ItemType = v.ItemType

				if _, ok := dataMap[v.LevelId]; ok {
					dataMap[v.LevelId] = append(dataMap[k], node)
				} else {
					dataMap[v.LevelId] = make([]proto.UserLevelItem, 0)
					dataMap[v.LevelId] = append(dataMap[k], node)
				}
			}
			for k, member := range responseMessage.Data {
				responseMessage.Data[k].ItemList = make([]proto.UserLevelItem, 0)
				if val, ok := dataMap[member.MemberLevel]; ok {
					if len(val) > 0 {
						responseMessage.Data[k].ItemList = append(responseMessage.Data[k].ItemList, val...)
					}
				}
			}
			kk_core.PushWorld(func() {
				rsp.WriteData(responseMessage)
				logger.Debugf(string(rsp.Bytes()))
				s.sendPacket(rsp)
				logger.Debugf("HandleGetMemberSys end")
			})
		})
	}

}

//获取用户当前等级状态
func (s *CSession) HandleGetUserLevel(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetUserLevel in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_USER_LEVEL_RSP)
	responseMessage := &proto.S2CGetUserLevel{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetUserLevel{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Message = ""
	reqRequest := make(map[string]interface{}, 0)
	reqRequest["userId"] = payLoad.UserId
	reqRequest["sysType"] = msg.SysType
	reqRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	reqRequest["sign"] = utils.GenTonken(reqRequest)
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetUserLevelUrl
	msgdata, err := utils.HttpGet(url, reqRequest)
	if err != nil {
		logger.Errorf("utils.HttpPost error, err=", err.Error())
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	if responseMessage.Code != 0 {
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		logger.Errorf(string(rsp.Bytes()))
		return
	}

	paramData := []string{"available_roles", "deblocked_roles", "level", "top_level", "user_type"}
	userInfo, err := db.RedisGame.HMGet(payLoad.UserId, paramData...).Result()
	if err != nil {
		logger.Errorf("HandleGetUserLevel HMGet(userid, 'paramData') failed err=", err.Error())
		responseMessage.Code = errcode.ERROR_REDIS
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_REDIS]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	aRoles := userInfo[0].(string)
	dRoles := userInfo[1].(string)
	strlevel := userInfo[2].(string)
	topLevelStr := userInfo[3].(string)
	userTypeStr := userInfo[4].(string)
	level, _ := strconv.Atoi(strlevel)
	topLevel, _ := strconv.Atoi(topLevelStr)
	userType, _ := strconv.Atoi(userTypeStr)
	logger.Debugf("level=%+v,responseMessage.Member.MemberLevel=%+v,topLevel=%+v,userType=%+v", level, responseMessage.Member.MemberLevel, topLevel, userType)
	if responseMessage.Member.MemberLevel > topLevel {
		s.UpdateLevel(topLevel, responseMessage.Member.MemberLevel)
	}
	availableRoles, deblockedRoles := GetRoles(payLoad.UserId, userType, responseMessage.Member.MemberLevel, aRoles, dRoles)
	if availableRoles != aRoles || deblockedRoles != dRoles {
		kk_core.PushMysql(func() {
			//更新available_roles、amount_blocked字段值
			rolesMap := make(map[string]interface{}, 0)
			rolesMap["available_roles"] = availableRoles
			_, err = db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, rolesMap)
			if err != nil {
				logger.Errorf("UpdateFields %+v failed err=%+v", db_service.UserTable, err.Error())
				kk_core.PushWorld(func() {
					responseMessage.Code = errcode.ERROR_MYSQL
					responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
					rsp.WriteData(responseMessage)
					s.sendPacket(rsp)
				})
				return
			}

			value, err := db.RedisGame.HMSet(payLoad.UserId, rolesMap).Result()
			if err != nil {
				logger.Errorf("RedisGame.HMSet failed uid=,err=", payLoad.UserId, err.Error())
				kk_core.PushWorld(func() {
					responseMessage.Code = errcode.ERROR_REDIS
					responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_REDIS]
					rsp.WriteData(responseMessage)
					s.sendPacket(rsp)
				})
				return
			}
			if value == "OK" {
				logger.Debugf("RedisGame.HMSet success, uid=", payLoad.UserId)
			}

			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			s.sendPacket(rsp)
		})

	} else {
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		s.sendPacket(rsp)
	}

	logger.Debugf("HandleGetUserLevel end")
	return
}

//用户升级时，赠送道具卡
func (s *CSession) UpdateLevel(level, newLevel int) {
	logger.Debugf("UpdateLevel in:", level, newLevel)
	var levelConfigs []model.UserLevelConfig
	var err error

	levelConfigs, err = db_service.UserLevelConfigIns.GetLevelItemByLevel(level, newLevel)
	if err != nil {
		logger.Errorf("UserLevelConfigIns.GetLevelItemByLevel failed err=", err.Error())
		return
	}

	if len(levelConfigs) == 0 {
		logger.Info("UserLevelConfigIns.GetLevelItemByLevel return 0 record")
		return
	}

	//token 检验
	user_info := db.RedisMgr.HGetAll(s.UserId)
	if user_info == nil {
		logger.Errorf("RedisGame.HMGet failed userId=,err=", s.UserId, err.Error())
		return
	}
	userItem := user_info["item_info"]
	itemMap := make(map[int]map[int]int, 0)
	err = json.Unmarshal([]byte(userItem), &itemMap)
	if err != nil {
		logger.Errorf("UpdateLevel json.Unmarshal failed err=", err.Error())
		return
	}

	logger.Debugf("UpdateLevel before itemMap", itemMap, levelConfigs)
	for _, config := range levelConfigs {
		bFlag := false
		if _, ok := itemMap[config.ItemType]; ok {
			if _, subok := itemMap[config.ItemType][config.ItemId]; subok {
				itemMap[config.ItemType][config.ItemId] = itemMap[config.ItemType][config.ItemId] + config.ItemNum
				bFlag = true
			} else {
				itemMap[config.ItemType][config.ItemId] = config.ItemNum
			}
		} else {
			itemMap[config.ItemType] = make(map[int]int, 0)
			itemMap[config.ItemType][config.ItemId] = 1
		}

		dataMap := make(map[string]interface{})
		dataMap["update_time"] = time.Now()
		dataMap["item_num"] = itemMap[config.ItemType][config.ItemId]

		if bFlag {
			logger.Debugf("UpdateLevel UpdateData UserKnapsackIns", s.UserId, config.ItemId, dataMap)
			_, err = db_service.UserKnapsackIns.UpdateData(s.UserId, config.ItemId, dataMap)
		} else {
			data := model.UserKnapsack{
				UserId:   s.UserId,
				ItemType: config.ItemType,
				ItemId:   config.ItemId,
				ItemNum:  config.ItemNum,
			}
			logger.Debugf("UpdateLevel Add UserKnapsackIns", s.UserId, config.ItemId, dataMap)
			_, err = db_service.UserKnapsackIns.Add(&data)
		}

	}
	if err != nil {
		logger.Errorf("UpdateLevel  update UserKnapsack failed err=", err.Error())
		return
	}

	logger.Debugf("UpdateLevel after itemMap", itemMap)
	userItemInfo, _ := json.Marshal(itemMap)
	userInfoMap := make(map[string]interface{})
	userInfoMap["update_time"] = time.Now()
	userInfoMap["level"] = newLevel
	userInfoMap["top_level"] = newLevel
	userInfoMap["item_info"] = userItemInfo

	_, err = db.RedisGame.HMSet(s.UserId, userInfoMap).Result()
	if err != nil {
		logger.Errorf("UpdateLevel RedisGame.HMSet failed err=", err.Error())
	}

	dataMap := make(map[string]interface{})
	dataMap["update_time"] = userInfoMap["update_time"]
	dataMap["level"] = userInfoMap["level"]
	dataMap["top_level"] = userInfoMap["top_level"]
	db_service.UpdateFields(db_service.UserTable, "user_id", s.UserId, dataMap)
	logger.Debugf("UpdateLevel end")

}

//抢战友
func (s *CSession) HandleGrapComrade(requestMsg *utils.Packet) {
	logger.Debugf("HandleGrapComrade in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GRAB_COMRADE_RES)

	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGrabComrade{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetUserLevelUrl
	msgdata, err := utils.HttpPost(url, string(requestMsg.Bytes()), proto.JSON)
	if err != nil {
		logger.Errorf("utils.HttpPost, err=", err.Error())
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
	}

	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleGrapComrade end")
	return
}

//同步用户系统
func SyncUserNickName(userId string, nickName string) error {
	logger.Debugf("syncUserNickName in request:", nickName, userId)

	dataMap := make(map[string]interface{}, 0)
	dataMap["nickName"] = nickName
	dataMap["userId"] = userId
	dataMap["sysType"] = proto.SysType
	dataMap["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	dataMap["sign"] = utils.GenTonken(dataMap)
	json_data, err := json.Marshal(dataMap)
	if err != nil {
		logger.Errorf("syncUserNickName json.Marshal, err=", err.Error())
		return err
	}

	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.UpdateUserUrl
	utils.HttpPost(url, string(json_data), proto.JSON)
	return nil
}

//修改昵称
func (s *CSession) HandleModifyNickName(requestMsg *utils.Packet) {
	logger.Debugf("HandleModifyNickNamein request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_MODIFY_NICKNAME_RSP)
	responseMessage := &proto.S2CModifyNickName{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SModifyNickName{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	isExists := db.RedisMgr.KeyExist(msg.NickName)
	if isExists {
		responseMessage.Code = errcode.ERROR_ROLE_EXITS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if msg.ModifyType == 0 {
		if user_info["modify_name_num"] == "0" {
			responseMessage.Code = errcode.ERROR_NO_FREE_MODIFY_COUNT
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}

		data_map := make(map[string]interface{}, 0)
		data_map["modify_name_num"] = 0
		data_map["nick_name"] = msg.NickName
		data_map["update_time"] = time.Now()
		_, err := db.RedisGame.HMSet(payLoad.UserId, data_map).Result()
		if err != nil {
			responseMessage.Code = errcode.ERROR_SYSTEM
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}

		_, err = db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, data_map)
		if err != nil {
			responseMessage.Code = errcode.ERROR_SYSTEM
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)

			data_map = make(map[string]interface{})
			data_map["modify_name_num"] = 1
			data_map["nick_name"] = user_info["nickName"]
			data_map["update_time"] = time.Now()
			_, err = db.RedisGame.HMSet(payLoad.UserId, data_map).Result()
			return
		}

		db.RedisMgr.SetCode(msg.NickName, "1")
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = ""
		responseMessage.NickName = msg.NickName
		responseMessage.ModifyNameNum = 0
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		s.sendPacket(rsp)
		SyncUserNickName(payLoad.UserId, msg.NickName)
	} else {
		bSuccess, code := UseItem(payLoad.UserId, user_info, msg.ItemId)
		if !bSuccess {
			responseMessage.Code = code
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			s.sendPacket(rsp)
			return
		}

		data_map := make(map[string]interface{}, 0)
		data_map["modify_name_num"] = 0
		data_map["nick_name"] = msg.NickName
		data_map["update_time"] = time.Now()
		_, err := db.RedisGame.HMSet(payLoad.UserId, data_map).Result()
		if err != nil {
			responseMessage.Code = errcode.ERROR_SYSTEM
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}
		_, err = db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, data_map)

		modifyNameNum, _ := strconv.Atoi(user_info["modify_name_num"])
		db.RedisMgr.SetCode(msg.NickName, "1")
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = ""
		responseMessage.NickName = msg.NickName
		responseMessage.ModifyNameNum = modifyNameNum
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		s.sendPacket(rsp)
		SyncUserNickName(payLoad.UserId, msg.NickName)
	}

	logger.Debugf("syncUserNickName end")
	return
}

//实名认证
func (s *CSession) HandleCertification(requestMsg *utils.Packet) {
	logger.Debugf("HandleCertification in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_CERTIFICATION_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SCertification{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	var photo proto.CertificationPhoto
	err = json.Unmarshal([]byte(msg.ObjectKey), &photo)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if photo.Front == "" || photo.Back == "" || photo.Other == "" {
		logger.Errorf("photo is null %+v", photo)
		responseMessage.Code = errcode.ERROR_PARAM_EMPTY
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	dbMsg := ""
	data, suberr := db_service.CertificationIns.GetDataByUid(payLoad.UserId)
	if suberr != nil {
		logger.Errorf("HandleCertification CertificationIns %+v failed(), err=%+v", dbMsg, err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(s.conn, rsp)
		return
	}

	if data.UserId == "" {
		data := model.Certification{
			UserId:      payLoad.UserId,
			Nationality: msg.Nationality,
			FirstName:   msg.FirstName,
			LastName:    msg.LastName,
			IdType:      msg.IdType,
			IdNumber:    msg.IdNumber,
			ObjectKey:   msg.ObjectKey,
			Suggestion:  0,
			Status:      0,
		}
		_, err = db_service.CertificationIns.Add(&data)
		dbMsg = "add"
	} else {
		data_map := make(map[string]interface{})
		data_map["nationality"] = msg.Nationality
		data_map["first_name"] = msg.FirstName
		data_map["last_name"] = msg.LastName
		data_map["id_type"] = msg.IdType
		data_map["id_number"] = msg.IdNumber
		data_map["object_key"] = msg.ObjectKey
		data_map["update_time"] = time.Now()
		data_map["status"] = 3
		_, err = db_service.UpdateFields(db_service.CertificationTable, "user_id", payLoad.UserId, data_map)
		dbMsg = "update"
	}
	if err != nil {
		logger.Errorf("HandleCertification CertificationIns %+v failed(), err=%+v", dbMsg, err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(s.conn, rsp)
		return
	}

	data_map := make(map[string]interface{})
	data_map["kyc_status"] = 0
	data_map["update_time"] = time.Now()
	db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, data_map)

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleCertification end")
	return
}

//绑定邀请码
func (s *CSession) HandleBindInvitationCode(requestMsg *utils.Packet) {
	logger.Debugf("HandleBindInvitationCode in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_BIND_INVITER_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SBindInviter{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	reqRequest := make(map[string]interface{}, 0)
	reqRequest["userId"] = payLoad.UserId
	reqRequest["inviteCode"] = msg.InviteCode
	reqRequest["sysType"] = proto.SysType
	reqRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	reqRequest["sign"] = utils.GenTonken(reqRequest)

	json_data, err := json.Marshal(reqRequest)
	if err != nil {
		logger.Errorf("json.Marshal, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.BindingInviterUrl
	msgdata, err := utils.HttpPost(url, string(json_data), proto.JSON)
	if err != nil {
		logger.Errorf("utils.HttpPost error, err=", err.Error())
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
	}
	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)

	logger.Debugf("HandleBindInvitationCode end")
}

//好友邀请收益
func (s *CSession) HandleDepositRebate(requestMsg *utils.Packet) {
	logger.Debugf("HandleDepositRebate in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_DEPOSIT_REBATE_RSP)
	responseMessage := &proto.S2CDepositRebate{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SDepositRebate{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	logger.Debugf("in request:", msg)

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	reqRequest := make(map[string]interface{}, 0)
	reqRequest["userId"] = payLoad.UserId
	reqRequest["page"] = msg.Page
	reqRequest["size"] = msg.Size
	reqRequest["sysType"] = proto.SysType
	reqRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	reqRequest["sign"] = utils.GenTonken(reqRequest)
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.InviteIncomelistUrl
	msgdata, err := utils.HttpGet(url, reqRequest)
	if err != nil {
		logger.Errorf("utils.HttpPost error, err=", err.Error())
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
	}
	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)

	logger.Debugf("HandleDepositRebate end")
}
