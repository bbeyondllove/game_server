package message

import (
	"encoding/json"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/proto"
	"strconv"

	"game_server/core/logger"
)

type OffLine struct {
	msgQueue *utils.SyncQueue
}

var G_OffLine = NewOffLine()

func NewOffLine() *OffLine {
	return &OffLine{msgQueue: utils.NewSyncQueue()}
}

//读取所有消息并处理
func (this *OffLine) HandleMsg() {
	pcks, ok := this.msgQueue.TryPopAll()
	if !ok || pcks == nil {
		return
	}

	for _, pck := range pcks {
		msg, ok := pck.(*utils.Packet)
		if !ok {
			logger.Errorf("pck.(*proto.C2SMessage) not ok")
			break
		}

		this.handleOffLine(msg)
	}
	logger.Debugf("OffLine HandleMsg() end")
	return
}

//处理单个离线消息
func (this *OffLine) handleOffLine(requestMsg *utils.Packet) {
	logger.Debugf("handleOffLine in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_BROAD_USER_OFFLINE)

	responseMessage := &proto.S2CCityUser{}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""

	msg := &proto.CityUser{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		return
	}

	if len(msg.UserList) == 0 {
		return
	}
	for k, v := range msg.UserList {
		uesrMap := make(map[string]proto.UserPosition, 0)
		for _, value := range v {
			redisUser := db.RedisMgr.HGetAll(value)
			if len(redisUser) == 0 {
				logger.Errorf("hanleOffLine user not found:", value)
				continue
			}
			x, _ := strconv.Atoi(redisUser["position_x"])
			y, _ := strconv.Atoi(redisUser["position_y"])

			uesrMap[value] = proto.UserPosition{
				NickName: redisUser["nick_name"],
				Sex:      redisUser["sex"],
				X:        x,
				Y:        y,
			}

			data_map := make(map[string]interface{})
			data_map["position_x"] = x
			data_map["position_y"] = y
			db_service.UpdateFields(db_service.UserTable, "user_id", value, data_map)
		}

		if len(uesrMap) == 0 {
			continue
		}

		responseMessage.UserList = uesrMap

		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		Sched.BroadCastMsg(k, "", rsp)
	}

	logger.Debugf("handleOffLine end")
	return
}

//QueuePacket 消息入队
func (this *OffLine) QueuePacket(msg *utils.Packet) {
	this.msgQueue.Push(msg)
}
