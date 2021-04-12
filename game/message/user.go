package message

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/message/activity/double_year"
	"game_server/game/message/statistical"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type UserMessage struct {
}

func get_DevicePlatform(platform string) int {
	logger.Errorf("get_DevicePlatform:", platform)
	if len(platform) == 0 {
		return 0
	}
	if platform == proto.PLATFORM_ANDROID_STR {
		return proto.PLATFORM_ANDROID
	}
	if platform == proto.PLATFORM_IOS_STR {
		return proto.PLATFORM_IOS
	}
	return proto.PLATFORM_WEB
}

//用户心跳
func (a *agent) HandleHEARTBEAT(requestMsg *utils.Packet) {
	//logger.Debugf("HandleHEARTBEAT in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_HEARTBEAT_RSP)
	responseMessage := &proto.HeartBeatRsp{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW

	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())

		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	flag, payLoad, userInfo := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("HandleHEARTBEAT token error：%+V", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	//在线任务处理
	taskKey := proto.ONLINE_KEY + payLoad.UserId
	curSecond := time.Now().Unix()
	_, oldTime, _ := utils.SetKeyValue(taskKey, "lastTime", curSecond, false, utils.ITEM_DAY)
	if oldTime > 0 {
		timeLong := curSecond - oldTime
		//logger.Debugf("HandleHEARTBEAT oldTime:%+v,timeLong:%+v", oldTime, timeLong)
		go taskProcess(&a.conn, payLoad.UserId, proto.MSG_HEARTBEAT, 0, int(timeLong), true)

		// 双旦活动－累计在线时长.
		go func() {
			rankList := double_year.NewRankList()
			rankList.UpdateProp(userInfo["user_id"], double_year.PropOnlineTime, int(timeLong))
		}()
	}

	statistical.StatisticsDotIns.Heartbeat(payLoad.UserId)

	// 检查用户当前角色是否是圣诞老人， 过期后通知前端
	go func() {
		//金童玉女
		RemoveExpireColdenCoupl(payLoad.UserId)
	}()

	if a.auth {
		if time.Now().Unix() < payLoad.Exp {
			SendPacket(a.conn, rsp)
			return
		}
		/*
			//超时，刷新token
			request_map := make(map[string]interface{}, 0)
			request_map["token"] = a.token
			request_map["sysType"] = proto.SysType
			request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
			request_map["sign"] = utils.GenTonken(request_map)
			buf, err := json.Marshal(request_map)

			url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.RefreshTokenUrl
			msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
			if err == nil {
				resp := &proto.S2C_HTTP{}
				err = json.Unmarshal(msgdata, resp)

				if err == nil {
					a.token = resp.Data["token"].(string)
					responseMessage.Token = a.token
					rsp.WriteData(responseMessage)

					db.RedisGame.HSet(payLoad.UserId, "token", a.token)
				}
			}
		*/
		SendPacket(a.conn, rsp)
	} else {
		logger.Errorf("agent is not auth")
		a.conn.Close()
	}
	return
}

//手机注册
func (a *agent) HandleMobileRegister(requestMsg *utils.Packet) {
	responseMessage := Handle_MobileRegister(requestMsg)
	SendPacket(a.conn, responseMessage)
	return
}

func Handle_MobileRegister(requestMsg *utils.Packet) *utils.Packet {
	logger.Debugf("Handle_MobileRegister in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_REGISTER_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SRegisterMobile{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		return rsp
	}

	if msg.NickNname == "" {
		msg.NickNname = msg.Mobile
	}

	request_map := utils.StructToMap(*msg)
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)
	buf, err := json.Marshal(request_map)

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.MobileRegisterUrl
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
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

	if responseMessage.Code == errcode.MSG_SUCCESS {
		dataMap := make(map[string]interface{})
		dataMap["role_id"] = 0
		dataMap["nick_name"] = msg.Mobile
		if msg.NickNname != "" {
			dataMap["nick_name"] = msg.NickNname
		}
		dataMap["sex"] = 0
		dataMap["country_code"] = msg.CountryCode
		dataMap["mobile"] = msg.Mobile
		dataMap["email"] = ""
		dataMap["token"] = ""
		dataMap["status"] = 0
		dataMap["invite_code"] = msg.Inviter

		userdata := model.User{
			SysType:        msg.SysType,
			UserId:         "0",
			UserType:       0,
			RoleId:         0,
			AvailableRoles: G_BaseCfg.PingminARoles,
			DeblockedRoles: G_BaseCfg.PingminDRoles,
			NickName:       msg.NickNname,
			Sex:            0,
			Level:          0,
			CountryCode:    msg.CountryCode,
			Mobile:         msg.Mobile,
			Email:          "",
			Token:          "",
			Status:         1,
			LocationId:     0,
			PositionX:      G_BaseCfg.BirthPlace.X,
			PositionY:      G_BaseCfg.BirthPlace.Y,
			HouseNum:       0,
			ModifyNameNum:  1,
			InviteCode:     msg.Inviter,
			InviterId:      "0",
			KycPassed:      0,
			KycStatus:      -1,
			TopLevel:       0,
			Point:          0,
			Platform:       get_DevicePlatform(msg.Platform),
		}
		_, err := db_service.UserIns.Add(&userdata)
		if err != nil {
			logger.Errorf("Handle_Login UserIns.Add() failed(), err=", err.Error())
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			return rsp
		}
	}
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("Handle_MobileRegister end")
	return rsp
}

//邮箱注册
func (a *agent) HandleEmailRegister(requestMsg *utils.Packet) {
	responseMessage := Handle_EmailRegister(requestMsg)
	SendPacket(a.conn, responseMessage)
	return
}

func Handle_EmailRegister(requestMsg *utils.Packet) *utils.Packet {
	logger.Debugf("Handle_EmailRegister in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_REGISTER_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SRegisterEmail{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		return rsp
	}

	if msg.NickNname == "" {
		msg.NickNname = msg.Email
	}

	request_map := utils.StructToMap(*msg)
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)
	buf, err := json.Marshal(request_map)

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.EmailRegisterUrl
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
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

	if responseMessage.Code == errcode.MSG_SUCCESS {
		dataMap := make(map[string]interface{})
		dataMap["role_id"] = 0
		dataMap["nick_name"] = msg.Email
		if msg.NickNname != "" {
			dataMap["nick_name"] = msg.NickNname
		}
		dataMap["sex"] = 0
		dataMap["country_code"] = ""
		dataMap["mobile"] = ""
		dataMap["email"] = msg.Email
		dataMap["token"] = ""
		dataMap["status"] = 0
		dataMap["invite_code"] = msg.Inviter

		userdata := model.User{
			SysType:        msg.SysType,
			UserId:         "0",
			UserType:       0,
			RoleId:         0,
			AvailableRoles: G_BaseCfg.PingminARoles,
			DeblockedRoles: G_BaseCfg.PingminDRoles,
			NickName:       msg.NickNname,
			Sex:            0,
			Level:          0,
			CountryCode:    0,
			Mobile:         "",
			Email:          msg.Email,
			Token:          "",
			Status:         1,
			LocationId:     0,
			PositionX:      G_BaseCfg.BirthPlace.X,
			PositionY:      G_BaseCfg.BirthPlace.Y,
			HouseNum:       0,
			ModifyNameNum:  1,
			InviteCode:     msg.Inviter,
			InviterId:      "0",
			KycPassed:      0,
			KycStatus:      -1,
			TopLevel:       0,
			Point:          0,
			Platform:       get_DevicePlatform(msg.Platform),
		}
		_, err := db_service.UserIns.Add(&userdata)
		if err != nil {
			logger.Errorf("Handle_Login UserIns.Add() failed(), err=", err.Error())
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			return rsp
		}
	}

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("Handle_EmailRegister end")
	return rsp
}

//正常用户登录
func normalLogin(a *agent, rsp *utils.Packet, responseMessage *proto.S2CLogin, msg *proto.C2SLogin) {
	request_map := make(map[string]interface{}, 0)
	request_map["sysType"] = msg.SysType
	request_map["password"] = msg.Password
	request_map["clientIp"] = msg.ClientIp
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	account := msg.Mobile
	url := ""
	if msg.CodeType == proto.CODE_TYPE_EMAIL {
		account = msg.Email
		url = base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.EmailLoginrUrl
		request_map["email"] = msg.Email

	} else {
		url = base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.MobileLoginUrl
		request_map["countryCode"] = msg.CountryCode
		request_map["mobile"] = msg.Mobile
	}
	request_map["sign"] = utils.GenTonken(request_map)
	buf, _ := json.Marshal(request_map)

	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	resp := &proto.S2C_HTTP{}
	err = json.Unmarshal(msgdata, resp)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	if resp.Code != errcode.MSG_SUCCESS {
		responseMessage.Code = resp.Code
		responseMessage.Message = resp.Message
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		logger.Debugf("HandleLogin end")
		return
	}

	setLoginInfo(a, rsp, responseMessage, msg, account, false, resp)
	logger.Debugf("HandleLogin end")
	return
}

//机器人登录
func robotLogin(a *agent, rsp *utils.Packet, responseMessage *proto.S2CLogin, msg *proto.C2SLogin) {
	account := msg.Mobile
	if msg.CodeType == proto.CODE_TYPE_EMAIL {
		account = msg.Email
	}
	setLoginInfo(a, rsp, responseMessage, msg, account, true, nil)
	logger.Debugf("HandleLogin end")
}

//设置登录缓存
func setLoginInfo(a *agent, rsp *utils.Packet, responseMessage *proto.S2CLogin, msg *proto.C2SLogin, account string, bRobot bool, resp *proto.S2C_HTTP) {
	userid := ""
	token := ""
	userType := 0
	if bRobot {
		userType = proto.ROBOT_USER
	}

	userid = resp.Data["userId"].(string)
	//expiresAt:=resp.ExpiresAt
	userInfo, err := db_service.UserIns.GetDataByAccount(msg.CodeType, account, userid)
	if err != nil {
		logger.Errorf("Handle_Login UserIns.GetDataByAccount() failed(), err=", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	if G_BaseCfg.Backstage.UserStatusSwitch == 1 {
		if len(userInfo.Id) > 0 && userInfo.Status == 0 {
			logger.Errorf("Handle_Login UserIns. reject login", userInfo.Id, userInfo.UserId, userInfo.Status)
			responseMessage.Code = errcode.ERROR_HTTP_USER_NOT_ALLOW
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}
	}

	oldUserId := userInfo.UserId
	if !bRobot {
		userid = resp.Data["userId"].(string)
		token = resp.Data["token"].(string)

		userinfoRequest := make(map[string]interface{}, 0)
		userinfoRequest["userId"] = userid
		userinfoRequest["sysType"] = msg.SysType
		userinfoRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
		userinfoRequest["sign"] = utils.GenTonken(userinfoRequest)

		url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetUserinfoUrl
		userdata, err := utils.HttpGet(url, userinfoRequest)
		userResp := &proto.S2C_HTTP{}

		if err == nil {
			err = json.Unmarshal(userdata, userResp)
			if err == nil {
				userInfo.UserId = userid
				userInfo.Token = token
				userInfo.InviteCode = userResp.Data["inviteCode"].(string)
				userInfo.InviterId = userResp.Data["inviterId"].(string)
				userInfo.Level = int(userResp.Data["vipLevel"].(float64))
				userInfo.CountryCode = int(userResp.Data["countryCode"].(float64))
				userInfo.NickName = userResp.Data["nickName"].(string)
				userInfo.Mobile = userResp.Data["mobile"].(string)
				userInfo.Email = userResp.Data["email"].(string)
				bFlag := userResp.Data["kycPassed"].(bool)
				if bFlag {
					userInfo.KycPassed = 1
				}
				userInfo.NickName = userResp.Data["nickName"].(string)
				if userInfo.Level >= 1 {
					userInfo.UserType = 1
				}
			}
		}
	} else {
		if oldUserId == "" {
			logger.Errorf("Handle_Login UserIns.GetDataByAccount() failed(), err=", err.Error())
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}
		userid = userInfo.UserId
		token = userInfo.Token

	}

	userInfo.Platform = get_DevicePlatform(msg.Platform)
	userInfo.Version = msg.Version
	userInfo.LoginIp = msg.ClientIp
	suggestion := 0
	availableRoles, deblockedRoles := GetRoles(userid, userType, userInfo.Level, userInfo.AvailableRoles, userInfo.DeblockedRoles)
	if oldUserId == "" { //插入记录
		userInfo.Mobile = msg.Mobile
		userInfo.Email = msg.Email
		userInfo.RoleId = 0
		userInfo.AvailableRoles = availableRoles
		userInfo.DeblockedRoles = deblockedRoles
		userInfo.Sex = 0
		userInfo.LocationId = 0
		userInfo.PositionX = G_BaseCfg.BirthPlace.X
		userInfo.PositionY = G_BaseCfg.BirthPlace.Y
		userInfo.HouseNum = 0
		userInfo.ModifyNameNum = 1
		userInfo.KycStatus = -1
		userInfo.Status = 1
		userInfo.TopLevel = 0
		userInfo.SysType = msg.SysType
		userInfo.Point = 0
		dataMap := make(map[string]interface{})

		dataMap["user_id"] = userid
		dataMap["user_type"] = userInfo.UserType
		dataMap["role_id"] = userInfo.RoleId
		dataMap["nick_name"] = userInfo.NickName
		dataMap["sex"] = userInfo.Sex

		dataMap["level"] = userInfo.Level
		dataMap["country_code"] = userInfo.CountryCode
		dataMap["mobile"] = userInfo.Mobile
		dataMap["email"] = userInfo.Email
		dataMap["token"] = token

		dataMap["status"] = userInfo.Status
		dataMap["location_id"] = userInfo.LocationId
		dataMap["position_x"] = userInfo.PositionX
		dataMap["position_y"] = userInfo.PositionY
		dataMap["house_num"] = userInfo.HouseNum

		dataMap["modify_name_num"] = userInfo.ModifyNameNum
		dataMap["login_ip"] = userInfo.LoginIp
		dataMap["invite_code"] = userInfo.InviteCode
		dataMap["inviter_id"] = userInfo.InviterId
		dataMap["kyc_passed"] = userInfo.KycPassed

		dataMap["kyc_status"] = userInfo.KycStatus
		dataMap["top_level"] = userInfo.TopLevel
		dataMap["sys_type"] = userInfo.SysType
		dataMap["available_roles"] = availableRoles
		dataMap["deblocked_roles"] = deblockedRoles
		dataMap["point"] = userInfo.Point
		dataMap["cdt"] = userInfo.Cdt
		dataMap["treasure_box_total_income"] = userInfo.TreasureBoxTotalIncome
		dataMap["platform"] = userInfo.Platform
		dataMap["version"] = userInfo.Version

		dataMap["item_info"] = "{}"
		userMoney := make(map[string]map[string]float64)
		for _, moneyType := range G_BaseCfg.TokenCode {
			userMoney[moneyType] = make(map[string]float64)
			userMoney[moneyType]["amount"] = 0
			userMoney[moneyType]["amount_available"] = 0
			userMoney[moneyType]["amount_blocked"] = 0
		}
		buf, _ := json.Marshal(userMoney)
		dataMap["money"] = buf
		dataMap["suggestion"] = 0

		userdata := model.User{
			UserId:         userid,
			UserType:       userInfo.UserType,
			RoleId:         userInfo.RoleId,
			AvailableRoles: availableRoles,
			DeblockedRoles: deblockedRoles,

			NickName:    userInfo.NickName,
			Sex:         userInfo.Sex,
			Level:       userInfo.Level,
			CountryCode: userInfo.CountryCode,
			Mobile:      userInfo.Mobile,

			Email:         userInfo.Email,
			Token:         userInfo.Token,
			Status:        userInfo.Status,
			ModifyNameNum: userInfo.ModifyNameNum,
			LoginIp:       userInfo.LoginIp,

			InviteCode: userInfo.InviteCode,
			InviterId:  userInfo.InviterId,
			KycPassed:  userInfo.KycPassed,
			KycStatus:  userInfo.KycStatus,

			SysType:                msg.SysType,
			Point:                  userInfo.Point,
			Cdt:                    userInfo.Cdt,
			TreasureBoxTotalIncome: userInfo.TreasureBoxTotalIncome,
			Platform:               userInfo.Platform,
		}
		_, err := db_service.UserIns.Add(&userdata)
		if err != nil {
			logger.Errorf("Handle_Login UserIns.Add() failed(), err=", err.Error())
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}

		for _, v := range G_BaseCfg.TokenCode {
			moneyInfo := &model.GameWallet{
				UserId:          userid,
				Amount:          0,
				AmountAvailable: 0,
				AmountBlocked:   0,
				TokenCode:       v}
			_, err = db_service.GameWalletIns.Add(moneyInfo)
			if err != nil {
				logger.Errorf("Handle_Login GameWalletIns.InitAmount() failed(), err=", err.Error())
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
				return
			}
		}

		value, err := db.RedisGame.HMSet(userid, dataMap).Result()
		if err != nil {
			logger.Errorf("Handle_Login redisGame.HMSet() failed(), err=", err.Error())
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}
		if value == "OK" {
			logger.Debugf("HMSet success userId=", userid)
		}

	} else { //更新值到redis
		if userInfo.PositionX == 0 && userInfo.PositionY == 0 {
			userInfo.PositionX = G_BaseCfg.BirthPlace.X
			userInfo.PositionY = G_BaseCfg.BirthPlace.Y
		}
		wallet, err := db_service.GameWalletIns.GetAllAmount(userid)
		if err != nil {
			logger.Errorf("Handle_Login GetAllAmount() failed(), err=", err.Error())
			responseMessage.Code = errcode.ERROR_WALLET
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}

		walletAry := make([]model.GameWallet, 0)
		for _, v := range G_BaseCfg.TokenCode {
			moneyInfo := model.GameWallet{
				UserId:          userid,
				Amount:          0,
				AmountAvailable: 0,
				AmountBlocked:   0,
				TokenCode:       v}
			walletAry = append(walletAry, moneyInfo)

			bFound := false
			for _, value := range wallet {
				if value.TokenCode == v {
					bFound = true
					break
				}
			}
			if bFound {
				continue
			}

			_, err = db_service.GameWalletIns.Add(&moneyInfo)
			if err != nil {
				logger.Errorf("Handle_Login GameWalletIns.InitAmount() failed(), err=", err.Error())
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
				return
			}
		}

		userKnapsack, suberr := db_service.UserKnapsackIns.GetDataByUid(userid)
		if suberr != nil {
			logger.Errorf("Handle_Login UserKnapsackIns GetDataByUid() failed(), err=", suberr.Error())
			responseMessage.Code = errcode.ERROR_KNAPSACK
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}

		userCert, certerr := db_service.CertificationIns.GetDataByUid(userid)
		if certerr != nil {
			logger.Errorf("Handle_Login CertificationIns GetDataByUid() failed(), err=", suberr.Error())
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}
		// 机器人不用检查
		if !bRobot {
			// 检测用户信息中是否包含  圣诞老人
			userInfo.RoleId, availableRoles, deblockedRoles = CheckSantaClaus(userid, userInfo.RoleId, availableRoles, deblockedRoles)
			// 金童玉女
			userInfo.RoleId, availableRoles, deblockedRoles = CheckColdenCouple(userid, userInfo.RoleId, availableRoles, deblockedRoles)
		}

		//更新redis
		dataMap := make(map[string]interface{})
		dataMap["user_id"] = userid
		dataMap["user_type"] = userInfo.UserType
		dataMap["role_id"] = userInfo.RoleId
		dataMap["nick_name"] = userInfo.NickName
		dataMap["sex"] = userInfo.Sex

		dataMap["level"] = userInfo.Level
		dataMap["country_code"] = userInfo.CountryCode
		dataMap["mobile"] = userInfo.Mobile
		dataMap["email"] = userInfo.Email
		dataMap["token"] = token

		dataMap["status"] = userInfo.Status
		dataMap["location_id"] = userInfo.LocationId
		dataMap["position_x"] = userInfo.PositionX
		dataMap["position_y"] = userInfo.PositionY
		dataMap["house_num"] = userInfo.HouseNum

		dataMap["modify_name_num"] = userInfo.ModifyNameNum
		dataMap["login_ip"] = userInfo.LoginIp
		dataMap["invite_code"] = userInfo.InviteCode
		dataMap["inviter_id"] = userInfo.InviterId
		dataMap["kyc_passed"] = userInfo.KycPassed

		dataMap["kyc_status"] = userInfo.KycStatus
		dataMap["top_level"] = userInfo.TopLevel
		dataMap["sys_type"] = userInfo.SysType
		dataMap["point"] = userInfo.Point
		dataMap["cdt"] = strconv.FormatFloat(float64(userInfo.Cdt), 'f', 4, 32)
		dataMap["treasure_box_total_income"] = strconv.FormatFloat(float64(userInfo.TreasureBoxTotalIncome), 'f', 4, 32)
		dataMap["platform"] = userInfo.Platform
		dataMap["version"] = userInfo.Version
		dataMap["available_roles"] = availableRoles
		dataMap["deblocked_roles"] = deblockedRoles
		dataMap["update_time"] = time.Now()

		if oldUserId == "0" {
			if msg.CodeType == proto.CODE_TYPE_EMAIL {
				_, err = db_service.UpdateFields(db_service.UserTable, "email", account, dataMap)
			} else {
				_, err = db_service.UpdateFields(db_service.UserTable, "mobile", account, dataMap)
			}
		} else {
			_, err = db_service.UpdateFields(db_service.UserTable, "user_id", userid, dataMap)
		}
		if err != nil {
			logger.Errorf("UpdateFields %+v failed err=%+v", db_service.UserTable, err.Error())
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
		}

		userItemMap := make(map[int]map[int]int, 0)
		for _, v := range userKnapsack {
			if _, ok := userItemMap[v.ItemType]; !ok {
				userItemMap[v.ItemType] = make(map[int]int)
			}
			userItemMap[v.ItemType][v.ItemId] = v.ItemNum
		}
		// 圣诞老人   清除背包
		//userItemMap = CheckUserKnapsacSantaClaus(userid, userItemMap)
		userItemInfo, _ := json.Marshal(userItemMap)
		dataMap["item_info"] = userItemInfo

		userMoney := make(map[string]map[string]float64)
		for _, v := range walletAry {
			userMoney[v.TokenCode] = make(map[string]float64)
			userMoney[v.TokenCode]["amount"] = v.Amount
			userMoney[v.TokenCode]["amount_available"] = v.AmountAvailable
			userMoney[v.TokenCode]["amount_blocked"] = v.AmountBlocked
		}
		buf, _ := json.Marshal(userMoney)
		dataMap["money"] = buf

		suggestion = userCert.Suggestion
		dataMap["suggestion"] = userCert.Suggestion
		suggestion = userCert.Suggestion
		dataMap["screen_x"] = 0
		dataMap["screen_y"] = 0
		dataMap["left_top_x"] = 0
		dataMap["left_top_y"] = 0
		dataMap["right_bottom_x"] = 0
		dataMap["right_bottom_y"] = 0
		value, err := db.RedisGame.HMSet(userid, dataMap).Result()
		if err != nil {
			logger.Errorf("Handle_Login redisGame.HMSet() failed(), err=", err.Error())
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}
		if value == "OK" {
			logger.Debugf("HMSet success userId=", userid)
		}
	}

	taskKey := proto.ONLINE_KEY + userid
	utils.SetKeyValue(taskKey, "lastTime", int64(0), false, utils.ITEM_DAY)

	// 登陆用户打点
	statistical.StatisticsDotIns.Login(userid, bRobot, userInfo.Platform, userInfo.Version)

	a.auth = true //登录成功
	a.token = token

	responseMessage.SysType = userInfo.SysType
	responseMessage.UserId = userid
	responseMessage.UserType = userInfo.UserType
	responseMessage.RoleId = userInfo.RoleId
	responseMessage.NickNname = userInfo.NickName
	responseMessage.Sex = userInfo.Sex
	responseMessage.InviteCode = userInfo.InviteCode
	responseMessage.Level = userInfo.Level
	responseMessage.CountryCode = userInfo.CountryCode
	responseMessage.Mobile = userInfo.Mobile
	responseMessage.Email = userInfo.Email
	responseMessage.Status = userInfo.Status
	responseMessage.Rank = "乐于助人" //todo
	responseMessage.Token = token
	responseMessage.LocationID = userInfo.LocationId
	responseMessage.PositsionX = G_BaseCfg.BirthPlace.X //userInfo.PositionX
	responseMessage.PositsionY = G_BaseCfg.BirthPlace.Y //userInfo.PositionY
	responseMessage.HouseNum = userInfo.HouseNum
	responseMessage.ModifyNameNum = userInfo.ModifyNameNum
	responseMessage.InviterId = userInfo.InviterId
	responseMessage.KycPassed = userInfo.KycPassed
	responseMessage.KycStatus = userInfo.KycStatus
	responseMessage.AvailableRoles = availableRoles
	responseMessage.DeblockedRoles = deblockedRoles
	responseMessage.Suggestion = suggestion
	responseMessage.Point = userInfo.Point
	responseMessage.Cdt = decimal.NewFromFloat32(userInfo.Cdt)
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.ActivityStatus = GetTaskStatus()
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)

	//准备广播活动信息
	if userInfo.RoleId > 0 {
		go preBrocastTaskInfo(a, userid)
	}

}

func preBrocastTaskInfo(a *agent, userId string) {
	//延时30秒广播
	time.Sleep(time.Second * 10)
	//通知玩家更新任务信息
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_UPDATE_TASK_INFO)
	responseMessage := &proto.S2CChangeTaskStatus{}
	responseMessage.Status = GetTaskStatus()

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)

	key := "task_daily_pushtask:" + userId
	valueRet, err := db.RedisGame.HGet(key, "count").Result()
	if err != nil && err.Error() != "redis: nil" {
		logger.Errorf("TaskProcess Get(", key, ") failed:%+v", err.Error())
		return
	}

	if valueRet != "" {
		return
	}

	utils.SetKeyValue(key, "count", int64(1), true, utils.ITEM_DAY)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)

}

//用户登录
func (a *agent) HandleLogin(requestMsg *utils.Packet) {
	logger.Debugf("Handle_Login in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_LOGIN_RSP)
	responseMessage := &proto.S2CLogin{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SLogin{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	if msg.SysType == proto.ROBOT_USER {
		robotLogin(a, rsp, responseMessage, msg)
	} else {
		normalLogin(a, rsp, responseMessage, msg)
	}
}

//获取可用角色列表和已解锁角色
func GetRoles(userId string, userType, level int, availableRoles, deblockedRoles string) (string, string) {
	logger.Debugf("GetRoles in level=,availableRoles=, deblockedRoles=", level, availableRoles, deblockedRoles)
	if userType == proto.ROBOT_USER {
		availableRoles = strings.Join(G_RoleList, "|")
		return availableRoles, availableRoles
	}
	userKnap, err := db_service.UserKnapsackIns.GetDataByUid(userId)
	if err != nil || len(userKnap) == 0 {
		if availableRoles == "" && deblockedRoles == "" {
			return G_BaseCfg.PingminARoles, G_BaseCfg.PingminDRoles
		}

		return availableRoles, deblockedRoles
	}

	aRolesAry := strings.Split(availableRoles, "|")
	for _, v := range userKnap {
		roleId := GetRoleId(v.ItemId)
		if !utils.IsExistInArrs(roleId, aRolesAry) {
			aRolesAry = append(aRolesAry, roleId)
		}
	}

	availableRoles = strings.Join(aRolesAry, "|")
	if !strings.Contains(availableRoles, G_BaseCfg.PingminARoles) {
		availableRoles = G_BaseCfg.PingminARoles + availableRoles
		deblockedRoles = G_BaseCfg.PingminDRoles + deblockedRoles
	}
	return availableRoles, deblockedRoles

}

//获取用户信息
func (s *CSession) HandleGetUserInfo(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetUserInfo in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_USERINFO_RSP)
	responseMessage := &proto.S2CGetUserInfo{}
	responseMessage.Code = errcode.ERROR_SYSTEM
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetUserInfo{}
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

	userInfo, err := db.RedisGame.HGetAll(msg.UserId).Result()
	if err != nil {
		logger.Errorf("RedisGame.HGetAll error uid=, err=", msg.UserId, err.Error())
		responseMessage.Code = errcode.ERROR_REDIS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if len(userInfo) == 0 {
		userinfoRequest := make(map[string]interface{}, 0)
		userinfoRequest["userId"] = msg.UserId
		userinfoRequest["sysType"] = proto.SysType
		userinfoRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
		userinfoRequest["sign"] = utils.GenTonken(userinfoRequest)

		url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetUserinfoUrl
		userdata, err := utils.HttpGet(url, userinfoRequest)
		userResp := &proto.S2C_HTTP{}
		if err != nil {
			logger.Errorf("HttpGet err=", err.Error())
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}

		err = json.Unmarshal(userdata, userResp)
		if err != nil {
			logger.Errorf("json.Unmarshal error, err=", err.Error())
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}

		responseMessage.Level = int(userResp.Data["vipLevel"].(float64))
		responseMessage.NickName = userResp.Data["nickName"].(string)
	} else {
		level, _ := strconv.Atoi(userInfo["level"])
		responseMessage.Level = level
		roldId, _ := strconv.Atoi(userInfo["role_id"])
		responseMessage.RoleId = roldId
		responseMessage.NickName = userInfo["nick_name"]
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.UserId = msg.UserId
	responseMessage.Rank = "乐于助人" //todo, userInfo["rank"]

	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	logger.Debugf("HandleGetUserInfo end")
	return
}

//断线重连
func (a *agent) HandleRebind(requestMsg *utils.Packet) {
	logger.Infof("HandleRebind in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_REBIND_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SRebind{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	// cs := Sched.GetUser(locationId, uid)
	// if cs != nil {
	// 	logger.Debugf("session is still in, cs.UserId=", cs.UserId)
	// 	// cs.conn.Close()
	// 	a.auth = true //登录成功
	// 	cs.conn = a.conn
	// 	responseMessage.Code = errcode.MSG_SUCCESS
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// 	rsp.WriteData(responseMessage)
	// SendPacket(a.conn, rsp)
	// 	return
	// }
	platform, _ := strconv.Atoi(user_info["platform"])

	sess := NewCSession(msg.LocationID, payLoad.UserId, platform, user_info["version"], a.conn)
	a.session = sess
	ret := Sched.addSession(sess, true)
	if !ret {
		logger.Errorf("sched.addSession failed, uid=", payLoad.UserId)
		responseMessage.Code = errcode.ERROR_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	logger.Debugf("addSession()")

	statistical.StatisticsDotIns.Rebind(payLoad.UserId)

	a.auth = true //登录成功
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	SendPacket(a.conn, rsp)
	logger.Infof(string(rsp.Bytes()))
	logger.Debugf("HandleRebind end")
	return
}

//重置密码
func (a *agent) HandleResetPwd(requestMsg *utils.Packet) {
	logger.Debugf("HandleResetPwd in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_RESET_PASSWORD_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SResetPwd{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	request_map := make(map[string]interface{}, 0)
	request_map["sysType"] = msg.SysType
	request_map["password"] = msg.Password
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	url := ""
	if msg.CodeType == proto.CODE_TYPE_EMAIL {
		url = base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.EmailRetsetPwdUrl
		request_map["email"] = msg.Email
		request_map["emailCode"] = msg.VerificationCode
	} else {
		url = base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.MobileRetsetPwdUrl
		request_map["countryCode"] = msg.CountryCode
		request_map["mobile"] = msg.Mobile
		request_map["smsCode"] = msg.VerificationCode
	}
	request_map["sign"] = utils.GenTonken(request_map)
	buf, _ := json.Marshal(request_map)

	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
	}

	err = json.Unmarshal(msgdata, responseMessage)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
	logger.Debugf("HandleResetPwd end")
	return
}
