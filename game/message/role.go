package message

import (
	"encoding/json"
	kk_core "game_server/core"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/logic"
	"game_server/game/message/activity/activity_roles"
	"game_server/game/message/activity/double_year"
	"game_server/game/message/activity/santa_claus"
	"game_server/game/proto"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type RoleMessage struct {
}

var G_RoleList []string
var (
	KEY_NICK_NAME = "nickname:"
)

func init() {
	RoleData, _ := db_service.RoleInfoIns.GetAllData()
	G_RoleList = make([]string, 0)
	for _, v := range RoleData {
		G_RoleList = append(G_RoleList, strconv.Itoa(v.Id))
	}
}

func NickNameInit() {
	//角色昵称放在redis
	MsgData, _ := db_service.UserIns.GetAllData()
	for _, v := range MsgData {
		if v.NickName != "" {
			db.RedisMgr.SetCode(KEY_NICK_NAME+v.NickName, "1")
		}
	}

}

//创建角色
func (a *agent) HandleCreateRole(requestMsg *utils.Packet) {
	logger.Debugf("HandleCreateRole in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_CREATER_ROLE_RSP)

	responseMessage := &proto.S2CCreateRole{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SCreateRole{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=:", err.Error())
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

	isExists := db.RedisMgr.KeyExist(KEY_NICK_NAME + msg.NickNname)
	if isExists {
		responseMessage.Code = errcode.ERROR_ROLE_EXITS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	//更新用户表
	data_map := make(map[string]interface{}, 0)
	data_map["location_id"] = msg.LocationID
	data_map["user_type"] = msg.UserType
	data_map["role_id"] = msg.RoleId
	data_map["nick_name"] = msg.NickNname
	data_map["sex"] = msg.Sex
	data_map["update_time"] = time.Now()
	data_map["role_create_time"] = time.Now()
	db.RedisMgr.BatchHashSet(payLoad.UserId, data_map)
	db.RedisMgr.SetCode(msg.NickNname, "1")

	logger.Debugf("before PushMysql11:")
	kk_core.PushMysql(func() {
		_, err := db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, data_map)

		if err != nil {
			logger.Errorf("err=", err.Error())
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_SYSTEM
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})
			return
		}
		logger.Debugf("before SyncUserNickName()")
		//更新信息到用户系统
		err = SyncUserNickName(payLoad.UserId, msg.NickNname)
		logger.Debugf("after SyncUserNickName()")
		if err != nil {
			logger.Errorf("err:", err.Error())
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})
			return
		}

		logger.Debugf("PushWorld ()")
		kk_core.PushWorld(func() {
			responseMessage.Message = ""
			responseMessage.Code = errcode.MSG_SUCCESS

			level, _ := strconv.Atoi(user_info["level"])
			status, _ := strconv.Atoi(user_info["status"])
			x, _ := strconv.Atoi(user_info["position_x"])
			y, _ := strconv.Atoi(user_info["position_y"])
			houseNum, _ := strconv.Atoi(user_info["house_num"])
			userType, _ := strconv.Atoi(user_info["user_type"])

			responseMessage.UserId = payLoad.UserId
			responseMessage.UserType = userType
			responseMessage.RoleId = msg.RoleId
			responseMessage.Sex = msg.Sex
			responseMessage.NickNname = msg.NickNname
			responseMessage.Level = level
			responseMessage.CountryCode, _ = strconv.Atoi(user_info["country_code"])
			responseMessage.Mobile = user_info["mobile"]
			responseMessage.Email = user_info["email"]
			responseMessage.Status = status
			responseMessage.LocationID = data_map["location_id"].(int)
			responseMessage.X = x
			responseMessage.Y = y
			responseMessage.HouseNum = houseNum

			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			SendPacket(a.conn, rsp)
		})
	})
	// 检查是否是邀请用户创建角色
	go CheckIsInviteUserCreateRole(payLoad.UserId)
	return
}

//位置请求改变
func (s *CSession) HandlePositionChange(requestMsg *utils.Packet) {
	logger.Debugf("HandlePositionChange in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_BROAD_POSITION)

	responseMessage := &proto.S2CPositionChange{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = ""

	msg := &proto.C2SPositionChange{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("HandlePositionChange token error：[%+V,userid:%+V]", msg.Token, s.UserId)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	if msg.CurPosX == 0 || msg.CurPosY == 0 {
		logger.Errorf("HandlePositionChange position [x,y] is error")
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	//更新用户坐标
	data_map := make(map[string]interface{}, 0)
	data_map["location_id"] = msg.LocationID
	data_map["position_x"] = msg.CurPosX
	data_map["position_y"] = msg.CurPosY
	data_map["screen_x"] = msg.ScreenX
	data_map["screen_y"] = msg.ScreenY

	data_map["update_time"] = time.Now()
	db.RedisMgr.BatchHashSet(payLoad.UserId, data_map)
	s.X = msg.ScreenX
	s.Y = msg.ScreenY

	responseMessage.Message = ""
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.UserId = payLoad.UserId
	responseMessage.X = msg.CurPosX
	responseMessage.Y = msg.CurPosY
	responseMessage.Type = msg.Type
	responseMessage.NickName = user_info["nick_name"]
	responseMessage.Sex = user_info["sex"]
	roldId, _ := strconv.Atoi(user_info["role_id"])
	responseMessage.RoleId = roldId

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))

	//Sched.BroadCastMsg(int32(msg.LocationID), payLoad.UserId, rsp)
	Sched.SendScreenUser(int32(msg.LocationID), payLoad.UserId, msg.LeftTopX, msg.LeftTopY, msg.RightBottomX, msg.RightBottomY, rsp)
	return
}

//获取当前位置
func (s *CSession) HandleGetPosition(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetPosition in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_POSITION_RSP)

	responseMessage := &proto.S2CGetPosition{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetPosition{}
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

	responseMessage.Message = ""
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.UserId = payLoad.UserId
	x, _ := strconv.Atoi(user_info["position_x"])
	y, _ := strconv.Atoi(user_info["position_y"])

	responseMessage.X = x
	responseMessage.Y = y

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(s.conn, rsp)
	return
}

//进入城市地图
func (a *agent) HandleEnterCity(requestMsg *utils.Packet) {
	logger.Debugf("HandleEnterCity in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_ENTER_CITY_RSP)
	responseMessage := &proto.S2CEnterCity{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SEnterCity{}
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

	ok := Sched.ExistLoacateId(int32(msg.LocationID))
	if !ok {
		logger.Errorf("location_id eror:", msg.LocationID)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	platform, _ := strconv.Atoi(user_info["platform"])

	sess := NewCSession(int32(msg.LocationID), payLoad.UserId, platform, user_info["version"], a.conn)
	if G_BaseCfg.BirthScreenPlaceInit != 0 {
		sess.X = G_BaseCfg.BirthScreenPlace.X
		sess.Y = G_BaseCfg.BirthScreenPlace.Y
	}
	a.session = sess
	ret := Sched.addSession(sess, false)
	if !ret {
		logger.Errorf("sched.addSession failed, uid=", payLoad.UserId)
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	dataMap := make(map[string]interface{})
	dataMap["location_id"] = msg.LocationID

	_, err = db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, dataMap)
	if err != nil {
		logger.Errorf("err=", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	responseMessage.LocationID = msg.LocationID
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)
	logger.Debugf("HandleEnterCity end")
	go PushActivityStatus(false, payLoad.UserId)
	//进入城市之后主动推送活动状态

	//金童玉女 用户红包
	go ColdenCoupleRedEnvelope(payLoad.UserId)
	return
}

func (s *CSession) HandleGetCityUser(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetCityUser in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_CITY_USER_RSP)
	responseMessage := &proto.S2CGetCityUser{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetCityUser{}
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
	allow_treasureBox := true
	if s.Platform != 0 && s.Platform != proto.PLATFORM_ANDROID {
		allow_treasureBox = false
	}
	userList := Sched.GetSameCityUser(int32(msg.LocationID))
	responseMessage.LocationID = msg.LocationID
	areaEevnt := getAreaEventData(msg.LocationID, allow_treasureBox)
	if len(areaEevnt) > 0 {
		responseMessage.EventInfo = areaEevnt
		uesrMap := make(map[string]proto.UserPosition, 0)
		for _, v := range userList {
			redisUser := db.RedisMgr.HGetAll(v)
			if len(redisUser) > 0 {
				positionX, _ := strconv.Atoi(redisUser["position_x"])
				positionY, _ := strconv.Atoi(redisUser["position_y"])
				roldId, _ := strconv.Atoi(redisUser["role_id"])

				uesrMap[v] = proto.UserPosition{
					NickName: redisUser["nick_name"],
					Sex:      redisUser["sex"],
					RoleId:   roldId,
					X:        positionX,
					Y:        positionY,
				}
			}
		}
		responseMessage.UserList = uesrMap
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	logger.Debugf("HandleGetCityUser end")
	return
}

//退出城市地图
func (s *CSession) HandleQuitCity(requestMsg *utils.Packet) {
	logger.Debugf("HandleQuitCity in request:", requestMsg.GetBuffer())
	broadMsg := &utils.Packet{}
	broadMsg.Initialize(proto.MSG_BROAD_CITY_USER)
	broadMsgUser := &proto.S2CCityUser{}
	broadMsgUser.Code = errcode.MSG_SUCCESS
	broadMsgUser.Message = errcode.ERROR_MSG[broadMsgUser.Code]
	broadMsgUser.ActionType = proto.QUIT_CITY

	msg := &proto.C2SQuitCity{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		broadMsg.WriteData(broadMsgUser)
		SendPacket(s.conn, broadMsg)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		broadMsg.WriteData(broadMsgUser)
		SendPacket(s.conn, broadMsg)
		return
	}

	data_map := make(map[string]interface{})
	data_map["position_x"] = user_info["position_x"]
	data_map["position_y"] = user_info["position_y"]
	db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, data_map)
	var curUser proto.UserPosition
	curUser.X, _ = strconv.Atoi(user_info["position_x"])
	curUser.Y, _ = strconv.Atoi(user_info["position_y"])
	curUser.NickName = user_info["nick_name"]
	curUser.RoleId, _ = strconv.Atoi(user_info["role_id"])

	curUser.Sex = user_info["sex"]
	broadMsgUser.UserList = make(map[string]proto.UserPosition, 0)
	broadMsgUser.UserList[payLoad.UserId] = curUser
	broadMsg.WriteData(broadMsgUser)
	logger.Debugf(string(broadMsg.Bytes()))
	Sched.BroadCastMsg(int32(msg.LocationID), "", broadMsg)
	Sched.delSession(s)

	logger.Debugf("HandleQuitCity end")
	return
}

//点击事件
func (s *CSession) HandleFinishEvent(requestMsg *utils.Packet) {
	logger.Debugf("HandleFinishEvent in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_FINISH_EVENT_RSP)

	responseMessage := &proto.S2CFinishEvent{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SFinishEvent{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	if !(msg.ActivityType == proto.ACTIVITY_TYPE_NOMAL || msg.ActivityType == proto.ACTIVITY_TYPE_DOUBLE_YEAR) {
		logger.Errorf("ActivityType error, ", msg.ActivityType)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, userInfo := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	event_data := getEventData(msg.ActivityType, msg.LocationID, msg.X, msg.Y)
	if event_data == nil {
		logger.Errorf("ERROR_OBJ_NOT_EXISTS :", msg.LocationID, msg.X, msg.Y)
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		responseMessage.Code = errcode.ERROR_OBJ_NOT_EXISTS
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	get_event := func(active_type, event_id int) *proto.EventRate {
		if active_type == proto.ACTIVITY_TYPE_NOMAL {
			value, ok := G_EventConf.Load(event_id)
			if !ok {
				return nil
			}

			event := value.(proto.EventRate)
			return &event
		} else {
			value, ok := G_ChristmaEventConf.Load(event_id)
			if !ok {
				return nil
			}

			event := value.(proto.EventRate)
			return &event
		}
	}
	times := 1 // 倍数
	eventRate := get_event(msg.ActivityType, event_data.EventId)
	changeCdt := float32(0)
	changeItem := make(map[int][]*proto.AwardItem)
	if msg.ActivityType == proto.ACTIVITY_TYPE_DOUBLE_YEAR {
		bRet := false
		logger.Debugf("DoubleYearFinishEvent eventRate=%+V", eventRate)
		bRet, changeItem = G_DoubleYearEvent.DoubleYearFinishEvent(payLoad.UserId, userInfo, eventRate)
		if !bRet {
			logger.Errorf("token error", msg.Token)
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}
	} else {
		if eventRate.LimitType == proto.LIMITE_COUNT {
			key := "task_daily_rubbish:" + strconv.Itoa(event_data.EventId) + ":" + payLoad.UserId
			field := "count"
			valueRet, err := db.RedisGame.HGet(key, field).Result()
			if err == nil && valueRet != "" {
				oldValue, _ := strconv.Atoi(valueRet)
				if oldValue >= eventRate.LimitNum {
					responseMessage.Code = errcode.ERROR_CDT_OUT_OF_TODAY_LIMIT
					responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
					logger.Infof("HandleFinishEvent out of today limit user_id=", payLoad.UserId)
					rsp.WriteData(responseMessage)
					s.sendPacket(rsp)
					return
				}
			}
			utils.SetKeyValue(key, "count", int64(1), true, utils.ITEM_DAY)
		} else {
			key := "task_daily_event_num:" + payLoad.UserId
			field := "num"
			valueRet, err := db.RedisGame.HGet(key, field).Result()
			if err == nil && valueRet != "" {
				oldValue, _ := strconv.Atoi(valueRet)
				if oldValue >= eventRate.LimitNum {
					responseMessage.Code = errcode.ERROR_ITEM_OUT_OF_TODAY_LIMIT
					responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
					logger.Infof("HandleFinishEvent out of today limit user_id=", payLoad.UserId)
					rsp.WriteData(responseMessage)
					s.sendPacket(rsp)
					return
				}
			}
			utils.SetKeyValue(key, "num", int64(eventRate.ItemNum), true, utils.ITEM_DAY)
		}
		if item, ok := G_ItemList.Load(eventRate.ItemId); ok {
			node := item.(*proto.ProductItem)
			award := new(proto.AwardInfo)
			award.ItemId = eventRate.ItemId
			award.ItemNum = eventRate.ItemNum
			award.ImgUrl = node.ImgUrl
			award.ItemName = node.ItemName
			award.Desc = node.Desc

			changeCdt, changeItem, times = sendAward(proto.MSG_FINISH_EVENT, payLoad.UserId, userInfo, award)
		}
		go func() {
			//捡垃圾任务处理
			taskProcess(&s.conn, payLoad.UserId, proto.MSG_FINISH_EVENT, event_data.EventId, 1, true)
		}()
	}

	//事件重置
	resetEventData(msg.ActivityType, msg.LocationID, getPositionNo(msg.X, msg.Y), event_data)
	//返回事件处理结果给当前玩家
	responseMessage.Message = ""
	responseMessage.Code = errcode.MSG_SUCCESS
	// responseMessage.MsgData = make(map[string]interface{})
	responseMessage.UserId = payLoad.UserId
	responseMessage.LocationID = msg.LocationID
	responseMessage.X = msg.X
	responseMessage.Y = msg.Y
	// 圣诞老人 倍数乘 2
	if times == 2 {
		changeCdt = changeCdt * 2
	}
	responseMessage.Cdt = decimal.NewFromFloat32(changeCdt)
	responseMessage.ItemInfos = changeItem
	responseMessage.ActivityType = msg.ActivityType
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	//通知其它玩家消除事件
	broadMsg := &utils.Packet{}
	broadMsg.Initialize(proto.MSG_BROAD_FINISH_EVENT)
	broadMsg.WriteData(responseMessage)

	logger.Debugf(string(broadMsg.Bytes()))
	Sched.BroadCastMsg(int32(msg.LocationID), payLoad.UserId, broadMsg)

	logger.Debugf("HandleFinishEvent end")
	return
}

func Msg() {

}

//昵称检测
func (a *agent) HandleCheckNickName(requestMsg *utils.Packet) {
	logger.Debugf("HandleCheckNickName in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_CHECK_NICK_NAME_RSP)

	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SCheckNickName{}
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

	isExists := db.RedisMgr.KeyExist(KEY_NICK_NAME + msg.NickName)
	if isExists {
		responseMessage.Code = errcode.ERROR_ROLE_EXITS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(a.conn, rsp)

	logger.Debugf("HandleCheckNickName end")
	return
}

//可用角色列表
func (a *agent) HandleGetAllRole(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetAllRole in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_ALL_ROLE_RSP)
	responseMessage := &proto.S2CGetAllRole{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetAllRole{}
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

	kk_core.PushMysql(func() {
		paramData := []string{"available_roles", "deblocked_roles"}
		userInfo, err := db.RedisGame.HMGet(payLoad.UserId, paramData...).Result()
		if err != nil {
			logger.Errorf("HandleGetAllRole HMGet(userid, 'paramData') failed err=", err.Error())
			responseMessage.Code = errcode.ERROR_REDIS
			responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_REDIS]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}

		availableRoles := strings.Split(userInfo[0].(string), "|")
		deblockedRoles := strings.Split(userInfo[1].(string), "|")
		roleInfos, err := db_service.RoleInfoIns.GetAllData()
		if err != nil {
			logger.Errorf("GetAllData() error, err=", err.Error())
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})
		}
		responseMessage.RoleInfos = make([]*proto.RoleInfo, 0)
		for _, role := range roleInfos {
			roleInfo := &proto.RoleInfo{
				RoleId:   role.Id,
				RoleName: role.RoleName,
				Sex:      role.Sex,
				State:    0,
			}
			if utils.IsExistInArrs(strconv.Itoa(role.Id), availableRoles) {
				roleInfo.State = 1
			}
			if utils.IsExistInArrs(strconv.Itoa(role.Id), deblockedRoles) {
				roleInfo.State = 2
			}
			// 判断是否是金童玉女
			goldenCouple := activity_roles.NewGoldenCouple()
			if goldenCouple.CheckIsGoldenCoupleRole(role.Id) {
				roleInfo.ExpireTime = goldenCouple.GetRoleExpireTime(payLoad.UserId, role.Id)
			}
			//判断是否是圣诞老人的角色
			if role.Id == base.Setting.Doubleyear.SantaClausRoleId {
				//获取过期时间
				roleInfo.ExpireTime = santa_claus.SantaClausLogic.GetRoleExpireTime(payLoad.UserId) // 过期时间
			}
			roleInfo.ItemId = GetItemId(role.Id)
			responseMessage.RoleInfos = append(responseMessage.RoleInfos, roleInfo)
			// logger.Debugf("roleInfo=", *roleInfo)
		}

		kk_core.PushWorld(func() {
			responseMessage.Code = errcode.MSG_SUCCESS
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

			rsp.WriteData(responseMessage)
			logger.Debugf("end, res=", string(rsp.Bytes()))
			SendPacket(a.conn, rsp)
		})

	})

}

//解锁角色
func (a *agent) HandleUserAddRole(requestMsg *utils.Packet) {
	logger.Debugf("HandleUserAddRole in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_USER_ADD_ROLE_RSP)
	responseMessage := &proto.S2CAddRole{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SAddRole{}
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

	// 检查需要解锁的卡片是否是可以正常使用
	baeRoles := activity_roles.NewActivityRoles()
	if baeRoles.CheckIsActivityRoleCard(msg.ItemId) {
		if baeRoles.CheckRoleCardIsUnLock(msg.ItemId) == false {
			// 卡片已过期
			responseMessage.Code = errcode.ERROR_ITMES_EXPIRED
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
			return
		}
	}

	deblocked_roles := strings.Split(user_info["deblocked_roles"], "|")
	if utils.IsExistInArrs(strconv.Itoa(msg.RoleId), deblocked_roles) {
		logger.Errorf("角色已解锁:%+v,%+v", payLoad.UserId, msg.RoleId)
		responseMessage.Code = errcode.ERROR_ROLE_IS_UNLOCK
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	userItem := user_info["item_info"]
	itemMap := make(map[int]map[int]int, 0)
	err = json.Unmarshal([]byte(userItem), &itemMap)
	if err != nil {
		responseMessage.Code = errcode.ERROR_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}

	kk_core.PushMysql(func() {
		itemInfo, err := db_service.ItemIns.GetDataById(msg.ItemId)
		if err != nil {
			logger.Errorf("ItemIns.GetDataById() failed(), err=", err.Error())
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})
			return
		}

		bFlag := false //用户是否有解锁卡
		if _, ok := itemMap[itemInfo.ItemType]; ok {
			if _, subok := itemMap[itemInfo.ItemType][itemInfo.Id]; subok {
				if itemMap[itemInfo.ItemType][itemInfo.Id] > 0 {
					bFlag = true
				}
			}
		}

		if !bFlag {
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_NO_CLEAR_CARD
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})
			return
		}

		userInfoMap := make(map[string]interface{}, 0)
		userInfoMap["update_time"] = time.Now()
		bDelete := false
		//使用解锁卡
		if itemMap[itemInfo.ItemType][itemInfo.Id] == 1 {
			delete(itemMap[itemInfo.ItemType], itemInfo.Id)
			bDelete = true
		} else {
			itemMap[itemInfo.ItemType][itemInfo.Id]--
		}
		userItemInfo, _ := json.Marshal(itemMap)

		//更新redis值
		userInfoMap["item_info"] = userItemInfo
		//金童玉女
		goldenCouple := activity_roles.NewGoldenCouple()
		if goldenCouple.CheckIsGoldenCoupleRole(msg.RoleId) {
			userInfoMap["deblocked_roles"] = user_info["deblocked_roles"] + goldenCouple.GetUnloc()
		} else {
			userInfoMap["deblocked_roles"] = user_info["deblocked_roles"] + "|" + strconv.Itoa(msg.RoleId)
		}
		_, err = db.RedisGame.HMSet(payLoad.UserId, userInfoMap).Result()
		if err != nil {
			logger.Errorf("RedisGame.HMSet error, err=", err.Error())
			responseMessage.Code = errcode.ERROR_SYSTEM
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)

			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_SYSTEM
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})
			return
		}

		//更新用户记录
		dataMap := make(map[string]interface{}, 0)
		if goldenCouple.CheckIsGoldenCoupleRole(msg.RoleId) {
			dataMap["deblocked_roles"] = user_info["deblocked_roles"] + goldenCouple.GetUnloc()
		} else {
			dataMap["deblocked_roles"] = user_info["deblocked_roles"] + "|" + strconv.Itoa(msg.RoleId)
		}
		dataMap["update_time"] = time.Now()

		_, err = db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, dataMap)
		if err != nil {
			logger.Errorf("db_service.UpdateFields() err=", err.Error())
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				SendPacket(a.conn, rsp)
			})

			//回滚
			if bDelete {
				itemMap[itemInfo.ItemType][itemInfo.Id] = 1
			} else {
				itemMap[itemInfo.ItemType][itemInfo.Id]++
			}
			userItemInfo, _ := json.Marshal(itemMap)
			userInfoMap["item_info"] = userItemInfo
			db.RedisGame.HMSet(payLoad.UserId, userInfoMap).Result()
			return
		}

		if bDelete {
			db_service.UserKnapsackIns.Delete(payLoad.UserId, itemInfo.Id)
		} else {
			data_map := make(map[string]interface{}, 0)
			data_map["update_time"] = time.Now()
			data_map["item_num"] = itemMap[itemInfo.ItemType][itemInfo.Id]
			db_service.UserKnapsackIns.UpdateData(payLoad.UserId, itemInfo.Id, data_map)
		}
		//// 判断角色卡是否圣诞老人角色卡
		//if msg.ItemId == base.Setting.Doubleyear.SantaClausCardId {
		//	//记录圣诞老人过期时间
		//	santa_claus.SantaClausLogic.AddRolseExpired(payLoad.UserId, msg.ItemId)
		//	responseMessage.ExpireTime = santa_claus.SantaClausLogic.GetRoleExpireTime(payLoad.UserId) // 过期时间
		//}
		// 添加金童玉女过期时间
		isTrue, _ := goldenCouple.AddGoldenCoupleRoleExpire(payLoad.UserId, msg.RoleId)
		if isTrue {
			responseMessage.ExpireTime = goldenCouple.GetRoleExpireTime(payLoad.UserId, msg.RoleId)
		}
		kk_core.PushWorld(func() {
			logger.Debugf("HandleUserAddRole 333, MSG_SUCCESS")
			responseMessage.Code = errcode.MSG_SUCCESS
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			responseMessage.RoleId = msg.RoleId
			rsp.WriteData(responseMessage)
			SendPacket(a.conn, rsp)
		})
	})

}

//角色选择
func (s *CSession) HandleRoleSelect(requestMsg *utils.Packet) {
	logger.Debugf("HandleRoleSelect in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SELECT_ROLE_RSP)
	responseMessage := &proto.S2CCommon{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SSelectRole{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	// 判断要切换的角色是否是活动角色
	baseRoles := activity_roles.NewActivityRoles()
	if baseRoles.CheckIsActivityRole(msg.RoleId) {
		// 判断活动角色是否过期
		isExpire, err := baseRoles.CheckRolesIsExpire(msg.UserId, msg.RoleId)
		if isExpire == false {
			logger.Errorf("roles is Expire roleid= [%v] userId:[%v], err=[%v]", msg.RoleId, msg.UserId, err)
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_ROLES_EXPIRED
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				s.sendPacket(rsp)
			})
			return
		}
		// 记录用户之前的角色ID
		roleId, _ := db.RedisGame.HGet(msg.UserId, "role_id").Result()
		UserRoleId, _ := strconv.Atoi(roleId)
		if baseRoles.CheckIsActivityRole(UserRoleId) == false {
			// 保存旧的ID
			goldenCouple := activity_roles.NewGoldenCouple()
			goldenCouple.SaveUserOldRole(msg.UserId, UserRoleId)
		}
	}

	//更新redis
	_, err = db.RedisGame.HSet(msg.UserId, "role_id", msg.RoleId).Result()
	if err != nil {
		logger.Errorf("db.RedisGame.HSet() err=", err.Error())
		kk_core.PushWorld(func() {
			responseMessage.Code = errcode.ERROR_REDIS
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
		})
		return
	}

	//更新mysql
	kk_core.PushMysql(func() {
		dataMap := make(map[string]interface{}, 0)
		dataMap["role_id"] = msg.RoleId
		dataMap["sex"] = msg.Sex
		dataMap["update_time"] = time.Now()
		_, err = db_service.UpdateFields(db_service.UserTable, "user_id", msg.UserId, dataMap)
		if err != nil {
			logger.Errorf("db_service.UpdateFields() err=", err.Error())
			kk_core.PushWorld(func() {
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				s.sendPacket(rsp)
			})
			return
		}

		kk_core.PushWorld(func() {
			responseMessage.Code = errcode.MSG_SUCCESS
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)

			broad := &utils.Packet{}
			broad.Initialize(proto.MSG_ROLE_CHANGE)
			broadBroad := &proto.S2CRoleChange{}
			broadBroad.UserID = msg.UserId
			broadBroad.RoleID = msg.RoleId
			broad.WriteData(broadBroad)
			Sched.BroadCastMsg(s.LocateId, msg.UserId, broad)
		})

	})
	if baseRoles.CheckIsActivityRole(msg.RoleId) {
		//金童玉女 用户红包
		go ColdenCoupleRedEnvelope(msg.UserId)
	}
	return
}

//获取好友详细信息
func (s *CSession) HandleGetFriendInfo(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetFriendInfo in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_FRIEND_INFO_RSP)
	responseMessage := &proto.S2CHTTP{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetFirendInfo{}
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

	userinfoRequest := make(map[string]interface{}, 0)
	userinfoRequest["userId"] = msg.UserId
	userinfoRequest["sysType"] = msg.SysType
	userinfoRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	userinfoRequest["sign"] = utils.GenTonken(userinfoRequest)

	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetUserinfoUrl
	userdata, err := utils.HttpGet(url, userinfoRequest)
	userResp := &proto.S2C_HTTP{}
	if err == nil {
		err = json.Unmarshal(userdata, userResp)
		if err == nil {
			rsp.WriteData(userResp)
			logger.Errorf(string(rsp.Bytes()))
			SendPacket(s.conn, rsp)
			return
		}
	}
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(s.conn, rsp)

	logger.Debugf("HandleGetFriendInfo end")
	return
}

//获取用户的邀请关系
func (s *CSession) HandleGetInvitation(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetInvitation in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_INVITATION_RSP)
	responseMessage := &proto.S2CHTTP{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetUserInvitation{}
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

	userinfoRequest := make(map[string]interface{}, 0)
	userinfoRequest["userId"] = msg.UserId
	userinfoRequest["sysType"] = msg.SysType
	userinfoRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	userinfoRequest["sign"] = utils.GenTonken(userinfoRequest)

	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.UserInvitationUrl
	userdata, err := utils.HttpGet(url, userinfoRequest)
	userResp := &proto.S2C_HTTP{}
	if err == nil {
		err = json.Unmarshal(userdata, userResp)
		if err == nil {
			rsp.WriteData(userResp)
			logger.Debugf(string(rsp.Bytes()))
			SendPacket(s.conn, rsp)
			return
		}
	}
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(s.conn, rsp)

	logger.Debugf("HandleGetInvitation end")
}

// 检查是否是邀请用户创建角色
func CheckIsInviteUserCreateRole(userId string) {
	userInvite := logic.NewUserInvitation()
	isCreateRole, inviteUserId, _ := userInvite.CheckUserInviteIsFirstCreateRole(userId)
	//给用户的邀请人增加积分
	if isCreateRole && inviteUserId != "" {
		rankList := double_year.NewRankList()
		rankList.UpdateProp(inviteUserId, double_year.PropInviteFriend, 1) // 添加积分
	}
}
